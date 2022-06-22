package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/cmd/server/handlers"
	"github.com/vivalavoka/go-exporter/cmd/server/http/middlewares"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"
)

type Server struct {
	storage *storage.Storage
}

func (s *Server) Start(cfg config.Config) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CompressHandle)
	r.Use(middlewares.DecompressHandle)

	handlers.UpdateMetricRoute(r)
	handlers.GetAllMetricsRoute(r)
	handlers.GetMetricRoute(r)

	http.ListenAndServe(cfg.Address, r)
}
