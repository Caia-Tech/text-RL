package rl

import (
	"math"
	"sort"
)

// Multi-objective optimization for balancing performance, accuracy, and resource usage
type MultiObjectiveOptimizer struct {
	objectives []Objective
	weights    []float64
	paretoFront []Solution
	
	// Optimization parameters
	populationSize int
	generations    int
	mutationRate   float64
	crossoverRate  float64
	
	// Performance tracking
	generationHistory []Generation
	convergenceData   []ConvergencePoint
}

type Objective struct {
	Name        string                                           `json:"name"`
	Type        string                                           `json:"type"` // "minimize" or "maximize"
	Weight      float64                                          `json:"weight"`
	Evaluator   func(solution Solution) float64                  `json:"-"`
	Normalizer  func(values []float64) []float64                 `json:"-"`
}

type Solution struct {
	Actions     []EnhancedAction       `json:"actions"`
	Parameters  map[string]interface{} `json:"parameters"`
	Objectives  []float64              `json:"objectives"`
	Fitness     float64                `json:"fitness"`
	Rank        int                    `json:"rank"`
	Crowding    float64                `json:"crowding_distance"`
	
	// Performance metrics
	TotalTime   float64 `json:"total_time"`
	Accuracy    float64 `json:"accuracy"`
	MemoryUsage int64   `json:"memory_usage"`
	CacheHits   int     `json:"cache_hits"`
	Cost        float64 `json:"cost"`
}

type Generation struct {
	Number      int        `json:"number"`
	Population  []Solution `json:"population"`
	ParetoFront []Solution `json:"pareto_front"`
	BestFitness float64    `json:"best_fitness"`
	Diversity   float64    `json:"diversity"`
}

type ConvergencePoint struct {
	Generation     int     `json:"generation"`
	HyperVolume    float64 `json:"hypervolume"`
	Spread         float64 `json:"spread"`
	Convergence    float64 `json:"convergence"`
	ParetoFrontSize int    `json:"pareto_front_size"`
}

func NewMultiObjectiveOptimizer() *MultiObjectiveOptimizer {
	return &MultiObjectiveOptimizer{
		objectives:        make([]Objective, 0),
		weights:          make([]float64, 0),
		paretoFront:      make([]Solution, 0),
		populationSize:   100,
		generations:      200,
		mutationRate:     0.1,
		crossoverRate:    0.8,
		generationHistory: make([]Generation, 0),
		convergenceData:  make([]ConvergencePoint, 0),
	}
}

// Add optimization objectives
func (moo *MultiObjectiveOptimizer) AddObjective(obj Objective) {
	moo.objectives = append(moo.objectives, obj)
	moo.weights = append(moo.weights, obj.Weight)
}

// Set up standard TextLib optimization objectives
func (moo *MultiObjectiveOptimizer) SetupTextLibObjectives() {
	// Objective 1: Minimize execution time
	moo.AddObjective(Objective{
		Name:   "execution_time",
		Type:   "minimize",
		Weight: 0.4,
		Evaluator: func(sol Solution) float64 {
			return sol.TotalTime
		},
	})
	
	// Objective 2: Maximize accuracy/quality
	moo.AddObjective(Objective{
		Name:   "accuracy",
		Type:   "maximize",
		Weight: 0.3,
		Evaluator: func(sol Solution) float64 {
			return sol.Accuracy
		},
	})
	
	// Objective 3: Minimize memory usage
	moo.AddObjective(Objective{
		Name:   "memory_usage",
		Type:   "minimize",
		Weight: 0.2,
		Evaluator: func(sol Solution) float64 {
			return float64(sol.MemoryUsage)
		},
	})
	
	// Objective 4: Minimize cost
	moo.AddObjective(Objective{
		Name:   "cost",
		Type:   "minimize",
		Weight: 0.1,
		Evaluator: func(sol Solution) float64 {
			return sol.Cost
		},
	})
}

