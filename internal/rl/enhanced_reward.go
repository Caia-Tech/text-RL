package rl

import (
	"math"
	"strings"
)

// EnhancedRewardCalculator provides sophisticated reward calculations based on task context
type EnhancedRewardCalculator struct {
	TaskWeights       map[string]float64
	QualityThresholds map[string]float64
	SequenceBonus     map[string]float64
}

func NewEnhancedRewardCalculator() *EnhancedRewardCalculator {
	return &EnhancedRewardCalculator{
		TaskWeights: map[string]float64{
			"technical_analysis":      1.2,
			"code_analysis":          1.3,
			"academic_analysis":      1.1,
			"business_communication": 1.0,
			"news_analysis":         0.9,
			"social_media_analysis": 0.8,
			"legal_analysis":        1.4,
			"medical_analysis":      1.5,
			"instructional_analysis": 0.7,
			"log_analysis":          1.2,
			"marketing_analysis":    0.9,
			"scientific_analysis":   1.3,
		},
		QualityThresholds: map[string]float64{
			"entity_extraction":    0.7,
			"readability_analysis": 0.8,
			"code_detection":      0.9,
			"sentiment_analysis":  0.75,
			"keyword_extraction":  0.8,
		},
		SequenceBonus: map[string]float64{
			"extract_entities->extract_keywords":     0.3,
			"detect_code->analyze_readability":      0.2,
			"extract_keywords->sentiment_analysis":  0.25,
			"analyze_readability->summarize_text":   0.4,
			"extract_entities->validate_output":     0.2,
		},
	}
}

func (erc *EnhancedRewardCalculator) CalculateReward(state State, action Action, result ActionResult, example TrainingExample) float64 {
	if !result.Success {
		// Penalize failures based on action cost
		return -1.0 - (float64(action.Cost) * 0.1)
	}
	
	baseReward := 1.0
	
	// 1. Task-specific weight
	if weight, exists := erc.TaskWeights[example.TaskType]; exists {
		baseReward *= weight
	}
	
	// 2. Action relevance to task
	relevanceBonus := erc.calculateRelevanceBonus(action, example)
	
	// 3. Output quality assessment
	qualityScore := erc.assessOutputQuality(action, result, example)
	
	// 4. Efficiency reward (inverse of cost with context)
	efficiencyScore := erc.calculateEfficiencyScore(state, action)
	
	// 5. Sequence bonus for good action combinations
	sequenceBonus := erc.calculateSequenceBonus(state, action)
	
	// 6. Progress reward based on expected outcomes
	progressReward := erc.calculateProgressReward(state, action, example)
	
	// 7. Penalty for redundant actions
	redundancyPenalty := erc.calculateRedundancyPenalty(state, action)
	
	// 8. Time efficiency bonus
	timeBonus := erc.calculateTimeBonus(result, action)
	
	// 9. Diversity encouragement
	diversityBonus := erc.calculateDiversityBonus(state)
	
	// Combine all factors
	totalReward := baseReward + relevanceBonus + qualityScore + efficiencyScore + 
		sequenceBonus + progressReward - redundancyPenalty + timeBonus + diversityBonus
	
	// Apply difficulty scaling
	totalReward *= (1.0 + example.Difficulty * 0.5)
	
	// Normalize to reasonable range
	return math.Max(-5.0, math.Min(10.0, totalReward))
}

func (erc *EnhancedRewardCalculator) calculateRelevanceBonus(action Action, example TrainingExample) float64 {
	relevance := 0.0
	
	// Check if action is relevant to task type
	switch example.TaskType {
	case "code_analysis":
		if action.FunctionName == "detect_code" || action.FunctionName == "analyze_readability" {
			relevance = 1.0
		}
	case "technical_analysis", "academic_analysis", "scientific_analysis":
		if action.FunctionName == "extract_entities" || action.FunctionName == "extract_keywords" {
			relevance = 0.8
		}
	case "business_communication", "marketing_analysis":
		if action.FunctionName == "sentiment_analysis" || action.FunctionName == "extract_entities" {
			relevance = 0.7
		}
	case "social_media_analysis":
		if action.FunctionName == "sentiment_analysis" {
			relevance = 0.9
		}
	case "legal_analysis", "medical_analysis":
		if action.FunctionName == "extract_entities" {
			relevance = 1.0
		}
	case "log_analysis":
		if action.FunctionName == "detect_code" || action.FunctionName == "extract_keywords" {
			relevance = 0.8
		}
	}
	
	return relevance
}

