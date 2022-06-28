package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func main() {
	var config Config
	flag.StringVar(&config.Address, "a", "127.0.0.1:8080", "server address")
	flag.DurationVar(&config.StoreInterval, "i", time.Duration(300*time.Millisecond), "store interval")
	flag.StringVar(&config.StoreFile, "f", "/tmp/devops-metrics-db.json", "store file name")
	flag.BoolVar(&config.Restore, "r", true, "need restore")
	flag.Parse()

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(config)

	storage := NewStorage(config)
	defer storage.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		storage.DropCache()
		os.Exit(1)
	}()

	server := Server{}
	server.Start(config)
}
