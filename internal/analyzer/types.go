package analyzer

import (
	"time"
	"textlib-rl-system/internal/logging"
)

type InsightAnalyzer struct {
	Logger         *logging.InsightLogger
	MetricsDB      *logging.MetricsDatabase
	AnalysisWindow int
}

type APIFeedbackReport struct {
	Timestamp          time.Time                      `json:"timestamp"`
	AnalysisPeriod     int                           `json:"analysis_period"`
	FunctionUsageStats map[string]FunctionStats      `json:"function_usage_stats"`
	SequencePatterns   []SequencePattern             `json:"sequence_patterns"`
	PerformanceMetrics OverallPerformance            `json:"performance_metrics"`
	LearningCurve      []LearningDataPoint           `json:"learning_curve"`
	OptimalSequences   map[string][]string           `json:"optimal_sequences"`
	Recommendations    []string                      `json:"recommendations"`
	FailureAnalysis    FailureAnalysis               `json:"failure_analysis"`
	UsageInsights      UsageInsights                 `json:"usage_insights"`
}

type FunctionStats struct {
	CallCount       int                 `json:"call_count"`
	SuccessRate     float64            `json:"success_rate"`
	AvgReward       float64            `json:"avg_reward"`
	AvgQValue       float64            `json:"avg_q_value"`
	AvgDuration     float64            `json:"avg_duration"`
	CommonContexts  []string           `json:"common_contexts"`
	ErrorTypes      map[string]int     `json:"error_types"`
	QualityMetrics  QualityMetrics     `json:"quality_metrics"`
}

type SequencePattern struct {
	Sequence     []string `json:"sequence"`
	Frequency    int      `json:"frequency"`
	AvgReward    float64  `json:"avg_reward"`
	SuccessRate  float64  `json:"success_rate"`
	Efficiency   float64  `json:"efficiency"`
	TaskTypes    []string `json:"task_types"`
}

type OverallPerformance struct {
	TotalEpisodes      int     `json:"total_episodes"`
	TotalSteps         int     `json:"total_steps"`
	AvgEpisodeReward   float64 `json:"avg_episode_reward"`
	OverallSuccessRate float64 `json:"overall_success_rate"`
	LearningEfficiency float64 `json:"learning_efficiency"`
	ConvergenceRate    float64 `json:"convergence_rate"`
	ExplorationBalance float64 `json:"exploration_balance"`
}

type LearningDataPoint struct {
	Episode         int     `json:"episode"`
	AvgReward       float64 `json:"avg_reward"`
	SuccessRate     float64 `json:"success_rate"`
	QValueVariance  float64 `json:"q_value_variance"`
	ExplorationRate float64 `json:"exploration_rate"`
	PolicyStability float64 `json:"policy_stability"`
}

type FailureAnalysis struct {
	CommonFailures     map[string]FailurePattern `json:"common_failures"`
	FailuresByFunction map[string][]string        `json:"failures_by_function"`
	RecoveryPatterns   []RecoveryPattern          `json:"recovery_patterns"`
	CriticalIssues     []string                   `json:"critical_issues"`
}

type FailurePattern struct {
	Pattern     string  `json:"pattern"`
	Frequency   int     `json:"frequency"`
	Impact      float64 `json:"impact"`
	Suggestions []string `json:"suggestions"`
}

type RecoveryPattern struct {
	FailureType    string   `json:"failure_type"`
	RecoveryAction string   `json:"recovery_action"`
	SuccessRate    float64  `json:"success_rate"`
	RecommendedFor []string `json:"recommended_for"`
}

type UsageInsights struct {
	PreferredSequences  map[string][]string        `json:"preferred_sequences"`
	TaskSpecificTrends  map[string]TaskTrend       `json:"task_specific_trends"`
	EfficiencyTrends    []EfficiencyDataPoint      `json:"efficiency_trends"`
	OptimizationAreas   []OptimizationArea         `json:"optimization_areas"`
}

type TaskTrend struct {
	TaskType            string             `json:"task_type"`
	OptimalActions      []string           `json:"optimal_actions"`
	AvoidableActions    []string           `json:"avoidable_actions"`
	SuccessFactors      []string           `json:"success_factors"`
	PerformanceMetrics  map[string]float64 `json:"performance_metrics"`
}

type EfficiencyDataPoint struct {
	Episode            int     `json:"episode"`
	ActionsPerEpisode  float64 `json:"actions_per_episode"`
	RewardPerAction    float64 `json:"reward_per_action"`
	SuccessfulActions  float64 `json:"successful_actions"`
	WastedEffort       float64 `json:"wasted_effort"`
}

type OptimizationArea struct {
	Area        string   `json:"area"`
	Impact      float64  `json:"impact"`
	Difficulty  float64  `json:"difficulty"`
	Priority    string   `json:"priority"`
	Actions     []string `json:"actions"`
	ExpectedGain float64 `json:"expected_gain"`
}

type QualityMetrics struct {
	AvgOutputQuality float64            `json:"avg_output_quality"`
	ConsistencyScore float64            `json:"consistency_score"`
	ReliabilityScore float64            `json:"reliability_score"`
	QualityTrend     []QualityDataPoint `json:"quality_trend"`
}

type QualityDataPoint struct {
	Episode int     `json:"episode"`
	Quality float64 `json:"quality"`
}