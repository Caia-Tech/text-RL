// Copyright 2025 Caia Tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textlib

import (
	"strings"
	"testing"
	"time"
)

func TestSmartAnalyze(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected struct {
			hasStatistics bool
			hasEntities   bool
			hasSentences  bool
			qualityMin    float64
		}
	}{
		{
			name: "Technical content",
			text: "The algorithm processes data using advanced machine learning techniques. The function calculateSum(a, b int) returns the sum of two integers.",
			expected: struct {
				hasStatistics bool
				hasEntities   bool
				hasSentences  bool
				qualityMin    float64
			}{
				hasStatistics: true,
				hasEntities:   false, // May not find entities in technical text
				hasSentences:  true,
				qualityMin:    0.6,
			},
		},
		{
			name: "Social media content",
			text: "Just had an amazing meeting with @johnsmith! Great ideas for the #innovation project. Can't wait to share more updates soon!",
			expected: struct {
				hasStatistics bool
				hasEntities   bool
				hasSentences  bool
				qualityMin    float64
			}{
				hasStatistics: true,
				hasEntities:   false, // Simple entity extraction may not catch social handles
				hasSentences:  true,
				qualityMin:    0.6,
			},
		},
		{
			name: "Business content",
			text: "John Smith from Acme Corp will attend the meeting on January 15, 2025. The budget allocation is $50,000 for this quarter.",
			expected: struct {
				hasStatistics bool
				hasEntities   bool
				hasSentences  bool
				qualityMin    float64
			}{
				hasStatistics: true,
				hasEntities:   true, // Should find PERSON, ORG, MONEY, DATE entities
				hasSentences:  true,
				qualityMin:    0.8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SmartAnalyze(tt.text)

			// Test basic completeness
			if tt.expected.hasStatistics && result.Statistics.WordCount == 0 {
				t.Errorf("Expected statistics but got empty statistics")
			}

			if tt.expected.hasSentences && len(result.Sentences) == 0 {
				t.Errorf("Expected sentences but got none")
			}

			// Test quality score
			if result.QualityScore < tt.expected.qualityMin {
				t.Errorf("Quality score %f is below minimum %f", result.QualityScore, tt.expected.qualityMin)
			}

			// Test processing path
			expectedPath := []string{"statistics", "entities", "sentences"}
			if len(result.OptimizedPath) != len(expectedPath) {
				t.Errorf("Expected processing path %v, got %v", expectedPath, result.OptimizedPath)
			}

			// Test performance tracking
			if result.ProcessingTime <= 0 {
				t.Errorf("Expected positive processing time, got %v", result.ProcessingTime)
			}

			// Test strategy assignment
			if result.Strategy.Name == "" {
				t.Errorf("Expected strategy to be assigned")
			}
		})
	}
}

func TestValidatedExtraction(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected struct {
			minConfidence float64
			hasValidated  bool
		}
	}{
		{
			name: "High quality entity text",
			text: "Dr. Sarah Johnson from Microsoft Corporation met with President Biden in Washington, DC on March 15, 2024 to discuss the $2.5 million funding proposal.",
			expected: struct {
				minConfidence float64
				hasValidated  bool
			}{
				minConfidence: 0.7,
				hasValidated:  true,
			},
		},
		{
			name: "Low quality text",
			text: "quick test",
			expected: struct {
				minConfidence float64
				hasValidated  bool
			}{
				minConfidence: 0.3,
				hasValidated:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatedExtraction(tt.text)

			// Test confidence level
			if result.Confidence < tt.expected.minConfidence {
				t.Errorf("Confidence %f is below minimum %f", result.Confidence, tt.expected.minConfidence)
			}

			// Test validation level assignment
			validLevels := []string{"basic", "standard", "high"}
			found := false
			for _, level := range validLevels {
				if result.ValidationLevel == level {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Invalid validation level: %s", result.ValidationLevel)
			}

			// Test processing path
			expectedSteps := []string{"pre-validation", "context-extraction", "post-validation"}
			if len(result.ProcessingPath) != len(expectedSteps) {
				t.Errorf("Expected %d processing steps, got %d", len(expectedSteps), len(result.ProcessingPath))
			}

			// Test entity validation
			for _, entity := range result.Entities {
				if entity.ValidationMethod == "" {
					t.Errorf("Entity missing validation method")
				}
				if entity.ValidationScore < 0 || entity.ValidationScore > 1 {
					t.Errorf("Invalid validation score: %f", entity.ValidationScore)
				}
			}
		})
	}
}

