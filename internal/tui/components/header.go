package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Header displays character information and personality stats
type Header struct {
	title         string
	characterName string
	characterID   string
	personality   PersonalityStats
	mood          string
	moodIcon      string
	width         int
	styles        headerStyles
}

type headerStyles struct {
	title       lipgloss.Style
	personality lipgloss.Style
	mood        lipgloss.Style
	character   lipgloss.Style
}

// NewHeader creates a new header component
func NewHeader(width int) *Header {
	return &Header{
		title:         "Chat",
		characterName: "Loading...",
		width:         width,
		mood:          "Unknown",
		moodIcon:      "🤔",
		styles: headerStyles{
			title: lipgloss.NewStyle().
				Foreground(GruvboxAqua).
				Bold(true).
				Padding(0, 1),
			personality: lipgloss.NewStyle().
				Foreground(GruvboxPurple),
			mood: lipgloss.NewStyle().
				Foreground(GruvboxYellow),
			character: lipgloss.NewStyle().
				Foreground(GruvboxOrange).
				Bold(true),
		},
	}
}

// Update handles messages for the header
func (h *Header) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width

	case CharacterUpdateMsg:
		h.characterName = msg.Name
		h.characterID = msg.ID
		h.mood = msg.Mood
		h.moodIcon = msg.MoodIcon
		h.personality = msg.Personality
		h.title = fmt.Sprintf("󰊕 Chat with %s", msg.Name)
	}

	return nil
}

// View renders the header
func (h *Header) View() string {
	titleLine := h.styles.title.Render(h.title)

	// Personality traits with icons
	personalityStr := fmt.Sprintf(
		" O:%.1f  C:%.1f  E:%.1f  A:%.1f  N:%.1f",
		h.personality.Openness,
		h.personality.Conscientiousness,
		h.personality.Extraversion,
		h.personality.Agreeableness,
		h.personality.Neuroticism,
	)

	personalityInfo := h.styles.personality.Render(personalityStr)
	moodInfo := h.styles.mood.Render(fmt.Sprintf(" %s %s", h.moodIcon, h.mood))

	infoLine := fmt.Sprintf("  %s • %s", personalityInfo, moodInfo)

	return lipgloss.JoinVertical(lipgloss.Left, titleLine, infoLine, "")
}

// SetSize updates the width of the header
func (h *Header) SetSize(width, _ int) {
	h.width = width
}
