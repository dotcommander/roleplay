package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:    "demo",
	Short:  "Run a caching demonstration",
	Long: `Demonstrates the prompt caching system with a series of interactions
that showcase cache hits, misses, and cost savings.`,
	RunE:   runDemo,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(demoCmd)
	demoCmd.Flags().String("character", "rick-c137", "Character ID to use for demo")
	demoCmd.Flags().Bool("create-character", true, "Create demo character if it doesn't exist")
}

func runDemo(cmd *cobra.Command, args []string) error {
	characterID, _ := cmd.Flags().GetString("character")
	createChar, _ := cmd.Flags().GetBool("create-character")

	// Initialize configuration
	cfg := GetConfig()

	// Create manager (bot is now fully initialized)
	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		return err
	}

	// Create or load demo character
	if createChar {
		if err := createDemoCharacter(mgr, characterID); err != nil {
			return err
		}
	}

	// Ensure character is loaded
	char, err := mgr.GetOrLoadCharacter(characterID)
	if err != nil {
		return fmt.Errorf("failed to load character: %w", err)
	}

	// Create demo session
	sessionID := fmt.Sprintf("demo-%d", time.Now().Unix())
	session := &repository.Session{
		ID:           sessionID,
		CharacterID:  characterID,
		UserID:       "demo-user",
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Messages:     []repository.SessionMessage{},
		CacheMetrics: repository.CacheMetrics{},
	}

	// Initialize styles
	styles := newDemoStyles()

	// Display demo header
	displayDemoHeader(styles, char)

	// Get demo messages
	demoMessages := getDemoMessages()

	// Run demo interactions
	ctx := context.Background()
	for i, demo := range demoMessages {
		if demo.delay > 0 {
			time.Sleep(demo.delay)
		}

		// Display interaction header
		fmt.Printf("\n%s[Message %d] %s\n",
			styles.separator.Render(""),
			i+1,
			demo.description,
		)
		fmt.Printf("%sUser: %s\n",
			styles.bold.Render(""),
			styles.message.Render(demo.message),
		)

		// Process request
		req := models.ConversationRequest{
			CharacterID: characterID,
			UserID:      "demo-user",
			Message:     demo.message,
		}

		resp, _, err := processDemoMessage(ctx, mgr.GetBot(), &req, char, styles)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Update session metrics
		updateSessionMetrics(session, demo.message, resp)
	}

	// Calculate and save final metrics
	if session.CacheMetrics.TotalRequests > 0 {
		session.CacheMetrics.HitRate = float64(session.CacheMetrics.CacheHits) /
			float64(session.CacheMetrics.TotalRequests)
	}
	session.CacheMetrics.CostSaved = float64(session.CacheMetrics.TokensSaved) * 0.000003 // Approximate cost per token
	session.LastActivity = time.Now()

	// Save session
	if err := mgr.GetSessionRepository().SaveSession(session); err != nil {
		fmt.Printf("\nWarning: Failed to save session: %v\n", err)
	}

	// Display summary
	displayDemoSummary(session, styles)

	return nil
}

func createDemoCharacter(mgr *manager.CharacterManager, characterID string) error {
	// Check if already exists
	if _, err := mgr.GetOrLoadCharacter(characterID); err == nil {
		return nil // Already exists
	}

	// Create Rick Sanchez for demo
	if characterID == "rick-c137" {
		char := &models.Character{
			ID:        "rick-c137",
			Name:      "Rick Sanchez",
			Backstory: `The smartest man in the universe from dimension C-137. Cynical, alcoholic mad scientist who drags his grandson Morty on dangerous adventures across dimensions. Inventor of portal gun technology. Believes that nothing matters and science is the only truth. Has complex family relationships and deep-seated emotional issues masked by nihilism and substance abuse.`,
			Personality: models.PersonalityTraits{
				Openness:          1.0,
				Conscientiousness: 0.2,
				Extraversion:      0.7,
				Agreeableness:     0.1,
				Neuroticism:       0.9,
			},
			CurrentMood: models.EmotionalState{
				Joy:      0.2,
				Surprise: 0.1,
				Anger:    0.4,
				Fear:     0.1,
				Sadness:  0.3,
				Disgust:  0.5,
			},
			Quirks: []string{
				"Burps frequently mid-sentence (*burp*)",
				"Uses people's names as punctuation when talking",
				"Drinks from a flask constantly",
				"Makes pop culture references from multiple dimensions",
				"Dismisses emotions as 'chemical reactions'",
				"Uses scientific terminology casually",
			},
			SpeechStyle: "Cynical, sarcastic, frequently interrupted by burps. Uses complex scientific terms mixed with crude language. Often goes on nihilistic rants.",
		}
		return mgr.CreateCharacter(char)
	}

	// Default demo character
	char := &models.Character{
		ID:   characterID,
		Name: "Cache Demo Assistant",
		Backstory: `A helpful AI assistant designed to demonstrate prompt caching capabilities. 
I have a consistent personality and knowledge base that can be efficiently cached.`,
		Personality: models.PersonalityTraits{
			Openness:          0.9,
			Conscientiousness: 0.8,
			Extraversion:      0.7,
			Agreeableness:     0.9,
			Neuroticism:       0.2,
		},
		CurrentMood: models.EmotionalState{
			Joy:      0.7,
			Surprise: 0.3,
			Anger:    0.1,
			Fear:     0.1,
			Sadness:  0.1,
			Disgust:  0.1,
		},
		Quirks: []string{"helpful", "efficient", "knowledgeable"},
	}

	return mgr.CreateCharacter(char)
}

