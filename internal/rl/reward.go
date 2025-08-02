package rl

import (
	"math"
)

func (rc *RewardCalculator) CalculateReward(state State, action Action, result ActionResult) float64 {
	if !result.Success {
		return -1.0 // Penalty for failed actions
	}
	
	baseReward := 1.0
	
	// Task-specific weight
	if weight, exists := rc.TaskWeights[state.TaskType]; exists {
		baseReward *= weight
	}
	
	// Efficiency bonus (inverse of cost)
	efficiencyBonus := 1.0 / float64(action.Cost)
	
	// Quality bonus based on output
	qualityBonus := rc.calculateQualityBonus(result.Output)
	
	// Time penalty (encourage faster execution)
	timePenalty := math.Min(0.5, result.Duration.Seconds()/10.0)
	
	// Step penalty (encourage fewer steps)
	stepPenalty := float64(state.StepCount) * 0.1
	
	// Diversity bonus (encourage trying different actions)
	diversityBonus := rc.calculateDiversityBonus(state.ActionsUsed, action.FunctionName)
	
	totalReward := baseReward + efficiencyBonus + qualityBonus + diversityBonus - timePenalty - stepPenalty
	
	// Normalize reward to reasonable range
	return math.Max(-5.0, math.Min(5.0, totalReward))
}

func (rc *RewardCalculator) calculateQualityBonus(output interface{}) float64 {
	if output == nil {
		return 0.0
	}
	
	// This is a simplified quality assessment
	// In a real implementation, this would be much more sophisticated
	outputMap, ok := output.(map[string]interface{})
	if !ok {
		return 0.5
	}
	
	qualityScore := 0.0
	
	// Check for specific quality indicators
	if score, exists := outputMap["score"]; exists {
		if scoreFloat, ok := score.(float64); ok {
			qualityScore += scoreFloat
		}
	}
	
	if confidence, exists := outputMap["confidence"]; exists {
		if confFloat, ok := confidence.(float64); ok {
			qualityScore += confFloat
		}
	}
	
	// Check for completeness (presence of expected fields)
	expectedFields := []string{"count", "entities", "keywords", "summary", "sentiment"}
	fieldCount := 0
	for _, field := range expectedFields {
		if _, exists := outputMap[field]; exists {
			fieldCount++
		}
	}
	
	qualityScore += float64(fieldCount) * 0.2
	
	return math.Min(2.0, qualityScore)
}

func (rc *RewardCalculator) calculateDiversityBonus(actionsUsed []string, currentAction string) float64 {
	// Encourage trying new actions
	for _, action := range actionsUsed {
		if action == currentAction {
			return 0.0 // No bonus for repeated actions
		}
	}
	
	// Bonus for trying new action
	return 0.5
}