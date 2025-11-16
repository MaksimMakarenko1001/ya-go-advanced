package agent

import (
	"errors"
	"log"
	"math/rand"
	"runtime"
	"sync"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func genCounters(pollCount *int64) []models.Metrics {
	*pollCount++
	return []models.Metrics{
		{
			ID:    "PollCount",
			MType: pkg.MetricTypeCounter,
			Delta: pollCount,
		},
	}
}

func genGauge(memStats *runtime.MemStats) []models.Metrics {
	runtime.ReadMemStats(memStats)
	return []models.Metrics{
		{
			ID:    "Alloc",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.Alloc)),
		},
		{
			ID:    "BuckHashSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.BuckHashSys)),
		},
		{
			ID:    "Frees",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.Frees)),
		},
		{
			ID:    "GCCPUFraction",
			MType: pkg.MetricTypeGauge,
			Value: &memStats.GCCPUFraction,
		},
		{
			ID:    "GCSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.GCSys)),
		},
		{
			ID:    "HeapAlloc",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.HeapAlloc)),
		},
		{
			ID:    "HeapIdle",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.HeapIdle)),
		},
		{
			ID:    "HeapInuse",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.HeapInuse)),
		},
		{
			ID:    "HeapObjects",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.HeapObjects)),
		},
		{
			ID:    "HeapReleased",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.HeapReleased)),
		},
		{
			ID:    "HeapSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.HeapSys)),
		},
		{
			ID:    "LastGC",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.LastGC)),
		},
		{
			ID:    "Lookups",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.Lookups)),
		},
		{
			ID:    "MCacheInuse",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.MCacheInuse)),
		},
		{
			ID:    "MCacheSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.MCacheSys)),
		},
		{
			ID:    "MSpanInuse",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.MSpanInuse)),
		},
		{
			ID:    "MSpanSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.MSpanSys)),
		},
		{
			ID:    "Mallocs",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.Mallocs)),
		},
		{
			ID:    "NextGC",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.NextGC)),
		},
		{
			ID:    "NumForcedGC",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.NumForcedGC)),
		},
		{
			ID:    "NumGC",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.NumGC)),
		},
		{
			ID:    "OtherSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.OtherSys)),
		},
		{
			ID:    "PauseTotalNs",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.PauseTotalNs)),
		},
		{
			ID:    "StackInuse",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.StackInuse)),
		},
		{
			ID:    "StackSys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.StackSys)),
		},
		{
			ID:    "Sys",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.Sys)),
		},
		{
			ID:    "TotalAlloc",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(float64(memStats.TotalAlloc)),
		},
		{
			ID:    "RandomValue",
			MType: pkg.MetricTypeGauge,
			Value: pkg.ToPtr(rand.Float64()),
		},
	}
}

func genExtraGauge() []models.Metrics {
	var errs []error
	slice := make([]models.Metrics, 0, 3)

	v, err := mem.VirtualMemory()
	if err != nil {
		errs = append(errs, err)
	} else {
		slice = append(slice, []models.Metrics{
			{
				ID:    "TotalMemory",
				MType: pkg.MetricTypeGauge,
				Value: pkg.ToPtr(float64(v.Total)),
			},
			{
				ID:    "FreeMemory",
				MType: pkg.MetricTypeGauge,
				Value: pkg.ToPtr(float64(v.Free)),
			},
		}...)
	}

	cpuPercentages, err := cpu.Percent(0, false)
	if err != nil {
		errs = append(errs, err)
	} else {
		slice = append(slice, models.Metrics{
			ID:    "CPUutilization1",
			MType: pkg.MetricTypeGauge,
			Value: &cpuPercentages[0],
		})
	}

	if len(errs) > 0 {
		log.Printf("gopsutil metrics not ok, %v\n", errors.Join(err))
	}
	return slice
}

func gen(doneCh <-chan struct{}, input []models.Metrics) <-chan models.Metrics {
	ch := make(chan models.Metrics)
	go func() {
		defer close(ch)
		for _, i := range input {
			select {
			case <-doneCh:
				return
			case ch <- i:
			}
		}
	}()

	return ch
}

func fanIn(doneCh <-chan struct{}, inChs []<-chan models.Metrics) <-chan models.Metrics {
	var wg sync.WaitGroup
	ch := make(chan models.Metrics)

	for _, inCh := range inChs {
		wg.Add(1)

		go func(channel <-chan models.Metrics) {
			defer wg.Done()

			for i := range channel {
				select {
				case <-doneCh:
					return
				case ch <- i:
				}
			}
		}(inCh)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func batched(inCh <-chan models.Metrics, size int) <-chan []models.Metrics {
	ch := make(chan []models.Metrics)
	batch := make([]models.Metrics, 0, size)

	go func() {
		defer close(ch)

		for metric := range inCh {
			batch = append(batch, metric)

			if len(batch) == max(size, 1) {
				ch <- batch
				batch = batch[:0]
			}
		}
		if len(batch) > 0 {
			ch <- batch
		}
	}()

	return ch
}
