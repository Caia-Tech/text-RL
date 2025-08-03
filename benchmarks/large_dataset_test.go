package benchmarks

import (
	"fmt"
	"strings"
	"testing"
	
	"github.com/caiatech/textlib"
)

// Generate large text datasets for scaling analysis
func generateLargeText(sizeMB float64) string {
	// Base text with entities for realistic testing
	baseText := `Apple Inc. CEO Tim Cook announced on January 15, 2024 that the company will invest $50 billion in artificial intelligence research. The announcement was made at their headquarters in Cupertino, California. Microsoft Corporation and Google LLC are also investing heavily in AI technology. The PostgreSQL database engine uses MVCC (Multi-Version Concurrency Control) to handle concurrent transactions. Configuration parameters like shared_buffers=256MB, effective_cache_size=4GB, and checkpoint_segments=32 significantly impact performance. Common issues include deadlocks, index bloat, and vacuum scheduling problems.`
	
	// Calculate how many repetitions needed for target size
	targetBytes := int(sizeMB * 1024 * 1024)
	baseBytes := len(baseText)
	repetitions := targetBytes / baseBytes
	
	// Create large text by repeating base text
	var builder strings.Builder
	builder.Grow(targetBytes)
	
	for i := 0; i < repetitions; i++ {
		builder.WriteString(baseText)
		if i < repetitions-1 {
			builder.WriteString(" ") // Add spacing between repetitions
		}
	}
	
	return builder.String()
}

// Test performance scaling with different text sizes
func BenchmarkScalingAnalysis(b *testing.B) {
	sizes := []struct {
		name string
		sizeMB float64
	}{
		{"Small_1KB", 0.001},
		{"Medium_10KB", 0.01},
		{"Large_100KB", 0.1},
		{"XLarge_1MB", 1.0},
		{"XXLarge_5MB", 5.0},
	}
	
	for _, size := range sizes {
		text := generateLargeText(size.sizeMB)
		actualSizeMB := float64(len(text)) / (1024 * 1024)
		
		b.Run(size.name+"_Minimal", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Minimal analysis (optimal)
				_ = textlib.ExtractNamedEntities(text)
			}
			b.Logf("Size: %.3f MB, Minimal analysis", actualSizeMB)
		})
		
		b.Run(size.name+"_Full", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Full analysis (inefficient)
				_ = textlib.CalculateTextStatistics(text)
				_ = textlib.ExtractNamedEntities(text)
				_ = textlib.SplitIntoSentences(text)
			}
			b.Logf("Size: %.3f MB, Full analysis", actualSizeMB)
		})
	}
}

// Test memory allocation patterns with large data
func BenchmarkMemoryScaling(b *testing.B) {
	text1MB := generateLargeText(1.0)
	
	b.Run("Minimal_1MB", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			entities := textlib.ExtractNamedEntities(text1MB)
			_ = entities // Use result to prevent optimization
		}
	})
	
	b.Run("Full_1MB", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			stats := textlib.CalculateTextStatistics(text1MB)
			entities := textlib.ExtractNamedEntities(text1MB)
			sentences := textlib.SplitIntoSentences(text1MB)
			
			// Use results to prevent optimization
			_ = stats
			_ = entities
			_ = sentences
		}
	})
}

// Test if optimization benefits increase with text size
func BenchmarkOptimizationScaling(b *testing.B) {
	testSizes := []float64{0.01, 0.1, 1.0, 5.0} // 10KB, 100KB, 1MB, 5MB
	
	for _, sizeMB := range testSizes {
		text := generateLargeText(sizeMB)
		actualSize := float64(len(text)) / (1024 * 1024)
		
		b.Run(fmt.Sprintf("Size_%.1fMB_Minimal", actualSize), func(b *testing.B) {
			var totalEntities int
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				entities := textlib.ExtractNamedEntities(text)
				totalEntities += len(entities)
			}
			
			b.Logf("Size: %.3f MB, Entities: %d", actualSize, totalEntities/b.N)
		})
		
		b.Run(fmt.Sprintf("Size_%.1fMB_WithStats", actualSize), func(b *testing.B) {
			var totalEntities int
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				// Add expensive statistics call
				_ = textlib.CalculateTextStatistics(text)
				entities := textlib.ExtractNamedEntities(text)
				totalEntities += len(entities)
			}
			
			b.Logf("Size: %.3f MB, Entities: %d", actualSize, totalEntities/b.N)
		})
	}
}

// Test entity extraction quality across different text sizes
func BenchmarkQualityScaling(b *testing.B) {
	sizes := []struct {
		name string
		sizeMB float64
	}{
		{"Small", 0.01},   // 10KB
		{"Medium", 0.1},   // 100KB  
		{"Large", 1.0},    // 1MB
	}
	
	for _, size := range sizes {
		text := generateLargeText(size.sizeMB)
		
		b.Run(size.name+"_EntityQuality", func(b *testing.B) {
			var totalEntities int
			var totalUnique int
			entitySet := make(map[string]bool)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entities := textlib.ExtractNamedEntities(text)
				totalEntities += len(entities)
				
				// Count unique entities
				for _, entity := range entities {
					entityStr := entity.Text + ":" + entity.Type
					if !entitySet[entityStr] {
						entitySet[entityStr] = true
						totalUnique++
					}
				}
			}
			
			avgEntities := float64(totalEntities) / float64(b.N)
			b.Logf("Size: %.1f MB, Avg entities: %.1f, Unique: %d", 
				size.sizeMB, avgEntities, totalUnique)
		})
	}
}