package bridge

import (
	"context"
	"fmt"
	"strings"
	"time"

	roleplayModels "github.com/dotcommander/roleplay/internal/models"
)

// CharactersCharacter represents a character from the Characters system.
// This is a simplified version that contains the fields we need for conversion.
type CharactersCharacter struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Age          int                    `json:"age"`
	Gender       string                 `json:"gender"`
	Archetype    string                 `json:"archetype"`
	Experiences  []string               `json:"experiences"`
	Traits       []string               `json:"traits"`
	Attributes   map[string]interface{} `json:"attributes"`
	Differentials map[string]interface{} `json:"-"` // Not exposed in JSON
	Narrative    string                 `json:"narrative,omitempty"`
	NSFW         bool                   `json:"nsfw"`
	Persona      *CharactersPersona     `json:"persona,omitempty"` // AI behavior configuration
}

// CharactersPersona represents the AI behavior configuration from Characters system.
type CharactersPersona struct {
	CoreIdentity  CharactersCoreIdentity  `json:"coreIdentity"`
	Communication CharactersCommunication `json:"communication"`
	Behavior      CharactersBehavior      `json:"behavior"`
	State         CharactersState         `json:"state"`
	Intimate      *CharactersIntimate     `json:"intimate,omitempty"`
}

// CharactersCoreIdentity represents core identity aspects.
type CharactersCoreIdentity struct {
	Worldview      string `json:"worldview"`
	CoreMotivation string `json:"coreMotivation"`
	CoreFear       string `json:"coreFear"`
	Secret         string `json:"secret"`
}

// CharactersCommunication defines communication style.
type CharactersCommunication struct {
	VoicePacing       string   `json:"voicePacing"`
	VocabularyTier    string   `json:"vocabularyTier"`
	SentenceStructure string   `json:"sentenceStructure"`
	VerbalTics        []string `json:"verbalTics"`
	ForbiddenTopics   []string `json:"forbiddenTopics"`
}

// CharactersBehavior defines behavioral patterns.
type CharactersBehavior struct {
	DecisionHeuristic string   `json:"decisionHeuristic"`
	ConflictStyle     string   `json:"conflictStyle"`
	Quirks            []string `json:"quirks"`
}

// CharactersState represents dynamic state.
type CharactersState struct {
	Mood             string  `json:"mood"`
	StressLevel      float64 `json:"stressLevel"`
	StressThreshold  float64 `json:"stressThreshold"`
	CurrentObjective string  `json:"currentObjective"`
}

// CharactersIntimate represents intimate persona (NSFW).
type CharactersIntimate struct {
	Boundaries CharactersIntimateBoundaries `json:"boundaries"`
}

// CharactersIntimateBoundaries defines boundaries.
type CharactersIntimateBoundaries struct {
	HardLimits []string `json:"hardLimits"`
	SoftLimits []string `json:"softLimits"`
}

// GetAttribute retrieves a value from the attributes using dot notation
func (c *CharactersCharacter) GetAttribute(path string) interface{} {
	keys := strings.Split(path, ".")
	var current interface{} = c.Attributes
	
	for _, key := range keys {
		// If current isn't a map, we can't go deeper
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		
		// Get the value at this level
		value, exists := currentMap[key]
		if !exists {
			return nil
		}
		
		// Move to the next level
		current = value
	}
	
	return current
}

// CharactersConverter converts between the Characters format and UniversalCharacter format.
type CharactersConverter struct {
	*BaseConverter
	analyzer *TraitAnalyzer
}

// NewCharactersConverter creates a new converter for the Characters format.
func NewCharactersConverter() *CharactersConverter {
	return &CharactersConverter{
		BaseConverter: NewBaseConverter("characters"),
		analyzer:      NewTraitAnalyzer(),
	}
}

// CanConvert checks if the data is in Characters format.
func (c *CharactersConverter) CanConvert(data interface{}) bool {
	switch v := data.(type) {
	case *CharactersCharacter:
		return true
	case CharactersCharacter:
		return true
	case map[string]interface{}:
		// Check for Characters-specific fields
		_, hasTraits := v["traits"]
		_, hasAttributes := v["attributes"]
		_, hasDifferentials := v["differentials"]
		_, hasArchetype := v["archetype"]
		// More flexible checking - just need some character-like fields
		return hasTraits || hasAttributes || hasDifferentials || hasArchetype
	default:
		return false
	}
}

