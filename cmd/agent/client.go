package main

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	restClient *resty.Client
}

func NewClient(client *resty.Client) *Client {
	return &Client{restClient: client}
}

func (c *Client) MakeRequest(metric *Metrics) (*resty.Response, error) {
	response, _ := json.Marshal(metric)

	return c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(response).
		SetPathParams(map[string]string{
			"address": reportAddress,
		}).
		Post("http://{address}/update/")

}
