package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vivalavoka/go-exporter/cmd/server/config"
	server "github.com/vivalavoka/go-exporter/cmd/server/http"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.New(config)
	if err != nil {
		log.Fatal(err)
	}

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
