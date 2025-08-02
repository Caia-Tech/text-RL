package textlib

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// Benchmark test data
var benchmarkTexts = map[string]string{
	"technical": `The Redis persistence mechanism offers two distinct approaches: RDB (Redis Database) 
snapshots and AOF (Append Only File) logging. RDB performs point-in-time snapshots of your dataset 
at specified intervals, while AOF logs every write operation received by the server. These methods 
can be used independently or combined for maximum data safety. The trade-off involves balancing 
performance impact against data durability requirements. When RDB snapshots are created, Redis 
forks the process using copy-on-write semantics, allowing the parent process to continue serving 
clients while the child process saves the data to disk.`,
	
	"business": `Dear Stakeholders,

Following our Q3 strategic review, I'm pleased to report significant progress across all key 
performance indicators. Revenue increased by 23% year-over-year, exceeding our projections by 
$2.3M. Customer acquisition costs decreased by 15% while maintaining a 94% retention rate.

Key achievements this quarter:
- Launched Customer Portal v2.0 with 10,000+ active users
- Reduced average response time from 48 to 12 hours  
- Expanded our engineering team by 8 senior developers
- Secured Series B funding of $45M at $280M valuation

Challenges ahead include scaling our infrastructure to support projected Q4 growth and maintaining 
service quality during the holiday season. We recommend increasing our DevOps budget by 30% and 
implementing the proposed auto-scaling solution.

Please review the attached detailed metrics and provide feedback by Friday, October 15th.

Best regards,
Sarah Chen
CEO, TechCorp Solutions`,

	"academic": `Abstract: This study investigates the application of transformer-based architectures 
in multi-modal learning environments, specifically focusing on the integration of visual and 
textual data for enhanced semantic understanding. We propose a novel cross-attention mechanism 
that dynamically adjusts feature weights based on contextual relevance, achieving state-of-the-art 
results on three benchmark datasets. Our approach, termed Adaptive Multi-Modal Transformer (AMMT), 
demonstrates a 17.3% improvement in accuracy over baseline models on the COCO-Captions dataset, 
while reducing computational requirements by 34%. Extensive ablation studies reveal that the 
performance gains primarily stem from the proposed relevance-gating mechanism and the hierarchical 
encoding strategy. Furthermore, we provide theoretical analysis showing that our method converges 
to optimal attention distributions under mild assumptions. The implications of this work extend 
to various applications including image captioning, visual question answering, and cross-modal 
retrieval systems.`,

	"social": `🚀 HUGE announcement! After months of hard work, we're thrilled to launch CloudSync 3.0! 

✨ What's new:
- Lightning-fast uploads (10x faster!) ⚡
- End-to-end encryption for maximum security 🔒
- Real-time collaboration with your team 👥
- Beautiful dark mode for late-night work 🌙
- AI-powered file organization 🤖

But that's not all! We're giving away 1TB of free storage to our first 1000 users who sign up 
with code LAUNCH2024! 🎁

Massive thanks to our beta testers, especially @TechGuru and @DevQueen for their invaluable 
feedback. You rock! 🙌

Try it now 👉 cloudsync.io/upgrade

What feature are you most excited about? Let us know below! 👇

#CloudStorage #TechLaunch #StartupLife #Innovation #ProductHunt`,
}

// Benchmark: Optimized patterns vs Naive approach
func BenchmarkOptimizedVsNaive(b *testing.B) {
	testCases := []struct {
		name     string
		textType string
	}{
		{"Technical", "technical"},
		{"Business", "business"},
		{"Academic", "academic"},
		{"Social", "social"},
	}

	for _, tc := range testCases {
		text := benchmarkTexts[tc.textType]

		b.Run(fmt.Sprintf("%s_Optimized", tc.name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = runOptimizedSequence(text, tc.textType)
			}
		})

		b.Run(fmt.Sprintf("%s_Naive", tc.name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = runNaiveSequence(text)
			}
		})

		b.Run(fmt.Sprintf("%s_Random", tc.name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = runRandomSequence(text)
			}
		})
	}
}

// Benchmark: Individual optimized functions
func BenchmarkOptimizedFunctions(b *testing.B) {
	text := benchmarkTexts["technical"]

	b.Run("OptimizedGeneralAnalysis", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = OptimizedGeneralAnalysis(text)
		}
	})

	b.Run("OptimizedTechnicalAnalysis", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = OptimizedTechnicalAnalysis(text)
		}
	})

	b.Run("OptimizedQuickInsights", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = OptimizedQuickInsights(text)
		}
	})
}

// Benchmark: Cost efficiency
func BenchmarkCostEfficiency(b *testing.B) {
	text := benchmarkTexts["business"]

	b.Run("OptimizedCost", func(b *testing.B) {
		totalCost := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cost, _ := calculateOptimizedCost(text)
			totalCost += cost
		}
		b.Logf("Average cost per operation: %.2f", float64(totalCost)/float64(b.N))
	})

	b.Run("NaiveCost", func(b *testing.B) {
		totalCost := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cost, _ := calculateNaiveCost(text)
			totalCost += cost
		}
		b.Logf("Average cost per operation: %.2f", float64(totalCost)/float64(b.N))
	})
}

// Benchmark: Memory usage
func BenchmarkMemoryUsage(b *testing.B) {
	text := benchmarkTexts["academic"]

	b.Run("OptimizedMemory", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = OptimizedGeneralAnalysis(text)
		}
	})

	b.Run("NaiveMemory", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = runNaiveSequence(text)
		}
	})
}

