package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/manager"
	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var (
	characterID string
	userID      string
	sessionID   string
	format      string
	scenarioID  string
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
	chatCmd.Flags().StringVar(&scenarioID, "scenario", "", "Scenario ID to set the interaction context (optional)")

	if err := chatCmd.MarkFlagRequired("character"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking character flag as required: %v\n", err)
	}
	if err := chatCmd.MarkFlagRequired("user"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking user flag as required: %v\n", err)
	}
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

	// Register provider using factory
	bot := mgr.GetBot()
	if err := factory.InitializeAndRegisterProvider(bot, config); err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}

	// Ensure character is loaded
	if _, err := mgr.GetOrLoadCharacter(characterID); err != nil {
		return fmt.Errorf("character %s not found. Create it first with 'roleplay character create'", characterID)
	}

	// Get session repository
	sessionRepo := mgr.GetSessionRepository()

	// Generate session ID if not provided
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%d", time.Now().Unix())
	}

	// Load existing session or create new one
	var session *repository.Session
	existingSession, err := sessionRepo.LoadSession(characterID, sessionID)
	if err != nil {
		// Create new session if it doesn't exist
		session = &repository.Session{
			ID:           sessionID,
			CharacterID:  characterID,
			UserID:       userID,
			StartTime:    time.Now(),
			LastActivity: time.Now(),
			Messages:     []repository.SessionMessage{},
			CacheMetrics: repository.CacheMetrics{},
		}
	} else {
		session = existingSession
	}

	// Convert session messages to conversation context
	var recentMessages []models.Message
	// Get last 10 messages for context
	startIdx := 0
	if len(session.Messages) > 10 {
		startIdx = len(session.Messages) - 10
	}
	for i := startIdx; i < len(session.Messages); i++ {
		msg := session.Messages[i]
		// Map "character" role to "assistant" for API compatibility
		role := msg.Role
		if role == "character" {
			role = "assistant"
		}
		recentMessages = append(recentMessages, models.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Create conversation request
	req := &models.ConversationRequest{
		CharacterID: characterID,
		UserID:      userID,
		Message:     message,
		ScenarioID:  scenarioID,
		Context: models.ConversationContext{
			SessionID:      sessionID,
			StartTime:      session.StartTime,
			RecentMessages: recentMessages,
		},
	}

	// Process request
	ctx := context.Background()
	resp, err := bot.ProcessRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	// Update session with new messages
	session.Messages = append(session.Messages, repository.SessionMessage{
		Timestamp:  time.Now(),
		Role:       "user",
		Content:    message,
		TokensUsed: 0, // User messages don't consume tokens
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

	// Update cache metrics
	session.CacheMetrics.TotalRequests++
	if resp.CacheMetrics.Hit {
		session.CacheMetrics.CacheHits++
	} else {
		session.CacheMetrics.CacheMisses++
	}
	session.CacheMetrics.TokensSaved += resp.CacheMetrics.SavedTokens
	session.CacheMetrics.HitRate = float64(session.CacheMetrics.CacheHits) / float64(session.CacheMetrics.TotalRequests)
	session.CacheMetrics.CostSaved = float64(session.CacheMetrics.TokensSaved) * 0.000003 // Approximate cost per token
	session.LastActivity = time.Now()

	// Save session BEFORE checking for profile updates
	// This ensures the async goroutine has access to the latest messages
	if err := sessionRepo.SaveSession(session); err != nil {
		// Log error but don't fail the command
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save session: %v\n", err)
		}
	}
	
	// Check if we should update the user profile
	// For the chat command, we do this synchronously to ensure it completes
	if config.UserProfileConfig.Enabled && len(session.Messages) > 0 && (len(session.Messages) % config.UserProfileConfig.UpdateFrequency == 0) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintf(os.Stderr, "Updating user profile...\n")
		}
		
		// Get character from manager
		char, err := mgr.GetOrLoadCharacter(characterID)
		if err == nil {
			// Call the bot's profile update method directly
			bot.UpdateUserProfile(userID, char, sessionID)
		}
	}

	// Display response based on format
	if format == "json" {
		output := map[string]interface{}{
			"session_id": sessionID,
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
			fmt.Fprintf(os.Stderr, "Session ID: %s\n", sessionID)
			fmt.Fprintf(os.Stderr, "Cache Hit: %v\n", resp.CacheMetrics.Hit)
			fmt.Fprintf(os.Stderr, "Tokens Used: %d (cached: %d)\n",
				resp.TokensUsed.Total, resp.TokensUsed.CachedPrompt)
			fmt.Fprintf(os.Stderr, "Tokens Saved: %d\n", resp.CacheMetrics.SavedTokens)
			fmt.Fprintf(os.Stderr, "Latency: %v\n", resp.CacheMetrics.Latency)
			fmt.Fprintf(os.Stderr, "Session Messages: %d\n", len(session.Messages))
		}
	}

	return nil
}
