package main

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

var config Config

func main() {
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(config)

	rand.Seed(time.Now().Unix())

	client := NewClient(resty.New())
	agent := NewAgent(client)

	agent.Start()
}
