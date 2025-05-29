// +build integration

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
)

// These tests require actual API access and should be run with:
// go test -tags=integration ./test/...

func TestEndToEndChat(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("Skipping integration test: no API key found")
	}

	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	// Create config
	cfg := &config.Config{
		DefaultProvider: "openai",
		Model:          "gpt-4o-mini",
		APIKey:         os.Getenv("OPENAI_API_KEY"),
		CacheConfig: config.CacheConfig{
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
			StabilityThreshold: 10,
		},
		UserProfileConfig: config.UserProfileConfig{
			Enabled:             true,
			UpdateFrequency:     5,
			TurnsToConsider:     20,
			ConfidenceThreshold: 0.5,
			PromptCacheTTL:      1 * time.Hour,
		},
	}

	// Create manager
	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create a test character
	character := &models.Character{
		ID:        "integration-test-char",
		Name:      "Integration Test Character",
		Backstory: "A helpful AI assistant created for integration testing. Always positive and encouraging.",
		Personality: models.PersonalityTraits{
			Openness:          0.9,
			Conscientiousness: 0.8,
			Extraversion:      0.7,
			Agreeableness:     0.95,
			Neuroticism:       0.2,
		},
		CurrentMood: models.EmotionalState{
			Joy:      0.8,
			Surprise: 0.2,
		},
		Quirks: []string{
			"Always starts responses with enthusiasm",
			"Loves to encourage learning",
			"Uses emojis occasionally",
		},
		SpeechStyle: "Friendly, encouraging, and clear. Uses positive language and helpful examples.",
	}

	err = mgr.CreateCharacter(character)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Start a conversation
	ctx := context.Background()
	sessionID := fmt.Sprintf("integration-test-%d", time.Now().Unix())

	// First message
	req1 := &models.ConversationRequest{
		CharacterID: character.ID,
		UserID:      "test-user",
		Message:     "Hello! I'm learning Go programming. Can you help me?",
		SessionID:   sessionID,
	}

	resp1, err := mgr.ProcessMessage(ctx, req1)
	if err != nil {
		t.Fatalf("Failed to process first message: %v", err)
	}

	if resp1.Content == "" {
		t.Error("Expected non-empty response")
	}

	t.Logf("First response: %s", resp1.Content)

	// Second message - should use context
	req2 := &models.ConversationRequest{
		CharacterID: character.ID,
		UserID:      "test-user",
		Message:     "What are the best practices for error handling?",
		SessionID:   sessionID,
	}

	resp2, err := mgr.ProcessMessage(ctx, req2)
	if err != nil {
		t.Fatalf("Failed to process second message: %v", err)
	}

	t.Logf("Second response: %s", resp2.Content)

	// Verify session was saved
	sessionRepo := repository.NewSessionRepository(filepath.Join(tempDir, ".config", "roleplay"))
	session, err := sessionRepo.LoadSession(character.ID, sessionID)
	if err != nil {
		t.Fatalf("Failed to load session: %v", err)
	}

	if len(session.Messages) != 4 { // 2 user + 2 assistant
		t.Errorf("Expected 4 messages in session, got %d", len(session.Messages))
	}

	// Verify cache metrics
	if session.CacheMetrics.TotalRequests != 2 {
		t.Errorf("Expected 2 total requests, got %d", session.CacheMetrics.TotalRequests)
	}

	t.Logf("Cache hit rate: %.1f%%", session.CacheMetrics.HitRate*100)
}

func TestCachingBehavior(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping caching test: OPENAI_API_KEY not set")
	}

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		DefaultProvider: "openai",
		Model:          "gpt-4o-mini",
		APIKey:         os.Getenv("OPENAI_API_KEY"),
		CacheConfig: config.CacheConfig{
			CleanupInterval:   5 * time.Minute,
			DefaultTTL:        10 * time.Minute,
			EnableAdaptiveTTL: true,
		},
	}

	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create character
	character := &models.Character{
		ID:        "cache-test-char",
		Name:      "Cache Test Character",
		Backstory: "A character designed to test caching behavior",
	}

	err = mgr.CreateCharacter(character)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	ctx := context.Background()

	// Send identical messages
	message := "What is the meaning of life?"
	responses := make([]*providers.AIResponse, 3)
	times := make([]time.Duration, 3)

	for i := 0; i < 3; i++ {
		start := time.Now()
		req := &models.ConversationRequest{
			CharacterID: character.ID,
			UserID:      "cache-test-user",
			Message:     message,
			SessionID:   fmt.Sprintf("cache-test-%d", i),
		}

		resp, err := mgr.ProcessMessage(ctx, req)
		if err != nil {
			t.Fatalf("Failed to process message %d: %v", i, err)
		}

		responses[i] = resp
		times[i] = time.Since(start)

		t.Logf("Request %d: %v (cache hit: %v)", i+1, times[i], resp.CacheMetrics.Hit)

		// Small delay to ensure cache is written
		time.Sleep(100 * time.Millisecond)
	}

	// First request should be slowest (cache miss)
	// Subsequent requests should be faster (cache hits)
	if times[1] >= times[0] {
		t.Error("Second request should be faster due to cache hit")
	}

	if times[2] >= times[0] {
		t.Error("Third request should be faster due to cache hit")
	}

	// Verify cache hits
	if responses[0].CacheMetrics.Hit {
		t.Error("First response should be cache miss")
	}

	if !responses[1].CacheMetrics.Hit || !responses[2].CacheMetrics.Hit {
		t.Error("Second and third responses should be cache hits")
	}

	// Responses should be identical
	if responses[0].Content != responses[1].Content || responses[1].Content != responses[2].Content {
		t.Error("Cached responses should be identical")
	}
}

