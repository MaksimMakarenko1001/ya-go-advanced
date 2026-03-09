package pkg

// MetricType represents the type of a metric.
type MetricType = string

const (
	// MetricTypeGauge represents a gauge metric type.
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeCounter represents a counter metric type.
	MetricTypeCounter MetricType = "counter"
)
