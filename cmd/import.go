package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcommander/roleplay/internal/importer"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [markdown-file]",
	Short: "Import a character from an unstructured markdown file",
	Long: `Import a character from an unstructured markdown file using AI to extract
character information and convert it to the roleplay format.

Example:
  roleplay import /path/to/character.md
  roleplay import ~/Library/Application\ Support/aichat/roles/rick.md`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	markdownPath := args[0]

	if strings.HasPrefix(markdownPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		markdownPath = filepath.Join(home, markdownPath[1:])
	}

	absPath, err := filepath.Abs(markdownPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("markdown file not found: %s", absPath)
	}

	config := GetConfig()

	if config.APIKey == "" {
		return fmt.Errorf("API key not configured. Set ROLEPLAY_API_KEY or use --api-key flag")
	}

	var provider providers.AIProvider
	switch config.DefaultProvider {
	case "anthropic":
		provider = providers.NewAnthropicProvider(config.APIKey)
	case "openai":
		model := config.Model
		if model == "" {
			model = "gpt-4o-mini"
		}
		provider = providers.NewOpenAIProvider(config.APIKey, model)
	default:
		return fmt.Errorf("unknown provider: %s", config.DefaultProvider)
	}

	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}
	characterImporter := importer.NewCharacterImporter(provider, repo)

	fmt.Printf("Importing character from: %s\n", absPath)
	fmt.Println("Analyzing markdown content with AI...")

	ctx := context.Background()
	character, err := characterImporter.ImportFromMarkdown(ctx, absPath)
	if err != nil {
		return fmt.Errorf("failed to import character: %w", err)
	}

	fmt.Printf("\nSuccessfully imported character: %s\n", character.Name)
	fmt.Printf("ID: %s\n", character.ID)
	fmt.Printf("Backstory: %s\n", character.Backstory)
	fmt.Printf("\nPersonality traits:\n")
	fmt.Printf("  Openness: %.2f\n", character.Personality.Openness)
	fmt.Printf("  Conscientiousness: %.2f\n", character.Personality.Conscientiousness)
	fmt.Printf("  Extraversion: %.2f\n", character.Personality.Extraversion)
	fmt.Printf("  Agreeableness: %.2f\n", character.Personality.Agreeableness)
	fmt.Printf("  Neuroticism: %.2f\n", character.Personality.Neuroticism)

	charactersDir := filepath.Join(dataDir, "characters")
	fmt.Printf("\nCharacter saved to: %s\n", filepath.Join(charactersDir, character.ID+".json"))
	fmt.Printf("\nYou can now chat with this character using:\n")
	fmt.Printf("  roleplay chat \"Hello!\" --character %s\n", character.ID)
	fmt.Printf("  roleplay interactive --character %s\n", character.ID)

	return nil
}
