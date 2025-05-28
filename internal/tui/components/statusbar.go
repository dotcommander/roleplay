package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusBar displays connection status and cache metrics
type StatusBar struct {
	width         int
	connected     bool
	cacheHits     int
	cacheMisses   int
	tokensSaved   int
	sessionID     string
	model         string
	lastError     error
	styles        statusBarStyles
}

type statusBarStyles struct {
	container    lipgloss.Style
	connected    lipgloss.Style
	disconnected lipgloss.Style
	metrics      lipgloss.Style
	error        lipgloss.Style
	session      lipgloss.Style
	model        lipgloss.Style
}

// NewStatusBar creates a new status bar component
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{
		width:     width,
		connected: true,
		sessionID: "session",
		model:     "gpt-4o-mini",
		styles: statusBarStyles{
			container:    lipgloss.NewStyle().Background(GruvboxAqua).Foreground(GruvboxBg),
			connected:    lipgloss.NewStyle().Foreground(GruvboxBg).Bold(true),
			disconnected: lipgloss.NewStyle().Foreground(GruvboxRed).Bold(true),
			metrics:      lipgloss.NewStyle().Foreground(GruvboxBg),
			error:        lipgloss.NewStyle().Foreground(GruvboxRed).Bold(true),
			session:      lipgloss.NewStyle().Foreground(GruvboxBg),
			model:        lipgloss.NewStyle().Foreground(GruvboxBg),
		},
	}
}

// Update handles messages for the status bar
func (s *StatusBar) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width

	case StatusUpdateMsg:
		s.connected = msg.Connected
		s.cacheHits = msg.CacheHits
		s.cacheMisses = msg.CacheMisses
		s.tokensSaved = msg.TokensSaved
		if msg.SessionID != "" {
			s.sessionID = msg.SessionID
		}
		if msg.Model != "" {
			s.model = msg.Model
		}
		s.lastError = msg.Error
	}

	return nil
}

// View renders the status bar
func (s *StatusBar) View() string {
	cacheRate := 0.0
	totalRequests := s.cacheHits + s.cacheMisses
	if totalRequests > 0 {
		cacheRate = float64(s.cacheHits) / float64(totalRequests) * 100
	}

	// Connection status
	statusStyle := s.styles.connected
	statusIcon := "●"
	if !s.connected {
		statusStyle = s.styles.disconnected
		statusIcon = "○"
	}

	// Session ID (truncated if too long)
	sessionDisplay := s.sessionID
	if len(sessionDisplay) > 8 {
		sessionDisplay = sessionDisplay[:8]
	}

	// Build status parts
	parts := []string{
		statusStyle.Render(fmt.Sprintf("%s %s", statusIcon, sessionDisplay)),
		s.styles.model.Render(s.model),
		s.styles.metrics.Render(fmt.Sprintf(" %d", totalRequests)),
		s.styles.metrics.Render(fmt.Sprintf("%s %.0f%%", statusIcon, cacheRate)),
		s.styles.metrics.Render(fmt.Sprintf(" %d tokens saved", s.tokensSaved)),
	}

	status := lipgloss.JoinHorizontal(lipgloss.Left, parts[0], " │ ", parts[1], " │", parts[2], " │ ", parts[3], " │", parts[4])

	// Add error if present
	if s.lastError != nil {
		errorText := s.styles.error.Render(fmt.Sprintf(" │ Error: %v", s.lastError))
		status = lipgloss.JoinHorizontal(lipgloss.Left, status, errorText)
	}

	return s.styles.container.Width(s.width).Render(" " + status)
}

// SetSize updates the width of the status bar
func (s *StatusBar) SetSize(width, _ int) {
	s.width = width
}