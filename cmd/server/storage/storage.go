package storage

import "fmt"

type Gauge float64
type Counter int64

type Storage struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

var singleInstance *Storage

func GetStorage() *Storage {
	if singleInstance == nil {
		singleInstance = &Storage{
			Gauges:   map[string]Gauge{},
			Counters: map[string]Counter{},
		}
	}

	return singleInstance
}

func (s *Storage) GetGaugeMetrics() map[string]Gauge {
	return s.Gauges
}

func (s *Storage) GetCounterMetrics() map[string]Counter {
	return s.Counters
}

func (s *Storage) GetMetricGauge(name string) (Gauge, error) {
	if value, ok := s.Gauges[name]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("there is no metric by name: %s", name)
}

func (s *Storage) GetMetricCounter(name string) (Counter, error) {
	if value, ok := s.Counters[name]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("there is no metric by name: %s", name)
}

func (s *Storage) SaveGauge(name string, value Gauge) error {
	s.Gauges[name] = value
	return nil
}

func (s *Storage) SaveCounter(name string, value Counter) error {
	s.Counters[name] += value
	return nil
}
