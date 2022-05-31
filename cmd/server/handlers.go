package main

import (
	"net/http"
	"strconv"
	"strings"
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

// MetricHandle — обработчик запроса.
func MetricHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	args := strings.Split(r.URL.Path, "/")

	if len(args) != 5 {
		http.Error(w, "Wrong request path", http.StatusBadRequest)
		return
	}

	if args[1] != "update" {
		http.Error(w, "Wrong request path", http.StatusBadRequest)
		return
	}

	params := UpdateParams{
		MetricType:  args[2],
		MetricName:  args[3],
		MetricValue: args[4],
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
		http.Error(w, "Wrong metric type", http.StatusBadRequest)
		return
	}

	PrintStorage()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
