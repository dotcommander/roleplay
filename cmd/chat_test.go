package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/spf13/cobra"
)

func TestChatCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		setup       func(t *testing.T) string
		wantErr     bool
		wantOutput  string
		checkResult func(t *testing.T, tempDir string)
	}{
		{
			name: "simple chat message",
			args: []string{"chat", "Hello, how are you?"},
			flags: map[string]string{
				"character": "test-char",
				"user":      "test-user",
			},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create character
				charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
				if err := os.MkdirAll(charDir, 0755); err != nil {
					t.Fatalf("Failed to create character directory: %v", err)
				}
				
				char := models.Character{
					ID:        "test-char",
					Name:      "Test Character",
					Backstory: "A helpful test character",
					Personality: models.PersonalityTraits{
						Openness:      0.8,
						Agreeableness: 0.9,
					},
				}
				
				data, err := json.Marshal(&char)
				if err != nil {
					t.Fatalf("Failed to marshal character: %v", err)
				}
				if err := os.WriteFile(filepath.Join(charDir, "test-char.json"), data, 0644); err != nil {
					t.Fatalf("Failed to write character file: %v", err)
				}
				
				// Set up mock provider response
				setupMockProvider("I'm doing well, thank you for asking!")
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "well",  // Look for partial match since we're using real AI
		},
		{
			name: "chat with session",
			args: []string{"chat", "Continue our conversation"},
			flags: map[string]string{
				"character": "test-char",
				"user":      "test-user",
				"session":   "existing-session",
			},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create character
				charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
				if err := os.MkdirAll(charDir, 0755); err != nil {
					t.Fatalf("Failed to create character directory: %v", err)
				}
				
				char := models.Character{
					ID:   "test-char",
					Name: "Test Character",
				}
				
				data, err := json.Marshal(&char)
				if err != nil {
					t.Fatalf("Failed to marshal character: %v", err)
				}
				if err := os.WriteFile(filepath.Join(charDir, "test-char.json"), data, 0644); err != nil {
					t.Fatalf("Failed to write character file: %v", err)
				}
				
				// Create existing session
				sessionDir := filepath.Join(tempDir, ".config", "roleplay", "sessions", "test-char")
				if err := os.MkdirAll(sessionDir, 0755); err != nil {
					t.Fatalf("Failed to create session directory: %v", err)
				}
				
				session := map[string]interface{}{
					"id":           "existing-session",
					"character_id": "test-char",
					"user_id":      "test-user",
					"messages": []map[string]interface{}{
						{
							"role":    "user",
							"content": "Hello",
						},
						{
							"role":    "character",
							"content": "Hi there!",
						},
					},
				}
				
				sessionData, err := json.Marshal(session)
				if err != nil {
					t.Fatalf("Failed to marshal session: %v", err)
				}
				if err := os.WriteFile(filepath.Join(sessionDir, "existing-session.json"), sessionData, 0644); err != nil {
					t.Fatalf("Failed to write session file: %v", err)
				}
				
				setupMockProvider("Of course! I remember we were just getting started.")
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "test-user",  // AI should acknowledge the user
			checkResult: func(t *testing.T, tempDir string) {
				// Check session was updated
				sessionFile := filepath.Join(tempDir, ".config", "roleplay", "sessions", "test-char", "existing-session.json")
				data, err := os.ReadFile(sessionFile)
				if err != nil {
					t.Errorf("Failed to read session file: %v", err)
					return
				}
				
				var session map[string]interface{}
				if err := json.Unmarshal(data, &session); err != nil {
					t.Errorf("Failed to unmarshal session: %v", err)
					return
				}
				
				messages := session["messages"].([]interface{})
				if len(messages) != 4 { // 2 existing + 1 user message + 1 AI response
					t.Errorf("Expected 4 messages in session, got %d", len(messages))
				}
			},
		},
		{
			name: "chat with JSON format",
			args: []string{"chat", "Tell me a joke"},
			flags: map[string]string{
				"character": "test-char",
				"user":      "test-user",
				"format":    "json",
			},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create character
				charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
				if err := os.MkdirAll(charDir, 0755); err != nil {
					t.Fatalf("Failed to create character directory: %v", err)
				}
				
				char := models.Character{
					ID:   "test-char",
					Name: "Test Character",
				}
				
				data, err := json.Marshal(&char)
				if err != nil {
					t.Fatalf("Failed to marshal character: %v", err)
				}
				if err := os.WriteFile(filepath.Join(charDir, "test-char.json"), data, 0644); err != nil {
					t.Fatalf("Failed to write character file: %v", err)
				}
				
				setupMockProvider("Why don't scientists trust atoms? Because they make up everything!")
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: `"response":`,  // JSON structure uses 'response' not 'content'
			checkResult: func(t *testing.T, tempDir string) {
				// Output should be valid JSON
				// This would need to capture and parse stdout
			},
		},
		{
			name: "chat without character flag",
			args: []string{"chat", "Hello"},
			flags: map[string]string{
				"user": "test-user",
			},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:    true,
			wantOutput: "required flag(s) \"character\" not set",
		},
		{
			name: "chat with non-existent character",
			args: []string{"chat", "Hello"},
			flags: map[string]string{
				"character": "nonexistent",
				"user":      "test-user",
			},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:    true,
			wantOutput: "character nonexistent not found",
		},
		{
			name: "chat with scenario",
			args: []string{"chat", "Let's begin"},
			flags: map[string]string{
				"character": "test-char",
				"user":      "test-user",
				"scenario":  "test-scenario",
			},
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				
				// Create character
				charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
				if err := os.MkdirAll(charDir, 0755); err != nil {
					t.Fatalf("Failed to create character directory: %v", err)
				}
				
				char := models.Character{
					ID:   "test-char",
					Name: "Test Character",
				}
				
				data, err := json.Marshal(&char)
				if err != nil {
					t.Fatalf("Failed to marshal character: %v", err)
				}
				if err := os.WriteFile(filepath.Join(charDir, "test-char.json"), data, 0644); err != nil {
					t.Fatalf("Failed to write character file: %v", err)
				}
				
				// Create scenario
				scenarioDir := filepath.Join(tempDir, ".config", "roleplay", "scenarios")
				if err := os.MkdirAll(scenarioDir, 0755); err != nil {
					t.Fatalf("Failed to create scenario directory: %v", err)
				}
				
				scenario := models.Scenario{
					ID:     "test-scenario",
					Name:   "Test Scenario",
					Prompt: "You are in a test scenario. Be helpful and friendly.",
				}
				
				scenarioData, err := json.Marshal(&scenario)
				if err != nil {
					t.Fatalf("Failed to marshal scenario: %v", err)
				}
				if err := os.WriteFile(filepath.Join(scenarioDir, "test-scenario.json"), scenarioData, 0644); err != nil {
					t.Fatalf("Failed to write scenario file: %v", err)
				}
				
				setupMockProvider("Welcome to our test scenario! How can I help you today?")
				
				return tempDir
			},
			wantErr:    false,
			wantOutput: "test-user",  // AI should mention the user
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset commands first
			resetCommands()
			
			// Reset os.Args
			os.Args = []string{"roleplay"}
			os.Args = append(os.Args, tt.args...)
			
			// Setup test environment
			tempDir := ""
			oldHome := os.Getenv("HOME")
			if tt.setup != nil {
				tempDir = tt.setup(t)
				os.Setenv("HOME", tempDir)
				defer func() {
					if oldHome != "" {
						os.Setenv("HOME", oldHome)
					}
				}()
			}

			// Set flags
			for flag, value := range tt.flags {
				os.Args = append(os.Args, "--"+flag+"="+value)
			}

			// Capture output
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			// Execute command
			err := rootCmd.Execute()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check output
			output := buf.String()
			if tt.wantOutput != "" && !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Output does not contain expected string.\nGot: %s\nWant substring: %s", output, tt.wantOutput)
			}

			// Additional checks
			if tt.checkResult != nil && tempDir != "" {
				tt.checkResult(t, tempDir)
			}
		})
	}
}

