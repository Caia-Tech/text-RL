package rl

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"textlib-rl-system/internal/logging"
)

func NewQLearningAgent(learningRate, discountFactor, explorationRate, minExploration, decayRate float64) *QLearningAgent {
	return &QLearningAgent{
		QTable:          make(map[string]map[string]float64),
		LearningRate:    learningRate,
		DiscountFactor:  discountFactor,
		ExplorationRate: explorationRate,
		MinExploration:  minExploration,
		DecayRate:       decayRate,
	}
}

func (agent *QLearningAgent) SelectActionWithMetrics(state State) (Action, logging.ActionMetrics) {
	
	var selectedAction Action
	var qValue float64
	var isExploration bool
	
	if rand.Float64() < agent.ExplorationRate {
		selectedAction = agent.selectRandomAction(state)
		qValue = agent.GetQValue(state, selectedAction)
		isExploration = true
	} else {
		selectedAction = agent.selectBestAction(state)
		qValue = agent.GetQValue(state, selectedAction)
		isExploration = false
	}
	
	metrics := logging.ActionMetrics{
		FunctionName:    selectedAction.FunctionName,
		Category:        selectedAction.Category,
		ComputeCost:     selectedAction.Cost,
		InputSize:       len(state.Text),
		ExpectedOutput:  "simulated_output",
		QValue:          qValue,
		ExplorationFlag: isExploration,
	}
	
	return selectedAction, metrics
}

func (agent *QLearningAgent) GetQValue(state State, action Action) float64 {
	stateKey := agent.getStateKey(state)
	actionKey := agent.getActionKey(action)
	
	if stateActions, exists := agent.QTable[stateKey]; exists {
		if qValue, exists := stateActions[actionKey]; exists {
			return qValue
		}
	}
	
	return 0.0
}

func (agent *QLearningAgent) UpdateQValue(state State, action Action, reward float64, nextState State) {
	stateKey := agent.getStateKey(state)
	actionKey := agent.getActionKey(action)
	
	if agent.QTable[stateKey] == nil {
		agent.QTable[stateKey] = make(map[string]float64)
	}
	
	currentQ := agent.GetQValue(state, action)
	maxNextQ := agent.getMaxQValue(nextState)
	
	newQ := currentQ + agent.LearningRate*(reward+agent.DiscountFactor*maxNextQ-currentQ)
	agent.QTable[stateKey][actionKey] = newQ
	
	agent.ExplorationRate = math.Max(agent.MinExploration, agent.ExplorationRate*agent.DecayRate)
}

func (agent *QLearningAgent) selectBestAction(state State) Action {
	availableActions := agent.getAvailableActions(state)
	
	if len(availableActions) == 0 {
		return Action{FunctionName: "no_op", Category: "utility", Cost: 0}
	}
	
	bestAction := availableActions[0]
	bestQValue := agent.GetQValue(state, bestAction)
	
	for _, action := range availableActions[1:] {
		qValue := agent.GetQValue(state, action)
		if qValue > bestQValue {
			bestQValue = qValue
			bestAction = action
		}
	}
	
	return bestAction
}

func (agent *QLearningAgent) selectRandomAction(state State) Action {
	availableActions := agent.getAvailableActions(state)
	if len(availableActions) == 0 {
		return Action{FunctionName: "no_op", Category: "utility", Cost: 0}
	}
	
	return availableActions[rand.Intn(len(availableActions))]
}

func (agent *QLearningAgent) getMaxQValue(state State) float64 {
	availableActions := agent.getAvailableActions(state)
	if len(availableActions) == 0 {
		return 0.0
	}
	
	maxQ := agent.GetQValue(state, availableActions[0])
	for _, action := range availableActions[1:] {
		qValue := agent.GetQValue(state, action)
		if qValue > maxQ {
			maxQ = qValue
		}
	}
	
	return maxQ
}

func (agent *QLearningAgent) getAvailableActions(state State) []Action {
	return []Action{
		{FunctionName: "extract_entities", Category: "analysis", Cost: 5},
		{FunctionName: "analyze_readability", Category: "analysis", Cost: 3},
		{FunctionName: "detect_code", Category: "analysis", Cost: 2},
		{FunctionName: "extract_keywords", Category: "analysis", Cost: 4},
		{FunctionName: "sentiment_analysis", Category: "analysis", Cost: 3},
		{FunctionName: "summarize_text", Category: "generation", Cost: 8},
		{FunctionName: "format_text", Category: "formatting", Cost: 2},
		{FunctionName: "validate_output", Category: "utility", Cost: 1},
	}
}

func (agent *QLearningAgent) getStateKey(state State) string {
	data := fmt.Sprintf("%s|%s|%d|%d", state.TaskType, state.Text[:min(100, len(state.Text))], 
		state.StepCount, state.RemainingBudget)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

func (agent *QLearningAgent) getActionKey(action Action) string {
	return fmt.Sprintf("%s_%s", action.FunctionName, action.Category)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}