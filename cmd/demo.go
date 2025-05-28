package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dotcommander/roleplay/internal/config"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/dotcommander/roleplay/internal/services"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a caching demonstration",
	Long: `Demonstrates the prompt caching system with a series of interactions
that showcase cache hits, misses, and cost savings.`,
	RunE: runDemo,
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

	// Create manager
	mgr, err := manager.NewCharacterManager(cfg)
	if err != nil {
		return err
	}

	// Setup provider
	setupDemoProvider(mgr.GetBot(), cfg)

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

	// Style definitions
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7c6f64")).
		Background(lipgloss.Color("#3c3836")).
		Padding(0, 1)

	cacheHitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#b8bb26")).
		Bold(true)

	cacheMissStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fb4934")).
		Bold(true)

	metricsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#83a598"))

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ebdbb2"))

	// Demo messages
	demoMessages := []struct {
		message     string
		description string
		delay       time.Duration
	}{
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

	fmt.Println(titleStyle.Render("🚀 Roleplay Prompt Caching Demo"))
	fmt.Printf("\nCharacter: %s (%s)\n", char.Name, char.ID)
	fmt.Println(strings.Repeat("─", 60))

	// Run demo interactions
	for i, demo := range demoMessages {
		if demo.delay > 0 {
			time.Sleep(demo.delay)
		}

		fmt.Printf("\n%s[Message %d] %s\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#665c54")).Render(""),
			i+1,
			demo.description,
		)
		fmt.Printf("%sUser: %s\n",
			lipgloss.NewStyle().Bold(true).Render(""),
			messageStyle.Render(demo.message),
		)

		// Process request
		req := models.ConversationRequest{
			CharacterID: characterID,
			UserID:      "demo-user",
			Message:     demo.message,
		}

		ctx := context.Background()
		start := time.Now()
		resp, err := mgr.GetBot().ProcessRequest(ctx, &req)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Display response with word wrapping
		fmt.Printf("%s%s:\n",
			lipgloss.NewStyle().Bold(true).Render(""),
			char.Name,
		)
		fmt.Printf("%s\n", messageStyle.Render(utils.WrapText(resp.Content, 80)))

		// Display cache metrics
		cacheStatus := "MISS"
		style := cacheMissStyle
		if resp.CacheMetrics.Hit {
			cacheStatus = fmt.Sprintf("HIT (%d layers)", len(resp.CacheMetrics.Layers))
			style = cacheHitStyle
		}

		fmt.Printf("\n%s\n", metricsStyle.Render(fmt.Sprintf(
			"  ⚡ Response Time: %v | Cache: %s | Tokens: %d (saved: %d)",
			elapsed,
			style.Render(cacheStatus),
			resp.TokensUsed.Total,
			resp.CacheMetrics.SavedTokens,
		)))

		// Update session
		session.Messages = append(session.Messages, repository.SessionMessage{
			Timestamp: time.Now(),
			Role:      "user",
			Content:   demo.message,
		})
		session.Messages = append(session.Messages, repository.SessionMessage{
			Timestamp:  time.Now(),
			Role:       "character",
			Content:    resp.Content,
			TokensUsed: resp.TokensUsed.Total,
			CacheHits: func() int {
				if resp.CacheMetrics.Hit {
					return 1
				} else {
					return 0
				}
			}(),
			CacheMisses: func() int {
				if !resp.CacheMetrics.Hit {
					return 1
				} else {
					return 0
				}
			}(),
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

	// Calculate final metrics
	session.CacheMetrics.HitRate = float64(session.CacheMetrics.CacheHits) /
		float64(session.CacheMetrics.CacheHits+session.CacheMetrics.CacheMisses)
	session.CacheMetrics.CostSaved = float64(session.CacheMetrics.TokensSaved) * 0.000003 // Approximate cost per token
	session.LastActivity = time.Now()

	// Save session
	if err := mgr.GetSessionRepository().SaveSession(session); err != nil {
		fmt.Printf("\nWarning: Failed to save session: %v\n", err)
	}

	// Display summary
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println(titleStyle.Render("📊 Demo Summary"))
	fmt.Printf("\nTotal Interactions: %d\n", len(demoMessages))
	fmt.Printf("Overall Cache Hit Rate: %.1f%%\n", session.CacheMetrics.HitRate*100)
	fmt.Printf("Total Tokens Saved: %d\n", session.CacheMetrics.TokensSaved)
	fmt.Printf("Estimated Cost Saved: $%.4f\n", session.CacheMetrics.CostSaved)
	fmt.Printf("\nSession saved as: %s\n", sessionID)
	fmt.Println("\nView detailed metrics with: roleplay session stats")

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

func setupDemoProvider(bot *services.CharacterBot, cfg *config.Config) {
	provider := cfg.DefaultProvider
	apiKey := cfg.APIKey

	if apiKey == "" && provider == "openai" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" && provider == "anthropic" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	switch provider {
	case "anthropic":
		if apiKey != "" {
			p := providers.NewAnthropicProvider(apiKey)
			bot.RegisterProvider("anthropic", p)
		}
	case "openai":
		if apiKey != "" {
			model := cfg.Model
			if model == "" {
				model = "gpt-4o-mini"
			}
			p := providers.NewOpenAIProvider(apiKey, model)
			bot.RegisterProvider("openai", p)
		}
	}
}
