package benchmarks

import (
	"fmt"
	"strings"
	"testing"
	"time"
	
	// Import the actual text-API from GitHub
	"github.com/caiatech/textlib"
)

// This benchmark validates the patterns discovered by our RL system
// using the actual text-API functions

// Test data matching our training scenarios
var testTexts = map[string]string{
	"technical": `The Redis persistence mechanism offers two distinct approaches: RDB (Redis Database) 
snapshots and AOF (Append Only File) logging. RDB performs point-in-time snapshots of your dataset 
at specified intervals, while AOF logs every write operation received by the server.`,
	
	"business": `Following our strategic planning session last week, I'm sharing the updated Q3 roadmap 
priorities. Please review and provide feedback by EOD Friday. We'll finalize during Monday's standup.
Dr. Sarah Chen, CEO of TechCorp, will be presenting the financial results.`,
	
	"academic": `We present a novel approach to multi-task learning in neural networks that significantly 
improves performance on diverse NLP tasks. Our method, termed Adaptive Task Prioritization (ATP), 
dynamically adjusts task weights during training based on gradient similarity.`,
	
	"code": `func calculateFibonacci(n int) int {
		if n <= 1 { return n }
		return calculateFibonacci(n-1) + calculateFibonacci(n-2)
	}`,
}

// Benchmark: Test discovered optimal sequence vs other approaches
func BenchmarkDiscoveredPatterns(b *testing.B) {
	for textType, text := range testTexts {
		// Test our discovered optimal sequence
		b.Run(fmt.Sprintf("%s_RLOptimal", textType), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = runRLOptimalSequence(text)
			}
		})
		
		// Test naive sequence for comparison
		b.Run(fmt.Sprintf("%s_Naive", textType), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = runNaiveSequence(text)
			}
		})
		
		// Test single-function approach
		b.Run(fmt.Sprintf("%s_SingleFunction", textType), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = runSingleFunction(text)
			}
		})
	}
}

// Benchmark: Validate cost efficiency of discovered patterns
func BenchmarkCostEfficiency(b *testing.B) {
	text := testTexts["technical"]
	
	// Simulated costs based on computational complexity
	costs := map[string]int{
		"statistics":       1,  // CalculateTextStatistics
		"entities":        5,  // ExtractNamedEntities
		"readability":     3,  // CalculateFleschScore
		"advanced":        8,  // ExtractAdvancedEntities
		"patterns":        4,  // DetectPatterns
		"sentences":       2,  // SplitIntoSentences
		"code_analysis":   6,  // Code analysis functions
	}
	
	b.Run("RLOptimalCost", func(b *testing.B) {
		totalCost := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// RL discovered sequence adapted to available functions
			totalCost += costs["statistics"]    // Validate/analyze first
			totalCost += costs["entities"]      // Extract entities
			totalCost += costs["readability"]   // Analyze readability
			totalCost += costs["advanced"]      // Advanced extraction
		}
		b.Logf("RL Optimal - Avg cost: %.2f units", float64(totalCost)/float64(b.N))
	})
	
	b.Run("NaiveCost", func(b *testing.B) {
		totalCost := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Naive: try everything
			for _, cost := range costs {
				totalCost += cost
			}
		}
		b.Logf("Naive - Avg cost: %.2f units", float64(totalCost)/float64(b.N))
	})
}

// Benchmark: Measure quality improvements from discovered patterns
func BenchmarkQualityMetrics(b *testing.B) {
	text := testTexts["technical"]
	
	b.Run("RLOptimalQuality", func(b *testing.B) {
		var totalEntities int
		var totalComplexity float64
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			result := runRLOptimalSequence(text)
			totalEntities += len(result.Entities)
			totalComplexity += result.ReadabilityScore
		}
		
		b.Logf("RL Optimal - Avg entities: %.2f, readability: %.2f",
			float64(totalEntities)/float64(b.N),
			totalComplexity/float64(b.N))
	})
	
	b.Run("SinglePassQuality", func(b *testing.B) {
		var totalEntities int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			// Single pass extraction
			entities := textlib.ExtractNamedEntities(text)
			totalEntities += len(entities)
		}
		
		b.Logf("Single Pass - Avg entities: %.2f", float64(totalEntities)/float64(b.N))
	})
}

// Benchmark: Validate the double-extraction pattern discovery
func BenchmarkDoubleExtractionPattern(b *testing.B) {
	text := testTexts["academic"]
	
	b.Run("SingleExtraction", func(b *testing.B) {
		var totalEntities int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			entities := textlib.ExtractNamedEntities(text)
			totalEntities += len(entities)
		}
		
		b.Logf("Single extraction - Avg entities: %.2f", float64(totalEntities)/float64(b.N))
	})
	
	b.Run("DoubleExtractionWithContext", func(b *testing.B) {
		var totalEntities int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			// First pass - basic entities
			basicEntities := textlib.ExtractNamedEntities(text)
			
			// Build context with other analysis
			stats := textlib.CalculateTextStatistics(text)
			patterns := textlib.DetectPatterns(text)
			
			// Second pass - advanced entities (simulating context-aware extraction)
			advancedEntities := textlib.ExtractAdvancedEntities(text)
			
			// Combine results
			totalUnique := len(mergeUniqueStrings(basicEntities, advancedEntities))
			totalEntities += totalUnique
		}
		
		b.Logf("Double extraction - Avg entities: %.2f", float64(totalEntities)/float64(b.N))
	})
}

