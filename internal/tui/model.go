package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dotcommander/roleplay/internal/cache"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/tui/components"
)

// Model represents the main TUI application state
type Model struct {
	// Components
	header      *components.Header
	messageList *components.MessageList
	inputArea   *components.InputArea
	statusBar   *components.StatusBar

	// UI state
	width        int
	height       int
	ready        bool
	currentFocus string // "input", "messages"

	// Business logic
	characterID string
	userID      string
	sessionID   string
	scenarioID  string
	bot         *services.CharacterBot
	character   *models.Character
	context     models.ConversationContext

	// Command history
	commandHistory []string
	historyIndex   int
	historyBuffer  string

	// Cache metrics
	lastCacheHit    bool
	lastTokensSaved int
	totalRequests   int
	cacheHits       int

	// Model info
	aiModel string
}

// NewModel creates a new TUI model
func NewModel(cfg Config) *Model {
	// Calculate initial sizes
	headerHeight := 3
	statusHeight := 1
	helpHeight := 1
	inputHeight := 3
	messagesHeight := 20 // Default, will be adjusted

	if cfg.Height > 0 {
		messagesHeight = cfg.Height - headerHeight - statusHeight - helpHeight - inputHeight - 2
	}

	width := 80 // Default width
	if cfg.Width > 0 {
		width = cfg.Width
	}

	return &Model{
		header:      components.NewHeader(width),
		messageList: components.NewMessageList(width-4, messagesHeight),
		inputArea:   components.NewInputArea(width),
		statusBar:   components.NewStatusBar(width),

		characterID: cfg.CharacterID,
		userID:      cfg.UserID,
		sessionID:   cfg.SessionID,
		scenarioID:  cfg.ScenarioID,
		bot:         cfg.Bot,
		context:     cfg.Context,
		aiModel:     cfg.Model,

		currentFocus: "input",

		commandHistory: []string{},
		historyIndex:   0,

		totalRequests: cfg.InitialMetrics.TotalRequests,
		cacheHits:     cfg.InitialMetrics.CacheHits,
	}
}

// Config holds configuration for creating a new Model
type Config struct {
	CharacterID    string
	UserID         string
	SessionID      string
	ScenarioID     string
	Bot            *services.CharacterBot
	Context        models.ConversationContext
	Model          string
	Width          int
	Height         int
	InitialMetrics struct {
		TotalRequests int
		CacheHits     int
		TokensSaved   int
	}
}

// Init initializes the TUI model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.inputArea.Init(),
		m.loadCharacterInfo(),
	)
}

