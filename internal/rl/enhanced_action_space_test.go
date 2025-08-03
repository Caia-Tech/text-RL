package rl

import (
	"reflect"
	"testing"
)

func TestParameterOptimizer_NewParameterOptimizer(t *testing.T) {
	optimizer := NewParameterOptimizer()
	
	if optimizer == nil {
		t.Fatal("NewParameterOptimizer returned nil")
	}
	
	if optimizer.populationSize != 20 {
		t.Errorf("Expected populationSize 20, got %d", optimizer.populationSize)
	}
	
	if optimizer.mutationRate != 0.1 {
		t.Errorf("Expected mutationRate 0.1, got %f", optimizer.mutationRate)
	}
	
	if optimizer.crossoverRate != 0.7 {
		t.Errorf("Expected crossoverRate 0.7, got %f", optimizer.crossoverRate)
	}
	
	if optimizer.generations != 50 {
		t.Errorf("Expected generations 50, got %d", optimizer.generations)
	}
	
	if optimizer.bestParameters == nil {
		t.Error("bestParameters map not initialized")
	}
	
	if optimizer.parameterHistory == nil {
		t.Error("parameterHistory map not initialized")
	}
}

func TestParameterOptimizer_generateRandomValue(t *testing.T) {
	optimizer := NewParameterOptimizer()
	
	tests := []struct {
		name          string
		paramRange    ParameterRange
		expectedType  string
		validateValue func(interface{}) bool
	}{
		{
			name: "integer range",
			paramRange: ParameterRange{
				Type: "int",
				Min:  1,
				Max:  10,
			},
			expectedType: "int",
			validateValue: func(val interface{}) bool {
				if v, ok := val.(int); ok {
					return v >= 1 && v <= 10
				}
				return false
			},
		},
		{
			name: "float range",
			paramRange: ParameterRange{
				Type: "float",
				Min:  0.1,
				Max:  0.9,
			},
			expectedType: "float64",
			validateValue: func(val interface{}) bool {
				if v, ok := val.(float64); ok {
					return v >= 0.1 && v <= 0.9
				}
				return false
			},
		},
		{
			name: "boolean",
			paramRange: ParameterRange{
				Type: "bool",
			},
			expectedType: "bool",
			validateValue: func(val interface{}) bool {
				_, ok := val.(bool)
				return ok
			},
		},
		{
			name: "enum",
			paramRange: ParameterRange{
				Type:    "enum",
				Options: []interface{}{"option1", "option2", "option3"},
			},
			expectedType: "string",
			validateValue: func(val interface{}) bool {
				validOptions := map[string]bool{"option1": true, "option2": true, "option3": true}
				if v, ok := val.(string); ok {
					return validOptions[v]
				}
				return false
			},
		},
		{
			name: "default fallback",
			paramRange: ParameterRange{
				Type:    "unknown",
				Default: "default_value",
			},
			expectedType: "string",
			validateValue: func(val interface{}) bool {
				return val == "default_value"
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := optimizer.generateRandomValue(tt.paramRange)
			
			if reflect.TypeOf(value).String() != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, reflect.TypeOf(value).String())
			}
			
			if !tt.validateValue(value) {
				t.Errorf("Generated value %v failed validation", value)
			}
		})
	}
}

