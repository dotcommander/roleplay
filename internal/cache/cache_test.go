package cache

import (
	"testing"
	"time"
)

func TestPromptCache(t *testing.T) {
	cache := NewPromptCache(5*time.Minute, 1*time.Minute, 10*time.Minute)

	// Test storing and retrieving
	cache.Store("test-key", CorePersonalityLayer, "test content", 5*time.Minute)

	entry, found := cache.Get("test-key")
	if !found {
		t.Fatal("Expected to find cached entry")
	}

	if len(entry.Breakpoints) != 1 {
		t.Errorf("Expected 1 breakpoint, got %d", len(entry.Breakpoints))
	}

	if entry.Breakpoints[0].Layer != CorePersonalityLayer {
		t.Errorf("Expected CorePersonalityLayer, got %s", entry.Breakpoints[0].Layer)
	}

	// Test hit count
	initialHits := entry.HitCount
	cache.Get("test-key")
	entry, _ = cache.Get("test-key")
	if entry.HitCount != initialHits+2 {
		t.Errorf("Expected hit count to increase, got %d", entry.HitCount)
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"Hello world", 2},    // 11 chars / 4 ≈ 2
		{"This is a test", 3}, // 14 chars / 4 ≈ 3
		{"", 0},               // Empty string
		{"A", 0},              // 1 char / 4 = 0
		{"1234", 1},           // 4 chars / 4 = 1
	}

	for _, tt := range tests {
		result := EstimateTokens(tt.text)
		if result != tt.expected {
			t.Errorf("EstimateTokens(%q) = %d, want %d", tt.text, result, tt.expected)
		}
	}
}

func TestAdaptiveTTL(t *testing.T) {
	cache := NewPromptCache(5*time.Minute, 1*time.Minute, 10*time.Minute)

	// Test without cached entry
	ttl := cache.CalculateAdaptiveTTL(nil, false)
	if ttl != 5*time.Minute {
		t.Errorf("Expected base TTL of 5m, got %v", ttl)
	}

	// Test with complexity bonus
	ttl = cache.CalculateAdaptiveTTL(nil, true)
	expectedTTL := time.Duration(float64(5*time.Minute) * 1.2) // 20% bonus
	if ttl != expectedTTL {
		t.Errorf("Expected TTL with complexity bonus of %v, got %v", expectedTTL, ttl)
	}

	// Test with recent access
	entry := &CacheEntry{
		LastAccess: time.Now(),
		HitCount:   5,
	}
	ttl = cache.CalculateAdaptiveTTL(entry, false)
	expectedTTL = time.Duration(float64(5*time.Minute) * 1.5) // 50% bonus
	if ttl != expectedTTL {
		t.Errorf("Expected TTL with active bonus of %v, got %v", expectedTTL, ttl)
	}

	// Test max TTL enforcement
	cache.ttl.BaseTTL = 20 * time.Minute
	ttl = cache.CalculateAdaptiveTTL(nil, false)
	if ttl != cache.ttl.MaxTTL {
		t.Errorf("Expected TTL to be capped at MaxTTL %v, got %v", cache.ttl.MaxTTL, ttl)
	}
}

func TestCacheCleanup(t *testing.T) {
	cache := NewPromptCache(100*time.Millisecond, 50*time.Millisecond, 200*time.Millisecond)

	// Store entry with short TTL
	breakpoints := []CacheBreakpoint{
		{
			Layer:    CorePersonalityLayer,
			Content:  "test",
			TTL:      100 * time.Millisecond,
			LastUsed: time.Now(),
		},
	}
	cache.StoreWithTTL("expire-key", breakpoints, 100*time.Millisecond)

	// Verify entry exists
	_, found := cache.Get("expire-key")
	if !found {
		t.Fatal("Expected to find entry before expiration")
	}

	// Wait for expiration and cleanup
	time.Sleep(150 * time.Millisecond)
	cache.cleanup()

	// Verify entry was cleaned up
	_, found = cache.Get("expire-key")
	if found {
		t.Error("Expected entry to be cleaned up after expiration")
	}
}

func TestCacheLayers(t *testing.T) {
	layers := []CacheLayer{
		CorePersonalityLayer,
		LearnedBehaviorLayer,
		EmotionalStateLayer,
		ConversationLayer,
	}

	// Verify all layers are distinct
	seen := make(map[CacheLayer]bool)
	for _, layer := range layers {
		if seen[layer] {
			t.Errorf("Duplicate layer found: %s", layer)
		}
		seen[layer] = true
	}
}
