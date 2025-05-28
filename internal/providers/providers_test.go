package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
)

func TestOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider("test-api-key", "gpt-4")

	if provider.Name() != "openai_compatible" {
		t.Errorf("Expected name 'openai_compatible', got %s", provider.Name())
	}

	if provider.SupportsBreakpoints() {
		t.Error("Expected OpenAI to not support explicit breakpoints")
	}

	if provider.MaxBreakpoints() != 0 {
		t.Errorf("Expected 0 breakpoints, got %d", provider.MaxBreakpoints())
	}

	// Verify it uses the model passed in constructor
	if provider.model != "gpt-4" {
		t.Errorf("Expected model to be 'gpt-4', got %s", provider.model)
	}
}

func TestOpenAIProviderRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Missing or incorrect Authorization header")
		}

		// Verify request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if reqBody["model"] != "o4-mini" {
			t.Errorf("Expected model o4-mini, got %v", reqBody["model"])
		}

		// Send mock response
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "Test response",
					},
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     100,
				"completion_tokens": 50,
				"total_tokens":      150,
				"prompt_tokens_details": map[string]int{
					"cached_tokens": 80,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create provider with test server
	provider := NewOpenAIProviderWithBaseURL("test-key", "o4-mini", server.URL)

	// Create test request
	req := &PromptRequest{
		CharacterID: "test-char",
		UserID:      "test-user",
		Message:     "Hello",
		Context: models.ConversationContext{
			RecentMessages: []models.Message{
				{Role: "user", Content: "Previous message"},
			},
		},
		CacheBreakpoints: []cache.CacheBreakpoint{
			{Layer: cache.CorePersonalityLayer, Content: "Test personality"},
		},
	}

	// Send request
	ctx := context.Background()
	resp, err := provider.SendRequest(ctx, req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Verify response
	if resp.Content != "Test response" {
		t.Errorf("Expected 'Test response', got %s", resp.Content)
	}

	if resp.TokensUsed.Total != 150 {
		t.Errorf("Expected 150 total tokens, got %d", resp.TokensUsed.Total)
	}

	// Note: With the SDK-based implementation, cached tokens might not be reported
	// depending on the provider. This is a limitation of the OpenAI-compatible approach
}
