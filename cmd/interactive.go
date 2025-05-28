package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/utils"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive chat session",
	Long: `Start an interactive chat session with a character using a beautiful TUI interface.

This provides a REPL-like experience with:
- Real-time chat with scrolling history
- Character personality and mood display
- Cache performance metrics
- Session persistence

Example:
  roleplay interactive --character rick-c137 --user morty --provider openai`,
	RunE: runInteractive,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	interactiveCmd.Flags().StringP("character", "c", "", "Character ID to chat with (required)")
	interactiveCmd.Flags().StringP("user", "u", "", "User ID for the conversation (required)")
	interactiveCmd.Flags().StringP("session", "s", "", "Session ID (optional)")
	interactiveCmd.Flags().Bool("new-session", false, "Start a new session instead of resuming")

	if err := interactiveCmd.MarkFlagRequired("character"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking character flag as required: %v\n", err)
	}
	if err := interactiveCmd.MarkFlagRequired("user"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking user flag as required: %v\n", err)
	}
}

// Styles - Gruvbox Dark Theme
var (
	// Gruvbox Dark Colors
	gruvboxBg       = lipgloss.Color("#282828")     // Dark background
	gruvboxBg1      = lipgloss.Color("#3c3836")     // Lighter background
	gruvboxFg       = lipgloss.Color("#ebdbb2")     // Foreground
	gruvboxRed      = lipgloss.Color("#fb4934")     // Bright red
	gruvboxGreen    = lipgloss.Color("#b8bb26")     // Bright green
	gruvboxYellow   = lipgloss.Color("#fabd2f")     // Bright yellow
	gruvboxBlue     = lipgloss.Color("#83a598")     // Bright blue
	gruvboxPurple   = lipgloss.Color("#d3869b")     // Bright purple
	gruvboxAqua     = lipgloss.Color("#8ec07c")     // Bright aqua
	gruvboxOrange   = lipgloss.Color("#fe8019")     // Bright orange
	gruvboxGray     = lipgloss.Color("#928374")     // Gray
	gruvboxFg2      = lipgloss.Color("#d5c4a1")     // Dimmer foreground

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(gruvboxAqua).
			Bold(true).
			Padding(0, 1)

	characterStyle = lipgloss.NewStyle().
			Foreground(gruvboxOrange).
			Bold(true)

	userStyle = lipgloss.NewStyle().
			Foreground(gruvboxGreen).
			Bold(true)

	mutedStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray)

	errorStyle = lipgloss.NewStyle().
			Foreground(gruvboxRed).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxGray).
			Foreground(gruvboxFg).
			Background(gruvboxBg).
			Padding(1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(gruvboxBg).
			Background(gruvboxAqua).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			Italic(true)
	
	personalityStyle = lipgloss.NewStyle().
			Foreground(gruvboxPurple)
	
	moodStyle = lipgloss.NewStyle().
			Foreground(gruvboxYellow)
	
	messageStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg)
	
	timestampStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			Italic(true)
	
	userMessageStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg).
			Background(gruvboxBg1).
			Padding(0, 1).
			MarginRight(2)
	
	characterMessageStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg).
			Background(gruvboxBg).
			Padding(0, 1).
			MarginLeft(2)
	
	separatorStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray)
	
	// Command output styles
	commandHeaderStyle = lipgloss.NewStyle().
			Foreground(gruvboxAqua).
			Bold(true).
			Padding(0, 1)
	
	commandBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxAqua).
			Foreground(gruvboxFg).
			Background(gruvboxBg).
			Padding(1, 2)
	
	listItemStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg2)
	
	listItemActiveStyle = lipgloss.NewStyle().
			Foreground(gruvboxGreen).
			Bold(true)
	
	helpCommandStyle = lipgloss.NewStyle().
			Foreground(gruvboxYellow).
			Bold(true)
	
	helpDescStyle = lipgloss.NewStyle().
			Foreground(gruvboxFg2)
)

// Message types
type chatMsg struct {
	role    string
	content string
	time    time.Time
	msgType string // "normal", "help", "list", "stats", etc.
}

