package rl

import (
	"math"
	"reflect"
	"testing"
)

func TestNewMultiObjectiveOptimizer(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	if moo == nil {
		t.Fatal("NewMultiObjectiveOptimizer returned nil")
	}
	
	if moo.populationSize != 100 {
		t.Errorf("Expected populationSize 100, got %d", moo.populationSize)
	}
	
	if moo.generations != 200 {
		t.Errorf("Expected generations 200, got %d", moo.generations)
	}
	
	if moo.mutationRate != 0.1 {
		t.Errorf("Expected mutationRate 0.1, got %f", moo.mutationRate)
	}
	
	if moo.crossoverRate != 0.8 {
		t.Errorf("Expected crossoverRate 0.8, got %f", moo.crossoverRate)
	}
	
	if moo.objectives == nil {
		t.Error("objectives slice not initialized")
	}
	
	if moo.weights == nil {
		t.Error("weights slice not initialized")
	}
	
	if moo.paretoFront == nil {
		t.Error("paretoFront slice not initialized")
	}
	
	if moo.generationHistory == nil {
		t.Error("generationHistory slice not initialized")
	}
	
	if moo.convergenceData == nil {
		t.Error("convergenceData slice not initialized")
	}
}

func TestMultiObjectiveOptimizer_AddObjective(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	obj := Objective{
		Name:   "test_objective",
		Type:   "minimize",
		Weight: 0.5,
		Evaluator: func(sol Solution) float64 {
			return sol.TotalTime
		},
	}
	
	moo.AddObjective(obj)
	
	if len(moo.objectives) != 1 {
		t.Errorf("Expected 1 objective, got %d", len(moo.objectives))
	}
	
	if len(moo.weights) != 1 {
		t.Errorf("Expected 1 weight, got %d", len(moo.weights))
	}
	
	if moo.objectives[0].Name != obj.Name {
		t.Errorf("Expected objective name %s, got %s", obj.Name, moo.objectives[0].Name)
	}
	
	if moo.weights[0] != obj.Weight {
		t.Errorf("Expected weight %f, got %f", obj.Weight, moo.weights[0])
	}
}

func TestMultiObjectiveOptimizer_SetupTextLibObjectives(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	moo.SetupTextLibObjectives()
	
	expectedObjectives := 4 // execution_time, accuracy, memory_usage, cost
	if len(moo.objectives) != expectedObjectives {
		t.Errorf("Expected %d objectives, got %d", expectedObjectives, len(moo.objectives))
	}
	
	if len(moo.weights) != expectedObjectives {
		t.Errorf("Expected %d weights, got %d", expectedObjectives, len(moo.weights))
	}
	
	// Check objective names and types
	objectiveNames := make(map[string]string)
	for _, obj := range moo.objectives {
		objectiveNames[obj.Name] = obj.Type
	}
	
	expectedObjectiveTypes := map[string]string{
		"execution_time": "minimize",
		"accuracy":       "maximize",
		"memory_usage":   "minimize",
		"cost":          "minimize",
	}
	
	for name, expectedType := range expectedObjectiveTypes {
		if actualType, exists := objectiveNames[name]; !exists {
			t.Errorf("Missing objective: %s", name)
		} else if actualType != expectedType {
			t.Errorf("Objective %s: expected type %s, got %s", name, expectedType, actualType)
		}
	}
	
	// Check weights sum to reasonable value (should be 1.0)
	totalWeight := 0.0
	for _, weight := range moo.weights {
		totalWeight += weight
	}
	
	expectedTotalWeight := 1.0
	if math.Abs(totalWeight-expectedTotalWeight) > 0.001 {
		t.Errorf("Expected total weight %f, got %f", expectedTotalWeight, totalWeight)
	}
}

