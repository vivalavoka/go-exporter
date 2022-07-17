package exporter

import (
	"sync"

	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type MetricCache struct {
	metrics map[string]*metrics.Metric
	mutex   sync.RWMutex
}

func NewCache() *MetricCache {
	return &MetricCache{
		metrics: make(map[string]*metrics.Metric),
	}
}

func (mc *MetricCache) Set(metric *metrics.Metric) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics[metric.ID] = metric
}

func (mc *MetricCache) Get(ID string) *metrics.Metric {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	metric, ok := mc.metrics[ID]
	if !ok {
		return nil
	}

	return metric
}

func (mc *MetricCache) GetAll() []*metrics.Metric {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	var metricList []*metrics.Metric
	for _, metric := range mc.metrics {
		metricList = append(metricList, metric)
	}
	return metricList
}