// Optimize using NSGA-II algorithm
func (moo *MultiObjectiveOptimizer) Optimize(evaluateFunction func([]EnhancedAction, map[string]interface{}) Solution) []Solution {
	// Initialize population
	population := moo.initializePopulation()
	
	for generation := 0; generation < moo.generations; generation++ {
		// Evaluate population
		for i := range population {
			if len(population[i].Objectives) == 0 || population[i].TotalTime == 0.0 {
				population[i] = evaluateFunction(population[i].Actions, population[i].Parameters)
			}
		}
		
		// Non-dominated sorting
		fronts := moo.nonDominatedSort(population)
		
		// Calculate crowding distance
		for _, front := range fronts {
			moo.calculateCrowdingDistance(front)
		}
		
		// Update Pareto front
		if len(fronts) > 0 {
			moo.paretoFront = make([]Solution, len(fronts[0]))
			copy(moo.paretoFront, fronts[0])
		}
		
		// Record generation data
		gen := Generation{
			Number:      generation,
			Population:  make([]Solution, len(population)),
			ParetoFront: make([]Solution, len(moo.paretoFront)),
			BestFitness: moo.calculateBestFitness(population),
			Diversity:   moo.calculateDiversity(population),
		}
		copy(gen.Population, population)
		copy(gen.ParetoFront, moo.paretoFront)
		moo.generationHistory = append(moo.generationHistory, gen)
		
		// Record convergence data
		convergence := ConvergencePoint{
			Generation:      generation,
			HyperVolume:     moo.calculateHyperVolume(moo.paretoFront),
			Spread:          moo.calculateSpread(moo.paretoFront),
			Convergence:     moo.calculateConvergence(generation),
			ParetoFrontSize: len(moo.paretoFront),
		}
		moo.convergenceData = append(moo.convergenceData, convergence)
		
		// Create next generation
		newPopulation := moo.createNextGeneration(fronts)
		population = newPopulation
	}
	
	return moo.paretoFront
}

func (moo *MultiObjectiveOptimizer) initializePopulation() []Solution {
	population := make([]Solution, moo.populationSize)
	
	for i := 0; i < moo.populationSize; i++ {
		// Generate random solution
		actions := moo.generateRandomActions()
		parameters := moo.generateRandomParameters()
		
		population[i] = Solution{
			Actions:    actions,
			Parameters: parameters,
			Objectives: make([]float64, len(moo.objectives)),
		}
	}
	
	return population
}

func (moo *MultiObjectiveOptimizer) generateRandomActions() []EnhancedAction {
	// Generate a random sequence of actions
	availableActions := []string{
		"ExtractNamedEntities",
		"CalculateTextStatistics", 
		"SplitIntoSentences",
		"ExtractAdvancedEntities",
		"DetectPatterns",
	}
	
	numActions := 1 + (len(availableActions) / 2) // 1 to 3 actions
	actions := make([]EnhancedAction, numActions)
	
	for i := 0; i < numActions; i++ {
		actionName := availableActions[i%len(availableActions)]
		actions[i] = EnhancedAction{
			FunctionName: actionName,
			Category:     "analysis",
			Cost:         1,
			Parameters:   make(map[string]interface{}),
		}
	}
	
	return actions
}

func (moo *MultiObjectiveOptimizer) generateRandomParameters() map[string]interface{} {
	return map[string]interface{}{
		"max_parallel":        2 + (5), // 2-7 parallel operations
		"timeout_ms":         1000 + (4000), // 1-5 second timeout
		"enable_caching":     true,
		"cache_size":         100 + (900), // 100-1000 cache entries
	}
}

// Non-dominated sorting (NSGA-II)
func (moo *MultiObjectiveOptimizer) nonDominatedSort(population []Solution) [][]Solution {
	n := len(population)
	fronts := make([][]Solution, 0)
	dominationCount := make([]int, n)
	dominatedSolutions := make([][]int, n)
	
	// Initialize first front
	firstFront := make([]int, 0)
	
	for i := 0; i < n; i++ {
		dominatedSolutions[i] = make([]int, 0)
		dominationCount[i] = 0
		
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			
			if moo.dominates(population[i], population[j]) {
				dominatedSolutions[i] = append(dominatedSolutions[i], j)
			} else if moo.dominates(population[j], population[i]) {
				dominationCount[i]++
			}
		}
		
		if dominationCount[i] == 0 {
			population[i].Rank = 0
			firstFront = append(firstFront, i)
		}
	}
	
	// Add first front
	if len(firstFront) > 0 {
		front := make([]Solution, len(firstFront))
		for i, idx := range firstFront {
			front[i] = population[idx]
		}
		fronts = append(fronts, front)
	}
	
	// Generate subsequent fronts
	frontIndex := 0
	for len(fronts[frontIndex]) > 0 {
		nextFront := make([]int, 0)
		
		for _, p := range firstFront {
			for _, q := range dominatedSolutions[p] {
				dominationCount[q]--
				if dominationCount[q] == 0 {
					population[q].Rank = frontIndex + 1
					nextFront = append(nextFront, q)
				}
			}
		}
		
		if len(nextFront) > 0 {
			front := make([]Solution, len(nextFront))
			for i, idx := range nextFront {
				front[i] = population[idx]
			}
			fronts = append(fronts, front)
			firstFront = nextFront
			frontIndex++
		} else {
			break
		}
	}
	
	return fronts
}

