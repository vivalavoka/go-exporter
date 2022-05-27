package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	reportAddress  = "localhost:8000"
)

type gauge float64
type counter int64

type GaugeItem struct {
	Name  string
	Value gauge
}

type CounterItem struct {
	Name  string
	Value counter
}

func GetMetrics(pollCount counter) ([]GaugeItem, []CounterItem) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

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
		{"RandomValue", gauge(random.Float64())},
	}

	return gaugeMetrics, counterMetrics
}

func ReportMetrics(gaugeMetrics []GaugeItem, counterMetrics []CounterItem) {
	for _, item := range gaugeMetrics {
		fmt.Printf("Name: %s, Type: %s, Value: %v\n", item.Name, reflect.TypeOf(item.Value), item.Value)
	}

	for _, item := range counterMetrics {
		fmt.Printf("Name: %s, Type: %s, Value: %v\n", item.Name, reflect.TypeOf(item.Value), item.Value)
	}
}

func main() {
	pollCount := counter(0)
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	var gaugeMetrics []GaugeItem
	var counterMetrics []CounterItem
	for {
		select {
		case <-reportTicker.C:
			ReportMetrics(gaugeMetrics, counterMetrics)
		case <-pollTicker.C:
			pollCount = pollCount + 1
			gaugeMetrics, counterMetrics = GetMetrics(pollCount)
		}
	}
}
