package store

import (
	"log"
	"math"
	"sync"
	"time"
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

type ExecPoint struct {
	Name      string `json:"name"`
	TimeStamp int64  `json:"timestamp"`
	PID       uint32 `json:"pid"`
	UID       uint32 `json:"uid"`
}

type TelemetryStore struct {
	mu     sync.RWMutex
	buffer map[string][]MetricPoint
	execMu sync.Mutex
	execs  []ExecPoint
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

/*Measures events per second /metrics/rate*/

func (s *TelemetryStore) Rate(windowsSeconds int64) map[string]float64 {
	now := time.Now().Unix()
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]float64)
	for name, points := range s.buffer {
		var count int
		for _, p := range points {
			if now-p.Timestamp <= windowsSeconds {
				count++
			}
		}
		result[name] = float64(count) / float64(windowsSeconds)
	}
	return result
}

func (s *TelemetryStore) StartFlushing(ttl time.Duration, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.mu.Lock()
			now := time.Now().Unix()
			for name, points := range s.buffer {
				var filtered []MetricPoint
				for _, p := range points {
					if now-p.Timestamp <= int64(ttl.Seconds()) {
						filtered = append(filtered, p)
					}
				}
				s.buffer[name] = filtered
			}
			s.mu.Unlock()
			log.Println("[TTL] Expired old metric points")
		}
	}()
}

func (s *TelemetryStore) AddExec(e ExecPoint) {
	s.execMu.Lock()
	defer s.execMu.Unlock()
	s.execs = append(s.execs, e)
}

func (s *TelemetryStore) GetExecEvents() []ExecPoint {
	s.execMu.Lock()
	defer s.execMu.Unlock()
	events := make([]ExecPoint, len(s.execs))
	copy(events, s.execs)
	return events
}

func (s *TelemetryStore) StartExecTTL(ttl time.Duration, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			now := time.Now().Unix()
			s.execMu.Lock()
			var fresh []ExecPoint
			for _, e := range s.execs {
				if now-e.TimeStamp <= int64(ttl.Seconds()) {
					fresh = append(fresh, e)
				}
			}
			s.execs = fresh
			s.execMu.Unlock()
		}
	}()
}