func TestParameterOptimizer_initializePopulation(t *testing.T) {
	optimizer := NewParameterOptimizer()
	
	ranges := map[string]ParameterRange{
		"param1": {
			Type: "int",
			Min:  1,
			Max:  10,
		},
		"param2": {
			Type: "float",
			Min:  0.1,
			Max:  0.9,
		},
		"param3": {
			Type: "bool",
		},
	}
	
	population := optimizer.initializePopulation(ranges)
	
	if len(population) != optimizer.populationSize {
		t.Errorf("Expected population size %d, got %d", optimizer.populationSize, len(population))
	}
	
	// Check each individual
	for i, individual := range population {
		if len(individual) != len(ranges) {
			t.Errorf("Individual %d has %d parameters, expected %d", i, len(individual), len(ranges))
		}
		
		// Check each parameter exists and has correct type
		for paramName, paramRange := range ranges {
			value, exists := individual[paramName]
			if !exists {
				t.Errorf("Individual %d missing parameter %s", i, paramName)
				continue
			}
			
			// Validate parameter type and range
			switch paramRange.Type {
			case "int":
				if v, ok := value.(int); ok {
					min := paramRange.Min.(int)
					max := paramRange.Max.(int)
					if v < min || v > max {
						t.Errorf("Individual %d parameter %s value %d outside range [%d, %d]", i, paramName, v, min, max)
					}
				} else {
					t.Errorf("Individual %d parameter %s not int type", i, paramName)
				}
			case "float":
				if v, ok := value.(float64); ok {
					min := paramRange.Min.(float64)
					max := paramRange.Max.(float64)
					if v < min || v > max {
						t.Errorf("Individual %d parameter %s value %f outside range [%f, %f]", i, paramName, v, min, max)
					}
				} else {
					t.Errorf("Individual %d parameter %s not float64 type", i, paramName)
				}
			case "bool":
				if _, ok := value.(bool); !ok {
					t.Errorf("Individual %d parameter %s not bool type", i, paramName)
				}
			}
		}
	}
}

func TestParameterOptimizer_crossover(t *testing.T) {
	optimizer := NewParameterOptimizer()
	
	parent1 := map[string]interface{}{
		"param1": 5,
		"param2": 0.5,
		"param3": true,
	}
	
	parent2 := map[string]interface{}{
		"param1": 8,
		"param2": 0.8,
		"param3": false,
	}
	
	ranges := map[string]ParameterRange{
		"param1": {Type: "int", Min: 1, Max: 10},
		"param2": {Type: "float", Min: 0.1, Max: 0.9},
		"param3": {Type: "bool"},
	}
	
	child := optimizer.crossover(parent1, parent2, ranges)
	
	// Check child has all parameters
	if len(child) != len(parent1) {
		t.Errorf("Child has %d parameters, expected %d", len(child), len(parent1))
	}
	
	// Check each parameter comes from one of the parents
	for paramName := range ranges {
		childValue, exists := child[paramName]
		if !exists {
			t.Errorf("Child missing parameter %s", paramName)
			continue
		}
		
		parent1Value := parent1[paramName]
		parent2Value := parent2[paramName]
		
		if !reflect.DeepEqual(childValue, parent1Value) && !reflect.DeepEqual(childValue, parent2Value) {
			t.Errorf("Child parameter %s value %v not from either parent (%v, %v)", 
				paramName, childValue, parent1Value, parent2Value)
		}
	}
}

func TestParameterOptimizer_mutate(t *testing.T) {
	optimizer := NewParameterOptimizer()
	optimizer.mutationRate = 1.0 // Force mutation for testing
	
	individual := map[string]interface{}{
		"param1": 5,
		"param2": 0.5,
		"param3": true,
	}
	
	ranges := map[string]ParameterRange{
		"param1": {Type: "int", Min: 1, Max: 10},
		"param2": {Type: "float", Min: 0.1, Max: 0.9},
		"param3": {Type: "bool"},
	}
	
	mutated := optimizer.mutate(individual, ranges)
	
	// Check mutated individual has all parameters
	if len(mutated) != len(individual) {
		t.Errorf("Mutated individual has %d parameters, expected %d", len(mutated), len(individual))
	}
	
	// With 100% mutation rate, at least some parameters should be different
	// (though random generation might occasionally produce same values)
	changedCount := 0
	for paramName := range ranges {
		if !reflect.DeepEqual(mutated[paramName], individual[paramName]) {
			changedCount++
		}
	}
	
	// We can't guarantee changes due to randomness, but let's check types are correct
	for paramName, paramRange := range ranges {
		value := mutated[paramName]
		switch paramRange.Type {
		case "int":
			if _, ok := value.(int); !ok {
				t.Errorf("Mutated parameter %s not int type", paramName)
			}
		case "float":
			if _, ok := value.(float64); !ok {
				t.Errorf("Mutated parameter %s not float64 type", paramName)
			}
		case "bool":
			if _, ok := value.(bool); !ok {
				t.Errorf("Mutated parameter %s not bool type", paramName)
			}
		}
	}
}

