package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	httpClient *http.Client
	host       string
}

func NewClient(cfg HTTPClientConfig) *Client {
	httpClient := &http.Client{Timeout: cfg.Timeout}

	return &Client{
		host:       cfg.Address,
		httpClient: httpClient,
	}
}

func (c *Client) SendGaugeMetric(metricName string, value float64) (err error) {
	valueStr := strconv.FormatFloat(value, 'f', -1, 64)

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/update/gauge/" + metricName + "/" + valueStr,
	}

	resp, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("error sending gauge metric: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send gauge metric %s: %d", metricName, resp.StatusCode)
	}

	return nil
}

func (c *Client) SendCounterMetric(metricName string, value int64) (err error) {
	valueStr := strconv.FormatInt(value, 10)

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "/update/counter/" + metricName + "/" + valueStr,
	}

	resp, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("error sending counter metric: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send counter metric %s: %d", metricName, resp.StatusCode)
	}

	return nil
}
