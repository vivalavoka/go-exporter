package main

import (
	"flag"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

var config Config

func main() {
	flag.StringVar(&config.Address, "a", "127.0.0.1:8080", "server address")
	flag.DurationVar(&config.ReportInterval, "r", time.Duration(10*time.Second), "report interval")
	flag.DurationVar(&config.PollInterval, "p", time.Duration(2*time.Second), "poll interval")
	flag.Parse()

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
