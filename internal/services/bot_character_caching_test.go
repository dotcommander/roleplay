package services

import (
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharacterSystemPromptCaching(t *testing.T) {
	// Create a bot with character prompt caching configured
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			MaxEntries:                   1000,
			CleanupInterval:              5 * time.Minute,
			DefaultTTL:                   1 * time.Hour,
			CoreCharacterSystemPromptTTL: 7 * 24 * time.Hour, // 7 days
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    10,
			MediumTermDuration: 24 * time.Hour,
			ConsolidationRate:  0.1,
		},
	}

	bot := NewCharacterBot(cfg)

	// Create a test character
	char := &models.Character{
		ID:   "test-wizard",
		Name: "Gandalf",
		Personality: models.PersonalityTraits{
			Openness:          0.9,
			Conscientiousness: 0.8,
			Extraversion:      0.6,
			Agreeableness:     0.7,
			Neuroticism:       0.2,
		},
		Backstory:   "Ancient wizard who has walked Middle-earth for thousands of years",
		SpeechStyle: "Wise, cryptic, uses archaic language",
		Quirks:      []string{"Speaks in riddles", "Fond of pipe-weed"},
	}

	// Create the character
	err := bot.CreateCharacter(char)
	require.NoError(t, err)

	// Test 1: Verify cache warmup happened
	cacheKey := bot.generateCharacterSystemPromptCacheKey(char.ID)
	cachedEntry, exists := bot.cache.Get(cacheKey)
	require.True(t, exists, "Cache entry should exist after character creation")
	require.NotNil(t, cachedEntry)
	
	// Find the core personality layer content
	var foundContent string
	for _, bp := range cachedEntry.Breakpoints {
		if bp.Layer == cache.CorePersonalityLayer {
			foundContent = bp.Content
			break
		}
	}
	require.NotEmpty(t, foundContent, "Should have cached core personality content")
	assert.Contains(t, foundContent, "Gandalf")
	assert.Contains(t, foundContent, "Ancient wizard")

	// Test 2: BuildPrompt should use cached content
	req := &models.ConversationRequest{
		CharacterID: char.ID,
		UserID:      "test-user",
		Message:     "Hello!",
		Context: models.ConversationContext{
			RecentMessages: []models.Message{},
		},
	}

	prompt1, breakpoints1, err := bot.BuildPrompt(req)
	require.NoError(t, err)
	assert.NotEmpty(t, prompt1)
	assert.NotEmpty(t, breakpoints1)

	// Find the core personality layer
	var personalityBreakpoint *cache.CacheBreakpoint
	for i, bp := range breakpoints1 {
		if bp.Layer == cache.CorePersonalityLayer {
			personalityBreakpoint = &breakpoints1[i]
			break
		}
	}
	require.NotNil(t, personalityBreakpoint)
	// Check if it contains the cached marker
	assert.Contains(t, personalityBreakpoint.Content, "<!-- cached:true -->", "Should be served from cache")

	// Test 3: Second call should also be cached
	prompt2, breakpoints2, err := bot.BuildPrompt(req)
	require.NoError(t, err)
	assert.Equal(t, prompt1, prompt2, "Prompts should be identical when using cache")

	// Find personality breakpoint again
	for i, bp := range breakpoints2 {
		if bp.Layer == cache.CorePersonalityLayer {
			personalityBreakpoint = &breakpoints2[i]
			break
		}
	}
	assert.Contains(t, personalityBreakpoint.Content, "<!-- cached:true -->", "Should still be served from cache")

	// Test 4: Cache invalidation
	err = bot.InvalidateCharacterCache(char.ID)
	require.NoError(t, err)

	// After invalidation, the cache should be rebuilt
	cachedEntry, exists = bot.cache.Get(cacheKey)
	require.True(t, exists, "Cache should be rebuilt after invalidation")
	require.NotNil(t, cachedEntry)
	
	// Verify it has content
	foundContent = ""
	for _, bp := range cachedEntry.Breakpoints {
		if bp.Layer == cache.CorePersonalityLayer {
			foundContent = bp.Content
			break
		}
	}
	require.NotEmpty(t, foundContent, "Cache should have been rebuilt with content")
}

