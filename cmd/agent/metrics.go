package main

import (
	"fmt"
	"math/rand"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type gauge float64
type counter int64

const GaugeType = "gauge"
const CounterType = "counter"

type GaugeItem struct {
	Name  string
	Value gauge
}

type CounterItem struct {
	Name  string
	Value counter
}

func ReportMetrics(client *Client, gaugeMetrics []GaugeItem, counterMetrics []CounterItem) {
	for _, item := range gaugeMetrics {
		params := UpdateParams{
			MetricType:  GaugeType,
			MetricName:  item.Name,
			MetricValue: fmt.Sprintf("%f", item.Value),
		}
		_, err := client.MakeRequest(&params)

		if err != nil {
			log.Error(err)
		}
	}

	for _, item := range counterMetrics {
		params := UpdateParams{
			MetricType:  CounterType,
			MetricName:  item.Name,
			MetricValue: fmt.Sprintf("%d", item.Value),
		}
		_, err := client.MakeRequest(&params)

		if err != nil {
			log.Error(err)
		}
	}
}

func GetMetrics(pollCount counter) ([]GaugeItem, []CounterItem) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	random := rand.Float64()
	log.Info(random)

	counterMetrics := []CounterItem{
		{"PollCount", pollCount},
	}

	gaugeMetrics := []GaugeItem{
		{"Alloc", gauge(stats.Alloc)},
		{"BuckHashSys", gauge(stats.BuckHashSys)},
		{"Frees", gauge(stats.Frees)},
		{"GCCPUFraction", gauge(stats.GCCPUFraction)},
		{"GCSys", gauge(stats.GCSys)},
		{"HeapAlloc", gauge(stats.HeapAlloc)},
		{"HeapIdle", gauge(stats.HeapIdle)},
		{"HeapInuse", gauge(stats.HeapInuse)},
		{"HeapObjects", gauge(stats.HeapObjects)},
		{"HeapReleased", gauge(stats.HeapReleased)},
		{"HeapSys", gauge(stats.HeapSys)},
		{"LastGC", gauge(stats.LastGC)},
		{"Lookups", gauge(stats.Lookups)},
		{"MCacheInuse", gauge(stats.MCacheInuse)},
		{"MCacheSys", gauge(stats.MCacheSys)},
		{"MSpanInuse", gauge(stats.MSpanInuse)},
		{"MSpanSys", gauge(stats.MSpanSys)},
		{"Mallocs", gauge(stats.Mallocs)},
		{"NextGC", gauge(stats.NextGC)},
		{"NumForcedGC", gauge(stats.NumForcedGC)},
		{"NumGC", gauge(stats.NumGC)},
		{"OtherSys", gauge(stats.OtherSys)},
		{"PauseTotalNs", gauge(stats.PauseTotalNs)},
		{"StackInuse", gauge(stats.StackInuse)},
		{"StackSys", gauge(stats.StackSys)},
		{"Sys", gauge(stats.Sys)},
		{"TotalAlloc", gauge(stats.TotalAlloc)},
		{"RandomValue", gauge(random)},
	}

	return gaugeMetrics, counterMetrics
}