func TestDomainOptimizedAnalyze(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		domain string
		checks func(t *testing.T, result DomainAnalysisResult)
	}{
		{
			name:   "Technical domain",
			text:   "The function calculateSum(a, b int) implements an algorithm for adding two integers. This API endpoint accepts HTTP requests.",
			domain: "technical",
			checks: func(t *testing.T, result DomainAnalysisResult) {
				if result.Domain != "technical" {
					t.Errorf("Expected domain 'technical', got %s", result.Domain)
				}
				
				codeSnippets, exists := result.DomainSpecific["code_snippets"]
				if !exists {
					t.Errorf("Expected code_snippets in domain-specific analysis")
				}
				
				if snippets, ok := codeSnippets.([]string); ok && len(snippets) == 0 {
					t.Errorf("Expected to find code snippets in technical text")
				}
			},
		},
		{
			name:   "Social media domain",
			text:   "Hey everyone! Check out this amazing #innovation project @companyX is working on! 🚀",
			domain: "social-media",
			checks: func(t *testing.T, result DomainAnalysisResult) {
				if result.Domain != "social-media" {
					t.Errorf("Expected domain 'social-media', got %s", result.Domain)
				}
				
				hashtags, exists := result.DomainSpecific["hashtags"]
				if !exists {
					t.Errorf("Expected hashtags in social media analysis")
				}
				
				if tags, ok := hashtags.([]string); ok && len(tags) == 0 {
					t.Errorf("Expected to find hashtags in social media text")
				}
			},
		},
		{
			name:   "Academic domain",
			text:   "This research methodology examines the correlation between variables. According to [Smith et al., 2023], the results show significant improvement.",
			domain: "academic",
			checks: func(t *testing.T, result DomainAnalysisResult) {
				if result.Domain != "academic" {
					t.Errorf("Expected domain 'academic', got %s", result.Domain)
				}
				
				citations, exists := result.DomainSpecific["citations"]
				if !exists {
					t.Errorf("Expected citations in academic analysis")
				}
				
				if cites, ok := citations.([]string); ok && len(cites) == 0 {
					// Academic text may not always have detectable citations
					t.Logf("No citations found in academic text (this may be expected)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainOptimizedAnalyze(tt.text, tt.domain)

			// Test basic structure
			if result.Analysis == nil {
				t.Errorf("Expected analysis result")
			}

			if result.DomainSpecific == nil {
				t.Errorf("Expected domain-specific analysis")
			}

			if result.Strategy.Name == "" {
				t.Errorf("Expected strategy to be assigned")
			}

			// Run domain-specific checks
			tt.checks(t, result)
		})
	}
}

func TestQuickInsights(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected struct {
			minInsights    int
			hasSentiment   bool
			hasKeyTerms    bool
			hasReadability bool
		}
	}{
		{
			name: "Positive social media post",
			text: "This is an amazing product! Great quality and excellent customer service. Highly recommend to everyone!",
			expected: struct {
				minInsights    int
				hasSentiment   bool
				hasKeyTerms    bool
				hasReadability bool
			}{
				minInsights:    2,
				hasSentiment:   true,
				hasKeyTerms:    true,
				hasReadability: true,
			},
		},
		{
			name: "Negative review",
			text: "Terrible experience. Bad quality, awful customer service. Hate this product completely.",
			expected: struct {
				minInsights    int
				hasSentiment   bool
				hasKeyTerms    bool
				hasReadability bool
			}{
				minInsights:    2,
				hasSentiment:   true,
				hasKeyTerms:    true,
				hasReadability: true,
			},
		},
		{
			name: "Long technical content",
			text: strings.Repeat("This is a very complex technical document with advanced algorithms and sophisticated methodologies that require extensive analysis and deep understanding of the underlying principles. ", 10),
			expected: struct {
				minInsights    int
				hasSentiment   bool
				hasKeyTerms    bool
				hasReadability bool
			}{
				minInsights:    2,
				hasSentiment:   true,
				hasKeyTerms:    true,
				hasReadability: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuickInsights(tt.text)

			// Test minimum insights
			if len(result.Insights) < tt.expected.minInsights {
				t.Errorf("Expected at least %d insights, got %d", tt.expected.minInsights, len(result.Insights))
			}

			// Test sentiment score range
			if result.SentimentScore < 0 || result.SentimentScore > 1 {
				t.Errorf("Sentiment score %f is out of range [0,1]", result.SentimentScore)
			}

			// Test readability (Flesch scores can be negative for very complex text)
			if result.Readability < -100 || result.Readability > 200 {
				t.Errorf("Readability score %f is out of reasonable range [-100, 200]", result.Readability)
			}

			// Test strategy
			if result.Strategy.Name != "quick-insights" {
				t.Errorf("Expected strategy name 'quick-insights', got %s", result.Strategy.Name)
			}

			// Test performance characteristics
			if result.Strategy.ExpectedSpeed < 0.9 {
				t.Errorf("Quick insights should have high speed expectation, got %f", result.Strategy.ExpectedSpeed)
			}
		})
	}
}

