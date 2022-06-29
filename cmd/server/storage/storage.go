package storage

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/cmd/server/metrics"
)

type Storage struct {
	config  config.Config
	fileDB  *FileDB
	metrics map[string]metrics.Metric
}

var storage *Storage

func NewStorage(config config.Config) *Storage {
	metrics := map[string]metrics.Metric{}
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

func (s *Storage) GetMetrics() map[string]metrics.Metric {
	return s.metrics
}

func (s *Storage) GetMetric(name string) (metrics.Metric, error) {
	if value, ok := s.metrics[name]; ok {
		return value, nil
	}
	return metrics.Metric{}, fmt.Errorf("there is no metric by name: %s", name)
}

func (s *Storage) Save(metric *metrics.Metric) error {
	value, ok := s.metrics[metric.ID]
	if metric.MType == metrics.CounterType && ok {
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
