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
