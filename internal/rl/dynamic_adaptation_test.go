package rl

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestNewDynamicAdaptationEngine(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	if engine == nil {
		t.Fatal("NewDynamicAdaptationEngine returned nil")
	}
	
	if engine.performanceMonitor == nil {
		t.Error("performanceMonitor not initialized")
	}
	
	if engine.resourceMonitor == nil {
		t.Error("resourceMonitor not initialized")
	}
	
	if engine.adaptationHistory == nil {
		t.Error("adaptationHistory not initialized")
	}
	
	if engine.strategies == nil {
		t.Error("strategies map not initialized")
	}
	
	if engine.strategyPerformance == nil {
		t.Error("strategyPerformance map not initialized")
	}
	
	if engine.adaptationThreshold != 0.2 {
		t.Errorf("Expected adaptationThreshold 0.2, got %f", engine.adaptationThreshold)
	}
	
	if engine.monitoringInterval != time.Second*5 {
		t.Errorf("Expected monitoringInterval 5s, got %v", engine.monitoringInterval)
	}
	
	if engine.adaptationCooldown != time.Minute*2 {
		t.Errorf("Expected adaptationCooldown 2m, got %v", engine.adaptationCooldown)
	}
	
	if engine.activeStrategy != "default" {
		t.Errorf("Expected activeStrategy 'default', got %s", engine.activeStrategy)
	}
}

func TestNewPerformanceMonitor(t *testing.T) {
	monitor := NewPerformanceMonitor()
	
	if monitor == nil {
		t.Fatal("NewPerformanceMonitor returned nil")
	}
	
	if monitor.metrics == nil {
		t.Error("metrics slice not initialized")
	}
	
	if monitor.maxHistory != 1000 {
		t.Errorf("Expected maxHistory 1000, got %d", monitor.maxHistory)
	}
	
	if monitor.currentWindow != time.Minute*5 {
		t.Errorf("Expected currentWindow 5m, got %v", monitor.currentWindow)
	}
}

func TestNewResourceMonitor(t *testing.T) {
	monitor := NewResourceMonitor()
	
	if monitor == nil {
		t.Fatal("NewResourceMonitor returned nil")
	}
	
	if monitor.cpuThreshold != 80.0 {
		t.Errorf("Expected cpuThreshold 80.0, got %f", monitor.cpuThreshold)
	}
	
	if monitor.memoryThreshold != 85.0 {
		t.Errorf("Expected memoryThreshold 85.0, got %f", monitor.memoryThreshold)
	}
	
	if monitor.diskThreshold != 90.0 {
		t.Errorf("Expected diskThreshold 90.0, got %f", monitor.diskThreshold)
	}
}

func TestDynamicAdaptationEngine_initializeStrategies(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	expectedStrategies := []string{
		"high_latency",
		"high_memory",
		"low_throughput",
		"high_errors",
	}
	
	for _, strategyName := range expectedStrategies {
		strategy, exists := engine.strategies[strategyName]
		if !exists {
			t.Errorf("Missing strategy: %s", strategyName)
			continue
		}
		
		if strategy.Name != strategyName {
			t.Errorf("Strategy %s has incorrect name: %s", strategyName, strategy.Name)
		}
		
		if strategy.Description == "" {
			t.Errorf("Strategy %s has empty description", strategyName)
		}
		
		if len(strategy.Conditions) == 0 {
			t.Errorf("Strategy %s has no conditions", strategyName)
		}
		
		if len(strategy.Actions) == 0 {
			t.Errorf("Strategy %s has no actions", strategyName)
		}
		
		if strategy.Priority < 1 {
			t.Errorf("Strategy %s has invalid priority: %d", strategyName, strategy.Priority)
		}
		
		if strategy.CooldownPeriod <= 0 {
			t.Errorf("Strategy %s has invalid cooldown period: %v", strategyName, strategy.CooldownPeriod)
		}
		
		// Check conditions have required fields
		for i, condition := range strategy.Conditions {
			if condition.Metric == "" {
				t.Errorf("Strategy %s condition %d has empty metric", strategyName, i)
			}
			
			if condition.Operator == "" {
				t.Errorf("Strategy %s condition %d has empty operator", strategyName, i)
			}
			
			if condition.Duration <= 0 {
				t.Errorf("Strategy %s condition %d has invalid duration: %v", strategyName, i, condition.Duration)
			}
			
			// Check operator is valid
			validOperators := map[string]bool{">": true, "<": true, ">=": true, "<=": true, "==": true}
			if !validOperators[condition.Operator] {
				t.Errorf("Strategy %s condition %d has invalid operator: %s", strategyName, i, condition.Operator)
			}
		}
		
		// Check actions have required fields
		for i, action := range strategy.Actions {
			if action.Type == "" {
				t.Errorf("Strategy %s action %d has empty type", strategyName, i)
			}
			
			if action.Target == "" {
				t.Errorf("Strategy %s action %d has empty target", strategyName, i)
			}
			
			if action.Parameters == nil {
				t.Errorf("Strategy %s action %d has nil parameters", strategyName, i)
			}
		}
	}
}

