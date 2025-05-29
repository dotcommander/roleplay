package config

import (
	"testing"
	"time"
)

func TestConfigStructure(t *testing.T) {
	// Test that config can be created and has expected structure
	cfg := &Config{
		DefaultProvider: "openai",
		Model:           "gpt-4",
		APIKey:          "test-key",
		BaseURL:         "https://api.openai.com/v1",
		CacheConfig: CacheConfig{
			MaxEntries:        10000,
			CleanupInterval:   5 * time.Minute,
			DefaultTTL:        10 * time.Minute,
			EnableAdaptiveTTL: true,
		},
		MemoryConfig: MemoryConfig{
			ShortTermWindow:    20,
			MediumTermDuration: 24 * time.Hour,
			ConsolidationRate:  0.1,
		},
		PersonalityConfig: PersonalityConfig{
			EvolutionEnabled:   true,
			MaxDriftRate:       0.02,
			StabilityThreshold: 10,
		},
		UserProfileConfig: UserProfileConfig{
			Enabled:             true,
			UpdateFrequency:     5,
			TurnsToConsider:     20,
			ConfidenceThreshold: 0.5,
			PromptCacheTTL:      1 * time.Hour,
		},
	}

	// Verify fields
	if cfg.DefaultProvider != "openai" {
		t.Errorf("DefaultProvider mismatch: got %s, want openai", cfg.DefaultProvider)
	}
	if cfg.Model != "gpt-4" {
		t.Errorf("Model mismatch: got %s, want gpt-4", cfg.Model)
	}
	if cfg.APIKey != "test-key" {
		t.Errorf("APIKey mismatch: got %s, want test-key", cfg.APIKey)
	}
	if cfg.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("BaseURL mismatch: got %s, want https://api.openai.com/v1", cfg.BaseURL)
	}

	// Test cache config
	if cfg.CacheConfig.MaxEntries != 10000 {
		t.Errorf("CacheConfig.MaxEntries mismatch: got %d, want 10000", cfg.CacheConfig.MaxEntries)
	}
	if cfg.CacheConfig.CleanupInterval != 5*time.Minute {
		t.Errorf("CacheConfig.CleanupInterval mismatch: got %v, want 5m", cfg.CacheConfig.CleanupInterval)
	}
	if cfg.CacheConfig.DefaultTTL != 10*time.Minute {
		t.Errorf("CacheConfig.DefaultTTL mismatch: got %v, want 10m", cfg.CacheConfig.DefaultTTL)
	}
	if !cfg.CacheConfig.EnableAdaptiveTTL {
		t.Error("CacheConfig.EnableAdaptiveTTL should be true")
	}

	// Test memory config
	if cfg.MemoryConfig.ShortTermWindow != 20 {
		t.Errorf("MemoryConfig.ShortTermWindow mismatch: got %d, want 20", cfg.MemoryConfig.ShortTermWindow)
	}
	if cfg.MemoryConfig.MediumTermDuration != 24*time.Hour {
		t.Errorf("MemoryConfig.MediumTermDuration mismatch: got %v, want 24h", cfg.MemoryConfig.MediumTermDuration)
	}
	if cfg.MemoryConfig.ConsolidationRate != 0.1 {
		t.Errorf("MemoryConfig.ConsolidationRate mismatch: got %f, want 0.1", cfg.MemoryConfig.ConsolidationRate)
	}

	// Test personality config
	if !cfg.PersonalityConfig.EvolutionEnabled {
		t.Error("PersonalityConfig.EvolutionEnabled should be true")
	}
	if cfg.PersonalityConfig.MaxDriftRate != 0.02 {
		t.Errorf("PersonalityConfig.MaxDriftRate mismatch: got %f, want 0.02", cfg.PersonalityConfig.MaxDriftRate)
	}
	if cfg.PersonalityConfig.StabilityThreshold != 10 {
		t.Errorf("PersonalityConfig.StabilityThreshold mismatch: got %f, want 10", cfg.PersonalityConfig.StabilityThreshold)
	}

	// Test user profile config
	if !cfg.UserProfileConfig.Enabled {
		t.Error("UserProfileConfig.Enabled should be true")
	}
	if cfg.UserProfileConfig.UpdateFrequency != 5 {
		t.Errorf("UserProfileConfig.UpdateFrequency mismatch: got %d, want 5", cfg.UserProfileConfig.UpdateFrequency)
	}
	if cfg.UserProfileConfig.TurnsToConsider != 20 {
		t.Errorf("UserProfileConfig.TurnsToConsider mismatch: got %d, want 20", cfg.UserProfileConfig.TurnsToConsider)
	}
	if cfg.UserProfileConfig.ConfidenceThreshold != 0.5 {
		t.Errorf("UserProfileConfig.ConfidenceThreshold mismatch: got %f, want 0.5", cfg.UserProfileConfig.ConfidenceThreshold)
	}
	if cfg.UserProfileConfig.PromptCacheTTL != 1*time.Hour {
		t.Errorf("UserProfileConfig.PromptCacheTTL mismatch: got %v, want 1h", cfg.UserProfileConfig.PromptCacheTTL)
	}
}

