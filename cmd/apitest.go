package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

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

	// Check for API key from environment if not set
	if apiKey == "" && provider == "openai" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" && provider == "anthropic" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	if apiKey == "" {
		return fmt.Errorf("API key not found. Set --api-key flag or %s_API_KEY environment variable",
			map[string]string{"openai": "OPENAI", "anthropic": "ANTHROPIC"}[provider])
	}

	fmt.Printf("Testing %s API...\n", provider)
	fmt.Printf("Model: %s\n", func() string {
		if model != "" {
			return model
		}
		if provider == "openai" {
			return "gpt-4o-mini"
		}
		return "claude-3-haiku-20240307"
	}())
	fmt.Printf("Message: %s\n\n", message)

	// Create provider
	var p providers.AIProvider
	switch provider {
	case "openai":
		if model == "" {
			model = "gpt-4o-mini"
		}
		p = providers.NewOpenAIProvider(apiKey, model)
	case "anthropic":
		p = providers.NewAnthropicProvider(apiKey)
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
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
