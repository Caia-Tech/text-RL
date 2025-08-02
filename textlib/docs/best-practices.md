# TextLib API Best Practices

## Quick Start: Optimal Usage Patterns

This guide provides proven patterns for getting the most out of TextLib API based on extensive performance analysis.

## Recommended Sequences

### 1. General Purpose Text Analysis
**Best for:** Articles, documentation, general content

```go
// Optimal sequence for comprehensive analysis
text := "Your content here..."

// Step 1: Validate input
validation := textlib.ValidateOutput(text)
if !validation.IsValid {
    return validation.Issues
}

// Step 2: Extract key information
entities := textlib.ExtractEntities(text)
readability := textlib.AnalyzeReadability(text)
keywords := textlib.ExtractKeywords(text)

// Step 3: Enhance with entities context (captures more entities)
enhancedEntities := textlib.ExtractEntities(text, 
    textlib.WithContext(entities, keywords))
```

### 2. Technical Documentation
**Best for:** API docs, code documentation, technical articles

```go
// Specialized sequence for technical content
func AnalyzeTechnicalDoc(text string) TechAnalysis {
    // Detect code first to identify technical content
    codeInfo := textlib.DetectCode(text)
    
    // Extract technical entities
    entities := textlib.ExtractEntities(text)
    
    // Get keywords relevant to technical domain
    keywords := textlib.ExtractKeywords(text, 
        textlib.WithDomain("technical"))
    
    return TechAnalysis{
        HasCode: codeInfo.HasCode,
        TechTerms: entities.Entities,
        Keywords: keywords.Keywords,
    }
}
```

### 3. Business Communications
**Best for:** Emails, reports, business documents

```go
// Optimized for business context
func AnalyzeBusinessDoc(text string) BusinessAnalysis {
    // Sentiment provides context
    sentiment := textlib.SentimentAnalysis(text)
    
    // Extract business entities (people, companies, products)
    entities := textlib.ExtractEntities(text,
        textlib.WithTypes("PERSON", "ORG", "PRODUCT"))
    
    // Keywords for action items and themes
    keywords := textlib.ExtractKeywords(text)
    
    return BusinessAnalysis{
        Tone: sentiment.Sentiment,
        Stakeholders: entities.People,
        Organizations: entities.Organizations,
        ActionItems: keywords.Keywords,
    }
}
```

### 4. Social Media Content
**Best for:** Posts, comments, user-generated content

```go
// Fast analysis for social content
func AnalyzeSocialPost(text string) SocialAnalysis {
    // Sentiment is primary signal
    sentiment := textlib.SentimentAnalysis(text)
    
    // Extract trending topics
    keywords := textlib.ExtractKeywords(text,
        textlib.WithLimit(5))
    
    // Clean up formatting
    formatted := textlib.FormatText(text)
    
    return SocialAnalysis{
        Sentiment: sentiment,
        Topics: keywords.Keywords,
        CleanText: formatted.Text,
    }
}
```

## Performance Tips

### 1. Order Matters
Functions perform better when called in specific sequences:
- ✅ `ValidateOutput` → `ExtractEntities` → `AnalyzeReadability`
- ❌ `SummarizeText` → `ExtractEntities` (loses information)

### 2. Avoid Redundancy
- Don't call the same function multiple times (except `ExtractEntities` with context)
- Each function should add new information

### 3. Budget-Conscious Usage
When working with rate limits or cost constraints:

```go
// Low-cost validation and analysis
func QuickAnalysis(text string) QuickResult {
    // These are the most cost-effective functions
    valid := textlib.ValidateOutput(text)      // Lowest cost
    code := textlib.DetectCode(text)          // Low cost
    readable := textlib.AnalyzeReadability(text) // Medium cost
    
    return QuickResult{valid, code, readable}
}
```

### 4. Context Enhancement Pattern
Some functions benefit from prior analysis:

```go
// Enhanced extraction pattern
basicEntities := textlib.ExtractEntities(text)
keywords := textlib.ExtractKeywords(text)

// Second pass with context finds 20-30% more entities
enhancedEntities := textlib.ExtractEntities(text,
    textlib.WithContext(basicEntities, keywords))
```

## Common Use Cases

### Document Classification
```go
func ClassifyDocument(text string) string {
    code := textlib.DetectCode(text)
    if code.HasCode {
        return "technical"
    }
    
    sentiment := textlib.SentimentAnalysis(text)
    if sentiment.Score > 0.3 || sentiment.Score < -0.3 {
        return "opinion"
    }
    
    readability := textlib.AnalyzeReadability(text)
    if readability.Score > 80 {
        return "casual"
    }
    
    return "formal"
}
```

### Comprehensive Analysis
```go
func FullAnalysis(text string) CompleteAnalysis {
    // Validate first
    validation := textlib.ValidateOutput(text)
    if !validation.IsValid {
        return CompleteAnalysis{Error: validation.Issues}
    }
    
    // Core analysis
    entities := textlib.ExtractEntities(text)
    readability := textlib.AnalyzeReadability(text)
    keywords := textlib.ExtractKeywords(text)
    
    // Enhanced extraction
    enhancedEntities := textlib.ExtractEntities(text,
        textlib.WithContext(entities, keywords))
    
    // Synthesis
    summary := textlib.SummarizeText(text,
        textlib.WithContext(readability))
    
    return CompleteAnalysis{
        Entities: enhancedEntities,
        Keywords: keywords,
        Readability: readability,
        Summary: summary,
    }
}
```

### Quick Insights
```go
func QuickInsights(text string) Insights {
    // Minimal calls for maximum insight
    entities := textlib.ExtractEntities(text)
    keywords := textlib.ExtractKeywords(text)
    validation := textlib.ValidateOutput(text)
    
    return Insights{
        MainTopics: keywords.Top(3),
        KeyEntities: entities.MostFrequent(3),
        Quality: validation.Score,
    }
}
```

## Error Handling

```go
func SafeAnalysis(text string) (Result, error) {
    // Always validate first
    validation := textlib.ValidateOutput(text)
    if !validation.IsValid {
        return Result{}, fmt.Errorf("invalid input: %v", validation.Issues)
    }
    
    // Handle API errors gracefully
    entities, err := textlib.ExtractEntities(text)
    if err != nil {
        // Fallback to keywords only
        keywords, _ := textlib.ExtractKeywords(text)
        return Result{Keywords: keywords}, nil
    }
    
    return Result{Entities: entities}, nil
}
```

## Performance Benchmarks

Typical processing times for 1KB text:
- `ValidateOutput`: ~20ms
- `DetectCode`: ~30ms
- `ExtractEntities`: ~60ms
- `AnalyzeReadability`: ~40ms
- `ExtractKeywords`: ~50ms
- `SentimentAnalysis`: ~40ms
- `SummarizeText`: ~90ms

## Summary

1. **Start with validation** - It's fast and prevents wasted processing
2. **Extract entities early** - They provide context for other functions
3. **Use the right sequence** - Order significantly impacts quality
4. **Consider your use case** - Different patterns for different content types
5. **Reuse context** - Some functions perform better with prior analysis

---
*TextLib API v2.0 - Best Practices Guide*