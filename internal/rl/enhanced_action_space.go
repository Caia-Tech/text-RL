package rl

import (
	"math/rand"
)

// Enhanced action space that includes parameter optimization
type EnhancedAction struct {
	FunctionName string                 `json:"function_name"`
	Parameters   map[string]interface{} `json:"parameters"`
	Category     string                 `json:"category"`
	Cost         int                   `json:"cost"`
	
	// Enhanced fields for parameter optimization
	ParameterRanges map[string]ParameterRange `json:"parameter_ranges"`
	Confidence      float64                   `json:"confidence"`
	LastPerformance float64                   `json:"last_performance"`
}

type ParameterRange struct {
	Type     string      `json:"type"`     // "int", "float", "bool", "string", "enum"
	Min      interface{} `json:"min"`      // For numeric types
	Max      interface{} `json:"max"`      // For numeric types
	Options  []interface{} `json:"options"` // For enum types
	Default  interface{} `json:"default"`
	Step     interface{} `json:"step"`     // For discretization
}

// Enhanced state that tracks parameter performance
type EnhancedState struct {
	Text               string                 `json:"text"`
	TaskType           string                 `json:"task_type"`
	ActionsUsed        []EnhancedAction      `json:"actions_used"`
	CurrentResults     map[string]interface{} `json:"current_results"`
	StepCount          int                   `json:"step_count"`
	RemainingBudget    int                   `json:"remaining_budget"`
	
	// Enhanced fields
	TextCharacteristics map[string]float64    `json:"text_characteristics"` // length, complexity, etc.
	PerformanceHistory  []PerformanceMetric   `json:"performance_history"`
	CacheHits          int                   `json:"cache_hits"`
	MemoryPressure     float64               `json:"memory_pressure"`
}

