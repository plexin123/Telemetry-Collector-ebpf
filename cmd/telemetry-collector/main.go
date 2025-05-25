package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

// for structures we use:
type MetricPoint struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type Stats struct {
	Count int     `json:"count"`
	Sum   float64 `json:"sum"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
}

//in memory store basically a map of values

type TelemetryStore struct {
	mu sync.RWMutex
	//a hashmap of string as key and values as [] of MeitricPoints
	buffer map[string][]MetricPoint
}

func NewTelemetryStore() *TelemetryStore {
	return &TelemetryStore{
		buffer: make(map[string][]MetricPoint),
	}
}

func (s *TelemetryStore) Add(m MetricPoint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buffer[m.Name] = append(s.buffer[m.Name], m)
}

func (s *TelemetryStore) Stats() map[string]Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]Stats)
	for name, metrics := range s.buffer {
		var sum float64
		var max float64 = math.Inf(-1)
		for _, m := range metrics {
			sum = sum + m.Value
			if m.Value > max {
				max = m.Value
			}
		}
		count := len(metrics)
		result[name] = Stats{
			Count: count,
			Sum:   sum,
			Max:   max,
			Avg:   sum / float64(count),
		}
	}
	return result
}

func main() {
	store := NewTelemetryStore()
	//http handlers, basically the endpoints, we are going to use the methods to understand tht
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var m MetricPoint
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if m.Timestamp == 0 {
			m.Timestamp = time.Now().Unix()

		}
		store.Add(m)
		w.WriteHeader(http.StatusAccepted)

	})
	http.HandleFunc("/metrics/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := store.Stats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
