package rl

import (
	"fmt"
	"strings"
	"time"
)

func NewActionSimulator() *ActionSimulator {
	return &ActionSimulator{
		Functions: map[string]SimulatedFunction{
			"extract_entities": {
				Name:            "extract_entities",
				Category:        "analysis",
				Cost:            5,
				BaseSuccessRate: 0.85,
				OutputGenerator: simulateEntityExtraction,
			},
			"analyze_readability": {
				Name:            "analyze_readability",
				Category:        "analysis",
				Cost:            3,
				BaseSuccessRate: 0.90,
				OutputGenerator: simulateReadabilityAnalysis,
			},
			"detect_code": {
				Name:            "detect_code",
				Category:        "analysis",
				Cost:            2,
				BaseSuccessRate: 0.95,
				OutputGenerator: simulateCodeDetection,
			},
			"extract_keywords": {
				Name:            "extract_keywords",
				Category:        "analysis",
				Cost:            4,
				BaseSuccessRate: 0.88,
				OutputGenerator: simulateKeywordExtraction,
			},
			"sentiment_analysis": {
				Name:            "sentiment_analysis",
				Category:        "analysis",
				Cost:            3,
				BaseSuccessRate: 0.82,
				OutputGenerator: simulateSentimentAnalysis,
			},
			"summarize_text": {
				Name:            "summarize_text",
				Category:        "generation",
				Cost:            8,
				BaseSuccessRate: 0.75,
				OutputGenerator: simulateTextSummary,
			},
			"format_text": {
				Name:            "format_text",
				Category:        "formatting",
				Cost:            2,
				BaseSuccessRate: 0.98,
				OutputGenerator: simulateTextFormatting,
			},
			"validate_output": {
				Name:            "validate_output",
				Category:        "utility",
				Cost:            1,
				BaseSuccessRate: 0.99,
				OutputGenerator: simulateOutputValidation,
			},
		},
	}
}

func (sim *ActionSimulator) ExecuteAction(action Action, input string, params map[string]interface{}) ActionResult {
	startTime := time.Now()
	
	function, exists := sim.Functions[action.FunctionName]
	if !exists {
		return ActionResult{
			Success:    false,
			Output:     nil,
			Error:      fmt.Sprintf("unknown function: %s", action.FunctionName),
			Duration:   time.Since(startTime),
			MemoryUsed: 1024,
		}
	}
	
	// Simulate execution time based on input size and function complexity
	executionTime := time.Duration(len(input)/100+action.Cost*10) * time.Millisecond
	time.Sleep(executionTime)
	
	// Determine success based on base success rate
	success := simulateSuccess(function.BaseSuccessRate)
	
	var output interface{}
	var errorMsg string
	
	if success {
		var err error
		output, err = function.OutputGenerator(input, params)
		if err != nil {
			success = false
			errorMsg = err.Error()
		}
	} else {
		errorMsg = fmt.Sprintf("simulated failure for %s", action.FunctionName)
	}
	
	return ActionResult{
		Success:    success,
		Output:     output,
		Error:      errorMsg,
		Duration:   time.Since(startTime),
		MemoryUsed: int64(len(input) * 2), // Simplified memory calculation
	}
}

func simulateSuccess(baseRate float64) bool {
	// Add some randomness to success rate
	return (baseRate + (0.1 * (0.5 - float64(time.Now().UnixNano()%1000)/1000.0))) > 0.5
}

func simulateEntityExtraction(input string, params map[string]interface{}) (interface{}, error) {
	words := strings.Fields(input)
	entities := []map[string]interface{}{}
	
	for i, word := range words {
		if len(word) > 5 && i%3 == 0 { // Simulate entity detection
			entities = append(entities, map[string]interface{}{
				"text":  word,
				"type":  "ENTITY",
				"start": i,
				"end":   i + len(word),
			})
		}
	}
	
	return map[string]interface{}{
		"entities": entities,
		"count":    len(entities),
	}, nil
}

func simulateReadabilityAnalysis(input string, params map[string]interface{}) (interface{}, error) {
	words := len(strings.Fields(input))
	sentences := strings.Count(input, ".") + strings.Count(input, "!") + strings.Count(input, "?")
	if sentences == 0 {
		sentences = 1
	}
	
	avgWordsPerSentence := float64(words) / float64(sentences)
	readabilityScore := 100.0 - (avgWordsPerSentence * 2.0)
	
	if readabilityScore < 0 {
		readabilityScore = 0
	}
	if readabilityScore > 100 {
		readabilityScore = 100
	}
	
	return map[string]interface{}{
		"readability_score":      readabilityScore,
		"avg_words_per_sentence": avgWordsPerSentence,
		"total_words":           words,
		"total_sentences":       sentences,
		"level":                 getReadabilityLevel(readabilityScore),
	}, nil
}

