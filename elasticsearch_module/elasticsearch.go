package elasticsearch_module

import (
	"bytes"
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"go-prometheus-to-elastic/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func SendMetrics(all *models.Global) (err error) {
	logs.Info("", all.Config.Host, fmt.Sprintf("create post request body \n"))

	body, _ := createBody(all)
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := all.Config.Elasticsearch.Url + "/_bulk"
	logs.Info("", all.Config.Host, fmt.Sprintf("Starting to post body \n"))

	resp, err := client.Post(postingUrl, "application/json" , bytes.NewBuffer([]byte(body)))
	if err != nil {
		logs.Error("", all.Config.Host, fmt.Sprintf("Failed to post request to elasticSearch %s \n", err))
	}
	temp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("", all.Config.Host, fmt.Sprintf("Failed to read response from elasticSearch %s \n", err))
	}
	logs.Info("", all.Config.Host, fmt.Sprintf("response Body : %s \n", temp))

	resp.Body.Close()
	logs.Info("", all.Config.Host, fmt.Sprintf("Metrics successfully sent to Elasticsearch \n"))

	return nil
}


func createBody(all *models.Global) (body string, err error) {
	for _, record := range all.StructuredData {
		body = body + "{ 'update': { '_id': '" + record.Timestamp+"_"+ record.Metric+"_"+record.Labels["host"] + "', '_type': 'metrics', '_index': 'prometheus' }}\n"
		body = body + "{ 'doc': { "
		i := 0
		for key, value := range record.Labels {
			if i != 0 {
				body = body + ", "
			}
			body = body + "'" + key + "': '" + value + "'"
			i = i + 1
		}
		body = body + ", 'value': " + record.Value
		body = body + ", 'timestamp': " + record.Timestamp
		body = body + ", 'timestamp2': '" + record.Timestamp2 + "'"
		body = body + "}, 'doc_as_upsert' : true }\n"
		}
	body = strings.Replace(body, "'", "\"", -1)
	return body, nil
}
