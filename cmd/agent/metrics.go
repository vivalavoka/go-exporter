package main

import (
	"encoding/json"
	"math/rand"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type gauge float64
type counter int64

const GaugeType = "gauge"
const CounterType = "counter"

type Metrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func ReportMetrics(client *Client, metrics []Metrics) {
	for _, item := range metrics {
		body, err := json.Marshal(&item)

		if err != nil {
			log.Error(err)
		}

		_, reqErr := client.MakeRequest(body)

		if reqErr != nil {
			log.Error(reqErr)
		}
	}
}

func GetMetrics(pollCount counter) []Metrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	random := rand.Float64()
	log.Info(random)

	metrics := []Metrics{
		{ID: "PollCount", MType: CounterType, Delta: pollCount},
		{ID: "Alloc", MType: GaugeType, Value: gauge(stats.Alloc)},
		{ID: "BuckHashSys", MType: GaugeType, Value: gauge(stats.BuckHashSys)},
		{ID: "Frees", MType: GaugeType, Value: gauge(stats.Frees)},
		{ID: "GCCPUFraction", MType: GaugeType, Value: gauge(stats.GCCPUFraction)},
		{ID: "GCSys", MType: GaugeType, Value: gauge(stats.GCSys)},
		{ID: "HeapAlloc", MType: GaugeType, Value: gauge(stats.HeapAlloc)},
		{ID: "HeapIdle", MType: GaugeType, Value: gauge(stats.HeapIdle)},
		{ID: "HeapInuse", MType: GaugeType, Value: gauge(stats.HeapInuse)},
		{ID: "HeapObjects", MType: GaugeType, Value: gauge(stats.HeapObjects)},
		{ID: "HeapReleased", MType: GaugeType, Value: gauge(stats.HeapReleased)},
		{ID: "HeapSys", MType: GaugeType, Value: gauge(stats.HeapSys)},
		{ID: "LastGC", MType: GaugeType, Value: gauge(stats.LastGC)},
		{ID: "Lookups", MType: GaugeType, Value: gauge(stats.Lookups)},
		{ID: "MCacheInuse", MType: GaugeType, Value: gauge(stats.MCacheInuse)},
		{ID: "MCacheSys", MType: GaugeType, Value: gauge(stats.MCacheSys)},
		{ID: "MSpanInuse", MType: GaugeType, Value: gauge(stats.MSpanInuse)},
		{ID: "MSpanSys", MType: GaugeType, Value: gauge(stats.MSpanSys)},
		{ID: "Mallocs", MType: GaugeType, Value: gauge(stats.Mallocs)},
		{ID: "NextGC", MType: GaugeType, Value: gauge(stats.NextGC)},
		{ID: "NumForcedGC", MType: GaugeType, Value: gauge(stats.NumForcedGC)},
		{ID: "NumGC", MType: GaugeType, Value: gauge(stats.NumGC)},
		{ID: "OtherSys", MType: GaugeType, Value: gauge(stats.OtherSys)},
		{ID: "PauseTotalNs", MType: GaugeType, Value: gauge(stats.PauseTotalNs)},
		{ID: "StackInuse", MType: GaugeType, Value: gauge(stats.StackInuse)},
		{ID: "StackSys", MType: GaugeType, Value: gauge(stats.StackSys)},
		{ID: "Sys", MType: GaugeType, Value: gauge(stats.Sys)},
		{ID: "TotalAlloc", MType: GaugeType, Value: gauge(stats.TotalAlloc)},
		{ID: "RandomValue", MType: GaugeType, Value: gauge(random)},
	}

	return metrics
}
