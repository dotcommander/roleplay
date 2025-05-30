package providers

import (
	"encoding/json"
	"fmt"
	"os"
)

// DebugResponse prints the raw API response if DEBUG_RESPONSE env var is set
func DebugResponse(response interface{}) {
	if os.Getenv("DEBUG_RESPONSE") != "true" {
		return
	}
	
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Failed to marshal response: %v\n", err)
		return
	}
	
	fmt.Fprintf(os.Stderr, "\n🔍 DEBUG - Raw API Response:\n%s\n\n", string(data))
}

// DebugCachedTokens prints cached token information if available
func DebugCachedTokens(promptTokens, cachedTokens int) {
	if os.Getenv("DEBUG_CACHE") != "true" {
		return
	}
	
	if cachedTokens > 0 {
		percentage := float64(cachedTokens) / float64(promptTokens) * 100
		fmt.Fprintf(os.Stderr, "🔍 DEBUG - Prompt Cache: %d/%d tokens cached (%.1f%%)\n", 
			cachedTokens, promptTokens, percentage)
	} else {
		fmt.Fprintf(os.Stderr, "🔍 DEBUG - Prompt Cache: No tokens cached (0/%d)\n", promptTokens)
	}
}