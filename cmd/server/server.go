package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vivalavoka/go-exporter/cmd/server/middlewares"
)

type Server struct {
	storage *Storage
}

func (s *Server) Start(config Config) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipHandle)

	UpdateMetricRoute(r)
	GetAllMetricsRoute(r)
	GetMetricRoute(r)

	http.ListenAndServe(config.Address, r)
}
