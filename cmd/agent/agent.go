package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	reportAddress  = "127.0.0.1:8080"
)

type Agent struct {
	client         *Client
	pollCount      counter
	gaugeMetrics   []GaugeItem
	counterMetrics []CounterItem
}

func NewAgent(client *Client) *Agent {
	return &Agent{
		client:    client,
		pollCount: counter(0),
	}
}

func (a *Agent) Start() {
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			log.Info("Report metrics")
			ReportMetrics(a.client, a.gaugeMetrics, a.counterMetrics)
		case <-pollTicker.C:
			log.Info("Get metrics")
			a.pollCount += 1
			a.gaugeMetrics, a.counterMetrics = GetMetrics(a.pollCount)
		}
	}
}
