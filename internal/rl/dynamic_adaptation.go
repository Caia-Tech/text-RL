package rl

import (
	"math"
	"sync"
	"time"
)

// Dynamic adaptation system that adjusts optimization strategies in real-time
type DynamicAdaptationEngine struct {
	// Runtime monitoring
	performanceMonitor   *PerformanceMonitor
	resourceMonitor      *ResourceMonitor
	adaptationHistory    []AdaptationEvent
	
	// Adaptation strategies
	strategies           map[string]AdaptationStrategy
	activeStrategy       string
	strategyPerformance  map[string]StrategyMetrics
	
	// Configuration
	adaptationThreshold  float64
	monitoringInterval   time.Duration
	adaptationCooldown   time.Duration
	lastAdaptation       time.Time
	
	// Thread safety
	mutex sync.RWMutex
}

type PerformanceMonitor struct {
	metrics        []PerformanceSnapshot
	maxHistory     int
	currentWindow  time.Duration
	
	// Performance indicators
	avgLatency     float64
	throughput     float64
	errorRate      float64
	memoryPressure float64
}

type ResourceMonitor struct {
	cpuUsage      float64
	memoryUsage   float64
	diskIO        float64
	networkIO     float64
	
	// Resource limits
	cpuThreshold    float64
	memoryThreshold float64
	diskThreshold   float64
}

type PerformanceSnapshot struct {
	Timestamp      time.Time `json:"timestamp"`
	Latency        float64   `json:"latency_ms"`
	Throughput     float64   `json:"throughput_ops_sec"`
	MemoryUsage    int64     `json:"memory_usage_bytes"`
	CPUUsage       float64   `json:"cpu_usage_percent"`
	ErrorCount     int       `json:"error_count"`
	CacheHitRate   float64   `json:"cache_hit_rate"`
	ActiveActions  int       `json:"active_actions"`
}

type AdaptationEvent struct {
	Timestamp      time.Time              `json:"timestamp"`
	Trigger        string                 `json:"trigger"`
	OldStrategy    string                 `json:"old_strategy"`
	NewStrategy    string                 `json:"new_strategy"`
	Reason         string                 `json:"reason"`
	ImpactMetrics  map[string]float64     `json:"impact_metrics"`
	Success        bool                   `json:"success"`
}

type AdaptationStrategy struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Conditions     []AdaptationCondition  `json:"conditions"`
	Actions        []AdaptationAction     `json:"actions"`
	Priority       int                    `json:"priority"`
	CooldownPeriod time.Duration          `json:"cooldown_period"`
}

type AdaptationCondition struct {
	Metric      string  `json:"metric"`
	Operator    string  `json:"operator"` // ">", "<", ">=", "<=", "=="
	Threshold   float64 `json:"threshold"`
	Duration    time.Duration `json:"duration"` // How long condition must persist
}

type AdaptationAction struct {
	Type        string                 `json:"type"`
	Target      string                 `json:"target"`
	Parameters  map[string]interface{} `json:"parameters"`
	Rollback    map[string]interface{} `json:"rollback"`
}

type StrategyMetrics struct {
	ActivationCount   int           `json:"activation_count"`
	SuccessCount      int           `json:"success_count"`
	FailureCount      int           `json:"failure_count"`
	AvgImprovementPct float64       `json:"avg_improvement_pct"`
	LastActivation    time.Time     `json:"last_activation"`
	AvgDuration       time.Duration `json:"avg_duration"`
}

func NewDynamicAdaptationEngine() *DynamicAdaptationEngine {
	engine := &DynamicAdaptationEngine{
		performanceMonitor:  NewPerformanceMonitor(),
		resourceMonitor:     NewResourceMonitor(),
		adaptationHistory:   make([]AdaptationEvent, 0),
		strategies:          make(map[string]AdaptationStrategy),
		strategyPerformance: make(map[string]StrategyMetrics),
		adaptationThreshold: 0.2, // 20% performance degradation
		monitoringInterval:  time.Second * 5,
		adaptationCooldown:  time.Minute * 2,
		activeStrategy:      "default",
	}
	
	// Initialize default strategies
	engine.initializeStrategies()
	
	return engine
}

