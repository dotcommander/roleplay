package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dotcommander/roleplay/internal/cache"
)

// AnthropicProvider implements the AIProvider interface for Claude
type AnthropicProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
	version    string
}

// NewAnthropicProvider creates a new Anthropic provider instance
func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com/v1",
		httpClient: &http.Client{Timeout: 60 * time.Second},
		model:      "claude-3-opus-20240229",
		version:    "2024-01-01",
	}
}

// SendRequest sends a request to the Anthropic API
func (a *AnthropicProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
	// Build the system prompt from cacheable layers
	systemPrompt := ""
	cacheableBreakpoints := make([]cache.CacheBreakpoint, 0)
	
	// Separate cacheable and non-cacheable content
	for _, bp := range req.CacheBreakpoints {
		if bp.Layer != cache.ConversationLayer {
			cacheableBreakpoints = append(cacheableBreakpoints, bp)
			if systemPrompt != "" {
				systemPrompt += "\n\n"
			}
			systemPrompt += bp.Content
		}
	}
	
	// Build messages with cache control
	messages := a.buildMessagesWithCache(req)

	payload := map[string]interface{}{
		"model":       a.model,
		"messages":    messages,
		"max_tokens":  2000,
		"temperature": 0.7,
	}

	// Add cache control to system prompt if we have cacheable content
	if systemPrompt != "" {
		payload["system"] = []map[string]interface{}{
			{
				"type": "text",
				"text": systemPrompt,
				"cache_control": map[string]string{"type": "ephemeral"},
			},
		}
	}

	// Add beta header for prompt caching
	headers := map[string]string{
		"anthropic-beta":    "prompt-caching-2024-07-31",
		"anthropic-version": a.version,
		"content-type":      "application/json",
		"x-api-key":         a.apiKey,
	}

	// Make request
	respData, err := a.makeRequestWithHeaders(ctx, "/messages", payload, headers)
	if err != nil {
		return nil, err
	}

	// Parse response
	return a.parseResponse(respData)
}

func (a *AnthropicProvider) buildMessagesWithCache(req *PromptRequest) []map[string]interface{} {
	messages := make([]map[string]interface{}, 0)

	// Add conversation history from breakpoints (if any)
	for _, bp := range req.CacheBreakpoints {
		if bp.Layer == cache.ConversationLayer && bp.Content != "" {
			// Parse conversation history and add as messages
			for _, msg := range req.Context.RecentMessages {
				messages = append(messages, map[string]interface{}{
					"role":    msg.Role,
					"content": msg.Content,
				})
			}
			break
		}
	}

	// Add current user message
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": req.Message,
	})

	return messages
}

func (a *AnthropicProvider) makeRequestWithHeaders(ctx context.Context, endpoint string, payload interface{}, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (a *AnthropicProvider) parseResponse(data []byte) (*AIResponse, error) {
	var resp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	// Extract content
	content := ""
	for _, c := range resp.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	// Calculate cache metrics
	cacheHit := resp.Usage.CacheReadInputTokens > 0
	savedTokens := resp.Usage.CacheReadInputTokens
	
	// Determine which layers were cached based on token counts
	cachedLayers := []cache.CacheLayer{}
	if cacheHit {
		// If we have cached tokens, assume at least personality layer was cached
		cachedLayers = append(cachedLayers, cache.CorePersonalityLayer)
		// Additional heuristics could be added here based on token counts
	}

	return &AIResponse{
		Content: content,
		TokensUsed: TokenUsage{
			Prompt:       resp.Usage.InputTokens,
			Completion:   resp.Usage.OutputTokens,
			CachedPrompt: resp.Usage.CacheReadInputTokens,
			Total:        resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
		CacheMetrics: cache.CacheMetrics{
			Hit:         cacheHit,
			Layers:      cachedLayers,
			SavedTokens: savedTokens,
		},
	}, nil
}

// SupportsBreakpoints indicates that Anthropic supports cache breakpoints
func (a *AnthropicProvider) SupportsBreakpoints() bool { return true }

// MaxBreakpoints returns the maximum number of breakpoints supported
func (a *AnthropicProvider) MaxBreakpoints() int { return 4 }

// Name returns the provider name
func (a *AnthropicProvider) Name() string { return "anthropic" }