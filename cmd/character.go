package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var characterCmd = &cobra.Command{
	Use:   "character",
	Short: "Manage characters",
	Long:  `Create, list, and manage character profiles for the roleplay bot.`,
}

var listCharactersCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all available characters",
	RunE:    runListCharacters,
}

var createCharacterCmd = &cobra.Command{
	Use:   "create [character-file.json]",
	Short: "Create a new character from a JSON file",
	Long: `Create a new character by loading their profile from a JSON file.

The JSON file should contain:
{
  "id": "warrior-123",
  "name": "Thorin Ironforge",
  "backstory": "A veteran dwarf warrior...",
  "personality": {
    "openness": 0.3,
    "conscientiousness": 0.8,
    "extraversion": 0.4,
    "agreeableness": 0.6,
    "neuroticism": 0.7
  },
  "current_mood": {
    "joy": 0.2,
    "anger": 0.4,
    "sadness": 0.3
  },
  "quirks": ["Always checks exits", "Touches scars when nervous"],
  "speech_style": "Formal, archaic. Uses 'ye' and 'aye'."
}`,
	Args: cobra.ExactArgs(1),
	RunE: runCreateCharacter,
}

var showCharacterCmd = &cobra.Command{
	Use:   "show [character-id]",
	Short: "Show character details",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowCharacter,
}

var exampleCharacterCmd = &cobra.Command{
	Use:   "example",
	Short: "Generate an example character JSON file",
	RunE:  runExampleCharacter,
}

func init() {
	rootCmd.AddCommand(characterCmd)
	characterCmd.AddCommand(createCharacterCmd)
	characterCmd.AddCommand(showCharacterCmd)
	characterCmd.AddCommand(exampleCharacterCmd)
	characterCmd.AddCommand(listCharactersCmd)
}

func runCreateCharacter(cmd *cobra.Command, args []string) error {
	filename := args[0]
	config := GetConfig()

	// Read character file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read character file: %w", err)
	}

	// Parse character
	var char models.Character
	if err := json.Unmarshal(data, &char); err != nil {
		return fmt.Errorf("failed to parse character JSON: %w", err)
	}

	// Initialize manager (bot is now fully initialized)
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	// Create character using manager (handles both bot and persistence)
	if err := mgr.CreateCharacter(&char); err != nil {
		return fmt.Errorf("failed to create character: %w", err)
	}

	fmt.Printf("Character '%s' (ID: %s) created and saved successfully!\n", char.Name, char.ID)
	return nil
}

func runShowCharacter(cmd *cobra.Command, args []string) error {
	characterID := args[0]

	// Initialize repository to load from disk
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Load character from repository
	char, err := repo.LoadCharacter(characterID)
	if err != nil {
		return fmt.Errorf("character %s not found", characterID)
	}

	// Display character
	output, err := json.MarshalIndent(char, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format character: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func runExampleCharacter(cmd *cobra.Command, args []string) error {
	example := models.Character{
		ID:   "warrior-123",
		Name: "Thorin Ironforge",
		Backstory: `A veteran dwarf warrior from the Mountain Kingdoms. 
Survived the Battle of Five Armies. Gruff exterior hiding a heart of gold.
Lost his brother in battle, carries survivor's guilt.`,
		Personality: models.PersonalityTraits{
			Openness:          0.3,
			Conscientiousness: 0.8,
			Extraversion:      0.4,
			Agreeableness:     0.6,
			Neuroticism:       0.7,
		},
		CurrentMood: models.EmotionalState{
			Joy:     0.2,
			Anger:   0.4,
			Sadness: 0.3,
		},
		Quirks: []string{
			"Always checks exits when entering a room",
			"Unconsciously touches battle scars when nervous",
			"Refuses to sit with back to the door",
		},
		SpeechStyle: "Formal, archaic. Uses 'ye' and 'aye'. Short sentences. Military precision.",
		Memories: []models.Memory{
			{
				Type:      models.LongTermMemory,
				Content:   "Brother's last words: 'Protect the clan'",
				Emotional: 0.95,
			},
		},
	}

	output, err := json.MarshalIndent(&example, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format example: %w", err)
	}

	fmt.Println(string(output))
	fmt.Fprintln(os.Stderr, "\nSave this to a file (e.g., thorin.json) and use 'roleplay character create thorin.json'")
	return nil
}

func runListCharacters(cmd *cobra.Command, args []string) error {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return err
	}

	characters, err := repo.GetCharacterInfo()
	if err != nil {
		return err
	}

	if len(characters) == 0 {
		fmt.Println("No characters found. Create one with 'roleplay character create <file.json>'")
		return nil
	}

	// Display characters in a visually appealing format
	fmt.Println("🎭 Available Characters")
	fmt.Println("======================")
	fmt.Println()

	for i, char := range characters {
		// Character header with emoji
		fmt.Printf("📚 %s (%s)\n", char.Name, char.ID)
		fmt.Println(strings.Repeat("─", 50))
		
		// Backstory with proper word wrapping
		if char.Description != "" {
			wrapped := wrapText(char.Description, 70)
			for _, line := range wrapped {
				fmt.Printf("   %s\n", line)
			}
		} else {
			fmt.Println("   [No backstory provided]")
		}
		
		// Display speech style if available
		if char.SpeechStyle != "" {
			fmt.Println()
			fmt.Println("   💬 Speech Style:")
			speechWrapped := wrapText(char.SpeechStyle, 66)
			for _, line := range speechWrapped {
				fmt.Printf("      %s\n", line)
			}
		}
		
		// Display quirks if available
		if len(char.Tags) > 0 {
			fmt.Println()
			fmt.Println("   🎯 Quirks:")
			for _, quirk := range char.Tags {
				fmt.Printf("      • %s\n", quirk)
			}
		}
		
		// Add spacing between characters except for the last one
		if i < len(characters)-1 {
			fmt.Println()
		}
	}
	
	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("Total: %d character(s)\n", len(characters))
	
	return nil
}

// wrapText wraps text to the specified width
func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	
	if len(words) == 0 {
		return lines
	}
	
	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) > width {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine += " " + word
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return lines
}
