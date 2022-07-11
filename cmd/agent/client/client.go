package client

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type Client struct {
	address    string
	restClient *resty.Client
}

func New(address string, client *resty.Client) *Client {
	return &Client{address: address, restClient: client}
}

func (c *Client) SendMetric(metric *metrics.Metric) (*resty.Response, error) {
	body, err := json.Marshal(&metric)

	if err != nil {
		return nil, err
	}

	return c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"address": c.address,
		}).
		Post("http://{address}/update/")
}

func (c *Client) SendMetrics(metricList []*metrics.Metric) (*resty.Response, error) {
	body, err := json.Marshal(&metricList)

	if err != nil {
		return nil, err
	}

	return c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"address": c.address,
		}).
		Post("http://{address}/updates/")
}
