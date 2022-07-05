package handlers

import "github.com/go-chi/chi/v5"

func (h *Handlers) SetRoutes(r chi.Router) chi.Router {
	r.Get("/", h.GetAllMetrics)
	r.Get("/ping", h.CheckConnection)
	r.Post("/update/{type}/{name}/{value}", h.MetricHandle)
	r.Post("/update/", h.MetricHandleFromBody)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Post("/value/", h.GetMetricFromBody)

	return r
}
