package rl

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"time"
	"textlib-rl-system/internal/logging"
	"textlib-rl-system/internal/telemetry"
)

func NewEnhancedRLSystem(config SystemConfig) *EnhancedRLSystem {
	return &EnhancedRLSystem{
		Agent: NewQLearningAgent(0.1, 0.95, 1.0, 0.01, 0.995),
		RewardCalc: &RewardCalculator{
			TaskWeights: map[string]float64{
				"entity_extraction":    1.0,
				"readability_analysis": 0.8,
				"code_analysis":       0.9,
				"comprehensive":       1.2,
			},
		},
		Config:           config,
		availableActions: getDefaultActions(),
		simulator:        NewActionSimulator(),
	}
}

func (system *EnhancedRLSystem) SetLogger(logger *logging.InsightLogger) {
	system.Logger = logger
}

func (system *EnhancedRLSystem) SetTelemetry(telemetry *telemetry.TelemetryClient) {
	system.Telemetry = telemetry
}

func (system *EnhancedRLSystem) LoadTrainingData(data []TrainingExample) {
	system.TrainingData = data
}

func (system *EnhancedRLSystem) TrainWithLogging() {
	sessionID := generateSessionID()
	system.Logger.StartSession(sessionID)

	defer system.Logger.EndSession()
	defer system.SaveFinalModel()

	for episode := 0; episode < system.Config.MaxEpisodes; episode++ {
		episodeID := fmt.Sprintf("%s-ep%d", sessionID, episode)
		episodeMetrics := system.runEpisodeWithLogging(episodeID)

		system.Logger.LogEpisodeSummary(episodeMetrics)

		if episode%system.Config.CheckpointInterval == 0 {
			system.checkpointModel(episode)
		}

		if episode%system.Config.LoggingInterval == 0 {
			insights := system.analyzeProgress()
			system.Logger.LogInsights(insights)
		}
	}
}

func (system *EnhancedRLSystem) runEpisodeWithLogging(episodeID string) logging.EpisodeMetrics {
	example := system.selectTrainingExample()
	state := system.createInitialState(example)

	episodeMetrics := logging.EpisodeMetrics{
		EpisodeID: episodeID,
		StartTime: time.Now(),
		Actions:   []logging.ActionMetrics{},
		Rewards:   []float64{},
		States:    []logging.StateMetrics{},
	}

	for step := 0; step < system.Config.MaxStepsPerEpisode; step++ {
		stepStartTime := time.Now()

		stateMetrics := system.extractStateMetrics(state)
		system.Logger.LogEvent(logging.LogEvent{
			Timestamp:     stepStartTime,
			EpisodeID:     episodeID,
			StepNumber:    step,
			EventType:     "state_observation",
			StateSnapshot: stateMetrics,
		})

		action, actionMetrics := system.Agent.SelectActionWithMetrics(state)

		system.Logger.LogEvent(logging.LogEvent{
			Timestamp:   time.Now(),
			EpisodeID:   episodeID,
			StepNumber:  step,
			EventType:   "action_selected",
			ActionTaken: actionMetrics,
		})

		result := system.simulateAction(state, action, example)

		// Use enhanced reward calculation
		enhancedCalc := NewEnhancedRewardCalculator()
		reward := enhancedCalc.CalculateReward(state, action, result, example)

		system.Logger.LogEvent(logging.LogEvent{
			Timestamp:     time.Now(),
			EpisodeID:     episodeID,
			StepNumber:    step,
			EventType:     "reward_calculated",
			ResultMetrics: system.extractResultMetrics(result),
			Performance: logging.PerformanceMetrics{
				CumulativeReward: reward,
			},
		})

		nextState := system.updateState(state, action, result)

		oldQValue := system.Agent.GetQValue(state, action)
		system.Agent.UpdateQValue(state, action, reward, nextState)
		newQValue := system.Agent.GetQValue(state, action)

		system.Logger.LogEvent(logging.LogEvent{
			Timestamp:  time.Now(),
			EpisodeID:  episodeID,
			StepNumber: step,
			EventType:  "q_value_updated",
			LearningMetrics: logging.LearningMetrics{
				QValueConvergence: math.Abs(newQValue - oldQValue),
			},
		})

		episodeMetrics.Actions = append(episodeMetrics.Actions, actionMetrics)
		episodeMetrics.Rewards = append(episodeMetrics.Rewards, reward)
		episodeMetrics.States = append(episodeMetrics.States, stateMetrics)

		if system.isTaskComplete(nextState) {
			break
		}

		state = nextState
	}

	episodeMetrics.EndTime = time.Now()
	episodeMetrics.TotalReward = sum(episodeMetrics.Rewards)

	return episodeMetrics
}

