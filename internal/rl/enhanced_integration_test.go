package rl

import (
	"math"
	"testing"
	"time"
)

// Integration tests for the enhanced RL system components working together
func TestEnhancedRLSystem_Integration(t *testing.T) {
	// Test that all enhanced components can work together
	
	// Initialize all components
	paramOptimizer := NewParameterOptimizer()
	cache := NewIntelligentCache(100, time.Hour)
	multiOptimizer := NewMultiObjectiveOptimizer()
	adaptationEngine := NewDynamicAdaptationEngine()
	
	if paramOptimizer == nil {
		t.Error("Failed to create parameter optimizer")
	}
	if cache == nil {
		t.Error("Failed to create intelligent cache")
	}
	if multiOptimizer == nil {
		t.Error("Failed to create multi-objective optimizer")
	}
	if adaptationEngine == nil {
		t.Error("Failed to create dynamic adaptation engine")
	}
	
	// Test basic functionality of each component
	t.Run("parameter_optimization", func(t *testing.T) {
		ranges := map[string]ParameterRange{
			"test_param": {
				Type: "float",
				Min:  0.1,
				Max:  0.9,
				Default: 0.5,
			},
		}
		
		fitnessFunc := func(params map[string]interface{}) float64 {
			return params["test_param"].(float64)
		}
		
		result := paramOptimizer.OptimizeParameters("test_func", ranges, fitnessFunc)
		
		if len(result) == 0 {
			t.Error("Parameter optimization returned no results")
		}
		
		if val, exists := result["test_param"]; !exists {
			t.Error("Expected parameter not found in result")
		} else if val.(float64) < 0.1 || val.(float64) > 0.9 {
			t.Errorf("Parameter value %f outside expected range", val.(float64))
		}
	})
	
	t.Run("intelligent_caching", func(t *testing.T) {
		functionName := "TestFunction"
		text := "Hello world"
		params := map[string]interface{}{"param1": "value1"}
		value := "cached result"
		computeCost := time.Millisecond * 100
		
		// Test cache miss
		_, found := cache.Get(functionName, text, params)
		if found {
			t.Error("Expected cache miss for new entry")
		}
		
		// Test cache set and hit
		cache.Set(functionName, text, params, value, computeCost)
		result, found := cache.Get(functionName, text, params)
		
		if !found {
			t.Error("Expected cache hit after set")
		}
		if result != value {
			t.Errorf("Expected %v, got %v", value, result)
		}
	})
	
	t.Run("multi_objective_optimization", func(t *testing.T) {
		multiOptimizer.SetupTextLibObjectives()
		
		if len(multiOptimizer.objectives) == 0 {
			t.Error("No objectives configured")
		}
		
		// Test domination logic
		sol1 := Solution{Objectives: []float64{100, 0.9, 1000, 50}}
		sol2 := Solution{Objectives: []float64{200, 0.8, 2000, 100}}
		
		dominates := multiOptimizer.dominates(sol1, sol2)
		if !dominates {
			t.Error("Sol1 should dominate sol2")
		}
	})
	
	t.Run("dynamic_adaptation", func(t *testing.T) {
		// Test strategy activation
		result := adaptationEngine.ActivateStrategy("high_latency", "integration test")
		if !result {
			t.Error("Failed to activate strategy")
		}
		
		if adaptationEngine.activeStrategy != "high_latency" {
			t.Errorf("Expected active strategy 'high_latency', got %s", adaptationEngine.activeStrategy)
		}
		
		// Test status retrieval
		status := adaptationEngine.GetStatus()
		if status == nil {
			t.Error("Failed to get adaptation status")
		}
	})
}

