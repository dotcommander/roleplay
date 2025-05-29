package services

import (
	"strings"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
)

// TestUniversalCachingPrefixConsistency verifies that identical prefixes are generated
// for the same character/user/scenario combination across multiple requests
func TestUniversalCachingPrefixConsistency(t *testing.T) {
	// Create a test bot
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 0, // Disable background cleanup for test
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    10,
			MediumTermDuration: time.Hour,
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled: false,
		},
	}

	bot := NewCharacterBot(cfg)

	// Create a test character
	char := &models.Character{
		ID:   "test-char",
		Name: "Test Character",
		Personality: models.PersonalityTraits{
			Openness:          0.7,
			Conscientiousness: 0.8,
			Extraversion:      0.6,
			Agreeableness:     0.5,
			Neuroticism:       0.3,
		},
		Backstory:   "A test character for caching validation",
		SpeechStyle: "Professional and clear",
		Quirks:      []string{"Always punctual", "Loves testing"},
		CurrentMood: models.EmotionalState{
			Joy: 0.5,
		},
	}

	err := bot.CreateCharacter(char)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Create multiple requests with the same character/user combination
	requests := []*models.ConversationRequest{
		{
			CharacterID: "test-char",
			UserID:      "test-user",
			Message:     "Hello",
			Context: models.ConversationContext{
				SessionID: "test-session",
			},
		},
		{
			CharacterID: "test-char",
			UserID:      "test-user",
			Message:     "How are you?",
			Context: models.ConversationContext{
				SessionID: "test-session",
				RecentMessages: []models.Message{
					{Role: "user", Content: "Hello"},
					{Role: "assistant", Content: "Hello! I'm doing well."},
				},
			},
		},
		{
			CharacterID: "test-char",
			UserID:      "test-user",
			Message:     "Tell me about yourself",
			Context: models.ConversationContext{
				SessionID: "test-session",
				RecentMessages: []models.Message{
					{Role: "user", Content: "Hello"},
					{Role: "assistant", Content: "Hello! I'm doing well."},
					{Role: "user", Content: "How are you?"},
					{Role: "assistant", Content: "I'm great, thanks for asking!"},
				},
			},
		},
	}

	// Build prompts for each request and extract prefixes
	prefixes := make([]string, len(requests))
	for i, req := range requests {
		prompt, breakpoints, err := bot.BuildPrompt(req)
		if err != nil {
			t.Fatalf("Failed to build prompt for request %d: %v", i, err)
		}

		// Extract the prefix (everything before the conversation context separator)
		parts := strings.Split(prompt, "===== CONVERSATION CONTEXT =====")
		if len(parts) != 2 {
			t.Fatalf("Prompt %d does not have expected structure with separator", i)
		}
		prefixes[i] = strings.TrimSpace(parts[0])

		// Verify the prompt has the expected structure
		if !strings.Contains(prompt, "[SYSTEM INSTRUCTIONS]") {
			t.Errorf("Prompt %d missing system instructions", i)
		}
		if !strings.Contains(prompt, "[CHARACTER PROFILE]") {
			t.Errorf("Prompt %d missing character profile", i)
		}
		if !strings.Contains(prompt, "[USER CONTEXT]") {
			t.Errorf("Prompt %d missing user context", i)
		}

		// Verify cache key generation is consistent
		cacheKey := bot.generateCacheKey(req.CharacterID, req.UserID, req.ScenarioID, breakpoints)
		t.Logf("Request %d cache key: %s", i, cacheKey)
	}

	// Verify all prefixes are identical
	for i := 1; i < len(prefixes); i++ {
		if prefixes[i] != prefixes[0] {
			t.Errorf("Prefix %d differs from prefix 0:\nPrefix 0:\n%s\n\nPrefix %d:\n%s", 
				i, prefixes[0], i, prefixes[i])
		}
	}

	// Test with scenario
	scenarioRequests := []*models.ConversationRequest{
		{
			CharacterID: "test-char",
			UserID:      "test-user",
			ScenarioID:  "test-scenario",
			Message:     "Start scenario",
			Context: models.ConversationContext{
				SessionID: "test-session-2",
			},
		},
		{
			CharacterID: "test-char",
			UserID:      "test-user",
			ScenarioID:  "test-scenario",
			Message:     "Continue scenario",
			Context: models.ConversationContext{
				SessionID: "test-session-2",
				RecentMessages: []models.Message{
					{Role: "user", Content: "Start scenario"},
					{Role: "assistant", Content: "Scenario started!"},
				},
			},
		},
	}

	// Since we don't have a real scenario repo, we'll just verify the structure
	scenarioPrefixes := make([]string, len(scenarioRequests))
	for i, req := range scenarioRequests {
		// This will fail to load the scenario but continue without it
		prompt, _, err := bot.BuildPrompt(req)
		if err != nil {
			t.Fatalf("Failed to build prompt for scenario request %d: %v", i, err)
		}

		parts := strings.Split(prompt, "===== CONVERSATION CONTEXT =====")
		if len(parts) == 2 {
			scenarioPrefixes[i] = strings.TrimSpace(parts[0])
		}
	}

	// Verify scenario prefixes are consistent with each other
	if len(scenarioPrefixes) > 1 && scenarioPrefixes[0] != "" {
		for i := 1; i < len(scenarioPrefixes); i++ {
			if scenarioPrefixes[i] != scenarioPrefixes[0] {
				t.Errorf("Scenario prefix %d differs from scenario prefix 0", i)
			}
		}
	}
}

