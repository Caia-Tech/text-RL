package analyzer

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"
	"textlib-rl-system/internal/logging"
)

func NewInsightAnalyzer(logger *logging.InsightLogger, metricsDB *logging.MetricsDatabase, analysisWindow int) *InsightAnalyzer {
	return &InsightAnalyzer{
		Logger:         logger,
		MetricsDB:      metricsDB,
		AnalysisWindow: analysisWindow,
	}
}

func (analyzer *InsightAnalyzer) GenerateInsights() APIFeedbackReport {
	log.Println("Starting comprehensive insight analysis...")
	
	report := APIFeedbackReport{
		Timestamp:          time.Now(),
		AnalysisPeriod:     analyzer.AnalysisWindow,
		FunctionUsageStats: analyzer.analyzeFunctionUsage(),
		SequencePatterns:   analyzer.analyzeActionSequences(),
		PerformanceMetrics: analyzer.analyzePerformance(),
		LearningCurve:      analyzer.analyzeLearningProgress(),
		OptimalSequences:   analyzer.findOptimalSequences(),
		FailureAnalysis:    analyzer.analyzeFailures(),
		UsageInsights:      analyzer.analyzeUsagePatterns(),
		Recommendations:    []string{},
	}
	
	// Generate recommendations based on analysis
	report.Recommendations = analyzer.generateRecommendations(report)
	
	log.Printf("Generated %d recommendations from analysis", len(report.Recommendations))
	return report
}

func (analyzer *InsightAnalyzer) analyzeFunctionUsage() map[string]FunctionStats {
	events := analyzer.MetricsDB.GetEvents()
	functionStats := make(map[string]FunctionStats)
	
	// Group events by function
	functionEvents := make(map[string][]logging.LogEvent)
	for _, event := range events {
		if event.EventType == "action_selected" {
			functionName := event.ActionTaken.FunctionName
			functionEvents[functionName] = append(functionEvents[functionName], event)
		}
	}
	
	// Analyze each function
	for functionName, funcEvents := range functionEvents {
		stats := analyzer.calculateFunctionStats(funcEvents)
		functionStats[functionName] = stats
	}
	
	return functionStats
}

func (analyzer *InsightAnalyzer) calculateFunctionStats(events []logging.LogEvent) FunctionStats {
	if len(events) == 0 {
		return FunctionStats{}
	}
	
	var totalReward, totalQValue, totalDuration float64
	var successCount int
	errorTypes := make(map[string]int)
	contexts := make(map[string]int)
	qualityPoints := []QualityDataPoint{}
	
	for i, event := range events {
		// Find corresponding result event
		resultEvent := analyzer.findResultEvent(event.EpisodeID, event.StepNumber)
		if resultEvent != nil {
			if resultEvent.ResultMetrics.Success {
				successCount++
			} else if resultEvent.ResultMetrics.ErrorType != "" {
				errorTypes[resultEvent.ResultMetrics.ErrorType]++
			}
			
			totalReward += resultEvent.Performance.CumulativeReward
			totalDuration += resultEvent.ResultMetrics.ExecutionTime
			
			qualityPoints = append(qualityPoints, QualityDataPoint{
				Episode: i,
				Quality: resultEvent.ResultMetrics.OutputQuality,
			})
		}
		
		totalQValue += event.ActionTaken.QValue
		
		// Extract context (task type)
		if event.StateSnapshot.Features != nil {
			for context := range event.StateSnapshot.Features {
				contexts[context]++
			}
		}
	}
	
	callCount := len(events)
	successRate := float64(successCount) / float64(callCount)
	avgReward := totalReward / float64(callCount)
	avgQValue := totalQValue / float64(callCount)
	avgDuration := totalDuration / float64(callCount)
	
	// Get common contexts
	commonContexts := analyzer.getTopContexts(contexts, 5)
	
	// Calculate quality metrics
	qualityMetrics := analyzer.calculateQualityMetrics(qualityPoints)
	
	return FunctionStats{
		CallCount:      callCount,
		SuccessRate:    successRate,
		AvgReward:      avgReward,
		AvgQValue:      avgQValue,
		AvgDuration:    avgDuration,
		CommonContexts: commonContexts,
		ErrorTypes:     errorTypes,
		QualityMetrics: qualityMetrics,
	}
}

