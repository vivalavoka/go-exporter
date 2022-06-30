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

func New(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) Start(cfg config.Config) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CompressHandle)
	r.Use(middlewares.DecompressHandle)

	h := handlers.New(cfg)
	h.SetRoutes(r)

	http.ListenAndServe(cfg.Address, r)
}
