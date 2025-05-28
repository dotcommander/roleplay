package factory

import (
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvider(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		envSetup    func()
		envCleanup  func()
		wantErr     bool
		errContains string
	}{
		{
			name: "anthropic provider with API key in config",
			cfg: &config.Config{
				DefaultProvider: "anthropic",
				APIKey:          "test-anthropic-key",
			},
			wantErr: false,
		},
		{
			name: "openai provider with API key and model in config",
			cfg: &config.Config{
				DefaultProvider: "openai",
				APIKey:          "test-openai-key",
				Model:           "gpt-4",
			},
			wantErr: false,
		},
		{
			name: "openai provider with default model",
			cfg: &config.Config{
				DefaultProvider: "openai",
				APIKey:          "test-openai-key",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			cfg: &config.Config{
				DefaultProvider: "openai",
			},
			wantErr:     true,
			errContains: "API key required for openai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envSetup != nil {
				tt.envSetup()
			}
			if tt.envCleanup != nil {
				defer tt.envCleanup()
			}

			provider, err := CreateProvider(tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				// All providers now return "openai_compatible"
				assert.Equal(t, "openai_compatible", provider.Name())
			}
		})
	}
}

func TestInitializeAndRegisterProvider(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "openai",
		APIKey:          "test-key",
		Model:           "gpt-4",
		CacheConfig: config.CacheConfig{
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	bot := services.NewCharacterBot(cfg)

	err := InitializeAndRegisterProvider(bot, cfg)
	assert.NoError(t, err)

	// Test with missing API key
	cfg2 := &config.Config{
		DefaultProvider: "anthropic",
		CacheConfig: config.CacheConfig{
			DefaultTTL:      5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		},
	}

	err = InitializeAndRegisterProvider(bot, cfg2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestCreateProviderWithFallback(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		apiKey      string
		model       string
		baseURL     string
		wantErr     bool
	}{
		{
			name:        "direct API key",
			profileName: "openai",
			apiKey:      "direct-key",
			model:       "gpt-4",
			wantErr:     false,
		},
		{
			name:        "ollama without API key",
			profileName: "ollama",
			apiKey:      "",
			baseURL:     "http://localhost:11434/v1",
			wantErr:     false,
		},
		{
			name:        "no API key for non-local provider",
			profileName: "openai",
			apiKey:      "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := CreateProviderWithFallback(tt.profileName, tt.apiKey, tt.model, tt.baseURL)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			}
		})
	}
}

func TestGetDefaultModel(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"openai", "gpt-4o-mini"},
		{"anthropic", "claude-3-haiku-20240307"},
		{"unknown", "gpt-4o-mini"},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			model := GetDefaultModel(tt.provider)
			assert.Equal(t, tt.expected, model)
		})
	}
}
