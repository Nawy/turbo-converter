package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Type is type of objects
type Type int

const (
	// IMAGE is full size image
	IMAGE = Type(1)
	// THUMBNAIL is small image
	THUMBNAIL = Type(2)
)

var conf AppConfigYAML

type AppConfigYAML struct {
	Port  string `yaml:"port"`
	Salt  string `yaml:"salt"`
	Image struct {
		StoragePath    string `yaml:"storage_path"`
		ResponsePath   string `yaml:"response_path"`
		MaxWidth       int    `yaml:"maxwidth"`
		MaxHeight      int    `yaml:"maxheight"`
		Quality        int    `yaml:"quality"`
		PostProcessing struct {
			Sharpen    float64 `yaml:"sharpen"`
			Brightness float64 `yaml:"brightness"`
			Contrast   float64 `yaml:"contrast"`
			Gamma      float64 `yaml:"gamma"`
		} `yaml:"post_processing"`
		Downscale string `yaml:"downscale"`
	}
	Thumbnail struct {
		StoragePath    string `yaml:"storage_path"`
		ResponsePath   string `yaml:"response_path"`
		MaxWidth       int    `yaml:"maxwidth"`
		MaxHeight      int    `yaml:"maxheight"`
		Quality        int    `yaml:"quality"`
		PostProcessing struct {
			Sharpen    float64 `yaml:"sharpen"`
			Brightness float64 `yaml:"brightness"`
			Contrast   float64 `yaml:"contrast"`
			Gamma      float64 `yaml:"gamma"`
		} `yaml:"post_processing"`
		Resize struct {
			Upscale   string `yaml:"upscale"`
			Downscale string `yaml:"downscale"`
		} `yaml:"resize"`
	}
	Logging struct {
		Path   string `yaml:"path"`
		Format string `yaml:"format"`
	}
}

func loadAppConfig(confPath string, conf *AppConfigYAML) *AppConfigYAML {
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic("Not found config.yaml")
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic("Cannot unmashal config")
	}

	return conf
}

func getPathFromArgs() string {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "-c" {
			return args[i+1]
		}
	}
	panic("cannot read config from -c option")
}

func main() {
	conf = *loadAppConfig(getPathFromArgs(), &conf)
	initHashGen()
	initLogging()

	http.HandleFunc("/upload", uploadImageHandler)
	http.HandleFunc("/delete", deleteImageHandler)
	http.HandleFunc("/status", statusHandler)

	fmt.Println("image-converter started on port ", conf.Port)
	http.ListenAndServe(conf.Port, nil)
}
