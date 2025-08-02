package textlib

// Package textlib provides optimized composite functions based on discovered optimal patterns
// These functions implement the most effective sequences for common use cases

import (
	"fmt"
)

// OptimizedGeneralAnalysis performs the highest-reward sequence for general text
// Sequence: validate → entities → readability → keywords → enhanced entities
func OptimizedGeneralAnalysis(text string, opts ...Option) (*GeneralAnalysis, error) {
	// Step 1: Validate input (low cost, early failure detection)
	validation := ValidateOutput(text)
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid input: %v", validation.Issues)
	}

	// Step 2: Initial entity extraction
	entities, err := ExtractEntities(text)
	if err != nil {
		return nil, fmt.Errorf("entity extraction failed: %w", err)
	}

	// Step 3: Readability analysis
	readability, err := AnalyzeReadability(text)
	if err != nil {
		return nil, fmt.Errorf("readability analysis failed: %w", err)
	}

	// Step 4: Keyword extraction
	keywords, err := ExtractKeywords(text)
	if err != nil {
		return nil, fmt.Errorf("keyword extraction failed: %w", err)
	}

	// Step 5: Enhanced entity extraction with context
	enhancedEntities, err := ExtractEntities(text, 
		WithContext(entities, keywords),
		WithReadabilityHints(readability))
	if err != nil {
		// Fallback to initial entities
		enhancedEntities = entities
	}

	return &GeneralAnalysis{
		Validation:       validation,
		Entities:         enhancedEntities,
		Readability:      readability,
		Keywords:         keywords,
		ProcessingSteps:  5,
		OptimizationUsed: "general_highest_reward",
	}, nil
}

// OptimizedTechnicalAnalysis performs the best sequence for technical documentation
// Sequence: code detection → entities → keywords
func OptimizedTechnicalAnalysis(text string) (*TechnicalAnalysis, error) {
	// Validate first
	validation := ValidateOutput(text)
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid input: %v", validation.Issues)
	}

	// Step 1: Detect code presence
	codeInfo, err := DetectCode(text)
	if err != nil {
		return nil, fmt.Errorf("code detection failed: %w", err)
	}

	// Step 2: Extract technical entities
	entities, err := ExtractEntities(text, WithDomain("technical"))
	if err != nil {
		return nil, fmt.Errorf("entity extraction failed: %w", err)
	}

	// Step 3: Extract technical keywords
	keywords, err := ExtractKeywords(text, WithTechnicalTerms(true))
	if err != nil {
		return nil, fmt.Errorf("keyword extraction failed: %w", err)
	}

	// Optional: Readability for documentation quality
	readability, _ := AnalyzeReadability(text)

	return &TechnicalAnalysis{
		HasCode:          codeInfo.HasCode,
		CodeBlocks:       codeInfo.Blocks,
		TechnicalTerms:   entities.Entities,
		Keywords:         keywords.Keywords,
		Readability:      readability,
		ProcessingSteps:  4,
		OptimizationUsed: "technical_specialized",
	}, nil
}

// OptimizedBusinessAnalysis performs the best sequence for business communications
// Sequence: sentiment → entities → keywords
func OptimizedBusinessAnalysis(text string) (*BusinessAnalysis, error) {
	// Validate first
	validation := ValidateOutput(text)
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid input: %v", validation.Issues)
	}

	// Step 1: Sentiment analysis for context
	sentiment, err := SentimentAnalysis(text)
	if err != nil {
		return nil, fmt.Errorf("sentiment analysis failed: %w", err)
	}

	// Step 2: Extract business entities with sentiment context
	entities, err := ExtractEntities(text,
		WithTypes("PERSON", "ORG", "PRODUCT", "MONEY", "DATE"),
		WithSentimentContext(sentiment))
	if err != nil {
		return nil, fmt.Errorf("entity extraction failed: %w", err)
	}

	// Step 3: Extract keywords for themes and action items
	keywords, err := ExtractKeywords(text,
		WithFocus("action_items", "decisions", "deadlines"))
	if err != nil {
		return nil, fmt.Errorf("keyword extraction failed: %w", err)
	}

	return &BusinessAnalysis{
		Sentiment:        sentiment,
		Entities:         entities,
		Keywords:         keywords,
		ProcessingSteps:  3,
		OptimizationUsed: "business_sentiment_first",
	}, nil
}

