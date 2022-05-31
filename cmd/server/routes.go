package main

import "github.com/go-chi/chi/v5"

func UpdateMetricRoute(r chi.Router) chi.Router {
	r.Post("/update/{type}/{name}/{value}", MetricHandle)
	return r
}

func GetAllMetricsRoute(r chi.Router) chi.Router {
	r.Get("/", GetAllMetrics)
	return r
}

func GetMetricRoute(r chi.Router) chi.Router {
	r.Get("/value/{type}/{name}", GetMetric)
	return r
}
