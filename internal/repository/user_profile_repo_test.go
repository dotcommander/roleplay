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

func TestUserProfilePersistence(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Create a user profile
	profile := &models.UserProfile{
		UserID:      "test-user",
		CharacterID: "test-char",
		Facts: []models.UserFact{
			{
				Key:         "profession",
				Value:       "software engineer",
				Confidence:  0.95,
				SourceTurn:  1,
				LastUpdated: time.Now(),
			},
			{
				Key:         "hobby",
				Value:       "hiking",
				Confidence:  0.8,
				SourceTurn:  2,
				LastUpdated: time.Now(),
			},
		},
		OverallSummary:   "A software engineer who enjoys outdoor activities",
		InteractionStyle: "professional",
		LastAnalyzed:     time.Now(),
		Version:          1,
	}

	// Save profile
	err := repo.SaveUserProfile(profile)
	if err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Verify file exists
	profileFile := filepath.Join(tempDir, fmt.Sprintf("%s_%s.json", profile.UserID, profile.CharacterID))
	if _, err := os.Stat(profileFile); os.IsNotExist(err) {
		t.Error("Profile file was not created")
	}

	// Load profile
	loaded, err := repo.LoadUserProfile(profile.UserID, profile.CharacterID)
	if err != nil {
		t.Fatalf("Failed to load profile: %v", err)
	}

	// Verify fields
	if loaded.UserID != profile.UserID {
		t.Errorf("UserID mismatch: got %s, want %s", loaded.UserID, profile.UserID)
	}
	if loaded.CharacterID != profile.CharacterID {
		t.Errorf("CharacterID mismatch: got %s, want %s", loaded.CharacterID, profile.CharacterID)
	}
	if len(loaded.Facts) != len(profile.Facts) {
		t.Errorf("Facts count mismatch: got %d, want %d", len(loaded.Facts), len(profile.Facts))
	}
	if loaded.OverallSummary != profile.OverallSummary {
		t.Errorf("Summary mismatch: got %s, want %s", loaded.OverallSummary, profile.OverallSummary)
	}

	// Verify fact details
	if len(loaded.Facts) > 0 {
		if loaded.Facts[0].Key != profile.Facts[0].Key {
			t.Errorf("Fact key mismatch: got %s, want %s", loaded.Facts[0].Key, profile.Facts[0].Key)
		}
		if loaded.Facts[0].Value != profile.Facts[0].Value {
			t.Errorf("Fact value mismatch: got %s, want %s", loaded.Facts[0].Value, profile.Facts[0].Value)
		}
		if loaded.Facts[0].Confidence != profile.Facts[0].Confidence {
			t.Errorf("Fact confidence mismatch: got %f, want %f", loaded.Facts[0].Confidence, profile.Facts[0].Confidence)
		}
	}
}

func TestUserProfileList(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Create profiles for different user-character combinations
	profiles := []*models.UserProfile{
		{
			UserID:      "user1",
			CharacterID: "char1",
			OverallSummary:     "User1 with Char1",
			LastAnalyzed:   time.Now(),
		},
		{
			UserID:      "user1",
			CharacterID: "char2",
			OverallSummary:     "User1 with Char2",
			LastAnalyzed:   time.Now(),
		},
		{
			UserID:      "user2",
			CharacterID: "char1",
			OverallSummary:     "User2 with Char1",
			LastAnalyzed:   time.Now(),
		},
	}

	for _, profile := range profiles {
		if err := repo.SaveUserProfile(profile); err != nil {
			t.Fatalf("Failed to save profile: %v", err)
		}
	}

	// List profiles for user1
	user1Profiles, err := repo.ListUserProfiles("user1")
	if err != nil {
		t.Fatalf("Failed to list profiles for user1: %v", err)
	}

	if len(user1Profiles) != 2 {
		t.Errorf("Expected 2 profiles for user1, got %d", len(user1Profiles))
	}

	// Verify profile IDs
	foundChar1 := false
	foundChar2 := false
	for _, profile := range user1Profiles {
		if profile.CharacterID == "char1" {
			foundChar1 = true
		}
		if profile.CharacterID == "char2" {
			foundChar2 = true
		}
	}

	if !foundChar1 || !foundChar2 {
		t.Error("Not all expected profiles found for user1")
	}

	// List profiles for user2
	user2Profiles, err := repo.ListUserProfiles("user2")
	if err != nil {
		t.Fatalf("Failed to list profiles for user2: %v", err)
	}

	if len(user2Profiles) != 1 {
		t.Errorf("Expected 1 profile for user2, got %d", len(user2Profiles))
	}

	// List profiles for non-existent user
	noProfiles, err := repo.ListUserProfiles("nonexistent")
	if err != nil {
		t.Errorf("Unexpected error listing profiles for non-existent user: %v", err)
	}

	if len(noProfiles) != 0 {
		t.Errorf("Expected 0 profiles for non-existent user, got %d", len(noProfiles))
	}
}

