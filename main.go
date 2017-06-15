package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var conf appconfig

type appconfig struct {
	Port    string `yaml:"port"`
	Storage struct {
		Path string `yaml:"path"`
	}
	Image struct {
		MaxWidth  int `yaml:"maxwidth"`
		MaxHeight int `yaml:"maxheight"`
		Quality   int `yaml:"quality"`
	}
	Thumbnail struct {
		MaxWidth  int `yaml:"maxwidth"`
		MaxHeight int `yaml:"maxheight"`
		Quality   int `yaml:"quality"`
	}
}

func loadAppConfig(confPath string, conf *appconfig) *appconfig {
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
	fmt.Println("image-converter started on port ", conf.Port)
	http.HandleFunc("/upload", uploadImageHandler)
	http.ListenAndServe(conf.Port, nil)
}