func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics:       make([]PerformanceSnapshot, 0),
		maxHistory:    1000,
		currentWindow: time.Minute * 5,
	}
}

func NewResourceMonitor() *ResourceMonitor {
	return &ResourceMonitor{
		cpuThreshold:    80.0,  // 80% CPU
		memoryThreshold: 85.0,  // 85% Memory
		diskThreshold:   90.0,  // 90% Disk
	}
}

func (dae *DynamicAdaptationEngine) initializeStrategies() {
	// Strategy 1: High Latency Response
	dae.strategies["high_latency"] = AdaptationStrategy{
		Name:        "high_latency",
		Description: "Optimize for latency when response times are high",
		Priority:    1,
		CooldownPeriod: time.Minute * 3,
		Conditions: []AdaptationCondition{
			{
				Metric:    "avg_latency",
				Operator:  ">",
				Threshold: 500.0, // 500ms
				Duration:  time.Second * 30,
			},
		},
		Actions: []AdaptationAction{
			{
				Type:   "adjust_parameters",
				Target: "cache_size",
				Parameters: map[string]interface{}{
					"multiplier": 1.5,
					"max_size":   2000,
				},
			},
			{
				Type:   "adjust_parameters",
				Target: "parallelism",
				Parameters: map[string]interface{}{
					"max_parallel": 1, // Reduce parallelism to avoid contention
				},
			},
		},
	}
	
	// Strategy 2: High Memory Usage Response
	dae.strategies["high_memory"] = AdaptationStrategy{
		Name:        "high_memory",
		Description: "Optimize for memory when usage is high",
		Priority:    2,
		CooldownPeriod: time.Minute * 2,
		Conditions: []AdaptationCondition{
			{
				Metric:    "memory_usage",
				Operator:  ">",
				Threshold: 80.0, // 80% memory usage
				Duration:  time.Second * 20,
			},
		},
		Actions: []AdaptationAction{
			{
				Type:   "adjust_parameters",
				Target: "cache_size",
				Parameters: map[string]interface{}{
					"multiplier": 0.5,
					"min_size":   50,
				},
			},
			{
				Type:   "enable_feature",
				Target: "aggressive_gc",
				Parameters: map[string]interface{}{
					"enabled": true,
				},
			},
		},
	}
	
	// Strategy 3: Low Throughput Response
	dae.strategies["low_throughput"] = AdaptationStrategy{
		Name:        "low_throughput",
		Description: "Optimize for throughput when processing rate is low",
		Priority:    3,
		CooldownPeriod: time.Minute * 4,
		Conditions: []AdaptationCondition{
			{
				Metric:    "throughput",
				Operator:  "<",
				Threshold: 10.0, // Less than 10 ops/sec
				Duration:  time.Second * 45,
			},
		},
		Actions: []AdaptationAction{
			{
				Type:   "adjust_parameters",
				Target: "parallelism",
				Parameters: map[string]interface{}{
					"max_parallel": 5, // Increase parallelism
				},
			},
			{
				Type:   "adjust_parameters",
				Target: "batch_size",
				Parameters: map[string]interface{}{
					"size": 10,
				},
			},
		},
	}
	
	// Strategy 4: High Error Rate Response
	dae.strategies["high_errors"] = AdaptationStrategy{
		Name:        "high_errors",
		Description: "Increase reliability when error rate is high",
		Priority:    1, // Highest priority
		CooldownPeriod: time.Minute * 5,
		Conditions: []AdaptationCondition{
			{
				Metric:    "error_rate",
				Operator:  ">",
				Threshold: 5.0, // 5% error rate
				Duration:  time.Second * 15,
			},
		},
		Actions: []AdaptationAction{
			{
				Type:   "adjust_parameters",
				Target: "retry_count",
				Parameters: map[string]interface{}{
					"retries": 3,
				},
			},
			{
				Type:   "adjust_parameters",
				Target: "timeout",
				Parameters: map[string]interface{}{
					"timeout_ms": 10000, // Increase timeout
				},
			},
		},
	}
}

