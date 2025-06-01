package handler

import (
	"encoding/json"
	"net/http"
	"telemetry-collector/internal/store"
)

func EventHandler(s *store.TelemetryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := s.GetExecEvents()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