func (moo *MultiObjectiveOptimizer) dominates(sol1, sol2 Solution) bool {
	better := false
	
	for i, obj := range moo.objectives {
		val1 := sol1.Objectives[i]
		val2 := sol2.Objectives[i]
		
		if obj.Type == "minimize" {
			if val1 > val2 {
				return false
			}
			if val1 < val2 {
				better = true
			}
		} else { // maximize
			if val1 < val2 {
				return false
			}
			if val1 > val2 {
				better = true
			}
		}
	}
	
	return better
}

func (moo *MultiObjectiveOptimizer) calculateCrowdingDistance(front []Solution) {
	if len(front) <= 2 {
		for i := range front {
			front[i].Crowding = math.Inf(1)
		}
		return
	}
	
	// Initialize crowding distance
	for i := range front {
		front[i].Crowding = 0
	}
	
	// Calculate for each objective
	for objIndex := 0; objIndex < len(moo.objectives); objIndex++ {
		// Sort by objective value
		sort.Slice(front, func(i, j int) bool {
			return front[i].Objectives[objIndex] < front[j].Objectives[objIndex]
		})
		
		// Set boundary points to infinity
		front[0].Crowding = math.Inf(1)
		front[len(front)-1].Crowding = math.Inf(1)
		
		// Calculate crowding distance for intermediate points
		if len(front) > 2 {
			objRange := front[len(front)-1].Objectives[objIndex] - front[0].Objectives[objIndex]
			
			if objRange > 0 {
				for i := 1; i < len(front)-1; i++ {
					distance := (front[i+1].Objectives[objIndex] - front[i-1].Objectives[objIndex]) / objRange
					front[i].Crowding += distance
				}
			}
		}
	}
}

func (moo *MultiObjectiveOptimizer) createNextGeneration(fronts [][]Solution) []Solution {
	newPopulation := make([]Solution, 0, moo.populationSize)
	
	// Add fronts until population is full
	for _, front := range fronts {
		if len(newPopulation)+len(front) <= moo.populationSize {
			newPopulation = append(newPopulation, front...)
		} else {
			// Sort by crowding distance and add remaining
			sort.Slice(front, func(i, j int) bool {
				return front[i].Crowding > front[j].Crowding
			})
			
			remaining := moo.populationSize - len(newPopulation)
			newPopulation = append(newPopulation, front[:remaining]...)
			break
		}
	}
	
	// Generate offspring through crossover and mutation
	offspring := moo.generateOffspring(newPopulation)
	
	return offspring
}

func (moo *MultiObjectiveOptimizer) generateOffspring(parents []Solution) []Solution {
	offspring := make([]Solution, moo.populationSize)
	
	for i := 0; i < moo.populationSize; i++ {
		// Tournament selection
		parent1 := moo.tournamentSelection(parents)
		parent2 := moo.tournamentSelection(parents)
		
		// Crossover
		child := moo.crossover(parent1, parent2)
		
		// Mutation
		child = moo.mutate(child)
		
		offspring[i] = child
	}
	
	return offspring
}

func (moo *MultiObjectiveOptimizer) tournamentSelection(population []Solution) Solution {
	tournamentSize := 2
	best := population[0]
	
	for i := 1; i < tournamentSize; i++ {
		candidate := population[i%len(population)]
		if moo.compare(candidate, best) > 0 {
			best = candidate
		}
	}
	
	return best
}

func (moo *MultiObjectiveOptimizer) compare(sol1, sol2 Solution) int {
	if sol1.Rank < sol2.Rank {
		return 1
	}
	if sol1.Rank > sol2.Rank {
		return -1
	}
	if sol1.Crowding > sol2.Crowding {
		return 1
	}
	if sol1.Crowding < sol2.Crowding {
		return -1
	}
	return 0
}

