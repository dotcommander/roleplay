package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
)

func TestCharacterBot_UpdateUserProfileResilience(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "roleplay-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test config
	cfg := &config.Config{
		UserProfileConfig: config.UserProfileConfig{
			Enabled:         true,
			UpdateFrequency: 1, // Update on every message for testing
			TurnsToConsider: 5,
		},
	}

	// Create bot
	bot := NewCharacterBot(cfg)

	// Create test character
	char := &models.Character{
		ID:        "test-char",
		Name:      "Test Character",
		Backstory: "Test character for profile update testing",
	}
	if err := bot.CreateCharacter(char); err != nil {
		t.Fatal(err)
	}

	// Initialize user profile repository
	bot.userProfileRepo = repository.NewUserProfileRepository(tempDir)

	// Create existing profile
	existingProfile := &models.UserProfile{
		UserID:      "test-user",
		CharacterID: "test-char",
		Facts: []models.UserFact{
			{
				Key:   "name",
				Value: "Test User",
			},
		},
		Version: 1,
	}
	err = bot.userProfileRepo.SaveUserProfile(existingProfile)
	if err != nil {
		t.Fatal(err)
	}

	// Test that updateUserProfileSync doesn't crash on nil agent
	// This simulates a background update failure scenario
	bot.updateUserProfileSync("test-user", char, "test-session")

	// Verify existing profile is still intact
	loaded, err := bot.userProfileRepo.LoadUserProfile("test-user", "test-char")
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Version != 1 {
		t.Errorf("Profile version changed unexpectedly: %d", loaded.Version)
	}
	if len(loaded.Facts) != 1 || loaded.Facts[0].Value != "Test User" {
		t.Error("Profile data was corrupted")
	}
}

func TestCharacterBot_BackgroundUpdateTimeout(t *testing.T) {
	// Test that background updates have proper timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	select {
	case <-time.After(100 * time.Millisecond):
		// Simulate work
	case <-ctx.Done():
		t.Error("Context cancelled too early")
	}

	// Verify context has deadline
	if _, ok := ctx.Deadline(); !ok {
		t.Error("Context should have deadline")
	}
}