package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type FileDB struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

var fileDB *FileDB

func NewDB(config Config) *FileDB {
	file, err := os.OpenFile(config.StoreFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	fileDB = &FileDB{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}
	return fileDB
}

func GetDB() *FileDB {
	return fileDB
}

func (db *FileDB) Close() {
	db.file.Close()
}

func (db *FileDB) Read() (map[string]Metric, error) {
	var metrics []Metric

	if err := db.decoder.Decode(&metrics); err != nil {
		return nil, err
	}

	metricMap := map[string]Metric{}
	for _, value := range metrics {
		metricMap[value.ID] = value
	}
	return metricMap, nil
}

func (db *FileDB) SyncWrite(metricMap map[string]Metric) error {
	var metrics []Metric
	for _, value := range metricMap {
		metrics = append(metrics, value)
	}

	db.file.Seek(0, 0)
	return db.encoder.Encode(&metrics)
}
