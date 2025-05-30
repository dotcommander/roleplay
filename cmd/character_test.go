package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/spf13/cobra"
)

func TestCharacterCommands(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) string // Returns temp dir
		wantErr     bool
		wantOutput  string
		checkResult func(t *testing.T, tempDir string)
	}{
		{
			name: "create character from JSON",
			args: []string{"character", "create"},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create test character JSON
				char := models.Character{
					ID:        "test-char",
					Name:      "Test Character",
					Backstory: "A test character",
					Personality: models.PersonalityTraits{
						Openness: 0.5,
					},
				}
				
				data, err := json.Marshal(&char)
				if err != nil {
					t.Fatalf("Failed to marshal character: %v", err)
				}
				charFile := filepath.Join(tempDir, "test.json")
				if err := os.WriteFile(charFile, data, 0644); err != nil {
					t.Fatalf("Failed to write character file: %v", err)
				}
				
				// Add file path to args
				os.Args = append(os.Args, charFile)
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "Character 'Test Character' (ID: test-char) created",
			checkResult: func(t *testing.T, tempDir string) {
				// Check character was saved
				charFile := filepath.Join(tempDir, ".config", "roleplay", "characters", "test-char.json")
				if _, err := os.Stat(charFile); os.IsNotExist(err) {
					t.Error("Character file was not created")
				}
			},
		},
		{
			name: "create character - invalid JSON",
			args: []string{"character", "create"},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create invalid JSON
				invalidFile := filepath.Join(tempDir, "invalid.json")
				if err := os.WriteFile(invalidFile, []byte("{ invalid json"), 0644); err != nil {
					t.Fatalf("Failed to write invalid JSON file: %v", err)
				}
				
				os.Args = append(os.Args, invalidFile)
				
				return tempDir
			},
			wantErr:    true,
			wantOutput: "failed to parse character JSON",
		},
		{
			name: "create character - missing file",
			args: []string{"character", "create", "/non/existent/file.json"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:    true,
			wantOutput: "failed to read character file",
		},
		{
			name: "list characters - empty",
			args: []string{"character", "list"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:    false,
			wantOutput: "No characters found",
		},
		{
			name: "list characters - with characters",
			args: []string{"character", "list"},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create character directory and files
				charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
				if err := os.MkdirAll(charDir, 0755); err != nil {
					t.Fatalf("Failed to create character directory: %v", err)
				}
				
				// Create test characters
				chars := []models.Character{
					{
						ID:          "char1",
						Name:        "Character One",
						Backstory:   "First character",
						Quirks:      []string{"quirk1"},
						SpeechStyle: "Formal",
					},
					{
						ID:        "char2",
						Name:      "Character Two",
						Backstory: "Second character",
					},
				}
				
				for i := range chars {
					data, err := json.Marshal(&chars[i])
					if err != nil {
						t.Fatalf("Failed to marshal character %s: %v", chars[i].ID, err)
					}
					if err := os.WriteFile(filepath.Join(charDir, chars[i].ID+".json"), data, 0644); err != nil {
						t.Fatalf("Failed to write character file %s: %v", chars[i].ID, err)
					}
				}
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "Available Characters",
			checkResult: func(t *testing.T, tempDir string) {
				// Output should contain both characters
				// This would need to capture stdout to verify
			},
		},
		{
			name: "show character",
			args: []string{"character", "show", "test-char"},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create character
				charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
				if err := os.MkdirAll(charDir, 0755); err != nil {
					t.Fatalf("Failed to create character directory: %v", err)
				}
				
				char := models.Character{
					ID:        "test-char",
					Name:      "Test Character",
					Backstory: "Test backstory",
				}
				
				data, err := json.Marshal(&char)
				if err != nil {
					t.Fatalf("Failed to marshal character: %v", err)
				}
				if err := os.WriteFile(filepath.Join(charDir, "test-char.json"), data, 0644); err != nil {
					t.Fatalf("Failed to write character file: %v", err)
				}
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "Test Character",
		},
		{
			name: "show character - not found",
			args: []string{"character", "show", "nonexistent"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:    true,
			wantOutput: "character nonexistent not found",
		},
		{
			name: "example character",
			args: []string{"character", "example"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:    false,
			wantOutput: "Thorin Ironforge", // Example character name
		},
		{
			name: "import character from markdown",
			args: []string{"character", "import"},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				// Check if prompts directory exists
				if _, err := os.Stat("prompts/character-import.md"); os.IsNotExist(err) {
					// Try from parent directory (when running from cmd/)
					if _, err := os.Stat("../prompts/character-import.md"); os.IsNotExist(err) {
						t.Skip("Skipping import test: prompts directory not found")
						return "" // This return is needed after t.Skip
					}
				}
				
				// Create markdown file
				mdFile := filepath.Join(tempDir, "character.md")
				mdContent := `# Rick Sanchez

The smartest man in the universe from dimension C-137. 
Cynical, alcoholic mad scientist who drags his grandson Morty on dangerous adventures.

## Personality
- Extremely intelligent
- Nihilistic
- Alcoholic
- Narcissistic

## Speech Style
Constantly burps mid-sentence (*burp*). Uses crude language mixed with scientific terms.
Often says "Morty" as punctuation.

## Quirks
- Drinks from flask constantly
- Portal gun inventor
- Believes nothing matters`
				
				if err := os.WriteFile(mdFile, []byte(mdContent), 0644); err != nil {
					t.Fatalf("Failed to write markdown file: %v", err)
				}
				os.Args = append(os.Args, mdFile)
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "Successfully imported character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset os.Args
			os.Args = []string{"roleplay"}
			os.Args = append(os.Args, tt.args...)
			
			// Setup test environment
			tempDir := ""
			oldHome := os.Getenv("HOME")
			if tt.setup != nil {
				tempDir = tt.setup(t)
				// Set HOME to temp dir for config
				os.Setenv("HOME", tempDir)
				defer func() {
					if oldHome != "" {
						os.Setenv("HOME", oldHome)
					}
				}()
			}

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Execute command
			err := rootCmd.Execute()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check output
			output := buf.String()
			if tt.wantOutput != "" && !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Output does not contain expected string.\nGot: %s\nWant substring: %s", output, tt.wantOutput)
			}

			// Additional checks
			if tt.checkResult != nil && tempDir != "" {
				tt.checkResult(t, tempDir)
			}

			// Reset command for next test
			rootCmd = &cobra.Command{
				Use:   "roleplay",
				Short: "A sophisticated character bot with psychological modeling",
			}
			initCommands()
		})
	}
}