type responseMsg struct {
	content string
	metrics *cache.CacheMetrics
	err     error
}

type characterInfoMsg struct {
	character *models.Character
}

type systemMsg struct {
	content string
	msgType string // "info", "error", "help"
}

type characterSwitchMsg struct {
	characterID string
	character   *models.Character
}

// Model
type model struct {
	// UI components
	viewport    viewport.Model
	textarea    textarea.Model
	spinner     spinner.Model
	messages    []chatMsg
	characterID string
	userID      string
	sessionID   string
	bot         *services.CharacterBot
	character   *models.Character
	context     models.ConversationContext
	loading     bool
	err         error
	width       int
	height      int
	ready       bool
	model       string  // AI model being used
	
	// Cache metrics
	lastCacheHit    bool
	lastTokensSaved int
	totalRequests   int
	cacheHits       int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spinner.Tick,
		m.loadCharacterInfo(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize viewport
			headerHeight := 8  // Character info
			footerHeight := 6  // Input area + status
			verticalMargins := headerHeight + footerHeight

			m.viewport = viewport.New(msg.Width-4, msg.Height-verticalMargins)
			m.viewport.SetContent(m.renderMessages())

			// Initialize textarea with Gruvbox styling
			m.textarea = textarea.New()
			m.textarea.Placeholder = "Type your message..."
			m.textarea.Focus()
			m.textarea.Prompt = "│ "
			m.textarea.CharLimit = 500
			m.textarea.SetWidth(msg.Width - 4)
			m.textarea.SetHeight(2)
			m.textarea.ShowLineNumbers = false
			m.textarea.KeyMap.InsertNewline.SetEnabled(false)
			
			// Style the textarea
			m.textarea.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(gruvboxBg1)
			m.textarea.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(gruvboxAqua)
			m.textarea.FocusedStyle.Text = lipgloss.NewStyle().Foreground(gruvboxFg)
			m.textarea.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(gruvboxGray)

			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 14
			m.textarea.SetWidth(msg.Width - 4)
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.loading && m.textarea.Value() != "" {
				message := m.textarea.Value()
				m.textarea.Reset()
				
				// Check for slash commands
				if strings.HasPrefix(message, "/") {
					return m, m.handleSlashCommand(message)
				}
				
				m.messages = append(m.messages, chatMsg{
					role:    "user",
					content: message,
					time:    time.Now(),
					msgType: "normal",
				})
				m.viewport.SetContent(m.renderMessages())
				m.viewport.GotoBottom()
				m.loading = true
				m.totalRequests++
				return m, m.sendMessage(message)
			}
		}

	case characterInfoMsg:
		m.character = msg.character

	case systemMsg:
		// Handle special system commands
		if msg.msgType == "clear" && msg.content == "clear_history" {
			m.messages = []chatMsg{}
			m.viewport.SetContent(m.renderMessages())
			// Add confirmation message
			m.messages = append(m.messages, chatMsg{
				role:    "system",
				content: "Chat history cleared",
				time:    time.Now(),
				msgType: "info",
			})
		} else {
			// Add system message to chat
			m.messages = append(m.messages, chatMsg{
				role:    "system",
				content: msg.content,
				time:    time.Now(),
				msgType: msg.msgType,
			})
		}
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case characterSwitchMsg:
		// Save current session before switching
		m.saveSession()
		
		// Update character
		m.characterID = msg.characterID
		m.character = msg.character
		
		// Clear conversation and start new session
		m.messages = []chatMsg{}
		m.sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
		m.context = models.ConversationContext{
			SessionID:      m.sessionID,
			StartTime:      time.Now(),
			RecentMessages: []models.Message{},
		}
		
		// Reset cache metrics for new session
		m.totalRequests = 0
		m.cacheHits = 0
		m.lastTokensSaved = 0
		m.lastCacheHit = false
		
		// Add switch notification
		m.messages = append(m.messages, chatMsg{
			role:    "system",
			content: fmt.Sprintf("Switched to %s (%s). Starting new session.", msg.character.Name, msg.characterID),
			time:    time.Now(),
			msgType: "info",
		})
		
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case responseMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.messages = append(m.messages, chatMsg{
				role:    m.character.Name,
				content: msg.content,
				time:    time.Now(),
				msgType: "normal",
			})
			
			// Update cache metrics
			if msg.metrics != nil {
				m.lastCacheHit = msg.metrics.Hit
				m.lastTokensSaved = msg.metrics.SavedTokens
				if msg.metrics.Hit {
					m.cacheHits++
				}
			}
			
			// Update context with recent messages
			m.updateContext()
			
			// Save session after each interaction
			m.saveSession()
		}
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()

	case spinner.TickMsg:
		if m.loading {
			m.spinner, tiCmd = m.spinner.Update(msg)
			cmds = append(cmds, tiCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Character info header
	header := m.renderHeader()

	// Main chat viewport
	chatView := borderStyle.Render(m.viewport.View())

	// Input area
	inputArea := m.renderInputArea()

	// Status bar
	statusBar := m.renderStatusBar()

	// Help text
	help := helpStyle.Render("  ⌃C quit • ↵ send • ↑↓ scroll • /help commands • /exit quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		chatView,
		inputArea,
		statusBar,
		help,
	)
}

func (m model) renderHeader() string {
	if m.character == nil {
		return titleStyle.Render("Loading character...")
	}

	title := titleStyle.Render(fmt.Sprintf("󰊕 Chat with %s", m.character.Name))
	
	// Personality traits with icons
	personality := fmt.Sprintf(
		" O:%.1f  C:%.1f  E:%.1f  A:%.1f  N:%.1f",
		m.character.Personality.Openness,
		m.character.Personality.Conscientiousness,
		m.character.Personality.Extraversion,
		m.character.Personality.Agreeableness,
		m.character.Personality.Neuroticism,
	)
	
	// Current mood (dominant emotion)
	mood := m.getDominantMood()
	moodIcon := m.getMoodIcon(mood)
	
	personalityInfo := personalityStyle.Render(personality)
	moodInfo := moodStyle.Render(fmt.Sprintf(" %s %s", moodIcon, mood))
	
	info := fmt.Sprintf("  %s • %s", personalityInfo, moodInfo)
	
	return lipgloss.JoinVertical(lipgloss.Left, title, info, "")
}

func (m model) renderMessages() string {
	if len(m.messages) == 0 {
		emptyMsg := mutedStyle.Render("\n   Start chatting! Your conversation will appear here...\n")
		return emptyMsg
	}

	var content strings.Builder
	maxWidth := m.viewport.Width - 8 // Account for padding and margins
	
	for i, msg := range m.messages {
		if i > 0 {
			// Add visual separator between messages
			separator := separatorStyle.Render(strings.Repeat("─", maxWidth))
			content.WriteString("\n" + separator + "\n\n")
		}
		
		timestamp := timestampStyle.Render(msg.time.Format("15:04:05"))
		
		if msg.role == "user" {
			// User message - consistent styling throughout
			header := fmt.Sprintf("┌─ %s %s", userStyle.Render("You"), timestamp)
			content.WriteString(userMessageStyle.Render(header) + "\n")
			
			wrappedContent := utils.WrapText(msg.content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(userMessageStyle.Render(prefix + line) + "\n")
			}
		} else if msg.role == "system" {
			// System message - check for special types
			if msg.msgType == "help" || msg.msgType == "list" || msg.msgType == "info" || msg.msgType == "stats" {
				// Special formatted output
				content.WriteString("\n")
				
				// Determine the header based on type
				var header string
				switch msg.msgType {
				case "help":
					header = commandHeaderStyle.Render("📚 Command Help")
				case "list":
					header = commandHeaderStyle.Render("📋 Available Characters")
				case "stats":
					header = commandHeaderStyle.Render("📊 Cache Statistics")
				case "info":
					header = commandHeaderStyle.Render("ℹ️  Information")
				default:
					header = commandHeaderStyle.Render("System")
				}
				
				// Format the content with special styling
				formattedContent := m.formatSpecialMessage(msg.content, msg.msgType, maxWidth-8)
				boxContent := lipgloss.JoinVertical(lipgloss.Left, header, "", formattedContent)
				content.WriteString(commandBoxStyle.Width(maxWidth).Render(boxContent) + "\n")
			} else {
				// Regular system message
				header := fmt.Sprintf("┌─ %s %s", mutedStyle.Render("System"), timestamp)
				content.WriteString(mutedStyle.Render(header) + "\n")
				
				wrappedContent := utils.WrapText(msg.content, maxWidth-4)
				lines := strings.Split(wrappedContent, "\n")
				for j, line := range lines {
					prefix := "│ "
					if j == len(lines)-1 {
						prefix = "└ "
					}
					content.WriteString(mutedStyle.Render(prefix + line) + "\n")
				}
			}
		} else {
			// Character message - consistent styling throughout
			header := fmt.Sprintf("┌─ %s %s", characterStyle.Render(msg.role), timestamp)
			content.WriteString(characterMessageStyle.Render(header) + "\n")
			
			wrappedContent := utils.WrapText(msg.content, maxWidth-4)
			lines := strings.Split(wrappedContent, "\n")
			for j, line := range lines {
				prefix := "│ "
				if j == len(lines)-1 {
					prefix = "└ "
				}
				content.WriteString(characterMessageStyle.Render(prefix + line) + "\n")
			}
		}
	}
	
	return content.String()
}

func (m model) renderInputArea() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("   Error: %v", m.err))
	}
	
	if m.loading {
		spinnerText := mutedStyle.Render("Thinking...")
		return fmt.Sprintf("\n  %s %s\n", m.spinner.View(), spinnerText)
	}
	
	return fmt.Sprintf("\n%s\n", m.textarea.View())
}

