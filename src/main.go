package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	conf "./config"
	"github.com/beevik/etree"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/viper"
)

func init() {
	flag.Usage = func() {
		h := []string{
			"Extract information from the APK file",
			"",
			"Options:",
			"  -d, --dir <path>       path to jadx output directory",
			"  -o, --output <filename>  Name of the output file",
			"",
		}

		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

// ConfigReader reads the configuration using viper
func ConfigReader() (*viper.Viper, error) {
	v := viper.New()
	// Set the file name of the configurations file
	v.SetConfigName("config")

	// Set the path to look for the configurations file
	v.AddConfigPath(".")
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var configuration conf.Configurations
	err := v.Unmarshal(&configuration)

	return v, err
}

//Get Intents of all those activites which are either
// exported or have some intent filters defined
func getIntents(intentFilters []*etree.Element) {
	var formatedIntent string
	for _, intents := range intentFilters {
		fmt.Println("\tIntent filters")
		for _, Type := range intents.ChildElements() {
			if Type.Tag == "data" {
				host := Type.SelectAttrValue("android:host", "*")
				scheme := Type.SelectAttrValue("android:scheme", "*")
				formatedIntent = fmt.Sprintf("\t\t - %s: %s://%s", Type.Tag, scheme, host)

			} else {
				formatedIntent = fmt.Sprintf("\t\t - %s: %s", Type.Tag, Type.SelectAttrValue("android:name", "no name"))
			}
			fmt.Println(formatedIntent)
		}
	}
}

// Check all the activities, receivers, broadcasts, if they are exported.
// If not exported then we check if the intent-filters are set
func exported(component *etree.Element) {
	exported := component.SelectAttrValue("android:exported", "none")

	if exported == "none" {
		// If the activity doesn't have android:exported defined
		// we check if it has any Intent-filerts, if yes that means
		// exported by default
		if intentFilter := component.SelectElements("intent-filter"); intentFilter != nil {
			fmt.Printf("\n- %s", component.SelectAttrValue("android:name", "name not defined"))
			// TODO: Null has to be printed in red.
			fmt.Println("\n\tPermission:", component.SelectAttrValue("android:permission", "null"))
			getIntents(intentFilter)
		}
	} else if exported == "true" {
		if intentFilter := component.SelectElements("intent-filter"); intentFilter != nil {
			fmt.Printf("\n- %s", component.SelectAttrValue("android:name", "name not defined"))
			// TODO: Null has to be printed in red.
			fmt.Println("\n\tPermission:", component.SelectAttrValue("android:permission", "null"))

			getIntents(intentFilter)
		}
	}
}

// Parse the AndroidManifest.xml file
func parseManifest(document *etree.Document) {
	root := document.SelectElement("manifest")
	for _, app := range root.SelectElements("application") {
		// Check if the backup is allowed or not
		backup := app.SelectAttrValue("android:allowBackup", "true")
		fmt.Println("Backup allowed? ", backup)

		//Check if the app is set debuggable
		debuggable := app.SelectAttrValue("android:debuggable", "false")
		fmt.Println("Debuggable? ", debuggable)

		var attackSurface = []string{"activity", "receiver", "service"}
		for _, com := range attackSurface {
			fmt.Printf("\n------%s------\n", com)
			for _, components := range app.SelectElements(com) {
				exported(components)
			}
		}
	}
}

// Parse the /res/values/strings.xml
func parseStrings(document *etree.Document, googleURL interface{}) {
	root := document.SelectElement("resources")

	for _, str := range root.SelectElements("string") {
		strValues := str.SelectAttrValue("name", "none")
		if strValues == "firebase_database_url" {
			firebaseURL := fmt.Sprintf("%s/.json", str.Text())
			req, err := http.Get(firebaseURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to connect with firebase DB: %s\n", err)
			}

			defer req.Body.Close()

			if req.StatusCode == 401 {
				fmt.Printf("%s is not publiclly accesible", firebaseURL)
				fmt.Println()
			}
		}
		if strValues == "google_api_key" || strValues == "google_map_keys" {
			for _, keys := range googleURL.([]interface{}) {
				for _, values := range keys.(map[interface{}]interface{}) {
					requestURL := fmt.Sprintf("%s%s", values, str.Text())

					req, err := http.Get(requestURL)
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to connect with firebase DB: %s\n", err)
					}

					defer req.Body.Close()
					body, err := ioutil.ReadAll(req.Body)
					if err != nil {
						return
					}

					if req.StatusCode != 403 && !strings.Contains(string(body), "API project is not authorized") {
						fmt.Printf("\n%s returns %d", requestURL, req.StatusCode)
					}
				}
			}
		}
		if strings.Contains(strings.ToLower(strValues), "api") {
			if strValues == "abc_capital_off" || strValues == "abc_capital_on" || strValues == "currentApiLevel" {
				continue
			}
			fmt.Printf("%s - %s\n", strValues, str.Text())
		}
	}
}

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "", "")
	flag.StringVar(&dir, "d", "", "")

	flag.Parse()

	myFigure := figure.NewColorFigure("Slicer", "big", "green", true)
	myFigure.Print()

	v, err := ConfigReader()
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	paths := v.Get("paths")
	gURL := v.Get("URLs")
	// fmt.Println(reflect.TypeOf(gURL))

	for _, key := range paths.([]interface{}) {
		for _, file := range key.(map[interface{}]interface{}) {
			filePath := fmt.Sprintf("%s/%s", dir, file)
			doc := etree.NewDocument()
			if err := doc.ReadFromFile(filePath); err != nil {
				panic(err)
			}
			if err := doc.SelectElement("manifest"); err != nil {
				parseManifest(doc)
			} else {
				fmt.Printf("\n------%s------\n", "Strings")
				parseStrings(doc, gURL)
			}
		}
	}
}
