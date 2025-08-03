package rl

import (
	"strings"
	"testing"
	"time"
)

func TestNewIntelligentCache(t *testing.T) {
	maxSize := 100
	ttl := time.Hour
	
	cache := NewIntelligentCache(maxSize, ttl)
	
	if cache == nil {
		t.Fatal("NewIntelligentCache returned nil")
	}
	
	if cache.maxSize != maxSize {
		t.Errorf("Expected maxSize %d, got %d", maxSize, cache.maxSize)
	}
	
	if cache.ttl != ttl {
		t.Errorf("Expected ttl %v, got %v", ttl, cache.ttl)
	}
	
	if cache.hitRateThreshold != 0.3 {
		t.Errorf("Expected hitRateThreshold 0.3, got %f", cache.hitRateThreshold)
	}
	
	if cache.cache == nil {
		t.Error("cache map not initialized")
	}
	
	if cache.metadata == nil {
		t.Error("metadata map not initialized")
	}
	
	if cache.accessPatterns == nil {
		t.Error("accessPatterns map not initialized")
	}
	
	if cache.costBenefitModel == nil {
		t.Error("costBenefitModel map not initialized")
	}
}

func TestIntelligentCache_hashText(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	tests := []struct {
		name        string
		text        string
		expectDirect bool
	}{
		{
			name:        "short text",
			text:        "Hello world",
			expectDirect: true,
		},
		{
			name:        "medium text",
			text:        strings.Repeat("a", 150),
			expectDirect: false,
		},
		{
			name:        "long text",
			text:        strings.Repeat("b", 1000),
			expectDirect: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := cache.hashText(tt.text)
			
			if tt.expectDirect {
				if hash != tt.text {
					t.Errorf("Expected direct text for short input, got hash: %s", hash)
				}
			} else {
				if hash == tt.text {
					t.Error("Expected hash for long text, got direct text")
				}
				
				// Hash should contain length and sample information
				if !strings.Contains(hash, ":") {
					t.Error("Hash should contain delimiters for structured format")
				}
			}
		})
	}
}

func TestIntelligentCache_generateCacheKey(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	functionName := "TestFunction"
	text := "Hello world"
	params := map[string]interface{}{
		"param1": "value1",
		"param2": 42,
	}
	
	key1 := cache.generateCacheKey(functionName, text, params)
	key2 := cache.generateCacheKey(functionName, text, params)
	
	// Same inputs should generate same key
	if key1 != key2 {
		t.Errorf("Same inputs generated different keys: %s vs %s", key1, key2)
	}
	
	// Different function should generate different key
	key3 := cache.generateCacheKey("DifferentFunction", text, params)
	if key1 == key3 {
		t.Error("Different function names generated same key")
	}
	
	// Different text should generate different key
	key4 := cache.generateCacheKey(functionName, "Different text", params)
	if key1 == key4 {
		t.Error("Different text generated same key")
	}
	
	// Different params should generate different key
	params2 := map[string]interface{}{
		"param1": "different_value",
		"param2": 42,
	}
	key5 := cache.generateCacheKey(functionName, text, params2)
	if key1 == key5 {
		t.Error("Different params generated same key")
	}
}

func TestIntelligentCache_SetAndGet(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	functionName := "TestFunction"
	text := "Hello world"
	params := map[string]interface{}{"param1": "value1"}
	value := "test result"
	computeCost := time.Millisecond * 100
	
	// Test Set
	cache.Set(functionName, text, params, value, computeCost)
	
	// Test Get
	result, found := cache.Get(functionName, text, params)
	
	if !found {
		t.Error("Expected to find cached value")
	}
	
	if result != value {
		t.Errorf("Expected %v, got %v", value, result)
	}
	
	// Verify cache entry was created
	key := cache.generateCacheKey(functionName, text, params)
	entry, exists := cache.cache[key]
	if !exists {
		t.Error("Cache entry not created")
	}
	
	if entry.Value != value {
		t.Errorf("Cache entry value: expected %v, got %v", value, entry.Value)
	}
	
	if entry.ComputeCost != computeCost {
		t.Errorf("Cache entry compute cost: expected %v, got %v", computeCost, entry.ComputeCost)
	}
	
	// Verify metadata was created
	metadata, exists := cache.metadata[key]
	if !exists {
		t.Error("Cache metadata not created")
	}
	
	if metadata.FunctionName != functionName {
		t.Errorf("Metadata function name: expected %s, got %s", functionName, metadata.FunctionName)
	}
}