func TestCharacterListFormatting(t *testing.T) {
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		}
	}()

	// Create character directory
	charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatalf("Failed to create character directory: %v", err)
	}

	// Create a character with all fields
	char := models.Character{
		ID:          "detailed-char",
		Name:        "Detailed Character",
		Backstory:   "This is a very long backstory that should be wrapped properly when displayed. It contains multiple sentences and should demonstrate the text wrapping functionality.",
		Quirks:      []string{"Always punctual", "Speaks in rhymes", "Afraid of spiders"},
		SpeechStyle: "Speaks in iambic pentameter with occasional modern slang mixed in for comedic effect.",
	}

	data, err := json.Marshal(&char)
	if err != nil {
		t.Fatalf("Failed to marshal character: %v", err)
	}
	if err := os.WriteFile(filepath.Join(charDir, char.ID+".json"), data, 0644); err != nil {
		t.Fatalf("Failed to write character file: %v", err)
	}

	// Run list command
	os.Args = []string{"roleplay", "character", "list"}
	
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute list command: %v", err)
	}

	output := buf.String()

	// Check formatting elements
	checks := []string{
		"🎭 Available Characters",
		"📚 Detailed Character (detailed-char)",
		"💬 Speech Style:",
		"🎯 Quirks:",
		"• Always punctual",
		"• Speaks in rhymes",
		"• Afraid of spiders",
		"Total: 1 character(s)",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Output missing expected element: %s", check)
		}
	}
}

func TestCharacterCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "create without file",
			args:    []string{"character", "create"},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:    "show without ID",
			args:    []string{"character", "show"},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:    "import without file",
			args:    []string{"character", "import"},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:    "invalid subcommand",
			args:    []string{"character", "invalid"},
			wantErr: false, // Cobra doesn't return error for unknown subcommand, just shows help
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset commands before each test
			rootCmd = &cobra.Command{
				Use:   "roleplay",
				Short: "A sophisticated character bot with psychological modeling",
			}
			initCommands()
			
			os.Args = []string{"roleplay"}
			os.Args = append(os.Args, tt.args...)

			var buf bytes.Buffer
			rootCmd.SetErr(&buf)
			rootCmd.SetOut(&buf)

			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errMsg != "" {
				errOutput := buf.String()
				// Also check stdout since cobra might write there
				if err != nil && err.Error() != "" {
					errOutput = err.Error() + "\n" + errOutput
				}
				if !strings.Contains(errOutput, tt.errMsg) {
					t.Errorf("Error output does not contain expected message.\nGot: %s\nWant substring: %s", errOutput, tt.errMsg)
				}
			}

			// Reset
			rootCmd = &cobra.Command{
				Use:   "roleplay",
				Short: "A sophisticated character bot with psychological modeling",
			}
			initCommands()
		})
	}
}

// Helper function to initialize commands for testing
func initCommands() {
	// Re-add all commands to rootCmd after reset
	characterCmd = &cobra.Command{
		Use:   "character",
		Short: "Manage characters",
		Long:  `Create, list, and manage character profiles for the roleplay bot.`,
	}
	
	createCharacterCmd = &cobra.Command{
		Use:   "create [character-file.json]",
		Short: "Create a new character from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateCharacter,
	}
	
	showCharacterCmd = &cobra.Command{
		Use:   "show [character-id]",
		Short: "Show details of a specific character",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowCharacter,
	}
	
	exampleCharacterCmd = &cobra.Command{
		Use:   "example",
		Short: "Show an example character JSON",
		RunE:  runExampleCharacter,
	}
	
	listCharactersCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available characters",
		RunE:    runListCharacters,
	}
	
	// Recreate importCharacterCmd as defined in import.go
	importCharacterCmd = &cobra.Command{
		Use:   "import [markdown-file]",
		Short: "Import a character from an unstructured markdown file",
		Args:  cobra.ExactArgs(1),
		RunE:  runImport,
	}
	
	// Add commands to hierarchy
	rootCmd.AddCommand(characterCmd)
	characterCmd.AddCommand(createCharacterCmd)
	characterCmd.AddCommand(showCharacterCmd)
	characterCmd.AddCommand(exampleCharacterCmd)
	characterCmd.AddCommand(listCharactersCmd)
	characterCmd.AddCommand(importCharacterCmd)
}