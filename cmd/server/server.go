package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct{}

func (s *Server) Start() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	UpdateMetricRoute(r)
	GetAllMetricsRoute(r)
	GetMetricRoute(r)

	http.ListenAndServe(":8080", r)
}
