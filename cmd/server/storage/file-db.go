package storage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/cmd/server/metrics"
)

type FileDB struct {
	config  config.Config
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

var fileDB *FileDB

func NewDB(config config.Config) (*FileDB, error) {
	if config.StoreFile == "" {
		return &FileDB{}, nil
	}

	file, err := os.OpenFile(config.StoreFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	fileDB = &FileDB{
		config:  config,
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}
	return fileDB, nil
}

func GetDB() *FileDB {
	return fileDB
}

func (db *FileDB) RunTicker() {
	if db.config.StoreInterval == 0 {
		return
	}

	ticker := time.NewTicker(db.config.StoreInterval)
	for range ticker.C {
		storage := GetStorage()
		metrics := storage.GetMetrics()
		db.Write(metrics)
	}
}

func (db *FileDB) Close() {
	if db.config.StoreFile != "" {
		db.file.Close()
	}
}

func (db *FileDB) Read() (map[string]metrics.Metric, error) {
	var metricList []metrics.Metric
	metricMap := map[string]metrics.Metric{}

	if db.config.StoreFile == "" {
		return metricMap, nil
	}

	if err := db.decoder.Decode(&metricList); err != nil {
		return nil, err
	}

	for _, value := range metricList {
		metricMap[value.ID] = value
	}
	return metricMap, nil
}

func (db *FileDB) Write(metricMap map[string]metrics.Metric) error {
	if db.config.StoreFile == "" {
		return nil
	}

	var metricList []metrics.Metric
	for _, value := range metricMap {
		metricList = append(metricList, value)
	}

	db.file.Seek(0, 0)
	return db.encoder.Encode(&metricList)
}