func (analyzer *InsightAnalyzer) analyzeActionSequences() []SequencePattern {
	episodes := analyzer.MetricsDB.GetEpisodes()
	sequenceCounts := make(map[string]SequenceInfo)
	
	for _, episode := range episodes {
		// Extract action sequence
		var sequence []string
		for _, action := range episode.Actions {
			sequence = append(sequence, action.FunctionName)
		}
		
		// Analyze subsequences of length 2-5
		for length := 2; length <= min(5, len(sequence)); length++ {
			for i := 0; i <= len(sequence)-length; i++ {
				subseq := sequence[i : i+length]
				key := strings.Join(subseq, "->")
				
				info := sequenceCounts[key]
				info.Sequence = subseq
				info.Frequency++
				info.TotalReward += episode.TotalReward
				if episode.TotalReward > 0 {
					info.SuccessCount++
				}
				sequenceCounts[key] = info
			}
		}
	}
	
	// Convert to sorted list
	var patterns []SequencePattern
	for _, info := range sequenceCounts {
		if info.Frequency >= 3 { // Only include patterns that appear at least 3 times
			pattern := SequencePattern{
				Sequence:    info.Sequence,
				Frequency:   info.Frequency,
				AvgReward:   info.TotalReward / float64(info.Frequency),
				SuccessRate: float64(info.SuccessCount) / float64(info.Frequency),
				Efficiency:  analyzer.calculateSequenceEfficiency(info),
				TaskTypes:   []string{"comprehensive"}, // Simplified
			}
			patterns = append(patterns, pattern)
		}
	}
	
	// Sort by frequency and reward
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency*int(patterns[i].AvgReward*10) > 
			patterns[j].Frequency*int(patterns[j].AvgReward*10)
	})
	
	return patterns
}

func (analyzer *InsightAnalyzer) analyzePerformance() OverallPerformance {
	episodes := analyzer.MetricsDB.GetEpisodes()
	events := analyzer.MetricsDB.GetEvents()
	
	if len(episodes) == 0 {
		return OverallPerformance{}
	}
	
	var totalReward float64
	var totalSteps int
	var successfulEpisodes int
	var qValueVariances []float64
	
	for _, episode := range episodes {
		totalReward += episode.TotalReward
		totalSteps += len(episode.Actions)
		
		if episode.TotalReward > 0 {
			successfulEpisodes++
		}
		
		// Calculate Q-value variance for this episode
		variance := analyzer.calculateQValueVariance(episode)
		qValueVariances = append(qValueVariances, variance)
	}
	
	totalEpisodes := len(episodes)
	avgEpisodeReward := totalReward / float64(totalEpisodes)
	overallSuccessRate := float64(successfulEpisodes) / float64(totalEpisodes)
	
	// Calculate learning efficiency (improvement rate)
	episodesList := make([]logging.EpisodeMetrics, 0, len(episodes))
	for _, ep := range episodes {
		episodesList = append(episodesList, ep)
	}
	learningEfficiency := analyzer.calculateLearningEfficiency(episodesList)
	
	// Calculate convergence rate
	convergenceRate := analyzer.calculateConvergenceRate(qValueVariances)
	
	// Calculate exploration balance
	explorationBalance := analyzer.calculateExplorationBalance(events)
	
	return OverallPerformance{
		TotalEpisodes:      totalEpisodes,
		TotalSteps:         totalSteps,
		AvgEpisodeReward:   avgEpisodeReward,
		OverallSuccessRate: overallSuccessRate,
		LearningEfficiency: learningEfficiency,
		ConvergenceRate:    convergenceRate,
		ExplorationBalance: explorationBalance,
	}
}