// Start monitoring and adaptation
func (dae *DynamicAdaptationEngine) Start() {
	go dae.monitoringLoop()
	go dae.adaptationLoop()
}

func (dae *DynamicAdaptationEngine) monitoringLoop() {
	ticker := time.NewTicker(dae.monitoringInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		dae.collectMetrics()
		dae.updatePerformanceIndicators()
	}
}

func (dae *DynamicAdaptationEngine) adaptationLoop() {
	ticker := time.NewTicker(time.Second * 10) // Check every 10 seconds
	defer ticker.Stop()
	
	for range ticker.C {
		dae.evaluateAdaptationNeeds()
	}
}

func (dae *DynamicAdaptationEngine) collectMetrics() {
	dae.mutex.Lock()
	defer dae.mutex.Unlock()
	
	snapshot := PerformanceSnapshot{
		Timestamp:     time.Now(),
		Latency:       dae.measureCurrentLatency(),
		Throughput:    dae.measureCurrentThroughput(),
		MemoryUsage:   dae.measureCurrentMemoryUsage(),
		CPUUsage:      dae.measureCurrentCPUUsage(),
		ErrorCount:    dae.measureCurrentErrorCount(),
		CacheHitRate:  dae.measureCurrentCacheHitRate(),
		ActiveActions: dae.measureActiveActions(),
	}
	
	dae.performanceMonitor.metrics = append(dae.performanceMonitor.metrics, snapshot)
	
	// Keep only recent metrics
	if len(dae.performanceMonitor.metrics) > dae.performanceMonitor.maxHistory {
		dae.performanceMonitor.metrics = dae.performanceMonitor.metrics[1:]
	}
}

func (dae *DynamicAdaptationEngine) updatePerformanceIndicators() {
	if len(dae.performanceMonitor.metrics) == 0 {
		return
	}
	
	// Calculate rolling averages for current window
	cutoff := time.Now().Add(-dae.performanceMonitor.currentWindow)
	recentMetrics := make([]PerformanceSnapshot, 0)
	
	for _, metric := range dae.performanceMonitor.metrics {
		if metric.Timestamp.After(cutoff) {
			recentMetrics = append(recentMetrics, metric)
		}
	}
	
	if len(recentMetrics) == 0 {
		return
	}
	
	// Calculate averages
	totalLatency := 0.0
	totalThroughput := 0.0
	totalErrors := 0
	totalOperations := len(recentMetrics)
	
	for _, metric := range recentMetrics {
		totalLatency += metric.Latency
		totalThroughput += metric.Throughput
		totalErrors += metric.ErrorCount
	}
	
	dae.performanceMonitor.avgLatency = totalLatency / float64(totalOperations)
	dae.performanceMonitor.throughput = totalThroughput / float64(totalOperations)
	dae.performanceMonitor.errorRate = (float64(totalErrors) / float64(totalOperations)) * 100.0
	
	// Update resource monitor
	if len(recentMetrics) > 0 {
		latest := recentMetrics[len(recentMetrics)-1]
		dae.resourceMonitor.cpuUsage = latest.CPUUsage
		dae.resourceMonitor.memoryUsage = float64(latest.MemoryUsage) / (1024 * 1024 * 1024) // Convert to GB
	}
}