// ToUniversal converts from Characters format to UniversalCharacter.
func (c *CharactersConverter) ToUniversal(ctx context.Context, data interface{}) (*UniversalCharacter, error) {
	var char *CharactersCharacter
	
	switch v := data.(type) {
	case *CharactersCharacter:
		char = v
	case CharactersCharacter:
		char = &v
	case map[string]interface{}:
		// Convert map to Character struct
		converted, err := c.mapToCharacter(v)
		if err != nil {
			return nil, &ConversionError{
				Source: "characters",
				Target: "universal",
				Err:    err,
			}
		}
		char = converted
	default:
		return nil, &ConversionError{
			Source: "characters",
			Target: "universal",
			Err:    fmt.Errorf("unsupported data type: %T", data),
		}
	}

	// Create UniversalCharacter
	uc := &UniversalCharacter{
		ID:          char.ID,
		Name:        char.Name,
		Description: c.buildDescription(char),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Source:      "characters",
		Version:     "1.0",
	}

	// Extract traits and behaviors
	uc.Traits = char.Traits
	uc.Behaviors = c.extractBehaviors(char)
	
	// Analyze traits to determine OCEAN personality
	allTraits := append([]string{}, char.Traits...)
	if char.Persona != nil {
		allTraits = append(allTraits, c.extractPersonaTraits(char.Persona)...)
	}
	
	uc.Personality = c.analyzer.AnalyzeWithContext(ctx, allTraits, uc.Behaviors, char.Narrative)

	// Extract background/backstory
	if char.Narrative != "" {
		uc.Background = char.Narrative
	} else {
		uc.Background = c.generateBackstory(char)
	}

	// Extract speech style and quirks
	uc.Quirks = c.extractAllQuirks(char)
	uc.Catchphrases = c.extractCatchphrases(char)
	if char.Persona != nil {
		uc.SpeechStyle = c.extractSpeechStyle(char.Persona)
	} else {
		uc.SpeechStyle = c.extractSpeechStyleFromAttributes(char)
	}

	// Extract attributes as metadata
	uc.Topics = c.extractTopics(char)
	uc.Motivations = c.extractMotivations(char)
	uc.Fears = c.extractFears(char)
	uc.Relationships = c.extractRelationships(char)

	// Set boundaries from persona
	if char.Persona != nil {
		uc.Boundaries = char.Persona.Communication.ForbiddenTopics
		if char.Persona.Intimate != nil {
			uc.Boundaries = append(uc.Boundaries, char.Persona.Intimate.Boundaries.HardLimits...)
		}
	}

	// Generate system prompt
	uc.SystemPrompt = c.generateSystemPrompt(char)

	// Add tags
	uc.Tags = c.generateTags(char)

	// Store original data
	uc.SourceData = map[string]interface{}{
		"archetype":     char.Archetype,
		"experiences":   char.Experiences,
		"attributes":    char.Attributes,
		"differentials": char.Differentials,
		"nsfw":          char.NSFW,
		"age":           char.Age,
		"gender":        char.Gender,
	}
	
	if char.Persona != nil {
		uc.SourceData["persona"] = char.Persona
	}

	return uc, nil
}