func (analyzer *InsightAnalyzer) analyzeLearningProgress() []LearningDataPoint {
	episodes := analyzer.MetricsDB.GetEpisodes()
	
	// Sort episodes by episode ID (assuming chronological order)
	episodeList := make([]logging.EpisodeMetrics, 0, len(episodes))
	for _, episode := range episodes {
		episodeList = append(episodeList, episode)
	}
	
	sort.Slice(episodeList, func(i, j int) bool {
		return episodeList[i].EpisodeID < episodeList[j].EpisodeID
	})
	
	var learningCurve []LearningDataPoint
	windowSize := 10 // Moving average window
	
	for i := windowSize; i < len(episodeList); i++ {
		window := episodeList[i-windowSize : i]
		
		var avgReward, avgSuccess float64
		successCount := 0
		
		for _, ep := range window {
			avgReward += ep.TotalReward
			if ep.TotalReward > 0 {
				successCount++
			}
		}
		
		avgReward /= float64(windowSize)
		avgSuccess = float64(successCount) / float64(windowSize)
		
		dataPoint := LearningDataPoint{
			Episode:         i,
			AvgReward:       avgReward,
			SuccessRate:     avgSuccess,
			QValueVariance:  analyzer.calculateQValueVariance(episodeList[i]),
			ExplorationRate: 0.5, // Simplified
			PolicyStability: analyzer.calculatePolicyStability(window),
		}
		
		learningCurve = append(learningCurve, dataPoint)
	}
	
	return learningCurve
}

func (analyzer *InsightAnalyzer) findOptimalSequences() map[string][]string {
	patterns := analyzer.analyzeActionSequences()
	optimal := make(map[string][]string)
	
	// Find best sequences for different criteria
	var bestReward, bestSuccess, bestEfficiency SequencePattern
	maxReward, maxSuccess, maxEfficiency := -math.Inf(1), 0.0, 0.0
	
	for _, pattern := range patterns {
		if pattern.AvgReward > maxReward {
			maxReward = pattern.AvgReward
			bestReward = pattern
		}
		if pattern.SuccessRate > maxSuccess {
			maxSuccess = pattern.SuccessRate
			bestSuccess = pattern
		}
		if pattern.Efficiency > maxEfficiency {
			maxEfficiency = pattern.Efficiency
			bestEfficiency = pattern
		}
	}
	
	optimal["highest_reward"] = bestReward.Sequence
	optimal["highest_success"] = bestSuccess.Sequence
	optimal["most_efficient"] = bestEfficiency.Sequence
	
	return optimal
}

func (analyzer *InsightAnalyzer) analyzeFailures() FailureAnalysis {
	events := analyzer.MetricsDB.GetEvents()
	
	commonFailures := make(map[string]FailurePattern)
	failuresByFunction := make(map[string][]string)
	
	for _, event := range events {
		if event.EventType == "reward_calculated" && !event.ResultMetrics.Success {
			errorType := event.ResultMetrics.ErrorType
			if errorType == "" {
				errorType = "unknown_failure"
			}
			
			// Track common failures
			pattern := commonFailures[errorType]
			pattern.Pattern = errorType
			pattern.Frequency++
			pattern.Impact += math.Abs(event.Performance.CumulativeReward)
			commonFailures[errorType] = pattern
			
			// Track failures by function
			if event.ActionTaken.FunctionName != "" {
				failuresByFunction[event.ActionTaken.FunctionName] = 
					append(failuresByFunction[event.ActionTaken.FunctionName], errorType)
			}
		}
	}
	
	// Generate suggestions for common failures
	for errorType, pattern := range commonFailures {
		pattern.Suggestions = analyzer.generateFailureSuggestions(errorType, pattern)
		commonFailures[errorType] = pattern
	}
	
	return FailureAnalysis{
		CommonFailures:     commonFailures,
		FailuresByFunction: failuresByFunction,
		RecoveryPatterns:   analyzer.findRecoveryPatterns(),
		CriticalIssues:     analyzer.identifyCriticalIssues(commonFailures),
	}
}