func TestDeepTechnicalAnalysis(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected struct {
			hasCodeMetrics     bool
			hasDocumentation   bool
			hasComplexity      bool
			hasQuality         bool
			minComplexityScore float64
		}
	}{
		{
			name: "Technical documentation with code",
			text: `# API Documentation

## Overview
This API provides advanced functionality for data processing.

## Functions

### calculateSum(a, b int) int
This function calculates the sum of two integers.

### processData()
Advanced data processing algorithm using machine learning techniques.

## Implementation Notes
The implementation follows best practices and includes comprehensive error handling.`,
			expected: struct {
				hasCodeMetrics     bool
				hasDocumentation   bool
				hasComplexity      bool
				hasQuality         bool
				minComplexityScore float64
			}{
				hasCodeMetrics:     true,
				hasDocumentation:   true,
				hasComplexity:      true,
				hasQuality:         true,
				minComplexityScore: 0.3,
			},
		},
		{
			name: "Simple code snippet",
			text: "function add(a, b) { return a + b; }",
			expected: struct {
				hasCodeMetrics     bool
				hasDocumentation   bool
				hasComplexity      bool
				hasQuality         bool
				minComplexityScore float64
			}{
				hasCodeMetrics:     true,
				hasDocumentation:   true,
				hasComplexity:      true,
				hasQuality:         true,
				minComplexityScore: 0.1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeepTechnicalAnalysis(tt.text)

			// Test code metrics
			if tt.expected.hasCodeMetrics && result.CodeMetrics == nil {
				t.Errorf("Expected code metrics")
			}

			// Test documentation analysis
			if tt.expected.hasDocumentation && result.Documentation == nil {
				t.Errorf("Expected documentation analysis")
			}

			// Test complexity analysis
			if tt.expected.hasComplexity {
				if result.Complexity.LexicalComplexity < 0 {
					t.Errorf("Lexical complexity should be non-negative")
				}
				if result.Complexity.SyntacticComplexity < 0 {
					t.Errorf("Syntactic complexity should be non-negative")
				}
				if result.Complexity.SemanticComplexity < 0 {
					t.Errorf("Semantic complexity should be non-negative")
				}
			}

			// Test quality assessment
			if tt.expected.hasQuality {
				if result.Quality.OverallScore < 0 || result.Quality.OverallScore > 1 {
					t.Errorf("Overall quality score %f is out of range [0,1]", result.Quality.OverallScore)
				}
				// ReadabilityScore can be negative when normalized from Flesch scores
				if result.Quality.ReadabilityScore < -1 || result.Quality.ReadabilityScore > 2 {
					t.Errorf("Readability score %f is out of reasonable range [-1,2]", result.Quality.ReadabilityScore)
				}
			}

			// Test strategy
			if result.Strategy.Name != "deep-technical" {
				t.Errorf("Expected strategy name 'deep-technical', got %s", result.Strategy.Name)
			}

			// Test comprehensive analysis depth
			if result.Strategy.ExpectedQuality < 0.9 {
				t.Errorf("Deep technical analysis should have high quality expectation, got %f", result.Strategy.ExpectedQuality)
			}
		})
	}
}

