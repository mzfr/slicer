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
