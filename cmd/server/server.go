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
	// r.Use(middleware.Compress(5))
	r.Use(middlewares.CompressHandle)
	r.Use(middlewares.DecompressHandle)

	UpdateMetricRoute(r)
	GetAllMetricsRoute(r)
	GetMetricRoute(r)

	http.ListenAndServe(config.Address, r)
}
