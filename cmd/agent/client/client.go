package client

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type Client struct {
	restClient *resty.Client
}

func New(client *resty.Client) *Client {
	return &Client{restClient: client}
}

func (c *Client) SendMetric(address string, metric *metrics.Metric) (*resty.Response, error) {
	body, err := json.Marshal(&metric)

	if err != nil {
		return nil, err
	}

	return c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"address": address,
		}).
		Post("http://{address}/update/")
}