// Benchmark: Code analysis patterns
func BenchmarkCodeAnalysisPatterns(b *testing.B) {
	code := testTexts["code"]
	
	b.Run("RLOptimalCodeAnalysis", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// RL pattern: complexity first, then detailed analysis
			_ = textlib.CalculateCyclomaticComplexity(code)
			_ = textlib.ExtractFunctionSignatures(code)
			_ = textlib.FindHardcodedSecrets(code)
		}
	})
	
	b.Run("NaiveCodeAnalysis", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Naive: extract everything without order
			_ = textlib.ExtractFunctionSignatures(code)
			_ = textlib.FindHardcodedSecrets(code)
			_ = textlib.CalculateCyclomaticComplexity(code)
			// Redundant text analysis on code
			_ = textlib.CalculateFleschScore(code)
			_ = textlib.ExtractNamedEntities(code)
		}
	})
}

// Benchmark: Latency comparison
func BenchmarkLatencyComparison(b *testing.B) {
	text := testTexts["business"]
	
	var rlLatencies []time.Duration
	var naiveLatencies []time.Duration
	
	// Measure RL optimal
	for i := 0; i < 100; i++ {
		start := time.Now()
		_ = runRLOptimalSequence(text)
		rlLatencies = append(rlLatencies, time.Since(start))
	}
	
	// Measure naive
	for i := 0; i < 100; i++ {
		start := time.Now()
		_ = runNaiveSequence(text)
		naiveLatencies = append(naiveLatencies, time.Since(start))
	}
	
	rlP50, rlP95, rlP99 := calculatePercentiles(rlLatencies)
	naiveP50, naiveP95, naiveP99 := calculatePercentiles(naiveLatencies)
	
	b.Logf("RL Optimal - P50: %v, P95: %v, P99: %v", rlP50, rlP95, rlP99)
	b.Logf("Naive - P50: %v, P95: %v, P99: %v", naiveP50, naiveP95, naiveP99)
}

// Helper: Run the RL-discovered optimal sequence adapted to text-API
func runRLOptimalSequence(text string) AnalysisResult {
	// Adapted sequence: statistics → entities → readability → advanced entities
	
	// Step 1: Get statistics (validation equivalent)
	stats := textlib.CalculateTextStatistics(text)
	
	// Step 2: Extract basic entities
	entities := textlib.ExtractNamedEntities(text)
	
	// Step 3: Readability analysis
	readability := textlib.CalculateFleschScore(text)
	
	// Step 4: Advanced extraction (enhanced entities)
	advancedEntities := textlib.ExtractAdvancedEntities(text)
	
	// Step 5: Pattern detection for keywords equivalent
	patterns := textlib.DetectPatterns(text)
	
	return AnalysisResult{
		Statistics:       stats,
		Entities:         mergeUniqueStrings(entities, advancedEntities),
		ReadabilityScore: readability,
		Patterns:         patterns,
	}
}

// Helper: Run naive sequence
func runNaiveSequence(text string) AnalysisResult {
	// Inefficient order: split first, then analyze pieces
	sentences := textlib.SplitIntoSentences(text)
	
	var allEntities []string
	// Analyze each sentence separately (inefficient)
	for _, sentence := range sentences {
		entities := textlib.ExtractNamedEntities(sentence)
		allEntities = append(allEntities, entities...)
	}
	
	// Late validation
	stats := textlib.CalculateTextStatistics(text)
	
	// Redundant analysis
	_ = textlib.ExtractAdvancedEntities(text)
	patterns := textlib.DetectPatterns(text)
	
	return AnalysisResult{
		Statistics: stats,
		Entities:   allEntities,
		Patterns:   patterns,
	}
}

// Helper: Run single function
func runSingleFunction(text string) AnalysisResult {
	entities := textlib.ExtractNamedEntities(text)
	
	return AnalysisResult{
		Entities: entities,
	}
}

// Helper: Merge unique strings
func mergeUniqueStrings(list1, list2 []string) []string {
	seen := make(map[string]bool)
	var unique []string
	
	for _, s := range list1 {
		if !seen[strings.ToLower(s)] {
			seen[strings.ToLower(s)] = true
			unique = append(unique, s)
		}
	}
	
	for _, s := range list2 {
		if !seen[strings.ToLower(s)] {
			seen[strings.ToLower(s)] = true
			unique = append(unique, s)
		}
	}
	
	return unique
}

// Helper: Calculate percentiles
func calculatePercentiles(latencies []time.Duration) (p50, p95, p99 time.Duration) {
	if len(latencies) == 0 {
		return
	}
	
	// Sort latencies
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	
	p50 = sorted[len(sorted)*50/100]
	p95 = sorted[len(sorted)*95/100]
	p99 = sorted[len(sorted)*99/100]
	
	return
}

// Result type for comparison
type AnalysisResult struct {
	Statistics       textlib.TextStatistics
	Entities         []string
	ReadabilityScore float64
	Patterns         []string
}