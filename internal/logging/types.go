package logging

import (
	"sync"
	"time"
)

type LogEvent struct {
	Timestamp       time.Time         `json:"timestamp"`
	SessionID       string            `json:"session_id"`
	EpisodeID       string            `json:"episode_id"`
	StepNumber      int               `json:"step_number"`
	EventType       string            `json:"event_type"`
	StateSnapshot   StateMetrics      `json:"state_snapshot,omitempty"`
	ActionTaken     ActionMetrics     `json:"action_taken,omitempty"`
	ResultMetrics   ResultMetrics     `json:"result_metrics,omitempty"`
	Performance     PerformanceMetrics `json:"performance,omitempty"`
	LearningMetrics LearningMetrics   `json:"learning_metrics,omitempty"`
}

type StateMetrics struct {
	TextLength      int                `json:"text_length"`
	TextComplexity  float64            `json:"text_complexity"`
	EntityDensity   float64            `json:"entity_density"`
	CodePresence    bool               `json:"code_presence"`
	MathPresence    bool               `json:"math_presence"`
	StateHash       string             `json:"state_hash"`
	Features        map[string]float64 `json:"features"`
}

type ActionMetrics struct {
	FunctionName    string  `json:"function_name"`
	Category        string  `json:"category"`
	ComputeCost     int     `json:"compute_cost"`
	InputSize       int     `json:"input_size"`
	ExpectedOutput  string  `json:"expected_output"`
	QValue          float64 `json:"q_value"`
	ExplorationFlag bool    `json:"exploration_flag"`
}

type ResultMetrics struct {
	Success         bool    `json:"success"`
	OutputQuality   float64 `json:"output_quality"`
	ExecutionTime   float64 `json:"execution_time"`
	MemoryUsed      int64   `json:"memory_used"`
	ErrorType       string  `json:"error_type,omitempty"`
	OutputSize      int     `json:"output_size"`
}

type PerformanceMetrics struct {
	CumulativeReward   float64 `json:"cumulative_reward"`
	AverageReward      float64 `json:"average_reward"`
	SuccessRate        float64 `json:"success_rate"`
	EfficiencyScore    float64 `json:"efficiency_score"`
	TaskCompletionRate float64 `json:"task_completion_rate"`
}

type LearningMetrics struct {
	QValueConvergence float64 `json:"q_value_convergence"`
	ExplorationRate   float64 `json:"exploration_rate"`
	PolicyStability   float64 `json:"policy_stability"`
	ActionDiversity   float64 `json:"action_diversity"`
	LearningProgress  float64 `json:"learning_progress"`
}

type EpisodeMetrics struct {
	EpisodeID   string          `json:"episode_id"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     time.Time       `json:"end_time"`
	Actions     []ActionMetrics `json:"actions"`
	Rewards     []float64       `json:"rewards"`
	States      []StateMetrics  `json:"states"`
	TotalReward float64         `json:"total_reward"`
}

type InsightLogger struct {
	LogPath       string
	MetricsDB     *MetricsDatabase
	EventStream   chan LogEvent
	BatchSize     int
	FlushInterval time.Duration
	
	mu        sync.RWMutex
	batch     []LogEvent
	sessionID string
	active    bool
	done      chan struct{}
}

type MetricsDatabase struct {
	mu      sync.RWMutex
	events  []LogEvent
	episodes map[string]EpisodeMetrics
}