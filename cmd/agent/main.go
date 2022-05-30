package main

import (
	"net/http"
)

func main() {
	agent := Agent{
		client:    &http.Client{},
		pollCount: counter(0),
	}

	agent.Start()
}
