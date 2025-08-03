package rl

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Intelligent caching system that learns optimal caching strategies
type IntelligentCache struct {
	cache     map[string]CacheEntry
	metadata  map[string]CacheMetadata
	mutex     sync.RWMutex
	
	// Cache optimization parameters
	maxSize          int
	ttl              time.Duration
	hitRateThreshold float64
	
	// Learning components
	accessPatterns   map[string][]time.Time
	costBenefitModel map[string]CostBenefit
	
	// Statistics
	hits     int64
	misses   int64
	evictions int64
}

type CacheEntry struct {
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"`
	Timestamp time.Time             `json:"timestamp"`
	AccessCount int                 `json:"access_count"`
	LastAccess  time.Time           `json:"last_access"`
	ComputeCost time.Duration       `json:"compute_cost"`
	MemorySize  int64               `json:"memory_size"`
}

type CacheMetadata struct {
	FunctionName    string        `json:"function_name"`
	Parameters      string        `json:"parameters"`
	TextHash        string        `json:"text_hash"`
	CreationTime    time.Time     `json:"creation_time"`
	HitCount        int           `json:"hit_count"`
	AverageLatency  time.Duration `json:"average_latency"`
	EvictionRisk    float64       `json:"eviction_risk"`
}

type CostBenefit struct {
	ComputeCost    time.Duration `json:"compute_cost"`
	StorageCost    int64         `json:"storage_cost"`
	HitProbability float64       `json:"hit_probability"`
	Benefit        float64       `json:"benefit"`
}

func NewIntelligentCache(maxSize int, ttl time.Duration) *IntelligentCache {
	return &IntelligentCache{
		cache:            make(map[string]CacheEntry),
		metadata:         make(map[string]CacheMetadata),
		maxSize:          maxSize,
		ttl:              ttl,
		hitRateThreshold: 0.3,
		accessPatterns:   make(map[string][]time.Time),
		costBenefitModel: make(map[string]CostBenefit),
	}
}

// Generate cache key with intelligent hashing
func (ic *IntelligentCache) generateCacheKey(functionName string, text string, params map[string]interface{}) string {
	// Create text hash for large texts
	textHash := ic.hashText(text)
	
	// Serialize parameters
	paramBytes, _ := json.Marshal(params)
	
	// Combine into cache key
	key := fmt.Sprintf("%s:%s:%s", functionName, textHash, string(paramBytes))
	return key
}

func (ic *IntelligentCache) hashText(text string) string {
	if len(text) <= 100 {
		return text // Small texts, use directly
	}
	
	// For large texts, use hash + length + sample
	hash := sha256.Sum256([]byte(text))
	sample := ""
	if len(text) > 200 {
		sample = text[:100] + text[len(text)-100:]
	} else {
		sample = text
	}
	
	return fmt.Sprintf("%x:%d:%s", hash[:8], len(text), sample)
}

// Smart cache lookup with learning
func (ic *IntelligentCache) Get(functionName string, text string, params map[string]interface{}) (interface{}, bool) {
	ic.mutex.RLock()
	defer ic.mutex.RUnlock()
	
	key := ic.generateCacheKey(functionName, text, params)
	
	// Record access pattern
	ic.recordAccess(key)
	
	entry, exists := ic.cache[key]
	if !exists {
		ic.misses++
		return nil, false
	}
	
	// Check TTL
	if time.Since(entry.Timestamp) > ic.ttl {
		ic.misses++
		// Don't delete here to avoid lock issues
		return nil, false
	}
	
	// Update access statistics
	entry.AccessCount++
	entry.LastAccess = time.Now()
	ic.cache[key] = entry
	
	ic.hits++
	return entry.Value, true
}

// Smart cache storage with eviction learning
func (ic *IntelligentCache) Set(functionName string, text string, params map[string]interface{}, 
	value interface{}, computeCost time.Duration) {
	
	ic.mutex.Lock()
	defer ic.mutex.Unlock()
	
	key := ic.generateCacheKey(functionName, text, params)
	
	// Calculate memory size estimate
	memorySize := ic.estimateMemorySize(value)
	
	// Check if we need to evict
	if len(ic.cache) >= ic.maxSize {
		ic.intelligentEviction()
	}
	
	// Create cache entry
	entry := CacheEntry{
		Key:         key,
		Value:       value,
		Timestamp:   time.Now(),
		AccessCount: 1,
		LastAccess:  time.Now(),
		ComputeCost: computeCost,
		MemorySize:  memorySize,
	}
	
	// Store entry and metadata
	ic.cache[key] = entry
	ic.metadata[key] = CacheMetadata{
		FunctionName:   functionName,
		Parameters:     fmt.Sprintf("%v", params),
		TextHash:       ic.hashText(text),
		CreationTime:   time.Now(),
		HitCount:       0,
		AverageLatency: computeCost,
		EvictionRisk:   ic.calculateEvictionRisk(entry),
	}
	
	// Update cost-benefit model
	ic.updateCostBenefit(functionName, computeCost, memorySize)
}

// Intelligent eviction using learned patterns
func (ic *IntelligentCache) intelligentEviction() {
	if len(ic.cache) == 0 {
		return
	}
	
	// Calculate eviction scores for all entries
	scores := make(map[string]float64)
	
	for key, entry := range ic.cache {
		metadata := ic.metadata[key]
		score := ic.calculateEvictionScore(entry, metadata)
		scores[key] = score
	}
	
	// Find entry with lowest score (highest eviction priority)
	var evictKey string
	minScore := float64(1e9)
	
	for key, score := range scores {
		if score < minScore {
			minScore = score
			evictKey = key
		}
	}
	
	// Evict the selected entry
	if evictKey != "" {
		delete(ic.cache, evictKey)
		delete(ic.metadata, evictKey)
		ic.evictions++
	}
}