// FromUniversal converts from UniversalCharacter to Characters format (roleplay.Character).
func (c *CharactersConverter) FromUniversal(ctx context.Context, uc *UniversalCharacter) (interface{}, error) {
	// Create roleplay Character
	char := &roleplayModels.Character{
		ID:          uc.ID,
		Name:        uc.Name,
		Backstory:   uc.Background,
		Personality: c.convertPersonality(uc.Personality),
		SpeechStyle: uc.SpeechStyle,
		Quirks:      uc.Quirks,
		CatchPhrases: uc.Catchphrases,
		LastModified: uc.UpdatedAt,
	}

	// Set current mood to neutral
	char.CurrentMood = roleplayModels.EmotionalState{
		Joy:      0.5,
		Surprise: 0.5,
		Anger:    0.0,
		Fear:     0.0,
		Sadness:  0.0,
		Disgust:  0.0,
	}

	// Extract additional fields from source data if available
	if uc.SourceData != nil {
		if age, ok := uc.SourceData["age"].(int); ok {
			char.Age = fmt.Sprintf("%d", age)
		} else if age, ok := uc.SourceData["age"].(float64); ok {
			char.Age = fmt.Sprintf("%d", int(age))
		}
		
		if gender, ok := uc.SourceData["gender"].(string); ok {
			char.Gender = gender
		}
		
		if archetype, ok := uc.SourceData["archetype"].(string); ok {
			char.Occupation = archetype // Use archetype as occupation
		}
		
		if experiences, ok := uc.SourceData["experiences"].([]string); ok {
			char.BehaviorPatterns = experiences
		} else if experiences, ok := uc.SourceData["experiences"].([]interface{}); ok {
			for _, exp := range experiences {
				if expStr, ok := exp.(string); ok {
					char.BehaviorPatterns = append(char.BehaviorPatterns, expStr)
				}
			}
		}
	}

	// Map universal fields to roleplay extended fields
	char.Skills = c.extractSkillsFromTopics(uc.Topics)
	char.Interests = c.extractInterestsFromTopics(uc.Topics)
	char.Fears = uc.Fears
	char.Goals = uc.Motivations
	char.Relationships = uc.Relationships
	
	// Extract beliefs and moral code from traits
	char.CoreBeliefs, char.MoralCode = c.extractBeliefsAndMorals(uc.Traits)
	
	// Extract flaws and strengths
	char.Flaws, char.Strengths = c.extractFlawsAndStrengths(uc.Traits)
	
	// Set dialogue examples if available
	if uc.Examples != nil {
		for _, example := range uc.Examples {
			char.DialogueExamples = append(char.DialogueExamples, 
				fmt.Sprintf("User: %s\n%s: %s", example.User, uc.Name, example.Character))
		}
	}

	// Map behaviors
	char.BehaviorPatterns = append(char.BehaviorPatterns, uc.Behaviors...)
	
	// Set emotional triggers based on persona data
	if persona, ok := uc.SourceData["persona"].(*CharactersPersona); ok && persona != nil {
		char.EmotionalTriggers = c.extractEmotionalTriggers(persona)
		char.DecisionMaking = persona.Behavior.DecisionHeuristic
		char.ConflictStyle = persona.Behavior.ConflictStyle
		char.WorldView = persona.CoreIdentity.Worldview
	}

	// Initialize empty memories
	char.Memories = []roleplayModels.Memory{}

	return char, nil
}

// Helper methods

