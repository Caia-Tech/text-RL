package benchmarks

import (
	"testing"
	"time"
	
	"github.com/caiatech/textlib"
)

// Better test data that should actually trigger entity recognition
var improvedTestTexts = map[string]string{
	"entities_rich": `Apple Inc. CEO Tim Cook announced on January 15, 2024 that the company will invest $50 billion in artificial intelligence research. The announcement was made at their headquarters in Cupertino, California. Microsoft Corporation and Google LLC are also investing heavily in AI technology.`,
	
	"code_heavy": `package main
import "fmt"
import "os"
import "net/http"

func main() {
	apiKey := "sk-1234567890abcdef"
	password := "admin123"
	server := "https://api.example.com/v1/users"
	fmt.Println("Starting server...")
}`,
	
	"technical_dense": `The PostgreSQL database engine uses MVCC (Multi-Version Concurrency Control) to handle concurrent transactions. Configuration parameters like shared_buffers=256MB, effective_cache_size=4GB, and checkpoint_segments=32 significantly impact performance. Common issues include deadlocks, index bloat, and vacuum scheduling problems.`,
}

// Test if different sequences affect entity extraction quality
func BenchmarkEntityExtractionQuality(b *testing.B) {
	text := improvedTestTexts["entities_rich"]
	
	b.Run("Stats_Then_Entities", func(b *testing.B) {
		var totalEntities int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			// Get statistics first (context)
			_ = textlib.CalculateTextStatistics(text)
			entities := textlib.ExtractNamedEntities(text)
			totalEntities += len(entities)
		}
		
		b.Logf("Stats-first approach - Avg entities: %.2f", float64(totalEntities)/float64(b.N))
	})
	
	b.Run("Direct_Entities", func(b *testing.B) {
		var totalEntities int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			// Direct entity extraction (no context)
			entities := textlib.ExtractNamedEntities(text)
			totalEntities += len(entities)
		}
		
		b.Logf("Direct approach - Avg entities: %.2f", float64(totalEntities)/float64(b.N))
	})
}

// Test memory allocation optimization
func BenchmarkMemoryOptimization(b *testing.B) {
	text := improvedTestTexts["technical_dense"]
	
	b.Run("Sequential_Calls", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Individual calls (potential allocation overhead)
			_ = textlib.CalculateTextStatistics(text)
			_ = textlib.ExtractNamedEntities(text)
			_ = textlib.SplitIntoSentences(text)
		}
	})
	
	b.Run("Batched_Processing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Process in batch (reuse allocations)
			stats := textlib.CalculateTextStatistics(text)
			entities := textlib.ExtractNamedEntities(text)
			sentences := textlib.SplitIntoSentences(text)
			
			// Use results to prevent optimization
			_ = stats
			_ = entities  
			_ = sentences
		}
	})
}

// Test code analysis improvements
func BenchmarkCodeAnalysisImprovement(b *testing.B) {
	code := improvedTestTexts["code_heavy"]
	
	b.Run("Complexity_First", func(b *testing.B) {
		var totalComplexity int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			// Analyze complexity first (may provide context)
			complexity := textlib.CalculateCyclomaticComplexity(code)
			_ = textlib.ExtractFunctionSignatures(code)
			totalComplexity += complexity
		}
		
		b.Logf("Complexity-first - Avg complexity: %.2f", float64(totalComplexity)/float64(b.N))
	})
	
	b.Run("Functions_First", func(b *testing.B) {
		var totalComplexity int
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			// Extract functions first
			_ = textlib.ExtractFunctionSignatures(code)
			complexity := textlib.CalculateCyclomaticComplexity(code)
			totalComplexity += complexity
		}
		
		b.Logf("Functions-first - Avg complexity: %.2f", float64(totalComplexity)/float64(b.N))
	})
}

// Test API responsiveness under load
func BenchmarkAPIResponsiveness(b *testing.B) {
	text := improvedTestTexts["entities_rich"]
	
	b.Run("Single_Call_Latency", func(b *testing.B) {
		var totalDuration time.Duration
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			start := time.Now()
			_ = textlib.ExtractNamedEntities(text)
			totalDuration += time.Since(start)
		}
		
		avgDuration := totalDuration / time.Duration(b.N)
		b.Logf("Single call - Avg duration: %v", avgDuration)
	})
	
	b.Run("Burst_Call_Pattern", func(b *testing.B) {
		var totalDuration time.Duration
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			start := time.Now()
			// Burst of 3 calls
			_ = textlib.ExtractNamedEntities(text)
			_ = textlib.CalculateTextStatistics(text)
			_ = textlib.SplitIntoSentences(text)
			totalDuration += time.Since(start)
		}
		
		avgDuration := totalDuration / time.Duration(b.N)
		b.Logf("Burst pattern - Avg duration: %v", avgDuration)
	})
}

// Test actual function combinations (not just sequences)
func BenchmarkFunctionCombinations(b *testing.B) {
	text := improvedTestTexts["technical_dense"]
	
	b.Run("Full_Analysis", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Full analysis suite
			_ = textlib.CalculateTextStatistics(text)
			_ = textlib.ExtractNamedEntities(text)
			_ = textlib.SplitIntoSentences(text)
		}
	})
	
	b.Run("Minimal_Analysis", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Only essential analysis
			_ = textlib.ExtractNamedEntities(text)
		}
	})
	
	b.Run("Stats_Only", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Statistics only
			_ = textlib.CalculateTextStatistics(text)
		}
	})
}