func TestIntelligentCache_TTLExpiration(t *testing.T) {
	cache := NewIntelligentCache(100, time.Millisecond*10) // Very short TTL
	
	functionName := "TestFunction"
	text := "Hello world"
	params := map[string]interface{}{"param1": "value1"}
	value := "test result"
	computeCost := time.Millisecond * 100
	
	// Set value
	cache.Set(functionName, text, params, value, computeCost)
	
	// Should be found immediately
	result, found := cache.Get(functionName, text, params)
	if !found {
		t.Error("Expected to find cached value immediately after set")
	}
	if result != value {
		t.Errorf("Expected %v, got %v", value, result)
	}
	
	// Wait for TTL to expire
	time.Sleep(time.Millisecond * 15)
	
	// Should not be found after TTL expires
	_, found = cache.Get(functionName, text, params)
	if found {
		t.Error("Expected cached value to expire after TTL")
	}
}

func TestIntelligentCache_CacheMiss(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	functionName := "TestFunction"
	text := "Hello world"
	params := map[string]interface{}{"param1": "value1"}
	
	// Get non-existent value
	result, found := cache.Get(functionName, text, params)
	
	if found {
		t.Error("Expected cache miss for non-existent value")
	}
	
	if result != nil {
		t.Errorf("Expected nil result for cache miss, got %v", result)
	}
}

func TestIntelligentCache_intelligentEviction(t *testing.T) {
	cache := NewIntelligentCache(3, time.Hour) // Small cache for testing eviction
	
	// Fill cache to capacity
	for i := 0; i < 3; i++ {
		functionName := "TestFunction"
		text := "Hello world"
		params := map[string]interface{}{"param1": i}
		value := i
		computeCost := time.Millisecond * time.Duration(100+i*50) // Varying costs
		
		cache.Set(functionName, text, params, value, computeCost)
	}
	
	// Verify cache is full
	if len(cache.cache) != 3 {
		t.Errorf("Expected cache size 3, got %d", len(cache.cache))
	}
	
	// Add one more item to trigger eviction
	cache.Set("TestFunction", "Hello world", map[string]interface{}{"param1": 999}, 999, time.Millisecond*100)
	
	// Cache should still be at max size
	if len(cache.cache) != 3 {
		t.Errorf("Expected cache size 3 after eviction, got %d", len(cache.cache))
	}
	
	// New item should be in cache
	result, found := cache.Get("TestFunction", "Hello world", map[string]interface{}{"param1": 999})
	if !found {
		t.Error("Expected to find newly added item after eviction")
	}
	if result != 999 {
		t.Errorf("Expected 999, got %v", result)
	}
}

func TestIntelligentCache_calculateEvictionScore(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	now := time.Now()
	
	// Test entry with high access count and recent access
	entry1 := CacheEntry{
		Key:         "key1",
		Value:       "value1",
		Timestamp:   now.Add(-time.Minute * 10),
		AccessCount: 10,
		LastAccess:  now.Add(-time.Minute * 1),
		ComputeCost: time.Millisecond * 500,
		MemorySize:  1024,
	}
	
	metadata1 := CacheMetadata{
		FunctionName: "TestFunction",
		TextHash:     "hash1",
		CreationTime: now.Add(-time.Minute * 10),
		HitCount:     10,
	}
	
	// Test entry with low access count and old access
	entry2 := CacheEntry{
		Key:         "key2",
		Value:       "value2",
		Timestamp:   now.Add(-time.Hour * 2),
		AccessCount: 1,
		LastAccess:  now.Add(-time.Hour * 1),
		ComputeCost: time.Millisecond * 100,
		MemorySize:  2048,
	}
	
	metadata2 := CacheMetadata{
		FunctionName: "TestFunction",
		TextHash:     "hash2",
		CreationTime: now.Add(-time.Hour * 2),
		HitCount:     1,
	}
	
	score1 := cache.calculateEvictionScore(entry1, metadata1)
	score2 := cache.calculateEvictionScore(entry2, metadata2)
	
	// Entry1 should have higher score (less likely to be evicted)
	if score1 <= score2 {
		t.Errorf("Entry1 should have higher eviction score than entry2: %f vs %f", score1, score2)
	}
}