func (dae *DynamicAdaptationEngine) evaluateAdaptationNeeds() {
	dae.mutex.RLock()
	defer dae.mutex.RUnlock()
	
	// Check cooldown period
	if time.Since(dae.lastAdaptation) < dae.adaptationCooldown {
		return
	}
	
	// Evaluate each strategy
	triggeredStrategies := make([]AdaptationStrategy, 0)
	
	for _, strategy := range dae.strategies {
		if dae.evaluateStrategy(strategy) {
			triggeredStrategies = append(triggeredStrategies, strategy)
		}
	}
	
	// Apply highest priority strategy
	if len(triggeredStrategies) > 0 {
		// Sort by priority
		bestStrategy := triggeredStrategies[0]
		for _, strategy := range triggeredStrategies[1:] {
			if strategy.Priority < bestStrategy.Priority { // Lower number = higher priority
				bestStrategy = strategy
			}
		}
		
		dae.applyStrategy(bestStrategy)
	}
}

func (dae *DynamicAdaptationEngine) evaluateStrategy(strategy AdaptationStrategy) bool {
	// Check if strategy is in cooldown
	if metrics, exists := dae.strategyPerformance[strategy.Name]; exists {
		if time.Since(metrics.LastActivation) < strategy.CooldownPeriod {
			return false
		}
	}
	
	// Evaluate all conditions
	for _, condition := range strategy.Conditions {
		if !dae.evaluateCondition(condition) {
			return false
		}
	}
	
	return true
}

func (dae *DynamicAdaptationEngine) evaluateCondition(condition AdaptationCondition) bool {
	currentValue := dae.getCurrentMetricValue(condition.Metric)
	
	switch condition.Operator {
	case ">":
		return currentValue > condition.Threshold
	case "<":
		return currentValue < condition.Threshold
	case ">=":
		return currentValue >= condition.Threshold
	case "<=":
		return currentValue <= condition.Threshold
	case "==":
		return math.Abs(currentValue-condition.Threshold) < 0.001
	default:
		return false
	}
}

func (dae *DynamicAdaptationEngine) getCurrentMetricValue(metric string) float64 {
	switch metric {
	case "avg_latency":
		return dae.performanceMonitor.avgLatency
	case "throughput":
		return dae.performanceMonitor.throughput
	case "error_rate":
		return dae.performanceMonitor.errorRate
	case "memory_usage":
		return dae.resourceMonitor.memoryUsage
	case "cpu_usage":
		return dae.resourceMonitor.cpuUsage
	default:
		return 0.0
	}
}

func (dae *DynamicAdaptationEngine) applyStrategy(strategy AdaptationStrategy) {
	dae.mutex.Lock()
	defer dae.mutex.Unlock()
	
	oldStrategy := dae.activeStrategy
	
	// Apply adaptation actions
	success := true
	impactMetrics := make(map[string]float64)
	
	// Record baseline metrics
	baselineLatency := dae.performanceMonitor.avgLatency
	baselineMemory := dae.resourceMonitor.memoryUsage
	
	for _, action := range strategy.Actions {
		actionSuccess := dae.executeAdaptationAction(action)
		if !actionSuccess {
			success = false
			break
		}
	}
	
	// Record adaptation event
	event := AdaptationEvent{
		Timestamp:   time.Now(),
		Trigger:     "automatic",
		OldStrategy: oldStrategy,
		NewStrategy: strategy.Name,
		Reason:      strategy.Description,
		Success:     success,
		ImpactMetrics: impactMetrics,
	}
	
	dae.adaptationHistory = append(dae.adaptationHistory, event)
	
	// Update strategy metrics
	if _, exists := dae.strategyPerformance[strategy.Name]; !exists {
		dae.strategyPerformance[strategy.Name] = StrategyMetrics{}
	}
	
	metrics := dae.strategyPerformance[strategy.Name]
	metrics.ActivationCount++
	metrics.LastActivation = time.Now()
	
	if success {
		metrics.SuccessCount++
		dae.activeStrategy = strategy.Name
		
		// Calculate improvement
		time.Sleep(time.Second * 5) // Wait for metrics to stabilize
		newLatency := dae.measureCurrentLatency()
		newMemory := dae.measureCurrentMemoryUsage()
		
		latencyImprovement := ((baselineLatency - newLatency) / baselineLatency) * 100.0
		memoryImprovement := ((float64(baselineMemory) - float64(newMemory)) / float64(baselineMemory)) * 100.0
		
		metrics.AvgImprovementPct = (metrics.AvgImprovementPct + latencyImprovement) / 2.0
		event.ImpactMetrics["latency_improvement_pct"] = latencyImprovement
		event.ImpactMetrics["memory_improvement_pct"] = memoryImprovement
	} else {
		metrics.FailureCount++
	}
	
	dae.strategyPerformance[strategy.Name] = metrics
	dae.lastAdaptation = time.Now()
}

