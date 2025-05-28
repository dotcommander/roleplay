package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dotcommander/roleplay/internal/utils"
)

// Message represents a chat message
type Message struct {
	Role    string
	Content string
	Time    time.Time
	MsgType string // "normal", "help", "list", "stats", etc.
}

// MessageList displays the chat history
type MessageList struct {
	messages []Message
	viewport viewport.Model
	width    int
	height   int
	styles   messageListStyles
}

type messageListStyles struct {
	userMessage      lipgloss.Style
	characterMessage lipgloss.Style
	systemMessage    lipgloss.Style
	timestamp        lipgloss.Style
	separator        lipgloss.Style
	emptyMessage     lipgloss.Style
	user             lipgloss.Style
	character        lipgloss.Style
	system           lipgloss.Style
	// Special command output styles
	commandHeader lipgloss.Style
	commandBox    lipgloss.Style
	helpCommand   lipgloss.Style
	helpDesc      lipgloss.Style
	listItemActive lipgloss.Style
	listItem       lipgloss.Style
}

// NewMessageList creates a new message list component
func NewMessageList(width, height int) *MessageList {
	vp := viewport.New(width, height)
	vp.SetContent("")

	return &MessageList{
		messages: []Message{},
		viewport: vp,
		width:    width,
		height:   height,
		styles: messageListStyles{
			userMessage: lipgloss.NewStyle().
				Foreground(GruvboxFg).
				Background(GruvboxBg1).
				Padding(0, 1).
				MarginRight(2),
			characterMessage: lipgloss.NewStyle().
				Foreground(GruvboxFg).
				Background(GruvboxBg).
				Padding(0, 1).
				MarginLeft(2),
			systemMessage: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			timestamp: lipgloss.NewStyle().
				Foreground(GruvboxGray).
				Italic(true),
			separator: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			emptyMessage: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			user: lipgloss.NewStyle().
				Foreground(GruvboxGreen).
				Bold(true),
			character: lipgloss.NewStyle().
				Foreground(GruvboxOrange).
				Bold(true),
			system: lipgloss.NewStyle().
				Foreground(GruvboxGray),
			// Command styles
			commandHeader: lipgloss.NewStyle().
				Foreground(GruvboxAqua).
				Bold(true).
				Padding(0, 1),
			commandBox: lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(GruvboxAqua).
				Foreground(GruvboxFg).
				Background(GruvboxBg).
				Padding(1, 2),
			helpCommand: lipgloss.NewStyle().
				Foreground(GruvboxYellow).
				Bold(true),
			helpDesc: lipgloss.NewStyle().
				Foreground(GruvboxFg2),
			listItemActive: lipgloss.NewStyle().
				Foreground(GruvboxGreen).
				Bold(true),
			listItem: lipgloss.NewStyle().
				Foreground(GruvboxFg2),
		},
	}
}

// Update handles messages for the message list
func (m *MessageList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

	case MessageAppendMsg:
		m.messages = append(m.messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
			Time:    time.Now(),
			MsgType: msg.MsgType,
		})
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

// View renders the message list
func (m *MessageList) View() string {
	return m.viewport.View()
}

// SetSize updates the size of the message list
func (m *MessageList) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
	m.viewport.SetContent(m.renderMessages())
}

