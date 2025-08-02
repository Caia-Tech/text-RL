package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func NewInsightLogger(logPath string, batchSize int, flushInterval time.Duration) *InsightLogger {
	return &InsightLogger{
		LogPath:       logPath,
		MetricsDB:     NewMetricsDatabase(),
		EventStream:   make(chan LogEvent, 1000),
		BatchSize:     batchSize,
		FlushInterval: flushInterval,
		batch:         make([]LogEvent, 0, batchSize),
		done:          make(chan struct{}),
	}
}

func NewMetricsDatabase() *MetricsDatabase {
	return &MetricsDatabase{
		events:   make([]LogEvent, 0),
		episodes: make(map[string]EpisodeMetrics),
	}
}

func (il *InsightLogger) Start() error {
	il.mu.Lock()
	defer il.mu.Unlock()
	
	if il.active {
		return fmt.Errorf("logger already active")
	}
	
	if err := os.MkdirAll(il.LogPath, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	
	il.active = true
	go il.processEvents()
	
	return nil
}

func (il *InsightLogger) Stop() {
	il.mu.Lock()
	defer il.mu.Unlock()
	
	if !il.active {
		return
	}
	
	il.active = false
	close(il.done)
	il.flushBatch()
}

func (il *InsightLogger) StartSession(sessionID string) {
	il.mu.Lock()
	defer il.mu.Unlock()
	
	il.sessionID = sessionID
	log.Printf("Started logging session: %s", sessionID)
}

func (il *InsightLogger) EndSession() {
	il.mu.Lock()
	defer il.mu.Unlock()
	
	il.flushBatch()
	log.Printf("Ended logging session: %s", il.sessionID)
}

func (il *InsightLogger) LogEvent(event LogEvent) {
	if !il.active {
		return
	}
	
	select {
	case il.EventStream <- event:
	default:
		log.Printf("Warning: event stream buffer full, dropping event")
	}
}

func (il *InsightLogger) LogEpisodeSummary(metrics EpisodeMetrics) {
	il.MetricsDB.mu.Lock()
	defer il.MetricsDB.mu.Unlock()
	
	il.MetricsDB.episodes[metrics.EpisodeID] = metrics
	
	filename := filepath.Join(il.LogPath, fmt.Sprintf("episode_%s.json", metrics.EpisodeID))
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		log.Printf("Error marshaling episode metrics: %v", err)
		return
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Error writing episode summary: %v", err)
	}
}

func (il *InsightLogger) LogInsights(insights interface{}) {
	filename := filepath.Join(il.LogPath, "insights.json")
	data, err := json.MarshalIndent(insights, "", "  ")
	if err != nil {
		log.Printf("Error marshaling insights: %v", err)
		return
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Error writing insights: %v", err)
	}
}

func (il *InsightLogger) processEvents() {
	ticker := time.NewTicker(il.FlushInterval)
	defer ticker.Stop()
	
	for {
		select {
		case event := <-il.EventStream:
			il.addEventToBatch(event)
			
		case <-ticker.C:
			il.flushBatch()
			
		case <-il.done:
			return
		}
	}
}

func (il *InsightLogger) addEventToBatch(event LogEvent) {
	il.mu.Lock()
	defer il.mu.Unlock()
	
	il.batch = append(il.batch, event)
	il.MetricsDB.events = append(il.MetricsDB.events, event)
	
	if len(il.batch) >= il.BatchSize {
		il.flushBatch()
	}
}

func (il *InsightLogger) flushBatch() {
	if len(il.batch) == 0 {
		return
	}
	
	filename := filepath.Join(il.LogPath, fmt.Sprintf("events_%d.json", time.Now().Unix()))
	data, err := json.MarshalIndent(il.batch, "", "  ")
	if err != nil {
		log.Printf("Error marshaling batch: %v", err)
		return
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Error writing batch: %v", err)
		return
	}
	
	il.batch = il.batch[:0]
}

func (md *MetricsDatabase) GetEvents() []LogEvent {
	md.mu.RLock()
	defer md.mu.RUnlock()
	
	events := make([]LogEvent, len(md.events))
	copy(events, md.events)
	return events
}

func (md *MetricsDatabase) GetEpisodes() map[string]EpisodeMetrics {
	md.mu.RLock()
	defer md.mu.RUnlock()
	
	episodes := make(map[string]EpisodeMetrics)
	for k, v := range md.episodes {
		episodes[k] = v
	}
	return episodes
}

func (md *MetricsDatabase) GetEventsByType(eventType string) []LogEvent {
	md.mu.RLock()
	defer md.mu.RUnlock()
	
	var filtered []LogEvent
	for _, event := range md.events {
		if event.EventType == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func (md *MetricsDatabase) GetEventsByEpisode(episodeID string) []LogEvent {
	md.mu.RLock()
	defer md.mu.RUnlock()
	
	var filtered []LogEvent
	for _, event := range md.events {
		if event.EpisodeID == episodeID {
			filtered = append(filtered, event)
		}
	}
	return filtered
}