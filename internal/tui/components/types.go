package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Component defines the interface for all TUI components
type Component interface {
	Update(tea.Msg) tea.Cmd
	View() string
}

// WithFocus allows components to track focus state
type WithFocus interface {
	Focus()
	Blur()
	IsFocused() bool
}

// WithSize allows components to be resized
type WithSize interface {
	SetSize(width, height int)
}

// Message types for component communication
type (
	// FocusChangeMsg indicates focus has changed to a new component
	FocusChangeMsg struct {
		ComponentID string
	}

	// StatusUpdateMsg updates the status bar
	StatusUpdateMsg struct {
		Connected   bool
		CacheHits   int
		CacheMisses int
		TokensSaved int
		SessionID   string
		Model       string
		Error       error
	}

	// CharacterUpdateMsg updates character info in header
	CharacterUpdateMsg struct {
		Name        string
		ID          string
		Mood        string
		MoodIcon    string
		Personality PersonalityStats
	}

	// MessageAppendMsg adds a new message to the chat
	MessageAppendMsg struct {
		Role    string
		Content string
		MsgType string // "normal", "help", "list", "stats", etc.
	}

	// ProcessingStateMsg updates the processing state
	ProcessingStateMsg struct {
		IsProcessing bool
	}
)

// PersonalityStats holds OCEAN personality traits
type PersonalityStats struct {
	Openness          float64
	Conscientiousness float64
	Extraversion      float64
	Agreeableness     float64
	Neuroticism       float64
}

// Common styles using Gruvbox Dark theme
var (
	// Gruvbox Dark Colors
	GruvboxBg     = lipgloss.Color("#282828")
	GruvboxBg1    = lipgloss.Color("#3c3836")
	GruvboxFg     = lipgloss.Color("#ebdbb2")
	GruvboxRed    = lipgloss.Color("#fb4934")
	GruvboxGreen  = lipgloss.Color("#b8bb26")
	GruvboxYellow = lipgloss.Color("#fabd2f")
	GruvboxPurple = lipgloss.Color("#d3869b")
	GruvboxAqua   = lipgloss.Color("#8ec07c")
	GruvboxOrange = lipgloss.Color("#fe8019")
	GruvboxGray   = lipgloss.Color("#928374")
	GruvboxFg2    = lipgloss.Color("#d5c4a1")
)
