package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"
	"github.com/vivalavoka/go-exporter/internal/crypto"
	"github.com/vivalavoka/go-exporter/internal/metrics"
)

type Handlers struct {
	hasher  *crypto.SHA256
	storage *storage.Storage
}

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

func New(cfg config.Config, storage *storage.Storage) *Handlers {
	hasher := crypto.New(cfg.SHAKey)
	return &Handlers{
		hasher:  hasher,
		storage: storage,
	}
}

func (h *Handlers) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	ex, err := os.Executable()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	layoutPath := path.Join(filepath.Dir(ex), "handlers/layouts/metrics.html")
	tmpl := template.Must(template.ParseFiles(layoutPath))
	data := MetricsPageData{
		PageTitle: "Exporter metrics",
	}

	metricList, err := h.storage.Repo.GetMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	for name, value := range metricList {
		if value.MType == metrics.GaugeType {
			data.Metrics = append(data.Metrics, MetricData{name, fmt.Sprintf("%.3f", *value.Value)})
		} else {
			data.Metrics = append(data.Metrics, MetricData{name, fmt.Sprintf("%d", *value.Delta)})
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func (h *Handlers) CheckConnection(w http.ResponseWriter, r *http.Request) {
	ok := h.storage.Repo.CheckConnection()
	if !ok {
		http.Error(w, "Connection refused", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *Handlers) GetMetric(w http.ResponseWriter, r *http.Request) {
	params := UpdateParams{
		MetricType: chi.URLParam(r, "type"),
		MetricName: chi.URLParam(r, "name"),
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch params.MetricType {
	case metrics.GaugeType:
		value, err := h.storage.Repo.GetMetric(params.MetricName, params.MetricType)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%.3f", *value.Value)))
	case metrics.CounterType:
		value, err := h.storage.Repo.GetMetric(params.MetricName, params.MetricType)
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

func (h *Handlers) GetMetricFromBody(w http.ResponseWriter, r *http.Request) {
	var params metrics.Metric

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if params.MType != metrics.CounterType && params.MType != metrics.GaugeType {
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	metric, err := h.storage.Repo.GetMetric(params.ID, params.MType)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if h.hasher.Enable {
		metric.Hash = h.hasher.GetSum(metric.String())
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
func (h *Handlers) MetricHandle(w http.ResponseWriter, r *http.Request) {
	params := UpdateParams{
		MetricType:  chi.URLParam(r, "type"),
		MetricName:  chi.URLParam(r, "name"),
		MetricValue: chi.URLParam(r, "value"),
	}

	switch params.MetricType {
	case metrics.GaugeType:
		value, err := strconv.ParseFloat(params.MetricValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		h.storage.Repo.Save(&metrics.Metric{
			ID:    params.MetricName,
			MType: params.MetricType,
			Value: (*metrics.Gauge)(&value),
		})
	case metrics.CounterType:
		value, err := strconv.ParseInt(params.MetricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		h.storage.Repo.Save(&metrics.Metric{
			ID:    params.MetricName,
			MType: params.MetricType,
			Delta: (*metrics.Counter)(&value),
		})
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (h *Handlers) MetricHandleFromBody(w http.ResponseWriter, r *http.Request) {
	var params *metrics.Metric

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch params.MType {
	case metrics.GaugeType:
		if params.Value == nil {
			var v metrics.Gauge
			params.Value = &v
		}
	case metrics.CounterType:
		if params.Delta == nil {
			var v metrics.Counter
			params.Delta = &v
		}
	default:
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	if h.hasher.Enable {
		hash := h.hasher.GetSum(params.String())
		if hash != params.Hash {
			http.Error(w, "Wrong hash", http.StatusBadRequest)
			return
		}
	}

	err = h.storage.Repo.Save(params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *Handlers) MetricBatchHandle(w http.ResponseWriter, r *http.Request) {

	var params []*metrics.Metric

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, param := range params {
		switch param.MType {
		case metrics.GaugeType:
			if param.Value == nil {
				var v metrics.Gauge
				param.Value = &v
			}
		case metrics.CounterType:
			if param.Delta == nil {
				var v metrics.Counter
				param.Delta = &v
			}
		default:
			http.Error(w, "Wrong metric type", http.StatusNotImplemented)
			return
		}

		if h.hasher.Enable {
			hash := h.hasher.GetSum(param.String())
			if hash != param.Hash {
				http.Error(w, "Wrong hash", http.StatusBadRequest)
				return
			}
		}
	}

	h.storage.Repo.SaveBatch(params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
