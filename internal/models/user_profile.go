package models

import "time"

// UserFact represents a piece of information known about the user.
type UserFact struct {
	Key         string    `json:"key"`          // e.g., "PreferredColor", "StatedGoal", "MentionedPetName"
	Value       string    `json:"value"`        // e.g., "Blue", "Learn Go programming", "Buddy"
	SourceTurn  int       `json:"source_turn"`  // Turn number in conversation where this was inferred/stated
	Confidence  float64   `json:"confidence"`   // LLM's confidence in this fact (0.0-1.0)
	LastUpdated time.Time `json:"last_updated"`
}

// UserProfile holds synthesized information about a user.
type UserProfile struct {
	UserID           string     `json:"user_id"`
	CharacterID      string     `json:"character_id"`      // Profile might be character-specific
	Facts            []UserFact `json:"facts"`
	OverallSummary   string     `json:"overall_summary"`   // A brief LLM-generated summary of the user
	InteractionStyle string     `json:"interaction_style"` // e.g., "formal", "inquisitive", "humorous"
	LastAnalyzed     time.Time  `json:"last_analyzed"`
	Version          int        `json:"version"`
}