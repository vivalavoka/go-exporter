package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type gauge float64
type counter int64

const GaugeType = "gauge"
const CounterType = "counter"

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
}

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	body := ""

	gauges := GetGaugeMetrics()
	for name, value := range gauges {
		body += fmt.Sprintf("<strong>%s:</strong> %.3f</br>", name, value)
	}

	counters := GetCounterMetrics()
	for name, value := range counters {
		body += fmt.Sprintf("<strong>%s:</strong> %d</br>", name, value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	params := UpdateParams{
		MetricType: chi.URLParam(r, "type"),
		MetricName: chi.URLParam(r, "name"),
	}

	switch params.MetricType {
	case GaugeType:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		value, err := GetMetricGauge(params.MetricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%.3f", value)))
	case CounterType:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		value, err := GetMetricCounter(params.MetricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", value)))
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
}

// MetricHandle — обработчик запроса.
func MetricHandle(w http.ResponseWriter, r *http.Request) {
	params := UpdateParams{
		MetricType:  chi.URLParam(r, "type"),
		MetricName:  chi.URLParam(r, "name"),
		MetricValue: chi.URLParam(r, "value"),
	}

	switch params.MetricType {
	case GaugeType:
		value, err := strconv.ParseFloat(params.MetricValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		SaveGauge(params.MetricName, gauge(value))
	case CounterType:
		value, err := strconv.ParseInt(params.MetricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		SaveCounter(params.MetricName, counter(value))
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