// Update handles all incoming messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.initializeLayout(msg.Width, msg.Height)
			m.ready = true
		} else {
			m.resizeComponents(msg.Width, msg.Height)
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.saveSession()
			return m, tea.Quit

		case tea.KeyTab:
			// Switch focus between components
			if m.currentFocus == "input" {
				m.currentFocus = "messages"
				m.inputArea.Blur()
			} else {
				m.currentFocus = "input"
				m.inputArea.Focus()
			}

		case tea.KeyEnter:
			if m.currentFocus == "input" && !m.inputArea.IsFocused() {
				// Processing state, ignore enter
				return m, nil
			}

			if m.inputArea.Value() != "" {
				message := m.inputArea.Value()

				// Add to command history
				m.commandHistory = append(m.commandHistory, message)
				m.historyIndex = len(m.commandHistory)
				m.historyBuffer = ""

				m.inputArea.Reset()

				// Check for slash commands
				if strings.HasPrefix(message, "/") {
					return m, m.handleSlashCommand(message)
				}

				// Regular message
				m.messageList.Update(components.MessageAppendMsg{
					Role:    "user",
					Content: message,
					MsgType: "normal",
				})

				m.inputArea.Update(components.ProcessingStateMsg{IsProcessing: true})
				m.totalRequests++

				return m, m.sendMessage(message)
			}

		case tea.KeyUp:
			// Navigate command history
			if len(m.commandHistory) > 0 && m.historyIndex > 0 {
				if m.historyIndex == len(m.commandHistory) {
					m.historyBuffer = m.inputArea.Value()
				}
				m.historyIndex--
				m.inputArea.SetValue(m.commandHistory[m.historyIndex])
				m.inputArea.CursorEnd()
			}

		case tea.KeyDown:
			// Navigate command history
			if len(m.commandHistory) > 0 && m.historyIndex < len(m.commandHistory) {
				m.historyIndex++

				if m.historyIndex == len(m.commandHistory) {
					m.inputArea.SetValue(m.historyBuffer)
				} else {
					m.inputArea.SetValue(m.commandHistory[m.historyIndex])
				}
				m.inputArea.CursorEnd()
			}
		}

	case characterInfoMsg:
		m.character = msg.character
		m.updateCharacterDisplay()

	case responseMsg:
		m.inputArea.Update(components.ProcessingStateMsg{IsProcessing: false})

		if msg.err != nil {
			m.statusBar.Update(components.StatusUpdateMsg{Error: msg.err})
		} else {
			// Add character response
			m.messageList.Update(components.MessageAppendMsg{
				Role:    m.character.Name,
				Content: msg.content,
				MsgType: "normal",
			})

			// Update cache metrics
			if msg.metrics != nil {
				m.lastCacheHit = msg.metrics.Hit
				m.lastTokensSaved = msg.metrics.SavedTokens
				if msg.metrics.Hit {
					m.cacheHits++
				}
			}

			// Update context
			m.updateContext()

			// Update status bar
			m.updateStatusBar()

			// Save session
			go m.saveSession()
		}

	case slashCommandResult:
		// Handle slash command results
		switch msg.cmdType {
		case "quit":
			m.saveSession()
			return m, tea.Quit

		case "clear":
			m.messageList.ClearMessages()
			m.messageList.Update(components.MessageAppendMsg{
				Role:    "system",
				Content: "Chat history cleared",
				MsgType: "info",
			})

		case "switch":
			if msg.err != nil {
				m.messageList.Update(components.MessageAppendMsg{
					Role:    "system",
					Content: msg.err.Error(),
					MsgType: "error",
				})
			} else {
				// Save current session before switching
				m.saveSession()

				// Update character
				m.characterID = msg.newCharacterID
				m.character = msg.newCharacter
				m.updateCharacterDisplay()

				// Clear conversation and start new session
				m.messageList.ClearMessages()
				m.sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
				m.context = models.ConversationContext{
					SessionID:      m.sessionID,
					StartTime:      time.Now(),
					RecentMessages: []models.Message{},
				}

				// Reset cache metrics
				m.totalRequests = 0
				m.cacheHits = 0
				m.lastTokensSaved = 0
				m.lastCacheHit = false

				// Add switch notification
				m.messageList.Update(components.MessageAppendMsg{
					Role:    "system",
					Content: fmt.Sprintf("Switched to %s (%s). Starting new session.", msg.newCharacter.Name, msg.newCharacterID),
					MsgType: "info",
				})

				m.updateStatusBar()
			}

		default:
			// Display command output
			m.messageList.Update(components.MessageAppendMsg{
				Role:    "system",
				Content: msg.content,
				MsgType: msg.msgType,
			})
		}
	}

	// Route updates to components
	cmds = append(cmds, m.routeToComponents(msg)...)

	return m, tea.Batch(cmds...)
}

// View renders the entire TUI
func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// Render all components
	header := m.header.View()

	// Main chat viewport with border
	chatView := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(components.GruvboxGray).
		Foreground(components.GruvboxFg).
		Background(components.GruvboxBg).
		Padding(1).
		Render(m.messageList.View())

	inputArea := m.inputArea.View()
	statusBar := m.statusBar.View()

	// Help text
	help := lipgloss.NewStyle().
		Foreground(components.GruvboxGray).
		Italic(true).
		Render("  ⌃C quit • ↵ send • ↑↓ history • /help commands • /exit quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		chatView,
		inputArea,
		statusBar,
		help,
	)
}

