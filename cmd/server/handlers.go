package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
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
	repo := GetStorage()
	tmpl := template.Must(template.ParseFiles("layouts/metrics.html"))
	data := MetricsPageData{
		PageTitle: "Exporter metrics",
	}

	metrics := repo.GetMetrics()
	for name, value := range metrics {
		if value.MType == GaugeType {
			data.Metrics = append(data.Metrics, MetricData{name, fmt.Sprintf("%.3f", *value.Value)})
		} else {
			data.Metrics = append(data.Metrics, MetricData{name, fmt.Sprintf("%d", *value.Delta)})
		}
	}
	log.Info(w.Header().Get("Content-Type"))
	log.Info(data)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	repo := GetStorage()
	params := UpdateParams{
		MetricType: chi.URLParam(r, "type"),
		MetricName: chi.URLParam(r, "name"),
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch params.MetricType {
	case GaugeType:
		value, err := repo.GetMetric(params.MetricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%.3f", *value.Value)))
	case CounterType:
		value, err := repo.GetMetric(params.MetricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", *value.Delta)))
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
}

func GetMetricFromBody(w http.ResponseWriter, r *http.Request) {
	repo := GetStorage()
	var params Metric

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if params.MType != CounterType && params.MType != GaugeType {
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	metric, err := repo.GetMetric(params.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// MetricHandle — обработчик запроса.
func MetricHandle(w http.ResponseWriter, r *http.Request) {
	repo := GetStorage()
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
		repo.Save(&Metric{
			ID:    params.MetricName,
			MType: params.MetricType,
			Value: (*Gauge)(&value),
		})
	case CounterType:
		value, err := strconv.ParseInt(params.MetricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		repo.Save(&Metric{
			ID:    params.MetricName,
			MType: params.MetricType,
			Delta: (*Counter)(&value),
		})
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func MetricHandleFromBody(w http.ResponseWriter, r *http.Request) {
	repo := GetStorage()

	var params *Metric

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch params.MType {
	case GaugeType:
		if params.Value == nil {
			var v Gauge
			params.Value = &v
		}
	case CounterType:
		if params.Delta == nil {
			var v Counter
			params.Delta = &v
		}
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
	repo.Save(params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
