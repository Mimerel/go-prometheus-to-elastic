package local_storage_module

import (
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
		all.Config.Logger.Error("Unable to yaml marshal local_storage file %+v", err)
	}
	err = ioutil.WriteFile(pathToFile, yamlFile, 0777)
	if err != nil {
		all.Config.Logger.Error("Unable to write local storage file %+v", err)
	} else {
		all.Config.Logger.Info("Values stored in local storage\n")
	}
	return nil
}