func TestParameterOptimizer_tournamentSelection(t *testing.T) {
	optimizer := NewParameterOptimizer()
	
	population := []map[string]interface{}{
		{"value": 1},
		{"value": 2},
		{"value": 3},
		{"value": 4},
		{"value": 5},
	}
	
	fitness := []float64{0.1, 0.2, 0.9, 0.4, 0.5} // index 2 has highest fitness
	
	// Run tournament selection multiple times
	selections := make(map[int]int)
	iterations := 100
	
	for i := 0; i < iterations; i++ {
		selected := optimizer.tournamentSelection(population, fitness)
		
		// Find which individual was selected
		for j, individual := range population {
			if reflect.DeepEqual(selected, individual) {
				selections[j]++
				break
			}
		}
	}
	
	// Individual with highest fitness (index 2) should be selected most often
	if selections[2] == 0 {
		t.Error("Highest fitness individual never selected")
	}
	
	// Should have selected some individuals
	totalSelected := 0
	for _, count := range selections {
		totalSelected += count
	}
	
	if totalSelected != iterations {
		t.Errorf("Expected %d total selections, got %d", iterations, totalSelected)
	}
}

func TestParameterOptimizer_sortByFitness(t *testing.T) {
	optimizer := NewParameterOptimizer()
	
	fitness := []float64{0.1, 0.8, 0.3, 0.9, 0.2}
	expected := []int{3, 1, 2, 4, 0} // Indices sorted by fitness descending
	
	result := optimizer.sortByFitness(fitness)
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	
	// Verify sorting is correct
	for i := 0; i < len(result)-1; i++ {
		currentFitness := fitness[result[i]]
		nextFitness := fitness[result[i+1]]
		if currentFitness < nextFitness {
			t.Errorf("Fitness not sorted descending: %f < %f at positions %d, %d", 
				currentFitness, nextFitness, i, i+1)
		}
	}
}

func TestParameterOptimizer_OptimizeParameters(t *testing.T) {
	optimizer := NewParameterOptimizer()
	optimizer.generations = 5 // Quick test
	optimizer.populationSize = 10
	
	ranges := map[string]ParameterRange{
		"param1": {
			Type: "int",
			Min:  1,
			Max:  10,
		},
		"param2": {
			Type: "float",
			Min:  0.1,
			Max:  0.9,
		},
	}
	
	// Simple fitness function that prefers higher param1 and param2 values
	fitnessFunc := func(params map[string]interface{}) float64 {
		p1 := params["param1"].(int)
		p2 := params["param2"].(float64)
		return float64(p1) + p2
	}
	
	result := optimizer.OptimizeParameters("test_function", ranges, fitnessFunc)
	
	// Check result has correct parameters
	if len(result) != len(ranges) {
		t.Errorf("Result has %d parameters, expected %d", len(result), len(ranges))
	}
	
	for paramName, paramRange := range ranges {
		value, exists := result[paramName]
		if !exists {
			t.Errorf("Result missing parameter %s", paramName)
			continue
		}
		
		// Check parameter type and range
		switch paramRange.Type {
		case "int":
			if v, ok := value.(int); ok {
				min := paramRange.Min.(int)
				max := paramRange.Max.(int)
				if v < min || v > max {
					t.Errorf("Result parameter %s value %d outside range [%d, %d]", paramName, v, min, max)
				}
			} else {
				t.Errorf("Result parameter %s not int type", paramName)
			}
		case "float":
			if v, ok := value.(float64); ok {
				min := paramRange.Min.(float64)
				max := paramRange.Max.(float64)
				if v < min || v > max {
					t.Errorf("Result parameter %s value %f outside range [%f, %f]", paramName, v, min, max)
				}
			} else {
				t.Errorf("Result parameter %s not float64 type", paramName)
			}
		}
	}
	
	// Check that best parameters were stored
	bestParams, exists := optimizer.bestParameters["test_function"]
	if !exists {
		t.Error("Best parameters not stored")
	} else if !reflect.DeepEqual(bestParams, result) {
		t.Error("Stored best parameters don't match returned result")
	}
	
	// Check that parameter history was recorded
	history, exists := optimizer.parameterHistory["test_function"]
	if !exists {
		t.Error("Parameter history not recorded")
	} else if len(history) != optimizer.generations {
		t.Errorf("Expected %d generations in history, got %d", optimizer.generations, len(history))
	}
	
	// Verify fitness generally improves over generations
	if len(history) > 1 {
		firstFitness := history[0].BestFitness
		lastFitness := history[len(history)-1].BestFitness
		
		// Allow for some variance due to randomness, but expect general improvement
		if lastFitness < firstFitness-1.0 {
			t.Errorf("Fitness regressed significantly: %f to %f", firstFitness, lastFitness)
		}
	}
}

