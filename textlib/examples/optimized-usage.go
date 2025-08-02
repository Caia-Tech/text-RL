package main

import (
	"fmt"
	"log"
	"github.com/yourusername/textlib"
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

	// Example 1: General Purpose Analysis
	fmt.Println("=== General Purpose Analysis ===")
	generalResult, err := textlib.OptimizedGeneralAnalysis(technicalDoc)
	if err != nil {
		log.Printf("General analysis error: %v", err)
	} else {
		fmt.Printf("Entities found: %d\n", len(generalResult.Entities.Entities))
		fmt.Printf("Readability score: %.2f\n", generalResult.Readability.Score)
		fmt.Printf("Top keywords: %v\n", generalResult.Keywords.Top(3))
		fmt.Printf("Optimization used: %s\n", generalResult.OptimizationUsed)
	}

	// Example 2: Technical Documentation Analysis
	fmt.Println("\n=== Technical Documentation Analysis ===")
	techResult, err := textlib.OptimizedTechnicalAnalysis(technicalDoc)
	if err != nil {
		log.Printf("Technical analysis error: %v", err)
	} else {
		fmt.Printf("Has code: %v\n", techResult.HasCode)
		fmt.Printf("Technical terms: %v\n", techResult.TechnicalTerms[:min(5, len(techResult.TechnicalTerms))])
		fmt.Printf("Keywords: %v\n", techResult.Keywords[:min(5, len(techResult.Keywords))])
	}

	// Example 3: Business Communication Analysis
	fmt.Println("\n=== Business Communication Analysis ===")
	businessResult, err := textlib.OptimizedBusinessAnalysis(businessEmail)
	if err != nil {
		log.Printf("Business analysis error: %v", err)
	} else {
		fmt.Printf("Sentiment: %s (score: %.2f)\n", 
			businessResult.Sentiment.Sentiment, 
			businessResult.Sentiment.Score)
		fmt.Printf("People mentioned: %v\n", businessResult.Entities.People)
		fmt.Printf("Action items: %v\n", businessResult.Keywords.Keywords)
	}

	// Example 4: Social Media Analysis
	fmt.Println("\n=== Social Media Analysis ===")
	socialResult, err := textlib.OptimizedSocialAnalysis(socialPost)
	if err != nil {
		log.Printf("Social analysis error: %v", err)
	} else {
		fmt.Printf("Sentiment: %s (confidence: %.2f)\n", 
			socialResult.Sentiment.Sentiment,
			socialResult.Sentiment.Confidence)
		fmt.Printf("Trending topics: %v\n", socialResult.TrendingTopics)
		fmt.Printf("Formatted: %s\n", socialResult.FormattedText.Text)
	}

	// Example 5: Quick Insights (minimal processing)
	fmt.Println("\n=== Quick Insights ===")
	quickResult, err := textlib.OptimizedQuickInsights(technicalDoc)
	if err != nil {
		log.Printf("Quick insights error: %v", err)
	} else {
		fmt.Printf("Top entities: %v\n", quickResult.TopEntities)
		fmt.Printf("Top keywords: %v\n", quickResult.TopKeywords)
		fmt.Printf("Quality score: %.2f\n", quickResult.ValidationScore)
		fmt.Printf("Processing steps: %d (optimized for speed)\n", quickResult.ProcessingSteps)
	}

	// Example 6: Enhanced Extraction (double-pass pattern)
	fmt.Println("\n=== Enhanced Extraction Pattern ===")
	enhancedResult, err := textlib.OptimizedEnhancedExtraction(technicalDoc)
	if err != nil {
		log.Printf("Enhanced extraction error: %v", err)
	} else {
		fmt.Printf("Initial entities: %d\n", len(enhancedResult.InitialEntities))
		fmt.Printf("Enhanced entities: %d\n", len(enhancedResult.EnhancedEntities))
		fmt.Printf("Improvement: %.1f%%\n", enhancedResult.Improvement)
		fmt.Printf("New entities found: %v\n", getNewEntities(
			enhancedResult.InitialEntities, 
			enhancedResult.EnhancedEntities))
	}

	// Example 7: Custom Sequence (when you need specific order)
	fmt.Println("\n=== Custom Sequence Example ===")
	customAnalysis(technicalDoc)

	// Example 8: Budget-Conscious Processing
	fmt.Println("\n=== Budget-Conscious Processing ===")
	budgetAnalysis(businessEmail, 10) // Max 10 cost units
}

// Custom sequence when you need specific processing order
func customAnalysis(text string) {
	// Following discovered optimal pattern: validate → entities → readability → keywords
	
	// Always validate first (cost: 1)
	validation := textlib.ValidateOutput(text)
	if !validation.IsValid {
		log.Printf("Invalid input: %v", validation.Issues)
		return
	}
	
	// Extract entities (cost: 5)
	entities, _ := textlib.ExtractEntities(text)
	
	// Analyze readability (cost: 3)
	readability, _ := textlib.AnalyzeReadability(text)
	
	// Extract keywords with context (cost: 4)
	keywords, _ := textlib.ExtractKeywords(text,
		textlib.WithContext(entities))
	
	fmt.Printf("Custom sequence completed - Total cost: 13\n")
	fmt.Printf("Found %d entities, readability: %.2f, %d keywords\n",
		len(entities.Entities), readability.Score, len(keywords.Keywords))
}

// Budget-conscious processing when you have cost constraints
func budgetAnalysis(text string, maxCost int) {
	costUsed := 0
	
	// Tier 1: Essential validation (cost: 1)
	validation := textlib.ValidateOutput(text)
	costUsed += 1
	if !validation.IsValid {
		return
	}
	
	// Tier 2: Low-cost analysis (cost: 2)
	if costUsed+2 <= maxCost {
		code, _ := textlib.DetectCode(text)
		costUsed += 2
		fmt.Printf("Has code: %v\n", code.HasCode)
	}
	
	// Tier 3: Medium-cost analysis (cost: 3-4)
	if costUsed+3 <= maxCost {
		readability, _ := textlib.AnalyzeReadability(text)
		costUsed += 3
		fmt.Printf("Readability: %.2f\n", readability.Score)
	}
	
	// Tier 4: Higher-cost analysis (cost: 4-5)
	if costUsed+4 <= maxCost {
		keywords, _ := textlib.ExtractKeywords(text, textlib.WithLimit(3))
		costUsed += 4
		fmt.Printf("Top keywords: %v\n", keywords.Keywords)
	}
	
	fmt.Printf("Budget analysis completed - Cost: %d/%d\n", costUsed, maxCost)
}

// Helper function to find new entities
func getNewEntities(initial, enhanced []Entity) []string {
	initialMap := make(map[string]bool)
	for _, e := range initial {
		initialMap[e.Text] = true
	}
	
	var newEntities []string
	for _, e := range enhanced {
		if !initialMap[e.Text] {
			newEntities = append(newEntities, e.Text)
		}
	}
	
	return newEntities
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}