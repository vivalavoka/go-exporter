package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
}

func MakeRequest(client *resty.Client, params *UpdateParams) {
	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody("").
		SetPathParams(map[string]string{
			"address": reportAddress,
			"type":    params.MetricType,
			"name":    params.MetricName,
			"value":   params.MetricValue,
		}).
		Post("http://{address}/update/{type}/{name}/{value}")

	if err != nil {
		fmt.Println(err)
	}
}
