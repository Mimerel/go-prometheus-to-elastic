package configuration

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"go-prometheus-to-elastic/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func ReadConfiguration() (models.Configuration) {
	pathToFile := os.Getenv("GO_PROMETHEUS_TO-ELASTIC_CONFIGURATION_FILE")
	if pathToFile == "" {
		pathToFile = "/home/pi/go/src/go-prometheus-to-elastic/configuration.yaml"
	}
	yamlFile, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		fmt.Printf("Unable to read configuration file %+v", err)
	}

	var config models.Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Unable to yaml unmarshal configuration file %+v", err)
	} else {
		config.Logger = logs.New(config.Elasticsearch.Url, config.Host)
		config.Logger.Info("Configuration Loaded : %+v \n", config)
	}
	return config
}