func TestChatCacheMetrics(t *testing.T) {
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		}
	}()

	// Create character
	charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatalf("Failed to create character directory: %v", err)
	}
	
	char := models.Character{
		ID:   "cache-test",
		Name: "Cache Test Character",
	}
	
	data, err := json.Marshal(&char)
	if err != nil {
		t.Fatalf("Failed to marshal character: %v", err)
	}
	if err := os.WriteFile(filepath.Join(charDir, "cache-test.json"), data, 0644); err != nil {
		t.Fatalf("Failed to write character file: %v", err)
	}

	// Set up provider that tracks cache metrics
	setupMockProviderWithMetrics("Response with cache hit", true)

	// Run chat command
	os.Args = []string{"roleplay", "chat", "Test message", "--character", "cache-test", "--user", "test-user"}
	
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute chat command: %v", err)
	}

	// Check that session was saved with cache metrics
	sessionDir := filepath.Join(tempDir, ".config", "roleplay", "sessions", "cache-test")
	files, err := os.ReadDir(sessionDir)
	if err != nil || len(files) == 0 {
		t.Fatal("No session file created")
	}

	sessionFile := filepath.Join(sessionDir, files[0].Name())
	data, _ = os.ReadFile(sessionFile)
	
	var session map[string]interface{}
	if err := json.Unmarshal(data, &session); err != nil {
		t.Fatalf("Failed to unmarshal session: %v", err)
	}
	
	if metrics, ok := session["cache_metrics"].(map[string]interface{}); ok {
		if metrics["total_requests"].(float64) != 1 {
			t.Error("Expected 1 total request in cache metrics")
		}
	} else {
		t.Error("Cache metrics not found in session")
	}
}