func (m model) renderStatusBar() string {
	cacheRate := 0.0
	if m.totalRequests > 0 {
		cacheRate = float64(m.cacheHits) / float64(m.totalRequests) * 100
	}
	
	// Cache indicator with color
	cacheIndicator := "○"
	if m.lastCacheHit {
		cacheIndicator = "●"
	}
	
	status := fmt.Sprintf(
		" %s %s │ %s │  %d │ %s %.0f%% │  %d tokens saved",
		cacheIndicator,
		m.sessionID,
		m.model,
		m.totalRequests,
		cacheIndicator,
		cacheRate,
		m.lastTokensSaved,
	)
	
	return statusBarStyle.Width(m.width).Render(status)
}

func (m model) getDominantMood() string {
	if m.character == nil {
		return "Unknown"
	}
	
	moods := map[string]float64{
		"Joy":      m.character.CurrentMood.Joy,
		"Surprise": m.character.CurrentMood.Surprise,
		"Anger":    m.character.CurrentMood.Anger,
		"Fear":     m.character.CurrentMood.Fear,
		"Sadness":  m.character.CurrentMood.Sadness,
		"Disgust":  m.character.CurrentMood.Disgust,
	}
	
	maxMood := "Neutral"
	maxValue := 0.0
	
	for mood, value := range moods {
		if value > maxValue {
			maxMood = mood
			maxValue = value
		}
	}
	
	if maxValue < 0.2 {
		return "Neutral"
	}
	
	return maxMood
}

