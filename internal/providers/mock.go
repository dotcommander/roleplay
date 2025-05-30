package providers

import (
	"context"
	"sync"

	"github.com/dotcommander/roleplay/internal/cache"
)

var (
	globalMock *MockProvider
	globalMu   sync.RWMutex
)

// MockProvider is a mock AI provider for testing
type MockProvider struct {
	mu                sync.RWMutex
	responses         []string
	responseIndex     int
	shouldError       bool
	errorToReturn     error
	lastRequest       *PromptRequest
	requestCount      int
	returnCacheHit    bool
	cacheHitPercentage float64
}

// NewMockProvider creates a new mock provider
func NewMockProvider() *MockProvider {
	globalMu.Lock()
	defer globalMu.Unlock()
	
	if globalMock == nil {
		globalMock = &MockProvider{
			responses: []string{"Mock response"},
		}
	}
	return globalMock
}

// SetResponses sets the responses the mock provider will return
func (m *MockProvider) SetResponses(responses []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = responses
	m.responseIndex = 0
}

// SetError configures the provider to return an error
func (m *MockProvider) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldError = true
	m.errorToReturn = err
}

// SetCacheHit configures the provider to simulate cache hits
func (m *MockProvider) SetCacheHit(hit bool, percentage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.returnCacheHit = hit
	m.cacheHitPercentage = percentage
}

// GetLastRequest returns the last request made to the provider
func (m *MockProvider) GetLastRequest() *PromptRequest {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastRequest
}

// GetRequestCount returns the number of requests made
func (m *MockProvider) GetRequestCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.requestCount
}

// Reset resets the mock provider state
func (m *MockProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responseIndex = 0
	m.shouldError = false
	m.errorToReturn = nil
	m.lastRequest = nil
	m.requestCount = 0
	m.returnCacheHit = false
	m.cacheHitPercentage = 0
}

// SendRequest implements AIProvider interface
func (m *MockProvider) SendRequest(ctx context.Context, req *PromptRequest) (*AIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.lastRequest = req
	m.requestCount++
	
	if m.shouldError {
		return nil, m.errorToReturn
	}
	
	// Get response
	response := "Mock response"
	if len(m.responses) > 0 {
		response = m.responses[m.responseIndex%len(m.responses)]
		m.responseIndex++
	}
	
	// Calculate token usage
	promptTokens := len(req.Message) / 4  // Rough token estimate
	completionTokens := len(response) / 4
	cachedTokens := 0
	
	if m.returnCacheHit {
		cachedTokens = int(float64(promptTokens) * m.cacheHitPercentage)
	}
	
	return &AIResponse{
		Content: response,
		TokensUsed: TokenUsage{
			Prompt:       promptTokens,
			Completion:   completionTokens,
			Total:        promptTokens + completionTokens,
			CachedPrompt: cachedTokens,
		},
		CacheMetrics: cache.CacheMetrics{
			Hit:         m.returnCacheHit,
			SavedTokens: cachedTokens,
		},
	}, nil
}

// SendStreamRequest implements AIProvider interface
func (m *MockProvider) SendStreamRequest(ctx context.Context, req *PromptRequest, out chan<- PartialAIResponse) error {
	defer close(out)
	
	// For mock, just send the full response at once
	response, err := m.SendRequest(ctx, req)
	if err != nil {
		return err
	}
	
	select {
	case out <- PartialAIResponse{
		Content: response.Content,
		Done:    true,
	}:
	case <-ctx.Done():
		return ctx.Err()
	}
	
	return nil
}

// Name implements AIProvider interface
func (m *MockProvider) Name() string {
	return "mock"
}

// Global functions for configuring mock provider from tests

// SetGlobalMockResponses sets responses for the global mock provider
func SetGlobalMockResponses(responses []string) {
	globalMu.Lock()
	defer globalMu.Unlock()
	
	if globalMock == nil {
		globalMock = &MockProvider{
			responses: responses,
		}
	} else {
		globalMock.SetResponses(responses)
	}
}

// SetGlobalMockError sets error for the global mock provider
func SetGlobalMockError(err error) {
	globalMu.Lock()
	defer globalMu.Unlock()
	
	if globalMock == nil {
		globalMock = &MockProvider{
			responses:     []string{"Mock response"},
			shouldError:   true,
			errorToReturn: err,
		}
	} else {
		globalMock.SetError(err)
	}
}

// ResetGlobalMock resets the global mock provider
func ResetGlobalMock() {
	globalMu.Lock()
	defer globalMu.Unlock()
	
	if globalMock != nil {
		globalMock.Reset()
	}
}