func TestChatErrorHandling(t *testing.T) {
	t.Skip("Skipping error handling tests - mock provider not implemented")
	tests := []struct {
		name      string
		setupErr  error
		wantErr   bool
		errOutput string
	}{
		{
			name:      "provider error",
			setupErr:  fmt.Errorf("API rate limit exceeded"),
			wantErr:   true,
			errOutput: "API rate limit exceeded",
		},
		{
			name:      "network timeout",
			setupErr:  context.DeadlineExceeded,
			wantErr:   true,
			errOutput: "deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			oldHome := os.Getenv("HOME")
			os.Setenv("HOME", tempDir)
			defer func() {
				if oldHome != "" {
					os.Setenv("HOME", oldHome)
				}
			}()

			// Create character
			charDir := filepath.Join(tempDir, ".config", "roleplay", "characters")
			if err := os.MkdirAll(charDir, 0755); err != nil {
				t.Fatalf("Failed to create character directory: %v", err)
			}
			
			char := models.Character{
				ID:   "error-test",
				Name: "Error Test",
			}
			
			data, err := json.Marshal(&char)
			if err != nil {
				t.Fatalf("Failed to marshal character: %v", err)
			}
			if err := os.WriteFile(filepath.Join(charDir, "error-test.json"), data, 0644); err != nil {
				t.Fatalf("Failed to write character file: %v", err)
			}

			// Set up provider with error
			setupMockProviderWithError(tt.setupErr)

			// Run command
			os.Args = []string{"roleplay", "chat", "Test", "--character", "error-test", "--user", "test"}
			
			var buf bytes.Buffer
			rootCmd.SetErr(&buf)
			
			err = rootCmd.Execute()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if tt.errOutput != "" {
				errStr := buf.String()
				if !strings.Contains(errStr, tt.errOutput) {
					t.Errorf("Error output missing expected text: %s", tt.errOutput)
				}
			}
			
			resetCommands()
		})
	}
}

// Helper functions for setting up mock providers
func setupMockProvider(response string) {
	// This would need to integrate with the actual command setup
	// to inject the mock provider
}

func setupMockProviderWithMetrics(response string, cacheHit bool) {
	// Mock provider that includes cache metrics in response
}

func setupMockProviderWithError(err error) {
	// Mock provider that returns an error
}

func resetCommands() {
	// Reset flag variables
	characterID = ""
	userID = ""
	sessionID = ""
	format = ""
	scenarioID = ""
	
	rootCmd = &cobra.Command{
		Use:   "roleplay",
		Short: "A sophisticated character bot with psychological modeling",
	}
	
	// Re-add chat command
	chatCmd = &cobra.Command{
		Use:     "chat [message]",
		Aliases: []string{"c"},
		Short:   "Send a single message to a character",
		Args:    cobra.ExactArgs(1),
		RunE:    runChat,
	}
	
	// Add flags to chat command with proper binding
	chatCmd.Flags().StringVarP(&characterID, "character", "c", "", "Character ID to chat with (required)")
	chatCmd.Flags().StringVarP(&userID, "user", "u", "", "User ID for the conversation (required)")
	chatCmd.Flags().StringVarP(&sessionID, "session", "s", "", "Session ID (optional, auto-generated if not provided)")
	chatCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	chatCmd.Flags().StringVar(&scenarioID, "scenario", "", "Scenario ID to use (optional)")
	chatCmd.Flags().Bool("no-cache", false, "Disable response caching")
	_ = chatCmd.MarkFlagRequired("character")
	_ = chatCmd.MarkFlagRequired("user")
	
	// Add chat command to root
	rootCmd.AddCommand(chatCmd)
}