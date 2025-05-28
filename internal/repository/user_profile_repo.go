package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotcommander/roleplay/internal/models"
)

// UserProfileRepository manages persistence of user profiles
type UserProfileRepository struct {
	dataDir string
}

// NewUserProfileRepository creates a new repository instance
func NewUserProfileRepository(dataDir string) *UserProfileRepository {
	return &UserProfileRepository{
		dataDir: dataDir,
	}
}

// profileFilename generates the filename for a user profile
func (r *UserProfileRepository) profileFilename(userID, characterID string) string {
	return fmt.Sprintf("%s_%s.json", userID, characterID)
}

// SaveUserProfile saves a user profile to disk
func (r *UserProfileRepository) SaveUserProfile(profile *models.UserProfile) error {
	if err := os.MkdirAll(r.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create user profiles directory: %w", err)
	}

	filename := r.profileFilename(profile.UserID, profile.CharacterID)
	filepath := filepath.Join(r.dataDir, filename)

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user profile: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write user profile file: %w", err)
	}

	return nil
}

// LoadUserProfile loads a user profile from disk
func (r *UserProfileRepository) LoadUserProfile(userID, characterID string) (*models.UserProfile, error) {
	filename := r.profileFilename(userID, characterID)
	filepath := filepath.Join(r.dataDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err // Let caller handle non-existence
		}
		return nil, fmt.Errorf("failed to read user profile file: %w", err)
	}

	var profile models.UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user profile: %w", err)
	}

	return &profile, nil
}

// DeleteUserProfile deletes a user profile from disk
func (r *UserProfileRepository) DeleteUserProfile(userID, characterID string) error {
	filename := r.profileFilename(userID, characterID)
	filepath := filepath.Join(r.dataDir, filename)

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}

	return nil
}

// ListUserProfiles returns all user profiles for a given user
func (r *UserProfileRepository) ListUserProfiles(userID string) ([]*models.UserProfile, error) {
	pattern := filepath.Join(r.dataDir, fmt.Sprintf("%s_*.json", userID))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list user profiles: %w", err)
	}

	var profiles []*models.UserProfile
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}

		var profile models.UserProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			continue // Skip invalid JSON files
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}
