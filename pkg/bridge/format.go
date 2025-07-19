// Package bridge provides a universal character format and converters
// for translating between different character definition systems.
package bridge

import (
	"time"
)

// UniversalCharacter represents a character definition that can be
// converted between different AI systems while preserving essential
// personality traits and characteristics.
type UniversalCharacter struct {
	// Basic Information
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Core Personality (OCEAN model)
	Personality PersonalityTraits `json:"personality"`

	// Character Definition
	Background   string   `json:"background,omitempty"`
	Traits       []string `json:"traits,omitempty"`       // Raw personality traits
	Behaviors    []string `json:"behaviors,omitempty"`    // Behavioral patterns
	SpeechStyle  string   `json:"speech_style,omitempty"`  // How they speak
	Motivations  []string `json:"motivations,omitempty"`  // What drives them
	Fears        []string `json:"fears,omitempty"`        // What they fear
	Relationships map[string]string `json:"relationships,omitempty"` // Key relationships

	// Interaction Guidelines
	Topics      []string `json:"topics,omitempty"`       // Topics they know about
	Boundaries  []string `json:"boundaries,omitempty"`   // Things they won't do
	Quirks      []string `json:"quirks,omitempty"`       // Unique behaviors
	Catchphrases []string `json:"catchphrases,omitempty"` // Common phrases

	// System Instructions
	SystemPrompt string `json:"system_prompt,omitempty"` // Full system prompt if available
	Examples     []ConversationExample `json:"examples,omitempty"` // Example interactions

	// Metadata
	Source      string            `json:"source"`      // Original system (e.g., "character.ai", "manual")
	SourceData  map[string]interface{} `json:"source_data,omitempty"` // Original format data
	Tags        []string          `json:"tags,omitempty"`
	Version     string            `json:"version,omitempty"`
}

// PersonalityTraits represents the OCEAN (Big Five) personality model.
type PersonalityTraits struct {
	// Openness to Experience (0.0-1.0)
	// High: Imaginative, curious, open to new ideas
	// Low: Practical, conventional, prefers routine
	Openness float64 `json:"openness"`

	// Conscientiousness (0.0-1.0)
	// High: Organized, responsible, dependable
	// Low: Impulsive, careless, disorganized
	Conscientiousness float64 `json:"conscientiousness"`

	// Extraversion (0.0-1.0)
	// High: Outgoing, energetic, seeks stimulation
	// Low: Reserved, solitary, thoughtful
	Extraversion float64 `json:"extraversion"`

	// Agreeableness (0.0-1.0)
	// High: Cooperative, trusting, helpful
	// Low: Competitive, skeptical, challenging
	Agreeableness float64 `json:"agreeableness"`

	// Neuroticism (0.0-1.0)
	// High: Emotionally reactive, stress-prone
	// Low: Calm, resilient, secure
	Neuroticism float64 `json:"neuroticism"`
}

// ConversationExample represents a sample interaction.
type ConversationExample struct {
	User      string `json:"user"`
	Character string `json:"character"`
	Context   string `json:"context,omitempty"`
}

// Validate checks if the UniversalCharacter has required fields.
func (uc *UniversalCharacter) Validate() error {
	// TODO: Implement validation logic
	return nil
}

// SetDefaults sets default values for empty fields.
func (uc *UniversalCharacter) SetDefaults() {
	if uc.CreatedAt.IsZero() {
		uc.CreatedAt = time.Now()
	}
	if uc.UpdatedAt.IsZero() {
		uc.UpdatedAt = time.Now()
	}
	if uc.Version == "" {
		uc.Version = "1.0"
	}
}