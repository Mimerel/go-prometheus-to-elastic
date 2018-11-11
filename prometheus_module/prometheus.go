package prometheus_module

import (
	"bufio"
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"go-prometheus-to-elastic/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ReadPrometheusData(all *models.Global) (err error) {
	err = GetData(all)
	if err != nil {
		logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to read configuration file %+v", err))
		return err
	}
	// CleanRecords(all)
	err = StructureData(all)
	if err != nil {
		logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to prepare data to send to ElasticSearch %+v", err))
		return err
	}
	logs.Info(all.Config.Elasticsearch.Url,  all.Config.Host, fmt.Sprintf("Data ready to be sent to ElasticSearch"))
	return nil
}


func StructureData(all *models.Global) (err error) {
	for _, value := range all.Records {
		record := new(models.StructuredData)
		record.Labels = make(map[string]string)
		record.Metric = value.Metric
		splitLabels := strings.Split(value.Labels, ",")
		for _, label := range splitLabels {
			labelData := strings.Split(label, "=")
			labelValue := strings.Replace(labelData[1], "\"", "", -1)
			record.Labels[labelData[0]] = labelValue
		}

		moment, err := strconv.ParseInt(value.Timestamp, 10, 64)
		if err != nil {
			logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to convert timestamp to int %+v", err))
			return err
		}
		record.Timestamp = value.Timestamp
		record.Timestamp2 = time.Unix(moment/1000, 0).Format(time.RFC3339)
		record.Value = value.Value
		if err != nil {
			logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to convert value to float %+v", err))
			return err
		}
		all.StructuredData = append(all.StructuredData, *record)
	}
	return nil
}


func GetData(all *models.Global) (err error) {
	logs.Info(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Starting GetData method"))

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(all.Config.Prometheus.Url)
	if err != nil {
		logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to get prometheus metrics %+v", err))
	}
	defer resp.Body.Close()
	logs.Info(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Get data request executed"))

	if resp.StatusCode == http.StatusOK {

		bodyLines, err := LinesFromReader(resp.Body)
		if err != nil {
			logs.Error(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Unable to read body %+v", err))
			return err
		}

		logs.Info(all.Config.Elasticsearch.Url, all.Config.Host, fmt.Sprintf("Starting to read body"))

		for _, line := range bodyLines {
			record := new(models.BodyStruc)
			data := strings.Split(line, "{")
			record.Metric = data[0]
			if len(data) > 1 {
				data = strings.Split(data[1], "}")
				record.Labels = data[0]
				data = strings.Split(strings.Trim(data[1], " "), " ")
				record.Value = data[0]
				record.Timestamp = data[1]
			all.Records = append(all.Records, *record)
			}
		}
	}
	return nil
}

//func CleanRecords(all *models.Global) {
//	for key , value := range all.Records {
//		all.Records[key].Metric = strings.TrimRight(value.Metric, "_01")
//	}
//}


func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}