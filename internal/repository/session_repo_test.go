package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSessionPersistence(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	// Create a session
	session := &Session{
		ID:           "test-session",
		CharacterID:  "test-char",
		UserID:       "test-user",
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Messages: []SessionMessage{
			{
				Timestamp: time.Now(),
				Role:      "user",
				Content:   "Hello",
			},
			{
				Timestamp:  time.Now(),
				Role:       "character",
				Content:    "Hi there!",
				TokensUsed: 50,
			},
		},
		CacheMetrics: CacheMetrics{
			TotalRequests: 2,
			CacheHits:     1,
			CacheMisses:   1,
			TokensSaved:   25,
			HitRate:       0.5,
			CostSaved:     0.00075,
		},
	}

	// Save session
	err := repo.SaveSession(session)
	if err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Verify file exists
	sessionFile := filepath.Join(tempDir, "sessions", session.CharacterID, session.ID+".json")
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		t.Error("Session file was not created")
	}

	// Load session
	loaded, err := repo.LoadSession(session.CharacterID, session.ID)
	if err != nil {
		t.Fatalf("Failed to load session: %v", err)
	}

	// Verify fields
	if loaded.ID != session.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, session.ID)
	}
	if loaded.CharacterID != session.CharacterID {
		t.Errorf("CharacterID mismatch: got %s, want %s", loaded.CharacterID, session.CharacterID)
	}
	if loaded.UserID != session.UserID {
		t.Errorf("UserID mismatch: got %s, want %s", loaded.UserID, session.UserID)
	}
	if len(loaded.Messages) != len(session.Messages) {
		t.Errorf("Message count mismatch: got %d, want %d", len(loaded.Messages), len(session.Messages))
	}
	if loaded.CacheMetrics.TotalRequests != session.CacheMetrics.TotalRequests {
		t.Errorf("Cache metrics mismatch")
	}
}

func TestSessionList(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	// Create multiple sessions for different characters
	sessions := []*Session{
		{
			ID:           "session1",
			CharacterID:  "char1",
			UserID:       "user1",
			StartTime:    time.Now().Add(-2 * time.Hour),
			LastActivity: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:           "session2",
			CharacterID:  "char1",
			UserID:       "user2",
			StartTime:    time.Now().Add(-1 * time.Hour),
			LastActivity: time.Now().Add(-30 * time.Minute),
		},
		{
			ID:           "session3",
			CharacterID:  "char2",
			UserID:       "user1",
			StartTime:    time.Now().Add(-30 * time.Minute),
			LastActivity: time.Now(),
		},
	}

	for _, session := range sessions {
		if err := repo.SaveSession(session); err != nil {
			t.Fatalf("Failed to save session %s: %v", session.ID, err)
		}
	}

	// List sessions for char1
	char1Sessions, err := repo.ListSessions("char1")
	if err != nil {
		t.Fatalf("Failed to list sessions for char1: %v", err)
	}

	if len(char1Sessions) != 2 {
		t.Errorf("Expected 2 sessions for char1, got %d", len(char1Sessions))
	}

	// Verify sessions are sorted by LastActivity (most recent first)
	if len(char1Sessions) >= 2 {
		if char1Sessions[0].LastActivity.Before(char1Sessions[1].LastActivity) {
			t.Error("Sessions not sorted by LastActivity")
		}
	}

	// List sessions for char2
	char2Sessions, err := repo.ListSessions("char2")
	if err != nil {
		t.Fatalf("Failed to list sessions for char2: %v", err)
	}

	if len(char2Sessions) != 1 {
		t.Errorf("Expected 1 session for char2, got %d", len(char2Sessions))
	}

	// List sessions for non-existent character
	noSessions, err := repo.ListSessions("nonexistent")
	if err != nil {
		t.Errorf("Unexpected error listing sessions for non-existent character: %v", err)
	}

	if len(noSessions) != 0 {
		t.Errorf("Expected 0 sessions for non-existent character, got %d", len(noSessions))
	}
}

// TestSessionStats is commented out as GetSessionStats is not implemented
// This would test aggregate statistics across multiple sessions

