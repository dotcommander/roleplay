package manager

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
)

func TestCharacterLifecycle(t *testing.T) {
	// Create temp directory for testing
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		DefaultProvider: "openai",
		Model:           "gpt-3.5-turbo",
		APIKey:          "test-key",
	}
	
	mgr, err := NewCharacterManager(cfg)
	if err != nil {
		// This might fail if provider initialization fails, which is expected in tests
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	// Test character creation
	char := &models.Character{
		ID:        "test-char",
		Name:      "Test Character",
		Backstory: "A character for testing",
		Personality: models.PersonalityTraits{
			Openness: 0.5,
		},
	}
	
	err = mgr.CreateCharacter(char)
	if err != nil {
		t.Errorf("Failed to create character: %v", err)
	}
	
	// Test loading character
	loaded, err := mgr.GetOrLoadCharacter("test-char")
	if err != nil {
		t.Errorf("Failed to load character: %v", err)
	}
	
	if loaded.Name != char.Name {
		t.Errorf("Loaded character name mismatch: got %s, want %s", loaded.Name, char.Name)
	}
}

func TestConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		DefaultProvider: "openai",
		Model:           "gpt-3.5-turbo",
		APIKey:          "test-key",
	}
	
	mgr, err := NewCharacterManager(cfg)
	if err != nil {
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	// Create initial character
	char := &models.Character{
		ID:   "concurrent-test",
		Name: "Concurrent Test",
		Personality: models.PersonalityTraits{
			Openness: 0.5,
		},
	}
	if err := mgr.CreateCharacter(char); err != nil {
		t.Fatalf("Failed to create initial character: %v", err)
	}
	
	// Test concurrent reads
	var wg sync.WaitGroup
	errors := make(chan error, 50)
	
	// 50 readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := mgr.GetOrLoadCharacter("concurrent-test")
			if err != nil {
				errors <- err
			}
		}()
	}
	
	wg.Wait()
	close(errors)
	
	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
		errorCount++
	}
	
	if errorCount > 0 {
		t.Errorf("Total concurrent errors: %d", errorCount)
	}
}

func TestListAvailableCharacters(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		DefaultProvider: "openai",
		Model:           "gpt-3.5-turbo",
		APIKey:          "test-key",
	}
	
	mgr, err := NewCharacterManager(cfg)
	if err != nil {
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	// Create a few characters
	chars := []*models.Character{
		{
			ID:          "char1",
			Name:        "Character 1",
			Backstory:   "First character",
			Quirks:      []string{"quirk1"},
			SpeechStyle: "formal",
		},
		{
			ID:          "char2",
			Name:        "Character 2",
			Backstory:   "Second character",
			Quirks:      []string{"quirk2", "quirk3"},
			SpeechStyle: "casual",
		},
	}
	
	for _, char := range chars {
		if err := mgr.CreateCharacter(char); err != nil {
			t.Fatalf("Failed to create character %s: %v", char.ID, err)
		}
	}
	
	// List characters
	infos, err := mgr.ListAvailableCharacters()
	if err != nil {
		t.Fatalf("Failed to list characters: %v", err)
	}
	
	if len(infos) != 2 {
		t.Errorf("Expected 2 characters, got %d", len(infos))
	}
	
	// Verify character info
	foundChar1 := false
	foundChar2 := false
	for _, info := range infos {
		if info.ID == "char1" {
			foundChar1 = true
			if info.Name != "Character 1" {
				t.Errorf("Character 1 name mismatch: got %s", info.Name)
			}
		}
		if info.ID == "char2" {
			foundChar2 = true
			if info.Name != "Character 2" {
				t.Errorf("Character 2 name mismatch: got %s", info.Name)
			}
		}
	}
	
	if !foundChar1 || !foundChar2 {
		t.Error("Not all characters found in list")
	}
}

func TestLoadAllCharacters(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		DefaultProvider: "openai",
		Model:           "gpt-3.5-turbo",
		APIKey:          "test-key",
	}
	
	mgr, err := NewCharacterManager(cfg)
	if err != nil {
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	// Create characters directly in repository
	chars := []*models.Character{
		{ID: "load1", Name: "Load Test 1"},
		{ID: "load2", Name: "Load Test 2"},
		{ID: "load3", Name: "Load Test 3"},
	}
	
	for _, char := range chars {
		if err := mgr.CreateCharacter(char); err != nil {
			t.Fatalf("Failed to create character: %v", err)
		}
	}
	
	// Create new manager instance to test loading
	mgr2, err := NewCharacterManager(cfg)
	if err != nil {
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	// Load all characters
	err = mgr2.LoadAllCharacters()
	if err != nil {
		t.Errorf("Failed to load all characters: %v", err)
	}
	
	// Verify all characters are loaded
	for _, char := range chars {
		loaded, err := mgr2.GetOrLoadCharacter(char.ID)
		if err != nil {
			t.Errorf("Failed to get character %s: %v", char.ID, err)
		}
		if loaded.Name != char.Name {
			t.Errorf("Character name mismatch for %s: got %s, want %s", char.ID, loaded.Name, char.Name)
		}
	}
}

func TestGetBot(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		DefaultProvider: "openai",
		Model:           "gpt-3.5-turbo",
		APIKey:          "test-key",
	}
	
	mgr, err := NewCharacterManager(cfg)
	if err != nil {
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	bot := mgr.GetBot()
	if bot == nil {
		t.Error("GetBot() returned nil")
	}
}

func TestGetSessionRepository(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	cfg := &config.Config{
		CacheConfig: config.CacheConfig{
			CleanupInterval: 5 * time.Minute,
			DefaultTTL:      10 * time.Minute,
		},
		DefaultProvider: "openai",
		Model:           "gpt-3.5-turbo",
		APIKey:          "test-key",
	}
	
	mgr, err := NewCharacterManager(cfg)
	if err != nil {
		t.Skipf("Skipping test - provider initialization failed: %v", err)
	}
	
	sessions := mgr.GetSessionRepository()
	if sessions == nil {
		t.Error("GetSessionRepository() returned nil")
	}
}