package main

import (
	"github.com/go-resty/resty/v2"
)

func main() {

	client := resty.New()
	agent := Agent{
		client:    client,
		pollCount: counter(0),
	}

	agent.Start()
}
