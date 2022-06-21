package main

import "fmt"

type Storage struct {
	metrics map[string]Metric
}

var singleInstance *Storage

func GetStorage() *Storage {
	if singleInstance == nil {
		singleInstance = &Storage{
			metrics: map[string]Metric{},
		}
	}

	return singleInstance
}

func (s *Storage) GetMetrics() map[string]Metric {
	return s.metrics
}

func (s *Storage) GetMetric(name string) (Metric, error) {
	if value, ok := s.metrics[name]; ok {
		return value, nil
	}
	return Metric{}, fmt.Errorf("there is no metric by name: %s", name)
}

func (s *Storage) Save(metric Metric) error {
	value, ok := s.metrics[metric.ID]
	if metric.MType == CounterType && ok {
		var delta Counter
		delta = *metric.Delta + *value.Delta
		value.Delta = &delta
	} else {
		s.metrics[metric.ID] = metric
	}
	return nil
}
