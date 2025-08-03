package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"textlib-rl-system/internal/analyzer"
	"textlib-rl-system/internal/logging"
	"textlib-rl-system/internal/rl"
	"textlib-rl-system/internal/telemetry"
)

func main() {
	// Parse command line flags
	var (
		mode          = flag.String("mode", "train", "Mode: train, generate-report, health-check, or cleanup-logs")
		maxEpisodes   = flag.Int("episodes", 10000, "Maximum training episodes")
		logLevel      = flag.String("log-level", "info", "Logging level")
		checkpointDir = flag.String("checkpoint-dir", "./models", "Checkpoint directory")
		enableProfile = flag.Bool("profile", false, "Enable CPU/memory profiling")
		configFile    = flag.String("config", "", "Configuration file path")
		inputFile     = flag.String("input", "", "Input file for report generation")
		outputFile    = flag.String("output", "", "Output file for report generation")
		modelFile     = flag.String("model", "", "Model file for report generation")
	)
	flag.Parse()

	// Set resource limits as specified in the design
	runtime.GOMAXPROCS(2)

	// Set memory limit from environment variable
	if memLimit := os.Getenv("GOMEMLIMIT"); memLimit != "" {
		log.Printf("Memory limit set to: %s", memLimit)
	}

	// Configure logging level
	configureLogging(*logLevel)

	switch *mode {
	case "train":
		runTraining(*maxEpisodes, *checkpointDir, *enableProfile, *configFile)
	case "generate-report":
		generateReport(*inputFile, *outputFile, *modelFile)
	case "health-check":
		healthCheck()
	case "cleanup-logs":
		cleanupLogs()
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

func runTraining(maxEpisodes int, checkpointDir string, enableProfiling bool, configFile string) {
	log.Println("Starting RL training with comprehensive logging...")

	// Perform automatic log cleanup before training
	log.Println("Checking log directory size...")
	if err := logging.AutoCleanup("./logs", 50.0, 200); err != nil {
		log.Printf("Warning: log cleanup failed: %v", err)
	}

	// Initialize logging system
	logger := logging.NewInsightLogger("./logs", 100, 5*time.Second)
	if err := logger.Start(); err != nil {
		log.Fatalf("Failed to start logger: %v", err)
	}
	defer logger.Stop()

	// Initialize telemetry
	telemetryEndpoint := os.Getenv("TELEMETRY_ENDPOINT")
	telemetry := telemetry.NewTelemetryClient(telemetryEndpoint, 1000, 10*time.Second)
	if err := telemetry.Start(); err != nil {
		log.Fatalf("Failed to start telemetry: %v", err)
	}
	defer telemetry.Stop()

	// Load configuration
	config := loadConfiguration(configFile, maxEpisodes, enableProfiling)

	// Initialize RL system
	system := rl.NewEnhancedRLSystem(config)
	system.SetLogger(logger)
	system.SetTelemetry(telemetry)

	// Load training data
	trainingData := loadTrainingData()
	system.LoadTrainingData(trainingData)

	// Start training
	log.Printf("Starting training with %d episodes...", maxEpisodes)
	system.TrainWithLogging()

	// Generate final insights
	analyzer := analyzer.NewInsightAnalyzer(logger, logger.MetricsDB, maxEpisodes)
	insights := analyzer.GenerateInsights()

	// Save insights
	if err := saveInsights(insights, "./logs/insights.json"); err != nil {
		log.Printf("Failed to save insights: %v", err)
	}

	// Save final model
	if err := saveFinalModel(system, checkpointDir); err != nil {
		log.Printf("Failed to save final model: %v", err)
	}

	log.Println("Training completed successfully.")
}

func generateReport(inputFile, outputFile, modelFile string) {
	log.Println("Generating API usage guide...")

	if inputFile == "" {
		inputFile = "./logs/insights.json"
	}
	if outputFile == "" {
		outputFile = "./api-usage-guide.md"
	}

	// Load insights
	insights, err := loadInsights(inputFile)
	if err != nil {
		log.Fatalf("Failed to load insights: %v", err)
	}

	// Load model if specified
	var model interface{}
	if modelFile != "" {
		model, err = loadModel(modelFile)
		if err != nil {
			log.Printf("Warning: Failed to load model: %v", err)
		}
	}

	// Generate report
	report := generateAPIUsageGuide(insights, model)

	// Save report
	if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
		log.Fatalf("Failed to save report: %v", err)
	}

	log.Printf("API usage guide generated successfully: %s", outputFile)
}

func cleanupLogs() {
	fmt.Println("🧹 Starting log cleanup...")
	
	// Perform automatic cleanup with sensible defaults
	err := logging.AutoCleanup("./logs", 10.0, 100) // 10MB max, 100 files max
	if err != nil {
		log.Fatalf("Log cleanup failed: %v", err)
	}
	
	fmt.Println("✅ Log cleanup completed successfully")
}

func healthCheck() {
	// Perform basic health checks
	checks := []struct {
		name string
		fn   func() error
	}{
		{"Memory", checkMemory},
		{"Disk Space", checkDiskSpace},
		{"Logging Directory", checkLoggingDirectory},
		{"Models Directory", checkModelsDirectory},
	}

	allPassed := true
	for _, check := range checks {
		if err := check.fn(); err != nil {
			log.Printf("FAIL: %s - %v", check.name, err)
			allPassed = false
		} else {
			log.Printf("PASS: %s", check.name)
		}
	}

	if allPassed {
		log.Println("All health checks passed")
		os.Exit(0)
	} else {
		log.Println("Some health checks failed")
		os.Exit(1)
	}
}

func configureLogging(level string) {
	// Configure log output based on level
	switch level {
	case "debug":
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	case "info":
		log.SetFlags(log.LstdFlags)
	case "warn", "error":
		log.SetFlags(log.LstdFlags)
	default:
		log.SetFlags(log.LstdFlags)
	}
}

func loadConfiguration(configFile string, maxEpisodes int, enableProfiling bool) rl.SystemConfig {
	config := rl.SystemConfig{
		MaxEpisodes:        maxEpisodes,
		MaxStepsPerEpisode: 15,
		LoggingInterval:    100,
		CheckpointInterval: 500,
		MetricsPort:        8080,
		EnableProfiling:    enableProfiling,
	}

	// Override with environment variables
	if envEpisodes := os.Getenv("MAX_EPISODES"); envEpisodes != "" {
		if episodes, err := strconv.Atoi(envEpisodes); err == nil {
			config.MaxEpisodes = episodes
		}
	}

	if envCheckpoint := os.Getenv("CHECKPOINT_INTERVAL"); envCheckpoint != "" {
		if interval, err := strconv.Atoi(envCheckpoint); err == nil {
			config.CheckpointInterval = interval
		}
	}

	// Load from config file if specified
	if configFile != "" {
		if err := loadConfigFromFile(configFile, &config); err != nil {
			log.Printf("Warning: Failed to load config file: %v", err)
		}
	}

	return config
}

func loadConfigFromFile(filename string, config *rl.SystemConfig) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, config)
}

