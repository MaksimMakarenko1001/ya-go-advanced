package v0

type MetricRepository interface {
	List() (items []MetricItem)
}

type MetricItem struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
