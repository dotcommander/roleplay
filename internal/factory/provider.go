package factory

import (
	"fmt"
	"strings"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/services"
)

// CreateProvider creates an AI provider based on the configuration
func CreateProvider(cfg *config.Config) (providers.AIProvider, error) {
	// Check for test mode
	if cfg.DefaultProvider == "mock" {
		return providers.NewMockProvider(), nil
	}

	// The API key might be optional for local services like Ollama
	apiKey := cfg.APIKey
	baseURL := cfg.BaseURL
	model := cfg.Model

	// Apply sensible defaults based on the profile name (DefaultProvider)
	profileName := strings.ToLower(cfg.DefaultProvider)

	// Model defaults based on profile
	if model == "" {
		switch profileName {
		case "ollama":
			model = "llama3"
		case "openai":
			model = "gpt-4o-mini"
		case "anthropic", "anthropic_compatible":
			model = "claude-3-haiku-20240307"
		default:
			model = "gpt-4o-mini" // Safe default
		}
	}

	// Validate API key for non-local endpoints
	if apiKey == "" && !isLocalEndpoint(profileName, baseURL) {
		return nil, fmt.Errorf("API key required for %s. Set api_key in config or environment variable", profileName)
	}

	// Always create the unified OpenAI-compatible provider
	return providers.NewOpenAIProviderWithBaseURL(apiKey, model, baseURL), nil
}

// InitializeAndRegisterProvider creates and registers a provider with the bot
func InitializeAndRegisterProvider(bot *services.CharacterBot, cfg *config.Config) error {
	provider, err := CreateProvider(cfg)
	if err != nil {
		return err
	}

	// Register using the profile name as the key
	bot.RegisterProvider(cfg.DefaultProvider, provider)

	// Initialize user profile agent after provider is registered
	bot.InitializeUserProfileAgent()

	return nil
}

// CreateProviderWithFallback creates a provider with sensible defaults
// This is useful for commands that don't use the full config structure
func CreateProviderWithFallback(profileName, apiKey, model, baseURL string) (providers.AIProvider, error) {
	// Apply model defaults based on profile
	if model == "" {
		model = GetDefaultModel(profileName)
	}

	// Validate API key for non-local endpoints
	if apiKey == "" && !isLocalEndpoint(profileName, baseURL) {
		return nil, fmt.Errorf(`API key for provider '%s' is missing.

To fix this:
1. Run 'roleplay init' to configure it interactively
2. Or, set the 'api_key' in ~/.config/roleplay/config.yaml
3. Or, set the %s environment variable
4. Or, use the --api-key flag

For more info: roleplay config where`, profileName, getProviderEnvVar(profileName))
	}

	// Always create the unified OpenAI-compatible provider
	return providers.NewOpenAIProviderWithBaseURL(apiKey, model, baseURL), nil
}

// GetDefaultModel returns the default model for a profile
func GetDefaultModel(profileName string) string {
	profileName = strings.ToLower(profileName)
	switch profileName {
	case "openai":
		return "gpt-4o-mini"
	case "anthropic", "anthropic_compatible":
		return "claude-3-haiku-20240307"
	case "ollama":
		return "llama3"
	case "gemini", "gemini_compatible":
		return "gemini-1.5-flash"
	default:
		return "gpt-4o-mini"
	}
}

// isLocalEndpoint determines if an endpoint is local and doesn't require API key
func isLocalEndpoint(profileName, baseURL string) bool {
	profileName = strings.ToLower(profileName)

	// Known local profiles
	if profileName == "ollama" || profileName == "lm_studio" || profileName == "local" {
		return true
	}

	// Check if baseURL indicates localhost
	if baseURL != "" && (strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") || strings.Contains(baseURL, "0.0.0.0")) {
		return true
	}

	return false
}

// getProviderEnvVar returns the appropriate environment variable name for a provider
func getProviderEnvVar(profileName string) string {
	profileName = strings.ToLower(profileName)
	switch profileName {
	case "openai":
		return "OPENAI_API_KEY or ROLEPLAY_API_KEY"
	case "anthropic":
		return "ANTHROPIC_API_KEY or ROLEPLAY_API_KEY"
	case "gemini":
		return "GEMINI_API_KEY or ROLEPLAY_API_KEY"
	case "groq":
		return "GROQ_API_KEY or ROLEPLAY_API_KEY"
	default:
		return "ROLEPLAY_API_KEY"
	}
}