// TestCachingLayerOrdering verifies that layers are always added in the same order
func TestCachingLayerOrdering(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			DefaultTTL: 5 * time.Minute,
		},
	}

	bot := NewCharacterBot(cfg)

	// Create a character
	char := &models.Character{
		ID:          "order-test",
		Name:        "Order Tester",
		Personality: models.PersonalityTraits{},
		Backstory:   "Testing layer ordering",
	}

	err := bot.CreateCharacter(char)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	req := &models.ConversationRequest{
		CharacterID: "order-test",
		UserID:      "test-user",
		Message:     "Test",
		Context:     models.ConversationContext{},
	}

	_, breakpoints, err := bot.BuildPrompt(req)
	if err != nil {
		t.Fatalf("Failed to build prompt: %v", err)
	}

	// Verify layer ordering
	expectedOrder := []string{
		"system_admin",
		"core_personality",
		"learned_behavior", // might be skipped if no learned behaviors
		"emotional_state",
		"user_memory",
		"conversation", // might be skipped if no conversation history
	}

	layerIndex := 0
	for _, bp := range breakpoints {
		// Find this layer in our expected order
		found := false
		for i := layerIndex; i < len(expectedOrder); i++ {
			if string(bp.Layer) == expectedOrder[i] {
				layerIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected layer %s or out of order", bp.Layer)
		}
	}
}

// TestPrefixDeterminism verifies that the prefix content is deterministic
func TestPrefixDeterminism(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			DefaultTTL: 5 * time.Minute,
		},
	}

	// Create multiple bots to ensure prefix is consistent across instances
	bot1 := NewCharacterBot(cfg)
	bot2 := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "determinism-test",
		Name: "Determinism Tester",
		Personality: models.PersonalityTraits{
			Openness:          0.5,
			Conscientiousness: 0.5,
			Extraversion:      0.5,
			Agreeableness:     0.5,
			Neuroticism:       0.5,
		},
		Backstory:   "A character for testing deterministic prompt generation",
		SpeechStyle: "Precise and consistent",
		Quirks:      []string{"Values consistency", "Dislikes randomness"},
	}

	// Create character in both bots
	err := bot1.CreateCharacter(char)
	if err != nil {
		t.Fatalf("Failed to create character in bot1: %v", err)
	}

	err = bot2.CreateCharacter(char)
	if err != nil {
		t.Fatalf("Failed to create character in bot2: %v", err)
	}

	req := &models.ConversationRequest{
		CharacterID: "determinism-test",
		UserID:      "test-user",
		Message:     "Test message",
		Context:     models.ConversationContext{},
	}

	// Build prompts from both bots
	prompt1, _, err := bot1.BuildPrompt(req)
	if err != nil {
		t.Fatalf("Failed to build prompt from bot1: %v", err)
	}

	prompt2, _, err := bot2.BuildPrompt(req)
	if err != nil {
		t.Fatalf("Failed to build prompt from bot2: %v", err)
	}

	// Extract prefixes
	parts1 := strings.Split(prompt1, "===== CONVERSATION CONTEXT =====")
	parts2 := strings.Split(prompt2, "===== CONVERSATION CONTEXT =====")

	if len(parts1) != 2 || len(parts2) != 2 {
		t.Fatal("Prompts do not have expected structure")
	}

	prefix1 := strings.TrimSpace(parts1[0])
	prefix2 := strings.TrimSpace(parts2[0])

	if prefix1 != prefix2 {
		t.Errorf("Prefixes from different bot instances differ:\nBot1:\n%s\n\nBot2:\n%s", 
			prefix1, prefix2)
	}
}