// Benchmark: Accuracy comparison
func BenchmarkAccuracy(b *testing.B) {
	// This would require ground truth data
	// For now, we'll measure completeness
	text := benchmarkTexts["technical"]

	b.Run("OptimizedCompleteness", func(b *testing.B) {
		var totalEntities, totalKeywords int
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, _ := OptimizedGeneralAnalysis(text)
			totalEntities += len(result.Entities.Entities)
			totalKeywords += len(result.Keywords.Keywords)
		}
		b.Logf("Avg entities: %.2f, keywords: %.2f", 
			float64(totalEntities)/float64(b.N),
			float64(totalKeywords)/float64(b.N))
	})

	b.Run("NaiveCompleteness", func(b *testing.B) {
		var totalEntities, totalKeywords int
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, _ := runNaiveSequence(text)
			totalEntities += len(result.Entities)
			totalKeywords += len(result.Keywords)
		}
		b.Logf("Avg entities: %.2f, keywords: %.2f",
			float64(totalEntities)/float64(b.N),
			float64(totalKeywords)/float64(b.N))
	})
}

// Benchmark: Latency percentiles
func BenchmarkLatencyPercentiles(b *testing.B) {
	text := benchmarkTexts["business"]
	latencies := make([]time.Duration, 0, b.N)

	b.Run("OptimizedP99", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			start := time.Now()
			_, _ = OptimizedBusinessAnalysis(text)
			latencies = append(latencies, time.Since(start))
		}
		
		// Calculate percentiles
		p50, p95, p99 := calculatePercentiles(latencies)
		b.Logf("P50: %v, P95: %v, P99: %v", p50, p95, p99)
	})
}

// Helper: Run optimized sequence based on text type
func runOptimizedSequence(text, textType string) (*GeneralAnalysis, error) {
	switch textType {
	case "technical":
		result, err := OptimizedTechnicalAnalysis(text)
		if err != nil {
			return nil, err
		}
		// Convert to GeneralAnalysis for comparison
		return &GeneralAnalysis{
			Entities:    &Entities{Entities: result.TechnicalTerms},
			Keywords:    &Keywords{Keywords: result.Keywords},
			Readability: result.Readability,
		}, nil
	case "business":
		result, err := OptimizedBusinessAnalysis(text)
		if err != nil {
			return nil, err
		}
		return &GeneralAnalysis{
			Entities: result.Entities,
			Keywords: result.Keywords,
		}, nil
	case "social":
		result, err := OptimizedSocialAnalysis(text)
		if err != nil {
			return nil, err
		}
		return &GeneralAnalysis{
			Keywords: &Keywords{Keywords: result.TrendingTopics},
		}, nil
	default:
		return OptimizedGeneralAnalysis(text)
	}
}

// Helper: Run naive sequence (all functions, inefficient order)
func runNaiveSequence(text string) (*NaiveResult, error) {
	// Worst case: summarize first (loses information)
	summary, _ := SummarizeText(text)
	
	// Then try to extract from summary (poor results)
	entities, _ := ExtractEntities(summary.Summary)
	keywords, _ := ExtractKeywords(summary.Summary)
	
	// Redundant calls
	entities2, _ := ExtractEntities(text)
	keywords2, _ := ExtractKeywords(text)
	
	// Late validation (wasted work if invalid)
	validation := ValidateOutput(text)
	
	// Unnecessary formatting
	formatted, _ := FormatText(text)
	
	return &NaiveResult{
		Entities: append(entities.Entities, entities2.Entities...),
		Keywords: append(keywords.Keywords, keywords2.Keywords...),
		Summary:  summary,
		Valid:    validation.IsValid,
		Formatted: formatted,
	}, nil
}

// Helper: Run random sequence
func runRandomSequence(text string) (*GeneralAnalysis, error) {
	functions := []string{
		"validate", "entities", "keywords", "readability",
		"sentiment", "summarize", "format", "detect_code",
	}
	
	// Shuffle functions
	rand.Shuffle(len(functions), func(i, j int) {
		functions[i], functions[j] = functions[j], functions[i]
	})
	
	var result GeneralAnalysis
	
	// Execute in random order
	for _, fn := range functions[:4] { // Only run 4 random functions
		switch fn {
		case "entities":
			entities, _ := ExtractEntities(text)
			result.Entities = entities
		case "keywords":
			keywords, _ := ExtractKeywords(text)
			result.Keywords = keywords
		case "readability":
			readability, _ := AnalyzeReadability(text)
			result.Readability = readability
		case "validate":
			validation := ValidateOutput(text)
			result.Validation = validation
		}
	}
	
	return &result, nil
}

// Helper: Calculate optimized cost
func calculateOptimizedCost(text string) (int, error) {
	// Optimized sequence costs
	costs := map[string]int{
		"validate":    1,
		"entities":    5,
		"readability": 3,
		"keywords":    4,
		"entities2":   5, // Second extraction
	}
	
	total := costs["validate"] + costs["entities"] + 
		costs["readability"] + costs["keywords"] + costs["entities2"]
	
	return total, nil
}

// Helper: Calculate naive cost
func calculateNaiveCost(text string) (int, error) {
	// Naive approach uses everything
	costs := map[string]int{
		"summarize":   8,
		"entities":    5,
		"keywords":    4,
		"entities2":   5,
		"keywords2":   4,
		"validate":    1,
		"format":      2,
		"readability": 3,
		"sentiment":   3,
	}
	
	total := 0
	for _, cost := range costs {
		total += cost
	}
	
	return total, nil
}

// Helper: Calculate percentiles
func calculatePercentiles(latencies []time.Duration) (p50, p95, p99 time.Duration) {
	if len(latencies) == 0 {
		return
	}
	
	// Sort latencies
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	
	// Simple bubble sort for clarity
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

// Helper types for naive approach
type NaiveResult struct {
	Entities  []Entity
	Keywords  []string
	Summary   *Summary
	Valid     bool
	Formatted *FormattedText
}