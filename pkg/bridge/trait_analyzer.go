package bridge

import (
	"context"
	"strings"
)

// TraitAnalyzer analyzes personality traits and maps them to OCEAN values.
type TraitAnalyzer struct {
	mappings *TraitMappings
}

// NewTraitAnalyzer creates a new trait analyzer with default mappings.
func NewTraitAnalyzer() *TraitAnalyzer {
	return &TraitAnalyzer{
		mappings: GetDefaultMappings(),
	}
}

// NewTraitAnalyzerWithMappings creates a new trait analyzer with custom mappings.
func NewTraitAnalyzerWithMappings(mappings *TraitMappings) *TraitAnalyzer {
	return &TraitAnalyzer{
		mappings: mappings,
	}
}

// AnalyzeTraits converts a list of personality traits to OCEAN values.
func (ta *TraitAnalyzer) AnalyzeTraits(traits []string) PersonalityTraits {
	result := PersonalityTraits{
		Openness:          0.5, // Start with neutral values
		Conscientiousness: 0.5,
		Extraversion:      0.5,
		Agreeableness:     0.5,
		Neuroticism:       0.5,
	}

	if len(traits) == 0 {
		return result
	}

	// Count how many traits affect each dimension
	counts := make(map[string]int)
	adjustments := make(map[string]float64)

	for _, trait := range traits {
		normalizedTrait := ta.normalizeTrait(trait)
		
		// Check direct mappings
		if mapping, exists := ta.mappings.Traits[normalizedTrait]; exists {
			ta.applyMapping(mapping, &adjustments, &counts)
			continue
		}

		// Check for partial matches
		for mappedTrait, mapping := range ta.mappings.Traits {
			if ta.isPartialMatch(normalizedTrait, mappedTrait) {
				ta.applyMapping(mapping, &adjustments, &counts)
				break
			}
		}
	}

	// Apply adjustments with averaging
	if counts["openness"] > 0 {
		result.Openness = 0.5 + (adjustments["openness"] / float64(counts["openness"]))
	}
	if counts["conscientiousness"] > 0 {
		result.Conscientiousness = 0.5 + (adjustments["conscientiousness"] / float64(counts["conscientiousness"]))
	}
	if counts["extraversion"] > 0 {
		result.Extraversion = 0.5 + (adjustments["extraversion"] / float64(counts["extraversion"]))
	}
	if counts["agreeableness"] > 0 {
		result.Agreeableness = 0.5 + (adjustments["agreeableness"] / float64(counts["agreeableness"]))
	}
	if counts["neuroticism"] > 0 {
		result.Neuroticism = 0.5 + (adjustments["neuroticism"] / float64(counts["neuroticism"]))
	}

	// Ensure values are within bounds [0.0, 1.0]
	result = ta.clampValues(result)

	return result
}

// AnalyzeWithContext performs trait analysis with additional context.
func (ta *TraitAnalyzer) AnalyzeWithContext(ctx context.Context, traits []string, behaviors []string, background string) PersonalityTraits {
	// Start with basic trait analysis
	result := ta.AnalyzeTraits(traits)

	// Analyze behaviors for additional insights
	if len(behaviors) > 0 {
		behaviorTraits := ta.extractTraitsFromBehaviors(behaviors)
		behaviorResult := ta.AnalyzeTraits(behaviorTraits)
		
		// Blend with behavior analysis (weighted average)
		result = ta.blendPersonalities(result, behaviorResult, 0.7, 0.3)
	}

	// Analyze background for contextual adjustments
	if background != "" {
		backgroundTraits := ta.extractTraitsFromText(background)
		if len(backgroundTraits) > 0 {
			backgroundResult := ta.AnalyzeTraits(backgroundTraits)
			result = ta.blendPersonalities(result, backgroundResult, 0.8, 0.2)
		}
	}

	return ta.clampValues(result)
}

// normalizeTrait converts a trait to lowercase and removes extra spaces.
func (ta *TraitAnalyzer) normalizeTrait(trait string) string {
	return strings.ToLower(strings.TrimSpace(trait))
}

// isPartialMatch checks if two traits are partial matches.
func (ta *TraitAnalyzer) isPartialMatch(trait1, trait2 string) bool {
	// Simple substring matching for now
	return strings.Contains(trait1, trait2) || strings.Contains(trait2, trait1)
}

