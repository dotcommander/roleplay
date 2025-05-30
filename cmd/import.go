package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/importer"
	"github.com/dotcommander/roleplay/internal/repository"

	"github.com/spf13/cobra"
)

var importCharacterCmd = &cobra.Command{
	Use:   "import [markdown-file]",
	Short: "Import a character from an unstructured markdown file",
	Long: `Import a character from an unstructured markdown file using AI to extract
character information and convert it to the roleplay format.

Example:
  roleplay character import /path/to/character.md
  roleplay character import ~/Library/Application\ Support/aichat/roles/rick.md`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	characterCmd.AddCommand(importCharacterCmd)
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

	// Create provider using factory
	provider, err := factory.CreateProvider(config)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}
	characterImporter := importer.NewCharacterImporter(provider, repo)

	cmd.Printf("Importing character from: %s\n", absPath)
	cmd.Println("Analyzing markdown content with AI...")

	ctx := context.Background()
	character, err := characterImporter.ImportFromMarkdown(ctx, absPath)
	if err != nil {
		return fmt.Errorf("failed to import character: %w", err)
	}

	cmd.Printf("\nSuccessfully imported character: %s\n", character.Name)
	cmd.Printf("ID: %s\n", character.ID)
	cmd.Printf("Backstory: %s\n", character.Backstory)
	cmd.Printf("\nPersonality traits:\n")
	cmd.Printf("  Openness: %.2f\n", character.Personality.Openness)
	cmd.Printf("  Conscientiousness: %.2f\n", character.Personality.Conscientiousness)
	cmd.Printf("  Extraversion: %.2f\n", character.Personality.Extraversion)
	cmd.Printf("  Agreeableness: %.2f\n", character.Personality.Agreeableness)
	cmd.Printf("  Neuroticism: %.2f\n", character.Personality.Neuroticism)

	charactersDir := filepath.Join(dataDir, "characters")
	cmd.Printf("\nCharacter saved to: %s\n", filepath.Join(charactersDir, character.ID+".json"))
	cmd.Printf("\nYou can now chat with this character using:\n")
	cmd.Printf("  roleplay chat \"Hello!\" --character %s\n", character.ID)
	cmd.Printf("  roleplay interactive --character %s\n", character.ID)

	return nil
}
