package main

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	restClient *resty.Client
}

func NewClient(client *resty.Client) *Client {
	return &Client{restClient: client}
}

func (c *Client) MakeRequest(body []byte) (*resty.Response, error) {
	return c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetPathParams(map[string]string{
			"address": reportAddress,
		}).
		Post("http://{address}/update/")
}
