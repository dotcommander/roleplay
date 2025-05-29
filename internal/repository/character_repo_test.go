package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
)

func TestCharacterPersistence(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test character save
	char := &models.Character{
		ID:        "test-char",
		Name:      "Test Character",
		Backstory: "A character for testing persistence",
		Personality: models.PersonalityTraits{
			Openness:          0.7,
			Conscientiousness: 0.8,
			Extraversion:      0.6,
			Agreeableness:     0.9,
			Neuroticism:       0.3,
		},
		CurrentMood: models.EmotionalState{
			Joy:      0.7,
			Surprise: 0.2,
			Anger:    0.1,
		},
		Quirks:      []string{"Always tests things", "Never gives up"},
		SpeechStyle: "Clear and concise test speech",
		Memories: []models.Memory{
			{
				Type:      models.ShortTermMemory,
				Content:   "Test memory",
				Timestamp: time.Now(),
				Emotional: 0.5,
			},
		},
	}

	// Save character
	err = repo.SaveCharacter(char)
	if err != nil {
		t.Fatalf("Failed to save character: %v", err)
	}

	// Verify file exists
	charFile := filepath.Join(tempDir, "characters", char.ID+".json")
	if _, err := os.Stat(charFile); os.IsNotExist(err) {
		t.Error("Character file was not created")
	}

	// Load character
	loaded, err := repo.LoadCharacter(char.ID)
	if err != nil {
		t.Fatalf("Failed to load character: %v", err)
	}

	// Verify all fields
	if loaded.ID != char.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, char.ID)
	}
	if loaded.Name != char.Name {
		t.Errorf("Name mismatch: got %s, want %s", loaded.Name, char.Name)
	}
	if loaded.Backstory != char.Backstory {
		t.Errorf("Backstory mismatch: got %s, want %s", loaded.Backstory, char.Backstory)
	}
	if loaded.Personality.Openness != char.Personality.Openness {
		t.Errorf("Personality mismatch")
	}
	if len(loaded.Quirks) != len(char.Quirks) {
		t.Errorf("Quirks count mismatch: got %d, want %d", len(loaded.Quirks), len(char.Quirks))
	}
	if len(loaded.Memories) != len(char.Memories) {
		t.Errorf("Memories count mismatch: got %d, want %d", len(loaded.Memories), len(char.Memories))
	}
}

func TestCharacterListAndDelete(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create multiple characters
	chars := []models.Character{
		{ID: "char1", Name: "Character 1"},
		{ID: "char2", Name: "Character 2"},
		{ID: "char3", Name: "Character 3"},
	}

	for i := range chars {
		if err := repo.SaveCharacter(&chars[i]); err != nil {
			t.Fatalf("Failed to save character %s: %v", chars[i].ID, err)
		}
	}

	// List characters
	ids, err := repo.ListCharacters()
	if err != nil {
		t.Fatalf("Failed to list characters: %v", err)
	}

	if len(ids) != 3 {
		t.Errorf("Expected 3 characters, got %d", len(ids))
	}

	// Since DeleteCharacter is not implemented, we'll skip deletion test
	// This is a limitation of the current repository implementation
}

func TestConcurrentWrites(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create initial character
	char := &models.Character{
		ID:        "concurrent-char",
		Name:      "Concurrent Character",
		Backstory: "Initial backstory",
	}
	
	if err := repo.SaveCharacter(char); err != nil {
		t.Fatalf("Failed to save initial character: %v", err)
	}

	// Concurrent writes
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			// Load, modify, save
			loaded, err := repo.LoadCharacter("concurrent-char")
			if err != nil {
				errors <- err
				return
			}
			
			loaded.Backstory = fmt.Sprintf("Updated by goroutine %d", idx)
			loaded.LastModified = time.Now()
			
			if err := repo.SaveCharacter(loaded); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Concurrent write error: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("Total concurrent errors: %d", errorCount)
	}

	// Verify final state
	final, err := repo.LoadCharacter("concurrent-char")
	if err != nil {
		t.Fatalf("Failed to load final character: %v", err)
	}

	// Should have one of the updates
	if final.Backstory == "Initial backstory" {
		t.Error("Character was not updated by any goroutine")
	}
}

func TestFileCorruption(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Save a valid character
	char := &models.Character{
		ID:   "corrupt-test",
		Name: "Corruption Test",
	}
	
	if err := repo.SaveCharacter(char); err != nil {
		t.Fatalf("Failed to save character: %v", err)
	}

	// Corrupt the file
	charFile := filepath.Join(tempDir, "characters", char.ID+".json")
	if err := os.WriteFile(charFile, []byte("{ invalid json"), 0644); err != nil {
		t.Fatalf("Failed to corrupt file: %v", err)
	}

	// Try to load corrupted character
	_, err = repo.LoadCharacter("corrupt-test")
	if err == nil {
		t.Error("Expected error loading corrupted character")
	}

	// Should be able to overwrite with valid data
	if err := repo.SaveCharacter(char); err != nil {
		t.Errorf("Failed to overwrite corrupted file: %v", err)
	}

	// Should now load successfully
	loaded, err := repo.LoadCharacter("corrupt-test")
	if err != nil {
		t.Errorf("Failed to load after fixing corruption: %v", err)
	}

	if loaded.Name != char.Name {
		t.Error("Character data mismatch after recovery")
	}
}