func loadTrainingData() []rl.TrainingExample {
	// Load comprehensive realistic training data
	return rl.GetRealisticTrainingData()
}

func saveInsights(insights analyzer.APIFeedbackReport, filename string) error {
	data, err := json.MarshalIndent(insights, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func saveFinalModel(system *rl.EnhancedRLSystem, checkpointDir string) error {
	// Create checkpoint directory if it doesn't exist
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		return err
	}

	// Save Q-table and configuration
	modelData := map[string]interface{}{
		"q_table":     system.Agent.QTable,
		"config":      system.Config,
		"timestamp":   time.Now(),
		"version":     "1.0",
	}

	data, err := json.MarshalIndent(modelData, "", "  ")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/final_model_%d.json", checkpointDir, time.Now().Unix())
	return os.WriteFile(filename, data, 0644)
}

func loadInsights(filename string) (analyzer.APIFeedbackReport, error) {
	var insights analyzer.APIFeedbackReport

	data, err := os.ReadFile(filename)
	if err != nil {
		return insights, err
	}

	err = json.Unmarshal(data, &insights)
	return insights, err
}

func loadModel(filename string) (interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var model map[string]interface{}
	err = json.Unmarshal(data, &model)
	return model, err
}

func generateAPIUsageGuide(insights analyzer.APIFeedbackReport, model interface{}) string {
	report := fmt.Sprintf(`# TextLib API Usage Guide

Generated: %s
Analysis Period: %d episodes

## Executive Summary

This report provides comprehensive insights into the optimal usage patterns for the TextLib API based on reinforcement learning analysis.

### Key Findings

- **Total Episodes Analyzed**: %d
- **Overall Success Rate**: %.2f%%
- **Average Episode Reward**: %.2f
- **Learning Efficiency**: %.2f

## Function Performance Analysis

`, insights.Timestamp.Format("2006-01-02 15:04:05"),
		insights.AnalysisPeriod,
		insights.PerformanceMetrics.TotalEpisodes,
		insights.PerformanceMetrics.OverallSuccessRate*100,
		insights.PerformanceMetrics.AvgEpisodeReward,
		insights.PerformanceMetrics.LearningEfficiency)

	// Add function statistics
	report += "### Function Usage Statistics\n\n"
	for functionName, stats := range insights.FunctionUsageStats {
		report += fmt.Sprintf("#### %s\n", functionName)
		report += fmt.Sprintf("- **Call Count**: %d\n", stats.CallCount)
		report += fmt.Sprintf("- **Success Rate**: %.2f%%\n", stats.SuccessRate*100)
		report += fmt.Sprintf("- **Average Reward**: %.2f\n", stats.AvgReward)
		report += fmt.Sprintf("- **Average Duration**: %.2f seconds\n", stats.AvgDuration)
		report += fmt.Sprintf("- **Quality Score**: %.2f\n\n", stats.QualityMetrics.AvgOutputQuality)
	}

	// Add optimal sequences
	report += "## Optimal Usage Patterns\n\n"
	for pattern, sequence := range insights.OptimalSequences {
		report += fmt.Sprintf("### %s\n", pattern)
		report += fmt.Sprintf("Sequence: %s\n\n", fmt.Sprintf("%v", sequence))
	}

	// Add recommendations
	report += "## Recommendations\n\n"
	for i, recommendation := range insights.Recommendations {
		report += fmt.Sprintf("%d. %s\n", i+1, recommendation)
	}

	// Add failure analysis
	report += "\n## Failure Analysis\n\n"
	for pattern, failure := range insights.FailureAnalysis.CommonFailures {
		report += fmt.Sprintf("### %s\n", pattern)
		report += fmt.Sprintf("- **Frequency**: %d\n", failure.Frequency)
		report += fmt.Sprintf("- **Impact**: %.2f\n", failure.Impact)
		report += "- **Suggestions**:\n"
		for _, suggestion := range failure.Suggestions {
			report += fmt.Sprintf("  - %s\n", suggestion)
		}
		report += "\n"
	}

	report += "\n---\n*Generated by TextLib RL System*\n"

	return report
}

// Health check functions
func checkMemory() error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check if memory usage is reasonable (under 400MB)
	if m.Alloc > 400*1024*1024 {
		return fmt.Errorf("memory usage too high: %d MB", m.Alloc/(1024*1024))
	}

	return nil
}

func checkDiskSpace() error {
	// Simple check - ensure we can write to current directory
	testFile := "./health_check_temp.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("cannot write to disk: %v", err)
	}

	os.Remove(testFile)
	return nil
}

func checkLoggingDirectory() error {
	if err := os.MkdirAll("./logs", 0755); err != nil {
		return fmt.Errorf("cannot create logs directory: %v", err)
	}

	return nil
}

func checkModelsDirectory() error {
	if err := os.MkdirAll("./models", 0755); err != nil {
		return fmt.Errorf("cannot create models directory: %v", err)
	}

	return nil
}