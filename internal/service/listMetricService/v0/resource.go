package v0

import "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"

type MetricItem struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type MetricData struct {
	Counters []entities.CounterItem `json:"counters"`
	Gauges   []entities.GaugeItem   `json:"gauges"`
}

func (m MetricData) convertToModel() []MetricItem {
	res := make([]MetricItem, 0, len(m.Counters)+len(m.Gauges))

	for _, item := range m.Counters {
		res = append(res, MetricItem{
			Name:  item.MetricName,
			Value: item.MetricValue,
		})
	}

	for _, item := range m.Gauges {
		res = append(res, MetricItem{
			Name:  item.MetricName,
			Value: item.MetricValue,
		})
	}

	return res
}