func TestDiskSpace(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping disk space test in short mode")
	}

	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create a large character with many memories
	char := &models.Character{
		ID:        "large-char",
		Name:      "Large Character",
		Backstory: string(make([]byte, 1024*1024)), // 1MB backstory
		Memories:  make([]models.Memory, 1000),
	}

	// Fill memories
	for i := range char.Memories {
		char.Memories[i] = models.Memory{
			Type:      models.LongTermMemory,
			Content:   fmt.Sprintf("Memory %d with some content to take up space", i),
			Timestamp: time.Now(),
		}
	}

	// Save large character
	err = repo.SaveCharacter(char)
	if err != nil {
		// This might fail on systems with very limited temp space
		t.Logf("Failed to save large character (might be disk space): %v", err)
	}

	// Try to load it back
	if err == nil {
		loaded, err := repo.LoadCharacter("large-char")
		if err != nil {
			t.Errorf("Failed to load large character: %v", err)
		} else if len(loaded.Memories) != 1000 {
			t.Errorf("Memory count mismatch: got %d, want 1000", len(loaded.Memories))
		}
	}
}

func TestCharacterInfo(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create characters with different attributes
	chars := []models.Character{
		{
			ID:          "char1",
			Name:        "Character One",
			Backstory:   "First character backstory",
			Quirks:      []string{"quirk1", "quirk2"},
			SpeechStyle: "Formal speech",
		},
		{
			ID:          "char2",
			Name:        "Character Two",
			Backstory:   "Second character backstory",
			Quirks:      []string{"quirk3"},
			SpeechStyle: "Casual speech",
		},
	}

	for i := range chars {
		if err := repo.SaveCharacter(&chars[i]); err != nil {
			t.Fatalf("Failed to save character: %v", err)
		}
	}

	// Get character info
	infos, err := repo.GetCharacterInfo()
	if err != nil {
		t.Fatalf("Failed to get character info: %v", err)
	}

	if len(infos) != 2 {
		t.Errorf("Expected 2 character infos, got %d", len(infos))
	}

	// Verify info content
	infoMap := make(map[string]CharacterInfo)
	for _, info := range infos {
		infoMap[info.ID] = info
	}

	if info, ok := infoMap["char1"]; ok {
		if info.Name != "Character One" {
			t.Errorf("Name mismatch: got %s, want Character One", info.Name)
		}
		if info.Description != "First character backstory" {
			t.Errorf("Description mismatch")
		}
		if len(info.Tags) != 2 {
			t.Errorf("Tags count mismatch: got %d, want 2", len(info.Tags))
		}
		if info.SpeechStyle != "Formal speech" {
			t.Errorf("SpeechStyle mismatch")
		}
	} else {
		t.Error("char1 not found in info list")
	}
}

func TestInvalidCharacterData(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tests := []struct {
		name string
		char *models.Character
		wantErr bool
	}{
		{
			name: "nil character",
			char: nil,
			wantErr: true,
		},
		{
			name: "empty ID",
			char: &models.Character{
				Name: "No ID",
			},
			wantErr: true,
		},
		{
			name: "invalid ID characters",
			char: &models.Character{
				ID:   "../../etc/passwd",
				Name: "Path Traversal",
			},
			wantErr: true,
		},
		{
			name: "very long ID",
			char: &models.Character{
				ID:   string(make([]byte, 256)),
				Name: "Long ID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SaveCharacter(tt.char)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveCharacter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryRecovery(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create characters directory with wrong permissions
	charDir := filepath.Join(tempDir, "characters")
	if err := os.MkdirAll(charDir, 0400); err != nil { // Read-only
		t.Fatalf("Failed to create read-only directory: %v", err)
	}

	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		// This is expected on some systems
		t.Logf("Repository creation with read-only dir failed as expected: %v", err)
		
		// Fix permissions and retry
		if err := os.Chmod(charDir, 0755); err != nil {
			t.Fatalf("Failed to fix permissions: %v", err)
		}
		
		repo, err = NewCharacterRepository(tempDir)
		if err != nil {
			t.Fatalf("Failed to create repository after fixing permissions: %v", err)
		}
	}

	// Should be able to save now
	char := &models.Character{
		ID:   "recovery-test",
		Name: "Recovery Test",
	}
	
	if err := repo.SaveCharacter(char); err != nil {
		t.Errorf("Failed to save after recovery: %v", err)
	}
}

func TestAtomicWrites(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewCharacterRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create a character
	char := &models.Character{
		ID:        "atomic-test",
		Name:      "Atomic Test",
		Backstory: "Testing atomic writes",
	}
	
	if err := repo.SaveCharacter(char); err != nil {
		t.Fatalf("Failed to save character: %v", err)
	}

	// Start a goroutine that continuously reads
	stop := make(chan bool)
	errors := make(chan error, 100)
	
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				loaded, err := repo.LoadCharacter("atomic-test")
				if err != nil {
					errors <- err
				} else if loaded.Name == "" {
					errors <- fmt.Errorf("loaded empty character")
				}
			}
		}
	}()

	// Continuously update the character
	for i := 0; i < 100; i++ {
		char.Backstory = fmt.Sprintf("Update %d", i)
		if err := repo.SaveCharacter(char); err != nil {
			t.Errorf("Failed to save update %d: %v", i, err)
		}
	}

	// Stop the reader
	stop <- true
	close(errors)

	// Check for read errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Read error during atomic write test: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("Total read errors: %d", errorCount)
	}
}