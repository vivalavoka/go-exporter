package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	http *http.Server
}

func (s *Server) Start() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", MetricHandle)
	// запуск сервера с адресом localhost, порт 8080
	s.http = &http.Server{
		Addr: "127.0.0.1:8080",
	}
	fmt.Printf("Server running on %s\n", s.http.Addr)
	s.http.ListenAndServe()
}