func (system *EnhancedRLSystem) selectTrainingExample() TrainingExample {
	if len(system.TrainingData) == 0 {
		return TrainingExample{
			ID:       "default",
			Text:     "Sample text for analysis and processing.",
			TaskType: "comprehensive",
			Expected: map[string]interface{}{
				"entities": []string{"text", "analysis"},
				"quality":  0.8,
			},
			Difficulty: 0.5,
		}
	}
	// Randomly select from training data for variety
	idx := int(time.Now().UnixNano()) % len(system.TrainingData)
	return system.TrainingData[idx]
}

func (system *EnhancedRLSystem) createInitialState(example TrainingExample) State {
	return State{
		Text:            example.Text,
		TaskType:        example.TaskType,
		ActionsUsed:     []string{},
		CurrentResults:  make(map[string]interface{}),
		StepCount:       0,
		RemainingBudget: 50,
	}
}

func (system *EnhancedRLSystem) extractStateMetrics(state State) logging.StateMetrics {
	return logging.StateMetrics{
		TextLength:     len(state.Text),
		TextComplexity: calculateTextComplexity(state.Text),
		EntityDensity:  calculateEntityDensity(state.Text),
		CodePresence:   detectCodePresence(state.Text),
		MathPresence:   detectMathPresence(state.Text),
		StateHash:      system.getStateHash(state),
		Features: map[string]float64{
			"steps_taken":       float64(state.StepCount),
			"remaining_budget":  float64(state.RemainingBudget),
			"actions_used":      float64(len(state.ActionsUsed)),
		},
	}
}

func (system *EnhancedRLSystem) simulateAction(state State, action Action, example TrainingExample) ActionResult {
	return system.simulator.ExecuteAction(action, state.Text, action.Parameters)
}

func (system *EnhancedRLSystem) extractResultMetrics(result ActionResult) logging.ResultMetrics {
	return logging.ResultMetrics{
		Success:       result.Success,
		OutputQuality: calculateOutputQuality(result.Output),
		ExecutionTime: result.Duration.Seconds(),
		MemoryUsed:    result.MemoryUsed,
		ErrorType:     result.Error,
		OutputSize:    calculateOutputSize(result.Output),
	}
}

func (system *EnhancedRLSystem) updateState(state State, action Action, result ActionResult) State {
	newState := state
	newState.StepCount++
	newState.RemainingBudget -= action.Cost
	newState.ActionsUsed = append(newState.ActionsUsed, action.FunctionName)
	
	if result.Success {
		newState.CurrentResults[action.FunctionName] = result.Output
	}
	
	return newState
}

func (system *EnhancedRLSystem) isTaskComplete(state State) bool {
	return state.RemainingBudget <= 0 || len(state.ActionsUsed) >= 10
}

func (system *EnhancedRLSystem) checkpointModel(episode int) {
	log.Printf("Checkpointing model at episode %d", episode)
}

func (system *EnhancedRLSystem) analyzeProgress() interface{} {
	return map[string]interface{}{
		"timestamp": time.Now(),
		"progress":  "ongoing",
	}
}

func (system *EnhancedRLSystem) SaveFinalModel() {
	log.Println("Saving final model")
}

func (system *EnhancedRLSystem) SaveInsights(insights interface{}) {
	log.Println("Saving insights")
}

func (system *EnhancedRLSystem) getStateHash(state State) string {
	data := fmt.Sprintf("%s_%s_%d", state.Text, state.TaskType, state.StepCount)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().Unix())
}

func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}

func calculateTextComplexity(text string) float64 {
	return math.Min(1.0, float64(len(text))/1000.0)
}

func calculateEntityDensity(text string) float64 {
	return 0.1 // Simplified
}

func detectCodePresence(text string) bool {
	return false // Simplified
}

func detectMathPresence(text string) bool {
	return false // Simplified
}

func calculateOutputQuality(output interface{}) float64 {
	return 0.8 // Simplified
}

func calculateOutputSize(output interface{}) int {
	return 100 // Simplified
}

func getDefaultActions() []Action {
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