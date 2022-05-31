package main

import "fmt"

type Storage struct {
	Gauges   map[string]gauge
	Counters map[string]counter
}

func PrintStorage() {
	fmt.Println(storage)
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
