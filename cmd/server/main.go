package main

import (
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300ms"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	var config Config

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(config)

	storage := NewStorage(config)
	defer storage.Close()

	server := Server{}
	server.Start(config)
}
