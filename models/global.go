package models


type Global struct {
	Config Configuration
	Records []BodyStruc
	StructuredData []StructuredData
}

type StructuredData struct {
	Metric string
	Labels map[string]string
	Timestamp string
	Timestamp2 string
	Value string
}