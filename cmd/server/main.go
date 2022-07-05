package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vivalavoka/go-exporter/cmd/server/config"
	server "github.com/vivalavoka/go-exporter/cmd/server/http"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"
)

func main() {
	config := config.Get()
	storage, _ := storage.New(config)
	defer storage.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		storage.Repo.Close()
		os.Exit(1)
	}()

	http := server.Server{}
	http.Start(config)
}
