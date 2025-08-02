package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"textlib-rl-system/internal/logging"
)

func NewTelemetryClient(endpoint string, bufferSize int, flushInterval time.Duration) *TelemetryClient {
	return &TelemetryClient{
		Endpoint:      endpoint,
		BufferSize:    bufferSize,
		FlushInterval: flushInterval,
		MetricsBuffer: make(chan Metric, bufferSize),
		done:          make(chan struct{}),
	}
}

func (tc *TelemetryClient) Start() error {
	if tc.active {
		return fmt.Errorf("telemetry client already active")
	}
	
	tc.active = true
	go tc.processMetrics()
	
	log.Println("Telemetry client started")
	return nil
}

func (tc *TelemetryClient) Stop() {
	if !tc.active {
		return
	}
	
	tc.active = false
	close(tc.done)
	log.Println("Telemetry client stopped")
}

func (tc *TelemetryClient) RecordFunctionCall(functionName string, success bool, duration time.Duration) {
	if !tc.active {
		return
	}
	
	// Record function call count
	tc.sendMetric(Metric{
		Name:  "textlib.function.call",
		Value: 1,
		Tags: map[string]string{
			"function": functionName,
			"success":  strconv.FormatBool(success),
		},
		Timestamp: time.Now(),
	})
	
	// Record function duration
	tc.sendMetric(Metric{
		Name:  "textlib.function.duration",
		Value: duration.Seconds(),
		Tags: map[string]string{
			"function": functionName,
		},
		Timestamp: time.Now(),
	})
}

func (tc *TelemetryClient) RecordLearningMetrics(metrics logging.LearningMetrics) {
	if !tc.active {
		return
	}
	
	tc.sendMetric(Metric{
		Name:      "rl.q_value.convergence",
		Value:     metrics.QValueConvergence,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.exploration.rate",
		Value:     metrics.ExplorationRate,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.policy.stability",
		Value:     metrics.PolicyStability,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.action.diversity",
		Value:     metrics.ActionDiversity,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.learning.progress",
		Value:     metrics.LearningProgress,
		Timestamp: time.Now(),
	})
}

func (tc *TelemetryClient) RecordPerformanceMetrics(metrics logging.PerformanceMetrics) {
	if !tc.active {
		return
	}
	
	tc.sendMetric(Metric{
		Name:      "rl.performance.cumulative_reward",
		Value:     metrics.CumulativeReward,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.performance.average_reward",
		Value:     metrics.AverageReward,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.performance.success_rate",
		Value:     metrics.SuccessRate,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.performance.efficiency_score",
		Value:     metrics.EfficiencyScore,
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:      "rl.performance.task_completion_rate",
		Value:     metrics.TaskCompletionRate,
		Timestamp: time.Now(),
	})
}

func (tc *TelemetryClient) RecordEpisodeMetrics(episodeID string, totalReward float64, stepCount int, duration time.Duration) {
	if !tc.active {
		return
	}
	
	tc.sendMetric(Metric{
		Name:  "rl.episode.total_reward",
		Value: totalReward,
		Tags: map[string]string{
			"episode_id": episodeID,
		},
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:  "rl.episode.step_count",
		Value: float64(stepCount),
		Tags: map[string]string{
			"episode_id": episodeID,
		},
		Timestamp: time.Now(),
	})
	
	tc.sendMetric(Metric{
		Name:  "rl.episode.duration",
		Value: duration.Seconds(),
		Tags: map[string]string{
			"episode_id": episodeID,
		},
		Timestamp: time.Now(),
	})
}

func (tc *TelemetryClient) sendMetric(metric Metric) {
	select {
	case tc.MetricsBuffer <- metric:
	default:
		log.Printf("Warning: metrics buffer full, dropping metric: %s", metric.Name)
	}
}

func (tc *TelemetryClient) processMetrics() {
	ticker := time.NewTicker(tc.FlushInterval)
	defer ticker.Stop()
	
	var batch []Metric
	
	for {
		select {
		case metric := <-tc.MetricsBuffer:
			batch = append(batch, metric)
			
			if len(batch) >= 100 { // Flush when batch is full
				tc.flushBatch(batch)
				batch = nil
			}
			
		case <-ticker.C:
			if len(batch) > 0 {
				tc.flushBatch(batch)
				batch = nil
			}
			
		case <-tc.done:
			if len(batch) > 0 {
				tc.flushBatch(batch)
			}
			return
		}
	}
}

func (tc *TelemetryClient) flushBatch(batch []Metric) {
	if tc.Endpoint == "" {
		// If no endpoint configured, just log locally
		log.Printf("Telemetry batch: %d metrics", len(batch))
		return
	}
	
	data, err := json.Marshal(batch)
	if err != nil {
		log.Printf("Error marshaling metrics batch: %v", err)
		return
	}
	
	resp, err := http.Post(tc.Endpoint+"/api/v1/metrics", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending metrics to %s: %v", tc.Endpoint, err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("Metrics endpoint returned status %d", resp.StatusCode)
	}
}

func NewMetricsCollector(client *TelemetryClient) *MetricsCollector {
	return &MetricsCollector{
		client:   client,
		counters: make(map[string]float64),
		gauges:   make(map[string]float64),
		timers:   make(map[string][]time.Duration),
	}
}

func (mc *MetricsCollector) IncrementCounter(name string, tags map[string]string) {
	mc.counters[name]++
	mc.client.sendMetric(Metric{
		Name:      name,
		Value:     1,
		Tags:      tags,
		Timestamp: time.Now(),
	})
}

func (mc *MetricsCollector) SetGauge(name string, value float64, tags map[string]string) {
	mc.gauges[name] = value
	mc.client.sendMetric(Metric{
		Name:      name,
		Value:     value,
		Tags:      tags,
		Timestamp: time.Now(),
	})
}

func (mc *MetricsCollector) RecordTimer(name string, duration time.Duration, tags map[string]string) {
	mc.timers[name] = append(mc.timers[name], duration)
	mc.client.sendMetric(Metric{
		Name:      name,
		Value:     duration.Seconds(),
		Tags:      tags,
		Timestamp: time.Now(),
	})
}

func (mc *MetricsCollector) GetSummary() map[string]interface{} {
	summary := make(map[string]interface{})
	
	summary["counters"] = mc.counters
	summary["gauges"] = mc.gauges
	
	timerSummary := make(map[string]map[string]float64)
	for name, durations := range mc.timers {
		if len(durations) == 0 {
			continue
		}
		
		var total time.Duration
		min := durations[0]
		max := durations[0]
		
		for _, d := range durations {
			total += d
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}
		
		avg := total / time.Duration(len(durations))
		
		timerSummary[name] = map[string]float64{
			"count": float64(len(durations)),
			"min":   min.Seconds(),
			"max":   max.Seconds(),
			"avg":   avg.Seconds(),
			"total": total.Seconds(),
		}
	}
	
	summary["timers"] = timerSummary
	summary["timestamp"] = time.Now()
	
	return summary
}