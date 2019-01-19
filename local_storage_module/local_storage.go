package local_storage_module

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"go-prometheus-to-elastic/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func WriteLastValues (all *models.Global) (error) {
	pathToFile := os.Getenv("GO_PROMETHEUS_TO-ELASTIC_LOCAL")
	if pathToFile == "" {
		pathToFile = "/home/pi/go/src/go-prometheus-to-elastic/local_storage.yaml"
	}

	yamlFile, err := yaml.Marshal(all.StructuredData)
	if err != nil {
		logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to yaml marshal local_storage file %+v", err))
	}
	err = ioutil.WriteFile(pathToFile, yamlFile, 0777)
	if err != nil {
		logs.Error("", "", fmt.Sprintf("Unable to write local storage file %+v", err))
	} else {
		logs.Info(all.Config.Elasticsearch.Url, all.Config.Host, "Values stored in local storage\n")
	}
	return nil
}




