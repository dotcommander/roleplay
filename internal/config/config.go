package config

import "time"

// Config holds all application configuration
type Config struct {
	DefaultProvider   string
	Model            string
	APIKey           string
	CacheConfig       CacheConfig
	MemoryConfig      MemoryConfig
	PersonalityConfig PersonalityConfig
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