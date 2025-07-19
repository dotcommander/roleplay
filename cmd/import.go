package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/importer"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/pkg/bridge"

	"github.com/spf13/cobra"
)

var (
	sourceFormat string
	verbose      bool
)

var importCharacterCmd = &cobra.Command{
	Use:   "import [character-file]",
	Short: "Import a character from markdown or Characters format",
	Long: `Import a character from various formats including:
- Unstructured markdown files (using AI extraction)
- Characters format JSON files (direct conversion)

The command will auto-detect the format by default, or you can specify it explicitly.

Examples:
  roleplay character import /path/to/character.md
  roleplay character import /path/to/character.json --source characters
  roleplay character import ~/Library/Application\ Support/aichat/roles/rick.md`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	characterCmd.AddCommand(importCharacterCmd)
	importCharacterCmd.Flags().StringVar(&sourceFormat, "source", "auto", "Source format: auto, markdown, characters")
	importCharacterCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed conversion information")
}

func runImport(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Expand ~ to home directory
	if strings.HasPrefix(filePath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		filePath = filepath.Join(home, filePath[1:])
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check file exists
	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("file not found: %s", absPath)
	}

	// Detect format
	format := sourceFormat
	if format == "auto" {
		format = detectFormat(absPath)
		if verbose {
			cmd.Printf("Auto-detected format: %s\n", format)
		}
	}

	config := GetConfig()
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")

	var character *models.Character
	var warnings []string

	switch format {
	case "characters":
		// Import from Characters format
		if verbose {
			cmd.Println("Importing from Characters format...")
		}
		character, warnings, err = importFromCharactersFormat(cmd, absPath, dataDir)
		if err != nil {
			return fmt.Errorf("failed to import from Characters format: %w", err)
		}

	case "markdown":
		// Import from markdown using AI
		if verbose {
			cmd.Println("Importing from markdown using AI extraction...")
		}
		character, err = importFromMarkdown(cmd, absPath, config, dataDir)
		if err != nil {
			return fmt.Errorf("failed to import from markdown: %w", err)
		}

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	// Display import results
	cmd.Printf("\nSuccessfully imported character: %s\n", character.Name)
	cmd.Printf("ID: %s\n", character.ID)
	
	if verbose {
		cmd.Printf("\nDetailed information:\n")
		cmd.Printf("Backstory: %s\n", character.Backstory)
		cmd.Printf("\nPersonality traits:\n")
		cmd.Printf("  Openness: %.2f\n", character.Personality.Openness)
		cmd.Printf("  Conscientiousness: %.2f\n", character.Personality.Conscientiousness)
		cmd.Printf("  Extraversion: %.2f\n", character.Personality.Extraversion)
		cmd.Printf("  Agreeableness: %.2f\n", character.Personality.Agreeableness)
		cmd.Printf("  Neuroticism: %.2f\n", character.Personality.Neuroticism)
		
		if len(character.Quirks) > 0 {
			cmd.Printf("\nQuirks:\n")
			for _, quirk := range character.Quirks {
				cmd.Printf("  - %s\n", quirk)
			}
		}
		
		if character.SpeechStyle != "" {
			cmd.Printf("\nSpeech Style: %s\n", character.SpeechStyle)
		}
	}

	// Display warnings if any
	if len(warnings) > 0 {
		cmd.Printf("\nConversion warnings:\n")
		for _, warning := range warnings {
			cmd.Printf("  ⚠️  %s\n", warning)
		}
	}

	charactersDir := filepath.Join(dataDir, "characters")
	cmd.Printf("\nCharacter saved to: %s\n", filepath.Join(charactersDir, character.ID+".json"))
	cmd.Printf("\nYou can now chat with this character using:\n")
	cmd.Printf("  roleplay chat \"Hello!\" --character %s\n", character.ID)
	cmd.Printf("  roleplay interactive --character %s\n", character.ID)

	return nil
}

// detectFormat attempts to detect the format of the input file
func detectFormat(filePath string) string {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	
	if ext == ".md" || ext == ".markdown" {
		return "markdown"
	}
	
	if ext == ".json" {
		// Try to read and detect JSON structure
		data, err := os.ReadFile(filePath)
		if err == nil {
			var jsonData map[string]interface{}
			if err := json.Unmarshal(data, &jsonData); err == nil {
				// Check for Characters format indicators
				if _, hasTraits := jsonData["traits"]; hasTraits {
					return "characters"
				}
				if _, hasAttributes := jsonData["attributes"]; hasAttributes {
					return "characters"
				}
				if _, hasArchetype := jsonData["archetype"]; hasArchetype {
					return "characters"
				}
				// Check for roleplay format indicators
				if _, hasPersonality := jsonData["personality"]; hasPersonality {
					return "roleplay"
				}
			}
		}
	}
	
	// Default to markdown for unknown formats
	return "markdown"
}

// importFromCharactersFormat imports a character from Characters format
func importFromCharactersFormat(cmd *cobra.Command, filePath string, dataDir string) (*models.Character, []string, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse as JSON
	var charData map[string]interface{}
	if err := json.Unmarshal(data, &charData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create converter
	converter := bridge.NewCharactersConverter()
	
	// Check if converter can handle this data
	if !converter.CanConvert(charData) {
		return nil, nil, fmt.Errorf("file does not appear to be in Characters format")
	}

	// Convert to universal format
	ctx := context.Background()
	universal, err := converter.ToUniversal(ctx, charData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert to universal format: %w", err)
	}

	// Convert from universal to roleplay format
	result, err := converter.FromUniversal(ctx, universal)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert to roleplay format: %w", err)
	}

	character, ok := result.(*models.Character)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected result type from converter")
	}

	// Create manager to save the character
	mgr, err := manager.NewCharacterManagerWithoutProvider(GetConfig())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize manager: %w", err)
	}

	// Save the character
	if err := mgr.CreateCharacter(character); err != nil {
		return nil, nil, fmt.Errorf("failed to save character: %w", err)
	}

	// Collect any warnings
	warnings := []string{}
	
	// Check for missing fields that might need manual adjustment
	if character.Backstory == "" {
		warnings = append(warnings, "No backstory found - you may want to add one manually")
	}
	
	if len(character.Quirks) == 0 {
		warnings = append(warnings, "No quirks found - consider adding some for more personality")
	}
	
	if character.SpeechStyle == "" {
		warnings = append(warnings, "No speech style defined - consider adding one for more authentic dialogue")
	}
	
	// Check if NSFW content was detected
	if nsfw, ok := universal.SourceData["nsfw"].(bool); ok && nsfw {
		warnings = append(warnings, "NSFW content detected in original - some content may have been filtered")
	}

	return character, warnings, nil
}

// importFromMarkdown imports a character from markdown using AI
func importFromMarkdown(cmd *cobra.Command, filePath string, config *config.Config, dataDir string) (*models.Character, error) {
	// Create provider using factory
	provider, err := factory.CreateProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	// Create repository and importer
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	
	characterImporter := importer.NewCharacterImporter(provider, repo)

	cmd.Printf("Importing character from: %s\n", filePath)
	cmd.Println("Analyzing markdown content with AI...")

	ctx := context.Background()
	character, err := characterImporter.ImportFromMarkdown(ctx, filePath)
	if err != nil {
		return nil, err
	}

	return character, nil
}
