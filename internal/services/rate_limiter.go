package services

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter manages request rates per user-character pair to optimize cache routing
type RateLimiter struct {
	mu        sync.RWMutex
	buckets   map[string]*bucket
	maxRate   int           // Max requests per window
	window    time.Duration // Time window (e.g., 1 minute)
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// bucket tracks requests for a specific user-character pair
type bucket struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:         make(map[string]*bucket),
		maxRate:         maxRate,
		window:          window,
		cleanupInterval: window * 2, // Cleanup stale buckets every 2 windows
		stopCleanup:     make(chan struct{}),
	}
	
	// Start background cleanup
	go rl.cleanup()
	
	return rl
}

// Allow checks if a request should be allowed for the given user-character pair
func (rl *RateLimiter) Allow(userID, characterID string) (bool, error) {
	key := fmt.Sprintf("%s:%s", userID, characterID)
	
	rl.mu.RLock()
	b, exists := rl.buckets[key]
	rl.mu.RUnlock()
	
	if !exists {
		rl.mu.Lock()
		b = &bucket{requests: make([]time.Time, 0, rl.maxRate)}
		rl.buckets[key] = b
		rl.mu.Unlock()
	}
	
	b.mu.Lock()
	defer b.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	// Remove expired requests
	validRequests := make([]time.Time, 0, len(b.requests))
	for _, t := range b.requests {
		if t.After(cutoff) {
			validRequests = append(validRequests, t)
		}
	}
	b.requests = validRequests
	
	// Check if we're under the limit
	if len(b.requests) >= rl.maxRate {
		return false, fmt.Errorf("rate limit exceeded: %d requests in %v for user %s with character %s", 
			len(b.requests), rl.window, userID, characterID)
	}
	
	// Add current request
	b.requests = append(b.requests, now)
	return true, nil
}

// GetCurrentRate returns the current request rate for a user-character pair
func (rl *RateLimiter) GetCurrentRate(userID, characterID string) int {
	key := fmt.Sprintf("%s:%s", userID, characterID)
	
	rl.mu.RLock()
	b, exists := rl.buckets[key]
	rl.mu.RUnlock()
	
	if !exists {
		return 0
	}
	
	b.mu.Lock()
	defer b.mu.Unlock()
	
	cutoff := time.Now().Add(-rl.window)
	count := 0
	for _, t := range b.requests {
		if t.After(cutoff) {
			count++
		}
	}
	
	return count
}

// cleanup removes stale buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			rl.cleanupStale()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanupStale removes buckets that haven't been used recently
func (rl *RateLimiter) cleanupStale() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	cutoff := time.Now().Add(-rl.window * 2)
	
	for key, b := range rl.buckets {
		b.mu.Lock()
		if len(b.requests) == 0 || (len(b.requests) > 0 && b.requests[len(b.requests)-1].Before(cutoff)) {
			delete(rl.buckets, key)
		}
		b.mu.Unlock()
	}
}

// Stop gracefully stops the rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// GetStats returns statistics about current rate limiting
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	
	stats := map[string]interface{}{
		"total_buckets": len(rl.buckets),
		"max_rate":      rl.maxRate,
		"window":        rl.window.String(),
		"active_pairs":  make([]map[string]interface{}, 0),
	}
	
	for key, b := range rl.buckets {
		b.mu.Lock()
		cutoff := time.Now().Add(-rl.window)
		activeCount := 0
		for _, t := range b.requests {
			if t.After(cutoff) {
				activeCount++
			}
		}
		b.mu.Unlock()
		
		if activeCount > 0 {
			stats["active_pairs"] = append(stats["active_pairs"].([]map[string]interface{}), map[string]interface{}{
				"key":   key,
				"count": activeCount,
				"rate":  fmt.Sprintf("%.1f%%", float64(activeCount)/float64(rl.maxRate)*100),
			})
		}
	}
	
	return stats
}