package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/vivalavoka/go-exporter/cmd/server/storage"

	"github.com/go-chi/chi/v5"
)

const GaugeType = "gauge"
const CounterType = "counter"

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
}

type Metrics struct {
	ID    string           `json:"id"`              // имя метрики
	MType string           `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *storage.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *storage.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricData struct {
	Name  string
	Value string
}

type MetricsPageData struct {
	PageTitle string
	Metrics   []MetricData
}

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	repo := storage.GetStorage()
	tmpl := template.Must(template.ParseFiles("layouts/metrics.html"))
	data := MetricsPageData{
		PageTitle: "Exporter metrics",
	}

	gauges := repo.GetGaugeMetrics()
	for name, value := range gauges {
		data.Metrics = append(data.Metrics, MetricData{name, fmt.Sprintf("%.3f", value)})
	}

	counters := repo.GetCounterMetrics()
	for name, value := range counters {
		data.Metrics = append(data.Metrics, MetricData{name, fmt.Sprintf("%d", value)})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	repo := storage.GetStorage()
	params := UpdateParams{
		MetricType: chi.URLParam(r, "type"),
		MetricName: chi.URLParam(r, "name"),
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch params.MetricType {
	case GaugeType:
		value, err := repo.GetMetricGauge(params.MetricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%.3f", value)))
	case CounterType:
		value, err := repo.GetMetricCounter(params.MetricName)
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

func GetMetricFromBody(w http.ResponseWriter, r *http.Request) {
	repo := storage.GetStorage()
	var metric Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch metric.MType {
	case GaugeType:
		value, err := repo.GetMetricGauge(metric.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		metric.Value = &value
		response, err := json.Marshal(metric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	case CounterType:
		value, err := repo.GetMetricCounter(metric.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		metric.Delta = &value
		response, err := json.Marshal(metric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
}

// MetricHandle — обработчик запроса.
func MetricHandle(w http.ResponseWriter, r *http.Request) {
	repo := storage.GetStorage()
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
		repo.SaveGauge(params.MetricName, storage.Gauge(value))
	case CounterType:
		value, err := strconv.ParseInt(params.MetricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		repo.SaveCounter(params.MetricName, storage.Counter(value))
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func MetricHandleFromBody(w http.ResponseWriter, r *http.Request) {
	repo := storage.GetStorage()
	metric := Metrics{}
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case GaugeType:
		if metric.Value == nil {
			var v storage.Gauge
			metric.Value = &v
		}
		repo.SaveGauge(metric.ID, storage.Gauge(*metric.Value))
	case CounterType:
		if metric.Delta == nil {
			var v storage.Counter
			metric.Delta = &v
		}
		repo.SaveCounter(metric.ID, storage.Counter(*metric.Delta))
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