func (c *CharactersConverter) mapToCharacter(data map[string]interface{}) (*CharactersCharacter, error) {
	char := &CharactersCharacter{
		Attributes:    make(map[string]interface{}),
		Differentials: make(map[string]interface{}),
	}

	// Map basic fields
	if id, ok := data["id"].(string); ok {
		char.ID = id
	}
	if name, ok := data["name"].(string); ok {
		char.Name = name
	}
	if age, ok := data["age"].(int); ok {
		char.Age = age
	} else if age, ok := data["age"].(float64); ok {
		char.Age = int(age)
	}
	if gender, ok := data["gender"].(string); ok {
		char.Gender = gender
	}
	if archetype, ok := data["archetype"].(string); ok {
		char.Archetype = archetype
	}
	if narrative, ok := data["narrative"].(string); ok {
		char.Narrative = narrative
	}
	// Map backstory field
	if backstory, ok := data["backstory"].(string); ok {
		char.Narrative = backstory
	}
	if nsfw, ok := data["nsfw"].(bool); ok {
		char.NSFW = nsfw
	}

	// Map arrays
	if traits, ok := data["traits"].([]interface{}); ok {
		for _, t := range traits {
			if trait, ok := t.(string); ok {
				char.Traits = append(char.Traits, trait)
			}
		}
	}
	if experiences, ok := data["experiences"].([]interface{}); ok {
		for _, e := range experiences {
			if exp, ok := e.(string); ok {
				char.Experiences = append(char.Experiences, exp)
			}
		}
	}

	// Map complex fields
	if attrs, ok := data["attributes"].(map[string]interface{}); ok {
		char.Attributes = attrs
		// Extract age from attributes.physical.age if available
		if physical, ok := attrs["physical"].(map[string]interface{}); ok {
			if age, ok := physical["age"].(int); ok {
				char.Age = age
			} else if age, ok := physical["age"].(float64); ok {
				char.Age = int(age)
			}
		}
	}
	if diffs, ok := data["differentials"].(map[string]interface{}); ok {
		char.Differentials = diffs
	}

	// Map persona if present
	if personaData, ok := data["persona"].(map[string]interface{}); ok {
		persona := &CharactersPersona{}
		
		// Handle flat persona structure (test-character.json format)
		if worldview, ok := personaData["worldview"].(string); ok {
			persona.CoreIdentity.Worldview = worldview
		}
		if voicePacing, ok := personaData["voice_pacing"].(string); ok {
			persona.Communication.VoicePacing = voicePacing
		}
		
		// Map catchphrases to verbalTics
		if catchphrases, ok := personaData["catchphrases"].([]interface{}); ok {
			for _, cp := range catchphrases {
				if phrase, ok := cp.(string); ok {
					persona.Communication.VerbalTics = append(persona.Communication.VerbalTics, phrase)
				}
			}
		}
		
		// Map quirks to behavior
		if quirks, ok := personaData["quirks"].([]interface{}); ok {
			for _, q := range quirks {
				if quirk, ok := q.(string); ok {
					persona.Behavior.Quirks = append(persona.Behavior.Quirks, quirk)
				}
			}
		}
		
		// Map forbidden_topics to forbiddenTopics
		if topics, ok := personaData["forbidden_topics"].([]interface{}); ok {
			for _, t := range topics {
				if topic, ok := t.(string); ok {
					persona.Communication.ForbiddenTopics = append(persona.Communication.ForbiddenTopics, topic)
				}
			}
		}
		
		// Map CoreIdentity (nested structure)
		if coreData, ok := personaData["coreIdentity"].(map[string]interface{}); ok {
			persona.CoreIdentity.Worldview, _ = coreData["worldview"].(string)
			persona.CoreIdentity.CoreMotivation, _ = coreData["coreMotivation"].(string)
			persona.CoreIdentity.CoreFear, _ = coreData["coreFear"].(string)
			persona.CoreIdentity.Secret, _ = coreData["secret"].(string)
		}
		
		// Map Communication (nested structure)
		if commData, ok := personaData["communication"].(map[string]interface{}); ok {
			persona.Communication.VoicePacing, _ = commData["voicePacing"].(string)
			persona.Communication.VocabularyTier, _ = commData["vocabularyTier"].(string)
			persona.Communication.SentenceStructure, _ = commData["sentenceStructure"].(string)
			
			if tics, ok := commData["verbalTics"].([]interface{}); ok {
				for _, t := range tics {
					if tic, ok := t.(string); ok {
						persona.Communication.VerbalTics = append(persona.Communication.VerbalTics, tic)
					}
				}
			}
			
			if topics, ok := commData["forbiddenTopics"].([]interface{}); ok {
				for _, t := range topics {
					if topic, ok := t.(string); ok {
						persona.Communication.ForbiddenTopics = append(persona.Communication.ForbiddenTopics, topic)
					}
				}
			}
		}
		
		// Map Behavior (nested structure)
		if behaviorData, ok := personaData["behavior"].(map[string]interface{}); ok {
			persona.Behavior.DecisionHeuristic, _ = behaviorData["decisionHeuristic"].(string)
			persona.Behavior.ConflictStyle, _ = behaviorData["conflictStyle"].(string)
			
			if quirks, ok := behaviorData["quirks"].([]interface{}); ok {
				for _, q := range quirks {
					if quirk, ok := q.(string); ok {
						persona.Behavior.Quirks = append(persona.Behavior.Quirks, quirk)
					}
				}
			}
		}
		
		// Map State (nested structure)
		if stateData, ok := personaData["state"].(map[string]interface{}); ok {
			persona.State.Mood, _ = stateData["mood"].(string)
			if stress, ok := stateData["stressLevel"].(float64); ok {
				persona.State.StressLevel = stress
			}
			if threshold, ok := stateData["stressThreshold"].(float64); ok {
				persona.State.StressThreshold = threshold
			}
			persona.State.CurrentObjective, _ = stateData["currentObjective"].(string)
		}
		
		char.Persona = persona
	}

	return char, nil
}

func (c *CharactersConverter) buildDescription(char *CharactersCharacter) string {
	parts := []string{}
	
	if char.Age > 0 {
		parts = append(parts, fmt.Sprintf("%d year old", char.Age))
	}
	if char.Gender != "" {
		parts = append(parts, char.Gender)
	}
	if char.Archetype != "" {
		parts = append(parts, char.Archetype)
	}
	
	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}
	return ""
}