func (ic *IntelligentCache) calculateEvictionScore(entry CacheEntry, metadata CacheMetadata) float64 {
	now := time.Now()
	
	// Time-based factors
	age := now.Sub(entry.Timestamp).Seconds()
	timeSinceLastAccess := now.Sub(entry.LastAccess).Seconds()
	
	// Access pattern factors
	accessFrequency := float64(entry.AccessCount) / (age + 1)
	
	// Cost-benefit factors
	computeCostSeconds := entry.ComputeCost.Seconds()
	memoryCostMB := float64(entry.MemorySize) / (1024 * 1024)
	
	// Predicted future access probability
	accessProbability := ic.predictAccessProbability(metadata.FunctionName, metadata.TextHash)
	
	// Combined score (higher = keep, lower = evict)
	score := (accessFrequency * computeCostSeconds * accessProbability) / (memoryCostMB * (timeSinceLastAccess + 1))
	
	return score
}

func (ic *IntelligentCache) calculateEvictionRisk(entry CacheEntry) float64 {
	// Simple eviction risk calculation
	age := time.Since(entry.Timestamp).Seconds()
	timeSinceAccess := time.Since(entry.LastAccess).Seconds()
	
	risk := (age + timeSinceAccess) / (float64(entry.AccessCount) + 1)
	return risk
}

func (ic *IntelligentCache) recordAccess(key string) {
	now := time.Now()
	if _, exists := ic.accessPatterns[key]; !exists {
		ic.accessPatterns[key] = make([]time.Time, 0)
	}
	
	ic.accessPatterns[key] = append(ic.accessPatterns[key], now)
	
	// Keep only recent access patterns (last 100 accesses)
	if len(ic.accessPatterns[key]) > 100 {
		ic.accessPatterns[key] = ic.accessPatterns[key][1:]
	}
}

func (ic *IntelligentCache) predictAccessProbability(functionName, textHash string) float64 {
	// Simple prediction based on historical patterns
	patterns, exists := ic.accessPatterns[fmt.Sprintf("%s:%s", functionName, textHash)]
	if !exists || len(patterns) == 0 {
		return 0.1 // Low default probability
	}
	
	now := time.Now()
	recentAccesses := 0
	
	// Count accesses in the last hour
	for _, accessTime := range patterns {
		if now.Sub(accessTime) < time.Hour {
			recentAccesses++
		}
	}
	
	// Convert to probability (0-1)
	probability := float64(recentAccesses) / 10.0
	if probability > 1.0 {
		probability = 1.0
	}
	
	return probability
}

func (ic *IntelligentCache) updateCostBenefit(functionName string, computeCost time.Duration, memorySize int64) {
	if _, exists := ic.costBenefitModel[functionName]; !exists {
		ic.costBenefitModel[functionName] = CostBenefit{}
	}
	
	model := ic.costBenefitModel[functionName]
	
	// Update average compute cost
	model.ComputeCost = (model.ComputeCost + computeCost) / 2
	model.StorageCost = (model.StorageCost + memorySize) / 2
	
	// Update hit probability based on function usage
	hitRate := ic.getHitRate()
	model.HitProbability = (model.HitProbability + hitRate) / 2
	
	// Calculate benefit score
	if model.StorageCost > 0 {
		model.Benefit = (model.ComputeCost.Seconds() * model.HitProbability) / float64(model.StorageCost)
	}
	
	ic.costBenefitModel[functionName] = model
}

func (ic *IntelligentCache) estimateMemorySize(value interface{}) int64 {
	// Simple memory estimation
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return 1024 // Default estimate
	}
	return int64(len(jsonBytes) * 2) // Estimate with overhead
}

func (ic *IntelligentCache) getHitRate() float64 {
	total := ic.hits + ic.misses
	if total == 0 {
		return 0.0
	}
	return float64(ic.hits) / float64(total)
}

// Should cache decision using learned patterns
func (ic *IntelligentCache) ShouldCache(functionName string, text string, computeCost time.Duration) bool {
	ic.mutex.RLock()
	defer ic.mutex.RUnlock()
	
	// Get cost-benefit model for this function
	model, exists := ic.costBenefitModel[functionName]
	if !exists {
		// Default: cache if compute cost is high
		return computeCost > time.Millisecond*100
	}
	
	// Cache if benefit score is above threshold
	return model.Benefit > 0.1
}

// Get cache statistics for optimization
func (ic *IntelligentCache) GetStats() map[string]interface{} {
	ic.mutex.RLock()
	defer ic.mutex.RUnlock()
	
	hitRate := ic.getHitRate()
	
	stats := map[string]interface{}{
		"hits":             ic.hits,
		"misses":           ic.misses,
		"evictions":        ic.evictions,
		"hit_rate":         hitRate,
		"cache_size":       len(ic.cache),
		"max_size":         ic.maxSize,
		"cost_benefit_models": ic.costBenefitModel,
	}
	
	return stats
}

// Clean expired entries
func (ic *IntelligentCache) CleanExpired() int {
	ic.mutex.Lock()
	defer ic.mutex.Unlock()
	
	cleaned := 0
	now := time.Now()
	
	for key, entry := range ic.cache {
		if now.Sub(entry.Timestamp) > ic.ttl {
			delete(ic.cache, key)
			delete(ic.metadata, key)
			cleaned++
		}
	}
	
	return cleaned
}