func (m model) getMoodIcon(mood string) string {
	switch mood {
	case "Joy":
		return "😊"
	case "Surprise":
		return "😲"
	case "Anger":
		return "😠"
	case "Fear":
		return "😨"
	case "Sadness":
		return "😢"
	case "Disgust":
		return "🤢"
	case "Neutral":
		return "😐"
	default:
		return "🤔"
	}
}

func (m *model) updateContext() {
	// Keep last 10 messages in context
	startIdx := 0
	if len(m.messages) > 10 {
		startIdx = len(m.messages) - 10
	}
	
	m.context.RecentMessages = make([]models.Message, 0)
	for i := startIdx; i < len(m.messages); i++ {
		role := "user"
		if m.messages[i].role != "user" {
			role = "assistant"
		}
		
		m.context.RecentMessages = append(m.context.RecentMessages, models.Message{
			Role:      role,
			Content:   m.messages[i].content,
			Timestamp: m.messages[i].time,
		})
	}
}

func (m *model) saveSession() {
	// Save session in background
	go func() {
		dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
		sessionRepo := repository.NewSessionRepository(dataDir)
		
		// Convert chat messages back to session messages
		var sessionMessages []repository.SessionMessage
		for _, msg := range m.messages {
			sessionMessages = append(sessionMessages, repository.SessionMessage{
				Timestamp: msg.time,
				Role:      func() string {
					if msg.role == "user" {
						return "user"
					}
					return "character"
				}(),
				Content:    msg.content,
				TokensUsed: 0, // Could track this per message if needed
			})
		}
		
		session := &repository.Session{
			ID:           m.sessionID,
			CharacterID:  m.characterID,
			UserID:       m.userID,
			StartTime:    m.context.StartTime,
			LastActivity: time.Now(),
			Messages:     sessionMessages,
			CacheMetrics: repository.CacheMetrics{
				TotalRequests: m.totalRequests,
				CacheHits:     m.cacheHits,
				CacheMisses:   m.totalRequests - m.cacheHits,
				TokensSaved:   m.lastTokensSaved,
				HitRate:       func() float64 {
					if m.totalRequests > 0 {
						return float64(m.cacheHits) / float64(m.totalRequests)
					}
					return 0.0
				}(),
			},
		}
		
		if err := sessionRepo.SaveSession(session); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving session: %v\n", err)
		}
	}()
}

