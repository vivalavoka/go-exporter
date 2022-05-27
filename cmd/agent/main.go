package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	reportAddress  = "127.0.0.1:8080"
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

func SplitType(valueType string) string {
	return strings.Split(valueType, ".")[1]
}

func GetMetrics(pollCount counter) (*[]GaugeItem, *[]CounterItem) {
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

	return &gaugeMetrics, &counterMetrics
}

func MakeRequest(client *http.Client, url string) {
	data := url
	request, err := http.NewRequest(http.MethodPost, "https://webhook.site/ef75d15f-48dd-48b5-a4be-0a98a4633099", bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(data)))

	client.Do(request)
}

func ReportMetrics(client *http.Client, gaugeMetrics *[]GaugeItem, counterMetrics *[]CounterItem) {
	for _, item := range *gaugeMetrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", reportAddress, SplitType(reflect.TypeOf(item.Value).String()), item.Name, item.Value)
		MakeRequest(client, url)
	}

	for _, item := range *counterMetrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", reportAddress, SplitType(reflect.TypeOf(item.Value).String()), item.Name, item.Value)
		MakeRequest(client, url)
	}
}

func main() {
	client := &http.Client{}

	pollCount := counter(0)
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	var gaugeMetrics *[]GaugeItem
	var counterMetrics *[]CounterItem

	for {
		select {
		case <-reportTicker.C:
			ReportMetrics(client, gaugeMetrics, counterMetrics)
		case <-pollTicker.C:
			pollCount = pollCount + 1
			gaugeMetrics, counterMetrics = GetMetrics(pollCount)
		}
	}
}
