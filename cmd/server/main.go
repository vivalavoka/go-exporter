package main

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func main() {
	var config Config

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(config)

	server := Server{}

	server.Start(config)
}
