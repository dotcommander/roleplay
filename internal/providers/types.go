package providers

import (
	"context"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
)

// AIProvider defines the interface for AI service providers
type AIProvider interface {
	SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error)
	SupportsBreakpoints() bool
	MaxBreakpoints() int
	Name() string
}

// PromptRequest represents a request to an AI provider
type PromptRequest struct {
	CharacterID      string
	UserID           string
	Message          string
	Context          models.ConversationContext
	SystemPrompt     string
	CacheBreakpoints []cache.CacheBreakpoint
}

// AIResponse represents a response from an AI provider
type AIResponse struct {
	Content      string
	TokensUsed   TokenUsage
	CacheMetrics cache.CacheMetrics
	Emotions     models.EmotionalState
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	Prompt       int
	Completion   int
	CachedPrompt int
	Total        int
}