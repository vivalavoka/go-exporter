package memorydb

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type InMemoryDB struct {
	config    config.Config
	fileUse   bool
	asyncFile bool
	file      *os.File
	encoder   *json.Encoder
	decoder   *json.Decoder
	metrics   map[string]metrics.Metric
}

func New(config config.Config) (*InMemoryDB, error) {
	metrics := map[string]metrics.Metric{}

	repo := &InMemoryDB{
		config:    config,
		metrics:   metrics,
		fileUse:   config.StoreFile != "",
		asyncFile: config.StoreInterval != 0,
	}

	if repo.fileUse {
		file, err := os.OpenFile(config.StoreFile, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
		repo.file = file
		repo.encoder = json.NewEncoder(file)
		repo.decoder = json.NewDecoder(file)

		if config.Restore {
			metricList, err := repo.readFile()
			if err != nil {
				log.Warn(err)
			} else {
				metrics = metricList
			}
		}
	}

	if repo.asyncFile {
		go repo.runTicker()
	}

	return repo, nil
}

func (r *InMemoryDB) runTicker() {
	ticker := time.NewTicker(r.config.StoreInterval)
	for range ticker.C {
		r.writeFile(r.metrics)
	}
}

func (r *InMemoryDB) Close() {
	if r.fileUse {
		r.dropCache()
		r.file.Close()
	}
}

func (r *InMemoryDB) GetMetrics() (map[string]metrics.Metric, error) {
	return r.metrics, nil
}

func (r *InMemoryDB) GetMetric(name string) (metrics.Metric, error) {
	if value, ok := r.metrics[name]; ok {
		return value, nil
	}
	return metrics.Metric{}, fmt.Errorf("there is no metric by name: %s", name)
}

func (r *InMemoryDB) Save(metric *metrics.Metric) error {
	value, ok := r.metrics[metric.ID]
	if metric.MType == metrics.CounterType && ok {
		metric.Delta += value.Delta
	}
	r.metrics[metric.ID] = *metric

	if !r.asyncFile {
		if err := r.writeFile(r.metrics); err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (r *InMemoryDB) dropCache() {
	if err := r.writeFile(r.metrics); err != nil {
		log.Error(err)
	}
}

func (r *InMemoryDB) readFile() (map[string]metrics.Metric, error) {
	var metricList []metrics.Metric
	metricMap := map[string]metrics.Metric{}

	if err := r.decoder.Decode(&metricList); err != nil {
		return nil, err
	}

	println(metricList)

	for _, value := range metricList {
		metricMap[value.ID] = value
	}
	return metricMap, nil
}

func (r *InMemoryDB) writeFile(metricMap map[string]metrics.Metric) error {
	if !r.fileUse {
		return nil
	}

	var metricList []metrics.Metric
	for _, value := range metricMap {
		metricList = append(metricList, value)
	}

	r.file.Seek(0, 0)
	return r.encoder.Encode(&metricList)
}
