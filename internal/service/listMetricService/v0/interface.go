package v0

import "context"

type MetricRepository interface {
	List(ctx context.Context) (items []MetricItem, err error)
}

type MetricItem struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
