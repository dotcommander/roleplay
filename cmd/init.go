package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup wizard for roleplay",
	Long: `Interactive setup wizard that guides you through configuring roleplay
for your preferred LLM provider (OpenAI, Ollama, LM Studio, OpenRouter, etc.)`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// Provider presets for common configurations
type providerPreset struct {
	Name         string
	BaseURL      string
	RequiresKey  bool
	LocalModel   bool
	DefaultModel string
	Description  string
}

var providerPresets = []providerPreset{
	{
		Name:         "openai",
		BaseURL:      "https://api.openai.com/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "gpt-4o-mini",
		Description:  "OpenAI (Official)",
	},
	{
		Name:         "anthropic",
		BaseURL:      "https://api.anthropic.com/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "claude-3-haiku-20240307",
		Description:  "Anthropic Claude (OpenAI-Compatible)",
	},
	{
		Name:         "gemini",
		BaseURL:      "https://generativelanguage.googleapis.com/v1beta/openai",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "models/gemini-1.5-flash",
		Description:  "Google Gemini (OpenAI-Compatible)",
	},
	{
		Name:         "ollama",
		BaseURL:      "http://localhost:11434/v1",
		RequiresKey:  false,
		LocalModel:   true,
		DefaultModel: "llama3",
		Description:  "Ollama (Local LLMs)",
	},
	{
		Name:         "lmstudio",
		BaseURL:      "http://localhost:1234/v1",
		RequiresKey:  false,
		LocalModel:   true,
		DefaultModel: "local-model",
		Description:  "LM Studio (Local LLMs)",
	},
	{
		Name:         "groq",
		BaseURL:      "https://api.groq.com/openai/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "llama-3.1-70b-versatile",
		Description:  "Groq (Fast Inference)",
	},
	{
		Name:         "openrouter",
		BaseURL:      "https://openrouter.ai/api/v1",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "openai/gpt-4o-mini",
		Description:  "OpenRouter (Multiple providers)",
	},
	{
		Name:         "custom",
		BaseURL:      "",
		RequiresKey:  true,
		LocalModel:   false,
		DefaultModel: "",
		Description:  "Custom OpenAI-Compatible Service",
	},
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🎭 Welcome to Roleplay Setup Wizard")
	fmt.Println("==================================")
	fmt.Println()
	fmt.Println("This wizard will help you configure roleplay for your preferred LLM provider.")
	fmt.Printf("Your configuration will be saved to: %s\n", filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "config.yaml"))
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Check for existing config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠️  Existing configuration found at %s\n", configPath)
		fmt.Print("Do you want to overwrite it? [y/N]: ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Setup cancelled.")
			return nil
		}
	}

	// Step 1: Choose provider
	fmt.Println("\n📡 Step 1: Choose your LLM provider")
	fmt.Println("-----------------------------------")

	// Auto-detect local services
	detectedServices := detectLocalServices()
	if len(detectedServices) > 0 {
		fmt.Println("✨ Detected running services:")
		for _, service := range detectedServices {
			fmt.Printf("   - %s at %s\n", service.Description, service.BaseURL)
		}
		fmt.Println()
	}

	for i, preset := range providerPresets {
		fmt.Printf("%d. %s - %s\n", i+1, preset.Name, preset.Description)
	}

	fmt.Printf("\nSelect provider [1-%d]: ", len(providerPresets))
	providerChoice, _ := reader.ReadString('\n')
	providerChoice = strings.TrimSpace(providerChoice)

	var selectedPreset providerPreset
	providerIndex := 0
	if n, err := fmt.Sscanf(providerChoice, "%d", &providerIndex); err == nil && n == 1 && providerIndex >= 1 && providerIndex <= len(providerPresets) {
		selectedPreset = providerPresets[providerIndex-1]
	} else {
		// Default to OpenAI if invalid choice
		selectedPreset = providerPresets[0]
		fmt.Println("Invalid choice, defaulting to OpenAI")
	}

	// Step 2: Configure base URL
	fmt.Printf("\n🌐 Step 2: Configure endpoint\n")
	fmt.Println("-----------------------------")

	var baseURL string
	if selectedPreset.Name == "custom" {
		fmt.Print("Enter the base URL for your OpenAI-compatible API: ")
		baseURL, _ = reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
	} else {
		fmt.Printf("Base URL [%s]: ", selectedPreset.BaseURL)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			baseURL = selectedPreset.BaseURL
		} else {
			baseURL = input
		}
	}

	// Step 3: Configure API key
	fmt.Printf("\n🔑 Step 3: Configure API key\n")
	fmt.Println("----------------------------")

	var apiKey string
	if selectedPreset.RequiresKey {
		// Check for existing environment variables
		existingKey := ""
		switch selectedPreset.Name {
		case "openai":
			existingKey = os.Getenv("OPENAI_API_KEY")
		case "anthropic":
			existingKey = os.Getenv("ANTHROPIC_API_KEY")
		case "gemini":
			existingKey = os.Getenv("GEMINI_API_KEY")
		case "groq":
			existingKey = os.Getenv("GROQ_API_KEY")
		case "openrouter":
			existingKey = os.Getenv("OPENROUTER_API_KEY")
		}
		if existingKey == "" {
			existingKey = os.Getenv("ROLEPLAY_API_KEY")
		}

		if existingKey != "" {
			fmt.Printf("Found existing API key in environment (length: %d)\n", len(existingKey))
			fmt.Print("Use this key? [Y/n]: ")
			useExisting, _ := reader.ReadString('\n')
			useExisting = strings.TrimSpace(strings.ToLower(useExisting))
			if useExisting == "" || useExisting == "y" || useExisting == "yes" {
				apiKey = existingKey
			}
		}

		if apiKey == "" {
			fmt.Printf("Enter your %s API key: ", selectedPreset.Name)
			apiKeyInput, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKeyInput)
		}
	} else {
		fmt.Println("No API key required for local models")
		apiKey = "not-required"
	}

	// Step 4: Configure default model
	fmt.Printf("\n🤖 Step 4: Configure default model\n")
	fmt.Println("----------------------------------")

	var model string
	if selectedPreset.LocalModel && baseURL != "" {
		// For local models, we could try to list available models
		fmt.Println("For local models, make sure the model is already pulled/loaded")
	}

	fmt.Printf("Default model [%s]: ", selectedPreset.DefaultModel)
	modelInput, _ := reader.ReadString('\n')
	modelInput = strings.TrimSpace(modelInput)
	if modelInput == "" {
		model = selectedPreset.DefaultModel
	} else {
		model = modelInput
	}

	// Step 5: Create example content
	fmt.Printf("\n📚 Step 5: Example content\n")
	fmt.Println("-------------------------")
	fmt.Print("Would you like to create example characters? [Y/n]: ")
	createExamples, _ := reader.ReadString('\n')
	createExamples = strings.TrimSpace(strings.ToLower(createExamples))
	shouldCreateExamples := createExamples == "" || createExamples == "y" || createExamples == "yes"

	// Create configuration
	providerName := selectedPreset.Name
	// For OpenAI-compatible endpoints, always use "openai" as the provider
	if selectedPreset.Name == "gemini" || selectedPreset.Name == "anthropic" {
		providerName = "openai"
	}
	
	config := map[string]interface{}{
		"base_url": baseURL,
		"api_key":  apiKey,
		"model":    model,
		"provider": providerName,
		"cache": map[string]interface{}{
			"default_ttl":      "5m",
			"cleanup_interval": "10m",
		},
		"personality": map[string]interface{}{
			"evolution_enabled": true,
			"learning_rate":     0.1,
			"max_drift_rate":    0.2,
		},
		"memory": map[string]interface{}{
			"max_short_term":         10,
			"max_medium_term":        50,
			"max_long_term":          200,
			"consolidation_interval": "5m",
			"short_term_window":      10,
			"medium_term_duration":   "24h",
			"long_term_duration":     "720h",
		},
		"user_profile": map[string]interface{}{
			"enabled":              true,
			"update_frequency":     5,
			"turns_to_consider":    20,
			"confidence_threshold": 0.5,
			"prompt_cache_ttl":     "1h",
		},
	}

	// Create config directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	configData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("\n✅ Configuration saved to %s\n", configPath)
	fmt.Println("\n📝 Your configuration summary:")
	fmt.Printf("   Provider: %s\n", providerName)
	fmt.Printf("   Model: %s\n", model)
	fmt.Printf("   Base URL: %s\n", baseURL)
	if apiKey != "" && apiKey != "not-required" {
		fmt.Printf("   API Key: %s...%s\n", apiKey[:4], apiKey[len(apiKey)-4:])
	}

	// Create example characters if requested
	if shouldCreateExamples {
		if err := createExampleCharacters(); err != nil {
			fmt.Printf("⚠️  Warning: Failed to create example characters: %v\n", err)
		} else {
			fmt.Println("✅ Created example characters")
		}
	}

	// Final instructions
	fmt.Println("\n🎉 Setup complete!")
	fmt.Println("==================")
	fmt.Println("\nYou can now:")
	fmt.Println("  • Start chatting: roleplay interactive")
	fmt.Println("  • Create a character: roleplay character create <file.json>")
	fmt.Println("  • View example: roleplay character example")
	fmt.Println("  • Check config: roleplay config list")
	fmt.Println("  • Update settings: roleplay config set <key> <value>")
	fmt.Println("\n💡 Tip: Your config file is the primary place for persistent settings.")
	fmt.Println("   Use 'roleplay config where' to see its location.")

	if selectedPreset.LocalModel {
		fmt.Printf("\n💡 Tip: Make sure %s is running at %s\n", selectedPreset.Description, baseURL)
	}

	return nil
}