func (c *CharactersConverter) extractBehaviors(char *CharactersCharacter) []string {
	behaviors := []string{}
	
	// Extract from experiences
	behaviors = append(behaviors, char.Experiences...)
	
	// Extract from persona quirks
	if char.Persona != nil {
		behaviors = append(behaviors, char.Persona.Behavior.Quirks...)
	}
	
	// Extract behavioral attributes
	if attrs, ok := char.Attributes["behaviors"].([]interface{}); ok {
		for _, b := range attrs {
			if behavior, ok := b.(string); ok {
				behaviors = append(behaviors, behavior)
			}
		}
	}
	
	return behaviors
}

func (c *CharactersConverter) extractPersonaTraits(persona *CharactersPersona) []string {
	traits := []string{}
	
	// Extract traits from core identity
	if persona.CoreIdentity.Worldview != "" {
		traits = append(traits, c.worldviewToTraits(persona.CoreIdentity.Worldview)...)
	}
	
	// Extract from communication style
	if persona.Communication.VoicePacing != "" {
		traits = append(traits, c.voicePacingToTraits(persona.Communication.VoicePacing)...)
	}
	
	// Extract from behavior
	if persona.Behavior.ConflictStyle != "" {
		traits = append(traits, c.conflictStyleToTraits(persona.Behavior.ConflictStyle)...)
	}
	
	return traits
}

func (c *CharactersConverter) extractSpeechStyle(persona *CharactersPersona) string {
	parts := []string{}
	
	if persona.Communication.VoicePacing != "" {
		// Check if it already contains descriptive text
		if strings.Contains(strings.ToLower(persona.Communication.VoicePacing), "sentence") {
			parts = append(parts, persona.Communication.VoicePacing)
		} else {
			parts = append(parts, persona.Communication.VoicePacing+" pacing")
		}
	}
	if persona.Communication.VocabularyTier != "" {
		parts = append(parts, persona.Communication.VocabularyTier+" vocabulary")
	}
	if persona.Communication.SentenceStructure != "" {
		parts = append(parts, persona.Communication.SentenceStructure+" sentences")
	}
	
	return strings.Join(parts, ", ")
}

func (c *CharactersConverter) generateBackstory(char *CharactersCharacter) string {
	parts := []string{}
	
	// Add archetype-based backstory
	if char.Archetype != "" {
		parts = append(parts, fmt.Sprintf("A %s by nature", char.Archetype))
	}
	
	// Add experiences
	if len(char.Experiences) > 0 {
		parts = append(parts, "with experience in "+strings.Join(char.Experiences, ", "))
	}
	
	// Add key attributes
	if origin, ok := char.GetAttribute("origin").(string); ok {
		parts = append(parts, fmt.Sprintf("Originally from %s", origin))
	}
	
	if len(parts) > 0 {
		return strings.Join(parts, ". ") + "."
	}
	
	return ""
}

func (c *CharactersConverter) extractTopics(char *CharactersCharacter) []string {
	topics := []string{}
	
	// Add archetype as a topic
	if char.Archetype != "" {
		topics = append(topics, char.Archetype)
	}
	
	// Add experiences as topics
	topics = append(topics, char.Experiences...)
	
	// Extract from attributes
	if skills, ok := char.GetAttribute("skills").([]interface{}); ok {
		for _, s := range skills {
			if skill, ok := s.(string); ok {
				topics = append(topics, skill)
			}
		}
	}
	
	return topics
}

func (c *CharactersConverter) extractMotivations(char *CharactersCharacter) []string {
	motivations := []string{}
	
	if char.Persona != nil && char.Persona.CoreIdentity.CoreMotivation != "" {
		motivations = append(motivations, char.Persona.CoreIdentity.CoreMotivation)
	}
	
	// Extract from attributes
	if goals, ok := char.GetAttribute("goals").([]interface{}); ok {
		for _, g := range goals {
			if goal, ok := g.(string); ok {
				motivations = append(motivations, goal)
			}
		}
	}
	
	return motivations
}

func (c *CharactersConverter) extractFears(char *CharactersCharacter) []string {
	fears := []string{}
	
	if char.Persona != nil && char.Persona.CoreIdentity.CoreFear != "" {
		fears = append(fears, char.Persona.CoreIdentity.CoreFear)
	}
	
	// Extract from attributes
	if attrFears, ok := char.GetAttribute("fears").([]interface{}); ok {
		for _, f := range attrFears {
			if fear, ok := f.(string); ok {
				fears = append(fears, fear)
			}
		}
	}
	
	return fears
}