func TestUserProfileDelete(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Create a profile
	profile := &models.UserProfile{
		UserID:      "delete-user",
		CharacterID: "delete-char",
		OverallSummary:     "Profile to be deleted",
		LastAnalyzed:   time.Now(),
	}

	if err := repo.SaveUserProfile(profile); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Verify it exists
	_, err := repo.LoadUserProfile("delete-user", "delete-char")
	if err != nil {
		t.Fatalf("Failed to load profile before delete: %v", err)
	}

	// Delete profile
	err = repo.DeleteUserProfile("delete-user", "delete-char")
	if err != nil {
		t.Errorf("Failed to delete profile: %v", err)
	}

	// Verify it's gone
	_, err = repo.LoadUserProfile("delete-user", "delete-char")
	if err == nil {
		t.Error("Expected error loading deleted profile")
	}

	// Delete non-existent profile should not error
	err = repo.DeleteUserProfile("nonexistent", "nonexistent")
	if err != nil {
		t.Errorf("Unexpected error deleting non-existent profile: %v", err)
	}
}

func TestUserProfileUpdate(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Create initial profile
	profile := &models.UserProfile{
		UserID:      "update-user",
		CharacterID: "update-char",
		Facts: []models.UserFact{
			{
				Key:        "initial",
				Value:      "Initial fact",
				Confidence: 0.7,
				SourceTurn: 1,
				LastUpdated: time.Now(),
			},
		},
		OverallSummary:     "Initial summary",
		LastAnalyzed:   time.Now().Add(-2 * time.Hour),
	}

	if err := repo.SaveUserProfile(profile); err != nil {
		t.Fatalf("Failed to save initial profile: %v", err)
	}

	// Update profile
	profile.Facts = append(profile.Facts, models.UserFact{
		Key:        "new",
		Value:      "New fact",
		Confidence: 0.9,
		SourceTurn: 2,
		LastUpdated: time.Now(),
	})
	profile.OverallSummary = "Updated summary"
	profile.LastAnalyzed = time.Now()

	if err := repo.SaveUserProfile(profile); err != nil {
		t.Fatalf("Failed to save updated profile: %v", err)
	}

	// Load and verify
	loaded, err := repo.LoadUserProfile("update-user", "update-char")
	if err != nil {
		t.Fatalf("Failed to load updated profile: %v", err)
	}

	if len(loaded.Facts) != 2 {
		t.Errorf("Expected 2 facts after update, got %d", len(loaded.Facts))
	}

	if loaded.OverallSummary != "Updated summary" {
		t.Errorf("Summary not updated: got %s, want %s", loaded.OverallSummary, "Updated summary")
	}

	// LastAnalyzed should be equal (we saved the same object)
	if loaded.LastAnalyzed.IsZero() {
		t.Error("LastAnalyzed should not be zero")
	}
}

func TestConcurrentProfileAccess(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Create initial profile
	profile := &models.UserProfile{
		UserID:      "concurrent-user",
		CharacterID: "concurrent-char",
		Facts:       []models.UserFact{},
		OverallSummary:     "Initial",
		LastAnalyzed:   time.Now(),
	}

	if err := repo.SaveUserProfile(profile); err != nil {
		t.Fatalf("Failed to save initial profile: %v", err)
	}

	// Concurrent reads and writes
	var wg sync.WaitGroup
	errors := make(chan error, 20)

	// 10 readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := repo.LoadUserProfile("concurrent-user", "concurrent-char")
			if err != nil {
				errors <- err
			}
		}()
	}

	// 10 writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			// Load, modify, save
			loaded, err := repo.LoadUserProfile("concurrent-user", "concurrent-char")
			if err != nil {
				errors <- err
				return
			}

			loaded.Facts = append(loaded.Facts, models.UserFact{
				Key:        fmt.Sprintf("fact%d", idx),
				Value:      fmt.Sprintf("Fact %d", idx),
				Confidence: 0.8,
				SourceTurn: idx,
				LastUpdated: time.Now(),
			})
			loaded.LastAnalyzed = time.Now()

			if err := repo.SaveUserProfile(loaded); err != nil {
				errors <- err
			}
		}(i)
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

	// Verify final state
	final, err := repo.LoadUserProfile("concurrent-user", "concurrent-char")
	if err != nil {
		t.Fatalf("Failed to load final profile: %v", err)
	}

	// Should have at least some facts (exact count depends on race conditions)
	if len(final.Facts) == 0 {
		t.Error("No facts were saved during concurrent access")
	}
}

