package cache

import (
	"time"
)

// CacheLayer represents different cache layers
type CacheLayer string

const (
	CorePersonalityLayer CacheLayer = "core_personality"
	LearnedBehaviorLayer CacheLayer = "learned_behavior"
	EmotionalStateLayer  CacheLayer = "emotional_state"
	UserMemoryLayer      CacheLayer = "user_memory"
	ConversationLayer    CacheLayer = "conversation"
)

// CacheBreakpoint represents a cache checkpoint
type CacheBreakpoint struct {
	Layer      CacheLayer    `json:"layer"`
	Content    string        `json:"content"`
	TokenCount int           `json:"token_count"`
	TTL        time.Duration `json:"ttl"`
	LastUsed   time.Time     `json:"last_used"`
}

// CacheEntry represents a cached prompt entry
type CacheEntry struct {
	Breakpoints []CacheBreakpoint
	Hash        string
	CreatedAt   time.Time
	LastAccess  time.Time
	HitCount    int
	UserID      string
}

// TTLManager handles dynamic TTL calculations
type TTLManager struct {
	BaseTTL         time.Duration
	ActiveBonus     float64 // 50% bonus for active conversations
	ComplexityBonus float64 // 20% bonus for complex characters
	MinTTL          time.Duration
	MaxTTL          time.Duration
}

// CacheMetrics tracks cache performance
type CacheMetrics struct {
	Hit         bool
	Layers      []CacheLayer
	SavedTokens int
	Latency     time.Duration
}