func (moo *MultiObjectiveOptimizer) crossover(parent1, parent2 Solution) Solution {
	child := Solution{
		Actions:    make([]EnhancedAction, 0),
		Parameters: make(map[string]interface{}),
		Objectives: make([]float64, len(moo.objectives)),
	}
	
	// Crossover actions (take from both parents)
	if len(parent1.Actions) > 0 {
		child.Actions = append(child.Actions, parent1.Actions[0])
	}
	if len(parent2.Actions) > 0 {
		child.Actions = append(child.Actions, parent2.Actions[0])
	}
	
	// Crossover parameters
	for key, val1 := range parent1.Parameters {
		if val2, exists := parent2.Parameters[key]; exists {
			// Average numeric values
			if f1, ok := val1.(float64); ok {
				if f2, ok := val2.(float64); ok {
					child.Parameters[key] = (f1 + f2) / 2.0
					continue
				}
			}
			if i1, ok := val1.(int); ok {
				if i2, ok := val2.(int); ok {
					child.Parameters[key] = (i1 + i2) / 2
					continue
				}
			}
		}
		child.Parameters[key] = val1
	}
	
	return child
}

func (moo *MultiObjectiveOptimizer) mutate(solution Solution) Solution {
	// Simple mutation - randomly modify parameters
	for key, val := range solution.Parameters {
		if f, ok := val.(float64); ok {
			solution.Parameters[key] = f * (0.9 + 0.2) // ±10% variation
		}
		if i, ok := val.(int); ok {
			variation := int(float64(i) * 0.1)
			if variation == 0 {
				variation = 1
			}
			solution.Parameters[key] = i + variation - (variation / 2)
		}
	}
	
	return solution
}

// Utility functions for convergence analysis
func (moo *MultiObjectiveOptimizer) calculateBestFitness(population []Solution) float64 {
	if len(population) == 0 {
		return 0
	}
	
	best := population[0].Fitness
	for _, sol := range population[1:] {
		if sol.Fitness > best {
			best = sol.Fitness
		}
	}
	return best
}

func (moo *MultiObjectiveOptimizer) calculateDiversity(population []Solution) float64 {
	if len(population) < 2 {
		return 0
	}
	
	// Calculate average crowding distance
	totalCrowding := 0.0
	for _, sol := range population {
		if !math.IsInf(sol.Crowding, 1) {
			totalCrowding += sol.Crowding
		}
	}
	return totalCrowding / float64(len(population))
}

func (moo *MultiObjectiveOptimizer) calculateHyperVolume(front []Solution) float64 {
	// Simplified hypervolume calculation
	if len(front) == 0 {
		return 0
	}
	
	volume := 0.0
	for _, sol := range front {
		solVolume := 1.0
		for _, objVal := range sol.Objectives {
			solVolume *= objVal
		}
		volume += solVolume
	}
	return volume
}

func (moo *MultiObjectiveOptimizer) calculateSpread(front []Solution) float64 {
	if len(front) < 2 {
		return 0
	}
	
	// Calculate spread as standard deviation of crowding distances
	mean := 0.0
	count := 0
	for _, sol := range front {
		if !math.IsInf(sol.Crowding, 1) {
			mean += sol.Crowding
			count++
		}
	}
	if count == 0 {
		return 0
	}
	mean /= float64(count)
	
	variance := 0.0
	for _, sol := range front {
		if !math.IsInf(sol.Crowding, 1) {
			variance += math.Pow(sol.Crowding-mean, 2)
		}
	}
	variance /= float64(count)
	
	return math.Sqrt(variance)
}

func (moo *MultiObjectiveOptimizer) calculateConvergence(generation int) float64 {
	if generation < 10 {
		return 1.0
	}
	
	// Simple convergence metric based on Pareto front stability
	currentSize := len(moo.paretoFront)
	if generation >= len(moo.generationHistory) {
		return 0.0
	}
	
	prevSize := len(moo.generationHistory[generation-1].ParetoFront)
	if prevSize == 0 {
		return 1.0
	}
	
	return 1.0 - math.Abs(float64(currentSize-prevSize))/float64(prevSize)
}

// Get optimization results and analysis
func (moo *MultiObjectiveOptimizer) GetResults() map[string]interface{} {
	return map[string]interface{}{
		"pareto_front":       moo.paretoFront,
		"generation_history": moo.generationHistory,
		"convergence_data":   moo.convergenceData,
		"final_hypervolume":  moo.calculateHyperVolume(moo.paretoFront),
		"final_spread":       moo.calculateSpread(moo.paretoFront),
		"objectives":         moo.objectives,
	}
}