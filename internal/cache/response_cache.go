package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// ResponseCache caches complete API responses
type ResponseCache struct {
	responses map[string]*CachedResponse
	mu        sync.RWMutex
	ttl       time.Duration
}

// CachedResponse represents a cached API response
type CachedResponse struct {
	Content    string
	TokensUsed TokenUsage
	CachedAt   time.Time
	ExpiresAt  time.Time
	HitCount   int
}

// TokenUsage represents token usage stats
type TokenUsage struct {
	Prompt       int
	Completion   int
	CachedPrompt int
	Total        int
}

// NewResponseCache creates a new response cache
func NewResponseCache(ttl time.Duration) *ResponseCache {
	cache := &ResponseCache{
		responses: make(map[string]*CachedResponse),
		ttl:       ttl,
	}

	// Start cleanup worker
	go cache.cleanupWorker()

	return cache
}

// GenerateKey creates a cache key from request parameters
func (rc *ResponseCache) GenerateKey(characterID, userID, message string) string {
	h := sha256.New()
	h.Write([]byte(characterID + "|" + userID + "|" + message))
	return hex.EncodeToString(h.Sum(nil))
}

// Get retrieves a cached response if available
func (rc *ResponseCache) Get(key string) (*CachedResponse, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	resp, exists := rc.responses[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(resp.ExpiresAt) {
		return nil, false
	}

	// Update hit count
	resp.HitCount++

	return resp, true
}

// Store adds a response to the cache
func (rc *ResponseCache) Store(key, content string, tokens TokenUsage) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.responses[key] = &CachedResponse{
		Content:    content,
		TokensUsed: tokens,
		CachedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(rc.ttl),
		HitCount:   0,
	}
}

// cleanupWorker removes expired entries
func (rc *ResponseCache) cleanupWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.mu.Lock()
		now := time.Now()
		for key, resp := range rc.responses {
			if now.After(resp.ExpiresAt) {
				delete(rc.responses, key)
			}
		}
		rc.mu.Unlock()
	}
}

// GetStats returns cache statistics
func (rc *ResponseCache) GetStats() (hits, misses int) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	for _, resp := range rc.responses {
		hits += resp.HitCount
	}

	return hits, 0 // Misses would need to be tracked separately
}