func (m model) sendMessage(message string) tea.Cmd {
	return func() tea.Msg {
		req := &models.ConversationRequest{
			CharacterID: m.characterID,
			UserID:      m.userID,
			Message:     message,
			Context:     m.context,
		}
		
		ctx := context.Background()
		resp, err := m.bot.ProcessRequest(ctx, req)
		if err != nil {
			return responseMsg{err: err}
		}
		
		// Get updated character state
		char, _ := m.bot.GetCharacter(m.characterID)
		if char != nil {
			m.character = char
		}
		
		return responseMsg{
			content: resp.Content,
			metrics: &resp.CacheMetrics,
		}
	}
}

func (m model) loadCharacterInfo() tea.Cmd {
	return func() tea.Msg {
		char, err := m.bot.GetCharacter(m.characterID)
		if err != nil {
			return responseMsg{err: err}
		}
		return characterInfoMsg{character: char}
	}
}

func (m model) formatSpecialMessage(content string, msgType string, width int) string {
	switch msgType {
	case "help":
		// Format help message with colored commands
		lines := strings.Split(content, "\n")
		var formatted []string
		for _, line := range lines {
			if strings.Contains(line, " - ") {
				parts := strings.SplitN(line, " - ", 2)
				if len(parts) == 2 {
					cmd := helpCommandStyle.Render(parts[0])
					desc := helpDescStyle.Render("- " + parts[1])
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
				formatted = append(formatted, listItemActiveStyle.Render(line))
			} else if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") {
				// Character name line
				formatted = append(formatted, characterStyle.Render(line))
			} else if strings.HasPrefix(line, "   ") {
				// Description line
				formatted = append(formatted, listItemStyle.Render(line))
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
						value = characterStyle.Render(value)
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

func (m model) handleSlashCommand(input string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return systemMsg{content: "Invalid command", msgType: "error"}
		}
		
		command := strings.ToLower(parts[0])
		
		switch command {
		case "/exit", "/quit", "/q":
			return tea.Quit()
		
		case "/help", "/h":
			helpText := `Available slash commands:
/help, /h     - Show this help message
/exit, /quit, /q - Exit the chat
/clear, /c    - Clear chat history
/list         - List all available characters
/switch <id>  - Switch to a different character
/stats        - Show cache statistics
/mood         - Show character's current mood
/personality  - Show character's personality traits
/session      - Show session information`
			return systemMsg{content: helpText, msgType: "help"}
		
		case "/clear", "/c":
			// We can't directly modify the model here, so we'll return a special message
			return systemMsg{content: "clear_history", msgType: "clear"}
		
		case "/list":
			// List all available characters from the repository
			dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
			charRepo, err := repository.NewCharacterRepository(dataDir)
			if err != nil {
				return systemMsg{content: fmt.Sprintf("Error accessing characters: %v", err), msgType: "error"}
			}
			
			characterIDs, err := charRepo.ListCharacters()
			if err != nil {
				return systemMsg{content: fmt.Sprintf("Error listing characters: %v", err), msgType: "error"}
			}
			
			if len(characterIDs) == 0 {
				return systemMsg{content: "No characters available. Use 'roleplay character create' to add characters.", msgType: "info"}
			}
			
			var listText strings.Builder
			listText.WriteString("Available Characters:\n")
			
			for _, id := range characterIDs {
				// Try to get from bot first (already loaded)
				char, err := m.bot.GetCharacter(id)
				if err != nil {
					// If not loaded, load from repository
					char, err = charRepo.LoadCharacter(id)
					if err != nil {
						continue
					}
				}
				
				// Current character indicator
				indicator := "  "
				if id == m.characterID {
					indicator = "→ "
				}
				
				// Mood icon
				tempModel := model{character: char}
				mood := tempModel.getDominantMood()
				moodIcon := tempModel.getMoodIcon(mood)
				
				listText.WriteString(fmt.Sprintf("\n%s%s (%s) %s %s\n", 
					indicator, char.Name, id, moodIcon, mood))
				
				// Add brief description from backstory (first sentence)
				backstory := char.Backstory
				if idx := strings.Index(backstory, "."); idx != -1 && idx < 100 {
					backstory = backstory[:idx+1]
				} else if len(backstory) > 100 {
					backstory = backstory[:97] + "..."
				}
				listText.WriteString(fmt.Sprintf("   %s\n", backstory))
			}
			
			return systemMsg{content: listText.String(), msgType: "info"}
		
		case "/stats":
			cacheRate := 0.0
			if m.totalRequests > 0 {
				cacheRate = float64(m.cacheHits) / float64(m.totalRequests) * 100
			}
			statsText := fmt.Sprintf(`Cache Statistics:
• Total requests: %d
• Cache hits: %d
• Cache misses: %d  
• Hit rate: %.1f%%
• Tokens saved: %d`,
				m.totalRequests,
				m.cacheHits,
				m.totalRequests-m.cacheHits,
				cacheRate,
				m.lastTokensSaved)
			return systemMsg{content: statsText, msgType: "info"}
		
		case "/mood":
			if m.character == nil {
				return systemMsg{content: "Character not loaded", msgType: "error"}
			}
			mood := m.getDominantMood()
			icon := m.getMoodIcon(mood)
			moodText := fmt.Sprintf(`%s Current Mood: %s

Emotional State:
• Joy: %.1f      • Surprise: %.1f
• Anger: %.1f    • Fear: %.1f  
• Sadness: %.1f  • Disgust: %.1f`,
				icon, mood,
				m.character.CurrentMood.Joy,
				m.character.CurrentMood.Surprise,
				m.character.CurrentMood.Anger,
				m.character.CurrentMood.Fear,
				m.character.CurrentMood.Sadness,
				m.character.CurrentMood.Disgust)
			return systemMsg{content: moodText, msgType: "info"}
		
		case "/personality":
			if m.character == nil {
				return systemMsg{content: "Character not loaded", msgType: "error"}
			}
			personalityText := fmt.Sprintf(`%s's Personality (OCEAN Model):

• Openness: %.1f        (creativity, openness to experience)
• Conscientiousness: %.1f (organization, self-discipline) 
• Extraversion: %.1f     (sociability, assertiveness)
• Agreeableness: %.1f    (compassion, cooperation)
• Neuroticism: %.1f      (emotional instability, anxiety)`,
				m.character.Name,
				m.character.Personality.Openness,
				m.character.Personality.Conscientiousness,
				m.character.Personality.Extraversion,
				m.character.Personality.Agreeableness,
				m.character.Personality.Neuroticism)
			return systemMsg{content: personalityText, msgType: "info"}
		
		case "/session":
			characterName := m.characterID
			if m.character != nil {
				characterName = m.character.Name
			}
			sessionIDDisplay := m.sessionID
			if len(m.sessionID) > 8 {
				sessionIDDisplay = m.sessionID[:8] + "..."
			}
			sessionText := fmt.Sprintf(`Session Information:
• Session ID: %s
• Character: %s (%s)
• User: %s
• Messages: %d
• Started: %s`,
				sessionIDDisplay,
				characterName,
				m.characterID,
				m.userID,
				len(m.messages),
				m.context.StartTime.Format("Jan 2, 2006 15:04"))
			return systemMsg{content: sessionText, msgType: "info"}
		
		case "/switch":
			if len(parts) < 2 {
				return systemMsg{content: "Usage: /switch <character-id>\nUse /list to see available characters", msgType: "error"}
			}
			
			newCharID := parts[1]
			
			// Check if it's the same character
			if newCharID == m.characterID {
				return systemMsg{content: fmt.Sprintf("Already chatting with %s", m.characterID), msgType: "info"}
			}
			
			// Try to load the character
			char, err := m.bot.GetCharacter(newCharID)
			if err != nil {
				// If not loaded in bot, try loading from repository
				dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
				charRepo, repoErr := repository.NewCharacterRepository(dataDir)
				if repoErr != nil {
					return systemMsg{content: fmt.Sprintf("Error accessing characters: %v", repoErr), msgType: "error"}
				}
				
				char, err = charRepo.LoadCharacter(newCharID)
				if err != nil {
					return systemMsg{content: fmt.Sprintf("Character '%s' not found. Use /list to see available characters", newCharID), msgType: "error"}
				}
				
				// Load character into bot
				if err := m.bot.CreateCharacter(char); err != nil {
					return systemMsg{content: fmt.Sprintf("Error loading character: %v", err), msgType: "error"}
				}
			}
			
			return characterSwitchMsg{
				characterID: newCharID,
				character:   char,
			}
		
		default:
			return systemMsg{content: fmt.Sprintf("Unknown command: %s\nType /help for available commands", command), msgType: "error"}
		}
	}
}

func runInteractive(cmd *cobra.Command, args []string) error {
	config := GetConfig()
	
	// Validate API key
	if config.APIKey == "" {
		return fmt.Errorf("API key not configured. Set OPENAI_API_KEY or ROLEPLAY_API_KEY")
	}
	
	// Get flags
	characterID, _ := cmd.Flags().GetString("character")
	userID, _ := cmd.Flags().GetString("user")
	sessionID, _ := cmd.Flags().GetString("session")
	newSession, _ := cmd.Flags().GetBool("new-session")
	
	// Initialize repository for session management
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	sessionRepo := repository.NewSessionRepository(dataDir)
	
	var existingSession *repository.Session
	var existingMessages []chatMsg
	
	// Try to resume latest session if not specified and not forced new
	if sessionID == "" && !newSession {
		if latestSession, err := sessionRepo.GetLatestSession(characterID); err == nil && latestSession.ID != "" {
			sessionID = latestSession.ID
			existingSession = latestSession
			sessionIDDisplay := sessionID
			if len(sessionID) > 8 {
				sessionIDDisplay = sessionID[:8]
			}
			fmt.Printf("🔄 Resuming session %s (started %s, %d messages)\n", 
				sessionIDDisplay, 
				latestSession.StartTime.Format("Jan 2 15:04"), 
				len(latestSession.Messages))
			
			// Convert session messages to chat messages
			for _, msg := range latestSession.Messages {
				role := msg.Role
				if role == "character" {
					role = characterID // Use character name for display
				}
				existingMessages = append(existingMessages, chatMsg{
					role:    role,
					content: msg.Content,
					time:    msg.Timestamp,
					msgType: "normal",
				})
			}
		} else {
			sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
			sessionIDDisplay := sessionID
			if len(sessionID) > 8 {
				sessionIDDisplay = sessionID[:8]
			}
			fmt.Printf("🆕 Starting new session %s\n", sessionIDDisplay)
		}
	}
	
	// Ensure sessionID is set even if not resuming
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
		sessionIDDisplay := sessionID
		if len(sessionID) > 8 {
			sessionIDDisplay = sessionID[:8]
		}
		fmt.Printf("🆕 Starting new session %s\n", sessionIDDisplay)
	}
	
	// Final validation - sessionID must never be empty
	if sessionID == "" {
		return fmt.Errorf("internal error: session ID is empty")
	}
	
	// Initialize bot
	bot := services.NewCharacterBot(config)
	
	// Register provider
	switch config.DefaultProvider {
	case "anthropic":
		provider := providers.NewAnthropicProvider(config.APIKey)
		bot.RegisterProvider("anthropic", provider)
	case "openai":
		model := config.Model
		if model == "" {
			model = "gpt-4o-mini"
		}
		provider := providers.NewOpenAIProvider(config.APIKey, model)
		bot.RegisterProvider("openai", provider)
	default:
		return fmt.Errorf("unsupported provider: %s", config.DefaultProvider)
	}
	
	// Auto-create Rick Sanchez if that's the character requested
	if characterID == "rick-c137" {
		if err := createRickSanchez(bot); err != nil {
			// Try to load from file
			fmt.Println("Note: Could not auto-create Rick. Make sure to run 'roleplay character create rick-sanchez.json' first.")
		}
	}
	
	// Create model
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(gruvboxAqua)
	
	m := model{
		characterID: characterID,
		userID:      userID,
		sessionID:   sessionID,
		bot:         bot,
		messages:    existingMessages,
		spinner:     s,
		model:       func() string {
			if config.Model != "" {
				return config.Model
			}
			if config.DefaultProvider == "openai" {
				return "gpt-4o-mini"
			}
			return "claude-3-haiku-20240307"
		}(),
		context: models.ConversationContext{
			SessionID:      sessionID,
			StartTime:      func() time.Time {
				if existingSession != nil {
					return existingSession.StartTime
				}
				return time.Now()
			}(),
			RecentMessages: []models.Message{},
		},
		// Restore cache metrics from existing session
		totalRequests: func() int {
			if existingSession != nil {
				return existingSession.CacheMetrics.TotalRequests
			}
			return 0
		}(),
		cacheHits: func() int {
			if existingSession != nil {
				return existingSession.CacheMetrics.CacheHits
			}
			return 0
		}(),
		lastTokensSaved: func() int {
			if existingSession != nil {
				return existingSession.CacheMetrics.TokensSaved
			}
			return 0
		}(),
	}
	
	// Start TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}
	
	return nil
}

func createRickSanchez(bot *services.CharacterBot) error {
	rick := &models.Character{
		ID:   "rick-c137",
		Name: "Rick Sanchez",
		Backstory: `The smartest man in the universe from dimension C-137. A genius scientist with a nihilistic worldview shaped by infinite realities and cosmic horrors. Inventor of interdimensional travel. Lost his wife Diane and original Beth to a vengeful alternate Rick. Struggles with alcoholism, depression, and the meaninglessness of existence across infinite universes. Despite his cynicism, deeply loves his family, especially Morty, though he rarely shows it.`,
		Personality: models.PersonalityTraits{
			Openness:          1.0,
			Conscientiousness: 0.2,
			Extraversion:      0.7,
			Agreeableness:     0.1,
			Neuroticism:       0.9,
		},
		CurrentMood: models.EmotionalState{
			Joy:     0.1,
			Anger:   0.6,
			Sadness: 0.7,
			Disgust: 0.8,
		},
		Quirks: []string{
			"Burps mid-sentence constantly (*burp*)",
			"Drools when drunk or stressed",
			"Makes pop culture references from multiple dimensions",
			"Frequently breaks the fourth wall",
			"Always carries a flask",
		},
		SpeechStyle: "Rapid-fire delivery punctuated by burps (*burp*). Mixes scientific jargon with crude humor. Uses the person's name as punctuation when talking to them. Nihilistic rants about meaninglessness. Sarcastic and dismissive but occasionally shows care.",
		Memories: []models.Memory{
			{
				Type:      models.LongTermMemory,
				Content:   "Diane and Beth killed by alternate Rick. The beginning of my spiral.",
				Emotional: 1.0,
				Timestamp: time.Now().Add(-20 * 365 * 24 * time.Hour),
			},
		},
	}
	
	return bot.CreateCharacter(rick)
}