package models

import "time"

// Message represents a single message in a conversation
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationContext holds the current conversation state
type ConversationContext struct {
	RecentMessages []Message
	SessionID      string
	StartTime      time.Time
}

// ConversationRequest represents a user request to the character bot
type ConversationRequest struct {
	CharacterID string
	UserID      string
	Message     string
	Context     ConversationContext
	ScenarioID  string // Optional scenario context
}
