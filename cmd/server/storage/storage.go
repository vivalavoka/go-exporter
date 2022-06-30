package storage

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type Storage struct {
	config  config.Config
	fileDB  *FileDB
	metrics map[string]metrics.Metric
}

var storage *Storage

func New(config config.Config) (*Storage, error) {
	metrics := map[string]metrics.Metric{}
	fileDB, err := NewDB(config)

	if err != nil {
		return nil, err
	}

	if config.Restore {
		metricList, err := fileDB.Read()
		if err != nil {
			log.Error(err)
		} else {
			metrics = metricList
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

	return storage, nil
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
		metric.Delta += value.Delta
	}
	s.metrics[metric.ID] = *metric

	if s.config.StoreInterval == 0 {
		if err := s.fileDB.Write(s.metrics); err != nil {
			log.Error(err)
		}
	}

	return nil
}