func (dae *DynamicAdaptationEngine) executeAdaptationAction(action AdaptationAction) bool {
	// This would integrate with the actual system parameters
	// For now, return success for all actions
	return true
}

// Measurement functions (would integrate with actual system monitoring)
func (dae *DynamicAdaptationEngine) measureCurrentLatency() float64 {
	// Integrate with actual latency monitoring
	return 100.0 + (50.0) // Simulated latency
}

func (dae *DynamicAdaptationEngine) measureCurrentThroughput() float64 {
	// Integrate with actual throughput monitoring
	return 20.0 + (30.0) // Simulated throughput
}

func (dae *DynamicAdaptationEngine) measureCurrentMemoryUsage() int64 {
	// Integrate with actual memory monitoring
	return 1024 * 1024 * 1024 // Simulated 1GB
}

func (dae *DynamicAdaptationEngine) measureCurrentCPUUsage() float64 {
	// Integrate with actual CPU monitoring
	return 50.0 + (40.0) // Simulated CPU usage
}

func (dae *DynamicAdaptationEngine) measureCurrentErrorCount() int {
	// Integrate with actual error monitoring
	return 0 // Simulated error count
}

func (dae *DynamicAdaptationEngine) measureCurrentCacheHitRate() float64 {
	// Integrate with actual cache monitoring
	return 0.8 + (0.15) // Simulated 80-95% hit rate
}

func (dae *DynamicAdaptationEngine) measureActiveActions() int {
	// Integrate with actual action monitoring
	return 2 + (3) // Simulated 2-5 active actions
}

// Get adaptation status and history
func (dae *DynamicAdaptationEngine) GetStatus() map[string]interface{} {
	dae.mutex.RLock()
	defer dae.mutex.RUnlock()
	
	return map[string]interface{}{
		"active_strategy":       dae.activeStrategy,
		"performance_metrics":   dae.performanceMonitor,
		"resource_metrics":      dae.resourceMonitor,
		"adaptation_history":    dae.adaptationHistory,
		"strategy_performance":  dae.strategyPerformance,
		"last_adaptation":       dae.lastAdaptation,
	}
}

// Manual strategy activation
func (dae *DynamicAdaptationEngine) ActivateStrategy(strategyName string, reason string) bool {
	dae.mutex.Lock()
	defer dae.mutex.Unlock()
	
	strategy, exists := dae.strategies[strategyName]
	if !exists {
		return false
	}
	
	// Apply strategy with manual trigger
	oldStrategy := dae.activeStrategy
	success := true
	
	for _, action := range strategy.Actions {
		actionSuccess := dae.executeAdaptationAction(action)
		if !actionSuccess {
			success = false
			break
		}
	}
	
	// Record manual adaptation event
	event := AdaptationEvent{
		Timestamp:   time.Now(),
		Trigger:     "manual",
		OldStrategy: oldStrategy,
		NewStrategy: strategyName,
		Reason:      reason,
		Success:     success,
	}
	
	dae.adaptationHistory = append(dae.adaptationHistory, event)
	
	if success {
		dae.activeStrategy = strategyName
	}
	
	return success
}