func TestCoreCharacterSystemPromptGeneration(t *testing.T) {
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CoreCharacterSystemPromptTTL: 7 * 24 * time.Hour,
		},
	}

	bot := NewCharacterBot(cfg)

	testCases := []struct {
		name     string
		char     *models.Character
		expected []string
	}{
		{
			name: "High openness character",
			char: &models.Character{
				ID:   "creative-1",
				Name: "Luna",
				Personality: models.PersonalityTraits{
					Openness: 0.9,
				},
				Backstory:   "An artist who sees the world differently",
				SpeechStyle: "Poetic and imaginative",
				Quirks:      []string{"Uses color metaphors"},
			},
			expected: []string{
				"Creative, curious, open to new experiences",
				"Luna",
				"An artist who sees the world differently",
				"Poetic and imaginative",
				"Uses color metaphors",
			},
		},
		{
			name: "Low neuroticism character",
			char: &models.Character{
				ID:   "stable-1",
				Name: "Marcus",
				Personality: models.PersonalityTraits{
					Neuroticism: 0.1,
				},
				Backstory:   "A calm meditation teacher",
				SpeechStyle: "Measured and peaceful",
				Quirks:      []string{"Never raises voice"},
			},
			expected: []string{
				"Emotionally stable, calm, resilient",
				"Marcus",
				"A calm meditation teacher",
				"Measured and peaceful",
				"Never raises voice",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompt := bot.buildCoreCharacterSystemPrompt(tc.char)
			
			// Check all expected content is present
			for _, expected := range tc.expected {
				assert.Contains(t, prompt, expected)
			}
			
			// Verify structure
			assert.Contains(t, prompt, "[CHARACTER FOUNDATION]")
			assert.Contains(t, prompt, "[PERSONALITY MATRIX - OCEAN MODEL]")
			assert.Contains(t, prompt, "[COMPREHENSIVE BACKSTORY]")
			assert.Contains(t, prompt, "[SPEECH CHARACTERISTICS]")
			assert.Contains(t, prompt, "[DEFINING QUIRKS AND MANNERISMS]")
			assert.Contains(t, prompt, "[CORE INTERACTION PRINCIPLES]")
		})
	}
}

func TestCharacterCacheTTL(t *testing.T) {
	// Test that the TTL is properly applied
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CoreCharacterSystemPromptTTL: 48 * time.Hour, // 2 days for testing
		},
	}

	bot := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "ttl-test",
		Name: "TTL Tester",
		Personality: models.PersonalityTraits{
			Openness: 0.5,
		},
		Backstory:   "Test character for TTL",
		SpeechStyle: "Normal",
		Quirks:      []string{},
	}

	// Create character (triggers warmup)
	err := bot.CreateCharacter(char)
	require.NoError(t, err)

	// Check cache entry TTL
	cacheKey := bot.generateCharacterSystemPromptCacheKey(char.ID)
	cachedEntry, exists := bot.cache.Get(cacheKey)
	require.True(t, exists, "Cache entry should exist")
	require.NotNil(t, cachedEntry)

	// Find the core personality breakpoint and check its TTL
	for _, bp := range cachedEntry.Breakpoints {
		if bp.Layer == cache.CorePersonalityLayer {
			// The TTL should be what we configured
			assert.Equal(t, cfg.CacheConfig.CoreCharacterSystemPromptTTL, bp.TTL)
			break
		}
	}
}

func TestCharacterPromptConsistency(t *testing.T) {
	// Test that the same character always generates the same prompt
	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CoreCharacterSystemPromptTTL: 24 * time.Hour,
		},
	}

	bot1 := NewCharacterBot(cfg)
	bot2 := NewCharacterBot(cfg)

	char := &models.Character{
		ID:   "consistency-test",
		Name: "Consistent Charlie",
		Personality: models.PersonalityTraits{
			Openness:          0.6,
			Conscientiousness: 0.7,
			Extraversion:      0.5,
			Agreeableness:     0.8,
			Neuroticism:       0.3,
		},
		Backstory:   "A reliable friend who values consistency",
		SpeechStyle: "Clear and direct",
		Quirks:      []string{"Always on time", "Loves routine"},
	}

	// Generate prompts from two different bot instances
	prompt1 := bot1.buildCoreCharacterSystemPrompt(char)
	prompt2 := bot2.buildCoreCharacterSystemPrompt(char)

	// They should be identical
	assert.Equal(t, prompt1, prompt2, "Same character should generate identical prompts")
}