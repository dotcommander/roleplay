package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the AIProvider interface for OpenAI-compatible models
type OpenAIProvider struct {
	client  *openai.Client
	model   string
	baseURL string // Store for debug logging
}

// NewOpenAIProvider creates a new OpenAI provider instance
func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return NewOpenAIProviderWithBaseURL(apiKey, model, "")
}

// NewOpenAIProviderWithBaseURL creates a new OpenAI-compatible provider with custom base URL
func NewOpenAIProviderWithBaseURL(apiKey, model, baseURL string) *OpenAIProvider {
	// Log the model being used for debugging
	if strings.HasPrefix(model, "o1-") || strings.HasPrefix(model, "o4-") {
		fmt.Printf("⚠️  Using o1/o4 model: %s (limited parameter support)\n", model)
	}

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		// Trust the user-provided base URL - don't modify it
		// This allows for endpoints like /v1beta (Gemini) or other custom paths
		config.BaseURL = strings.TrimRight(baseURL, "/")
		
		// Add debug transport if DEBUG_HTTP env var is set
		if os.Getenv("DEBUG_HTTP") == "true" {
			fmt.Printf("🔧 Debug: OpenAI provider configured with base URL: %s\n", config.BaseURL)
			config.HTTPClient = &http.Client{
				Transport: &debugTransport{RoundTripper: http.DefaultTransport},
			}
		}
	}

	return &OpenAIProvider{
		client:  openai.NewClientWithConfig(config),
		model:   model,
		baseURL: config.BaseURL,
	}
}

// SendRequest sends a request to the OpenAI-compatible API
func (o *OpenAIProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
	// Build messages from the request
	messages := o.buildMessages(req)

	// Create the API request
	apiReq := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
		User:     req.UserID, // Add user parameter for better cache routing
	}

	// o1 models have restrictions on parameters
	if !strings.HasPrefix(o.model, "o1-") && !strings.HasPrefix(o.model, "o4-") {
		// Standard models support these parameters
		apiReq.Temperature = 0.7
		// Use more tokens for user profile updates which can be lengthy JSON
		if req.CharacterID == "system-user-profiler" {
			apiReq.MaxTokens = 4000
		} else {
			apiReq.MaxTokens = 2000
		}
	}

	// Send the request
	if os.Getenv("DEBUG_HTTP") == "true" {
		fmt.Printf("🔧 Debug: Sending chat completion request with model: %s\n", o.model)
	}
	resp, err := o.client.CreateChatCompletion(ctx, apiReq)
	if err != nil {
		// Add debug info if enabled
		if os.Getenv("DEBUG_HTTP") == "true" {
			fmt.Printf("🔧 Debug: Request failed - Error: %v\n", err)
			if o.baseURL != "" {
				fmt.Printf("🔧 Debug: Expected URL: %s/chat/completions\n", o.baseURL)
			}
		}
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	// Debug response if enabled
	DebugResponse(resp)

	// Parse the response
	return o.parseResponse(resp)
}

func (o *OpenAIProvider) buildMessages(req *PromptRequest) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	// Use SystemPrompt if provided (bot service assembles this from cache breakpoints)
	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	} else {
		// Fallback: Combine all breakpoints into system message
		systemContent := ""
		for _, bp := range req.CacheBreakpoints {
			systemContent += bp.Content + "\n\n"
		}
		if systemContent != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: strings.TrimSpace(systemContent),
			})
		}
	}

	// Add conversation history
	for _, msg := range req.Context.RecentMessages {
		role := openai.ChatMessageRoleUser
		if msg.Role == "assistant" || msg.Role == "character" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add current message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Message,
	})

	return messages
}