func simulateCodeDetection(input string, params map[string]interface{}) (interface{}, error) {
	codeIndicators := []string{"function", "class", "def", "import", "var", "const", "let", "if", "for", "while"}
	codeBlocks := []map[string]interface{}{}
	
	for _, indicator := range codeIndicators {
		if strings.Contains(strings.ToLower(input), indicator) {
			codeBlocks = append(codeBlocks, map[string]interface{}{
				"type":      "code_snippet",
				"language":  "unknown",
				"indicator": indicator,
			})
		}
	}
	
	return map[string]interface{}{
		"has_code":    len(codeBlocks) > 0,
		"code_blocks": codeBlocks,
		"confidence":  float64(len(codeBlocks)) / float64(len(codeIndicators)),
	}, nil
}

func simulateKeywordExtraction(input string, params map[string]interface{}) (interface{}, error) {
	words := strings.Fields(input)
	keywords := []map[string]interface{}{}
	
	for _, word := range words {
		if len(word) > 4 && !isCommonWord(word) {
			keywords = append(keywords, map[string]interface{}{
				"keyword": word,
				"score":   float64(len(word)) / 10.0,
			})
		}
	}
	
	// Limit to top 10 keywords
	if len(keywords) > 10 {
		keywords = keywords[:10]
	}
	
	return map[string]interface{}{
		"keywords": keywords,
		"count":    len(keywords),
	}, nil
}

func simulateSentimentAnalysis(input string, params map[string]interface{}) (interface{}, error) {
	positiveWords := []string{"good", "great", "excellent", "amazing", "wonderful", "fantastic"}
	negativeWords := []string{"bad", "terrible", "awful", "horrible", "disappointing", "poor"}
	
	positiveCount := 0
	negativeCount := 0
	
	lowerInput := strings.ToLower(input)
	for _, word := range positiveWords {
		positiveCount += strings.Count(lowerInput, word)
	}
	for _, word := range negativeWords {
		negativeCount += strings.Count(lowerInput, word)
	}
	
	sentiment := "neutral"
	score := 0.0
	
	if positiveCount > negativeCount {
		sentiment = "positive"
		score = float64(positiveCount) / float64(positiveCount+negativeCount+1)
	} else if negativeCount > positiveCount {
		sentiment = "negative"
		score = -float64(negativeCount) / float64(positiveCount+negativeCount+1)
	}
	
	return map[string]interface{}{
		"sentiment":       sentiment,
		"score":          score,
		"positive_count": positiveCount,
		"negative_count": negativeCount,
		"confidence":     0.75,
	}, nil
}

func simulateTextSummary(input string, params map[string]interface{}) (interface{}, error) {
	sentences := strings.Split(input, ".")
	if len(sentences) <= 2 {
		return map[string]interface{}{
			"summary": input,
			"ratio":   1.0,
		}, nil
	}
	
	// Take first and last sentence as a simple summary
	summary := strings.TrimSpace(sentences[0]) + ". " + strings.TrimSpace(sentences[len(sentences)-2]) + "."
	
	return map[string]interface{}{
		"summary":          summary,
		"original_length":  len(input),
		"summary_length":   len(summary),
		"compression_ratio": float64(len(summary)) / float64(len(input)),
	}, nil
}

func simulateTextFormatting(input string, params map[string]interface{}) (interface{}, error) {
	formatted := strings.TrimSpace(input)
	formatted = strings.ReplaceAll(formatted, "  ", " ") // Remove double spaces
	
	return map[string]interface{}{
		"formatted_text": formatted,
		"changes_made":   []string{"trimmed_whitespace", "normalized_spaces"},
	}, nil
}

func simulateOutputValidation(input string, params map[string]interface{}) (interface{}, error) {
	isValid := len(input) > 0 && len(input) < 10000
	issues := []string{}
	
	if len(input) == 0 {
		issues = append(issues, "empty_input")
	}
	if len(input) > 10000 {
		issues = append(issues, "input_too_long")
	}
	
	return map[string]interface{}{
		"is_valid": isValid,
		"issues":   issues,
		"score":    0.95,
	}, nil
}

func getReadabilityLevel(score float64) string {
	if score >= 90 {
		return "very_easy"
	} else if score >= 80 {
		return "easy"
	} else if score >= 70 {
		return "fairly_easy"
	} else if score >= 60 {
		return "standard"
	} else if score >= 50 {
		return "fairly_difficult"
	} else if score >= 30 {
		return "difficult"
	}
	return "very_difficult"
}

func isCommonWord(word string) bool {
	commonWords := []string{"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"}
	lowerWord := strings.ToLower(word)
	for _, common := range commonWords {
		if lowerWord == common {
			return true
		}
	}
	return false
}