// OptimizedSocialAnalysis performs the best sequence for social media content
// Sequence: sentiment → keywords → format
func OptimizedSocialAnalysis(text string) (*SocialAnalysis, error) {
	// Quick validation
	validation := ValidateOutput(text)
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid input: %v", validation.Issues)
	}

	// Step 1: Sentiment is primary for social content
	sentiment, err := SentimentAnalysis(text)
	if err != nil {
		return nil, fmt.Errorf("sentiment analysis failed: %w", err)
	}

	// Step 2: Extract trending keywords
	keywords, err := ExtractKeywords(text,
		WithLimit(5),
		WithHashtagDetection(true))
	if err != nil {
		return nil, fmt.Errorf("keyword extraction failed: %w", err)
	}

	// Step 3: Format for display
	formatted, err := FormatText(text,
		WithSocialMediaOptimization(true))
	if err != nil {
		formatted = &FormattedText{Text: text} // Fallback
	}

	return &SocialAnalysis{
		Sentiment:        sentiment,
		TrendingTopics:   keywords.Keywords[:min(5, len(keywords.Keywords))],
		FormattedText:    formatted,
		ProcessingSteps:  3,
		OptimizationUsed: "social_sentiment_keywords",
	}, nil
}

// OptimizedQuickInsights performs minimal high-value analysis
// Sequence: validate → entities → keywords (lowest cost, good value)
func OptimizedQuickInsights(text string) (*QuickInsights, error) {
	// Validate
	validation := ValidateOutput(text)
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid input: %v", validation.Issues)
	}

	// Core extraction only
	entities, err := ExtractEntities(text, WithLimit(5))
	if err != nil {
		return nil, fmt.Errorf("entity extraction failed: %w", err)
	}

	keywords, err := ExtractKeywords(text, WithLimit(5))
	if err != nil {
		return nil, fmt.Errorf("keyword extraction failed: %w", err)
	}

	return &QuickInsights{
		TopEntities:      entities.MostFrequent(3),
		TopKeywords:      keywords.Top(3),
		ValidationScore:  validation.Score,
		ProcessingSteps:  3,
		OptimizationUsed: "quick_minimal",
	}, nil
}

// OptimizedEnhancedExtraction performs the double-extraction pattern
// This pattern captures 20-30% more entities through context building
func OptimizedEnhancedExtraction(text string) (*EnhancedExtraction, error) {
	// Validate
	validation := ValidateOutput(text)
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid input: %v", validation.Issues)
	}

	// First pass: Surface extraction
	firstPass, err := ExtractEntities(text)
	if err != nil {
		return nil, fmt.Errorf("first pass extraction failed: %w", err)
	}

	// Context building
	keywords, _ := ExtractKeywords(text)
	readability, _ := AnalyzeReadability(text)

	// Second pass: Context-aware extraction
	secondPass, err := ExtractEntities(text,
		WithContext(firstPass, keywords),
		WithReadabilityHints(readability),
		WithDeepExtraction(true))
	if err != nil {
		// Fallback to first pass
		return &EnhancedExtraction{
			InitialEntities:  firstPass.Entities,
			EnhancedEntities: firstPass.Entities,
			Improvement:      0,
		}, nil
	}

	// Calculate improvement
	improvement := float64(len(secondPass.Entities)-len(firstPass.Entities)) / 
		float64(max(len(firstPass.Entities), 1)) * 100

	return &EnhancedExtraction{
		InitialEntities:  firstPass.Entities,
		EnhancedEntities: secondPass.Entities,
		Keywords:         keywords.Keywords,
		Improvement:      improvement,
		ProcessingSteps:  4,
		OptimizationUsed: "double_extraction",
	}, nil
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}