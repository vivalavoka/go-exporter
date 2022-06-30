package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	SHAKey         string        `env:"KEY"`
}

func Get() (Config, error) {
	var config Config
	flag.StringVar(&config.Address, "a", "127.0.0.1:8080", "server address")
	flag.DurationVar(&config.ReportInterval, "r", time.Duration(10*time.Second), "report interval")
	flag.DurationVar(&config.PollInterval, "p", time.Duration(2*time.Second), "poll interval")
	flag.StringVar(&config.SHAKey, "k", "", "sha256 key")
	flag.Parse()

	err := env.Parse(&config)
	if err != nil {
		return config, err
	}

	log.Info(config)
	return config, nil
}