func TestIntelligentCache_estimateMemorySize(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	tests := []struct {
		name     string
		value    interface{}
		minSize  int64
	}{
		{
			name:    "string value",
			value:   "hello world",
			minSize: 10, // Should be at least a few bytes
		},
		{
			name:    "integer value",
			value:   42,
			minSize: 2,
		},
		{
			name: "map value",
			value: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
			minSize: 20,
		},
		{
			name:    "slice value",
			value:   []string{"a", "b", "c"},
			minSize: 10,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := cache.estimateMemorySize(tt.value)
			
			if size < tt.minSize {
				t.Errorf("Estimated size %d is too small, expected at least %d", size, tt.minSize)
			}
			
			if size <= 0 {
				t.Errorf("Estimated size should be positive, got %d", size)
			}
		})
	}
}

func TestIntelligentCache_ShouldCache(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	tests := []struct {
		name         string
		functionName string
		text         string
		computeCost  time.Duration
		setupModel   bool
		expectedResult bool
	}{
		{
			name:         "high cost no model",
			functionName: "ExpensiveFunction",
			text:         "test text",
			computeCost:  time.Millisecond * 200,
			setupModel:   false,
			expectedResult: true, // Default: cache if cost > 100ms
		},
		{
			name:         "low cost no model",
			functionName: "CheapFunction",
			text:         "test text",
			computeCost:  time.Millisecond * 50,
			setupModel:   false,
			expectedResult: false, // Default: don't cache if cost <= 100ms
		},
		{
			name:         "with beneficial model",
			functionName: "TestFunction",
			text:         "test text",
			computeCost:  time.Millisecond * 150,
			setupModel:   true,
			expectedResult: true, // Model says benefit > 0.1
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupModel {
				// Setup a cost-benefit model with good benefit
				cache.costBenefitModel[tt.functionName] = CostBenefit{
					ComputeCost:    tt.computeCost,
					StorageCost:    1024,
					HitProbability: 0.8,
					Benefit:        0.2, // Above threshold of 0.1
				}
			}
			
			result := cache.ShouldCache(tt.functionName, tt.text, tt.computeCost)
			
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestIntelligentCache_CleanExpired(t *testing.T) {
	cache := NewIntelligentCache(100, time.Millisecond*10) // Very short TTL
	
	// Add some entries
	for i := 0; i < 5; i++ {
		functionName := "TestFunction"
		text := "Hello world"
		params := map[string]interface{}{"param1": i}
		value := i
		computeCost := time.Millisecond * 100
		
		cache.Set(functionName, text, params, value, computeCost)
	}
	
	// Verify all entries are in cache
	if len(cache.cache) != 5 {
		t.Errorf("Expected 5 cache entries, got %d", len(cache.cache))
	}
	
	// Wait for TTL to expire
	time.Sleep(time.Millisecond * 15)
	
	// Clean expired entries
	cleaned := cache.CleanExpired()
	
	if cleaned != 5 {
		t.Errorf("Expected to clean 5 entries, cleaned %d", cleaned)
	}
	
	if len(cache.cache) != 0 {
		t.Errorf("Expected empty cache after cleaning, got %d entries", len(cache.cache))
	}
	
	if len(cache.metadata) != 0 {
		t.Errorf("Expected empty metadata after cleaning, got %d entries", len(cache.metadata))
	}
}

func TestIntelligentCache_GetStats(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	// Generate some cache activity
	functionName := "TestFunction"
	text := "Hello world"
	params := map[string]interface{}{"param1": "value1"}
	value := "test result"
	computeCost := time.Millisecond * 100
	
	// Set and get to generate hits/misses
	cache.Set(functionName, text, params, value, computeCost)
	cache.Get(functionName, text, params) // Hit
	cache.Get(functionName, text, map[string]interface{}{"param1": "different"}) // Miss
	
	stats := cache.GetStats()
	
	// Check that stats contain expected keys
	expectedKeys := []string{
		"hits", "misses", "evictions", "hit_rate", 
		"cache_size", "max_size", "cost_benefit_models",
	}
	
	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Stats missing key: %s", key)
		}
	}
	
	// Check specific values
	if stats["hits"].(int64) != 1 {
		t.Errorf("Expected 1 hit, got %v", stats["hits"])
	}
	
	if stats["misses"].(int64) != 1 {
		t.Errorf("Expected 1 miss, got %v", stats["misses"])
	}
	
	if stats["cache_size"].(int) != 1 {
		t.Errorf("Expected cache size 1, got %v", stats["cache_size"])
	}
	
	if stats["max_size"].(int) != 100 {
		t.Errorf("Expected max size 100, got %v", stats["max_size"])
	}
	
	hitRate := stats["hit_rate"].(float64)
	expectedHitRate := 0.5 // 1 hit out of 2 total requests
	if hitRate != expectedHitRate {
		t.Errorf("Expected hit rate %f, got %f", expectedHitRate, hitRate)
	}
}

func TestIntelligentCache_recordAccess(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	key := "test_key"
	
	// Record some accesses
	for i := 0; i < 5; i++ {
		cache.recordAccess(key)
		time.Sleep(time.Millisecond) // Small delay to ensure different timestamps
	}
	
	// Check access patterns were recorded
	patterns, exists := cache.accessPatterns[key]
	if !exists {
		t.Error("Access patterns not recorded")
	}
	
	if len(patterns) != 5 {
		t.Errorf("Expected 5 access records, got %d", len(patterns))
	}
	
	// Record many more accesses to test trimming
	for i := 0; i < 150; i++ {
		cache.recordAccess(key)
	}
	
	patterns = cache.accessPatterns[key]
	if len(patterns) > 100 {
		t.Errorf("Access patterns should be trimmed to 100, got %d", len(patterns))
	}
}

func TestIntelligentCache_predictAccessProbability(t *testing.T) {
	cache := NewIntelligentCache(100, time.Hour)
	
	functionName := "TestFunction"
	textHash := "hash123"
	
	// Test with no history
	prob1 := cache.predictAccessProbability(functionName, textHash)
	if prob1 != 0.1 {
		t.Errorf("Expected default probability 0.1, got %f", prob1)
	}
	
	// Add some recent access patterns
	key := functionName + ":" + textHash
	now := time.Now()
	
	// Add recent accesses (within last hour)
	for i := 0; i < 5; i++ {
		cache.accessPatterns[key] = append(cache.accessPatterns[key], now.Add(-time.Minute*time.Duration(i)))
	}
	
	// Add old accesses (beyond last hour)
	for i := 0; i < 3; i++ {
		cache.accessPatterns[key] = append(cache.accessPatterns[key], now.Add(-time.Hour*time.Duration(i+2)))
	}
	
	prob2 := cache.predictAccessProbability(functionName, textHash)
	
	// Should be 5 recent accesses / 10 = 0.5
	expectedProb := 0.5
	if prob2 != expectedProb {
		t.Errorf("Expected probability %f, got %f", expectedProb, prob2)
	}
	
	// Test capping at 1.0
	for i := 0; i < 20; i++ {
		cache.accessPatterns[key] = append(cache.accessPatterns[key], now.Add(-time.Minute*time.Duration(i)))
	}
	
	prob3 := cache.predictAccessProbability(functionName, textHash)
	if prob3 != 1.0 {
		t.Errorf("Expected probability capped at 1.0, got %f", prob3)
	}
}

// Benchmark tests
func BenchmarkIntelligentCache_Get(b *testing.B) {
	cache := NewIntelligentCache(1000, time.Hour)
	
	// Pre-populate cache
	for i := 0; i < 100; i++ {
		functionName := "TestFunction"
		text := "Hello world"
		params := map[string]interface{}{"param1": i}
		value := i
		computeCost := time.Millisecond * 100
		
		cache.Set(functionName, text, params, value, computeCost)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("TestFunction", "Hello world", map[string]interface{}{"param1": i % 100})
	}
}

func BenchmarkIntelligentCache_Set(b *testing.B) {
	cache := NewIntelligentCache(10000, time.Hour)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		functionName := "TestFunction"
		text := "Hello world"
		params := map[string]interface{}{"param1": i}
		value := i
		computeCost := time.Millisecond * 100
		
		cache.Set(functionName, text, params, value, computeCost)
	}
}

func BenchmarkIntelligentCache_generateCacheKey(b *testing.B) {
	cache := NewIntelligentCache(100, time.Hour)
	
	functionName := "TestFunction"
	text := "Hello world"
	params := map[string]interface{}{
		"param1": "value1",
		"param2": 42,
		"param3": true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.generateCacheKey(functionName, text, params)
	}
}