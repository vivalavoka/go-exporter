package main

import (
	"encoding/json"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type FileDB struct {
	config  Config
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

var fileDB *FileDB

func NewDB(config Config) *FileDB {
	if config.StoreFile == "" {
		return &FileDB{}
	}

	file, err := os.OpenFile(config.StoreFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	fileDB = &FileDB{
		config:  config,
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}
	return fileDB
}

func GetDB() *FileDB {
	return fileDB
}

func (db *FileDB) RunTicker() {
	if db.config.StoreInterval == 0 {
		return
	}

	storage := GetStorage()
	ticker := time.NewTicker(db.config.StoreInterval)
	for {
		select {
		case <-ticker.C:
			metrics := storage.GetMetrics()
			db.Write(metrics)
		}
	}
}

func (db *FileDB) Close() {
	if db.config.StoreFile != "" {
		db.file.Close()
	}
}

func (db *FileDB) Read() (map[string]Metric, error) {
	var metrics []Metric
	metricMap := map[string]Metric{}

	if db.config.StoreFile == "" {
		return metricMap, nil
	}

	if err := db.decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	for _, value := range metrics {
		metricMap[value.ID] = value
	}
	return metricMap, nil
}

func (db *FileDB) Write(metricMap map[string]Metric) error {
	if db.config.StoreFile == "" {
		return nil
	}

	var metrics []Metric
	for _, value := range metricMap {
		metrics = append(metrics, value)
	}

	db.file.Seek(0, 0)
	return db.encoder.Encode(&metrics)
}