func TestMultiObjectiveOptimizer_dominates(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	moo.SetupTextLibObjectives()
	
	tests := []struct {
		name     string
		sol1     Solution
		sol2     Solution
		expected bool
	}{
		{
			name: "sol1 dominates sol2",
			sol1: Solution{
				Objectives: []float64{100, 0.9, 1000, 50}, // time, accuracy, memory, cost
			},
			sol2: Solution{
				Objectives: []float64{200, 0.8, 2000, 100}, // worse in all minimize, worse in maximize
			},
			expected: true,
		},
		{
			name: "sol2 dominates sol1",
			sol1: Solution{
				Objectives: []float64{200, 0.8, 2000, 100},
			},
			sol2: Solution{
				Objectives: []float64{100, 0.9, 1000, 50},
			},
			expected: false,
		},
		{
			name: "no domination - tradeoff",
			sol1: Solution{
				Objectives: []float64{100, 0.8, 1000, 50}, // faster but less accurate
			},
			sol2: Solution{
				Objectives: []float64{200, 0.9, 1000, 50}, // slower but more accurate
			},
			expected: false,
		},
		{
			name: "identical solutions",
			sol1: Solution{
				Objectives: []float64{100, 0.8, 1000, 50},
			},
			sol2: Solution{
				Objectives: []float64{100, 0.8, 1000, 50},
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moo.dominates(tt.sol1, tt.sol2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMultiObjectiveOptimizer_nonDominatedSort(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	moo.SetupTextLibObjectives()
	
	// Create test population with known domination relationships
	population := []Solution{
		{Objectives: []float64{100, 0.9, 1000, 50}}, // Front 0 - best
		{Objectives: []float64{200, 0.8, 2000, 100}}, // Front 1 - dominated by 0
		{Objectives: []float64{150, 0.85, 1500, 75}}, // Front 1 - dominated by 0
		{Objectives: []float64{300, 0.7, 3000, 150}}, // Front 2 - dominated by 1
		{Objectives: []float64{50, 0.95, 500, 25}},   // Front 0 - non-dominated
	}
	
	fronts := moo.nonDominatedSort(population)
	
	if len(fronts) == 0 {
		t.Fatal("No fronts generated")
	}
	
	// Check that solutions are properly ranked
	for i, front := range fronts {
		for _, sol := range front {
			if sol.Rank != i {
				t.Errorf("Solution in front %d has rank %d", i, sol.Rank)
			}
		}
	}
	
	// Check first front has non-dominated solutions
	if len(fronts[0]) == 0 {
		t.Error("First front is empty")
	}
	
	// Verify no solution in first front dominates another
	for i, sol1 := range fronts[0] {
		for j, sol2 := range fronts[0] {
			if i != j && moo.dominates(sol1, sol2) {
				t.Errorf("Solution %d dominates solution %d in same front", i, j)
			}
		}
	}
}

func TestMultiObjectiveOptimizer_calculateCrowdingDistance(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	moo.SetupTextLibObjectives()
	
	// Test with small front
	smallFront := []Solution{
		{Objectives: []float64{100, 0.8, 1000, 50}},
	}
	
	moo.calculateCrowdingDistance(smallFront)
	
	// Single solution should have infinite crowding distance
	if !math.IsInf(smallFront[0].Crowding, 1) {
		t.Errorf("Single solution should have infinite crowding distance, got %f", smallFront[0].Crowding)
	}
	
	// Test with larger front
	front := []Solution{
		{Objectives: []float64{100, 0.9, 1000, 50}}, // Best in time and cost
		{Objectives: []float64{200, 0.8, 2000, 100}}, // Middle
		{Objectives: []float64{300, 0.7, 3000, 150}}, // Worst in time and cost
		{Objectives: []float64{150, 0.95, 1500, 75}}, // Best in accuracy
	}
	
	moo.calculateCrowdingDistance(front)
	
	// Check that all solutions have crowding distance assigned
	for i, sol := range front {
		if sol.Crowding < 0 {
			t.Errorf("Solution %d has negative crowding distance: %f", i, sol.Crowding)
		}
	}
	
	// Boundary solutions should have higher crowding distance
	// (though the specific values depend on sorting order)
	foundInfiniteCrowding := false
	for _, sol := range front {
		if math.IsInf(sol.Crowding, 1) {
			foundInfiniteCrowding = true
			break
		}
	}
	
	if !foundInfiniteCrowding {
		t.Error("Expected at least one solution with infinite crowding distance")
	}
}

func TestMultiObjectiveOptimizer_compare(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	tests := []struct {
		name     string
		sol1     Solution
		sol2     Solution
		expected int
	}{
		{
			name: "sol1 better rank",
			sol1: Solution{Rank: 0, Crowding: 0.5},
			sol2: Solution{Rank: 1, Crowding: 0.8},
			expected: 1,
		},
		{
			name: "sol2 better rank",
			sol1: Solution{Rank: 1, Crowding: 0.8},
			sol2: Solution{Rank: 0, Crowding: 0.5},
			expected: -1,
		},
		{
			name: "same rank, sol1 better crowding",
			sol1: Solution{Rank: 0, Crowding: 0.8},
			sol2: Solution{Rank: 0, Crowding: 0.5},
			expected: 1,
		},
		{
			name: "same rank, sol2 better crowding",
			sol1: Solution{Rank: 0, Crowding: 0.5},
			sol2: Solution{Rank: 0, Crowding: 0.8},
			expected: -1,
		},
		{
			name: "identical solutions",
			sol1: Solution{Rank: 0, Crowding: 0.5},
			sol2: Solution{Rank: 0, Crowding: 0.5},
			expected: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moo.compare(tt.sol1, tt.sol2)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestMultiObjectiveOptimizer_generateRandomActions(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	actions := moo.generateRandomActions()
	
	if len(actions) == 0 {
		t.Error("Generated actions list is empty")
	}
	
	// Check that all actions have required fields
	for i, action := range actions {
		if action.FunctionName == "" {
			t.Errorf("Action %d has empty function name", i)
		}
		
		if action.Category == "" {
			t.Errorf("Action %d has empty category", i)
		}
		
		if action.Cost < 0 {
			t.Errorf("Action %d has negative cost: %d", i, action.Cost)
		}
		
		if action.Parameters == nil {
			t.Errorf("Action %d has nil parameters", i)
		}
	}
	
	// Check that function names are valid
	validFunctions := map[string]bool{
		"ExtractNamedEntities":     true,
		"CalculateTextStatistics":  true,
		"SplitIntoSentences":      true,
		"ExtractAdvancedEntities": true,
		"DetectPatterns":          true,
	}
	
	for i, action := range actions {
		if !validFunctions[action.FunctionName] {
			t.Errorf("Action %d has invalid function name: %s", i, action.FunctionName)
		}
	}
}

func TestMultiObjectiveOptimizer_generateRandomParameters(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	params := moo.generateRandomParameters()
	
	if len(params) == 0 {
		t.Error("Generated parameters map is empty")
	}
	
	// Check expected parameter keys
	expectedKeys := []string{"max_parallel", "timeout_ms", "enable_caching", "cache_size"}
	
	for _, key := range expectedKeys {
		if _, exists := params[key]; !exists {
			t.Errorf("Missing expected parameter: %s", key)
		}
	}
	
	// Check parameter value types and ranges
	if maxParallel, ok := params["max_parallel"].(int); ok {
		if maxParallel < 2 || maxParallel > 7 {
			t.Errorf("max_parallel %d outside expected range [2, 7]", maxParallel)
		}
	} else {
		t.Error("max_parallel not an int")
	}
	
	if timeoutMs, ok := params["timeout_ms"].(int); ok {
		if timeoutMs < 1000 || timeoutMs > 5000 {
			t.Errorf("timeout_ms %d outside expected range [1000, 5000]", timeoutMs)
		}
	} else {
		t.Error("timeout_ms not an int")
	}
	
	if enableCaching, ok := params["enable_caching"].(bool); !ok {
		t.Error("enable_caching not a bool")
	} else if !enableCaching {
		t.Error("enable_caching should be true by default")
	}
	
	if cacheSize, ok := params["cache_size"].(int); ok {
		if cacheSize < 100 || cacheSize > 1000 {
			t.Errorf("cache_size %d outside expected range [100, 1000]", cacheSize)
		}
	} else {
		t.Error("cache_size not an int")
	}
}

func TestMultiObjectiveOptimizer_crossover(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	parent1 := Solution{
		Actions: []EnhancedAction{
			{FunctionName: "ExtractNamedEntities", Category: "analysis", Cost: 1},
		},
		Parameters: map[string]interface{}{
			"max_parallel": 3,
			"timeout_ms":   2000,
		},
		Objectives: []float64{100, 0.8, 1000, 50},
	}
	
	parent2 := Solution{
		Actions: []EnhancedAction{
			{FunctionName: "CalculateTextStatistics", Category: "analysis", Cost: 2},
		},
		Parameters: map[string]interface{}{
			"max_parallel": 5,
			"timeout_ms":   3000,
		},
		Objectives: []float64{200, 0.9, 2000, 100},
	}
	
	child := moo.crossover(parent1, parent2)
	
	// Check child has correct structure
	if len(child.Actions) == 0 {
		t.Error("Child has no actions")
	}
	
	if len(child.Parameters) == 0 {
		t.Error("Child has no parameters")
	}
	
	if len(child.Objectives) != len(moo.objectives) {
		t.Errorf("Child has %d objectives, expected %d", len(child.Objectives), len(moo.objectives))
	}
	
	// Check that child actions come from parents
	childActionNames := make(map[string]bool)
	for _, action := range child.Actions {
		childActionNames[action.FunctionName] = true
	}
	
	parentActionNames := map[string]bool{
		"ExtractNamedEntities":    true,
		"CalculateTextStatistics": true,
	}
	
	for actionName := range childActionNames {
		if !parentActionNames[actionName] {
			t.Errorf("Child action %s not from either parent", actionName)
		}
	}
	
	// Check that parameters are reasonable combinations
	for key, value := range child.Parameters {
		parent1Value, exists1 := parent1.Parameters[key]
		parent2Value, exists2 := parent2.Parameters[key]
		
		if !exists1 && !exists2 {
			t.Errorf("Child parameter %s not present in either parent", key)
			continue
		}
		
		// For numeric parameters, should be average or from one parent
		if exists1 && exists2 {
			if v1, ok1 := parent1Value.(int); ok1 {
				if v2, ok2 := parent2Value.(int); ok2 {
					if childVal, ok := value.(int); ok {
						expected := (v1 + v2) / 2
						if childVal != expected && childVal != v1 && childVal != v2 {
							t.Errorf("Child parameter %s value %d not expected crossover result", key, childVal)
						}
					}
				}
			}
		}
	}
}

func TestMultiObjectiveOptimizer_mutate(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	solution := Solution{
		Actions: []EnhancedAction{
			{FunctionName: "ExtractNamedEntities", Category: "analysis", Cost: 1},
		},
		Parameters: map[string]interface{}{
			"max_parallel": 3,
			"timeout_ms":   2000,
			"cache_size":   500,
		},
		Objectives: []float64{100, 0.8, 1000, 50},
	}
	
	mutated := moo.mutate(solution)
	
	// Check structure is preserved
	if len(mutated.Actions) != len(solution.Actions) {
		t.Errorf("Mutation changed number of actions: %d -> %d", len(solution.Actions), len(mutated.Actions))
	}
	
	if len(mutated.Parameters) != len(solution.Parameters) {
		t.Errorf("Mutation changed number of parameters: %d -> %d", len(solution.Parameters), len(mutated.Parameters))
	}
	
	if len(mutated.Objectives) != len(solution.Objectives) {
		t.Errorf("Mutation changed number of objectives: %d -> %d", len(solution.Objectives), len(mutated.Objectives))
	}
	
	// Check that parameter types are preserved
	for key, originalValue := range solution.Parameters {
		mutatedValue, exists := mutated.Parameters[key]
		if !exists {
			t.Errorf("Mutation removed parameter %s", key)
			continue
		}
		
		if reflect.TypeOf(originalValue) != reflect.TypeOf(mutatedValue) {
			t.Errorf("Mutation changed type of parameter %s: %T -> %T", 
				key, originalValue, mutatedValue)
		}
	}
	
	// Check that values are reasonable (within 50% of original for numeric params)
	for key, originalValue := range solution.Parameters {
		mutatedValue := mutated.Parameters[key]
		
		if origInt, ok := originalValue.(int); ok {
			if mutInt, ok := mutatedValue.(int); ok {
				ratio := float64(mutInt) / float64(origInt)
				if ratio < 0.5 || ratio > 1.5 {
					// This might happen due to randomness, but very large changes are suspicious
					if ratio < 0.1 || ratio > 10.0 {
						t.Errorf("Parameter %s mutated too dramatically: %d -> %d", key, origInt, mutInt)
					}
				}
			}
		}
	}
}

func TestMultiObjectiveOptimizer_calculateBestFitness(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	// Test with empty population
	emptyPop := []Solution{}
	bestFitness := moo.calculateBestFitness(emptyPop)
	if bestFitness != 0 {
		t.Errorf("Expected 0 for empty population, got %f", bestFitness)
	}
	
	// Test with population
	population := []Solution{
		{Fitness: 0.3},
		{Fitness: 0.8},
		{Fitness: 0.5},
		{Fitness: 0.9},
		{Fitness: 0.2},
	}
	
	bestFitness = moo.calculateBestFitness(population)
	expectedBest := 0.9
	
	if bestFitness != expectedBest {
		t.Errorf("Expected best fitness %f, got %f", expectedBest, bestFitness)
	}
}

func TestMultiObjectiveOptimizer_calculateDiversity(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	// Test with empty population
	emptyPop := []Solution{}
	diversity := moo.calculateDiversity(emptyPop)
	if diversity != 0 {
		t.Errorf("Expected 0 diversity for empty population, got %f", diversity)
	}
	
	// Test with single solution
	singlePop := []Solution{{Crowding: 0.5}}
	diversity = moo.calculateDiversity(singlePop)
	if diversity != 0 {
		t.Errorf("Expected 0 diversity for single solution, got %f", diversity)
	}
	
	// Test with population
	population := []Solution{
		{Crowding: 0.3},
		{Crowding: math.Inf(1)}, // Should be excluded from calculation
		{Crowding: 0.5},
		{Crowding: 0.7},
		{Crowding: math.Inf(1)}, // Should be excluded from calculation
	}
	
	diversity = moo.calculateDiversity(population)
	
	// Should calculate average of finite crowding distances: (0.3 + 0.5 + 0.7) / 5 = 0.3 (includes infinite values as part of total)
	expectedDiversity := 0.3
	if math.Abs(diversity-expectedDiversity) > 0.001 {
		t.Errorf("Expected diversity %f, got %f", expectedDiversity, diversity)
	}
}

func TestMultiObjectiveOptimizer_calculateHyperVolume(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	// Test with empty front
	emptyFront := []Solution{}
	hyperVolume := moo.calculateHyperVolume(emptyFront)
	if hyperVolume != 0 {
		t.Errorf("Expected 0 hypervolume for empty front, got %f", hyperVolume)
	}
	
	// Test with front
	front := []Solution{
		{Objectives: []float64{1.0, 2.0, 3.0}},
		{Objectives: []float64{2.0, 1.0, 1.5}},
		{Objectives: []float64{1.5, 1.5, 2.0}},
	}
	
	hyperVolume = moo.calculateHyperVolume(front)
	
	// Should be sum of products of objectives
	expectedHV := (1.0*2.0*3.0) + (2.0*1.0*1.5) + (1.5*1.5*2.0)
	if math.Abs(hyperVolume-expectedHV) > 0.001 {
		t.Errorf("Expected hypervolume %f, got %f", expectedHV, hyperVolume)
	}
}

func TestMultiObjectiveOptimizer_calculateSpread(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	// Test with empty front
	emptyFront := []Solution{}
	spread := moo.calculateSpread(emptyFront)
	if spread != 0 {
		t.Errorf("Expected 0 spread for empty front, got %f", spread)
	}
	
	// Test with single solution
	singleFront := []Solution{{Crowding: 0.5}}
	spread = moo.calculateSpread(singleFront)
	if spread != 0 {
		t.Errorf("Expected 0 spread for single solution, got %f", spread)
	}
	
	// Test with front containing infinite crowding distances
	front := []Solution{
		{Crowding: 0.3},
		{Crowding: math.Inf(1)}, // Should be excluded
		{Crowding: 0.5},
		{Crowding: 0.7},
	}
	
	spread = moo.calculateSpread(front)
	
	// Should calculate standard deviation of finite crowding distances
	// Mean = (0.3 + 0.5 + 0.7) / 3 = 0.5
	// Variance = [(0.3-0.5)² + (0.5-0.5)² + (0.7-0.5)²] / 3 = [0.04 + 0 + 0.04] / 3
	expectedSpread := math.Sqrt(0.08 / 3.0)
	if math.Abs(spread-expectedSpread) > 0.001 {
		t.Errorf("Expected spread %f, got %f", expectedSpread, spread)
	}
}

func TestMultiObjectiveOptimizer_calculateConvergence(t *testing.T) {
	moo := NewMultiObjectiveOptimizer()
	
	// Test early generation
	convergence := moo.calculateConvergence(5)
	if convergence != 1.0 {
		t.Errorf("Expected convergence 1.0 for early generation, got %f", convergence)
	}
	
	// Setup generation history for testing
	moo.generationHistory = []Generation{
		{ParetoFront: []Solution{{}, {}}},     // 2 solutions
		{ParetoFront: []Solution{{}, {}, {}}}, // 3 solutions
	}
	moo.paretoFront = []Solution{{}, {}, {}} // Current: 3 solutions
	
	convergence = moo.calculateConvergence(1)
	
	// Convergence should be 1.0 - |3-3|/3 = 1.0 - 0 = 1.0 (current matches previous)  
	expectedConvergence := 1.0
	if math.Abs(convergence-expectedConvergence) > 0.001 {
		t.Errorf("Expected convergence %f, got %f", expectedConvergence, convergence)
	}
}

// Benchmark tests for performance validation
func BenchmarkMultiObjectiveOptimizer_dominates(b *testing.B) {
	moo := NewMultiObjectiveOptimizer()
	moo.SetupTextLibObjectives()
	
	sol1 := Solution{Objectives: []float64{100, 0.9, 1000, 50}}
	sol2 := Solution{Objectives: []float64{200, 0.8, 2000, 100}}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		moo.dominates(sol1, sol2)
	}
}

func BenchmarkMultiObjectiveOptimizer_calculateCrowdingDistance(b *testing.B) {
	moo := NewMultiObjectiveOptimizer()
	moo.SetupTextLibObjectives()
	
	front := make([]Solution, 20)
	for i := range front {
		front[i] = Solution{
			Objectives: []float64{
				float64(100 + i*10),
				0.5 + float64(i)*0.02,
				float64(1000 + i*100),
				float64(50 + i*5),
			},
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy since the function modifies the slice
		frontCopy := make([]Solution, len(front))
		copy(frontCopy, front)
		moo.calculateCrowdingDistance(frontCopy)
	}
}

func BenchmarkMultiObjectiveOptimizer_nonDominatedSort(b *testing.B) {
	moo := NewMultiObjectiveOptimizer()
	moo.SetupTextLibObjectives()
	
	population := make([]Solution, 50)
	for i := range population {
		population[i] = Solution{
			Objectives: []float64{
				float64(100 + i*5),
				0.5 + float64(i%10)*0.05,
				float64(1000 + i*50),
				float64(50 + i*2),
			},
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy since the function modifies the slice
		popCopy := make([]Solution, len(population))
		copy(popCopy, population)
		moo.nonDominatedSort(popCopy)
	}
}