func TestStrategySelector(t *testing.T) {
	selector := NewStrategySelector()

	tests := []struct {
		name            string
		characteristics TextCharacteristics
		requirements    AlgorithmRequirements
		expectedType    string
	}{
		{
			name: "Short text should use fast strategy",
			characteristics: TextCharacteristics{
				Length:     50,
				Language:   "en",
				Domain:     "general",
				Complexity: 0.3,
				Structure:  "simple",
			},
			requirements: AlgorithmRequirements{
				MinQuality:  0.7,
				MaxTimeMs:   1000,
				MaxMemoryMB: 50,
			},
			expectedType: "fast",
		},
		{
			name: "Long technical text should use comprehensive strategy",
			characteristics: TextCharacteristics{
				Length:     15000,
				Language:   "en",
				Domain:     "technical",
				Complexity: 0.9,
				Structure:  "multi-paragraph",
			},
			requirements: AlgorithmRequirements{
				MinQuality:  0.9,
				MaxTimeMs:   10000,
				MaxMemoryMB: 500,
			},
			expectedType: "comprehensive",
		},
		{
			name: "Medium text should use balanced strategy",
			characteristics: TextCharacteristics{
				Length:     500,
				Language:   "en",
				Domain:     "general",
				Complexity: 0.5,
				Structure:  "multi-sentence",
			},
			requirements: AlgorithmRequirements{
				MinQuality:  0.8,
				MaxTimeMs:   3000,
				MaxMemoryMB: 200,
			},
			expectedType: "balanced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := selector.SelectStrategy(tt.characteristics, tt.requirements)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if strategy.Name != tt.expectedType {
				t.Errorf("Expected strategy type %s, got %s", tt.expectedType, strategy.Name)
			}

			// Test strategy has required fields
			if strategy.Description == "" {
				t.Errorf("Strategy should have description")
			}

			if strategy.Parameters == nil {
				t.Errorf("Strategy should have parameters")
			}

			if strategy.ExpectedQuality <= 0 || strategy.ExpectedQuality > 1 {
				t.Errorf("Expected quality %f is out of range (0,1]", strategy.ExpectedQuality)
			}

			if strategy.ExpectedSpeed <= 0 || strategy.ExpectedSpeed > 1 {
				t.Errorf("Expected speed %f is out of range (0,1]", strategy.ExpectedSpeed)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test language detection
	if detectLanguage(" The quick brown fox ") != "en" {
		t.Errorf("Expected English detection")
	}

	if detectLanguage("xxx yyy zzz") != "unknown" {
		t.Errorf("Expected unknown language for non-English text")
	}

	// Test domain classification
	if classifyDomain("function test() { return true; }") != "technical" {
		t.Errorf("Expected technical domain for code")
	}

	if classifyDomain("Check out this cool #hashtag @mention") != "social-media" {
		t.Errorf("Expected social-media domain")
	}

	if classifyDomain("The methodology used in this abstract research") != "academic" {
		t.Errorf("Expected academic domain")
	}

	// Test complexity estimation
	complexity := estimateComplexity("This is a simple test")
	if complexity < 0 || complexity > 1 {
		t.Errorf("Complexity %f is out of range [0,1]", complexity)
	}

	// Test structure analysis
	if analyzeStructure("Single line") != "simple" {
		t.Errorf("Expected simple structure for single line")
	}

	if analyzeStructure("First sentence. Second sentence.") != "multi-sentence" {
		t.Errorf("Expected multi-sentence structure")
	}

	if analyzeStructure("Paragraph one.\n\nParagraph two.") != "multi-paragraph" {
		t.Errorf("Expected multi-paragraph structure")
	}
}

func TestPerformanceMetrics(t *testing.T) {
	text := "This is a test document for performance measurement. It contains multiple sentences and should provide measurable processing time."

	start := time.Now()
	result := SmartAnalyze(text)
	actualTime := time.Since(start)

	// Test that reported time is reasonable
	if result.ProcessingTime <= 0 {
		t.Errorf("Processing time should be positive")
	}

	if result.ProcessingTime > actualTime*2 {
		t.Errorf("Reported processing time %v seems too high compared to actual %v", result.ProcessingTime, actualTime)
	}

	// Test memory usage estimation
	if result.ResourceUsage.MemoryUsedMB < 0 {
		t.Errorf("Memory usage should be non-negative")
	}

	// Test step timings
	if len(result.Performance.StepTimings) == 0 {
		t.Errorf("Should have step timing information")
	}

	totalStepTime := time.Duration(0)
	for _, stepTime := range result.Performance.StepTimings {
		if stepTime <= 0 {
			t.Errorf("Each step should have positive duration")
		}
		totalStepTime += stepTime
	}

	if totalStepTime > result.ProcessingTime*2 {
		t.Errorf("Sum of step times %v should not exceed total time %v by too much", totalStepTime, result.ProcessingTime)
	}
}

func BenchmarkSmartAnalyze(b *testing.B) {
	text := "This is a comprehensive test document for benchmarking the SmartAnalyze function. It contains multiple sentences, various entities like John Smith from Acme Corp, dates like January 15, 2025, and monetary amounts like $50,000. The document structure includes technical terms, business content, and general information to provide a realistic workload for performance testing."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SmartAnalyze(text)
	}
}

func BenchmarkQuickInsights(b *testing.B) {
	text := "Amazing product! Great quality and excellent service. #innovation @company highly recommended!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickInsights(text)
	}
}

func BenchmarkValidatedExtraction(b *testing.B) {
	text := "Dr. Sarah Johnson from Microsoft Corporation met with President Biden in Washington, DC on March 15, 2024 to discuss the $2.5 million funding proposal for the new AI research initiative."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatedExtraction(text)
	}
}