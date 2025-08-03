package benchmarks

import (
	"fmt"
	"testing"
	"github.com/caiatech/textlib"
)

// Test to discover actual API function signatures and return types
func TestAPIDiscovery(t *testing.T) {
	text := "Apple Inc. CEO Tim Cook announced on January 15, 2024."
	
	// Test basic functions
	stats := textlib.CalculateTextStatistics(text)
	fmt.Printf("CalculateTextStatistics returns: %T\n", stats)
	
	entities := textlib.ExtractNamedEntities(text)
	fmt.Printf("ExtractNamedEntities returns: %T, len=%d\n", entities, len(entities))
	
	advanced := textlib.ExtractAdvancedEntities(text)
	fmt.Printf("ExtractAdvancedEntities returns: %T, len=%d\n", advanced, len(advanced))
	
	sentences := textlib.SplitIntoSentences(text)
	fmt.Printf("SplitIntoSentences returns: %T, len=%d\n", sentences, len(sentences))
	
	paragraphs := textlib.SplitIntoParagraphs(text)
	fmt.Printf("SplitIntoParagraphs returns: %T, len=%d\n", paragraphs, len(paragraphs))
	
	patterns := textlib.DetectPatterns(text)
	fmt.Printf("DetectPatterns returns: %T\n", patterns)
	
	// Test code functions with simple code
	code := "func main() { fmt.Println(\"hello\") }"
	complexity := textlib.CalculateCyclomaticComplexity(code)
	fmt.Printf("CalculateCyclomaticComplexity returns: %T, value=%d\n", complexity, complexity)
	
	sigs := textlib.ExtractFunctionSignatures(code)
	fmt.Printf("ExtractFunctionSignatures returns: %T, len=%d\n", sigs, len(sigs))
}