// SendStreamRequest sends a streaming request to the OpenAI-compatible API
func (o *OpenAIProvider) SendStreamRequest(ctx context.Context, req *PromptRequest, out chan<- PartialAIResponse) error {
	defer close(out)

	// Build messages from the request
	messages := o.buildMessages(req)

	// Create the API request
	apiReq := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
		User:     req.UserID, // Add user parameter for better cache routing
		Stream:   true,
	}

	// o1 models have restrictions on parameters
	if !strings.HasPrefix(o.model, "o1-") && !strings.HasPrefix(o.model, "o4-") {
		// Standard models support these parameters
		apiReq.Temperature = 0.7
		// Use more tokens for user profile updates which can be lengthy JSON
		if req.CharacterID == "system-user-profiler" {
			apiReq.MaxTokens = 4000
		} else {
			apiReq.MaxTokens = 2000
		}
	}

	// Create the stream
	stream, err := o.client.CreateChatCompletionStream(ctx, apiReq)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	// Process stream chunks
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// Stream finished
			out <- PartialAIResponse{
				Done: true,
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}

		// Extract content from the chunk
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			out <- PartialAIResponse{
				Content: response.Choices[0].Delta.Content,
				Done:    false,
			}
		}
	}
}

func (o *OpenAIProvider) parseResponse(resp openai.ChatCompletionResponse) (*AIResponse, error) {
	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
	}

	// Extract cached tokens from prompt_tokens_details if available
	// This field is optional and may not be present in all OpenAI-compatible APIs
	cachedTokens := 0
	cachedLayers := []cache.CacheLayer{}
	
	// Safely check for cached tokens - not all providers return this
	if resp.Usage.PromptTokensDetails != nil {
		cachedTokens = resp.Usage.PromptTokensDetails.CachedTokens
		
		// Debug cached tokens if enabled
		DebugCachedTokens(resp.Usage.PromptTokens, cachedTokens)
	}

	// Determine which cache layers were hit based on cached token count
	// OpenAI caches in 128-token increments starting at 1024 tokens
	if cachedTokens >= 1024 {
		// Core character system prompt is cached (Layer 2)
		cachedLayers = append(cachedLayers, cache.CorePersonalityLayer)
		
		// If significantly more tokens are cached, likely includes user profile (Layer 3)
		if cachedTokens >= 1536 {
			cachedLayers = append(cachedLayers, cache.UserMemoryLayer)
		}
	}

	// Build token usage, handling potential nil or missing fields
	tokenUsage := TokenUsage{
		Prompt:       resp.Usage.PromptTokens,
		Completion:   resp.Usage.CompletionTokens,
		CachedPrompt: cachedTokens,
		Total:        resp.Usage.TotalTokens,
	}
	
	// Calculate saved tokens based on OpenAI's pricing model
	// OpenAI offers 50% discount on cached tokens
	savedTokens := 0
	if cachedTokens > 0 {
		savedTokens = cachedTokens / 2
	}

	return &AIResponse{
		Content:    content,
		TokensUsed: tokenUsage,
		CacheMetrics: cache.CacheMetrics{
			Hit:         cachedTokens > 0,
			Layers:      cachedLayers,
			SavedTokens: savedTokens,
		},
	}, nil
}

// SupportsBreakpoints indicates that OpenAI-compatible APIs handle caching server-side
func (o *OpenAIProvider) SupportsBreakpoints() bool { return false }

// MaxBreakpoints returns 0 as caching is handled server-side
func (o *OpenAIProvider) MaxBreakpoints() int { return 0 }

// Name returns the provider name
func (o *OpenAIProvider) Name() string { return "openai_compatible" }

// debugTransport is an HTTP transport that logs all requests
type debugTransport struct {
	http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request URL
	fmt.Printf("🔧 Debug: HTTP %s %s\n", req.Method, req.URL.String())
	
	// Make the actual request
	resp, err := d.RoundTripper.RoundTrip(req)
	
	// Log response status
	if err != nil {
		fmt.Printf("🔧 Debug: Request failed with error: %v\n", err)
	} else {
		fmt.Printf("🔧 Debug: Response status: %d %s\n", resp.StatusCode, resp.Status)
	}
	
	return resp, err
}
