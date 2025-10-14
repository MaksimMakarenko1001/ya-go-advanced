package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
)

type Client struct {
	httpClient *http.Client
	host       string
	memStats   runtime.MemStats
	pollCount  int64
}

func NewClient(cfg HTTPClientConfig) *Client {
	httpClient := &http.Client{Timeout: cfg.Timeout}

	return &Client{
		host:       cfg.Address,
		httpClient: httpClient,
	}
}

func (c *Client) sendGaugeMetric(metricName string, value float64) (err error) {
	valueStr := strconv.FormatFloat(value, 'f', -1, 64)

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/update/gauge/" + metricName + "/" + valueStr,
	}

	resp, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("error sending gauge metric: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send gauge metric %s: %d", metricName, resp.StatusCode)
	}

	return nil
}

func (c *Client) sendGaugeMetricJSON(metricName string, value float64) (err error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/update/",
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(models.Metrics{
		ID:    metricName,
		MType: "gauge",
		Value: &value,
	})
	if err != nil {
		return fmt.Errorf("gauge encoder not ok, %w", err)
	}

	r, err := newGZipRequest(http.MethodPost, u.String(), buf.Bytes())
	if err != nil {
		return fmt.Errorf("gauge request not ok, %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("gauge http not ok, %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gauge response status not ok, %s: %d", metricName, resp.StatusCode)
	}

	return nil
}

func (c *Client) sendCounterMetric(metricName string, value int64) (err error) {
	valueStr := strconv.FormatInt(value, 10)

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/update/counter/" + metricName + "/" + valueStr,
	}

	resp, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("error sending counter metric: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send counter metric %s: %d", metricName, resp.StatusCode)
	}

	return nil
}

func (c *Client) sendCounterMetricJSON(metricName string, value int64) (err error) {

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/update/",
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(models.Metrics{
		ID:    metricName,
		MType: "counter",
		Delta: &value,
	})

	if err != nil {
		return fmt.Errorf("counter encoder not ok, %w", err)
	}

	r, err := newGZipRequest(http.MethodPost, u.String(), buf.Bytes())
	if err != nil {
		return fmt.Errorf("counter request not ok, %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("counter http not ok, %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("counter response status not ok, %s: %d", metricName, resp.StatusCode)
	}

	return nil
}

func (c *Client) collectCounterMetrics() map[string]int64 {
	c.pollCount += 1
	return map[string]int64{
		"PollCount": c.pollCount,
	}
}

func (c *Client) collectGaugeMetrics() map[string]float64 {
	runtime.ReadMemStats(&c.memStats)
	return map[string]float64{
		"Alloc":         float64(c.memStats.Alloc),
		"BuckHashSys":   float64(c.memStats.BuckHashSys),
		"Frees":         float64(c.memStats.Frees),
		"GCCPUFraction": c.memStats.GCCPUFraction,
		"GCSys":         float64(c.memStats.GCSys),
		"HeapAlloc":     float64(c.memStats.HeapAlloc),
		"HeapIdle":      float64(c.memStats.HeapIdle),
		"HeapInuse":     float64(c.memStats.HeapInuse),
		"HeapObjects":   float64(c.memStats.HeapObjects),
		"HeapReleased":  float64(c.memStats.HeapReleased),
		"HeapSys":       float64(c.memStats.HeapSys),
		"LastGC":        float64(c.memStats.LastGC),
		"Lookups":       float64(c.memStats.Lookups),
		"MCacheInuse":   float64(c.memStats.MCacheInuse),
		"MCacheSys":     float64(c.memStats.MCacheSys),
		"MSpanInuse":    float64(c.memStats.MSpanInuse),
		"MSpanSys":      float64(c.memStats.MSpanSys),
		"Mallocs":       float64(c.memStats.Mallocs),
		"NextGC":        float64(c.memStats.NextGC),
		"NumForcedGC":   float64(c.memStats.NumForcedGC),
		"NumGC":         float64(c.memStats.NumGC),
		"OtherSys":      float64(c.memStats.OtherSys),
		"PauseTotalNs":  float64(c.memStats.PauseTotalNs),
		"StackInuse":    float64(c.memStats.StackInuse),
		"StackSys":      float64(c.memStats.StackSys),
		"Sys":           float64(c.memStats.Sys),
		"TotalAlloc":    float64(c.memStats.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

}

func (c *Client) Start(pollInterval time.Duration, reportInterval time.Duration) error {
	ticker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer ticker.Stop()
	defer reportTicker.Stop()

	var gaugeMetrics map[string]float64
	var counterMetrics map[string]int64
	var err error

	for {
		select {
		case <-ticker.C:
			log.Println("Collecting metrics")
			gaugeMetrics = c.collectGaugeMetrics()
			counterMetrics = c.collectCounterMetrics()

		case <-reportTicker.C:
			log.Println("Reporting metrics")
			// TODO implement fan-out technique
			for name, value := range gaugeMetrics {
				err = c.sendGaugeMetricJSON(name, value)
				if err != nil {
					log.Println(err.Error())
				}
			}
			for name, value := range counterMetrics {
				err = c.sendCounterMetricJSON(name, value)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}

}

func newGZipRequest(method string, url string, body []byte) (*http.Request, error) {
	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)

	_, err := zw.Write(body)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	return request, nil
}
