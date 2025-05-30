package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
)

// Session represents a conversation session
type Session struct {
	ID           string           `json:"id"`
	CharacterID  string           `json:"character_id"`
	UserID       string           `json:"user_id"`
	StartTime    time.Time        `json:"start_time"`
	LastActivity time.Time        `json:"last_activity"`
	Messages     []SessionMessage `json:"messages"`
	Memories     []models.Memory  `json:"memories"`
	CacheMetrics CacheMetrics     `json:"cache_metrics"`
}

// SessionMessage represents a single message in a session
type SessionMessage struct {
	Timestamp    time.Time `json:"timestamp"`
	Role         string    `json:"role"` // "user" or "character"
	Content      string    `json:"content"`
	TokensUsed   int       `json:"tokens_used,omitempty"`
	CachedTokens int       `json:"cached_tokens,omitempty"` // Tokens served from OpenAI cache
	CacheHits    int       `json:"cache_hits,omitempty"`
	CacheMisses  int       `json:"cache_misses,omitempty"`
}

// CacheMetrics tracks cache performance for the session
type CacheMetrics struct {
	TotalRequests       int     `json:"total_requests"`
	CacheHits           int     `json:"cache_hits"`
	CacheMisses         int     `json:"cache_misses"`
	TokensSaved         int     `json:"tokens_saved"`
	CachedTokensTotal   int     `json:"cached_tokens_total"`   // Total OpenAI cached tokens
	PromptCacheHitRate  float64 `json:"prompt_cache_hit_rate"` // OpenAI prompt cache hit rate
	CostSaved           float64 `json:"cost_saved"`
	HitRate             float64 `json:"hit_rate"`
}

// SessionRepository manages session persistence
type SessionRepository struct {
	dataDir string
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(dataDir string) *SessionRepository {
	return &SessionRepository{dataDir: dataDir}
}

// SaveSession persists a session to disk
func (s *SessionRepository) SaveSession(session *Session) error {
	sessionDir := filepath.Join(s.dataDir, "sessions", session.CharacterID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	filename := filepath.Join(sessionDir, fmt.Sprintf("%s.json", session.ID))

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadSession loads a session from disk
func (s *SessionRepository) LoadSession(characterID, sessionID string) (*Session, error) {
	filename := filepath.Join(s.dataDir, "sessions", characterID, fmt.Sprintf("%s.json", sessionID))

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session %s not found", sessionID)
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// ListSessions returns all sessions for a character
func (s *SessionRepository) ListSessions(characterID string) ([]SessionInfo, error) {
	sessionDir := filepath.Join(s.dataDir, "sessions", characterID)

	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SessionInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}

	var sessions []SessionInfo
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			sessionID := filepath.Base(entry.Name())
			sessionID = sessionID[:len(sessionID)-5] // Remove .json

			session, err := s.LoadSession(characterID, sessionID)
			if err != nil {
				continue
			}

			sessions = append(sessions, SessionInfo{
				ID:           session.ID,
				CharacterID:  session.CharacterID,
				StartTime:    session.StartTime,
				LastActivity: session.LastActivity,
				MessageCount: len(session.Messages),
				CacheHitRate: session.CacheMetrics.HitRate,
			})
		}
	}

	return sessions, nil
}

// SessionInfo provides basic session information
type SessionInfo struct {
	ID           string    `json:"id"`
	CharacterID  string    `json:"character_id"`
	StartTime    time.Time `json:"start_time"`
	LastActivity time.Time `json:"last_activity"`
	MessageCount int       `json:"message_count"`
	CacheHitRate float64   `json:"cache_hit_rate"`
}

// GetLatestSession returns the most recent session for a character
func (s *SessionRepository) GetLatestSession(characterID string) (*Session, error) {
	sessions, err := s.ListSessions(characterID)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions found for character %s", characterID)
	}

	// Find most recent session
	var latest SessionInfo
	for _, session := range sessions {
		if session.LastActivity.After(latest.LastActivity) {
			latest = session
		}
	}

	return s.LoadSession(characterID, latest.ID)
}
