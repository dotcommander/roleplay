package components

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStatusBar(t *testing.T) {
	tests := []struct {
		name     string
		updates  []tea.Msg
		contains []string
	}{
		{
			name:     "Initial state shows connected",
			updates:  []tea.Msg{},
			contains: []string{"●", "session", "0"},
		},
		{
			name: "Updates cache metrics",
			updates: []tea.Msg{
				StatusUpdateMsg{
					Connected:   true,
					CacheHits:   10,
					CacheMisses: 5,
					TokensSaved: 1500,
				},
			},
			contains: []string{"15", "67%", "1500 tokens saved"},
		},
		{
			name: "Shows error when present",
			updates: []tea.Msg{
				StatusUpdateMsg{
					Error: errors.New("API rate limit exceeded"),
				},
			},
			contains: []string{"Error: API rate limit"},
		},
		{
			name: "Updates session ID",
			updates: []tea.Msg{
				StatusUpdateMsg{
					SessionID: "session-1234567890",
				},
			},
			contains: []string{"session-"}, // Should truncate
		},
		{
			name: "Shows disconnected state",
			updates: []tea.Msg{
				StatusUpdateMsg{
					Connected: false,
				},
			},
			contains: []string{"○"}, // Empty circle for disconnected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create status bar
			statusBar := NewStatusBar(80)

			// Apply updates
			for _, msg := range tt.updates {
				statusBar.Update(msg)
			}

			// Get view
			view := statusBar.View()

			// Check contains
			for _, expected := range tt.contains {
				if !strings.Contains(view, expected) {
					t.Errorf("Expected view to contain %q, got: %s", expected, view)
				}
			}
		})
	}
}

func TestStatusBarResize(t *testing.T) {
	statusBar := NewStatusBar(80)

	// Update with new size
	statusBar.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Check that width was updated
	if statusBar.width != 120 {
		t.Errorf("Expected width to be 120, got %d", statusBar.width)
	}

	// View should render at new width
	view := statusBar.View()
	// The view should have styling that attempts to use the full width
	// We can't easily test the exact width due to terminal rendering,
	// but we can verify it doesn't panic and produces output
	if len(view) == 0 {
		t.Error("Expected non-empty view after resize")
	}
}

func TestStatusBarCacheRate(t *testing.T) {
	tests := []struct {
		name         string
		hits         int
		misses       int
		expectedRate string
	}{
		{
			name:         "Perfect cache rate",
			hits:         10,
			misses:       0,
			expectedRate: "100%",
		},
		{
			name:         "No cache hits",
			hits:         0,
			misses:       10,
			expectedRate: "0%",
		},
		{
			name:         "50% cache rate",
			hits:         5,
			misses:       5,
			expectedRate: "50%",
		},
		{
			name:         "No requests",
			hits:         0,
			misses:       0,
			expectedRate: "0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusBar := NewStatusBar(80)
			statusBar.Update(StatusUpdateMsg{
				CacheHits:   tt.hits,
				CacheMisses: tt.misses,
			})

			view := statusBar.View()
			if !strings.Contains(view, tt.expectedRate) {
				t.Errorf("Expected cache rate %s in view, got: %s", tt.expectedRate, view)
			}
		})
	}
}
