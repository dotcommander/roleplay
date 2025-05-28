package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/dotcommander/roleplay/internal/factory"
	"github.com/dotcommander/roleplay/internal/providers"
	"github.com/dotcommander/roleplay/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiTestCmd = &cobra.Command{
	Use:   "api-test",
	Short: "Test API connection with a simple completion request",
	Long:  `Sends a simple test message to verify API connectivity and configuration.`,
	RunE:  runAPITest,
}

func init() {
	rootCmd.AddCommand(apiTestCmd)
	apiTestCmd.Flags().String("message", "Hello! Please respond with a brief greeting.", "Test message to send")
}

func runAPITest(cmd *cobra.Command, args []string) error {
	// Get configuration
	provider := viper.GetString("provider")
	apiKey := viper.GetString("api_key")
	model := viper.GetString("model")
	message, _ := cmd.Flags().GetString("message")

	// No need to check environment here - factory will handle it

	fmt.Printf("Testing %s API...\n", provider)
	// Use factory to get default model if needed
	if model == "" {
		model = factory.GetDefaultModel(provider)
	}
	fmt.Printf("Model: %s\n", model)
	fmt.Printf("Message: %s\n\n", message)

	// Create provider using factory
	p, err := factory.CreateProviderWithFallback(provider, apiKey, model)
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
	fmt.Printf("Tokens used: %d\n", resp.TokensUsed)
	fmt.Printf("\nResponse:\n%s\n", utils.WrapText(resp.Content, 80))

	return nil
}