func (analyzer *InsightAnalyzer) analyzeUsagePatterns() UsageInsights {
	// Placeholder implementation
	return UsageInsights{
		PreferredSequences: map[string][]string{
			"analysis_tasks": {"extract_entities", "analyze_readability"},
			"generation_tasks": {"extract_keywords", "summarize_text"},
		},
		TaskSpecificTrends: make(map[string]TaskTrend),
		EfficiencyTrends: []EfficiencyDataPoint{},
		OptimizationAreas: []OptimizationArea{},
	}
}

func (analyzer *InsightAnalyzer) generateRecommendations(report APIFeedbackReport) []string {
	var recommendations []string
	
	// Analyze function effectiveness
	for functionName, stats := range report.FunctionUsageStats {
		if stats.SuccessRate < 0.7 {
			recommendations = append(recommendations,
				fmt.Sprintf("Function '%s' has low success rate (%.2f%%). Consider improving documentation or error handling.",
					functionName, stats.SuccessRate*100))
		}
		
		if stats.AvgReward < 0 {
			recommendations = append(recommendations,
				fmt.Sprintf("Function '%s' shows negative average reward (%.2f). Review implementation or usage patterns.",
					functionName, stats.AvgReward))
		}
	}
	
	// Analyze common failure patterns
	for pattern, failure := range report.FailureAnalysis.CommonFailures {
		if failure.Frequency > 10 {
			recommendations = append(recommendations,
				fmt.Sprintf("Common failure pattern detected: '%s' (occurred %d times). Impact: %.2f",
					pattern, failure.Frequency, failure.Impact))
		}
	}
	
	// Performance recommendations
	if report.PerformanceMetrics.OverallSuccessRate < 0.8 {
		recommendations = append(recommendations,
			fmt.Sprintf("Overall success rate is low (%.2f%%). Consider reviewing task complexity or function implementations.",
				report.PerformanceMetrics.OverallSuccessRate*100))
	}
	
	return recommendations
}

// Helper functions

func (analyzer *InsightAnalyzer) findResultEvent(episodeID string, stepNumber int) *logging.LogEvent {
	events := analyzer.MetricsDB.GetEventsByEpisode(episodeID)
	for _, event := range events {
		if event.EventType == "reward_calculated" && event.StepNumber == stepNumber {
			return &event
		}
	}
	return nil
}

func (analyzer *InsightAnalyzer) getTopContexts(contexts map[string]int, limit int) []string {
	type contextPair struct {
		context string
		count   int
	}
	
	var pairs []contextPair
	for context, count := range contexts {
		pairs = append(pairs, contextPair{context, count})
	}
	
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	
	var result []string
	for i := 0; i < min(limit, len(pairs)); i++ {
		result = append(result, pairs[i].context)
	}
	
	return result
}

func (analyzer *InsightAnalyzer) calculateQualityMetrics(points []QualityDataPoint) QualityMetrics {
	if len(points) == 0 {
		return QualityMetrics{}
	}
	
	var totalQuality float64
	for _, point := range points {
		totalQuality += point.Quality
	}
	
	avgQuality := totalQuality / float64(len(points))
	
	return QualityMetrics{
		AvgOutputQuality: avgQuality,
		ConsistencyScore: 0.8, // Simplified
		ReliabilityScore: 0.9, // Simplified
		QualityTrend:     points,
	}
}

func (analyzer *InsightAnalyzer) calculateSequenceEfficiency(info SequenceInfo) float64 {
	if info.Frequency == 0 {
		return 0
	}
	
	avgReward := info.TotalReward / float64(info.Frequency)
	successRate := float64(info.SuccessCount) / float64(info.Frequency)
	
	return avgReward * successRate
}

