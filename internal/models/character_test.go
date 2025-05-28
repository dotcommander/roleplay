package models

import (
	"testing"
	"time"
)

func TestNormalizePersonality(t *testing.T) {
	tests := []struct {
		name     string
		input    PersonalityTraits
		expected PersonalityTraits
	}{
		{
			name: "values within range",
			input: PersonalityTraits{
				Openness:          0.5,
				Conscientiousness: 0.6,
				Extraversion:      0.7,
				Agreeableness:     0.8,
				Neuroticism:       0.9,
			},
			expected: PersonalityTraits{
				Openness:          0.5,
				Conscientiousness: 0.6,
				Extraversion:      0.7,
				Agreeableness:     0.8,
				Neuroticism:       0.9,
			},
		},
		{
			name: "values above 1 should be clamped",
			input: PersonalityTraits{
				Openness:          1.5,
				Conscientiousness: 2.0,
				Extraversion:      1.1,
				Agreeableness:     1.3,
				Neuroticism:       1.8,
			},
			expected: PersonalityTraits{
				Openness:          1.0,
				Conscientiousness: 1.0,
				Extraversion:      1.0,
				Agreeableness:     1.0,
				Neuroticism:       1.0,
			},
		},
		{
			name: "values below 0 should be clamped",
			input: PersonalityTraits{
				Openness:          -0.5,
				Conscientiousness: -1.0,
				Extraversion:      -0.1,
				Agreeableness:     -0.3,
				Neuroticism:       -0.8,
			},
			expected: PersonalityTraits{
				Openness:          0.0,
				Conscientiousness: 0.0,
				Extraversion:      0.0,
				Agreeableness:     0.0,
				Neuroticism:       0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePersonality(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePersonality() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCharacterLocking(t *testing.T) {
	char := &Character{
		ID:   "test-123",
		Name: "Test Character",
	}

	// Test write lock
	char.Lock()
	char.Name = "Modified Character"
	char.Unlock()

	if char.Name != "Modified Character" {
		t.Errorf("Expected name to be 'Modified Character', got %s", char.Name)
	}

	// Test read lock
	char.RLock()
	name := char.Name
	char.RUnlock()

	if name != "Modified Character" {
		t.Errorf("Expected to read 'Modified Character', got %s", name)
	}
}

func TestMemoryTypes(t *testing.T) {
	memories := []Memory{
		{
			Type:      ShortTermMemory,
			Content:   "Recent conversation",
			Timestamp: time.Now(),
			Emotional: 0.5,
		},
		{
			Type:      MediumTermMemory,
			Content:   "Important interaction",
			Timestamp: time.Now().Add(-time.Hour),
			Emotional: 0.8,
		},
		{
			Type:      LongTermMemory,
			Content:   "Core memory",
			Timestamp: time.Now().Add(-24 * time.Hour),
			Emotional: 0.95,
		},
	}

	for _, mem := range memories {
		switch mem.Type {
		case ShortTermMemory, MediumTermMemory, LongTermMemory:
			// Valid memory type
		default:
			t.Errorf("Invalid memory type: %s", mem.Type)
		}
	}
}