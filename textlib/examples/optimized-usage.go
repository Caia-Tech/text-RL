package main

import (
	"fmt"
	"textlib-rl-system/textlib"
)

// This example demonstrates how to use the optimized TextLib functions
// for maximum efficiency and quality

func main() {
	// Example texts
	technicalDoc := `The Redis persistence mechanism offers two distinct approaches: RDB (Redis Database) 
snapshots and AOF (Append Only File) logging. RDB performs point-in-time snapshots of your dataset 
at specified intervals, while AOF logs every write operation received by the server.`

	businessEmail := `Hi Team, Following our strategic planning session last week, I'm sharing the 
updated Q3 roadmap priorities. Please review the Customer Portal v2.0 requirements and provide 
feedback by EOD Friday. We need to finalize during Monday's standup. Best regards, Sarah`

	socialPost := `🚀 Excited to announce that our team just shipped the biggest update yet! 
10x faster uploads, end-to-end encryption, and dark mode! Thank you beta testers! 
#ProductLaunch #TechNews`

	// Example 1: Smart Analysis (RL-optimized comprehensive analysis)
	fmt.Println("=== Smart Analysis ===")
	smartResult := textlib.SmartAnalyze(technicalDoc)
	fmt.Printf("Entities found: %d\n", len(smartResult.Entities))
	fmt.Printf("Sentences: %d\n", len(smartResult.Sentences))
	fmt.Printf("Quality score: %.2f\n", smartResult.QualityScore)
	fmt.Printf("Processing path: %v\n", smartResult.OptimizedPath)
	fmt.Printf("Strategy: %s\n", smartResult.Strategy.Name)

	// Example 2: Deep Technical Analysis
	fmt.Println("\n=== Deep Technical Analysis ===")
	techResult := textlib.DeepTechnicalAnalysis(technicalDoc)
	if codeMetrics, ok := techResult.CodeMetrics["code_blocks"].([]string); ok {
		fmt.Printf("Code snippets found: %d\n", len(codeMetrics))
	}
	fmt.Printf("Lexical complexity: %.2f\n", techResult.Complexity.LexicalComplexity)
	fmt.Printf("Overall quality: %.2f\n", techResult.Quality.OverallScore)
	fmt.Printf("Strategy: %s\n", techResult.Strategy.Name)

	// Example 3: Domain-Optimized Analysis
	fmt.Println("\n=== Domain-Optimized Analysis ===")
	businessResult := textlib.DomainOptimizedAnalyze(businessEmail, "general")
	fmt.Printf("Domain: %s\n", businessResult.Domain)
	fmt.Printf("Analysis completed\n")
	if domainSpecific := businessResult.DomainSpecific; domainSpecific != nil {
		fmt.Printf("Domain-specific insights available\n")
	}

	// Example 4: Quick Insights for Social Media
	fmt.Println("\n=== Quick Insights for Social Media ===")
	socialResult := textlib.QuickInsights(socialPost)
	fmt.Printf("Insights: %v\n", socialResult.Insights)
	fmt.Printf("Sentiment score: %.2f\n", socialResult.SentimentScore)
	fmt.Printf("Key terms: %v\n", socialResult.KeyTerms)
	fmt.Printf("Readability: %.2f\n", socialResult.Readability)

	// Example 5: Validated Entity Extraction
	fmt.Println("\n=== Validated Entity Extraction ===")
	validatedResult := textlib.ValidatedExtraction(businessEmail)
	fmt.Printf("Entities found: %d\n", len(validatedResult.Entities))
	fmt.Printf("Validation level: %s\n", validatedResult.ValidationLevel)
	fmt.Printf("Confidence: %.2f\n", validatedResult.Confidence)
	fmt.Printf("Processing path: %v\n", validatedResult.ProcessingPath)

	// Example 6: Performance Comparison
	fmt.Println("\n=== Performance Comparison ===")
	performanceComparison(technicalDoc)
}

// Performance comparison demonstration
func performanceComparison(text string) {
	// Compare different analysis approaches
	
	// Quick insights (fastest)
	fmt.Println("Quick insights (optimized for speed):")
	quickResult := textlib.QuickInsights(text)
	fmt.Printf("  Processing time: %v\n", quickResult.Performance.TotalTime)
	fmt.Printf("  Insights count: %d\n", len(quickResult.Insights))
	
	// Smart analysis (balanced)
	fmt.Println("Smart analysis (balanced quality/speed):")
	smartResult := textlib.SmartAnalyze(text)
	fmt.Printf("  Processing time: %v\n", smartResult.ProcessingTime)
	fmt.Printf("  Quality score: %.2f\n", smartResult.QualityScore)
	
	// Deep technical (most comprehensive)
	fmt.Println("Deep technical analysis (highest quality):")
	deepResult := textlib.DeepTechnicalAnalysis(text)
	fmt.Printf("  Processing time: %v\n", deepResult.Performance.TotalTime)
	fmt.Printf("  Overall quality: %.2f\n", deepResult.Quality.OverallScore)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}