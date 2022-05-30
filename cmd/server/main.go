package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const gaugeType = "gauge"
const counterType = "counter"

type gauge float64
type counter int64

type GagueMetric struct {
	Name  string
	Value gauge
}

type CounterMetric struct {
	Name  string
	Value counter
}

type Storage struct {
	Gauges   map[string]gauge
	Counters map[string]counter
}

func (s *Storage) SaveGauge(metric *GagueMetric) error {
	s.Gauges[metric.Name] = metric.Value
	return nil
}

func (s *Storage) SaveCounter(metric CounterMetric) error {
	if _, ok := s.Counters[metric.Name]; ok {
		fmt.Println(ok)
		s.Counters[metric.Name] += counter(metric.Value)
	} else {
		fmt.Println(ok)
		s.Counters[metric.Name] = counter(metric.Value)
	}
	return nil
}

type UpdateParams struct {
	MetricName  string
	MetricType  string
	MetricValue string
}

var storage = Storage{
	Gauges:   map[string]gauge{},
	Counters: map[string]counter{},
}

// MetricHandle — обработчик запроса.
func MetricHandle(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
	// 	return
	// }

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
	case gaugeType:
		value, err := strconv.ParseFloat(params.MetricValue, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
		}
		storage.SaveGauge(&GagueMetric{params.MetricName, gauge(value)})
	case counterType:
		value, err := strconv.ParseInt(params.MetricValue, 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
		}
		storage.SaveCounter(CounterMetric{params.MetricName, counter(value)})
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
