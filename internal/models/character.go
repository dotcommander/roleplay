package models

import (
	"sync"
	"time"
)

// PersonalityTraits represents OCEAN model traits
type PersonalityTraits struct {
	Openness          float64 `json:"openness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Extraversion      float64 `json:"extraversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Neuroticism       float64 `json:"neuroticism"`
}

// EmotionalState represents current emotional context
type EmotionalState struct {
	Joy      float64 `json:"joy"`
	Surprise float64 `json:"surprise"`
	Anger    float64 `json:"anger"`
	Fear     float64 `json:"fear"`
	Sadness  float64 `json:"sadness"`
	Disgust  float64 `json:"disgust"`
}

// MemoryType represents different types of memories
type MemoryType string

const (
	ShortTermMemory  MemoryType = "short_term"
	MediumTermMemory MemoryType = "medium_term"
	LongTermMemory   MemoryType = "long_term"
)

// Memory represents different memory types
type Memory struct {
	Type      MemoryType `json:"type"`
	Content   string     `json:"content"`
	Timestamp time.Time  `json:"timestamp"`
	Emotional float64    `json:"emotional_weight"`
}

// Character represents a complete character profile
type Character struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Backstory    string            `json:"backstory"`
	Personality  PersonalityTraits `json:"personality"`
	CurrentMood  EmotionalState    `json:"current_mood"`
	Quirks       []string          `json:"quirks"`
	SpeechStyle  string            `json:"speech_style"`
	Memories     []Memory          `json:"memories"`
	LastModified time.Time         `json:"last_modified"`
	
	// Extended fields for richer character definition (OpenAI 1024+ token caching)
	Age              string                 `json:"age,omitempty"`
	Gender           string                 `json:"gender,omitempty"`
	Occupation       string                 `json:"occupation,omitempty"`
	Education        string                 `json:"education,omitempty"`
	Nationality      string                 `json:"nationality,omitempty"`
	Ethnicity        string                 `json:"ethnicity,omitempty"`
	PhysicalTraits   []string               `json:"physical_traits,omitempty"`
	Skills           []string               `json:"skills,omitempty"`
	Interests        []string               `json:"interests,omitempty"`
	Fears            []string               `json:"fears,omitempty"`
	Goals            []string               `json:"goals,omitempty"`
	Relationships    map[string]string      `json:"relationships,omitempty"`
	CoreBeliefs      []string               `json:"core_beliefs,omitempty"`
	MoralCode        []string               `json:"moral_code,omitempty"`
	Flaws            []string               `json:"flaws,omitempty"`
	Strengths        []string               `json:"strengths,omitempty"`
	CatchPhrases     []string               `json:"catch_phrases,omitempty"`
	DialogueExamples []string               `json:"dialogue_examples,omitempty"`
	BehaviorPatterns []string               `json:"behavior_patterns,omitempty"`
	EmotionalTriggers map[string]string     `json:"emotional_triggers,omitempty"`
	DecisionMaking   string                 `json:"decision_making,omitempty"`
	ConflictStyle    string                 `json:"conflict_style,omitempty"`
	WorldView        string                 `json:"world_view,omitempty"`
	LifePhilosophy   string                 `json:"life_philosophy,omitempty"`
	DailyRoutines    []string               `json:"daily_routines,omitempty"`
	Hobbies          []string               `json:"hobbies,omitempty"`
	PetPeeves        []string               `json:"pet_peeves,omitempty"`
	Secrets          []string               `json:"secrets,omitempty"`
	Regrets          []string               `json:"regrets,omitempty"`
	Achievements     []string               `json:"achievements,omitempty"`
	
	mu           sync.RWMutex
}

// Lock acquires write lock
func (c *Character) Lock() {
	c.mu.Lock()
}

// Unlock releases write lock
func (c *Character) Unlock() {
	c.mu.Unlock()
}

// RLock acquires read lock
func (c *Character) RLock() {
	c.mu.RLock()
}

// RUnlock releases read lock
func (c *Character) RUnlock() {
	c.mu.RUnlock()
}

// NormalizePersonality ensures all personality traits are within [0, 1] range
func NormalizePersonality(p PersonalityTraits) PersonalityTraits {
	return PersonalityTraits{
		Openness:          clamp(p.Openness, 0, 1),
		Conscientiousness: clamp(p.Conscientiousness, 0, 1),
		Extraversion:      clamp(p.Extraversion, 0, 1),
		Agreeableness:     clamp(p.Agreeableness, 0, 1),
		Neuroticism:       clamp(p.Neuroticism, 0, 1),
	}
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
