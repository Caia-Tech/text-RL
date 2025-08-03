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
	"fmt"
	"time"
)

// Core types for RL-optimized API extensions

// ProcessingMetrics tracks performance characteristics of API calls
type ProcessingMetrics struct {
	TimeElapsed    time.Duration `json:"time_elapsed"`
	MemoryPeak     int64         `json:"memory_peak"`
	AlgorithmSteps int           `json:"algorithm_steps"`
	CacheHits      int           `json:"cache_hits"`
	TotalTime      time.Duration `json:"total_time"`
	StepTimings    map[string]time.Duration `json:"step_timings"`
	MemoryUsage    int64         `json:"memory_usage"`
	CacheUtilization float64      `json:"cache_utilization"`
}

// PerformanceMetrics tracks API performance (alias for ProcessingMetrics)
type PerformanceMetrics = ProcessingMetrics

// QualityMetrics represents the quality assessment of results
type QualityMetrics struct {
	Accuracy   float64 `json:"accuracy"`
	Confidence float64 `json:"confidence"`
	Coverage   float64 `json:"coverage"`
}

// ResourceUsage tracks computational resources used
type ResourceUsage struct {
	MemoryUsedMB     int   `json:"memory_used_mb"`
	CPUTimeMs        int64 `json:"cpu_time_ms"`
	NetworkCallsMade int   `json:"network_calls_made"`
	CacheHits        int   `json:"cache_hits"`
}

// ProcessingCost represents the computational cost of an operation
type ProcessingCost struct {
	TimeMs    int64 `json:"time_ms"`
	MemoryKB  int64 `json:"memory_kb"`
	CPUCycles int64 `json:"cpu_cycles"`
}

// ComplexityReport represents detailed text complexity analysis
type ComplexityReport struct {
	LexicalComplexity   float64            `json:"lexical_complexity"`
	SyntacticComplexity float64            `json:"syntactic_complexity"`
	SemanticComplexity  float64            `json:"semantic_complexity"`
	ReadabilityScores   map[string]float64 `json:"readability_scores"`
	ProcessingTime      time.Duration      `json:"processing_time"`
	MemoryUsed          int64              `json:"memory_used"`
	AlgorithmUsed       string             `json:"algorithm_used"`
	QualityMetrics      QualityMetrics     `json:"quality_metrics"`
}

// QualityAssessment represents document quality analysis
type QualityAssessment struct {
	OverallScore       float64          `json:"overall_score"`
	ReadabilityScore   float64          `json:"readability_score"`
	CompletenessScore  float64          `json:"completeness_score"`
	ConsistencyScore   float64          `json:"consistency_score"`
	Issues             []QualityIssue   `json:"issues"`
	Recommendations    []string         `json:"recommendations"`
}

// QualityIssue represents a specific quality problem
type QualityIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"` // low/medium/high
	Description string `json:"description"`
	Location    RLPosition `json:"location"`
	Suggestion  string `json:"suggestion"`
}

// Position represents the location of text within a document
type RLPosition struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Line  int `json:"line"`
}

// ProcessingStrategy represents a strategy for text processing
type ProcessingStrategy struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Parameters       map[string]interface{} `json:"parameters"`
	ExpectedQuality  float64                `json:"expected_quality"`
	ExpectedSpeed    float64                `json:"expected_speed"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements"`
}

// ResourceRequirements specifies resource needs
type ResourceRequirements struct {
	MinMemoryMB      int     `json:"min_memory_mb"`
	MaxMemoryMB      int     `json:"max_memory_mb"`
	EstimatedCPUTime int64   `json:"estimated_cpu_time"`
	NetworkRequired  bool    `json:"network_required"`
	CacheRecommended bool    `json:"cache_recommended"`
}

// TextCharacteristics describes text properties for strategy selection
type TextCharacteristics struct {
	Length      int     `json:"length"`
	Language    string  `json:"language"`
	Domain      string  `json:"domain"`
	Complexity  float64 `json:"complexity"`
	Structure   string  `json:"structure"`
}

// StrategySelector selects optimal processing strategies
type StrategySelector struct {
	// Strategy history for learning
	history []StrategyOutcome
	
	// Learned preferences
	preferences map[string]float64
}

