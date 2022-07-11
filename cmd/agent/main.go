package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vivalavoka/go-exporter/cmd/agent/client"
	"github.com/vivalavoka/go-exporter/cmd/agent/config"
	"github.com/vivalavoka/go-exporter/cmd/agent/exporter"
)

func main() {
	rand.Seed(time.Now().Unix())

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	client := client.New(cfg.Address, resty.New())

	exporter := exporter.New(cfg, client)
	exporter.Start()
}
