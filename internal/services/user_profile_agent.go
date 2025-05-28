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
) (*models.UserProfile, error) {

	if len(sessionMessages) == 0 {
		return nil, fmt.Errorf("no conversation history provided to update user profile")
	}

	// Load existing profile or create new one
	existingProfile, err := upa.repo.LoadUserProfile(userID, character.ID)
	if err != nil {
		if os.IsNotExist(err) {
			existingProfile = &models.UserProfile{
				UserID:      userID,
				CharacterID: character.ID,
				Facts:       []models.UserFact{},
				Version:     0,
			}
		} else {
			return nil, fmt.Errorf("failed to load existing user profile: %w", err)
		}
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

	// Clean and parse JSON response
	jsonContent := strings.TrimSpace(response.Content)

	// Remove markdown code blocks if present
	if strings.HasPrefix(jsonContent, "```json") {
		jsonContent = strings.TrimPrefix(jsonContent, "```json")
		jsonContent = strings.TrimSuffix(jsonContent, "```")
		jsonContent = strings.TrimSpace(jsonContent)
	} else if strings.HasPrefix(jsonContent, "```") {
		jsonContent = strings.TrimPrefix(jsonContent, "```")
		jsonContent = strings.TrimSuffix(jsonContent, "```")
		jsonContent = strings.TrimSpace(jsonContent)
	}

	var updatedProfile models.UserProfile
	if err := json.Unmarshal([]byte(jsonContent), &updatedProfile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse LLM response for user profile. Raw content:\n%s\n", jsonContent)
		return nil, fmt.Errorf("failed to parse LLM response as JSON for user profile: %w", err)
	}

	// Validate the response
	if updatedProfile.UserID != userID || updatedProfile.CharacterID != character.ID {
		return nil, fmt.Errorf("LLM returned profile for incorrect user/character. Expected %s/%s, got %s/%s",
			userID, character.ID, updatedProfile.UserID, updatedProfile.CharacterID)
	}

	// Save the updated profile
	if err := upa.repo.SaveUserProfile(&updatedProfile); err != nil {
		return nil, fmt.Errorf("failed to save updated user profile: %w", err)
	}

	return &updatedProfile, nil
}
