package services

import (
	"context"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
)

// Mock provider for testing
type mockProvider struct {
	name        string
	breakpoints bool
	maxBreaks   int
	response    *providers.AIResponse
	err         error
}

func (m *mockProvider) SendRequest(ctx context.Context, req *providers.PromptRequest) (*providers.AIResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.response != nil {
		return m.response, nil
	}
	return &providers.AIResponse{
		Content: "Mock response",
		TokensUsed: providers.TokenUsage{
			Prompt:     100,
			Completion: 50,
			Total:      150,
		},
	}, nil
}

func (m *mockProvider) SupportsBreakpoints() bool { return m.breakpoints }
func (m *mockProvider) MaxBreakpoints() int      { return m.maxBreaks }
func (m *mockProvider) Name() string              { return m.name }

func TestCharacterBot(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "mock",
		CacheConfig: config.CacheConfig{
			MaxEntries:        100,
			CleanupInterval:   5 * time.Minute,
			DefaultTTL:        10 * time.Minute,
			EnableAdaptiveTTL: true,
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    10,
			MediumTermDuration: 24 * time.Hour,
			ConsolidationRate:  0.1,
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled:   true,
			MaxDriftRate:       0.02,
			StabilityThreshold: 5,
		},
	}

	bot := NewCharacterBot(cfg)

	// Register mock provider
	mockProv := &mockProvider{name: "mock", breakpoints: true, maxBreaks: 4}
	bot.RegisterProvider("mock", mockProv)

	// Create character
	char := &models.Character{
		ID:        "test-char",
		Name:      "Test Character",
		Backstory: "A test character for unit tests",
		Personality: models.PersonalityTraits{
			Openness:          0.5,
			Conscientiousness: 0.5,
			Extraversion:      0.5,
			Agreeableness:     0.5,
			Neuroticism:       0.5,
		},
		CurrentMood: models.EmotionalState{
			Joy: 0.5,
		},
		SpeechStyle: "Test speech",
		Quirks:      []string{"Test quirk"},
	}

	err := bot.CreateCharacter(char)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Test duplicate character
	err = bot.CreateCharacter(char)
	if err == nil {
		t.Error("Expected error creating duplicate character")
	}

	// Get character
	retrieved, err := bot.GetCharacter("test-char")
	if err != nil {
		t.Fatalf("Failed to get character: %v", err)
	}

	if retrieved.Name != "Test Character" {
		t.Errorf("Expected name 'Test Character', got %s", retrieved.Name)
	}
}

func TestBuildPrompt(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "mock",
		CacheConfig: config.CacheConfig{
			DefaultTTL:      10 * time.Minute,
			CleanupInterval: 5 * time.Minute,
		},
	}

	bot := NewCharacterBot(cfg)

	// Create character
	char := &models.Character{
		ID:        "prompt-test",
		Name:      "Prompt Test",
		Backstory: "Testing prompt building",
		Personality: models.PersonalityTraits{
			Openness: 0.7,
		},
		CurrentMood: models.EmotionalState{
			Joy: 0.8,
		},
		SpeechStyle: "Test style",
		Memories: []models.Memory{
			{
				Type:      models.MediumTermMemory,
				Content:   "Test pattern",
				Emotional: 0.5,
			},
		},
	}

	if err := bot.CreateCharacter(char); err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Build prompt
	req := &models.ConversationRequest{
		CharacterID: "prompt-test",
		UserID:      "user-123",
		Message:     "Hello",
		Context: models.ConversationContext{
			RecentMessages: []models.Message{
				{Role: "user", Content: "Previous message"},
			},
		},
	}

	prompt, breakpoints, err := bot.BuildPrompt(req)
	if err != nil {
		t.Fatalf("Failed to build prompt: %v", err)
	}

	// Verify prompt contains expected content
	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	// Verify breakpoints
	if len(breakpoints) < 3 {
		t.Errorf("Expected at least 3 breakpoints, got %d", len(breakpoints))
	}

	// Check for personality layer
	foundPersonality := false
	for _, bp := range breakpoints {
		if bp.Layer == "core_personality" {
			foundPersonality = true
			break
		}
	}
	if !foundPersonality {
		t.Error("Expected to find core personality layer in breakpoints")
	}
}

func TestMemoryConsolidation(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    3,
			MediumTermDuration: 1 * time.Hour,
			ConsolidationRate:  0.1,
		},
	}

	bot := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "memory-test",
		Name: "Memory Test",
		Memories: []models.Memory{
			{Type: models.ShortTermMemory, Content: "Memory 1", Emotional: 0.8, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 2", Emotional: 0.9, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 3", Emotional: 0.7, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 4", Emotional: 0.85, Timestamp: time.Now()},
			{Type: models.ShortTermMemory, Content: "Memory 5", Emotional: 0.75, Timestamp: time.Now()},
		},
	}

	// Consolidate memories
	bot.consolidateMemories(char)

	// Check for consolidated memory
	hasMediumTerm := false
	for _, mem := range char.Memories {
		if mem.Type == models.MediumTermMemory {
			hasMediumTerm = true
			break
		}
	}

	if !hasMediumTerm {
		t.Error("Expected to find consolidated medium-term memory")
	}
}

func TestPersonalityEvolution(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled:   true,
			MaxDriftRate:       0.1,
			StabilityThreshold: 0,
		},
	}

	bot := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "evolution-test",
		Name: "Evolution Test",
		Personality: models.PersonalityTraits{
			Openness:          0.5,
			Conscientiousness: 0.5,
			Extraversion:      0.5,
			Agreeableness:     0.5,
			Neuroticism:       0.5,
		},
	}

	// Create response with emotional impact
	resp := &providers.AIResponse{
		Emotions: models.EmotionalState{
			Joy:      1.0, // High joy should increase extraversion
			Surprise: 1.0, // High surprise should increase openness
		},
	}

	originalOpenness := char.Personality.Openness
	originalExtraversion := char.Personality.Extraversion

	bot.evolvePersonality(char, resp)

	// Verify personality changed but within bounds
	if char.Personality.Openness <= originalOpenness {
		t.Error("Expected openness to increase with high surprise")
	}

	if char.Personality.Extraversion <= originalExtraversion {
		t.Error("Expected extraversion to increase with high joy")
	}

	// Verify changes are bounded
	maxChange := cfg.PersonalityConfig.MaxDriftRate
	if char.Personality.Openness-originalOpenness > maxChange {
		t.Error("Personality change exceeded max drift rate")
	}
}