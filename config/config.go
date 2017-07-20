package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var conf = AppConfigYAML{}

func InitConfig() {
	loadAppConfig(getPathFromArgs(), &conf)
}

func loadAppConfig(confPath string, conf *AppConfigYAML) {
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic("Not found config.yaml")
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic("Cannot unmashal config")
	}
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

func GlobalConfig() *AppConfigYAML {
	return &conf
}
