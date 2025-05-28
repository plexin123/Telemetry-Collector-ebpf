package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"telemetry-collector/internal/store"
)

func RateHandler(s *store.TelemetryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		window := int64(10)
		if param := r.URL.Query().Get("window"); param != "" {
			if parsed, err := strconv.ParseInt(param, 10, 64); err == nil {
				window = parsed
			}
		}
		rates := s.Rate(window)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rates)
	}
}
