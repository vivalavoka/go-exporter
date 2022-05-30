package main

import (
	"fmt"
	"go-exporter/cmd/server/metrics"
	"net/http"
	"strconv"
	"strings"
)

type Storage struct {
	Gauges   map[string]metrics.Gauge
	Counters map[string]metrics.Counter
}

func (s *Storage) SaveGauge(metric *metrics.GagueMetric) error {
	s.Gauges[metric.Name] = metric.Value
	return nil
}

func (s *Storage) SaveCounter(metric metrics.CounterMetric) error {
	if _, ok := s.Counters[metric.Name]; ok {
		fmt.Println(ok)
		s.Counters[metric.Name] += metrics.Counter(metric.Value)
	} else {
		fmt.Println(ok)
		s.Counters[metric.Name] = metrics.Counter(metric.Value)
	}
	return nil
}

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
}

var storage = Storage{
	Gauges:   map[string]metrics.Gauge{},
	Counters: map[string]metrics.Counter{},
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
	case metrics.GaugeType:
		value, err := strconv.ParseFloat(params.MetricValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
		}
		storage.SaveGauge(&metrics.GagueMetric{params.MetricName, metrics.Gauge(value)})
	case metrics.CounterType:
		value, err := strconv.ParseInt(params.MetricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
		}
		storage.SaveCounter(metrics.CounterMetric{params.MetricName, metrics.Counter(value)})
	default:
		http.Error(w, "Wrong metric type", http.StatusBadRequest)
	}

	fmt.Println(storage.Gauges)
	fmt.Println(storage.Counters)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", MetricHandle)
	// запуск сервера с адресом localhost, порт 8080
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	server.ListenAndServe()
}
