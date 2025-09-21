package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/agent"
)

func main() {
	log.Println("agent starting")

	if err := run(); err != nil {
		panic(err)
	}

	log.Println("agent stoped")
}

func run() error {
	cfg := &agent.Config{}
	cfg.LoadConfig()

	ticker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer ticker.Stop()
	defer reportTicker.Stop()

	collectGaugeMetricsFunc := gaugeMetricsCollector()
	collectCounterMetricsFunc := counterMetricsCollector()

	var gaugeMetrics map[string]float64
	var counterMetrics map[string]int64

	cli := agent.NewClient(cfg.HTTP)

	for {
		select {
		case <-ticker.C:
			log.Printf("Collecting metrics")
			gaugeMetrics = collectGaugeMetricsFunc()
			counterMetrics = collectCounterMetricsFunc()

		case <-reportTicker.C:
			log.Printf("Reporting metrics")
			// TODO implement fan-out technique
			for name, value := range gaugeMetrics {
				cli.SendGaugeMetric(name, value)
			}
			for name, value := range counterMetrics {
				cli.SendCounterMetric(name, value)
			}
		}
	}
}

func gaugeMetricsCollector() func() map[string]float64 {
	memStats := runtime.MemStats{}
	return func() map[string]float64 {
		runtime.ReadMemStats(&memStats)
		return map[string]float64{
			"Alloc":         float64(memStats.Alloc),
			"BuckHashSys":   float64(memStats.BuckHashSys),
			"Frees":         float64(memStats.Frees),
			"GCCPUFraction": memStats.GCCPUFraction,
			"GCSys":         float64(memStats.GCSys),
			"HeapAlloc":     float64(memStats.HeapAlloc),
			"HeapIdle":      float64(memStats.HeapIdle),
			"HeapInuse":     float64(memStats.HeapInuse),
			"HeapObjects":   float64(memStats.HeapObjects),
			"HeapReleased":  float64(memStats.HeapReleased),
			"HeapSys":       float64(memStats.HeapSys),
			"LastGC":        float64(memStats.LastGC),
			"Lookups":       float64(memStats.Lookups),
			"MCacheInuse":   float64(memStats.MCacheInuse),
			"MCacheSys":     float64(memStats.MCacheSys),
			"MSpanInuse":    float64(memStats.MSpanInuse),
			"MSpanSys":      float64(memStats.MSpanSys),
			"Mallocs":       float64(memStats.Mallocs),
			"NextGC":        float64(memStats.NextGC),
			"NumForcedGC":   float64(memStats.NumForcedGC),
			"NumGC":         float64(memStats.NumGC),
			"OtherSys":      float64(memStats.OtherSys),
			"PauseTotalNs":  float64(memStats.PauseTotalNs),
			"StackInuse":    float64(memStats.StackInuse),
			"StackSys":      float64(memStats.StackSys),
			"Sys":           float64(memStats.Sys),
			"TotalAlloc":    float64(memStats.TotalAlloc),
			"RandomValue":   rand.Float64(),
		}
	}

}

func counterMetricsCollector() func() map[string]int64 {
	var pollCount int64
	return func() map[string]int64 {
		pollCount += 1
		return map[string]int64{
			"PollCount": pollCount,
		}
	}
}