// Helper methods

func (m *Model) initializeLayout(width, height int) {
	headerHeight := 3
	statusHeight := 1
	helpHeight := 1
	inputHeight := 3
	borderHeight := 2
	verticalMargins := headerHeight + statusHeight + helpHeight + inputHeight + borderHeight

	messagesHeight := height - verticalMargins
	if messagesHeight < 5 {
		messagesHeight = 5
	}

	m.header.SetSize(width, headerHeight)
	m.messageList.SetSize(width-4, messagesHeight)
	m.inputArea.SetSize(width, inputHeight)
	m.statusBar.SetSize(width, statusHeight)

	m.inputArea.Focus()
}

func (m *Model) resizeComponents(width, height int) {
	headerHeight := 3
	statusHeight := 1
	helpHeight := 1
	inputHeight := 3
	borderHeight := 2
	verticalMargins := headerHeight + statusHeight + helpHeight + inputHeight + borderHeight

	messagesHeight := height - verticalMargins
	if messagesHeight < 5 {
		messagesHeight = 5
	}

	m.header.SetSize(width, headerHeight)
	m.messageList.SetSize(width-4, messagesHeight)
	m.inputArea.SetSize(width, inputHeight)
	m.statusBar.SetSize(width, statusHeight)
}

func (m *Model) routeToComponents(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd

	// Update all components
	if cmd := m.header.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.messageList.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.inputArea.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.statusBar.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m *Model) updateCharacterDisplay() {
	if m.character == nil {
		return
	}

	mood := m.getDominantMood()
	moodIcon := m.getMoodIcon(mood)

	m.header.Update(components.CharacterUpdateMsg{
		Name:     m.character.Name,
		ID:       m.character.ID,
		Mood:     mood,
		MoodIcon: moodIcon,
		Personality: components.PersonalityStats{
			Openness:          m.character.Personality.Openness,
			Conscientiousness: m.character.Personality.Conscientiousness,
			Extraversion:      m.character.Personality.Extraversion,
			Agreeableness:     m.character.Personality.Agreeableness,
			Neuroticism:       m.character.Personality.Neuroticism,
		},
	})
}

func (m *Model) updateStatusBar() {
	m.statusBar.Update(components.StatusUpdateMsg{
		Connected:   true,
		CacheHits:   m.cacheHits,
		CacheMisses: m.totalRequests - m.cacheHits,
		TokensSaved: m.lastTokensSaved,
		SessionID:   m.sessionID,
		Model:       m.aiModel,
	})
}

func (m *Model) updateContext() {
	// Implementation would update conversation context
	// This is simplified for the example
}

func (m *Model) getDominantMood() string {
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

func (m *Model) getMoodIcon(mood string) string {
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

// Message types for tea.Cmd results

type characterInfoMsg struct {
	character *models.Character
}

type responseMsg struct {
	content string
	metrics *cache.CacheMetrics
	err     error
}

type slashCommandResult struct {
	cmdType        string // "help", "list", "stats", etc.
	content        string
	msgType        string // for display formatting
	err            error
	newCharacterID string // for switch command
	newCharacter   *models.Character
}

// Tea commands

func (m *Model) loadCharacterInfo() tea.Cmd {
	return func() tea.Msg {
		char, err := m.bot.GetCharacter(m.characterID)
		if err != nil {
			return responseMsg{err: err}
		}
		return characterInfoMsg{character: char}
	}
}

func (m *Model) sendMessage(message string) tea.Cmd {
	return func() tea.Msg {
		req := &models.ConversationRequest{
			CharacterID: m.characterID,
			UserID:      m.userID,
			Message:     message,
			ScenarioID:  m.scenarioID,
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

// Placeholder for session saving
func (m *Model) saveSession() {
	// Implementation would save the current session
	// This is left as a placeholder for the refactored version
}