func (erc *EnhancedRewardCalculator) assessOutputQuality(action Action, result ActionResult, example TrainingExample) float64 {
	outputMap, ok := result.Output.(map[string]interface{})
	if !ok {
		return 0.3
	}
	
	quality := 0.0
	
	// Check against expected outcomes
	if _, exists := example.Expected[action.FunctionName]; exists {
		// Compare output with expected (simplified)
		quality += 0.5
	}
	
	// Assess based on output completeness
	switch action.FunctionName {
	case "extract_entities":
		if entities, ok := outputMap["entities"].([]map[string]interface{}); ok {
			quality += float64(len(entities)) * 0.1
		}
	case "analyze_readability":
		if score, ok := outputMap["readability_score"].(float64); ok {
			// Reward scores in reasonable range
			if score >= 30 && score <= 90 {
				quality += 0.4
			}
		}
	case "sentiment_analysis":
		if _, hasScore := outputMap["score"]; hasScore {
			if _, hasConfidence := outputMap["confidence"]; hasConfidence {
				quality += 0.5
			}
		}
	}
	
	// Check quality threshold
	if threshold, exists := erc.QualityThresholds[strings.Replace(action.FunctionName, "_", "_", -1)]; exists {
		if quality >= threshold {
			quality *= 1.2
		}
	}
	
	return math.Min(1.5, quality)
}

func (erc *EnhancedRewardCalculator) calculateEfficiencyScore(state State, action Action) float64 {
	// Reward efficiency based on remaining budget and action cost
	budgetRatio := float64(state.RemainingBudget) / 50.0
	costEfficiency := 1.0 / float64(action.Cost)
	
	// Penalize expensive actions when budget is low
	if budgetRatio < 0.3 && action.Cost > 5 {
		return -0.5
	}
	
	return costEfficiency * budgetRatio
}

func (erc *EnhancedRewardCalculator) calculateSequenceBonus(state State, action Action) float64 {
	if len(state.ActionsUsed) == 0 {
		return 0
	}
	
	lastAction := state.ActionsUsed[len(state.ActionsUsed)-1]
	sequenceKey := lastAction + "->" + action.FunctionName
	
	if bonus, exists := erc.SequenceBonus[sequenceKey]; exists {
		return bonus
	}
	
	// Check for generally good sequences
	switch lastAction {
	case "extract_entities":
		if action.FunctionName == "extract_keywords" || action.FunctionName == "validate_output" {
			return 0.2
		}
	case "detect_code":
		if action.FunctionName == "analyze_readability" {
			return 0.15
		}
	case "analyze_readability":
		if action.FunctionName == "summarize_text" {
			return 0.25
		}
	}
	
	return 0
}

func (erc *EnhancedRewardCalculator) calculateProgressReward(state State, action Action, example TrainingExample) float64 {
	// Reward actions that make progress toward expected outcomes
	progress := 0.0
	
	// Check if this action helps achieve expected results
	for key := range example.Expected {
		if strings.Contains(key, strings.Split(action.FunctionName, "_")[0]) {
			progress += 0.3
		}
	}
	
	// Reward completing different types of analysis
	actionTypes := make(map[string]bool)
	for _, usedAction := range state.ActionsUsed {
		category := "analysis"
		if strings.Contains(usedAction, "format") {
			category = "formatting"
		} else if strings.Contains(usedAction, "validate") {
			category = "validation"
		}
		actionTypes[category] = true
	}
	
	if len(actionTypes) >= 2 {
		progress += 0.2
	}
	
	return progress
}

func (erc *EnhancedRewardCalculator) calculateRedundancyPenalty(state State, action Action) float64 {
	// Penalize repeated actions
	count := 0
	for _, usedAction := range state.ActionsUsed {
		if usedAction == action.FunctionName {
			count++
		}
	}
	
	if count > 0 {
		// Increasingly penalize repetition
		return float64(count) * 0.3
	}
	
	return 0
}

func (erc *EnhancedRewardCalculator) calculateTimeBonus(result ActionResult, action Action) float64 {
	// Reward fast execution
	executionMs := result.Duration.Milliseconds()
	expectedMs := int64(action.Cost * 10)
	
	if executionMs < expectedMs {
		return 0.1
	} else if executionMs > expectedMs*2 {
		return -0.1
	}
	
	return 0
}

func (erc *EnhancedRewardCalculator) calculateDiversityBonus(state State) float64 {
	// Reward trying different actions
	uniqueActions := make(map[string]bool)
	for _, action := range state.ActionsUsed {
		uniqueActions[action] = true
	}
	
	diversityRatio := float64(len(uniqueActions)) / float64(len(state.ActionsUsed)+1)
	return diversityRatio * 0.3
}