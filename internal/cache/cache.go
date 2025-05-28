package cache

import (
	"sync"
	"time"
)

// PromptCache manages cached prompts with TTL
type PromptCache struct {
	entries map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     TTLManager
}

// NewPromptCache creates a new cache with the given TTL configuration
func NewPromptCache(baseTTL, minTTL, maxTTL time.Duration) *PromptCache {
	return &PromptCache{
		entries: make(map[string]*CacheEntry),
		ttl: TTLManager{
			BaseTTL:         baseTTL,
			ActiveBonus:     0.5,
			ComplexityBonus: 0.2,
			MinTTL:          minTTL,
			MaxTTL:          maxTTL,
		},
	}
}

// Store adds a new cache entry for a specific layer
func (pc *PromptCache) Store(key string, layer CacheLayer, content string, ttl time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry, exists := pc.entries[key]
	if !exists {
		entry = &CacheEntry{
			CreatedAt:   time.Now(),
			Breakpoints: make([]CacheBreakpoint, 0),
		}
		pc.entries[key] = entry
	}

	breakpoint := CacheBreakpoint{
		Layer:      layer,
		Content:    content,
		TokenCount: EstimateTokens(content),
		TTL:        ttl,
		LastUsed:   time.Now(),
	}

	entry.Breakpoints = append(entry.Breakpoints, breakpoint)
	entry.LastAccess = time.Now()
}

// StoreWithTTL stores a complete cache entry with breakpoints
func (pc *PromptCache) StoreWithTTL(key string, breakpoints []CacheBreakpoint, ttl time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry := &CacheEntry{
		Breakpoints: breakpoints,
		CreatedAt:   time.Now(),
		LastAccess:  time.Now(),
		HitCount:    0,
	}

	pc.entries[key] = entry
}

// Get retrieves a cache entry if it exists and is not expired
func (pc *PromptCache) Get(key string) (*CacheEntry, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	entry, exists := pc.entries[key]
	if !exists {
		return nil, false
	}

	// Update access time and hit count
	entry.LastAccess = time.Now()
	entry.HitCount++

	return entry, true
}

// CalculateAdaptiveTTL determines the effective TTL based on usage patterns
func (pc *PromptCache) CalculateAdaptiveTTL(cached *CacheEntry, hasComplexCharacter bool) time.Duration {
	baseTTL := pc.ttl.BaseTTL

	// Active conversation bonus
	if cached != nil && time.Since(cached.LastAccess) < 5*time.Minute {
		baseTTL = time.Duration(float64(baseTTL) * (1 + pc.ttl.ActiveBonus))
	}

	// Character complexity bonus
	if hasComplexCharacter {
		baseTTL = time.Duration(float64(baseTTL) * (1 + pc.ttl.ComplexityBonus))
	}

	// Enforce limits
	if baseTTL < pc.ttl.MinTTL {
		baseTTL = pc.ttl.MinTTL
	}
	if baseTTL > pc.ttl.MaxTTL {
		baseTTL = pc.ttl.MaxTTL
	}

	return baseTTL
}

// CleanupWorker runs periodic cleanup of expired entries
func (pc *PromptCache) CleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		pc.cleanup()
	}
}

// cleanup removes expired cache entries
func (pc *PromptCache) cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for key, entry := range pc.entries {
		// Check if any breakpoint has expired
		expired := false
		for _, bp := range entry.Breakpoints {
			if now.Sub(bp.LastUsed) > bp.TTL {
				expired = true
				break
			}
		}

		if expired {
			delete(pc.entries, key)
		}
	}
}

// EstimateTokens provides a rough estimation of token count
func EstimateTokens(text string) int {
	// Rough estimation: ~4 chars per token
	return len(text) / 4
}