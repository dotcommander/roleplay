package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dotcommander/roleplay/internal/cache"
)

// OpenAIProvider implements the AIProvider interface for OpenAI models
type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
}

// NewOpenAIProvider creates a new OpenAI provider instance
func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return NewOpenAIProviderWithBaseURL(apiKey, model, "https://api.openai.com/v1")
}

// NewOpenAIProviderWithBaseURL creates a new OpenAI provider with custom base URL
func NewOpenAIProviderWithBaseURL(apiKey, model, baseURL string) *OpenAIProvider {
	// Log the model being used for debugging
	if strings.HasPrefix(model, "o1-") || strings.HasPrefix(model, "o4-") {
		fmt.Printf("⚠️  Using o1/o4 model: %s (limited parameter support)\n", model)
	}
	
	// Ensure baseURL doesn't have trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &OpenAIProvider{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		model:      model,
	}
}

// SendRequest sends a request to the OpenAI API
func (o *OpenAIProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
	// OpenAI uses automatic caching, so we just need to structure prompts consistently
	messages := o.buildMessages(req)

	payload := map[string]interface{}{
		"model":    o.model,
		"messages": messages,
	}

	// o1 models have restrictions on parameters
	if strings.HasPrefix(o.model, "o1-") || strings.HasPrefix(o.model, "o4-") {
		// o1 models don't support temperature or max_tokens
		// They use default values
	} else {
		// Standard models support these parameters
		payload["temperature"] = 0.7
		payload["max_tokens"] = 2000
	}

	respData, err := o.makeRequest(ctx, "/chat/completions", payload)
	if err != nil {
		return nil, err
	}

	return o.parseResponse(respData)
}

func (o *OpenAIProvider) buildMessages(req *PromptRequest) []map[string]string {
	messages := []map[string]string{}

	// Combine all breakpoints into system message for consistent caching
	systemContent := ""
	for _, bp := range req.CacheBreakpoints {
		systemContent += bp.Content + "\n\n"
	}

	if systemContent != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": systemContent,
		})
	}

	// Add conversation history
	for _, msg := range req.Context.RecentMessages {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// Add current message
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.Message,
	})

	return messages
}

func (o *OpenAIProvider) makeRequest(ctx context.Context, endpoint string, payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
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

func (o *OpenAIProvider) parseResponse(data []byte) (*AIResponse, error) {
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
	}

	// Determine cached layers
	cachedLayers := []cache.CacheLayer{}
	if resp.Usage.PromptTokensDetails.CachedTokens > 0 {
		// OpenAI's automatic caching likely cached the system prompt
		cachedLayers = append(cachedLayers, cache.CorePersonalityLayer)
	}

	return &AIResponse{
		Content: content,
		TokensUsed: TokenUsage{
			Prompt:       resp.Usage.PromptTokens,
			Completion:   resp.Usage.CompletionTokens,
			CachedPrompt: resp.Usage.PromptTokensDetails.CachedTokens,
			Total:        resp.Usage.TotalTokens,
		},
		CacheMetrics: cache.CacheMetrics{
			Hit:         resp.Usage.PromptTokensDetails.CachedTokens > 0,
			Layers:      cachedLayers,
			SavedTokens: resp.Usage.PromptTokensDetails.CachedTokens / 2, // 50% discount
		},
	}, nil
}

// SupportsBreakpoints indicates that OpenAI uses automatic caching
func (o *OpenAIProvider) SupportsBreakpoints() bool { return false }

// MaxBreakpoints returns 0 as OpenAI handles caching automatically
func (o *OpenAIProvider) MaxBreakpoints() int { return 0 }

// Name returns the provider name
func (o *OpenAIProvider) Name() string { return "openai" }
