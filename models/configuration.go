package models

import "time"

type Configuration struct {
	Elasticsearch Elasticsearch `yaml:"elasticSearch,omitempty"`
	Prometheus Prometheus `yaml:"prometheus,omitempty"`
	Host string `yaml:"go-log-to-elastic,omitempty"`
	Scrape_Interval time.Duration `yaml:"scrape_interval"`
}

type Elasticsearch struct {
	Url string `yaml:"url,omitempty"`
}

type Prometheus struct {
	Url string `yaml:"url,omitempty"`
}