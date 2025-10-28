package entities

type CounterItem struct {
	MetricName  string `json:"metric_name"`
	MetricValue int64  `json:"metric_value"`
}

type GaugeItem struct {
	MetricName  string  `json:"metric_name"`
	MetricValue float64 `json:"metric_value"`
}
