package main

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"go-prometheus-to-elastic/configuration"
	"go-prometheus-to-elastic/elasticsearch_module"
	"go-prometheus-to-elastic/local_storage_module"
	"go-prometheus-to-elastic/models"
	"go-prometheus-to-elastic/prometheus_module"
	"time"
)


func main() {
	config := configuration.ReadConfiguration()
	for _ = range time.Tick(config.Scrape_Interval*time.Second) {
		all := new(models.Global)
		logs.Info(all.Config.ElasticsearchLogs.Url, config.Host, fmt.Sprintf("Requesting Metrics"))
		all.Config = config
		err := prometheus_module.ReadPrometheusData(all)
		if err != nil {
			logs.Error(all.Config.ElasticsearchLogs.Url, all.Config.Host, fmt.Sprintf("Unable to get prometheus metrics %+v", err))
		}
		err = local_storage_module.WriteLastValues(all)
		if err != nil {
			logs.Error(all.Config.ElasticsearchLogs.Url, all.Config.Host, fmt.Sprintf("Unable to write local storage %+v", err))
		}
		err = elasticsearch_module.SendMetrics(all)
		if err != nil {
			logs.Error(all.Config.ElasticsearchLogs.Url, all.Config.Host, fmt.Sprintf("Unable to send metrics to elasticsearch %+v", err))
		}
	}

}