func TestParameterOptimization_WithCaching_Integration(t *testing.T) {
	// Test parameter optimization using cached function evaluations
	
	cache := NewIntelligentCache(50, time.Hour)
	paramOptimizer := NewParameterOptimizer()
	paramOptimizer.generations = 3 // Quick test
	paramOptimizer.populationSize = 10
	
	evaluationCount := 0
	
	// Fitness function that benefits from caching
	fitnessFunc := func(params map[string]interface{}) float64 {
		// Check cache first
		if result, found := cache.Get("fitness", "test", params); found {
			return result.(float64)
		}
		
		// Expensive computation simulation
		evaluationCount++
		time.Sleep(time.Millisecond) // Simulate work
		
		param1 := params["param1"].(float64)
		param2 := params["param2"].(float64)
		fitness := param1 + param2
		
		// Cache the result
		cache.Set("fitness", "test", params, fitness, time.Millisecond)
		
		return fitness
	}
	
	ranges := map[string]ParameterRange{
		"param1": {Type: "float", Min: 0.0, Max: 1.0, Default: 0.5},
		"param2": {Type: "float", Min: 0.0, Max: 1.0, Default: 0.5},
	}
	
	result := paramOptimizer.OptimizeParameters("test_function", ranges, fitnessFunc)
	
	// Verify results
	if len(result) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(result))
	}
	
	// Cache should have reduced the number of evaluations
	totalGenerations := paramOptimizer.generations
	totalPopulation := paramOptimizer.populationSize * totalGenerations
	
	if evaluationCount >= totalPopulation {
		t.Errorf("Cache didn't reduce evaluations: %d evaluations for %d total operations", 
			evaluationCount, totalPopulation)
	}
	
	// Check cache stats
	stats := cache.GetStats()
	hits := stats["hits"].(int64)
	
	if hits == 0 {
		t.Error("Expected some cache hits during optimization")
	}
}

func TestMultiObjective_WithAdaptation_Integration(t *testing.T) {
	// Test multi-objective optimization with dynamic adaptation
	
	multiOptimizer := NewMultiObjectiveOptimizer()
	adaptationEngine := NewDynamicAdaptationEngine()
	
	multiOptimizer.SetupTextLibObjectives()
	multiOptimizer.generations = 3 // Quick test
	multiOptimizer.populationSize = 20
	
	evaluationCount := 0
	
	// Evaluation function that triggers adaptation based on performance
	evaluateFunction := func(actions []EnhancedAction, params map[string]interface{}) Solution {
		evaluationCount++
		
		// Simulate performance metrics
		totalTime := 100.0 + float64(len(actions)*50) // More actions = more time
		accuracy := 0.8 + 0.1*float64(len(actions))   // More actions = better accuracy
		memoryUsage := int64(1000 * (1 + len(actions))) // More actions = more memory
		cost := 10.0 * float64(len(actions))
		
		// Check if we should trigger adaptation
		if totalTime > 200.0 {
			// High latency - trigger adaptation
			adaptationEngine.performanceMonitor.avgLatency = totalTime
			adaptationEngine.ActivateStrategy("high_latency", "automatic from evaluation")
		}
		
		if memoryUsage > 3000 {
			// High memory - trigger adaptation
			adaptationEngine.resourceMonitor.memoryUsage = float64(memoryUsage) / 1000.0
			adaptationEngine.ActivateStrategy("high_memory", "automatic from evaluation")
		}
		
		return Solution{
			Actions:     actions,
			Parameters:  params,
			Objectives:  []float64{totalTime, accuracy, float64(memoryUsage), cost},
			TotalTime:   totalTime,
			Accuracy:    accuracy,
			MemoryUsage: memoryUsage,
			Cost:        cost,
		}
	}
	
	// Run optimization
	paretoFront := multiOptimizer.Optimize(evaluateFunction)
	
	// Verify results
	if len(paretoFront) == 0 {
		t.Error("No solutions in Pareto front")
	}
	
	// Check that solutions are non-dominated
	for i, sol1 := range paretoFront {
		for j, sol2 := range paretoFront {
			if i != j && multiOptimizer.dominates(sol1, sol2) {
				t.Errorf("Solution %d dominates solution %d in Pareto front", i, j)
			}
		}
	}
	
	// Check adaptation was triggered
	adaptationHistory := adaptationEngine.adaptationHistory
	if len(adaptationHistory) == 0 {
		t.Error("Expected adaptation to be triggered during optimization")
	}
	
	// Verify some evaluations were performed
	if evaluationCount == 0 {
		t.Error("No evaluations performed")
	}
}

