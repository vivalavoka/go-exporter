package main

import "fmt"

type Storage struct {
	Gauges   map[string]gauge
	Counters map[string]counter
}

func PrintStorage() {
	fmt.Println(storage)
}

func GetGaugeMetrics() map[string]gauge {
	return storage.Gauges
}

func GetCounterMetrics() map[string]counter {
	return storage.Counters
}

func GetMetricGauge(name string) (gauge, error) {
	if value, ok := storage.Gauges[name]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("there is no metric by name: %s", name)
}

func GetMetricCounter(name string) (counter, error) {
	if value, ok := storage.Counters[name]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("there is no metric by name: %s", name)
}

func SaveGauge(name string, value gauge) error {
	storage.Gauges[name] = value
	return nil
}

func SaveCounter(name string, value counter) error {
	storage.Counters[name] += value
	return nil
}

var storage = Storage{
	Gauges:   map[string]gauge{},
	Counters: map[string]counter{},
}
