package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	SHAKey        string        `env:"KEY"`
}

func Get() (Config, error) {
	var config Config
	flag.StringVar(&config.Address, "a", "127.0.0.1:8080", "server address")
	flag.DurationVar(&config.StoreInterval, "i", time.Duration(300*time.Millisecond), "store interval")
	flag.StringVar(&config.StoreFile, "f", "/tmp/devops-metrics-db.json", "store file name")
	flag.BoolVar(&config.Restore, "r", true, "need restore")
	flag.StringVar(&config.SHAKey, "k", "", "sha key")
	flag.Parse()

	err := env.Parse(&config)
	if err != nil {
		return config, err
	}

	log.Info(config)
	return config, nil
}
