package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/spf13/cobra"
)

var configTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test API connection with a simple completion request",
	Long:  `Sends a simple test message to verify API connectivity and configuration.`,
	RunE:  runAPITest,
}

func init() {
	configCmd.AddCommand(configTestCmd)
	configTestCmd.Flags().String("message", "Hello! Please respond with a brief greeting.", "Test message to send")
}

func runAPITest(cmd *cobra.Command, args []string) error {
	// Get configuration from global config (already resolved in initConfig)
	cfg := GetConfig()
	message, _ := cmd.Flags().GetString("message")

	fmt.Printf("Testing %s endpoint...\n", cfg.DefaultProvider)
	if cfg.BaseURL != "" {
		fmt.Printf("Base URL: %s\n", cfg.BaseURL)
	}
	// Use factory to get default model if needed
	if cfg.Model == "" {
		cfg.Model = factory.GetDefaultModel(cfg.DefaultProvider)
	}
	fmt.Printf("Model: %s\n", cfg.Model)
	fmt.Printf("Message: %s\n\n", message)

	// Create provider using factory
	p, err := factory.CreateProviderWithFallback(cfg.DefaultProvider, cfg.APIKey, cfg.Model, cfg.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// Create request
	req := &providers.PromptRequest{
		SystemPrompt: "You are a helpful assistant. Keep responses brief.",
		Message:      message,
	}

	// Send request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	resp, err := p.SendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	elapsed := time.Since(start)

	// Display results
	fmt.Println("✓ API test successful!")
	fmt.Printf("Response time: %v\n", elapsed)
	fmt.Printf("Tokens used: %d (prompt: %d, completion: %d)\n",
		resp.TokensUsed.Total, resp.TokensUsed.Prompt, resp.TokensUsed.Completion)
	if resp.TokensUsed.CachedPrompt > 0 {
		fmt.Printf("Cached tokens: %d\n", resp.TokensUsed.CachedPrompt)
	}
	fmt.Printf("\nResponse:\n%s\n", utils.WrapText(resp.Content, 80))

	return nil
}