func (analyzer *InsightAnalyzer) calculateQValueVariance(episode logging.EpisodeMetrics) float64 {
	if len(episode.Actions) == 0 {
		return 0
	}
	
	var sum, sumSquares float64
	for _, action := range episode.Actions {
		qValue := action.QValue
		sum += qValue
		sumSquares += qValue * qValue
	}
	
	n := float64(len(episode.Actions))
	mean := sum / n
	variance := (sumSquares / n) - (mean * mean)
	
	return variance
}

func (analyzer *InsightAnalyzer) calculateLearningEfficiency(episodes []logging.EpisodeMetrics) float64 {
	if len(episodes) < 2 {
		return 0
	}
	
	// Calculate improvement rate over time
	firstHalf := episodes[:len(episodes)/2]
	secondHalf := episodes[len(episodes)/2:]
	
	var firstAvg, secondAvg float64
	for _, ep := range firstHalf {
		firstAvg += ep.TotalReward
	}
	firstAvg /= float64(len(firstHalf))
	
	for _, ep := range secondHalf {
		secondAvg += ep.TotalReward
	}
	secondAvg /= float64(len(secondHalf))
	
	return (secondAvg - firstAvg) / math.Abs(firstAvg)
}

func (analyzer *InsightAnalyzer) calculateConvergenceRate(variances []float64) float64 {
	if len(variances) < 2 {
		return 0
	}
	
	// Calculate how quickly variance decreases
	firstVariance := variances[0]
	lastVariance := variances[len(variances)-1]
	
	return (firstVariance - lastVariance) / firstVariance
}

func (analyzer *InsightAnalyzer) calculateExplorationBalance(events []logging.LogEvent) float64 {
	explorationCount := 0
	totalActions := 0
	
	for _, event := range events {
		if event.EventType == "action_selected" {
			totalActions++
			if event.ActionTaken.ExplorationFlag {
				explorationCount++
			}
		}
	}
	
	if totalActions == 0 {
		return 0
	}
	
	return float64(explorationCount) / float64(totalActions)
}

func (analyzer *InsightAnalyzer) calculatePolicyStability(episodes []logging.EpisodeMetrics) float64 {
	// Simplified calculation
	return 0.8
}

func (analyzer *InsightAnalyzer) generateFailureSuggestions(errorType string, pattern FailurePattern) []string {
	suggestions := []string{}
	
	switch errorType {
	case "timeout":
		suggestions = append(suggestions, "Consider increasing timeout values or optimizing function performance")
	case "invalid_input":
		suggestions = append(suggestions, "Add input validation and preprocessing steps")
	case "memory_limit":
		suggestions = append(suggestions, "Optimize memory usage or increase available memory")
	default:
		suggestions = append(suggestions, "Review function implementation and error handling")
	}
	
	return suggestions
}

func (analyzer *InsightAnalyzer) findRecoveryPatterns() []RecoveryPattern {
	// Placeholder implementation
	return []RecoveryPattern{
		{
			FailureType:    "timeout",
			RecoveryAction: "retry_with_shorter_input",
			SuccessRate:    0.7,
			RecommendedFor: []string{"text_processing", "analysis"},
		},
	}
}

func (analyzer *InsightAnalyzer) identifyCriticalIssues(failures map[string]FailurePattern) []string {
	var critical []string
	
	for pattern, failure := range failures {
		if failure.Frequency > 20 && failure.Impact > 10.0 {
			critical = append(critical, fmt.Sprintf("Critical issue: %s (frequency: %d, impact: %.2f)",
				pattern, failure.Frequency, failure.Impact))
		}
	}
	
	return critical
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type SequenceInfo struct {
	Sequence     []string
	Frequency    int
	TotalReward  float64
	SuccessCount int
}