package main

import (
	"log"
	"net/http"
	"telemetry-collector/internal/handler"
	"telemetry-collector/internal/store"
	"time"
)

func main() {
	s := store.NewTelemetryStore()
	http.HandleFunc("/metrics", handler.MetricHandler(s))
	http.HandleFunc("/metrics/stats", handler.StatsHandler(s))
	http.HandleFunc("/health", handler.HealthHandler)
	http.HandleFunc("/metrics/rate", handler.RateHandler(s))
	http.HandleFunc("/events", handler.EventHandler(s))
	handler.AddMockExecEvents(s)
	s.StartExecTTL(30*time.Second, 5*time.Second)
	s.StartFlushing(30*time.Second, 5*time.Second)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
