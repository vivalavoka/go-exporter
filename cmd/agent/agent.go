package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Agent struct {
	client    *Client
	pollCount counter
	metrics   []Metrics
}

func NewAgent(client *Client) *Agent {
	return &Agent{
		client:    client,
		pollCount: counter(0),
	}
}

func (a *Agent) Start() {
	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			log.Info("Report metrics")
			ReportMetrics(a.client, a.metrics)
		case <-pollTicker.C:
			log.Info("Get metrics")
			a.pollCount += 1
			a.metrics = GetMetrics(a.pollCount)
		}
	}
}