// StrategyOutcome represents the outcome of a strategy selection
type StrategyOutcome struct {
	Context    TextCharacteristics
	Strategy   ProcessingStrategy
	Metrics    OptimizationMetrics
	Successful bool
}

// OptimizationMetrics for RL integration
type OptimizationMetrics struct {
	QualityScore     float64       `json:"quality_score"`     // 0-1, higher is better
	PerformanceScore float64       `json:"performance_score"` // 0-1, higher is better (inverse of time)
	ResourceScore    float64       `json:"resource_score"`    // 0-1, higher is better (inverse of memory)
	UserSatisfaction float64       `json:"user_satisfaction"` // 0-1, based on user feedback
	WeightedTotal    float64       `json:"weighted_total"`    // Combined score
}

// NewStrategySelector creates a new strategy selector
func NewStrategySelector() *StrategySelector {
	return &StrategySelector{
		history:     make([]StrategyOutcome, 0),
		preferences: make(map[string]float64),
	}
}

// SelectStrategy selects the optimal strategy for given text characteristics
func (ss *StrategySelector) SelectStrategy(characteristics TextCharacteristics, requirements AlgorithmRequirements) (ProcessingStrategy, error) {
	// Analyze text characteristics
	strategyType := ss.determineStrategyType(characteristics)
	
	// Build strategy based on requirements
	strategy := ProcessingStrategy{
		Name:        strategyType,
		Description: fmt.Sprintf("Optimized strategy for %s text", characteristics.Domain),
		Parameters:  make(map[string]interface{}),
	}
	
	// Set parameters based on text characteristics
	switch strategyType {
	case "fast":
		strategy.Parameters["depth"] = 1
		strategy.Parameters["algorithms"] = []string{"flesch", "gunning-fog"}
		strategy.Parameters["quality"] = 0.7
		strategy.ExpectedQuality = 0.7
		strategy.ExpectedSpeed = 0.95
		strategy.ResourceRequirements = ResourceRequirements{
			MinMemoryMB:      10,
			MaxMemoryMB:      50,
			EstimatedCPUTime: 50,
			NetworkRequired:  false,
			CacheRecommended: true,
		}
		
	case "balanced":
		strategy.Parameters["depth"] = 2
		strategy.Parameters["algorithms"] = []string{"flesch", "gunning-fog", "coleman-liau", "ari"}
		strategy.Parameters["quality"] = 0.85
		strategy.ExpectedQuality = 0.85
		strategy.ExpectedSpeed = 0.7
		strategy.ResourceRequirements = ResourceRequirements{
			MinMemoryMB:      50,
			MaxMemoryMB:      200,
			EstimatedCPUTime: 200,
			NetworkRequired:  false,
			CacheRecommended: true,
		}
		
	case "comprehensive":
		strategy.Parameters["depth"] = 3
		strategy.Parameters["algorithms"] = []string{"all"}
		strategy.Parameters["quality"] = 0.95
		strategy.ExpectedQuality = 0.95
		strategy.ExpectedSpeed = 0.4
		strategy.ResourceRequirements = ResourceRequirements{
			MinMemoryMB:      100,
			MaxMemoryMB:      500,
			EstimatedCPUTime: 500,
			NetworkRequired:  false,
			CacheRecommended: false, // Too many variations for effective caching
		}
		
	default:
		strategy.Parameters["depth"] = 2
		strategy.Parameters["quality"] = 0.8
		strategy.ExpectedQuality = 0.8
		strategy.ExpectedSpeed = 0.8
	}
	
	return strategy, nil
}

// determineStrategyType determines the base strategy type
func (ss *StrategySelector) determineStrategyType(characteristics TextCharacteristics) string {
	// Simple heuristics for strategy selection
	if characteristics.Length < 100 {
		return "fast"
	}
	
	if characteristics.Length > 10000 || characteristics.Complexity > 0.8 {
		return "comprehensive"
	}
	
	if characteristics.Domain == "technical" || characteristics.Domain == "academic" {
		return "comprehensive"
	}
	
	if characteristics.Domain == "social-media" || characteristics.Domain == "chat" {
		return "fast"
	}
	
	return "balanced"
}