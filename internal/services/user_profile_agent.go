package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/utils"
)

// UserProfileAgent handles AI-powered user profile extraction and updates
type UserProfileAgent struct {
	provider   providers.AIProvider
	repo       *repository.UserProfileRepository
	promptPath string
}

// NewUserProfileAgent creates a new user profile agent
func NewUserProfileAgent(provider providers.AIProvider, repo *repository.UserProfileRepository) *UserProfileAgent {
	// Find prompt file relative to executable
	promptFile := "prompts/user-profile-extraction.md"

	// Try multiple locations for the prompt file
	possiblePaths := []string{
		promptFile,
		filepath.Join(".", promptFile),
		filepath.Join("..", promptFile),
		filepath.Join(os.Getenv("HOME"), "go", "src", "roleplay", promptFile),
	}

	var finalPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			finalPath = path
			break
		}
	}

	if finalPath == "" {
		// Fallback to the first option
		finalPath = promptFile
	}

	return &UserProfileAgent{
		provider:   provider,
		repo:       repo,
		promptPath: finalPath,
	}
}

type profilePromptData struct {
	ExistingProfileJSON string
	HistoryTurnCount    int
	Messages            []profileMessageData
	CharacterName       string
	CharacterID         string
	UserID              string
	NextVersion         int
	CurrentTimestamp    string
}

type profileMessageData struct {
	Role       string
	Content    string
	TurnNumber int
	Timestamp  string
}

// UpdateUserProfile analyzes conversation history and updates the user profile
func (upa *UserProfileAgent) UpdateUserProfile(
	ctx context.Context,
	userID string,
	character *models.Character,
	sessionMessages []repository.SessionMessage,
	turnsToConsider int,
	currentProfile *models.UserProfile,
) (*models.UserProfile, error) {

	// If no current profile provided, create a new one
	if currentProfile == nil {
		currentProfile = &models.UserProfile{
			UserID:      userID,
			CharacterID: character.ID,
			Facts:       []models.UserFact{},
			Version:     0,
		}
	}

	// Use existingProfile to refer to currentProfile for consistency with rest of code
	existingProfile := currentProfile

	if len(sessionMessages) == 0 {
		// No conversation history - return existing profile unchanged
		return existingProfile, nil
	}

	existingProfileJSON, _ := json.MarshalIndent(existingProfile, "", "  ")

	// Prepare recent conversation history
	startIndex := 0
	if len(sessionMessages) > turnsToConsider {
		startIndex = len(sessionMessages) - turnsToConsider
	}
	recentHistory := sessionMessages[startIndex:]

	var promptMessages []profileMessageData
	for i, msg := range recentHistory {
		promptMessages = append(promptMessages, profileMessageData{
			Role:       msg.Role,
			Content:    msg.Content,
			TurnNumber: startIndex + i + 1,
			Timestamp:  msg.Timestamp.Format(time.RFC3339),
		})
	}

	// Load and parse prompt template
	promptTemplateBytes, err := os.ReadFile(upa.promptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user profile prompt template '%s': %w", upa.promptPath, err)
	}

	tmpl, err := template.New("userProfile").Parse(string(promptTemplateBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse user profile prompt template: %w", err)
	}

	data := profilePromptData{
		ExistingProfileJSON: string(existingProfileJSON),
		Messages:            promptMessages,
		HistoryTurnCount:    len(promptMessages),
		CharacterName:       character.Name,
		CharacterID:         character.ID,
		UserID:              userID,
		NextVersion:         existingProfile.Version + 1,
		CurrentTimestamp:    time.Now().Format(time.RFC3339),
	}

	var promptBuilder strings.Builder
	if err := tmpl.Execute(&promptBuilder, data); err != nil {
		return nil, fmt.Errorf("failed to execute user profile prompt template: %w", err)
	}

	// Make LLM call
	request := &providers.PromptRequest{
		CharacterID:  "system-user-profiler",
		UserID:       userID,
		Message:      promptBuilder.String(),
		SystemPrompt: "You are an analytical AI. Your task is to extract and update user profile information based on the provided conversation history and existing profile. Respond ONLY with the updated JSON profile.",
	}

	response, err := upa.provider.SendRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("LLM call for user profile extraction failed: %w", err)
	}

	// Use robust JSON extraction to handle LLM output quirks
	extractedJSON, err := utils.ExtractValidJSON(response.Content)
	if err != nil {
		// Log the error but return existing profile to maintain stability
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE: Failed to extract valid JSON from LLM response for %s/%s.\nError: %v\nRaw content (first 500 chars): %s\n", 
			userID, character.ID, err, truncateString(response.Content, 500))
		return existingProfile, fmt.Errorf("failed to extract valid JSON for user profile update: %w", err)
	}

	var updatedProfile models.UserProfile
	if err := json.Unmarshal([]byte(extractedJSON), &updatedProfile); err != nil {
		// This should be rare now that we're using ExtractValidJSON
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE: Failed to parse extracted JSON for %s/%s.\nError: %v\nExtracted JSON (first 500 chars): %s\n", 
			userID, character.ID, err, truncateString(extractedJSON, 500))
		return existingProfile, fmt.Errorf("failed to parse extracted JSON for user profile update: %w", err)
	}

	// Validate the response
	if updatedProfile.UserID != userID || updatedProfile.CharacterID != character.ID {
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE: LLM returned profile for incorrect user/character. Expected %s/%s, got %s/%s\n",
			userID, character.ID, updatedProfile.UserID, updatedProfile.CharacterID)
		return existingProfile, fmt.Errorf("LLM returned profile for incorrect user/character")
	}

	// Save the updated profile
	if err := upa.repo.SaveUserProfile(&updatedProfile); err != nil {
		// Log save failure but return the in-memory updated profile
		fmt.Fprintf(os.Stderr, "BACKGROUND PROFILE UPDATE: Updated user profile for %s/%s in memory but FAILED TO SAVE: %v\n",
			userID, character.ID, err)
		return &updatedProfile, fmt.Errorf("updated user profile in memory but FAILED TO SAVE: %w", err)
	}

	return &updatedProfile, nil
}

// truncateString truncates a string to maxLen characters for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
