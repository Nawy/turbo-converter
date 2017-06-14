package main

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// func getAvgColor(img Image) {
//
// }

type appconfig struct {
	Port     int32  `yaml:"port"`
	SavePath string `yaml:"storage_path"`
}

func getAppConfig(conf *appconfig) *appconfig {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic("Not found config.yaml")
	}

	err = yaml.Unmarshal(yamlFile, *conf)
	if err != nil {
		panic("")
	}
	return conf
}

func main() {
	var conf appconfig
	conf = getAppConfig(&conf)

	http.HandleFunc("/upload/", uploadImageHandler)
	http.ListenAndServe(addr, handler)
}