func TestDynamicAdaptationEngine_getCurrentMetricValue(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	// Set some test values
	engine.performanceMonitor.avgLatency = 150.0
	engine.performanceMonitor.throughput = 25.5
	engine.performanceMonitor.errorRate = 2.5
	engine.resourceMonitor.memoryUsage = 75.0
	engine.resourceMonitor.cpuUsage = 60.0
	
	tests := []struct {
		metric   string
		expected float64
	}{
		{"avg_latency", 150.0},
		{"throughput", 25.5},
		{"error_rate", 2.5},
		{"memory_usage", 75.0},
		{"cpu_usage", 60.0},
		{"unknown_metric", 0.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.metric, func(t *testing.T) {
			value := engine.getCurrentMetricValue(tt.metric)
			if value != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, value)
			}
		})
	}
}

func TestDynamicAdaptationEngine_evaluateCondition(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	// Set test metric values
	engine.performanceMonitor.avgLatency = 600.0 // High latency
	engine.performanceMonitor.throughput = 8.0   // Low throughput
	engine.performanceMonitor.errorRate = 6.0    // High error rate
	engine.resourceMonitor.memoryUsage = 85.0    // High memory
	
	tests := []struct {
		name      string
		condition AdaptationCondition
		expected  bool
	}{
		{
			name: "high latency condition met",
			condition: AdaptationCondition{
				Metric:    "avg_latency",
				Operator:  ">",
				Threshold: 500.0,
			},
			expected: true,
		},
		{
			name: "high latency condition not met",
			condition: AdaptationCondition{
				Metric:    "avg_latency",
				Operator:  ">",
				Threshold: 700.0,
			},
			expected: false,
		},
		{
			name: "low throughput condition met",
			condition: AdaptationCondition{
				Metric:    "throughput",
				Operator:  "<",
				Threshold: 10.0,
			},
			expected: true,
		},
		{
			name: "memory usage condition met",
			condition: AdaptationCondition{
				Metric:    "memory_usage",
				Operator:  ">=",
				Threshold: 85.0,
			},
			expected: true,
		},
		{
			name: "equality condition met",
			condition: AdaptationCondition{
				Metric:    "error_rate",
				Operator:  "==",
				Threshold: 6.0,
			},
			expected: true,
		},
		{
			name: "less than or equal condition met",
			condition: AdaptationCondition{
				Metric:    "error_rate",
				Operator:  "<=",
				Threshold: 6.0,
			},
			expected: true,
		},
		{
			name: "invalid operator",
			condition: AdaptationCondition{
				Metric:    "avg_latency",
				Operator:  "invalid",
				Threshold: 500.0,
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.evaluateCondition(tt.condition)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDynamicAdaptationEngine_evaluateStrategy(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	// Create test strategy
	strategy := AdaptationStrategy{
		Name:           "test_strategy",
		CooldownPeriod: time.Minute,
		Conditions: []AdaptationCondition{
			{
				Metric:    "avg_latency",
				Operator:  ">",
				Threshold: 500.0,
			},
			{
				Metric:    "error_rate",
				Operator:  "<",
				Threshold: 10.0,
			},
		},
	}
	
	// Set metric values that satisfy conditions
	engine.performanceMonitor.avgLatency = 600.0
	engine.performanceMonitor.errorRate = 5.0
	
	// Test without cooldown restriction
	result := engine.evaluateStrategy(strategy)
	if !result {
		t.Error("Expected strategy to be triggered when conditions are met")
	}
	
	// Test with cooldown restriction
	engine.strategyPerformance["test_strategy"] = StrategyMetrics{
		LastActivation: time.Now().Add(-time.Second * 30), // Recent activation
	}
	
	result = engine.evaluateStrategy(strategy)
	if result {
		t.Error("Expected strategy to be blocked by cooldown")
	}
	
	// Test with expired cooldown
	engine.strategyPerformance["test_strategy"] = StrategyMetrics{
		LastActivation: time.Now().Add(-time.Minute * 2), // Old activation
	}
	
	result = engine.evaluateStrategy(strategy)
	if !result {
		t.Error("Expected strategy to be triggered after cooldown expires")
	}
	
	// Test with condition not met
	engine.performanceMonitor.avgLatency = 400.0 // Below threshold
	
	result = engine.evaluateStrategy(strategy)
	if result {
		t.Error("Expected strategy not to be triggered when condition not met")
	}
}

func TestDynamicAdaptationEngine_collectMetrics(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	initialCount := len(engine.performanceMonitor.metrics)
	
	engine.collectMetrics()
	
	// Should have added one metric
	newCount := len(engine.performanceMonitor.metrics)
	if newCount != initialCount+1 {
		t.Errorf("Expected %d metrics, got %d", initialCount+1, newCount)
	}
	
	// Check the metric has required fields
	if newCount > 0 {
		metric := engine.performanceMonitor.metrics[newCount-1]
		
		if metric.Timestamp.IsZero() {
			t.Error("Metric timestamp not set")
		}
		
		if metric.Latency < 0 {
			t.Errorf("Invalid latency: %f", metric.Latency)
		}
		
		if metric.Throughput < 0 {
			t.Errorf("Invalid throughput: %f", metric.Throughput)
		}
		
		if metric.MemoryUsage < 0 {
			t.Errorf("Invalid memory usage: %d", metric.MemoryUsage)
		}
		
		if metric.CPUUsage < 0 || metric.CPUUsage > 100 {
			t.Errorf("Invalid CPU usage: %f", metric.CPUUsage)
		}
		
		if metric.ErrorCount < 0 {
			t.Errorf("Invalid error count: %d", metric.ErrorCount)
		}
		
		if metric.CacheHitRate < 0 || metric.CacheHitRate > 1 {
			t.Errorf("Invalid cache hit rate: %f", metric.CacheHitRate)
		}
		
		if metric.ActiveActions < 0 {
			t.Errorf("Invalid active actions: %d", metric.ActiveActions)
		}
	}
	
	// Test metric history limit
	for i := 0; i < engine.performanceMonitor.maxHistory+10; i++ {
		engine.collectMetrics()
	}
	
	finalCount := len(engine.performanceMonitor.metrics)
	if finalCount > engine.performanceMonitor.maxHistory {
		t.Errorf("Metrics exceed max history: %d > %d", finalCount, engine.performanceMonitor.maxHistory)
	}
}

func TestDynamicAdaptationEngine_updatePerformanceIndicators(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	// Add some test metrics
	now := time.Now()
	testMetrics := []PerformanceSnapshot{
		{
			Timestamp:    now.Add(-time.Minute * 1),
			Latency:      100.0,
			Throughput:   20.0,
			ErrorCount:   1,
			MemoryUsage:  1024 * 1024 * 1024, // 1GB
		},
		{
			Timestamp:    now.Add(-time.Minute * 2),
			Latency:      200.0,
			Throughput:   15.0,
			ErrorCount:   0,
			MemoryUsage:  2 * 1024 * 1024 * 1024, // 2GB
		},
		{
			Timestamp:    now.Add(-time.Minute * 3),
			Latency:      150.0,
			Throughput:   25.0,
			ErrorCount:   2,
			MemoryUsage:  1.5 * 1024 * 1024 * 1024, // 1.5GB
		},
	}
	
	engine.performanceMonitor.metrics = testMetrics
	
	engine.updatePerformanceIndicators()
	
	// Check calculated averages
	expectedAvgLatency := (100.0 + 200.0 + 150.0) / 3.0
	if math.Abs(engine.performanceMonitor.avgLatency-expectedAvgLatency) > 0.001 {
		t.Errorf("Expected avg latency %f, got %f", expectedAvgLatency, engine.performanceMonitor.avgLatency)
	}
	
	expectedThroughput := (20.0 + 15.0 + 25.0) / 3.0
	if math.Abs(engine.performanceMonitor.throughput-expectedThroughput) > 0.001 {
		t.Errorf("Expected throughput %f, got %f", expectedThroughput, engine.performanceMonitor.throughput)
	}
	
	expectedErrorRate := (float64(1+0+2) / 3.0) * 100.0
	if math.Abs(engine.performanceMonitor.errorRate-expectedErrorRate) > 0.001 {
		t.Errorf("Expected error rate %f, got %f", expectedErrorRate, engine.performanceMonitor.errorRate)
	}
	
	// Memory usage should be from latest metric (converted to GB)
	expectedMemoryGB := 1.5
	if math.Abs(engine.resourceMonitor.memoryUsage-expectedMemoryGB) > 0.001 {
		t.Errorf("Expected memory usage %f GB, got %f", expectedMemoryGB, engine.resourceMonitor.memoryUsage)
	}
}

func TestDynamicAdaptationEngine_executeAdaptationAction(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	tests := []struct {
		name   string
		action AdaptationAction
	}{
		{
			name: "adjust parameters action",
			action: AdaptationAction{
				Type:   "adjust_parameters",
				Target: "cache_size",
				Parameters: map[string]interface{}{
					"multiplier": 1.5,
					"max_size":   2000,
				},
			},
		},
		{
			name: "enable feature action",
			action: AdaptationAction{
				Type:   "enable_feature",
				Target: "aggressive_gc",
				Parameters: map[string]interface{}{
					"enabled": true,
				},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := engine.executeAdaptationAction(tt.action)
			
			// Currently all actions return true (stub implementation)
			if !success {
				t.Error("Expected action to succeed")
			}
		})
	}
}

func TestDynamicAdaptationEngine_ActivateStrategy(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	// Test activating existing strategy
	result := engine.ActivateStrategy("high_latency", "manual test")
	if !result {
		t.Error("Expected manual strategy activation to succeed")
	}
	
	// Check that strategy was activated
	if engine.activeStrategy != "high_latency" {
		t.Errorf("Expected activeStrategy 'high_latency', got %s", engine.activeStrategy)
	}
	
	// Check that adaptation event was recorded
	if len(engine.adaptationHistory) == 0 {
		t.Error("Expected adaptation event to be recorded")
	} else {
		event := engine.adaptationHistory[len(engine.adaptationHistory)-1]
		if event.Trigger != "manual" {
			t.Errorf("Expected trigger 'manual', got %s", event.Trigger)
		}
		if event.NewStrategy != "high_latency" {
			t.Errorf("Expected new strategy 'high_latency', got %s", event.NewStrategy)
		}
		if event.Reason != "manual test" {
			t.Errorf("Expected reason 'manual test', got %s", event.Reason)
		}
	}
	
	// Test activating non-existent strategy
	result = engine.ActivateStrategy("nonexistent", "test")
	if result {
		t.Error("Expected activation of non-existent strategy to fail")
	}
}

func TestDynamicAdaptationEngine_GetStatus(t *testing.T) {
	engine := NewDynamicAdaptationEngine()
	
	// Add some test data
	engine.activeStrategy = "test_strategy"
	engine.adaptationHistory = []AdaptationEvent{
		{
			Timestamp:   time.Now(),
			Trigger:     "automatic",
			OldStrategy: "default",
			NewStrategy: "test_strategy",
			Success:     true,
		},
	}
	
	status := engine.GetStatus()
	
	expectedKeys := []string{
		"active_strategy",
		"performance_metrics",
		"resource_metrics",
		"adaptation_history",
		"strategy_performance",
		"last_adaptation",
	}
	
	for _, key := range expectedKeys {
		if _, exists := status[key]; !exists {
			t.Errorf("Status missing key: %s", key)
		}
	}
	
	if status["active_strategy"] != "test_strategy" {
		t.Errorf("Expected active_strategy 'test_strategy', got %v", status["active_strategy"])
	}
	
	if reflect.TypeOf(status["performance_metrics"]) != reflect.TypeOf(&PerformanceMonitor{}) {
		t.Error("performance_metrics has wrong type")
	}
	
	if reflect.TypeOf(status["resource_metrics"]) != reflect.TypeOf(&ResourceMonitor{}) {
		t.Error("resource_metrics has wrong type")
	}
	
	adaptationHistory, ok := status["adaptation_history"].([]AdaptationEvent)
	if !ok {
		t.Error("adaptation_history has wrong type")
	} else if len(adaptationHistory) != 1 {
		t.Errorf("Expected 1 adaptation event, got %d", len(adaptationHistory))
	}
}

func TestAdaptationCondition_Validation(t *testing.T) {
	validCondition := AdaptationCondition{
		Metric:    "avg_latency",
		Operator:  ">",
		Threshold: 500.0,
		Duration:  time.Second * 30,
	}
	
	// Test that condition has all required fields
	if validCondition.Metric == "" {
		t.Error("Valid condition should have metric")
	}
	
	if validCondition.Operator == "" {
		t.Error("Valid condition should have operator")
	}
	
	if validCondition.Duration <= 0 {
		t.Error("Valid condition should have positive duration")
	}
}

func TestAdaptationAction_Validation(t *testing.T) {
	validAction := AdaptationAction{
		Type:   "adjust_parameters",
		Target: "cache_size",
		Parameters: map[string]interface{}{
			"multiplier": 1.5,
		},
		Rollback: map[string]interface{}{
			"multiplier": 1.0,
		},
	}
	
	// Test that action has all required fields
	if validAction.Type == "" {
		t.Error("Valid action should have type")
	}
	
	if validAction.Target == "" {
		t.Error("Valid action should have target")
	}
	
	if validAction.Parameters == nil {
		t.Error("Valid action should have parameters")
	}
}

// Benchmark tests for performance validation
func BenchmarkDynamicAdaptationEngine_collectMetrics(b *testing.B) {
	engine := NewDynamicAdaptationEngine()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.collectMetrics()
	}
}

func BenchmarkDynamicAdaptationEngine_updatePerformanceIndicators(b *testing.B) {
	engine := NewDynamicAdaptationEngine()
	
	// Pre-populate with metrics
	now := time.Now()
	for i := 0; i < 100; i++ {
		engine.performanceMonitor.metrics = append(engine.performanceMonitor.metrics, PerformanceSnapshot{
			Timestamp:    now.Add(-time.Duration(i) * time.Second),
			Latency:      float64(100 + i),
			Throughput:   float64(20 + i%10),
			ErrorCount:   i % 5,
			MemoryUsage:  int64(1024 * 1024 * (1024 + i)),
			CPUUsage:     float64(50 + i%30),
			CacheHitRate: 0.8 + float64(i%20)*0.01,
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.updatePerformanceIndicators()
	}
}

func BenchmarkDynamicAdaptationEngine_evaluateCondition(b *testing.B) {
	engine := NewDynamicAdaptationEngine()
	engine.performanceMonitor.avgLatency = 600.0
	
	condition := AdaptationCondition{
		Metric:    "avg_latency",
		Operator:  ">",
		Threshold: 500.0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.evaluateCondition(condition)
	}
}

func BenchmarkDynamicAdaptationEngine_evaluateStrategy(b *testing.B) {
	engine := NewDynamicAdaptationEngine()
	engine.performanceMonitor.avgLatency = 600.0
	engine.performanceMonitor.errorRate = 5.0
	
	strategy := AdaptationStrategy{
		Name:           "test_strategy",
		CooldownPeriod: time.Minute,
		Conditions: []AdaptationCondition{
			{
				Metric:    "avg_latency",
				Operator:  ">",
				Threshold: 500.0,
			},
			{
				Metric:    "error_rate",
				Operator:  "<",
				Threshold: 10.0,
			},
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.evaluateStrategy(strategy)
	}
}