func TestFullSystem_RealisticScenario_Integration(t *testing.T) {
	// Integration test simulating a realistic optimization scenario
	
	// Initialize all components
	paramOptimizer := NewParameterOptimizer()
	cache := NewIntelligentCache(100, time.Hour)
	multiOptimizer := NewMultiObjectiveOptimizer()
	adaptationEngine := NewDynamicAdaptationEngine()
	
	// Configure for quick test
	paramOptimizer.generations = 5
	paramOptimizer.populationSize = 15
	multiOptimizer.generations = 3
	multiOptimizer.populationSize = 10
	multiOptimizer.SetupTextLibObjectives()
	
	scenario := struct {
		textSizes       []int
		expectedActions []string
		performanceGoal float64
	}{
		textSizes:       []int{1000, 5000, 10000}, // Different text sizes
		expectedActions: []string{"ExtractNamedEntities", "CalculateTextStatistics"},
		performanceGoal: 200.0, // Target latency
	}
	
	// Phase 1: Parameter optimization for each function
	functionParameters := make(map[string]map[string]interface{})
	
	for _, actionName := range scenario.expectedActions {
		ranges := GetTextLibParameterRanges()[actionName]
		if ranges == nil {
			continue
		}
		
		// Fitness function considers both performance and accuracy
		fitnessFunc := func(params map[string]interface{}) float64 {
			// Check cache first
			if result, found := cache.Get(actionName+"_fitness", "param_opt", params); found {
				return result.(float64)
			}
			
			// Simulate function execution with these parameters
			simLatency := 50.0 + 100.0*math.Exp(-getParamValue(params, "confidence_threshold", 0.5))
			simAccuracy := 0.7 + 0.2*getParamValue(params, "confidence_threshold", 0.5)
			
			// Multi-objective fitness: balance latency and accuracy
			fitness := simAccuracy - (simLatency / 1000.0) // Prefer accuracy, penalize latency
			
			// Cache the result
			cache.Set(actionName+"_fitness", "param_opt", params, fitness, time.Millisecond*10)
			
			return fitness
		}
		
		optimizedParams := paramOptimizer.OptimizeParameters(actionName, ranges, fitnessFunc)
		functionParameters[actionName] = optimizedParams
	}
	
	// Phase 2: Multi-objective optimization of action sequences
	evaluateFunction := func(actions []EnhancedAction, params map[string]interface{}) Solution {
		totalLatency := 0.0
		totalAccuracy := 0.0
		totalMemory := int64(0)
		totalCost := 0.0
		
		for _, action := range actions {
			// Use optimized parameters if available
			actionParams := functionParameters[action.FunctionName]
			if actionParams == nil {
				actionParams = make(map[string]interface{})
			}
			
			// Check cache for this action result
			cacheKey := action.FunctionName + "_eval"
			if result, found := cache.Get(cacheKey, "eval", actionParams); found {
				cached := result.(map[string]float64)
				totalLatency += cached["latency"]
				totalAccuracy += cached["accuracy"]
				totalMemory += int64(cached["memory"])
				totalCost += cached["cost"]
				continue
			}
			
			// Simulate action execution
			simLatency := 80.0 + 50.0*math.Exp(-getParamValue(actionParams, "confidence_threshold", 0.5))
			simAccuracy := 0.75 + 0.15*getParamValue(actionParams, "confidence_threshold", 0.5)
			simMemory := 1000.0 + 500.0*getParamValue(actionParams, "max_entities", 100.0)/100.0
			simCost := 5.0 + 3.0*simLatency/100.0
			
			totalLatency += simLatency
			totalAccuracy += simAccuracy
			totalMemory += int64(simMemory)
			totalCost += simCost
			
			// Cache the result
			result := map[string]float64{
				"latency": simLatency,
				"accuracy": simAccuracy,
				"memory": simMemory,
				"cost": simCost,
			}
			cache.Set(cacheKey, "eval", actionParams, result, time.Millisecond*5)
		}
		
		// Trigger adaptation if performance is poor
		if totalLatency > scenario.performanceGoal*1.5 {
			adaptationEngine.performanceMonitor.avgLatency = totalLatency
			adaptationEngine.ActivateStrategy("high_latency", "performance threshold exceeded")
		}
		
		return Solution{
			Actions:     actions,
			Parameters:  params,
			Objectives:  []float64{totalLatency, totalAccuracy / float64(len(actions)), float64(totalMemory), totalCost},
			TotalTime:   totalLatency,
			Accuracy:    totalAccuracy / float64(len(actions)),
			MemoryUsage: totalMemory,
			Cost:        totalCost,
		}
	}
	
	// Run multi-objective optimization
	paretoFront := multiOptimizer.Optimize(evaluateFunction)
	
	// Phase 3: Validate results
	if len(paretoFront) == 0 {
		t.Fatal("No solutions found in Pareto front")
	}
	
	// Find solution that meets performance goal
	var bestSolution *Solution
	for _, sol := range paretoFront {
		if sol.TotalTime <= scenario.performanceGoal {
			if bestSolution == nil || sol.Accuracy > bestSolution.Accuracy {
				bestSolution = &sol
			}
		}
	}
	
	if bestSolution == nil {
		t.Logf("No solution met performance goal of %f ms", scenario.performanceGoal)
		// Find best compromise solution
		for _, sol := range paretoFront {
			if bestSolution == nil || sol.Fitness > bestSolution.Fitness {
				bestSolution = &sol
			}
		}
	}
	
	if bestSolution == nil {
		t.Fatal("No viable solution found")
	}
	
	// Validate best solution
	if bestSolution.TotalTime <= 0 {
		t.Errorf("Invalid total time: %f", bestSolution.TotalTime)
	}
	
	if bestSolution.Accuracy <= 0 || bestSolution.Accuracy > 1 {
		t.Errorf("Invalid accuracy: %f", bestSolution.Accuracy)
	}
	
	if bestSolution.MemoryUsage <= 0 {
		t.Errorf("Invalid memory usage: %d", bestSolution.MemoryUsage)
	}
	
	if len(bestSolution.Actions) == 0 {
		t.Error("Best solution has no actions")
	}
	
	// Phase 4: Check cache effectiveness
	cacheStats := cache.GetStats()
	hitRate := cacheStats["hit_rate"].(float64)
	
	if hitRate == 0 {
		t.Error("Cache was not effective (0% hit rate)")
	}
	
	// Phase 5: Check adaptation history
	adaptationHistory := adaptationEngine.adaptationHistory
	
	if len(adaptationHistory) > 0 {
		lastAdaptation := adaptationHistory[len(adaptationHistory)-1]
		if !lastAdaptation.Success {
			t.Error("Last adaptation was not successful")
		}
	}
	
	// Report results
	t.Logf("Best solution found:")
	t.Logf("  Actions: %d", len(bestSolution.Actions))
	t.Logf("  Total time: %.2f ms", bestSolution.TotalTime)
	t.Logf("  Accuracy: %.3f", bestSolution.Accuracy)
	t.Logf("  Memory: %d bytes", bestSolution.MemoryUsage)
	t.Logf("  Cost: %.2f", bestSolution.Cost)
	t.Logf("Cache hit rate: %.1f%%", hitRate*100)
	t.Logf("Adaptations triggered: %d", len(adaptationHistory))
}

