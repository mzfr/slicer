package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func exported(activity *etree.Element) {
	// Choose all those activity which are exported
	// If not exported then we check if the intent-filters are set
	exported := activity.SelectAttrValue("android:exported", "none")

	if exported == "none" {
		// If the activity doesn't have android:exported defined
		// we check if it has any Intent-filerts, if yes that means
		// exported by default
		if intentFilter := activity.SelectElements("intent-filter"); intentFilter != nil {
			fmt.Println("Activity Name:", activity.SelectAttrValue("android:name", "name not defined"))
			// TODO: Null has to be printed in red.
			fmt.Println("Permission:", activity.SelectAttrValue("android:permission", "null"))
			getIntents(intentFilter)
		}
	} else if exported == "true" {
		if intentFilter := activity.SelectElements("intent-filter"); intentFilter != nil {
			fmt.Println("Activity Name:", activity.SelectAttrValue("android:name", "name not defined"))
			// TODO: Null has to be printed in red.
			fmt.Println("Permission:", activity.SelectAttrValue("android:permission", "null"))

			getIntents(intentFilter)
		}
	}
}

func parseManifest(document *etree.Document) {
	root := document.SelectElement("manifest")
	for _, app := range root.SelectElements("application") {
		// Check if the backup is allowed or not
		backup := app.SelectAttrValue("android:allowBackup", "true")
		fmt.Println("Backup allowed? ", backup)

		//Check if the app is set debuggable
		debuggable := app.SelectAttrValue("android:debuggable", "false")
		fmt.Println("Debuggable? ", debuggable)
		fmt.Println()
		// Select all the activities
		for _, activity := range app.SelectElements("activity") {
			exported(activity)
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

	for _, key := range paths.([]interface{}) {
		for _, file := range key.(map[interface{}]interface{}) {
			filePath := fmt.Sprintf("%s/%s", dir, file)
			doc := etree.NewDocument()
			if err := doc.ReadFromFile(filePath); err != nil {
				panic(err)
			}
			if err := doc.SelectElement("manifest"); err != nil {
				parseManifest(doc)
			}
		}
	}
}
