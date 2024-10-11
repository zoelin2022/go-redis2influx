package models

type Point struct {
	Name   string             `json:"name"`
	Time   int64              `json:"timestamp"`
	Tags   map[string]string  `json:"tags"`
	Fields map[string]float64 `json:"fields"`
}

type MetricsData struct {
	Metrics []Point `json:"metrics"`
}