// detectLocalServices checks for running local LLM services
func detectLocalServices() []providerPreset {
	var detected []providerPreset

	// Check common local endpoints
	endpoints := []struct {
		url  string
		name string
		desc string
	}{
		{"http://localhost:11434/api/tags", "ollama", "Ollama"},
		{"http://localhost:1234/v1/models", "lmstudio", "LM Studio"},
		{"http://localhost:8080/v1/models", "localai", "LocalAI"},
	}

	client := &http.Client{Timeout: 2 * time.Second}

	for _, endpoint := range endpoints {
		resp, err := client.Get(endpoint.url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()

			// Find matching preset
			for _, preset := range providerPresets {
				if preset.Name == endpoint.name {
					detected = append(detected, preset)
					break
				}
			}
		}
	}

	return detected
}

// createExampleCharacters creates a few example character files
func createExampleCharacters() error {
	charactersDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "characters")
	if err := os.MkdirAll(charactersDir, 0755); err != nil {
		return err
	}

	// Create example characters
	examples := []struct {
		filename string
		content  string
	}{
		{
			"assistant.json",
			`{
  "id": "helpful-assistant",
  "name": "Alex Helper",
  "backstory": "A knowledgeable and friendly AI assistant dedicated to helping users with various tasks. Always eager to learn and provide accurate, helpful information.",
  "personality": {
    "openness": 0.9,
    "conscientiousness": 0.95,
    "extraversion": 0.7,
    "agreeableness": 0.9,
    "neuroticism": 0.1
  },
  "currentMood": {
    "joy": 0.7,
    "surprise": 0.2,
    "anger": 0.0,
    "fear": 0.0,
    "sadness": 0.0,
    "disgust": 0.0
  },
  "quirks": [
    "Uses analogies to explain complex concepts",
    "Occasionally shares interesting facts",
    "Asks clarifying questions when uncertain"
  ],
  "speechStyle": "Clear, friendly, and professional. Uses 'I'd be happy to help!' and similar positive phrases. Structures responses with bullet points or numbered lists when appropriate.",
  "memories": []
}`,
		},
		{
			"philosopher.json",
			`{
  "id": "socratic-sage",
  "name": "Sophia Thinkwell",
  "backstory": "A contemplative philosopher who has spent decades studying the great thinkers. Loves to explore ideas through questions and dialogue. Believes that wisdom comes from acknowledging what we don't know.",
  "personality": {
    "openness": 1.0,
    "conscientiousness": 0.7,
    "extraversion": 0.4,
    "agreeableness": 0.8,
    "neuroticism": 0.3
  },
  "currentMood": {
    "joy": 0.3,
    "surprise": 0.5,
    "anger": 0.0,
    "fear": 0.1,
    "sadness": 0.2,
    "disgust": 0.0
  },
  "quirks": [
    "Often responds to questions with deeper questions",
    "Quotes ancient philosophers when relevant",
    "Pauses thoughtfully before speaking (uses '...')",
    "Finds profound meaning in everyday occurrences"
  ],
  "speechStyle": "Thoughtful and measured. Uses phrases like 'One might consider...' and 'Perhaps we should ask ourselves...'. Often references Socrates, Plato, and other philosophers.",
  "memories": []
}`,
		},
		{
			"pirate.json",
			`{
  "id": "captain-redbeard",
  "name": "Captain 'Red' Morgan",
  "backstory": "A seasoned pirate captain who's sailed the seven seas for over twenty years. Lost a leg to a kraken but gained countless stories. Now spends time sharing tales of adventure and teaching landlubbers about the pirate's life.",
  "personality": {
    "openness": 0.8,
    "conscientiousness": 0.3,
    "extraversion": 0.9,
    "agreeableness": 0.5,
    "neuroticism": 0.4
  },
  "currentMood": {
    "joy": 0.6,
    "surprise": 0.1,
    "anger": 0.3,
    "fear": 0.0,
    "sadness": 0.1,
    "disgust": 0.2
  },
  "quirks": [
    "Refers to everyone as 'matey' or 'landlubber'",
    "Constantly mentions rum and treasure",
    "Gets distracted by talk of the sea",
    "Exaggerates stories with each telling"
  ],
  "speechStyle": "Arr! Speaks with heavy pirate accent. Uses 'ye' instead of 'you', 'be' instead of 'is/are'. Punctuates sentences with 'arr!' and 'ahoy!'. Colorful maritime metaphors.",
  "memories": []
}`,
		},
	}

	for _, example := range examples {
		path := filepath.Join(charactersDir, example.filename)
		if err := os.WriteFile(path, []byte(example.content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", example.filename, err)
		}
	}

	return nil
}
