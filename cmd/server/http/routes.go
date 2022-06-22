package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vivalavoka/go-exporter/cmd/server/handlers"
)

func UpdateMetricRoute(r chi.Router) chi.Router {
	r.Post("/update/{type}/{name}/{value}", handlers.MetricHandle)
	r.Post("/update/", handlers.MetricHandleFromBody)
	return r
}

func GetAllMetricsRoute(r chi.Router) chi.Router {
	r.Get("/", handlers.GetAllMetrics)
	return r
}

func GetMetricRoute(r chi.Router) chi.Router {
	r.Get("/value/{type}/{name}", handlers.GetMetric)
	r.Post("/value/", handlers.GetMetricFromBody)
	return r
}
