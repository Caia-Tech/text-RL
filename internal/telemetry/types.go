package telemetry

import (
	"time"
)

type TelemetryClient struct {
	Endpoint      string
	BufferSize    int
	FlushInterval time.Duration
	MetricsBuffer chan Metric
	
	active bool
	done   chan struct{}
}

type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
	Timestamp time.Time         `json:"timestamp"`
}

type MetricsCollector struct {
	client   *TelemetryClient
	counters map[string]float64
	gauges   map[string]float64
	timers   map[string][]time.Duration
}