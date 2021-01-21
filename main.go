package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	conf "github.com/mzfr/slicer/config"
	"github.com/mzfr/slicer/extractor"

	"github.com/beevik/etree"
	"github.com/spf13/viper"
)

var (
	dir    string
	banner bool
)

// NotVulnerable exported
var NotVulnerable = map[string]bool{
	"net.openid.appauth.RedirectUriReceiverActivity": true,
}

func init() {
	flag.Usage = func() {
		h := []string{
			"Extract information from Manifest and strings of an APK\n",
			"Usage:",
			"\tslicer [OPTION] [Extracted APK directory]\n",
			"Options:",
			"\n  -d, --dir		path to jadx output directory",
			"  -o, --output		Name of the output file(not implemented)",
			"\nExamples:\n",
			"slicer -d /path/to/the/extract/apk",
		}

		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

// Function to test if the given
// directory/file exists or not
func dirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// list all the files of a given directory
func visit(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	return files, err
}

func sendRequest(url string) (*http.Response, error) {
	req, err := http.Get(url)
	return req, err
}

// ConfigReader reads the configuration using viper
func ConfigReader() (*viper.Viper, error) {
	v := viper.New()
	// Set the file name of the configurations file
	v.SetConfigName("config")

	// Set the path to look for the configurations file
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.slicer/")
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
		fmt.Println("\tIntent-filters:")
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
	activityName := component.SelectAttrValue("android:name", "name not defined")
	// If the activity is present in unhackable
	// kind of list then no point in reporting it
	// see issue #25 on github.com/mzfr/slicer
	if NotVulnerable[activityName] {
		return
	}
	permission := component.SelectAttrValue("android:permission", "null")
	acitvityCode := strings.ReplaceAll(activityName, ".", "/")

	// If the permissions is not "null" than there is no use because we can't access that from outside.
	if permission != "null" {
		return
	}

	if exported == "none" {
		// If the activity doesn't have android:exported defined
		// we check if it has any Intent-filerts, if yes that means
		// exported by default
		if intentFilter := component.SelectElements("intent-filter"); intentFilter != nil {
			fmt.Printf("\t%s:", activityName)
			fmt.Printf("\n\tFile-Path: %s/sources/%s.java", dir, acitvityCode)
			getIntents(intentFilter)
		}
	} else if exported == "true" {
		fmt.Printf("\t%s:", activityName)
		fmt.Printf("\n\tFile-Path: %s/sources/%s.java", dir, acitvityCode)
		if intentFilter := component.SelectElements("intent-filter"); intentFilter != nil {
			getIntents(intentFilter)
		}
	}
}

// Parse the AndroidManifest.xml file
func parseManifest(document *etree.Document) {
	root := document.SelectElement("manifest")

	// Show the name of the package
	packageName := root.SelectAttrValue("package", "none")
	fmt.Println("Package: ", packageName)

	// Keep the version of the app as well
	appVersion := root.SelectAttrValue("android:versionName", "none")
	fmt.Println("Version: ", appVersion)

	for _, app := range root.SelectElements("application") {
		// Check if the backup is allowed or not
		backup := app.SelectAttrValue("android:allowBackup", "false")
		fmt.Println("Backup allowed: ", backup)

		//Check if the app is set debuggable
		debuggable := app.SelectAttrValue("android:debuggable", "false")
		fmt.Println("Debuggable: ", debuggable)

		var attackSurface = []string{"activity", "receiver", "service", "provider"}
		for _, com := range attackSurface {
			fmt.Printf("\n%s:\n", com)
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
		// Get Firebase DB URL and check if /.json trick works or not
		if strValues == "firebase_database_url" {
			firebaseURL := fmt.Sprintf("%s/.json", str.Text())
			req, err := sendRequest(firebaseURL)
			if err != nil {
				fmt.Println("Couldn't connect to Firebase")
				continue
			}
			if req.StatusCode == 401 {
				fmt.Printf("\n\t- %s: Permission Denied", firebaseURL)
				fmt.Println()
			} else {
				fmt.Printf("\n\t- %s: Is open to public", firebaseURL)
				fmt.Println()
			}
		}

		// Get Google API keys and see if they are restricted or not
		if strValues == "google_api_key" || strValues == "google_map_keys" {
			for _, keys := range googleURL.([]interface{}) {
				for _, values := range keys.(map[interface{}]interface{}) {
					requestURL := fmt.Sprintf("%s%s", values, str.Text())

					req, err := sendRequest(requestURL)
					if err != nil {
						fmt.Println("Unable to connect to the the google map")
						continue
					}
					body, err := ioutil.ReadAll(req.Body)
					if err != nil {
						return
					}

					if req.StatusCode != 403 && !strings.Contains(string(body), "API project is not authorized") && str.Text() != "" {
						fmt.Printf("\n\t- %s: %d\n", requestURL, req.StatusCode)
					}
				}
			}
		}

		// Some keys that I have found in loads of strings.xml and they have nothing important
		// so just filter those out.
		if strings.Contains(strings.ToLower(strValues), "api") && str.Text() != "" && strings.Contains(strings.ToLower(strValues), "key") && strings.Contains(strings.ToLower(strValues), "tokens") {
			if strValues == "abc_capital_off" || strValues == "abc_capital_on" || strValues == "currentApiLevel" {
				continue
			}
			fmt.Printf("\t- %s: %s\n", strValues, str.Text())
		}
	}
}

func main() {
	flag.StringVar(&dir, "d", "", "")

	v, _ := ConfigReader()

	paths := v.Get("paths")
	googleURL := v.Get("URLs")

	for _, key := range paths.([]interface{}) {
		for _, file := range key.(map[interface{}]interface{}) {
			PathofFile := fmt.Sprintf("%s/%s", dir, file)
			checkDir, err := dirExist(PathofFile)
			// in some app directories like raw or xml doesn't exists
			if err != nil {
				fmt.Printf("Some issue with directories")
				continue
			}
			if checkDir == false {
				continue
			}
			fi, err := os.Stat(PathofFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Reading the config.yml if the input is of a file
			// Than use the case 1 i.e mode.IsRegular
			// if the input was a directory than the case 2
			switch mode := fi.Mode(); {
			case mode.IsRegular():
				doc := etree.NewDocument()
				if err := doc.ReadFromFile(PathofFile); err != nil {
					panic(err)
				}
				// Parsing AndroidManifest for any API keys in there
				if err := doc.SelectElement("manifest"); err != nil {
					parseManifest(doc)
					fmt.Printf("%s:\n", "Apikeys-in-manifest")
					for _, elem := range doc.FindElements("./manifest/application/meta-data") {
						name := elem.SelectAttrValue("android:name", "")
						if name != "" && strings.Contains(strings.ToLower(name), "api") {
							values := elem.SelectAttrValue("android:value", "none")
							fmt.Printf("\t- %s: %s\n", name, values)
						}
					}
				} else {
					fmt.Printf("%s:\n", "Strings")
					parseStrings(doc, googleURL)
				}
			// Just listing files of raw and xml directory
			case mode.IsDir():
				if filepath.Base(PathofFile) == "xml" {
					fmt.Printf("%s:\n", "XML-files")
				} else {
					fmt.Printf("%s:\n", "raw-files")
				}
				files, err := visit(PathofFile)
				if err != nil {
					fmt.Println(err)
				}
				for _, file := range files {
					fmt.Printf("\t- %s\n", file)
				}
			}
		}
	}
	// This should extract all the URLs
	extractor.Extract(dir)
}
