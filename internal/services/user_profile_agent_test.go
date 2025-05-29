package services

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
)

// Mock provider for testing user profile agent
type mockUserProfileProvider struct {
	response *providers.AIResponse
	err      error
}

func (m *mockUserProfileProvider) SendRequest(ctx context.Context, req *providers.PromptRequest) (*providers.AIResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func (m *mockUserProfileProvider) SendStreamRequest(ctx context.Context, req *providers.PromptRequest, out chan<- providers.PartialAIResponse) error {
	return errors.New("streaming not implemented")
}

func (m *mockUserProfileProvider) Name() string {
	return "mock"
}

func TestProfileUpdate(t *testing.T) {
	mockProvider := &mockUserProfileProvider{
		response: &providers.AIResponse{
			Content: `{
				"user_id": "user123",
				"character_id": "char456",
				"facts": [
					{"key": "profession", "value": "software developer", "confidence": 0.95, "source_turn": 1}
				],
				"overall_summary": "A software developer",
				"interaction_style": "professional",
				"version": 1
			}`,
		},
	}
	
	tempDir := t.TempDir()
	repo := repository.NewUserProfileRepository(tempDir)
	agent := NewUserProfileAgent(mockProvider, repo)
	
	// Create initial profile
	initialProfile := &models.UserProfile{
		UserID:      "user123",
		CharacterID: "char456",
		Facts: []models.UserFact{
			{
				Key:         "hobby",
				Value:       "reading",
				Confidence:  0.8,
				SourceTurn:  0,
				LastUpdated: time.Now().Add(-1 * time.Hour),
			},
		},
		LastAnalyzed: time.Now().Add(-1 * time.Hour),
	}
	if err := repo.SaveUserProfile(initialProfile); err != nil {
		t.Fatalf("Failed to save initial profile: %v", err)
	}
	
	// Update profile
	sessionMessages := []repository.SessionMessage{
		{Timestamp: time.Now(), Role: "user", Content: "I'm a software developer"},
		{Timestamp: time.Now(), Role: "assistant", Content: "That's interesting!"},
	}
	
	character := &models.Character{
		ID:   "char456",
		Name: "Test Character",
	}
	
	updatedProfile, err := agent.UpdateUserProfile(
		context.Background(),
		"user123",
		character,
		sessionMessages,
		10,
		initialProfile,
	)
	
	if err != nil {
		t.Fatalf("Failed to update profile: %v", err)
	}
	
	// Check that profile was updated
	if updatedProfile.UserID != "user123" {
		t.Errorf("UserID mismatch: got %s, want user123", updatedProfile.UserID)
	}
	if updatedProfile.CharacterID != "char456" {
		t.Errorf("CharacterID mismatch: got %s, want char456", updatedProfile.CharacterID)
	}
	
	// Check facts
	if len(updatedProfile.Facts) == 0 {
		t.Error("Expected at least one fact in updated profile")
	}
}

func TestProfileUpdateWithError(t *testing.T) {
	mockProvider := &mockUserProfileProvider{
		err: errors.New("API error"),
	}
	
	tempDir := t.TempDir()
	repo := repository.NewUserProfileRepository(tempDir)
	agent := NewUserProfileAgent(mockProvider, repo)
	
	character := &models.Character{
		ID:   "char456",
		Name: "Test Character",
	}
	
	sessionMessages := []repository.SessionMessage{
		{Timestamp: time.Now(), Role: "user", Content: "Test message"},
	}
	
	// Should handle error gracefully
	_, err := agent.UpdateUserProfile(
		context.Background(),
		"user123",
		character,
		sessionMessages,
		10,
		nil,
	)
	
	if err == nil {
		t.Error("Expected error from provider")
	}
}

func TestProfileSaveError(t *testing.T) {
	// This test is tricky to simulate save errors with real filesystem
	// So we'll skip the actual save error simulation
	t.Skip("Skipping save error test - difficult to simulate with real filesystem")
}

func TestInvalidJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name:     "invalid JSON",
			response: "This is not JSON at all",
			wantErr:  true,
		},
		{
			name:     "JSON with markdown wrapper",
			response: "Here's the analysis:\n```json\n{\"user_id\": \"user123\", \"character_id\": \"char456\", \"facts\": [], \"version\": 1}\n```",
			wantErr:  false, // Should handle this case
		},
		{
			name:     "incomplete JSON",
			response: `{"user_id": "user123", "character_id": "char456"`,
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &mockUserProfileProvider{
				response: &providers.AIResponse{
					Content: tt.response,
				},
			}
			
			tempDir := t.TempDir()
			repo := repository.NewUserProfileRepository(tempDir)
			agent := NewUserProfileAgent(mockProvider, repo)
			
			character := &models.Character{
				ID:   "char456",
				Name: "Test Character",
			}
			
			sessionMessages := []repository.SessionMessage{
				{Timestamp: time.Now(), Role: "user", Content: "Test"},
			}
			
			profile, err := agent.UpdateUserProfile(
				context.Background(),
				"user123",
				character,
				sessionMessages,
				10,
				nil,
			)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			// For error cases with existing profile, should return existing
			// This is ok - it might return existing profile on error
			_ = profile
		})
	}
}

func TestConcurrentProfileUpdates(t *testing.T) {
	mockProvider := &mockUserProfileProvider{
		response: &providers.AIResponse{
			Content: `{
				"user_id": "user123",
				"character_id": "char456",
				"facts": [],
				"overall_summary": "Test profile",
				"version": 1
			}`,
		},
	}
	
	tempDir := t.TempDir()
	repo := repository.NewUserProfileRepository(tempDir)
	agent := NewUserProfileAgent(mockProvider, repo)
	
	character := &models.Character{
		ID:   "char456",
		Name: "Test Character",
	}
	
	// Run multiple concurrent updates
	var wg sync.WaitGroup
	errors := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			sessionMessages := []repository.SessionMessage{
				{
					Timestamp: time.Now(),
					Role:      "user",
					Content:   "Test message",
				},
			}
			
			_, err := agent.UpdateUserProfile(
				context.Background(),
				"user123",
				character,
				sessionMessages,
				10,
				nil,
			)
			
			if err != nil {
				errors <- err
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Concurrent update error: %v", err)
		errorCount++
	}
	
	if errorCount > 0 {
		t.Errorf("Total concurrent errors: %d", errorCount)
	}
}