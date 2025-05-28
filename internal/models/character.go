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