func TestSessionPersistence(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	// Use mock provider for this test
	cfg := &config.Config{
		DefaultProvider: "mock",
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}

	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Register mock provider
	mockProvider := &mockIntegrationProvider{
		responses: map[string]string{
			"Hello":          "Hi there! How can I help you?",
			"Remember this":  "I'll remember that for you.",
			"What did I say": "You asked me to remember something.",
		},
	}
	mgr.RegisterProvider("mock", mockProvider)

	// Create character
	character := &models.Character{
		ID:   "persist-test-char",
		Name: "Persistence Test Character",
	}
	mgr.CreateCharacter(character)

	ctx := context.Background()
	sessionID := "persist-test-session"

	// First conversation
	messages := []string{"Hello", "Remember this"}
	for _, msg := range messages {
		req := &models.ConversationRequest{
			CharacterID: character.ID,
			UserID:      "persist-user",
			Message:     msg,
			SessionID:   sessionID,
		}

		_, err := mgr.ProcessMessage(ctx, req)
		if err != nil {
			t.Fatalf("Failed to process message '%s': %v", msg, err)
		}
	}

	// Create new manager instance (simulating restart)
	mgr2, err := manager.NewCharacterManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create second manager: %v", err)
	}
	mgr2.RegisterProvider("mock", mockProvider)

	// Load character into new manager
	mgr2.CreateCharacter(character)

	// Continue conversation with same session
	req := &models.ConversationRequest{
		CharacterID: character.ID,
		UserID:      "persist-user",
		Message:     "What did I say",
		SessionID:   sessionID,
	}

	// Load session history
	sessionRepo := repository.NewSessionRepository(filepath.Join(tempDir, ".config", "roleplay"))
	session, err := sessionRepo.LoadSession(character.ID, sessionID)
	if err != nil {
		t.Fatalf("Failed to load session: %v", err)
	}

	// Build context from session
	var contextMessages []models.Message
	for _, msg := range session.Messages {
		contextMessages = append(contextMessages, models.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	req.Context.RecentMessages = contextMessages

	resp, err := mgr2.ProcessMessage(ctx, req)
	if err != nil {
		t.Fatalf("Failed to process message after restart: %v", err)
	}

	// Should have context from previous conversation
	if resp.Content != mockProvider.responses["What did I say"] {
		t.Errorf("Expected response with context, got: %s", resp.Content)
	}

	// Verify session has all messages
	finalSession, err := sessionRepo.LoadSession(character.ID, sessionID)
	if err != nil {
		t.Fatalf("Failed to load final session: %v", err)
	}

	expectedMessages := 6 // 3 pairs of user/assistant messages
	if len(finalSession.Messages) != expectedMessages {
		t.Errorf("Expected %d messages, got %d", expectedMessages, len(finalSession.Messages))
	}
}

func TestMultiCharacterConversations(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		DefaultProvider: "mock",
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}

	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Register mock provider
	mgr.RegisterProvider("mock", &mockIntegrationProvider{
		responses: map[string]string{
			"default": "Response from character",
		},
	})

	// Create multiple characters
	characters := []*models.Character{
		{
			ID:        "char1",
			Name:      "Character One",
			Backstory: "First test character",
			SpeechStyle: "Formal and polite",
		},
		{
			ID:        "char2",
			Name:      "Character Two",
			Backstory: "Second test character",
			SpeechStyle: "Casual and friendly",
		},
		{
			ID:        "char3",
			Name:      "Character Three",
			Backstory: "Third test character",
			SpeechStyle: "Technical and precise",
		},
	}

	for _, char := range characters {
		if err := mgr.CreateCharacter(char); err != nil {
			t.Fatalf("Failed to create character %s: %v", char.ID, err)
		}
	}

	ctx := context.Background()
	userID := "multi-char-user"

	// Have conversations with each character
	for _, char := range characters {
		sessionID := fmt.Sprintf("session-%s", char.ID)
		
		for i := 0; i < 3; i++ {
			req := &models.ConversationRequest{
				CharacterID: char.ID,
				UserID:      userID,
				Message:     fmt.Sprintf("Message %d to %s", i+1, char.Name),
				SessionID:   sessionID,
			}

			resp, err := mgr.ProcessMessage(ctx, req)
			if err != nil {
				t.Fatalf("Failed to process message for %s: %v", char.ID, err)
			}

			if resp.Content == "" {
				t.Errorf("Empty response from %s", char.ID)
			}
		}
	}

	// Verify each character has separate sessions
	sessionRepo := repository.NewSessionRepository(filepath.Join(tempDir, ".config", "roleplay"))
	
	for _, char := range characters {
		sessions, err := sessionRepo.ListSessions(char.ID)
		if err != nil {
			t.Fatalf("Failed to list sessions for %s: %v", char.ID, err)
		}

		if len(sessions) != 1 {
			t.Errorf("Expected 1 session for %s, got %d", char.ID, len(sessions))
		}

		session := sessions[0]
		if len(session.Messages) != 6 { // 3 user + 3 assistant
			t.Errorf("Expected 6 messages for %s, got %d", char.ID, len(session.Messages))
		}
	}
}

