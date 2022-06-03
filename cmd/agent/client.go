package main

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	restClient *resty.Client
}

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
}

func NewClient(client *resty.Client) *Client {
	return &Client{restClient: client}
}

func (c *Client) MakeRequest(params *UpdateParams) (*resty.Response, error) {
	return c.restClient.R().
		SetHeader("Content-Type", "text/plain").
		SetBody("").
		SetPathParams(map[string]string{
			"address": reportAddress,
			"type":    params.MetricType,
			"name":    params.MetricName,
			"value":   params.MetricValue,
		}).
		Post("http://{address}/update/{type}/{name}/{value}")

}
