package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current configuration and status",
	Long:  `Display the current provider, model, and other configuration settings.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Get configuration
	provider := viper.GetString("provider")
	model := viper.GetString("model")
	apiKey := viper.GetString("api_key")

	// Check for API key from environment if not set
	if apiKey == "" && provider == "openai" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" && provider == "anthropic" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	// Determine actual model that will be used
	actualModel := model
	if actualModel == "" {
		switch provider {
		case "openai":
			actualModel = "gpt-4o-mini"
		case "anthropic":
			actualModel = "claude-3-haiku-20240307"
		}
	}

	fmt.Println("🤖 Roleplay Status")
	fmt.Println("==================")
	fmt.Printf("Provider: %s\n", provider)
	fmt.Printf("Model: %s\n", actualModel)
	if model != "" && model != actualModel {
		fmt.Printf("  (configured: %s, using default: %s)\n", model, actualModel)
	}
	fmt.Printf("API Key: %s\n", func() string {
		if apiKey == "" {
			return "❌ Not configured"
		}
		if len(apiKey) > 8 {
			return "✓ " + apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		}
		return "✓ Configured"
	}())

	// Show cache configuration
	fmt.Printf("\nCache Configuration:\n")
	fmt.Printf("  Default TTL: %v\n", viper.GetDuration("cache.default_ttl"))
	fmt.Printf("  Adaptive TTL: %v\n", viper.GetBool("cache.adaptive_ttl"))

	// Show data directory
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	fmt.Printf("\nData Directory: %s\n", dataDir)

	// Show character count
	charRepo, err := repository.NewCharacterRepository(dataDir)
	if err == nil {
		chars, _ := charRepo.ListCharacters()
		fmt.Printf("Characters: %d available\n", len(chars))
	}

	// Show session count
	sessionRepo := repository.NewSessionRepository(dataDir)
	totalSessions := 0
	if charRepo != nil {
		chars, _ := charRepo.ListCharacters()
		for _, charID := range chars {
			sessions, _ := sessionRepo.ListSessions(charID)
			totalSessions += len(sessions)
		}
	}
	fmt.Printf("Sessions: %d total\n", totalSessions)

	return nil
}