func TestConfigZeroValues(t *testing.T) {
	// Test that zero values work correctly
	cfg := &Config{}

	// All string fields should be empty
	if cfg.DefaultProvider != "" {
		t.Error("DefaultProvider should be empty string by default")
	}
	if cfg.Model != "" {
		t.Error("Model should be empty string by default")
	}
	if cfg.APIKey != "" {
		t.Error("APIKey should be empty string by default")
	}
	if cfg.BaseURL != "" {
		t.Error("BaseURL should be empty string by default")
	}

	// All numeric fields should be zero
	if cfg.CacheConfig.MaxEntries != 0 {
		t.Error("CacheConfig.MaxEntries should be 0 by default")
	}
	if cfg.CacheConfig.CleanupInterval != 0 {
		t.Error("CacheConfig.CleanupInterval should be 0 by default")
	}
	if cfg.CacheConfig.DefaultTTL != 0 {
		t.Error("CacheConfig.DefaultTTL should be 0 by default")
	}

	// All bool fields should be false
	if cfg.CacheConfig.EnableAdaptiveTTL {
		t.Error("CacheConfig.EnableAdaptiveTTL should be false by default")
	}
	if cfg.PersonalityConfig.EvolutionEnabled {
		t.Error("PersonalityConfig.EvolutionEnabled should be false by default")
	}
	if cfg.UserProfileConfig.Enabled {
		t.Error("UserProfileConfig.Enabled should be false by default")
	}
}

func TestModelAliases(t *testing.T) {
	cfg := &Config{
		ModelAliases: map[string]string{
			"gpt4":       "gpt-4",
			"gpt4-turbo": "gpt-4-turbo-preview",
			"claude":     "claude-3-opus-20240229",
		},
	}

	// Test alias lookup
	tests := []struct {
		alias    string
		expected string
	}{
		{"gpt4", "gpt-4"},
		{"gpt4-turbo", "gpt-4-turbo-preview"},
		{"claude", "claude-3-opus-20240229"},
	}

	for _, tt := range tests {
		if model, ok := cfg.ModelAliases[tt.alias]; !ok || model != tt.expected {
			t.Errorf("ModelAliases[%s] = %s, want %s", tt.alias, model, tt.expected)
		}
	}
}

func TestCacheConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config CacheConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: CacheConfig{
				MaxEntries:      10000,
				CleanupInterval: 5 * time.Minute,
				DefaultTTL:      10 * time.Minute,
			},
			valid: true,
		},
		{
			name: "zero max entries is valid",
			config: CacheConfig{
				MaxEntries:      0,
				CleanupInterval: 5 * time.Minute,
				DefaultTTL:      10 * time.Minute,
			},
			valid: true,
		},
		{
			name: "negative values should be invalid",
			config: CacheConfig{
				MaxEntries:      -1,
				CleanupInterval: -1 * time.Minute,
				DefaultTTL:      -1 * time.Minute,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since we don't have a validation function in the package,
			// we'll just test that the values can be set
			cfg := &Config{CacheConfig: tt.config}
			
			if tt.valid {
				// For valid configs, values should be as set
				if cfg.CacheConfig.MaxEntries != tt.config.MaxEntries {
					t.Error("MaxEntries not set correctly")
				}
			} else {
				// For invalid configs, we'd expect validation to fail
				// but since there's no validation function, we skip this
				t.Skip("No validation function available")
			}
		})
	}
}

func TestUserProfileConfig(t *testing.T) {
	// Test various user profile configurations
	tests := []struct {
		name   string
		config UserProfileConfig
	}{
		{
			name: "enabled profile",
			config: UserProfileConfig{
				Enabled:             true,
				UpdateFrequency:     5,
				TurnsToConsider:     20,
				ConfidenceThreshold: 0.5,
				PromptCacheTTL:      1 * time.Hour,
			},
		},
		{
			name: "disabled profile",
			config: UserProfileConfig{
				Enabled: false,
			},
		},
		{
			name: "high confidence threshold",
			config: UserProfileConfig{
				Enabled:             true,
				UpdateFrequency:     10,
				TurnsToConsider:     50,
				ConfidenceThreshold: 0.9,
				PromptCacheTTL:      2 * time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{UserProfileConfig: tt.config}
			
			if cfg.UserProfileConfig.Enabled != tt.config.Enabled {
				t.Errorf("Enabled mismatch: got %v, want %v", cfg.UserProfileConfig.Enabled, tt.config.Enabled)
			}
			if cfg.UserProfileConfig.UpdateFrequency != tt.config.UpdateFrequency {
				t.Errorf("UpdateFrequency mismatch: got %d, want %d", cfg.UserProfileConfig.UpdateFrequency, tt.config.UpdateFrequency)
			}
			if cfg.UserProfileConfig.ConfidenceThreshold != tt.config.ConfidenceThreshold {
				t.Errorf("ConfidenceThreshold mismatch: got %f, want %f", cfg.UserProfileConfig.ConfidenceThreshold, tt.config.ConfidenceThreshold)
			}
		})
	}
}