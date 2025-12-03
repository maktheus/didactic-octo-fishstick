package metrics

import (
	"sort"
	"strings"
	"sync"
)

// Recorder exposes primitives for capturing telemetry used by the services.
type Recorder interface {
	AddCounter(name string, labels map[string]string, delta float64)
	ObserveHistogram(name string, labels map[string]string, value float64)
	Snapshot() Snapshot
}

// Snapshot holds a thread-safe copy of the captured metrics.
type Snapshot struct {
	Counters   map[string]float64
	Histograms map[string][]float64
}

// InMemoryRecorder stores metrics in-memory for inspection during tests.
type InMemoryRecorder struct {
	mu         sync.RWMutex
	counters   map[string]float64
	histograms map[string][]float64
}

// NewInMemory creates a recorder instance backed by local state only.
func NewInMemory() *InMemoryRecorder {
	return &InMemoryRecorder{
		counters:   make(map[string]float64),
		histograms: make(map[string][]float64),
	}
}

// AddCounter records the provided delta under the counter name and labels.
func (r *InMemoryRecorder) AddCounter(name string, labels map[string]string, delta float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := formatKey(name, labels)
	r.counters[key] += delta
}

// ObserveHistogram records the provided value in the histogram bucket.
func (r *InMemoryRecorder) ObserveHistogram(name string, labels map[string]string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := formatKey(name, labels)
	r.histograms[key] = append(r.histograms[key], value)
}

// Snapshot returns a deep copy of the aggregated metrics to avoid data races.
func (r *InMemoryRecorder) Snapshot() Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counters := make(map[string]float64, len(r.counters))
	for k, v := range r.counters {
		counters[k] = v
	}
	histograms := make(map[string][]float64, len(r.histograms))
	for k, values := range r.histograms {
		copied := make([]float64, len(values))
		copy(copied, values)
		histograms[k] = copied
	}
	return Snapshot{Counters: counters, Histograms: histograms}
}

func formatKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	// deterministic ordering
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b := strings.Builder{}
	b.WriteString(name)
	b.WriteString("{")
	for i, k := range keys {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(labels[k])
	}
	b.WriteString("}")
	return b.String()
}