func (c *CharactersConverter) extractRelationships(char *CharactersCharacter) map[string]string {
	relationships := make(map[string]string)
	
	// Extract from attributes
	if rels, ok := char.GetAttribute("relationships").(map[string]interface{}); ok {
		for name, desc := range rels {
			if descStr, ok := desc.(string); ok {
				relationships[name] = descStr
			}
		}
	}
	
	return relationships
}

func (c *CharactersConverter) generateSystemPrompt(char *CharactersCharacter) string {
	prompt := fmt.Sprintf("You are %s", char.Name)
	
	if desc := c.buildDescription(char); desc != "" {
		prompt += ", " + desc
	}
	
	if char.Narrative != "" {
		prompt += ". " + char.Narrative
	}
	
	if char.Persona != nil {
		if char.Persona.CoreIdentity.Worldview != "" {
			prompt += fmt.Sprintf("\n\nWorldview: %s", char.Persona.CoreIdentity.Worldview)
		}
		if char.Persona.CoreIdentity.CoreMotivation != "" {
			prompt += fmt.Sprintf("\nCore motivation: %s", char.Persona.CoreIdentity.CoreMotivation)
		}
	}
	
	return prompt
}

func (c *CharactersConverter) generateTags(char *CharactersCharacter) []string {
	tags := []string{}
	
	if char.Archetype != "" {
		tags = append(tags, char.Archetype)
	}
	if char.Gender != "" {
		tags = append(tags, char.Gender)
	}
	if char.NSFW {
		tags = append(tags, "nsfw")
	}
	
	return tags
}

func (c *CharactersConverter) convertPersonality(p PersonalityTraits) roleplayModels.PersonalityTraits {
	return roleplayModels.PersonalityTraits{
		Openness:          p.Openness,
		Conscientiousness: p.Conscientiousness,
		Extraversion:      p.Extraversion,
		Agreeableness:     p.Agreeableness,
		Neuroticism:       p.Neuroticism,
	}
}

func (c *CharactersConverter) extractSkillsFromTopics(topics []string) []string {
	// Filter topics that represent skills
	skills := []string{}
	for _, topic := range topics {
		if c.isSkill(topic) {
			skills = append(skills, topic)
		}
	}
	return skills
}

func (c *CharactersConverter) extractInterestsFromTopics(topics []string) []string {
	// Filter topics that represent interests
	interests := []string{}
	for _, topic := range topics {
		if !c.isSkill(topic) {
			interests = append(interests, topic)
		}
	}
	return interests
}

