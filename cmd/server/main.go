package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vivalavoka/go-exporter/cmd/server/config"
	server "github.com/vivalavoka/go-exporter/cmd/server/http"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func main() {
	config := config.Get()
	storage := storage.NewStorage(config)
	defer storage.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		storage.DropCache()
		os.Exit(1)
	}()

	http := server.Server{}
	http.Start(config)
}
