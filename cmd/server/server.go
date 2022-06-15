package main

import (
	"net/http"

	"github.com/vivalavoka/go-exporter/cmd/server/handlers"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	storage *storage.Storage
}

func (s *Server) Start(config Config) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.UpdateMetricRoute(r)
	handlers.GetAllMetricsRoute(r)
	handlers.GetMetricRoute(r)

	http.ListenAndServe(config.Address, r)
}
