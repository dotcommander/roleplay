package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ExtractValidJSON attempts to find and extract a valid JSON object from a raw string.
// It's designed to handle common LLM output issues like:
// - Extraneous text before/after JSON
// - Markdown code blocks
// - Incomplete JSON due to token limits
// - Multiple JSON objects (returns the first complete one)
func ExtractValidJSON(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	// Remove common markdown code blocks first
	if strings.HasPrefix(raw, "```json") && strings.HasSuffix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	} else if strings.HasPrefix(raw, "```") && strings.HasSuffix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}

	// Find the first '{' which indicates the start of a JSON object
	startIndex := strings.Index(raw, "{")
	if startIndex == -1 {
		return "", fmt.Errorf("no JSON object start '{' found in response")
	}

	// Try to find a matching closing brace by counting braces
	braceCount := 0
	inString := false
	escapeNext := false
	endIndex := -1

	for i := startIndex; i < len(raw); i++ {
		char := raw[i]

		// Handle escape sequences in strings
		if escapeNext {
			escapeNext = false
			continue
		}

		if char == '\\' && inString {
			escapeNext = true
			continue
		}

		// Toggle string state when we encounter quotes
		if char == '"' && !escapeNext {
			inString = !inString
			continue
		}

		// Only count braces outside of strings
		if !inString {
			switch char {
			case '{':
				braceCount++
			case '}':
				braceCount--
				if braceCount == 0 {
					endIndex = i
					break
				}
			}
		}

		if endIndex != -1 {
			break
		}
	}

	// If we didn't find a matching closing brace, try to be more lenient
	if endIndex == -1 {
		// Look for the last '}' in the string as a fallback
		lastBrace := strings.LastIndex(raw, "}")
		if lastBrace > startIndex {
			endIndex = lastBrace
		} else {
			return "", fmt.Errorf("no matching JSON object end '}' found (incomplete JSON)")
		}
	}

	potentialJSON := raw[startIndex : endIndex+1]

	// Validate that the extracted substring is valid JSON
	var jsonData interface{}
	decoder := json.NewDecoder(strings.NewReader(potentialJSON))
	decoder.UseNumber() // Preserve number precision
	
	if err := decoder.Decode(&jsonData); err != nil {
		// If it fails, it might be incomplete. Try to repair common issues.
		repairedJSON := attemptJSONRepair(potentialJSON)
		if repairedJSON != potentialJSON {
			// Try parsing the repaired version
			decoder = json.NewDecoder(strings.NewReader(repairedJSON))
			decoder.UseNumber()
			if err2 := decoder.Decode(&jsonData); err2 == nil {
				return repairedJSON, nil
			}
		}
		
		return "", fmt.Errorf("extracted substring is not valid JSON: %w. Substring: %s", err, potentialJSON)
	}

	return potentialJSON, nil
}

// attemptJSONRepair tries to fix common JSON truncation issues
func attemptJSONRepair(jsonStr string) string {
	// Count open brackets/braces
	openBraces := strings.Count(jsonStr, "{")
	closeBraces := strings.Count(jsonStr, "}")
	openBrackets := strings.Count(jsonStr, "[")
	closeBrackets := strings.Count(jsonStr, "]")

	repaired := jsonStr

	// Add missing closing brackets/braces
	for i := 0; i < openBrackets-closeBrackets; i++ {
		repaired += "]"
	}
	for i := 0; i < openBraces-closeBraces; i++ {
		repaired += "}"
	}

	// Check if the JSON ends mid-string (common truncation point)
	// Look for an unclosed quoted string at the end
	if !strings.HasSuffix(strings.TrimSpace(repaired), "}") && !strings.HasSuffix(strings.TrimSpace(repaired), "]") {
		// Count quotes to see if we're in an unclosed string
		quoteCount := 0
		inEscape := false
		for _, char := range repaired {
			if inEscape {
				inEscape = false
				continue
			}
			if char == '\\' {
				inEscape = true
				continue
			}
			if char == '"' {
				quoteCount++
			}
		}

		if quoteCount%2 == 1 {
			// Odd number of quotes means unclosed string
			// Close the string and any necessary JSON structure
			repaired += "\""
			
			// Try to intelligently close the structure
			// This is a heuristic - look at what came before the string
			lastComma := strings.LastIndex(repaired, ",")
			lastOpenBrace := strings.LastIndex(repaired, "{")
			lastOpenBracket := strings.LastIndex(repaired, "[")
			
			if lastComma > lastOpenBrace && lastComma > lastOpenBracket {
				// We're likely in the middle of an object or array
				if lastOpenBracket > lastOpenBrace {
					repaired += "]"
				}
				repaired += "}"
			}
		}
	}

	return repaired
}