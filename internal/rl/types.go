package rl

import (
	"time"
	"textlib-rl-system/internal/logging"
	"textlib-rl-system/internal/telemetry"
)

type QLearningAgent struct {
	QTable          map[string]map[string]float64
	LearningRate    float64
	DiscountFactor  float64
	ExplorationRate float64
	MinExploration  float64
	DecayRate       float64
}

type RewardCalculator struct {
	TaskWeights map[string]float64
}

type TrainingExample struct {
	ID          string                 `json:"id"`
	Text        string                 `json:"text"`
	TaskType    string                 `json:"task_type"`
	Expected    map[string]interface{} `json:"expected"`
	Difficulty  float64               `json:"difficulty"`
}

type State struct {
	Text            string                 `json:"text"`
	TaskType        string                 `json:"task_type"`
	ActionsUsed     []string              `json:"actions_used"`
	CurrentResults  map[string]interface{} `json:"current_results"`
	StepCount       int                   `json:"step_count"`
	RemainingBudget int                   `json:"remaining_budget"`
}

type Action struct {
	FunctionName string                 `json:"function_name"`
	Parameters   map[string]interface{} `json:"parameters"`
	Category     string                 `json:"category"`
	Cost         int                   `json:"cost"`
}

type ActionResult struct {
	Success    bool                   `json:"success"`
	Output     interface{}            `json:"output"`
	Error      string                 `json:"error,omitempty"`
	Duration   time.Duration          `json:"duration"`
	MemoryUsed int64                  `json:"memory_used"`
}

type SystemConfig struct {
	MaxEpisodes        int
	MaxStepsPerEpisode int
	LoggingInterval    int
	CheckpointInterval int
	MetricsPort        int
	EnableProfiling    bool
}

type EnhancedRLSystem struct {
	Agent        *QLearningAgent
	Logger       *logging.InsightLogger
	Telemetry    *telemetry.TelemetryClient
	RewardCalc   *RewardCalculator
	TrainingData []TrainingExample
	Config       SystemConfig
	
	availableActions []Action
	simulator        *ActionSimulator
}

type ActionSimulator struct {
	Functions map[string]SimulatedFunction
}

type SimulatedFunction struct {
	Name           string
	Category       string
	Cost           int
	BaseSuccessRate float64
	OutputGenerator func(input string, params map[string]interface{}) (interface{}, error)
}