// renderMessages renders all messages with proper formatting
func (m *MessageList) renderMessages() string {
	if len(m.messages) == 0 {
		emptyMsg := m.styles.emptyMessage.Render("\n   Start chatting! Your conversation will appear here...\n")
		return emptyMsg
	}

	var content strings.Builder
	maxWidth := m.viewport.Width - 8 // Account for padding and margins

	for i, msg := range m.messages {
		if i > 0 {
			// Add visual separator between messages
			separator := m.styles.separator.Render(strings.Repeat("─", maxWidth))
			content.WriteString("\n" + separator + "\n\n")
		}

		timestamp := m.styles.timestamp.Render(msg.Time.Format("15:04:05"))

		if msg.Role == "user" {
			// User message
			header := fmt.Sprintf("┌─ %s %s", m.styles.user.Render("You"), timestamp)
			content.WriteString(m.styles.userMessage.Render(header) + "\n")

			wrappedContent := utils.WrapText(msg.Content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(m.styles.userMessage.Render(prefix+line) + "\n")
			}
		} else if msg.Role == "system" {
			// System message - check for special types
			if msg.MsgType == "help" || msg.MsgType == "list" || msg.MsgType == "info" || msg.MsgType == "stats" {
				// Special formatted output
				content.WriteString("\n")

				// Determine the header based on type
				var header string
				switch msg.MsgType {
				case "help":
					header = m.styles.commandHeader.Render("📚 Command Help")
				case "list":
					header = m.styles.commandHeader.Render("📋 Available Characters")
				case "stats":
					header = m.styles.commandHeader.Render("📊 Cache Statistics")
				case "info":
					header = m.styles.commandHeader.Render("ℹ️  Information")
				default:
					header = m.styles.commandHeader.Render("System")
				}

				// Format the content with special styling
				formattedContent := m.formatSpecialMessage(msg.Content, msg.MsgType, maxWidth-8)
				boxContent := lipgloss.JoinVertical(lipgloss.Left, header, "", formattedContent)
				content.WriteString(m.styles.commandBox.Width(maxWidth).Render(boxContent) + "\n")
			} else {
				// Regular system message
				header := fmt.Sprintf("┌─ %s %s", m.styles.system.Render("System"), timestamp)
				content.WriteString(m.styles.systemMessage.Render(header) + "\n")

				wrappedContent := utils.WrapText(msg.Content, maxWidth-4)
				lines := strings.Split(wrappedContent, "\n")
				for j, line := range lines {
					prefix := "│ "
					if j == len(lines)-1 {
						prefix = "└ "
					}
					content.WriteString(m.styles.systemMessage.Render(prefix+line) + "\n")
				}
			}
		} else {
			// Character message
			header := fmt.Sprintf("┌─ %s %s", m.styles.character.Render(msg.Role), timestamp)
			content.WriteString(m.styles.characterMessage.Render(header) + "\n")

			wrappedContent := utils.WrapText(msg.Content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(m.styles.characterMessage.Render(prefix+line) + "\n")
			}
		}
	}

	return content.String()
}

// formatSpecialMessage formats messages based on their type
func (m *MessageList) formatSpecialMessage(content string, msgType string, width int) string {
	switch msgType {
	case "help":
		// Format help message with colored commands
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, " - ") {
				parts := strings.SplitN(line, " - ", 2)
				if len(parts) == 2 {
					cmd := m.styles.helpCommand.Render(parts[0])
					desc := m.styles.helpDesc.Render("- " + parts[1])
					formatted = append(formatted, fmt.Sprintf("%s %s", cmd, desc))
				} else {
					formatted = append(formatted, line)
				}
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	case "list":
		// Format character list with special styling
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.HasPrefix(line, "→ ") {
				// Active character
				formatted = append(formatted, m.styles.listItemActive.Render(line))
			} else if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") {
				// Character name line
				formatted = append(formatted, m.styles.character.Render(line))
			} else if strings.HasPrefix(line, "   ") {
				// Description line
				formatted = append(formatted, m.styles.listItem.Render(line))
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	case "stats":
		// Format stats with colored numbers
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					label := parts[0]
					value := strings.TrimSpace(parts[1])
					// Color numbers
					if strings.Contains(value, "%") || strings.Contains(value, "tokens") {
						value = m.styles.character.Render(value)
					}
					formatted = append(formatted, fmt.Sprintf("%s: %s", label, value))
				} else {
					formatted = append(formatted, line)
				}
			} else {
				formatted = append(formatted, line)
			}
		}
		return strings.Join(formatted, "\n")

	default:
		// Default formatting
		return utils.WrapText(content, width)
	}
}

// ClearMessages clears all messages
func (m *MessageList) ClearMessages() {
	m.messages = []Message{}
	m.viewport.SetContent("")
}