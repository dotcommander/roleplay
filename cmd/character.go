package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/spf13/cobra"
)

var characterCmd = &cobra.Command{
	Use:   "character",
	Short: "Manage characters",
	Long:  `Create, list, and manage character profiles for the roleplay bot.`,
}

var listCharactersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available characters",
	RunE:  runListCharacters,
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

	// Initialize bot
	bot := services.NewCharacterBot(config)

	// Register provider using factory (needed for character creation warmup)
	if err := factory.InitializeAndRegisterProvider(bot, config); err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}

	// Create character
	if err := bot.CreateCharacter(&char); err != nil {
		return fmt.Errorf("failed to create character: %w", err)
	}

	// Save to repository for persistence
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	if err := repo.SaveCharacter(&char); err != nil {
		return fmt.Errorf("failed to save character: %w", err)
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION")

	for _, char := range characters {
		fmt.Fprintf(w, "%s\t%s\t%s\n", char.ID, char.Name, char.Description)
	}

	w.Flush()
	return nil
}
