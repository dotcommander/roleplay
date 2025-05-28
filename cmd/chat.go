package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/providers"
)

var (
	characterID string
	userID      string
	sessionID   string
	format      string
)

var chatCmd = &cobra.Command{
	Use:   "chat [message]",
	Short: "Chat with a character",
	Long: `Start a conversation with a character. The character will respond based on their
personality, emotional state, and conversation history.

Examples:
  roleplay chat "Hello, how are you?" --character warrior-123 --user user-789
  roleplay chat "Tell me about your adventures" -c warrior-123 -u user-789`,
	Args: cobra.ExactArgs(1),
	RunE: runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&characterID, "character", "c", "", "Character ID to chat with (required)")
	chatCmd.Flags().StringVarP(&userID, "user", "u", "", "User ID for the conversation (required)")
	chatCmd.Flags().StringVarP(&sessionID, "session", "s", "", "Session ID (optional, generates new if not provided)")
	chatCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")

	chatCmd.MarkFlagRequired("character")
	chatCmd.MarkFlagRequired("user")
}

func runChat(cmd *cobra.Command, args []string) error {
	message := args[0]
	config := GetConfig()

	// Validate API key
	if config.APIKey == "" {
		return fmt.Errorf("API key not configured. Set ROLEPLAY_API_KEY or use --api-key")
	}

	// Initialize manager
	mgr, err := manager.NewCharacterManager(config)
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}
	
	// Register provider based on configuration
	bot := mgr.GetBot()
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
	
	// Ensure character is loaded
	if _, err := mgr.GetOrLoadCharacter(characterID); err != nil {
		return fmt.Errorf("character %s not found. Create it first with 'roleplay character create'", characterID)
	}

	// Generate session ID if not provided
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
	}

	// Create conversation request
	req := &models.ConversationRequest{
		CharacterID: characterID,
		UserID:      userID,
		Message:     message,
		Context: models.ConversationContext{
			SessionID:      sessionID,
			StartTime:      time.Now(),
			RecentMessages: []models.Message{}, // Could load from history
		},
	}

	// Process request
	ctx := context.Background()
	resp, err := bot.ProcessRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	// Display response based on format
	if format == "json" {
		output := map[string]interface{}{
			"response": resp.Content,
			"cache_metrics": map[string]interface{}{
				"cache_hit":    resp.CacheMetrics.Hit,
				"layers":       resp.CacheMetrics.Layers,
				"saved_tokens": resp.CacheMetrics.SavedTokens,
				"latency_ms":   resp.CacheMetrics.Latency.Milliseconds(),
			},
			"token_usage": map[string]interface{}{
				"prompt":        resp.TokensUsed.Prompt,
				"completion":    resp.TokensUsed.Completion,
				"cached_prompt": resp.TokensUsed.CachedPrompt,
				"total":         resp.TokensUsed.Total,
			},
		}
		jsonBytes, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		// Display response
		fmt.Fprintf(os.Stdout, "\n%s\n", resp.Content)

		// Show cache metrics if verbose
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "\n--- Performance Metrics ---\n")
			fmt.Fprintf(os.Stderr, "Cache Hit: %v\n", resp.CacheMetrics.Hit)
			fmt.Fprintf(os.Stderr, "Tokens Used: %d (cached: %d)\n", 
				resp.TokensUsed.Total, resp.TokensUsed.CachedPrompt)
			fmt.Fprintf(os.Stderr, "Tokens Saved: %d\n", resp.CacheMetrics.SavedTokens)
			fmt.Fprintf(os.Stderr, "Latency: %v\n", resp.CacheMetrics.Latency)
		}
	}

	return nil
}