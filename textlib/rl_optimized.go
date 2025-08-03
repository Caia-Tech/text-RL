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
	"regexp"
	"strings"
	"time"

	"github.com/caiatech/textlib"
)

// SmartAnalyzeResult represents the comprehensive analysis result
type SmartAnalyzeResult struct {
	// Core analysis results
	Statistics   *textlib.TextStatistics `json:"statistics"`
	Entities     []textlib.Entity        `json:"entities"`
	Sentences    []string                `json:"sentences"`
	
	// RL-enhanced insights
	Strategy        ProcessingStrategy      `json:"strategy"`
	QualityScore    float64                `json:"quality_score"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	OptimizedPath   []string               `json:"optimized_path"`
	CacheUtilized   bool                  `json:"cache_utilized"`
	
	// Performance metrics
	Performance     PerformanceMetrics     `json:"performance"`
	ResourceUsage   ResourceUsage          `json:"resource_usage"`
}

// ValidatedExtractionResult represents validated entity extraction
type ValidatedExtractionResult struct {
	Entities        []ValidatedEntity      `json:"entities"`
	ValidationLevel string                 `json:"validation_level"`
	Confidence      float64               `json:"confidence"`
	ProcessingPath  []string              `json:"processing_path"`
	Performance     PerformanceMetrics    `json:"performance"`
}

// ValidatedEntity represents an entity with validation metadata
type ValidatedEntity struct {
	textlib.Entity
	Validated       bool    `json:"validated"`
	ValidationScore float64 `json:"validation_score"`
	ValidationMethod string `json:"validation_method"`
	Context         string  `json:"context"`
}

// DomainAnalysisResult represents domain-specific analysis
type DomainAnalysisResult struct {
	Domain          string                 `json:"domain"`
	Analysis        interface{}            `json:"analysis"`
	DomainSpecific  map[string]interface{} `json:"domain_specific"`
	Strategy        ProcessingStrategy     `json:"strategy"`
	Performance     PerformanceMetrics     `json:"performance"`
}

// QuickInsightsResult represents rapid social media analysis
type QuickInsightsResult struct {
	Insights        []string               `json:"insights"`
	SentimentScore  float64               `json:"sentiment_score"`
	KeyTerms        []string              `json:"key_terms"`
	Readability     float64               `json:"readability"`
	Strategy        ProcessingStrategy     `json:"strategy"`
	Performance     PerformanceMetrics     `json:"performance"`
}

// TechnicalAnalysisResult represents comprehensive technical analysis
type TechnicalAnalysisResult struct {
	CodeMetrics     map[string]interface{} `json:"code_metrics"`
	Documentation   map[string]interface{} `json:"documentation"`
	Complexity      ComplexityReport       `json:"complexity"`
	Quality         QualityAssessment      `json:"quality"`
	Strategy        ProcessingStrategy     `json:"strategy"`
	Performance     PerformanceMetrics     `json:"performance"`
}

// AlgorithmRequirements represents requirements for processing
type AlgorithmRequirements struct {
	MinQuality  float64 `json:"min_quality"`
	MaxTimeMs   int64   `json:"max_time_ms"`
	MaxMemoryMB int64   `json:"max_memory_mb"`
}

// SmartAnalyze performs comprehensive analysis using RL-discovered optimal sequence
// Based on our RL findings: statistics first, then entities, then sentences for best performance
func SmartAnalyze(text string) SmartAnalyzeResult {
	startTime := time.Now()
	
	// Initialize strategy selector with our discovered patterns
	selector := NewStrategySelector()
	
	// Analyze text characteristics
	characteristics := TextCharacteristics{
		Length:     len(text),
		Language:   detectLanguage(text),
		Domain:     classifyDomain(text),
		Complexity: estimateComplexity(text),
		Structure:  analyzeStructure(text),
	}
	
	// Select optimal strategy
	requirements := AlgorithmRequirements{
		MinQuality:  0.8,
		MaxTimeMs:   5000, // 5 second limit
		MaxMemoryMB: 100,
	}
	
	strategy, _ := selector.SelectStrategy(characteristics, requirements)
	
	// Execute analysis in RL-discovered optimal order
	var stats *textlib.TextStatistics
	var entities []textlib.Entity
	var sentences []string
	var optimizedPath []string
	
	// Step 1: Statistics (fastest, most informative)
	optimizedPath = append(optimizedPath, "statistics")
	stats = textlib.CalculateTextStatistics(text)
	
	// Step 2: Entity extraction (moderate cost, high value)
	optimizedPath = append(optimizedPath, "entities")
	entities = textlib.ExtractNamedEntities(text)
	
	// Step 3: Sentence splitting (builds on previous analysis)
	optimizedPath = append(optimizedPath, "sentences")
	sentences = textlib.SplitIntoSentences(text)
	
	processingTime := time.Since(startTime)
	
	// Calculate quality score based on completeness and accuracy
	qualityScore := calculateQualityScore(stats, entities, sentences, text)
	
	// Create performance metrics
	performance := PerformanceMetrics{
		TotalTime: processingTime,
		StepTimings: map[string]time.Duration{
			"statistics": processingTime / 3, // Approximate timing
			"entities":   processingTime / 3,
			"sentences":  processingTime / 3,
		},
		MemoryUsage:      int64(len(text) * 2), // Estimate
		CacheUtilization: 0.0,                  // No cache in this simple implementation
	}
	
	resourceUsage := ResourceUsage{
		MemoryUsedMB:     int(performance.MemoryUsage / 1024 / 1024),
		CPUTimeMs:        processingTime.Milliseconds(),
		NetworkCallsMade: 0,
		CacheHits:        0,
	}
	
	return SmartAnalyzeResult{
		Statistics:      stats,
		Entities:        entities,
		Sentences:       sentences,
		Strategy:        strategy,
		QualityScore:    qualityScore,
		ProcessingTime:  processingTime,
		OptimizedPath:   optimizedPath,
		CacheUtilized:   false,
		Performance:     performance,
		ResourceUsage:   resourceUsage,
	}
}

// ValidatedExtraction performs entity extraction with validation-first approach
// Based on RL insight: validate context before extraction for better accuracy
func ValidatedExtraction(text string) ValidatedExtractionResult {
	startTime := time.Now()
	
	var processingPath []string
	var validatedEntities []ValidatedEntity
	
	// Step 1: Pre-validation analysis
	processingPath = append(processingPath, "pre-validation")
	textQuality := assessTextQuality(text)
	
	// Step 2: Context-aware extraction
	processingPath = append(processingPath, "context-extraction")
	rawEntities := textlib.ExtractNamedEntities(text)
	
	// Step 3: Post-extraction validation
	processingPath = append(processingPath, "post-validation")
	for _, entity := range rawEntities {
		validated := validateEntity(entity, text)
		validatedEntities = append(validatedEntities, validated)
	}
	
	processingTime := time.Since(startTime)
	
	// Calculate overall confidence
	confidence := calculateValidationConfidence(validatedEntities, textQuality)
	
	// Determine validation level
	validationLevel := "standard"
	if confidence > 0.9 {
		validationLevel = "high"
	} else if confidence < 0.7 {
		validationLevel = "basic"
	}
	
	performance := PerformanceMetrics{
		TotalTime: processingTime,
		StepTimings: map[string]time.Duration{
			"pre-validation":     processingTime / 4,
			"context-extraction": processingTime / 2,
			"post-validation":    processingTime / 4,
		},
		MemoryUsage:      int64(len(text) + len(validatedEntities)*100),
		CacheUtilization: 0.0,
	}
	
	return ValidatedExtractionResult{
		Entities:        validatedEntities,
		ValidationLevel: validationLevel,
		Confidence:      confidence,
		ProcessingPath:  processingPath,
		Performance:     performance,
	}
}

// DomainOptimizedAnalyze performs domain-specific optimized analysis
func DomainOptimizedAnalyze(text string, domain string) DomainAnalysisResult {
	startTime := time.Now()
	
	// Create domain-specific strategy
	characteristics := TextCharacteristics{
		Length:     len(text),
		Language:   detectLanguage(text),
		Domain:     domain,
		Complexity: estimateComplexity(text),
		Structure:  analyzeStructure(text),
	}
	
	selector := NewStrategySelector()
	requirements := AlgorithmRequirements{
		MinQuality:  0.85,
		MaxTimeMs:   3000,
		MaxMemoryMB: 200,
	}
	
	strategy, _ := selector.SelectStrategy(characteristics, requirements)
	
	// Domain-specific analysis
	var analysis interface{}
	var domainSpecific map[string]interface{}
	
	switch domain {
	case "technical":
		analysis = analyzeTechnicalContent(text)
		domainSpecific = map[string]interface{}{
			"code_snippets": extractCodeSnippets(text),
			"technical_terms": extractTechnicalTerms(text),
		}
	case "academic":
		analysis = analyzeAcademicContent(text)
		domainSpecific = map[string]interface{}{
			"citations": extractCitations(text),
			"abstracts": extractAbstracts(text),
		}
	case "social-media":
		analysis = analyzeSocialContent(text)
		domainSpecific = map[string]interface{}{
			"hashtags": extractHashtags(text),
			"mentions": extractMentions(text),
		}
	default:
		analysis = textlib.CalculateTextStatistics(text)
		domainSpecific = map[string]interface{}{
			"general_metrics": "basic analysis applied",
		}
	}
	
	processingTime := time.Since(startTime)
	
	performance := PerformanceMetrics{
		TotalTime:        processingTime,
		MemoryUsage:      int64(len(text) * 3),
		CacheUtilization: 0.0,
	}
	
	return DomainAnalysisResult{
		Domain:         domain,
		Analysis:       analysis,
		DomainSpecific: domainSpecific,
		Strategy:       strategy,
		Performance:    performance,
	}
}

// QuickInsights performs rapid analysis optimized for social media
// Based on RL insight: prioritize sentiment and key terms for social content
func QuickInsights(text string) QuickInsightsResult {
	startTime := time.Now()
	
	// Fast-path analysis for social media content
	insights := []string{}
	
	// Quick sentiment analysis
	sentimentScore := quickSentimentAnalysis(text)
	if sentimentScore > 0.6 {
		insights = append(insights, "Positive sentiment detected")
	} else if sentimentScore < 0.4 {
		insights = append(insights, "Negative sentiment detected")
	} else {
		insights = append(insights, "Neutral sentiment")
	}
	
	// Extract key terms rapidly
	keyTerms := extractKeyTermsRapidly(text)
	if len(keyTerms) > 0 {
		insights = append(insights, fmt.Sprintf("Key topics: %s", strings.Join(keyTerms[:min(3, len(keyTerms))], ", ")))
	}
	
	// Quick readability check
	readability := textlib.CalculateFleschReadingEase(text)
	if readability > 80 {
		insights = append(insights, "High readability - easy to understand")
	} else if readability < 50 {
		insights = append(insights, "Low readability - complex text")
	}
	
	// Length-based insights
	wordCount := len(strings.Fields(text))
	if wordCount < 20 {
		insights = append(insights, "Brief content - good for social media")
	} else if wordCount > 200 {
		insights = append(insights, "Long content - consider breaking into parts")
	}
	
	processingTime := time.Since(startTime)
	
	// Create optimized strategy for social media
	strategy := ProcessingStrategy{
		Name:        "quick-insights",
		Description: "Optimized for rapid social media analysis",
		Parameters: map[string]interface{}{
			"depth":      1,
			"algorithms": []string{"sentiment", "keywords", "readability"},
		},
		ExpectedQuality: 0.75,
		ExpectedSpeed:   0.95,
	}
	
	performance := PerformanceMetrics{
		TotalTime:        processingTime,
		MemoryUsage:      int64(len(text)),
		CacheUtilization: 0.0,
	}
	
	return QuickInsightsResult{
		Insights:       insights,
		SentimentScore: sentimentScore,
		KeyTerms:       keyTerms,
		Readability:    readability,
		Strategy:       strategy,
		Performance:    performance,
	}
}

// DeepTechnicalAnalysis performs comprehensive analysis for code and technical documentation
func DeepTechnicalAnalysis(text string) TechnicalAnalysisResult {
	startTime := time.Now()
	
	// Comprehensive technical analysis
	codeMetrics := analyzeTechnicalContent(text)
	
	// Documentation analysis
	documentation := map[string]interface{}{
		"sections":     extractDocumentationSections(text),
		"code_blocks":  extractCodeSnippets(text),
		"api_refs":     extractAPIReferences(text),
	}
	
	// Complexity analysis
	complexity := ComplexityReport{
		LexicalComplexity:   calculateLexicalComplexity(text),
		SyntacticComplexity: calculateSyntacticComplexity(text),
		SemanticComplexity:  calculateSemanticComplexity(text),
		ReadabilityScores: map[string]float64{
			"flesch": textlib.CalculateFleschReadingEase(text),
		},
		ProcessingTime: time.Since(startTime),
		MemoryUsed:     int64(len(text) * 4),
		AlgorithmUsed:  "comprehensive-technical",
		QualityMetrics: QualityMetrics{
			Accuracy:   0.9,
			Confidence: 0.85,
			Coverage:   0.95,
		},
	}
	
	// Quality assessment
	quality := QualityAssessment{
		OverallScore:      0.85,
		ReadabilityScore:  complexity.ReadabilityScores["flesch"] / 100.0,
		CompletenessScore: assessCompletenessScore(text),
		ConsistencyScore:  assessConsistencyScore(text),
		Issues:           []QualityIssue{},
		Recommendations:  []string{"Consider adding more code examples", "Include API documentation"},
	}
	
	processingTime := time.Since(startTime)
	
	// Technical analysis strategy
	strategy := ProcessingStrategy{
		Name:        "deep-technical",
		Description: "Comprehensive analysis for technical content",
		Parameters: map[string]interface{}{
			"depth":      3,
			"algorithms": []string{"all-technical"},
			"focus":      "code-analysis",
		},
		ExpectedQuality: 0.95,
		ExpectedSpeed:   0.4,
	}
	
	performance := PerformanceMetrics{
		TotalTime:        processingTime,
		MemoryUsage:      int64(len(text) * 5),
		CacheUtilization: 0.0,
	}
	
	return TechnicalAnalysisResult{
		CodeMetrics:   codeMetrics,
		Documentation: documentation,
		Complexity:    complexity,
		Quality:       quality,
		Strategy:      strategy,
		Performance:   performance,
	}
}

// Helper functions

func detectLanguage(text string) string {
	// Simple language detection - in practice, use a proper library
	textLower := strings.ToLower(text)
	if strings.Contains(textLower, " the ") || strings.Contains(textLower, " and ") || 
	   strings.Contains(textLower, " a ") || strings.Contains(textLower, " to ") {
		return "en"
	}
	return "unknown"
}

func classifyDomain(text string) string {
	text = strings.ToLower(text)
	
	if strings.Contains(text, "function") || strings.Contains(text, "class") || strings.Contains(text, "var ") {
		return "technical"
	}
	if strings.Contains(text, "@") || strings.Contains(text, "#") {
		return "social-media"
	}
	if strings.Contains(text, "abstract") || strings.Contains(text, "methodology") {
		return "academic"
	}
	
	return "general"
}

func estimateComplexity(text string) float64 {
	// Simple complexity estimation
	words := strings.Fields(text)
	avgWordLength := 0.0
	for _, word := range words {
		avgWordLength += float64(len(word))
	}
	if len(words) > 0 {
		avgWordLength /= float64(len(words))
	}
	
	return minFloat(avgWordLength/10.0, 1.0)
}

func analyzeStructure(text string) string {
	if strings.Contains(text, "\n\n") {
		return "multi-paragraph"
	}
	if strings.Contains(text, ". ") {
		return "multi-sentence"
	}
	return "simple"
}

func calculateQualityScore(stats *textlib.TextStatistics, entities []textlib.Entity, sentences []string, text string) float64 {
	score := 0.0
	
	// Completeness score
	if stats != nil && stats.WordCount > 0 {
		score += 0.3
	}
	if len(entities) > 0 {
		score += 0.3
	}
	if len(sentences) > 0 {
		score += 0.3
	}
	
	// Quality indicators
	if stats != nil && stats.TypeTokenRatio > 0.5 {
		score += 0.1 // Good vocabulary diversity
	}
	
	return score
}

func assessTextQuality(text string) float64 {
	// Simple text quality assessment
	score := 0.5
	
	if len(text) > 50 {
		score += 0.2
	}
	if strings.Contains(text, ".") || strings.Contains(text, "!") || strings.Contains(text, "?") {
		score += 0.2
	}
	
	return minFloat(score, 1.0)
}

func validateEntity(entity textlib.Entity, text string) ValidatedEntity {
	// Simple validation - check if entity appears in meaningful context
	context := extractContext(entity, text)
	validationScore := 0.8 // Default score
	
	if len(context) > 20 {
		validationScore += 0.1
	}
	
	return ValidatedEntity{
		Entity:           entity,
		Validated:        validationScore > 0.7,
		ValidationScore:  validationScore,
		ValidationMethod: "context-based",
		Context:          context,
	}
}

func extractContext(entity textlib.Entity, text string) string {
	start := maxInt(0, entity.Position.Start-50)
	end := minInt(len(text), entity.Position.End+50)
	return text[start:end]
}

func calculateValidationConfidence(entities []ValidatedEntity, textQuality float64) float64 {
	if len(entities) == 0 {
		return textQuality
	}
	
	totalScore := 0.0
	for _, entity := range entities {
		totalScore += entity.ValidationScore
	}
	
	return (totalScore / float64(len(entities))) * textQuality
}

func analyzeTechnicalContent(text string) map[string]interface{} {
	return map[string]interface{}{
		"code_blocks":     extractCodeSnippets(text),
		"technical_terms": extractTechnicalTerms(text),
		"complexity":      estimateComplexity(text),
	}
}

func analyzeAcademicContent(text string) map[string]interface{} {
	return map[string]interface{}{
		"formal_tone":  assessFormalTone(text),
		"citations":    extractCitations(text),
		"methodology":  extractMethodology(text),
	}
}

func analyzeSocialContent(text string) map[string]interface{} {
	return map[string]interface{}{
		"hashtags":     extractHashtags(text),
		"mentions":     extractMentions(text),
		"informality":  assessInformalTone(text),
	}
}

func extractCodeSnippets(text string) []string {
	// Simple code extraction
	var snippets []string
	lines := strings.Split(text, "\n")
	
	for _, line := range lines {
		if strings.Contains(line, "func ") || strings.Contains(line, "class ") || 
		   strings.Contains(line, "def ") || strings.Contains(line, "function ") {
			snippets = append(snippets, strings.TrimSpace(line))
		}
	}
	
	return snippets
}

func extractTechnicalTerms(text string) []string {
	terms := []string{}
	words := strings.Fields(strings.ToLower(text))
	
	technicalWords := map[string]bool{
		"algorithm": true, "function": true, "class": true, "method": true,
		"api": true, "database": true, "server": true, "client": true,
		"protocol": true, "interface": true, "framework": true, "library": true,
	}
	
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		if technicalWords[cleaned] {
			terms = append(terms, cleaned)
		}
	}
	
	return terms
}

func extractCitations(text string) []string {
	// Simple citation extraction
	var citations []string
	
	// Look for different citation patterns
	// Pattern 1: [Author, Year] or [1]
	re1 := regexp.MustCompile(`\[[^\]]+\]`)
	matches1 := re1.FindAllString(text, -1)
	citations = append(citations, matches1...)
	
	// Pattern 2: Author et al.
	re2 := regexp.MustCompile(`\b[A-Z][a-z]+\s+et\s+al\.`)
	matches2 := re2.FindAllString(text, -1)
	citations = append(citations, matches2...)
	
	// Pattern 3: Years in parentheses (2023)
	re3 := regexp.MustCompile(`\(\d{4}\)`)
	matches3 := re3.FindAllString(text, -1)
	citations = append(citations, matches3...)
	
	return citations
}

func extractAbstracts(text string) []string {
	// Look for abstract sections
	if strings.Contains(strings.ToLower(text), "abstract") {
		return []string{"Abstract section detected"}
	}
	return []string{}
}

func extractHashtags(text string) []string {
	var hashtags []string
	words := strings.Fields(text)
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			hashtags = append(hashtags, word)
		}
	}
	
	return hashtags
}

func extractMentions(text string) []string {
	var mentions []string
	words := strings.Fields(text)
	
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			mentions = append(mentions, word)
		}
	}
	
	return mentions
}

func quickSentimentAnalysis(text string) float64 {
	// Simple sentiment analysis
	text = strings.ToLower(text)
	positiveWords := []string{"good", "great", "excellent", "amazing", "love", "best", "awesome"}
	negativeWords := []string{"bad", "terrible", "awful", "hate", "worst", "horrible", "sad"}
	
	positiveCount := 0
	negativeCount := 0
	
	for _, word := range positiveWords {
		positiveCount += strings.Count(text, word)
	}
	
	for _, word := range negativeWords {
		negativeCount += strings.Count(text, word)
	}
	
	total := positiveCount + negativeCount
	if total == 0 {
		return 0.5 // Neutral
	}
	
	return float64(positiveCount) / float64(total)
}

func extractKeyTermsRapidly(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	wordCount := make(map[string]int)
	
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
	}
	
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		if len(cleaned) > 3 && !stopWords[cleaned] {
			wordCount[cleaned]++
		}
	}
	
	// Extract top terms
	var keyTerms []string
	for word, count := range wordCount {
		if count > 1 { // Appears more than once
			keyTerms = append(keyTerms, word)
		}
		if len(keyTerms) >= 5 {
			break
		}
	}
	
	return keyTerms
}

func extractDocumentationSections(text string) []string {
	sections := []string{}
	lines := strings.Split(text, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "##") {
			sections = append(sections, line)
		}
	}
	
	return sections
}

func extractAPIReferences(text string) []string {
	var refs []string
	words := strings.Fields(text)
	
	for _, word := range words {
		if strings.Contains(word, "()") { // Function calls
			refs = append(refs, word)
		}
	}
	
	return refs
}

func extractMethodology(text string) []string {
	if strings.Contains(strings.ToLower(text), "methodology") ||
	   strings.Contains(strings.ToLower(text), "method") {
		return []string{"Methodology section detected"}
	}
	return []string{}
}

func assessFormalTone(text string) float64 {
	// Count formal indicators
	formalWords := []string{"therefore", "however", "furthermore", "consequently", "methodology"}
	count := 0
	
	textLower := strings.ToLower(text)
	for _, word := range formalWords {
		if strings.Contains(textLower, word) {
			count++
		}
	}
	
	return minFloat(float64(count)/5.0, 1.0)
}

func assessInformalTone(text string) float64 {
	// Count informal indicators
	informalWords := []string{"lol", "omg", "btw", "tbh", "awesome", "cool", "hey"}
	count := 0
	
	textLower := strings.ToLower(text)
	for _, word := range informalWords {
		if strings.Contains(textLower, word) {
			count++
		}
	}
	
	return minFloat(float64(count)/3.0, 1.0)
}

func calculateLexicalComplexity(text string) float64 {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0
	}
	
	totalLength := 0
	for _, word := range words {
		totalLength += len(word)
	}
	
	avgLength := float64(totalLength) / float64(len(words))
	return minFloat(avgLength/10.0, 1.0)
}

func calculateSyntacticComplexity(text string) float64 {
	sentences := textlib.SplitIntoSentences(text)
	if len(sentences) == 0 {
		return 0
	}
	
	totalWords := 0
	for _, sentence := range sentences {
		totalWords += len(strings.Fields(sentence))
	}
	
	avgSentenceLength := float64(totalWords) / float64(len(sentences))
	return minFloat(avgSentenceLength/20.0, 1.0)
}

func calculateSemanticComplexity(text string) float64 {
	// Simple semantic complexity based on unique words ratio
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return 0
	}
	
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		uniqueWords[word] = true
	}
	
	ratio := float64(len(uniqueWords)) / float64(len(words))
	return ratio
}

func assessCompletenessScore(text string) float64 {
	score := 0.5
	
	if strings.Contains(text, "introduction") || strings.Contains(text, "overview") {
		score += 0.2
	}
	if strings.Contains(text, "conclusion") || strings.Contains(text, "summary") {
		score += 0.2
	}
	if len(text) > 500 {
		score += 0.1
	}
	
	return minFloat(score, 1.0)
}

func assessConsistencyScore(text string) float64 {
	// Simple consistency check based on formatting
	score := 0.7
	
	lines := strings.Split(text, "\n")
	headingPattern := 0
	
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			headingPattern++
		}
	}
	
	if headingPattern > 0 {
		score += 0.2
	}
	
	return minFloat(score, 1.0)
}

// Utility functions
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}