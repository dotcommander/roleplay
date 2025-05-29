package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dotcommander/roleplay/internal/models"
)

func TestCharacterRepository_GetCharacterInfo(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "roleplay-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test character with full backstory and speech style
	testChar := &models.Character{
		ID:          "test-char",
		Name:        "Test Character",
		Backstory:   "This is a very long backstory that should not be truncated. It contains detailed information about the character's past, their motivations, and their goals. The backstory is essential for understanding the character and should be displayed in full when listing characters.",
		SpeechStyle: "Speaks in a formal, articulate manner with occasional humor.",
		Quirks:      []string{"Always polite", "Loves tea"},
	}

	// Save character
	if err := repo.SaveCharacter(testChar); err != nil {
		t.Fatal(err)
	}

	// Get character info
	infos, err := repo.GetCharacterInfo()
	if err != nil {
		t.Fatal(err)
	}

	// Verify results
	if len(infos) != 1 {
		t.Fatalf("Expected 1 character, got %d", len(infos))
	}

	info := infos[0]
	
	// Test that full backstory is returned (not truncated)
	if info.Description != testChar.Backstory {
		t.Errorf("Backstory was truncated or modified.\nExpected: %s\nGot: %s", 
			testChar.Backstory, info.Description)
	}

	// Test that speech style is included
	if info.SpeechStyle != testChar.SpeechStyle {
		t.Errorf("Speech style mismatch.\nExpected: %s\nGot: %s",
			testChar.SpeechStyle, info.SpeechStyle)
	}

	// Test that quirks are included
	if len(info.Tags) != len(testChar.Quirks) {
		t.Errorf("Quirks count mismatch. Expected %d, got %d",
			len(testChar.Quirks), len(info.Tags))
	}
}

func TestCharacterRepository_SaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "roleplay-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create repository
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test character
	testChar := &models.Character{
		ID:        "test-123",
		Name:      "Test Character",
		Backstory: "Test backstory",
		Personality: models.PersonalityTraits{
			Openness: 0.5,
		},
	}

	// Save character
	if err := repo.SaveCharacter(testChar); err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	expectedPath := filepath.Join(tempDir, "characters", "test-123.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("Character file was not created")
	}

	// Load character
	loaded, err := repo.LoadCharacter("test-123")
	if err != nil {
		t.Fatal(err)
	}

	// Verify loaded data
	if loaded.ID != testChar.ID {
		t.Errorf("ID mismatch: expected %s, got %s", testChar.ID, loaded.ID)
	}
	if loaded.Name != testChar.Name {
		t.Errorf("Name mismatch: expected %s, got %s", testChar.Name, loaded.Name)
	}
}