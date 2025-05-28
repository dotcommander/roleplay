package config

import "time"

// Config holds all application configuration
type Config struct {
	DefaultProvider   string
	Model             string
	APIKey            string
	BaseURL           string            // OpenAI-compatible endpoint
	ModelAliases      map[string]string // Aliases for models
	CacheConfig       CacheConfig
	MemoryConfig      MemoryConfig
	PersonalityConfig PersonalityConfig
	UserProfileConfig UserProfileConfig
}

// CacheConfig holds cache-related configuration
type CacheConfig struct {
	MaxEntries        int
	CleanupInterval   time.Duration
	DefaultTTL        time.Duration
	EnableAdaptiveTTL bool
}

// MemoryConfig holds memory management configuration
type MemoryConfig struct {
	ShortTermWindow    int           // Number of messages
	MediumTermDuration time.Duration // How long to keep
	ConsolidationRate  float64       // Learning rate for personality evolution
}

// PersonalityConfig holds personality evolution configuration
type PersonalityConfig struct {
	EvolutionEnabled   bool
	MaxDriftRate       float64 // Maximum personality change per interaction
	StabilityThreshold float64 // Minimum interactions before evolution
}

// UserProfileConfig holds user profile agent configuration
type UserProfileConfig struct {
	Enabled              bool          `mapstructure:"enabled"`
	UpdateFrequency      int           `mapstructure:"update_frequency_messages"` // Update every N messages
	TurnsToConsider      int           `mapstructure:"turns_to_consider"`         // How many past turns to analyze
	ConfidenceThreshold  float64       `mapstructure:"confidence_threshold"`      // Min confidence for facts
	PromptCacheTTL       time.Duration `mapstructure:"prompt_cache_ttl"`
}
