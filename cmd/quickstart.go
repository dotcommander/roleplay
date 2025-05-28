package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quick start with zero configuration",
	Long: `Automatically detects local LLM services or uses environment variables
to start chatting immediately with Rick Sanchez`,
	RunE: runQuickstart,
}

func init() {
	rootCmd.AddCommand(quickstartCmd)
}

func runQuickstart(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Roleplay Quickstart")
	fmt.Println("=====================")

	// Try to auto-detect configuration
	config := GetConfig()

	// Check if we have a base URL configured
	if config.BaseURL == "" {
		// Try to detect local services
		fmt.Println("\n🔍 Detecting local LLM services...")

		localEndpoints := []struct {
			name    string
			baseURL string
			testURL string
			model   string
		}{
			{"Ollama", "http://localhost:11434/v1", "http://localhost:11434/api/tags", "llama3"},
			{"LM Studio", "http://localhost:1234/v1", "http://localhost:1234/v1/models", "local-model"},
			{"LocalAI", "http://localhost:8080/v1", "http://localhost:8080/v1/models", "gpt-4"},
		}

		client := &http.Client{Timeout: 2 * time.Second}

		for _, endpoint := range localEndpoints {
			resp, err := client.Get(endpoint.testURL)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				fmt.Printf("✅ Found %s running at %s\n", endpoint.name, endpoint.baseURL)
				config.BaseURL = endpoint.baseURL
				config.Model = endpoint.model
				config.APIKey = "not-required"
				break
			}
		}

		// If no local service found, check for API keys
		if config.BaseURL == "" {
			if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
				fmt.Println("✅ Found OPENAI_API_KEY in environment")
				config.BaseURL = "https://api.openai.com/v1"
				config.APIKey = apiKey
				config.Model = "gpt-4o-mini"
			} else if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
				fmt.Println("✅ Found ANTHROPIC_API_KEY in environment")
				config.DefaultProvider = "anthropic"
				config.APIKey = apiKey
				config.Model = "claude-3-haiku-20240307"
			} else {
				return fmt.Errorf("no LLM service detected. Please run 'roleplay init' to configure")
			}
		}
	} else {
		fmt.Printf("✅ Using configured endpoint: %s\n", config.BaseURL)
	}

	// Initialize manager
	fmt.Println("\n🎭 Initializing character system...")
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Create or load Rick
	characterID := "rick-c137"
	if _, err := mgr.GetOrLoadCharacter(characterID); err != nil {
		// Create Rick
		fmt.Println("🧬 Creating Rick Sanchez...")
		rick := createQuickstartRick()
		if err := mgr.CreateCharacter(rick); err != nil {
			return fmt.Errorf("failed to create character: %w", err)
		}
	}

	// Get user ID
	userID := os.Getenv("USER")
	if userID == "" {
		userID = os.Getenv("USERNAME") // Windows
	}
	if userID == "" {
		userID = "user"
	}

	// Create a quick session
	sessionID := fmt.Sprintf("quickstart-%d", time.Now().Unix())

	fmt.Printf("\n💬 Starting chat with Rick Sanchez\n")
	fmt.Printf("   User: %s\n", userID)
	fmt.Printf("   Model: %s\n", config.Model)
	fmt.Println("\n" + strings.Repeat("─", 60) + "\n")

	// Send a greeting
	req := &models.ConversationRequest{
		CharacterID: characterID,
		UserID:      userID,
		Message:     "Hello Rick!",
		Context: models.ConversationContext{
			SessionID:      sessionID,
			StartTime:      time.Now(),
			RecentMessages: []models.Message{},
		},
	}

	fmt.Printf("You: Hello Rick!\n\n")

	ctx := context.Background()
	resp, err := mgr.GetBot().ProcessRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	fmt.Printf("Rick: %s\n", resp.Content)

	// Show next steps
	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Println("\n✨ Quickstart successful!")
	fmt.Println("\nNext steps:")
	fmt.Println("  • Continue chatting: roleplay interactive")
	fmt.Println("  • Configure properly: roleplay init")
	fmt.Println("  • Create more characters: roleplay character create")
	fmt.Println("  • Check configuration: roleplay config list")

	if resp.CacheMetrics.Hit {
		fmt.Printf("\n💡 Tip: This response used cached prompts, saving %d tokens!\n", resp.CacheMetrics.SavedTokens)
	}

	return nil
}

func createQuickstartRick() *models.Character {
	return &models.Character{
		ID:        "rick-c137",
		Name:      "Rick Sanchez",
		Backstory: `The smartest man in the universe from dimension C-137. A genius scientist with a nihilistic worldview. Inventor of portal gun technology.`,
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
			"Burps mid-sentence (*burp*)",
			"Uses Morty's name as punctuation",
		},
		SpeechStyle: "Rapid-fire with burps. Mixes science with crude humor. Sarcastic.",
		Memories:    []models.Memory{},
	}
}
