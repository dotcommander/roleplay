package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
)

// handleSlashCommand processes slash commands and returns a tea.Cmd
func (m *Model) handleSlashCommand(input string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return slashCommandResult{
				cmdType: "error",
				content: "Invalid command",
				msgType: "error",
			}
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "/exit", "/quit", "/q":
			return slashCommandResult{cmdType: "quit"}

		case "/help", "/h":
			return slashCommandResult{
				cmdType: "help",
				content: `Available slash commands:
/help, /h     - Show this help message
/exit, /quit, /q - Exit the chat
/clear, /c    - Clear chat history
/list         - List all available characters
/switch <id>  - Switch to a different character
/stats        - Show cache statistics
/mood         - Show character's current mood
/personality  - Show character's personality traits
/session      - Show session information`,
				msgType: "help",
			}

		case "/clear", "/c":
			return slashCommandResult{cmdType: "clear"}

		case "/list":
			return m.listCharacters()

		case "/stats":
			return m.showStats()

		case "/mood":
			return m.showMood()

		case "/personality":
			return m.showPersonality()

		case "/session":
			return m.showSession()

		case "/switch":
			if len(parts) < 2 {
				return slashCommandResult{
					cmdType: "error",
					content: "Usage: /switch <character-id>\nUse /list to see available characters",
					msgType: "error",
				}
			}
			return m.switchCharacter(parts[1])

		default:
			return slashCommandResult{
				cmdType: "error",
				content: fmt.Sprintf("Unknown command: %s\nType /help for available commands", command),
				msgType: "error",
			}
		}
	}
}

func (m *Model) listCharacters() slashCommandResult {
	// List all available characters from the repository
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	charRepo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return slashCommandResult{
			cmdType: "error",
			content: fmt.Sprintf("Error accessing characters: %v", err),
			msgType: "error",
		}
	}

	characterIDs, err := charRepo.ListCharacters()
	if err != nil {
		return slashCommandResult{
			cmdType: "error",
			content: fmt.Sprintf("Error listing characters: %v", err),
			msgType: "error",
		}
	}

	if len(characterIDs) == 0 {
		return slashCommandResult{
			cmdType: "info",
			content: "No characters available. Use 'roleplay character create' to add characters.",
			msgType: "info",
		}
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
			// Load into bot for future use
			_ = m.bot.CreateCharacter(char)
		}

		// Current character indicator
		indicator := "  "
		if id == m.characterID {
			indicator = "→ "
		}

		// Mood icon
		tempMood := m.calculateMood(char)
		moodIcon := m.getMoodIcon(tempMood)

		listText.WriteString(fmt.Sprintf("\n%s%s (%s) %s %s\n",
			indicator, char.Name, id, moodIcon, tempMood))

		// Add brief description from backstory (first sentence)
		backstory := char.Backstory
		if idx := strings.Index(backstory, "."); idx != -1 && idx < 100 {
			backstory = backstory[:idx+1]
		} else if len(backstory) > 100 {
			backstory = backstory[:97] + "..."
		}
		listText.WriteString(fmt.Sprintf("   %s\n", backstory))
	}

	return slashCommandResult{
		cmdType: "list",
		content: listText.String(),
		msgType: "list",
	}
}

func (m *Model) showStats() slashCommandResult {
	cacheRate := 0.0
	if m.totalRequests > 0 {
		cacheRate = float64(m.cacheHits) / float64(m.totalRequests) * 100
	}

	content := fmt.Sprintf(`Cache Statistics:
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

	return slashCommandResult{
		cmdType: "stats",
		content: content,
		msgType: "stats",
	}
}

func (m *Model) showMood() slashCommandResult {
	if m.character == nil {
		return slashCommandResult{
			cmdType: "error",
			content: "Character not loaded",
			msgType: "error",
		}
	}

	mood := m.getDominantMood()
	icon := m.getMoodIcon(mood)
	content := fmt.Sprintf(`%s Current Mood: %s

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

	return slashCommandResult{
		cmdType: "mood",
		content: content,
		msgType: "info",
	}
}

func (m *Model) showPersonality() slashCommandResult {
	if m.character == nil {
		return slashCommandResult{
			cmdType: "error",
			content: "Character not loaded",
			msgType: "error",
		}
	}

	content := fmt.Sprintf(`%s's Personality (OCEAN Model):

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

	return slashCommandResult{
		cmdType: "personality",
		content: content,
		msgType: "info",
	}
}

func (m *Model) showSession() slashCommandResult {
	characterName := m.characterID
	if m.character != nil {
		characterName = m.character.Name
	}
	sessionIDDisplay := m.sessionID
	if len(m.sessionID) > 8 {
		sessionIDDisplay = m.sessionID[:8] + "..."
	}

	messageCount := 0 // Would be tracked in the refactored version

	content := fmt.Sprintf(`Session Information:
• Session ID: %s
• Character: %s (%s)
• User: %s
• Messages: %d
• Started: %s`,
		sessionIDDisplay,
		characterName,
		m.characterID,
		m.userID,
		messageCount,
		m.context.StartTime.Format("Jan 2, 2006 15:04"))

	return slashCommandResult{
		cmdType: "session",
		content: content,
		msgType: "info",
	}
}

func (m *Model) switchCharacter(newCharID string) slashCommandResult {
	// Check if it's the same character
	if newCharID == m.characterID {
		return slashCommandResult{
			cmdType: "info",
			content: fmt.Sprintf("Already chatting with %s", m.characterID),
			msgType: "info",
		}
	}

	// Try to load the character
	char, err := m.bot.GetCharacter(newCharID)
	if err != nil {
		// If not loaded in bot, try loading from repository
		dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
		charRepo, repoErr := repository.NewCharacterRepository(dataDir)
		if repoErr != nil {
			return slashCommandResult{
				cmdType: "error",
				content: fmt.Sprintf("Error accessing characters: %v", repoErr),
				msgType: "error",
			}
		}

		char, err = charRepo.LoadCharacter(newCharID)
		if err != nil {
			return slashCommandResult{
				cmdType: "error",
				content: fmt.Sprintf("Character '%s' not found. Use /list to see available characters", newCharID),
				msgType: "error",
			}
		}

		// Load character into bot
		if err := m.bot.CreateCharacter(char); err != nil {
			return slashCommandResult{
				cmdType: "error",
				content: fmt.Sprintf("Error loading character: %v", err),
				msgType: "error",
			}
		}
	}

	return slashCommandResult{
		cmdType:        "switch",
		newCharacterID: newCharID,
		newCharacter:   char,
	}
}

// Helper function to calculate mood for a character
func (m *Model) calculateMood(char *models.Character) string {
	if char == nil {
		return "Unknown"
	}

	moods := map[string]float64{
		"Joy":      char.CurrentMood.Joy,
		"Surprise": char.CurrentMood.Surprise,
		"Anger":    char.CurrentMood.Anger,
		"Fear":     char.CurrentMood.Fear,
		"Sadness":  char.CurrentMood.Sadness,
		"Disgust":  char.CurrentMood.Disgust,
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
