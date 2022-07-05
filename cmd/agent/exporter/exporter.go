package exporter

import (
	"math/rand"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/agent/client"
	"github.com/vivalavoka/go-exporter/cmd/agent/config"
	"github.com/vivalavoka/go-exporter/internal/crypto"
	"github.com/vivalavoka/go-exporter/internal/metrics"
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

		_, err := a.client.SendMetric(a.config.Address, &item)

		if err != nil {
			log.Error(err)
		}
	}
}

func (a *Agent) GetMetrics() []metrics.Metric {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	random := rand.Float64()

	alloc := metrics.Gauge(stats.Alloc)
	buckHashSys := metrics.Gauge(stats.BuckHashSys)
	frees := metrics.Gauge(stats.Frees)
	gCCPUFraction := metrics.Gauge(stats.GCCPUFraction)
	gCSys := metrics.Gauge(stats.GCSys)
	heapAlloc := metrics.Gauge(stats.HeapAlloc)
	heapIdle := metrics.Gauge(stats.HeapIdle)
	heapInuse := metrics.Gauge(stats.HeapInuse)
	heapObjects := metrics.Gauge(stats.HeapObjects)
	heapReleased := metrics.Gauge(stats.HeapReleased)
	heapSys := metrics.Gauge(stats.HeapSys)
	lastGC := metrics.Gauge(stats.LastGC)
	lookups := metrics.Gauge(stats.Lookups)
	mCacheInuse := metrics.Gauge(stats.MCacheInuse)
	mCacheSys := metrics.Gauge(stats.MCacheSys)
	mSpanInuse := metrics.Gauge(stats.MSpanInuse)
	mSpanSys := metrics.Gauge(stats.MSpanSys)
	mallocs := metrics.Gauge(stats.Mallocs)
	nextGC := metrics.Gauge(stats.NextGC)
	numForcedGC := metrics.Gauge(stats.NumForcedGC)
	numGC := metrics.Gauge(stats.NumGC)
	otherSys := metrics.Gauge(stats.OtherSys)
	pauseTotalNs := metrics.Gauge(stats.PauseTotalNs)
	stackInuse := metrics.Gauge(stats.StackInuse)
	stackSys := metrics.Gauge(stats.StackSys)
	sys := metrics.Gauge(stats.Sys)
	totalAlloc := metrics.Gauge(stats.TotalAlloc)
	randomValue := metrics.Gauge(random)

	metrics := []metrics.Metric{
		{ID: "PollCount", MType: metrics.CounterType, Delta: &a.pollCount},
		{ID: "Alloc", MType: metrics.GaugeType, Value: &alloc},
		{ID: "BuckHashSys", MType: metrics.GaugeType, Value: &buckHashSys},
		{ID: "Frees", MType: metrics.GaugeType, Value: &frees},
		{ID: "GCCPUFraction", MType: metrics.GaugeType, Value: &gCCPUFraction},
		{ID: "GCSys", MType: metrics.GaugeType, Value: &gCSys},
		{ID: "HeapAlloc", MType: metrics.GaugeType, Value: &heapAlloc},
		{ID: "HeapIdle", MType: metrics.GaugeType, Value: &heapIdle},
		{ID: "HeapInuse", MType: metrics.GaugeType, Value: &heapInuse},
		{ID: "HeapObjects", MType: metrics.GaugeType, Value: &heapObjects},
		{ID: "HeapReleased", MType: metrics.GaugeType, Value: &heapReleased},
		{ID: "HeapSys", MType: metrics.GaugeType, Value: &heapSys},
		{ID: "LastGC", MType: metrics.GaugeType, Value: &lastGC},
		{ID: "Lookups", MType: metrics.GaugeType, Value: &lookups},
		{ID: "MCacheInuse", MType: metrics.GaugeType, Value: &mCacheInuse},
		{ID: "MCacheSys", MType: metrics.GaugeType, Value: &mCacheSys},
		{ID: "MSpanInuse", MType: metrics.GaugeType, Value: &mSpanInuse},
		{ID: "MSpanSys", MType: metrics.GaugeType, Value: &mSpanSys},
		{ID: "Mallocs", MType: metrics.GaugeType, Value: &mallocs},
		{ID: "NextGC", MType: metrics.GaugeType, Value: &nextGC},
		{ID: "NumForcedGC", MType: metrics.GaugeType, Value: &numForcedGC},
		{ID: "NumGC", MType: metrics.GaugeType, Value: &numGC},
		{ID: "OtherSys", MType: metrics.GaugeType, Value: &otherSys},
		{ID: "PauseTotalNs", MType: metrics.GaugeType, Value: &pauseTotalNs},
		{ID: "StackInuse", MType: metrics.GaugeType, Value: &stackInuse},
		{ID: "StackSys", MType: metrics.GaugeType, Value: &stackSys},
		{ID: "Sys", MType: metrics.GaugeType, Value: &sys},
		{ID: "TotalAlloc", MType: metrics.GaugeType, Value: &totalAlloc},
		{ID: "RandomValue", MType: metrics.GaugeType, Value: &randomValue},
	}

	return metrics
}
