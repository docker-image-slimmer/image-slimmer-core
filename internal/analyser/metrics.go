package analyzer

import (
	"sync"
	"time"
)

// Metrics represents immutable execution metrics snapshot.
// It is safe to expose externally.
type Metrics struct {
	FetchDuration time.Duration
	BuildDuration time.Duration
	TotalDuration time.Duration

	FetchAttempts int
	DigestPinned  bool
	Success       bool
}

// metricsCollector accumulates execution timings internally.
// It is not exposed outside the analyzer boundary.
type metricsCollector struct {
	mu sync.Mutex

	start time.Time

	fetchStart time.Time
	buildStart time.Time

	fetchDuration time.Duration
	buildDuration time.Duration

	fetchAttempts int
	digestPinned  bool
	success       bool
}

// newMetricsCollector initializes timing collection.
func newMetricsCollector() *metricsCollector {
	return &metricsCollector{
		start: time.Now(),
	}
}

// ---- FETCH PHASE ----

func (m *metricsCollector) startFetch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fetchStart = time.Now()
}

func (m *metricsCollector) endFetch(attempts int, digestPinned bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.fetchStart.IsZero() {
		m.fetchDuration += time.Since(m.fetchStart)
	}

	m.fetchAttempts = attempts
	m.digestPinned = digestPinned
}

// ---- BUILD PHASE ----

func (m *metricsCollector) startBuild() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.buildStart = time.Now()
}

func (m *metricsCollector) endBuild() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.buildStart.IsZero() {
		m.buildDuration += time.Since(m.buildStart)
	}
}

// ---- FINALIZATION ----

// markSuccess marks the execution result.
func (m *metricsCollector) markSuccess(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.success = success
}

// snapshot produces an immutable Metrics struct.
// TotalDuration is always computed from start time.
func (m *metricsCollector) snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Metrics{
		FetchDuration: m.fetchDuration,
		BuildDuration: m.buildDuration,
		TotalDuration: time.Since(m.start),
		FetchAttempts: m.fetchAttempts,
		DigestPinned:  m.digestPinned,
		Success:       m.success,
	}
}