func TestGetTextLibParameterRanges(t *testing.T) {
	ranges := GetTextLibParameterRanges()
	
	expectedFunctions := []string{
		"ExtractNamedEntities",
		"CalculateTextStatistics", 
		"SplitIntoSentences",
	}
	
	for _, funcName := range expectedFunctions {
		funcRanges, exists := ranges[funcName]
		if !exists {
			t.Errorf("Missing parameter ranges for function %s", funcName)
			continue
		}
		
		if len(funcRanges) == 0 {
			t.Errorf("Function %s has no parameter ranges", funcName)
		}
		
		// Check each parameter range has required fields
		for paramName, paramRange := range funcRanges {
			if paramRange.Type == "" {
				t.Errorf("Function %s parameter %s missing type", funcName, paramName)
			}
			
			if paramRange.Default == nil {
				t.Errorf("Function %s parameter %s missing default", funcName, paramName)
			}
			
			// Type-specific validation
			switch paramRange.Type {
			case "int", "float":
				if paramRange.Min == nil || paramRange.Max == nil {
					t.Errorf("Function %s parameter %s missing min/max for numeric type", funcName, paramName)
				}
			case "enum":
				if len(paramRange.Options) == 0 {
					t.Errorf("Function %s parameter %s missing options for enum type", funcName, paramName)
				}
			}
		}
	}
}

func TestCopyParams(t *testing.T) {
	original := map[string]interface{}{
		"param1": 42,
		"param2": 3.14,
		"param3": "test",
		"param4": true,
	}
	
	copied := copyParams(original)
	
	// Check all parameters copied
	if len(copied) != len(original) {
		t.Errorf("Copied map has %d parameters, expected %d", len(copied), len(original))
	}
	
	// Check each parameter
	for key, value := range original {
		copiedValue, exists := copied[key]
		if !exists {
			t.Errorf("Copied map missing key %s", key)
			continue
		}
		
		if !reflect.DeepEqual(value, copiedValue) {
			t.Errorf("Key %s: expected %v, got %v", key, value, copiedValue)
		}
	}
	
	// Verify it's a deep copy by modifying original
	original["param1"] = 999
	
	if copied["param1"] == 999 {
		t.Error("Copy is not independent of original")
	}
}

// Benchmark tests for performance validation
func BenchmarkParameterOptimizer_generateRandomValue(b *testing.B) {
	optimizer := NewParameterOptimizer()
	paramRange := ParameterRange{
		Type: "float",
		Min:  0.0,
		Max:  1.0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.generateRandomValue(paramRange)
	}
}

func BenchmarkParameterOptimizer_crossover(b *testing.B) {
	optimizer := NewParameterOptimizer()
	
	parent1 := map[string]interface{}{
		"param1": 5,
		"param2": 0.5,
		"param3": true,
	}
	
	parent2 := map[string]interface{}{
		"param1": 8,
		"param2": 0.8,
		"param3": false,
	}
	
	ranges := map[string]ParameterRange{
		"param1": {Type: "int", Min: 1, Max: 10},
		"param2": {Type: "float", Min: 0.1, Max: 0.9},
		"param3": {Type: "bool"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.crossover(parent1, parent2, ranges)
	}
}

func BenchmarkParameterOptimizer_mutate(b *testing.B) {
	optimizer := NewParameterOptimizer()
	
	individual := map[string]interface{}{
		"param1": 5,
		"param2": 0.5,
		"param3": true,
	}
	
	ranges := map[string]ParameterRange{
		"param1": {Type: "int", Min: 1, Max: 10},
		"param2": {Type: "float", Min: 0.1, Max: 0.9},
		"param3": {Type: "bool"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.mutate(individual, ranges)
	}
}