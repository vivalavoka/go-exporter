package exporter

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/agent/client"
	"github.com/vivalavoka/go-exporter/cmd/agent/config"
	"github.com/vivalavoka/go-exporter/cmd/agent/crypto"
	"github.com/vivalavoka/go-exporter/cmd/agent/metrics"
)

type Agent struct {
	config       config.Config
	client       *client.Client
	pollCount    metrics.Counter
	metrics      []metrics.Metric
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
	hasher       *crypto.SHA256
}

func New(config config.Config, client *client.Client) *Agent {
	hasher := crypto.New(config.SHAKey)
	return &Agent{
		config:    config,
		client:    client,
		pollCount: metrics.Counter(0),
		hasher:    hasher,
	}
}

func (a *Agent) Start() {
	a.pollTicker = time.NewTicker(a.config.PollInterval)
	a.reportTicker = time.NewTicker(a.config.ReportInterval)
	defer a.pollTicker.Stop()
	defer a.reportTicker.Stop()

	for {
		select {
		case <-a.reportTicker.C:
			log.Info("Report metrics")
			a.ReportMetrics()
		case <-a.pollTicker.C:
			log.Info("Get metrics")
			a.pollCount += 1
			a.metrics = a.GetMetrics()
		}
	}
}

func (a *Agent) ReportMetrics() {
	for _, item := range a.metrics {
		if a.hasher.Enable {
			item.Hash = a.hasher.GetSum(item.String())
		}
		body, err := json.Marshal(&item)

		if err != nil {
			log.Error(err)
			continue
		}

		_, reqErr := a.client.MakeRequest(a.config.Address, body)

		if reqErr != nil {
			log.Error(reqErr)
		}
	}
}

func (a *Agent) GetMetrics() []metrics.Metric {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	random := rand.Float64()

	metrics := []metrics.Metric{
		{ID: "PollCount", MType: metrics.CounterType, Delta: a.pollCount},
		{ID: "Alloc", MType: metrics.GaugeType, Value: metrics.Gauge(stats.Alloc)},
		{ID: "BuckHashSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.BuckHashSys)},
		{ID: "Frees", MType: metrics.GaugeType, Value: metrics.Gauge(stats.Frees)},
		{ID: "GCCPUFraction", MType: metrics.GaugeType, Value: metrics.Gauge(stats.GCCPUFraction)},
		{ID: "GCSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.GCSys)},
		{ID: "HeapAlloc", MType: metrics.GaugeType, Value: metrics.Gauge(stats.HeapAlloc)},
		{ID: "HeapIdle", MType: metrics.GaugeType, Value: metrics.Gauge(stats.HeapIdle)},
		{ID: "HeapInuse", MType: metrics.GaugeType, Value: metrics.Gauge(stats.HeapInuse)},
		{ID: "HeapObjects", MType: metrics.GaugeType, Value: metrics.Gauge(stats.HeapObjects)},
		{ID: "HeapReleased", MType: metrics.GaugeType, Value: metrics.Gauge(stats.HeapReleased)},
		{ID: "HeapSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.HeapSys)},
		{ID: "LastGC", MType: metrics.GaugeType, Value: metrics.Gauge(stats.LastGC)},
		{ID: "Lookups", MType: metrics.GaugeType, Value: metrics.Gauge(stats.Lookups)},
		{ID: "MCacheInuse", MType: metrics.GaugeType, Value: metrics.Gauge(stats.MCacheInuse)},
		{ID: "MCacheSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.MCacheSys)},
		{ID: "MSpanInuse", MType: metrics.GaugeType, Value: metrics.Gauge(stats.MSpanInuse)},
		{ID: "MSpanSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.MSpanSys)},
		{ID: "Mallocs", MType: metrics.GaugeType, Value: metrics.Gauge(stats.Mallocs)},
		{ID: "NextGC", MType: metrics.GaugeType, Value: metrics.Gauge(stats.NextGC)},
		{ID: "NumForcedGC", MType: metrics.GaugeType, Value: metrics.Gauge(stats.NumForcedGC)},
		{ID: "NumGC", MType: metrics.GaugeType, Value: metrics.Gauge(stats.NumGC)},
		{ID: "OtherSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.OtherSys)},
		{ID: "PauseTotalNs", MType: metrics.GaugeType, Value: metrics.Gauge(stats.PauseTotalNs)},
		{ID: "StackInuse", MType: metrics.GaugeType, Value: metrics.Gauge(stats.StackInuse)},
		{ID: "StackSys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.StackSys)},
		{ID: "Sys", MType: metrics.GaugeType, Value: metrics.Gauge(stats.Sys)},
		{ID: "TotalAlloc", MType: metrics.GaugeType, Value: metrics.Gauge(stats.TotalAlloc)},
		{ID: "RandomValue", MType: metrics.GaugeType, Value: metrics.Gauge(random)},
	}

	return metrics
}
