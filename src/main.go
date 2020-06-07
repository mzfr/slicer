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
	patterns := v.Get("patterns")
	fmt.Println(patterns)

	// firebase := v.GetBool("checks.firebase")
	// gmap := v.GetBool("checks.gmaps")

	// fmt.Println(firebase, gmap)

	if dir != "" {
		fmt.Printf("Searching path: %s", dir)

		var files []string
		err := filepath.Walk(dir, visit(&files))
		if err != nil {
			panic(err)
		}
		// fmt.Println(files)

		var cmd *exec.Cmd

		for _, key := range patterns.(map[string]interface{}) {
			k := key.([]interface{})
			command := fmt.Sprintf("%v %v", k[1], k[0])
			// path := filepath.Join(dir, "*.*")
			for _, f := range files {
				// fmt.Println(f)
				cmd = exec.Command("grep", command, f)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				// cmd.Stderr = os.Stderr
				cmd.Run()
			}

		}

	}
}