type PerformanceMetric struct {
	Action     string    `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
	Duration   float64   `json:"duration"`
	Quality    float64   `json:"quality"`
	MemoryUsed int64     `json:"memory_used"`
	CacheHit   bool      `json:"cache_hit"`
}

// Parameter optimization using evolutionary strategies
type ParameterOptimizer struct {
	populationSize int
	mutationRate   float64
	crossoverRate  float64
	generations    int
	
	// Track best parameters for each function
	bestParameters map[string]map[string]interface{}
	parameterHistory map[string][]ParameterGeneration
}

type ParameterGeneration struct {
	Generation int                           `json:"generation"`
	Population []map[string]interface{}     `json:"population"`
	Fitness    []float64                    `json:"fitness"`
	BestParams map[string]interface{}       `json:"best_params"`
	BestFitness float64                     `json:"best_fitness"`
}

func NewParameterOptimizer() *ParameterOptimizer {
	return &ParameterOptimizer{
		populationSize:   20,
		mutationRate:     0.1,
		crossoverRate:    0.7,
		generations:      50,
		bestParameters:   make(map[string]map[string]interface{}),
		parameterHistory: make(map[string][]ParameterGeneration),
	}
}

// Generate optimized parameters for a function
func (po *ParameterOptimizer) OptimizeParameters(functionName string, ranges map[string]ParameterRange, 
	fitnessFunc func(params map[string]interface{}) float64) map[string]interface{} {
	
	// Initialize population
	population := po.initializePopulation(ranges)
	
	for generation := 0; generation < po.generations; generation++ {
		// Evaluate fitness
		fitness := make([]float64, len(population))
		for i, individual := range population {
			fitness[i] = fitnessFunc(individual)
		}
		
		// Track best in this generation
		bestIdx := 0
		for i, f := range fitness {
			if f > fitness[bestIdx] {
				bestIdx = i
			}
		}
		
		// Store generation history
		po.parameterHistory[functionName] = append(po.parameterHistory[functionName], ParameterGeneration{
			Generation:  generation,
			Population:  make([]map[string]interface{}, len(population)),
			Fitness:     make([]float64, len(fitness)),
			BestParams:  copyParams(population[bestIdx]),
			BestFitness: fitness[bestIdx],
		})
		copy(po.parameterHistory[functionName][generation].Population, population)
		copy(po.parameterHistory[functionName][generation].Fitness, fitness)
		
		// Create next generation
		newPopulation := po.evolvePopulation(population, fitness, ranges)
		population = newPopulation
	}
	
	// Return best parameters found
	bestParams := po.selectBest(population, fitnessFunc)
	po.bestParameters[functionName] = bestParams
	return bestParams
}

func (po *ParameterOptimizer) initializePopulation(ranges map[string]ParameterRange) []map[string]interface{} {
	population := make([]map[string]interface{}, po.populationSize)
	
	for i := 0; i < po.populationSize; i++ {
		individual := make(map[string]interface{})
		for paramName, paramRange := range ranges {
			individual[paramName] = po.generateRandomValue(paramRange)
		}
		population[i] = individual
	}
	
	return population
}

func (po *ParameterOptimizer) generateRandomValue(paramRange ParameterRange) interface{} {
	switch paramRange.Type {
	case "int":
		min := paramRange.Min.(int)
		max := paramRange.Max.(int)
		return rand.Intn(max-min+1) + min
		
	case "float":
		min := paramRange.Min.(float64)
		max := paramRange.Max.(float64)
		return min + rand.Float64()*(max-min)
		
	case "bool":
		return rand.Float64() < 0.5
		
	case "enum":
		options := paramRange.Options
		return options[rand.Intn(len(options))]
		
	default:
		return paramRange.Default
	}
}

func (po *ParameterOptimizer) evolvePopulation(population []map[string]interface{}, 
	fitness []float64, ranges map[string]ParameterRange) []map[string]interface{} {
	
	newPopulation := make([]map[string]interface{}, po.populationSize)
	
	// Elitism - keep best individuals
	eliteCount := po.populationSize / 10
	indices := po.sortByFitness(fitness)
	for i := 0; i < eliteCount; i++ {
		newPopulation[i] = copyParams(population[indices[i]])
	}
	
	// Generate rest through crossover and mutation
	for i := eliteCount; i < po.populationSize; i++ {
		parent1 := po.tournamentSelection(population, fitness)
		parent2 := po.tournamentSelection(population, fitness)
		
		child := po.crossover(parent1, parent2, ranges)
		child = po.mutate(child, ranges)
		
		newPopulation[i] = child
	}
	
	return newPopulation
}

func (po *ParameterOptimizer) tournamentSelection(population []map[string]interface{}, 
	fitness []float64) map[string]interface{} {
	
	tournamentSize := 3
	bestIdx := rand.Intn(len(population))
	bestFitness := fitness[bestIdx]
	
	for i := 1; i < tournamentSize; i++ {
		idx := rand.Intn(len(population))
		if fitness[idx] > bestFitness {
			bestIdx = idx
			bestFitness = fitness[idx]
		}
	}
	
	return population[bestIdx]
}

func (po *ParameterOptimizer) crossover(parent1, parent2 map[string]interface{}, 
	ranges map[string]ParameterRange) map[string]interface{} {
	
	child := make(map[string]interface{})
	
	for paramName := range ranges {
		if rand.Float64() < po.crossoverRate {
			// Take from parent1
			child[paramName] = parent1[paramName]
		} else {
			// Take from parent2
			child[paramName] = parent2[paramName]
		}
	}
	
	return child
}

func (po *ParameterOptimizer) mutate(individual map[string]interface{}, 
	ranges map[string]ParameterRange) map[string]interface{} {
	
	mutated := copyParams(individual)
	
	for paramName, paramRange := range ranges {
		if rand.Float64() < po.mutationRate {
			mutated[paramName] = po.generateRandomValue(paramRange)
		}
	}
	
	return mutated
}

func (po *ParameterOptimizer) selectBest(population []map[string]interface{}, 
	fitnessFunc func(params map[string]interface{}) float64) map[string]interface{} {
	
	bestIdx := 0
	bestFitness := fitnessFunc(population[0])
	
	for i := 1; i < len(population); i++ {
		fitness := fitnessFunc(population[i])
		if fitness > bestFitness {
			bestIdx = i
			bestFitness = fitness
		}
	}
	
	return population[bestIdx]
}

func (po *ParameterOptimizer) sortByFitness(fitness []float64) []int {
	indices := make([]int, len(fitness))
	for i := range indices {
		indices[i] = i
	}
	
	// Sort indices by fitness (descending)
	for i := 0; i < len(indices)-1; i++ {
		for j := i + 1; j < len(indices); j++ {
			if fitness[indices[j]] > fitness[indices[i]] {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}
	
	return indices
}

func copyParams(params map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range params {
		copy[k] = v
	}
	return copy
}

// Define parameter ranges for TextLib functions
func GetTextLibParameterRanges() map[string]map[string]ParameterRange {
	return map[string]map[string]ParameterRange{
		"ExtractNamedEntities": {
			"confidence_threshold": {
				Type: "float",
				Min:  0.1,
				Max:  0.9,
				Default: 0.5,
				Step: 0.1,
			},
			"max_entities": {
				Type: "int",
				Min:  10,
				Max:  1000,
				Default: 100,
				Step: 10,
			},
		},
		"CalculateTextStatistics": {
			"include_advanced": {
				Type: "bool",
				Default: false,
			},
			"sample_size": {
				Type: "int",
				Min:  100,
				Max:  10000,
				Default: 1000,
				Step: 100,
			},
		},
		"SplitIntoSentences": {
			"min_sentence_length": {
				Type: "int",
				Min:  5,
				Max:  50,
				Default: 10,
				Step: 5,
			},
			"delimiter_style": {
				Type: "enum",
				Options: []interface{}{"standard", "aggressive", "conservative"},
				Default: "standard",
			},
		},
	}
}