func (c *CharactersConverter) isSkill(topic string) bool {
	// Simple heuristic to determine if a topic is a skill
	skillKeywords := []string{"combat", "magic", "craft", "tech", "medical", "social", "survival"}
	lower := strings.ToLower(topic)
	for _, keyword := range skillKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

func (c *CharactersConverter) extractBeliefsAndMorals(traits []string) ([]string, []string) {
	beliefs := []string{}
	morals := []string{}
	
	for _, trait := range traits {
		lower := strings.ToLower(trait)
		if strings.Contains(lower, "believe") || strings.Contains(lower, "faith") {
			beliefs = append(beliefs, trait)
		} else if strings.Contains(lower, "honor") || strings.Contains(lower, "moral") || 
			strings.Contains(lower, "ethic") || strings.Contains(lower, "principle") {
			morals = append(morals, trait)
		}
	}
	
	return beliefs, morals
}

func (c *CharactersConverter) extractFlawsAndStrengths(traits []string) ([]string, []string) {
	flaws := []string{}
	strengths := []string{}
	
	// Negative trait keywords
	negativeKeywords := []string{"arrogant", "stubborn", "reckless", "impulsive", "lazy", 
		"coward", "selfish", "dishonest", "cruel", "weak"}
	
	// Positive trait keywords
	positiveKeywords := []string{"brave", "loyal", "honest", "kind", "strong", 
		"intelligent", "wise", "creative", "determined", "compassionate"}
	
	for _, trait := range traits {
		lower := strings.ToLower(trait)
		
		isNegative := false
		for _, keyword := range negativeKeywords {
			if strings.Contains(lower, keyword) {
				flaws = append(flaws, trait)
				isNegative = true
				break
			}
		}
		
		if !isNegative {
			for _, keyword := range positiveKeywords {
				if strings.Contains(lower, keyword) {
					strengths = append(strengths, trait)
					break
				}
			}
		}
	}
	
	return flaws, strengths
}

func (c *CharactersConverter) extractEmotionalTriggers(persona *CharactersPersona) map[string]string {
	triggers := make(map[string]string)
	
	// Map stress level to trigger
	if persona.State.StressLevel > persona.State.StressThreshold {
		triggers["high_stress"] = "becomes agitated and reactive"
	}
	
	// Map forbidden topics to triggers
	for _, topic := range persona.Communication.ForbiddenTopics {
		triggers[topic] = "becomes uncomfortable and defensive"
	}
	
	// Map core fear to trigger
	if persona.CoreIdentity.CoreFear != "" {
		triggers[persona.CoreIdentity.CoreFear] = "experiences deep anxiety"
	}
	
	return triggers
}

// Trait mapping helpers

func (c *CharactersConverter) worldviewToTraits(worldview string) []string {
	lower := strings.ToLower(worldview)
	traits := []string{}
	
	if strings.Contains(lower, "optimist") {
		traits = append(traits, "optimistic", "hopeful")
	} else if strings.Contains(lower, "pessimist") {
		traits = append(traits, "pessimistic", "cynical")
	}
	
	if strings.Contains(lower, "pragmat") {
		traits = append(traits, "pragmatic", "practical")
	} else if strings.Contains(lower, "idealist") {
		traits = append(traits, "idealistic", "visionary")
	}
	
	return traits
}

func (c *CharactersConverter) voicePacingToTraits(pacing string) []string {
	lower := strings.ToLower(pacing)
	traits := []string{}
	
	if strings.Contains(lower, "fast") || strings.Contains(lower, "quick") {
		traits = append(traits, "energetic", "impulsive")
	} else if strings.Contains(lower, "slow") || strings.Contains(lower, "deliberate") {
		traits = append(traits, "thoughtful", "patient")
	} else if strings.Contains(lower, "short") || strings.Contains(lower, "clipped") {
		traits = append(traits, "direct", "efficient", "serious")
	} else if strings.Contains(lower, "direct") {
		traits = append(traits, "straightforward", "honest")
	}
	
	return traits
}

func (c *CharactersConverter) conflictStyleToTraits(style string) []string {
	lower := strings.ToLower(style)
	traits := []string{}
	
	if strings.Contains(lower, "avoid") {
		traits = append(traits, "conflict-avoidant", "peaceful")
	} else if strings.Contains(lower, "confront") || strings.Contains(lower, "aggressive") {
		traits = append(traits, "confrontational", "assertive")
	} else if strings.Contains(lower, "collaborat") {
		traits = append(traits, "collaborative", "diplomatic")
	}
	
	return traits
}

// extractAllQuirks extracts quirks from both attributes and persona
func (c *CharactersConverter) extractAllQuirks(char *CharactersCharacter) []string {
	quirks := []string{}
	
	// Extract from attributes.personality.quirks
	if personality, ok := char.GetAttribute("personality.quirks").([]interface{}); ok {
		for _, q := range personality {
			if quirk, ok := q.(string); ok {
				quirks = append(quirks, quirk)
			}
		}
	}
	
	// Extract from persona.quirks
	if char.Persona != nil {
		quirks = append(quirks, char.Persona.Behavior.Quirks...)
	}
	
	return quirks
}

// extractCatchphrases extracts catchphrases from persona
func (c *CharactersConverter) extractCatchphrases(char *CharactersCharacter) []string {
	catchphrases := []string{}
	
	if char.Persona != nil {
		catchphrases = append(catchphrases, char.Persona.Communication.VerbalTics...)
	}
	
	return catchphrases
}

// extractSpeechStyleFromAttributes extracts speech style from attributes when persona is not available
func (c *CharactersConverter) extractSpeechStyleFromAttributes(char *CharactersCharacter) string {
	parts := []string{}
	
	// Try to extract from various attribute paths
	if style, ok := char.GetAttribute("speech.style").(string); ok {
		parts = append(parts, style)
	}
	if pacing, ok := char.GetAttribute("speech.pacing").(string); ok {
		parts = append(parts, pacing+" pacing")
	}
	
	return strings.Join(parts, ", ")
}