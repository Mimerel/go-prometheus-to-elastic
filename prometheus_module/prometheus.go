package prometheus_module

import (
	"bufio"
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
		all.Config.Logger.Error("Unable to read configuration file %+v", err)
		return err
	}
	// CleanRecords(all)
	err = StructureData(all)
	if err != nil {
		all.Config.Logger.Error("Unable to prepare data to send to ElasticSearch %+v", err)
		return err
	}
	all.Config.Logger.Info("Data ready to be sent to ElasticSearch")
	return nil
}

/**
Converts the line text array to structured data
and enriches the data
 */
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
			all.Config.Logger.Error("Unable to convert timestamp to int %+v", err)
			return err
		}
		record.Timestamp = value.Timestamp
		record.Timestamp2 = time.Unix(moment/1000, 0).Format(time.RFC3339)
		record.Value = value.Value
		if err != nil {
			all.Config.Logger.Error("Unable to convert value to float %+v", err)
			return err
		}
		all.StructuredData = append(all.StructuredData, *record)
	}
	return nil
}

/**
Requests data from Prometheus
The data arrives as text
The text received is then converted into a structure for future use
 */
func GetData(all *models.Global) (err error) {
	all.Config.Logger.Info("Starting GetData method")

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(all.Config.Prometheus.Url)
	if err != nil {
		all.Config.Logger.Error("Unable to get prometheus metrics %+v", err)
	}
	defer resp.Body.Close()
	all.Config.Logger.Info("Get data request executed")

	if resp.StatusCode == http.StatusOK {

		bodyLines, err := LinesFromReader(resp.Body)
		if err != nil {
			all.Config.Logger.Error("Unable to read body %+v", err)
			return err
		}

		all.Config.Logger.Info("Starting to read body")

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

/**
Reads all text lines received by prometheus
and adds each line to an array of strings
 */
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