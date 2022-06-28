package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Storage struct {
	config  Config
	fileDB  *FileDB
	metrics map[string]Metric
}

var storage *Storage

func NewStorage(config Config) *Storage {
	metrics := map[string]Metric{}
	fileDB := NewDB(config)

	if config.Restore {
		_metrics, err := fileDB.Read()
		if err != nil {
			log.Error(err)
		} else {
			metrics = _metrics
		}
	}

	storage = &Storage{
		config:  config,
		fileDB:  fileDB,
		metrics: metrics,
	}

	if config.StoreInterval != 0 {
		go fileDB.RunTicker()
	}

	return storage
}

func GetStorage() *Storage {
	return storage
}

func (s *Storage) DropCache() {
	if err := s.fileDB.Write(s.metrics); err != nil {
		log.Error(err)
	}
}

func (s *Storage) Close() {
	s.fileDB.Close()
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

func (s *Storage) Save(metric *Metric) error {
	value, ok := s.metrics[metric.ID]
	if metric.MType == CounterType && ok {
		*metric.Delta += *value.Delta
	}
	s.metrics[metric.ID] = *metric

	if s.config.StoreInterval == 0 {
		if err := s.fileDB.Write(s.metrics); err != nil {
			log.Error(err)
		}
	}

	return nil
}