func TestProviderFailover(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		DefaultProvider: "primary",
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
	}

	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Register providers
	failingProvider := &mockIntegrationProvider{
		failAfter: 2,
		responses: map[string]string{
			"default": "Response from primary",
		},
	}
	
	fallbackProvider := &mockIntegrationProvider{
		responses: map[string]string{
			"default": "Response from fallback",
		},
	}

	mgr.RegisterProvider("primary", failingProvider)
	mgr.RegisterProvider("fallback", fallbackProvider)

	// Create character
	character := &models.Character{
		ID:   "failover-test",
		Name: "Failover Test Character",
	}
	mgr.CreateCharacter(character)

	ctx := context.Background()

	// First two requests should succeed
	for i := 0; i < 2; i++ {
		req := &models.ConversationRequest{
			CharacterID: character.ID,
			UserID:      "test-user",
			Message:     fmt.Sprintf("Message %d", i+1),
		}

		resp, err := mgr.ProcessMessage(ctx, req)
		if err != nil {
			t.Fatalf("Request %d failed unexpectedly: %v", i+1, err)
		}

		if resp.Content != "Response from primary" {
			t.Errorf("Expected response from primary provider, got: %s", resp.Content)
		}
	}

	// Third request should fail
	req := &models.ConversationRequest{
		CharacterID: character.ID,
		UserID:      "test-user",
		Message:     "Message 3",
	}

	_, err = mgr.ProcessMessage(ctx, req)
	if err == nil {
		t.Error("Expected error from failing provider")
	}

	// In a real implementation, we would switch to fallback provider here
	// For this test, we'll manually switch
	cfg.DefaultProvider = "fallback"

	// Retry with fallback
	resp, err := mgr.ProcessMessage(ctx, req)
	if err != nil {
		t.Fatalf("Fallback request failed: %v", err)
	}

	if resp.Content != "Response from fallback" {
		t.Errorf("Expected response from fallback provider, got: %s", resp.Content)
	}
}

// Mock provider for integration tests
type mockIntegrationProvider struct {
	responses   map[string]string
	failAfter   int
	callCount   int
}

func (m *mockIntegrationProvider) SendRequest(ctx context.Context, req *providers.PromptRequest) (*providers.AIResponse, error) {
	m.callCount++
	
	if m.failAfter > 0 && m.callCount > m.failAfter {
		return nil, fmt.Errorf("provider failed after %d calls", m.failAfter)
	}

	// Simple response based on message content
	response := m.responses["default"]
	for key, resp := range m.responses {
		if strings.Contains(req.Message, key) {
			response = resp
			break
		}
	}

	return &providers.AIResponse{
		Content: response,
		TokensUsed: providers.TokenUsage{
			Prompt:     100,
			Completion: 50,
			Total:      150,
		},
		CacheMetrics: providers.CacheMetrics{
			Hit: false,
		},
	}, nil
}

func (m *mockIntegrationProvider) SendStreamRequest(ctx context.Context, req *providers.PromptRequest, out chan<- providers.PartialAIResponse) error {
	defer close(out)
	
	resp, err := m.SendRequest(ctx, req)
	if err != nil {
		return err
	}

	out <- providers.PartialAIResponse{
		Content: resp.Content,
		Done:    true,
	}
	return nil
}

func (m *mockIntegrationProvider) Name() string {
	return "mock-integration"
}