func TestConcurrentSessionWrites(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	// Create initial session
	session := &Session{
		ID:           "concurrent-session",
		CharacterID:  "test-char",
		UserID:       "test-user",
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Messages:     []SessionMessage{},
	}

	if err := repo.SaveSession(session); err != nil {
		t.Fatalf("Failed to save initial session: %v", err)
	}

	// Concurrent updates
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Load session
			loaded, err := repo.LoadSession("test-char", "concurrent-session")
			if err != nil {
				errors <- err
				return
			}

			// Add message
			loaded.Messages = append(loaded.Messages, SessionMessage{
				Timestamp: time.Now(),
				Role:      "user",
				Content:   fmt.Sprintf("Message %d", idx),
			})
			loaded.LastActivity = time.Now()

			// Save session
			if err := repo.SaveSession(loaded); err != nil {
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
	final, err := repo.LoadSession("test-char", "concurrent-session")
	if err != nil {
		t.Fatalf("Failed to load final session: %v", err)
	}

	// Should have at least some messages (exact count depends on race conditions)
	if len(final.Messages) == 0 {
		t.Error("No messages were saved")
	}
}

func TestSessionCorruption(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	// Save a valid session
	session := &Session{
		ID:          "corrupt-test",
		CharacterID: "test-char",
		UserID:      "test-user",
		StartTime:   time.Now(),
		Messages: []SessionMessage{
			{
				Timestamp: time.Now(),
				Role:      "user",
				Content:   "Test message",
			},
		},
	}

	if err := repo.SaveSession(session); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Corrupt the file
	sessionFile := filepath.Join(tempDir, "sessions", session.CharacterID, session.ID+".json")
	if err := os.WriteFile(sessionFile, []byte("{ corrupt json"), 0644); err != nil {
		t.Fatalf("Failed to corrupt file: %v", err)
	}

	// Try to load corrupted session
	_, err := repo.LoadSession("test-char", "corrupt-test")
	if err == nil {
		t.Error("Expected error loading corrupted session")
	}

	// Should be able to overwrite with valid data
	if err := repo.SaveSession(session); err != nil {
		t.Errorf("Failed to overwrite corrupted file: %v", err)
	}

	// Should now load successfully
	loaded, err := repo.LoadSession("test-char", "corrupt-test")
	if err != nil {
		t.Errorf("Failed to load after fixing corruption: %v", err)
	}

	if len(loaded.Messages) != 1 {
		t.Error("Session data mismatch after recovery")
	}
}

func TestLargeSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large session test in short mode")
	}

	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	// Create a large session with many messages
	session := &Session{
		ID:          "large-session",
		CharacterID: "test-char",
		UserID:      "test-user",
		StartTime:   time.Now().Add(-24 * time.Hour),
		Messages:    make([]SessionMessage, 10000),
	}

	// Fill messages
	for i := range session.Messages {
		session.Messages[i] = SessionMessage{
			Timestamp:  time.Now().Add(time.Duration(i) * time.Second),
			Role:       "user",
			Content:    fmt.Sprintf("Message %d with some content to simulate a real conversation", i),
			TokensUsed: 50,
		}
	}

	// Save large session
	start := time.Now()
	err := repo.SaveSession(session)
	saveTime := time.Since(start)
	
	if err != nil {
		t.Fatalf("Failed to save large session: %v", err)
	}

	t.Logf("Saved 10,000 messages in %v", saveTime)

	// Load it back
	start = time.Now()
	loaded, err := repo.LoadSession("test-char", "large-session")
	loadTime := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to load large session: %v", err)
	}

	t.Logf("Loaded 10,000 messages in %v", loadTime)

	if len(loaded.Messages) != 10000 {
		t.Errorf("Message count mismatch: got %d, want 10000", len(loaded.Messages))
	}
}

func TestSessionResume(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	// Create a session
	originalSession := &Session{
		ID:           "resume-test",
		CharacterID:  "test-char",
		UserID:       "test-user",
		StartTime:    time.Now().Add(-1 * time.Hour),
		LastActivity: time.Now().Add(-30 * time.Minute),
		Messages: []SessionMessage{
			{
				Timestamp: time.Now().Add(-45 * time.Minute),
				Role:      "user",
				Content:   "First message",
			},
			{
				Timestamp: time.Now().Add(-40 * time.Minute),
				Role:      "character",
				Content:   "First response",
			},
		},
		CacheMetrics: CacheMetrics{
			TotalRequests: 1,
			CacheHits:     0,
			CacheMisses:   1,
		},
	}

	if err := repo.SaveSession(originalSession); err != nil {
		t.Fatalf("Failed to save original session: %v", err)
	}

	// Simulate resuming the session
	resumed, err := repo.LoadSession("test-char", "resume-test")
	if err != nil {
		t.Fatalf("Failed to load session for resume: %v", err)
	}

	// Add new messages
	resumed.Messages = append(resumed.Messages, SessionMessage{
		Timestamp: time.Now(),
		Role:      "user",
		Content:   "Resumed message",
	})
	resumed.LastActivity = time.Now()
	resumed.CacheMetrics.TotalRequests++
	resumed.CacheMetrics.CacheHits++

	// Save resumed session
	if err := repo.SaveSession(resumed); err != nil {
		t.Fatalf("Failed to save resumed session: %v", err)
	}

	// Load again to verify
	final, err := repo.LoadSession("test-char", "resume-test")
	if err != nil {
		t.Fatalf("Failed to load final session: %v", err)
	}

	if len(final.Messages) != 3 {
		t.Errorf("Expected 3 messages after resume, got %d", len(final.Messages))
	}

	if final.CacheMetrics.TotalRequests != 2 {
		t.Errorf("Expected 2 total requests, got %d", final.CacheMetrics.TotalRequests)
	}

	if final.CacheMetrics.CacheHits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", final.CacheMetrics.CacheHits)
	}
}

func TestInvalidSessionData(t *testing.T) {
	tempDir := t.TempDir()
	repo := NewSessionRepository(tempDir)

	tests := []struct {
		name    string
		session *Session
		wantErr bool
	}{
		{
			name:    "nil session",
			session: nil,
			wantErr: true,
		},
		{
			name: "empty ID",
			session: &Session{
				CharacterID: "test-char",
				UserID:      "test-user",
			},
			wantErr: true,
		},
		{
			name: "empty character ID",
			session: &Session{
				ID:     "test-session",
				UserID: "test-user",
			},
			wantErr: true,
		},
		{
			name: "invalid ID characters",
			session: &Session{
				ID:          "../../../etc/passwd",
				CharacterID: "test-char",
				UserID:      "test-user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SaveSession(tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}