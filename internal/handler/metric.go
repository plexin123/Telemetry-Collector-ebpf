package handler

import (
	"encoding/json"
	"net/http"
	"telemetry-collector/internal/store"
	"time"
)

func MetricHandler(s *store.TelemetryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var m store.MetricPoint
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if m.Timestamp == 0 {
			m.Timestamp = time.Now().Unix()
		}
		s.Add(m)
		w.WriteHeader(http.StatusAccepted)

	}
}

func StatsHandler(s *store.TelemetryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := s.Stats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
