package entities

import "time"

type CounterItem struct {
	MetricType  string    `json:"metric_type"`
	MetricName  string    `json:"metric_name"`
	MetricValue int64     `json:"metric_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GaugeItem struct {
	MetricType  string    `json:"metric_type"`
	MetricName  string    `json:"metric_name"`
	MetricValue float64   `json:"metric_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