func TestProfileFactDeduplication(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Create profile with duplicate facts
	profile := &models.UserProfile{
		UserID:      "dedup-user",
		CharacterID: "dedup-char",
		Facts: []models.UserFact{
			{
				Key:        "food_preference",
				Value:      "Likes pizza",
				Confidence: 0.8,
				SourceTurn: 1,
				LastUpdated: time.Now(),
			},
			{
				Key:        "food_preference",
				Value:      "Likes pizza", // Duplicate
				Confidence: 0.9,
				SourceTurn: 2,
				LastUpdated: time.Now(),
			},
			{
				Key:        "location",
				Value:      "Lives in New York",
				Confidence: 0.7,
				SourceTurn: 3,
				LastUpdated: time.Now(),
			},
		},
		LastAnalyzed: time.Now(),
	}

	// Save and load
	if err := repo.SaveUserProfile(profile); err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	loaded, err := repo.LoadUserProfile("dedup-user", "dedup-char")
	if err != nil {
		t.Fatalf("Failed to load profile: %v", err)
	}

	// Repository doesn't deduplicate - that's the service's job
	// Just verify all facts are saved
	if len(loaded.Facts) != 3 {
		t.Errorf("Expected 3 facts (including duplicates), got %d", len(loaded.Facts))
	}
}

func TestInvalidProfileData(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	tests := []struct {
		name    string
		profile *models.UserProfile
		wantErr bool
	}{
		{
			name:    "nil profile",
			profile: nil,
			wantErr: true,
		},
		{
			name: "empty user ID",
			profile: &models.UserProfile{
				CharacterID: "test-char",
			},
			wantErr: true,
		},
		{
			name: "empty character ID",
			profile: &models.UserProfile{
				UserID: "test-user",
			},
			wantErr: true,
		},
		{
			name: "invalid user ID characters",
			profile: &models.UserProfile{
				UserID:      "../../../etc",
				CharacterID: "test-char",
			},
			wantErr: true,
		},
		{
			name: "very long ID",
			profile: &models.UserProfile{
				UserID:      string(make([]byte, 256)),
				CharacterID: "test-char",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SaveUserProfile(tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUserProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileCorruption(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewUserProfileRepository(tempDir)

	// Save a valid profile
	profile := &models.UserProfile{
		UserID:      "corrupt-user",
		CharacterID: "corrupt-char",
		OverallSummary:     "Test profile",
		LastAnalyzed:   time.Now(),
	}

	err := repo.SaveUserProfile(profile)
	if err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Corrupt the file
	profileFile := filepath.Join(tempDir, fmt.Sprintf("%s_%s.json", profile.UserID, profile.CharacterID))
	if err := os.WriteFile(profileFile, []byte("{ invalid json"), 0644); err != nil {
		t.Fatalf("Failed to corrupt file: %v", err)
	}

	// Try to load corrupted profile
	_, err = repo.LoadUserProfile("corrupt-user", "corrupt-char")
	if err == nil {
		t.Error("Expected error loading corrupted profile")
	}

	// Should be able to overwrite with valid data
	if err := repo.SaveUserProfile(profile); err != nil {
		t.Errorf("Failed to overwrite corrupted file: %v", err)
	}

	// Should now load successfully
	loaded, err := repo.LoadUserProfile("corrupt-user", "corrupt-char")
	if err != nil {
		t.Errorf("Failed to load after fixing corruption: %v", err)
	}

	if loaded.OverallSummary != profile.OverallSummary {
		t.Error("Profile data mismatch after recovery")
	}
}