// demoStyles holds all the Lipgloss styles used in the demo command
type demoStyles struct {
	title     lipgloss.Style
	cacheHit  lipgloss.Style
	cacheMiss lipgloss.Style
	metrics   lipgloss.Style
	message   lipgloss.Style
	separator lipgloss.Style
	bold      lipgloss.Style
}

// newDemoStyles creates and initializes all demo styles
func newDemoStyles() *demoStyles {
	return &demoStyles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7c6f64")).
			Background(lipgloss.Color("#3c3836")).
			Padding(0, 1),

		cacheHit: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b8bb26")).
			Bold(true),

		cacheMiss: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fb4934")).
			Bold(true),

		metrics: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#83a598")),

		message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ebdbb2")),

		separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#665c54")),

		bold: lipgloss.NewStyle().Bold(true),
	}
}

// demoMessage represents a single demo interaction
type demoMessage struct {
	message     string
	description string
	delay       time.Duration
}

// getDemoMessages returns the predefined demo message sequence
func getDemoMessages() []demoMessage {
	return []demoMessage{
		{
			"Tell me about yourself",
			"Initial request - establishes cache layers",
			0,
		},
		{
			"Tell me about yourself",
			"Exact repeat - should hit response cache",
			1 * time.Second,
		},
		{
			"What are your core values?",
			"New question - cache miss",
			1 * time.Second,
		},
		{
			"What are your core values?",
			"Repeat question - should hit cache",
			1 * time.Second,
		},
		{
			"Tell me about yourself",
			"Third repeat - should hit cache with high savings",
			1 * time.Second,
		},
	}
}

// displayDemoHeader shows the demo title and character info
func displayDemoHeader(styles *demoStyles, char *models.Character) {
	fmt.Println(styles.title.Render("🚀 Roleplay Prompt Caching Demo"))
	fmt.Printf("\nCharacter: %s (%s)\n", char.Name, char.ID)
	fmt.Println(strings.Repeat("─", 60))
}

// processDemoMessage handles a single demo interaction
func processDemoMessage(
	ctx context.Context,
	bot *services.CharacterBot,
	req *models.ConversationRequest,
	char *models.Character,
	styles *demoStyles,
) (*providers.AIResponse, time.Duration, error) {
	start := time.Now()
	resp, err := bot.ProcessRequest(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, elapsed, err
	}

	// Display response
	fmt.Printf("%s%s:\n", styles.bold.Render(""), char.Name)
	fmt.Printf("%s\n", styles.message.Render(utils.WrapText(resp.Content, 80)))

	// Display cache metrics
	cacheStatus := "MISS"
	style := styles.cacheMiss
	if resp.CacheMetrics.Hit {
		cacheStatus = fmt.Sprintf("HIT (%d layers)", len(resp.CacheMetrics.Layers))
		style = styles.cacheHit
	}

	fmt.Printf("\n%s\n", styles.metrics.Render(fmt.Sprintf(
		"  ⚡ Response Time: %v | Cache: %s | Tokens: %d (saved: %d)",
		elapsed,
		style.Render(cacheStatus),
		resp.TokensUsed.Total,
		resp.CacheMetrics.SavedTokens,
	)))

	return resp, elapsed, nil
}

// updateSessionMetrics updates the session with response metrics
func updateSessionMetrics(
	session *repository.Session,
	userMessage string,
	resp *providers.AIResponse,
) {
	// Add messages to session
	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp: time.Now(),
		Role:      "user",
		Content:   userMessage,
	})

	cacheHits := 0
	cacheMisses := 0
	if resp.CacheMetrics.Hit {
		cacheHits = 1
	} else {
		cacheMisses = 1
	}

	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp:   time.Now(),
		Role:        "character",
		Content:     resp.Content,
		TokensUsed:  resp.TokensUsed.Total,
		CacheHits:   cacheHits,
		CacheMisses: cacheMisses,
	})

	// Update cumulative metrics
	session.CacheMetrics.TotalRequests++
	if resp.CacheMetrics.Hit {
		session.CacheMetrics.CacheHits++
	} else {
		session.CacheMetrics.CacheMisses++
	}
	session.CacheMetrics.TokensSaved += resp.CacheMetrics.SavedTokens
}

// displayDemoSummary shows the final summary of the demo
func displayDemoSummary(session *repository.Session, styles *demoStyles) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println(styles.title.Render("📊 Demo Summary"))
	fmt.Printf("\nTotal Interactions: %d\n", session.CacheMetrics.TotalRequests)
	fmt.Printf("Overall Cache Hit Rate: %.1f%%\n", session.CacheMetrics.HitRate*100)
	fmt.Printf("Total Tokens Saved: %d\n", session.CacheMetrics.TokensSaved)
	fmt.Printf("Estimated Cost Saved: $%.4f\n", session.CacheMetrics.CostSaved)
	fmt.Printf("\nSession saved as: %s\n", session.ID)
	fmt.Println("\nView detailed metrics with: roleplay session stats")
}