// applyMapping applies a trait mapping to the adjustments.
func (ta *TraitAnalyzer) applyMapping(mapping TraitMapping, adjustments *map[string]float64, counts *map[string]int) {
	if mapping.Openness != 0 {
		(*adjustments)["openness"] += mapping.Openness
		(*counts)["openness"]++
	}
	if mapping.Conscientiousness != 0 {
		(*adjustments)["conscientiousness"] += mapping.Conscientiousness
		(*counts)["conscientiousness"]++
	}
	if mapping.Extraversion != 0 {
		(*adjustments)["extraversion"] += mapping.Extraversion
		(*counts)["extraversion"]++
	}
	if mapping.Agreeableness != 0 {
		(*adjustments)["agreeableness"] += mapping.Agreeableness
		(*counts)["agreeableness"]++
	}
	if mapping.Neuroticism != 0 {
		(*adjustments)["neuroticism"] += mapping.Neuroticism
		(*counts)["neuroticism"]++
	}
}

// clampValues ensures all personality values are within [0.0, 1.0].
func (ta *TraitAnalyzer) clampValues(p PersonalityTraits) PersonalityTraits {
	return PersonalityTraits{
		Openness:          ta.clamp(p.Openness),
		Conscientiousness: ta.clamp(p.Conscientiousness),
		Extraversion:      ta.clamp(p.Extraversion),
		Agreeableness:     ta.clamp(p.Agreeableness),
		Neuroticism:       ta.clamp(p.Neuroticism),
	}
}

// clamp restricts a value to [0.0, 1.0].
func (ta *TraitAnalyzer) clamp(value float64) float64 {
	if value < 0.0 {
		return 0.0
	}
	if value > 1.0 {
		return 1.0
	}
	return value
}

// blendPersonalities combines two personality profiles with weights.
func (ta *TraitAnalyzer) blendPersonalities(p1, p2 PersonalityTraits, weight1, weight2 float64) PersonalityTraits {
	return PersonalityTraits{
		Openness:          p1.Openness*weight1 + p2.Openness*weight2,
		Conscientiousness: p1.Conscientiousness*weight1 + p2.Conscientiousness*weight2,
		Extraversion:      p1.Extraversion*weight1 + p2.Extraversion*weight2,
		Agreeableness:     p1.Agreeableness*weight1 + p2.Agreeableness*weight2,
		Neuroticism:       p1.Neuroticism*weight1 + p2.Neuroticism*weight2,
	}
}

// extractTraitsFromBehaviors extracts implicit traits from behavior descriptions.
func (ta *TraitAnalyzer) extractTraitsFromBehaviors(behaviors []string) []string {
	var traits []string
	
	for _, behavior := range behaviors {
		lower := strings.ToLower(behavior)
		
		// Extract traits based on behavior keywords
		if strings.Contains(lower, "careful") || strings.Contains(lower, "meticulous") {
			traits = append(traits, "conscientious")
		}
		if strings.Contains(lower, "social") || strings.Contains(lower, "talkative") {
			traits = append(traits, "extraverted")
		}
		if strings.Contains(lower, "creative") || strings.Contains(lower, "imaginative") {
			traits = append(traits, "open-minded")
		}
		if strings.Contains(lower, "helpful") || strings.Contains(lower, "supportive") {
			traits = append(traits, "agreeable")
		}
		if strings.Contains(lower, "anxious") || strings.Contains(lower, "worried") {
			traits = append(traits, "neurotic")
		}
	}
	
	return traits
}

// extractTraitsFromText extracts personality traits from free-form text.
func (ta *TraitAnalyzer) extractTraitsFromText(text string) []string {
	var traits []string
	lower := strings.ToLower(text)
	
	// Look for trait indicators in the text
	traitKeywords := map[string]string{
		"adventurous":  "adventurous",
		"creative":     "creative",
		"organized":    "organized",
		"friendly":     "friendly",
		"outgoing":     "outgoing",
		"analytical":   "analytical",
		"empathetic":   "empathetic",
		"disciplined":  "disciplined",
		"curious":      "curious",
		"reliable":     "reliable",
	}
	
	for keyword, trait := range traitKeywords {
		if strings.Contains(lower, keyword) {
			traits = append(traits, trait)
		}
	}
	
	return traits
}