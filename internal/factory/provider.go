package factory

import (
	"fmt"
	"os"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/services"
)

// CreateProvider creates an AI provider based on the configuration
func CreateProvider(cfg *config.Config) (providers.AIProvider, error) {
	apiKey := getAPIKey(cfg)
	if apiKey == "" {
		return nil, fmt.Errorf("API key for provider %s not found. Set api_key in config or %s environment variable",
			cfg.DefaultProvider, getEnvVarName(cfg.DefaultProvider))
	}

	switch cfg.DefaultProvider {
	case "anthropic":
		return providers.NewAnthropicProvider(apiKey), nil
	case "openai":
		model := cfg.Model
		if model == "" {
			model = "gpt-4o-mini" // Centralized default
		}
		return providers.NewOpenAIProvider(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
	}
}

// InitializeAndRegisterProvider creates and registers a provider with the bot
func InitializeAndRegisterProvider(bot *services.CharacterBot, cfg *config.Config) error {
	provider, err := CreateProvider(cfg)
	if err != nil {
		return err
	}

	bot.RegisterProvider(cfg.DefaultProvider, provider)
	
	// Initialize user profile agent after provider is registered
	bot.InitializeUserProfileAgent()
	
	return nil
}

// CreateProviderWithFallback creates a provider with environment variable fallback
// This is useful for commands that don't use the full config structure
func CreateProviderWithFallback(providerName, apiKey, model string) (providers.AIProvider, error) {
	// Try environment variable if API key not provided
	if apiKey == "" {
		apiKey = os.Getenv(getEnvVarName(providerName))
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key for provider %s not found", providerName)
	}

	switch providerName {
	case "anthropic":
		return providers.NewAnthropicProvider(apiKey), nil
	case "openai":
		if model == "" {
			model = "gpt-4o-mini"
		}
		return providers.NewOpenAIProvider(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

// GetDefaultModel returns the default model for a provider
func GetDefaultModel(providerName string) string {
	switch providerName {
	case "openai":
		return "gpt-4o-mini"
	case "anthropic":
		return "claude-3-haiku-20240307"
	default:
		return ""
	}
}

// getAPIKey retrieves the API key from config or environment
func getAPIKey(cfg *config.Config) string {
	if cfg.APIKey != "" {
		return cfg.APIKey
	}

	// Fall back to environment variable
	return os.Getenv(getEnvVarName(cfg.DefaultProvider))
}

// getEnvVarName returns the environment variable name for a provider
func getEnvVarName(provider string) string {
	switch provider {
	case "openai":
		return "OPENAI_API_KEY"
	case "anthropic":
		return "ANTHROPIC_API_KEY"
	default:
		return ""
	}
}
