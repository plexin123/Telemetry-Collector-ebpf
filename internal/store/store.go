package store

import (
	"math"
	"sync"
)

/* Basically we are trying to architecture the design of the app according to different parts, in this
case we have a file store
*/

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

type TelemetryStore struct {
	mu     sync.RWMutex
	buffer map[string][]MetricPoint
}

/* Basically we need to create a func to initialize a TelemetryStore" */

func NewTelemetryStore() *TelemetryStore {
	return &TelemetryStore{
		buffer: make(map[string][]MetricPoint),
	}
}

func (s *TelemetryStore) Add(m MetricPoint) {
	s.mu.Lock()
	/*this will execute at the end */
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
