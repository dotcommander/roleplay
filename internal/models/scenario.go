package models

import "time"

// Scenario represents a high-level operational framework or meta-prompt
// that defines the overarching context for an entire class of interactions.
// This is the highest cache layer, sitting above even system prompts.
type Scenario struct {
	ID          string    `json:"id"`          // e.g., "starship_bridge_crisis_v1"
	Name        string    `json:"name"`        // User-friendly name
	Description string    `json:"description"` // What this scenario is for
	Prompt      string    `json:"prompt"`      // The actual meta-prompt content
	Version     int       `json:"version"`     // Version number for tracking changes
	Tags        []string  `json:"tags"`        // Tags for categorization
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastUsed    time.Time `json:"last_used"`
}

// ScenarioRequest represents a request that includes scenario context
type ScenarioRequest struct {
	ScenarioID string `json:"scenario_id,omitempty"`
}