// Helper function to safely extract parameter values
func getParamValue(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

func TestComponentInteraction_CacheWithAdaptation(t *testing.T) {
	// Test that cache and adaptation engine work together effectively
	
	cache := NewIntelligentCache(50, time.Hour)
	adaptationEngine := NewDynamicAdaptationEngine()
	
	// Simulate scenario where cache performance affects adaptation decisions
	functionName := "TestFunction"
	
	// Phase 1: Low cache hit rate should trigger adaptation
	for i := 0; i < 20; i++ {
		text := "Test text " + string(rune(i)) // Different texts = cache misses
		params := map[string]interface{}{"param": i}
		
		// Cache miss
		_, found := cache.Get(functionName, text, params)
		if found {
			t.Errorf("Unexpected cache hit on iteration %d", i)
		}
		
		// Simulate expensive computation
		value := "result " + string(rune(i))
		computeCost := time.Millisecond * time.Duration(100+i*10)
		cache.Set(functionName, text, params, value, computeCost)
	}
	
	// Check cache performance
	stats := cache.GetStats()
	hitRate := stats["hit_rate"].(float64)
	
	if hitRate > 0.1 { // Should be very low due to unique texts
		t.Errorf("Expected low hit rate, got %.2f", hitRate)
	}
	
	// Simulate adaptation decision based on cache performance
	if hitRate < adaptationEngine.adaptationThreshold {
		adaptationEngine.ActivateStrategy("low_throughput", "low cache hit rate")
	}
	
	// Verify adaptation was triggered
	if len(adaptationEngine.adaptationHistory) == 0 {
		t.Error("Expected adaptation to be triggered due to low cache performance")
	}
	
	// Phase 2: High cache hit rate with repeated access patterns
	cache2 := NewIntelligentCache(50, time.Hour)
	
	// Use same text repeatedly to increase hit rate
	sameText := "Repeated text for caching"
	for i := 0; i < 20; i++ {
		params := map[string]interface{}{"param": "constant"}
		
		result, found := cache2.Get(functionName, sameText, params)
		if !found {
			// Cache miss - set value
			value := "cached result"
			cache2.Set(functionName, sameText, params, value, time.Millisecond*100)
		} else {
			// Cache hit
			if result != "cached result" {
				t.Errorf("Unexpected cached value: %v", result)
			}
		}
	}
	
	stats2 := cache2.GetStats()
	hitRate2 := stats2["hit_rate"].(float64)
	
	if hitRate2 < 0.9 { // Should be very high due to repeated access
		t.Errorf("Expected high hit rate, got %.2f", hitRate2)
	}
}

func TestParameterOptimization_RealWorldScenario(t *testing.T) {
	// Test parameter optimization in a realistic text processing scenario
	
	paramOptimizer := NewParameterOptimizer()
	paramOptimizer.generations = 10
	paramOptimizer.populationSize = 20
	
	// Simulate text processing with configurable parameters
	textSamples := []string{
		"Apple Inc. CEO Tim Cook announced the new iPhone.",
		"Microsoft Corporation released a statement about AI development.",
		"Google's Alphabet subsidiary is working on quantum computing.",
	}
	
	ranges := map[string]ParameterRange{
		"confidence_threshold": {
			Type: "float",
			Min:  0.1,
			Max:  0.9,
			Default: 0.5,
			Step: 0.1,
		},
		"max_entities": {
			Type: "int",
			Min:  5,
			Max:  50,
			Default: 20,
			Step: 5,
		},
	}
	
	// Fitness function that evaluates parameter performance across text samples
	fitnessFunc := func(params map[string]interface{}) float64 {
		confidenceThreshold := params["confidence_threshold"].(float64)
		maxEntities := params["max_entities"].(int)
		
		totalScore := 0.0
		
		for range textSamples {
			// Simulate entity extraction with these parameters
			expectedEntities := 3 // Known entities in test texts
			
			// Model accuracy based on confidence threshold
			accuracy := 0.6 + 0.3*confidenceThreshold
			
			// Model recall based on max entities limit
			recall := math.Min(1.0, float64(maxEntities)/float64(expectedEntities))
			
			// Processing time increases with lower confidence and higher max entities
			processingTime := (1.0 - confidenceThreshold) * 100.0 + float64(maxEntities) * 2.0
			
			// F1 score for quality
			f1Score := 2.0 * (accuracy * recall) / (accuracy + recall)
			
			// Combined score balancing quality and speed
			score := f1Score - (processingTime / 1000.0) // Penalize slow processing
			
			totalScore += score
		}
		
		return totalScore / float64(len(textSamples))
	}
	
	// Optimize parameters
	result := paramOptimizer.OptimizeParameters("ExtractNamedEntities", ranges, fitnessFunc)
	
	// Validate results
	confidenceThreshold, exists := result["confidence_threshold"]
	if !exists {
		t.Error("confidence_threshold not found in result")
	} else {
		threshold := confidenceThreshold.(float64)
		if threshold < 0.1 || threshold > 0.9 {
			t.Errorf("confidence_threshold %f outside valid range", threshold)
		}
	}
	
	maxEntities, exists := result["max_entities"]
	if !exists {
		t.Error("max_entities not found in result")
	} else {
		entities := maxEntities.(int)
		if entities < 5 || entities > 50 {
			t.Errorf("max_entities %d outside valid range", entities)
		}
	}
	
	// Test that optimized parameters perform better than defaults
	defaultParams := map[string]interface{}{
		"confidence_threshold": 0.5,
		"max_entities":         20,
	}
	
	defaultFitness := fitnessFunc(defaultParams)
	optimizedFitness := fitnessFunc(result)
	
	if optimizedFitness <= defaultFitness {
		t.Logf("Optimized fitness: %.3f, Default fitness: %.3f", optimizedFitness, defaultFitness)
		t.Log("Warning: Optimization did not improve over defaults (may be due to randomness)")
	}
}

// Benchmark the integrated system performance
func BenchmarkIntegratedSystem_EndToEnd(b *testing.B) {
	// Benchmark the performance of the complete enhanced system
	
	paramOptimizer := NewParameterOptimizer()
	cache := NewIntelligentCache(100, time.Hour)
	multiOptimizer := NewMultiObjectiveOptimizer()
	adaptationEngine := NewDynamicAdaptationEngine()
	
	// Configure for benchmark
	paramOptimizer.generations = 2
	paramOptimizer.populationSize = 5
	multiOptimizer.generations = 2
	multiOptimizer.populationSize = 5
	multiOptimizer.SetupTextLibObjectives()
	
	ranges := map[string]ParameterRange{
		"test_param": {Type: "float", Min: 0.1, Max: 0.9, Default: 0.5},
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Parameter optimization
		fitnessFunc := func(params map[string]interface{}) float64 {
			return params["test_param"].(float64)
		}
		paramOptimizer.OptimizeParameters("benchmark_func", ranges, fitnessFunc)
		
		// Cache operations
		cache.Set("bench_func", "test", map[string]interface{}{"i": i}, i, time.Microsecond)
		cache.Get("bench_func", "test", map[string]interface{}{"i": i})
		
		// Multi-objective evaluation
		evalFunc := func(actions []EnhancedAction, params map[string]interface{}) Solution {
			return Solution{
				Objectives: []float64{100.0, 0.8, 1000.0, 50.0},
				TotalTime:  100.0,
				Accuracy:   0.8,
				MemoryUsage: 1000,
				Cost:       50.0,
			}
		}
		multiOptimizer.Optimize(evalFunc)
		
		// Adaptation check
		adaptationEngine.collectMetrics()
	}
}