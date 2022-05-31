package main

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	reportAddress  = "127.0.0.1:8080"
)

type Agent struct {
	client         *resty.Client
	pollCount      counter
	gaugeMetrics   *[]GaugeItem
	counterMetrics *[]CounterItem
}

func (a *Agent) Start() {
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			fmt.Println("Report metrics")
			ReportMetrics(a.client, a.gaugeMetrics, a.counterMetrics)
		case <-pollTicker.C:
			fmt.Println("Get metrics")
			a.pollCount = a.pollCount + 1
			a.gaugeMetrics, a.counterMetrics = GetMetrics(a.pollCount)
		}
	}
}
