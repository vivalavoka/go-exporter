package exporter

import (
	"sync"

	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type MetricCache struct {
	metrics []*metrics.Metric
	mutex   sync.RWMutex
}

func NewCache() *MetricCache {
	return &MetricCache{}
}

func (mc *MetricCache) Set(metrics []*metrics.Metric) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics = metrics
}

func (mc *MetricCache) Get() []*metrics.Metric {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	return mc.metrics
}
