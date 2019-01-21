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
		logs.Error("", "", fmt.Sprintf("Unable to read configuration file %+v", err))
	}

	var config models.Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logs.Error("", config.Host, fmt.Sprintf("Unable to yaml unmarshal configuration file %+v", err))
	} else {
		logs.Info(config.ElasticsearchLogs.Url, config.Host, fmt.Sprint("Configuration Loaded : %+v \n", config))
	}
	return config
}
