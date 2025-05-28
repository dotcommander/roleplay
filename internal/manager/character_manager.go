package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
)

// CharacterManager handles character lifecycle and persistence
type CharacterManager struct {
	bot      *services.CharacterBot
	repo     *repository.CharacterRepository
	sessions *repository.SessionRepository
	mu       sync.RWMutex
	dataDir  string
}

// NewCharacterManager creates a new character manager with fully initialized bot
func NewCharacterManager(cfg *config.Config) (*CharacterManager, error) {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")

	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	sessions := repository.NewSessionRepository(dataDir)
	bot := services.NewCharacterBot(cfg)

	// Initialize provider and UserProfileAgent for the bot
	provider, err := createProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	bot.RegisterProvider(cfg.DefaultProvider, provider)
	bot.InitializeUserProfileAgent()

	return &CharacterManager{
		bot:      bot,
		repo:     repo,
		sessions: sessions,
		dataDir:  dataDir,
	}, nil
}

// LoadAllCharacters loads all persisted characters into memory
func (m *CharacterManager) LoadAllCharacters() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	characters, err := m.repo.ListCharacters()
	if err != nil {
		return err
	}

	for _, id := range characters {
		char, err := m.repo.LoadCharacter(id)
		if err != nil {
			continue
		}

		if err := m.bot.CreateCharacter(char); err != nil {
			return fmt.Errorf("failed to load character %s: %w", id, err)
		}
	}

	return nil
}

// LoadCharacter loads a specific character
func (m *CharacterManager) LoadCharacter(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already loaded
	if _, err := m.bot.GetCharacter(id); err == nil {
		return nil
	}

	// Load from repository
	char, err := m.repo.LoadCharacter(id)
	if err != nil {
		return err
	}

	return m.bot.CreateCharacter(char)
}

// CreateCharacter creates and persists a new character
func (m *CharacterManager) CreateCharacter(char *models.Character) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create in bot
	if err := m.bot.CreateCharacter(char); err != nil {
		return err
	}

	// Persist to disk
	return m.repo.SaveCharacter(char)
}

// GetOrLoadCharacter ensures a character is loaded
func (m *CharacterManager) GetOrLoadCharacter(id string) (*models.Character, error) {
	// First try to get from memory
	char, err := m.bot.GetCharacter(id)
	if err == nil {
		return char, nil
	}

	// Try to load from disk
	if err := m.LoadCharacter(id); err != nil {
		return nil, fmt.Errorf("character %s not found", id)
	}

	return m.bot.GetCharacter(id)
}

// ListAvailableCharacters returns all characters (loaded and unloaded)
func (m *CharacterManager) ListAvailableCharacters() ([]repository.CharacterInfo, error) {
	return m.repo.GetCharacterInfo()
}

// GetBot returns the underlying character bot
func (m *CharacterManager) GetBot() *services.CharacterBot {
	return m.bot
}

// GetSessionRepository returns the session repository
func (m *CharacterManager) GetSessionRepository() *repository.SessionRepository {
	return m.sessions
}

// createProvider creates an AI provider based on the configuration
func createProvider(cfg *config.Config) (providers.AIProvider, error) {
	apiKey := getAPIKey(cfg)
	if apiKey == "" {
		return nil, fmt.Errorf("API key for provider %s not found. Set api_key in config or %s environment variable",
			cfg.DefaultProvider, getEnvVarName(cfg.DefaultProvider))
	}

	switch cfg.DefaultProvider {
	case "anthropic":
		return providers.NewAnthropicProvider(apiKey), nil
	case "openai":
		model := resolveModel(cfg)
		if cfg.BaseURL != "" {
			return providers.NewOpenAIProviderWithBaseURL(apiKey, model, cfg.BaseURL), nil
		}
		return providers.NewOpenAIProvider(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.DefaultProvider)
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

// resolveModel resolves model aliases and returns the actual model name
func resolveModel(cfg *config.Config) string {
	model := cfg.Model
	
	// Check if it's an alias
	if cfg.ModelAliases != nil {
		if resolved, ok := cfg.ModelAliases[model]; ok {
			return resolved
		}
	}
	
	// Return the model as-is if not an alias, or use default
	if model == "" {
		return "gpt-4o-mini" // Default model
	}
	return model
}
