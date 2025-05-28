package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage roleplay configuration",
	Long:  `View and modify roleplay configuration settings`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings",
	Long:  `Display all current configuration values, showing the merged result from config file, environment variables, and defaults`,
	RunE:  runConfigList,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a specific configuration value",
	Long:  `Retrieve the value of a specific configuration key`,
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Update a configuration value in the config file`,
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

var configWhereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show configuration file location",
	Long:  `Display the path to the active configuration file`,
	RunE:  runConfigWhere,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configWhereCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	// Get all settings
	settings := viper.AllSettings()

	// Show where config is loaded from
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Configuration file: %s\n", configFile)
	} else {
		fmt.Println("No configuration file found, using defaults and environment variables")
	}
	fmt.Println()

	// Display settings in a readable format
	displaySettings(settings, "")

	// Show which environment variables are set
	fmt.Println("\nEnvironment variables:")
	envVars := []string{
		"ROLEPLAY_API_KEY",
		"ROLEPLAY_BASE_URL",
		"ROLEPLAY_MODEL",
		"ROLEPLAY_DEFAULT_PROVIDER",
		"OPENAI_API_KEY",
		"OPENAI_BASE_URL",
		"ANTHROPIC_API_KEY",
		"GEMINI_API_KEY",
		"GROQ_API_KEY",
		"OLLAMA_HOST",
	}

	anySet := false
	for _, env := range envVars {
		if val := os.Getenv(env); val != "" {
			// Mask API keys
			if strings.Contains(env, "KEY") && len(val) > 8 {
				val = val[:4] + "****" + val[len(val)-4:]
			}
			fmt.Printf("  %s = %s\n", env, val)
			anySet = true
		}
	}

	if !anySet {
		fmt.Println("  (none set)")
	}

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	if !viper.IsSet(key) {
		return fmt.Errorf("configuration key '%s' not found", key)
	}

	value := viper.Get(key)

	// Mask API keys when displaying
	if strings.Contains(strings.ToLower(key), "key") || strings.Contains(strings.ToLower(key), "api_key") {
		if str, ok := value.(string); ok && len(str) > 8 {
			value = str[:4] + "****" + str[len(str)-4:]
		}
	}

	// Print just the value for easy scripting
	fmt.Println(value)

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Set the value in viper
	viper.Set(key, value)

	// Get the config file path
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// Create default config path
		configDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		configFile = filepath.Join(configDir, "config.yaml")
	}

	// Read existing config or create new
	configData := make(map[string]interface{})
	if data, err := os.ReadFile(configFile); err == nil {
		if err := yaml.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse existing config: %w", err)
		}
	}

	// Update the specific key
	setNestedValue(configData, key, value)

	// Write back to file
	data, err := yaml.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✅ Set %s = %s\n", key, value)
	fmt.Printf("Configuration saved to %s\n", configFile)

	return nil
}

func runConfigWhere(cmd *cobra.Command, args []string) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = filepath.Join(os.Getenv("HOME"), ".config", "roleplay", "config.yaml")
		fmt.Printf("Default location (not created yet): %s\n", configFile)
	} else {
		fmt.Println(configFile)
	}
	return nil
}

// displaySettings recursively displays configuration settings
func displaySettings(settings map[string]interface{}, prefix string) {
	for key, value := range settings {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Printf("%s:\n", fullKey)
			displaySettings(v, fullKey)
		case string:
			// Mask sensitive values
			if strings.Contains(strings.ToLower(key), "key") && len(v) > 8 {
				v = v[:4] + "****" + v[len(v)-4:]
			}
			fmt.Printf("  %s = %s\n", fullKey, v)
		default:
			fmt.Printf("  %s = %v\n", fullKey, v)
		}
	}
}

// setNestedValue sets a value in a nested map using dot notation
func setNestedValue(m map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, ".")
	current := m

	// Navigate to the nested location
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			// Key exists but is not a map, overwrite it
			current[part] = make(map[string]interface{})
			current = current[part].(map[string]interface{})
		}
	}

	// Set the final value
	finalKey := parts[len(parts)-1]

	// Try to parse value to appropriate type
	switch {
	case value == "true":
		current[finalKey] = true
	case value == "false":
		current[finalKey] = false
	case isNumeric(value.(string)):
		// Keep as string for now, let viper handle type conversion
		current[finalKey] = value
	default:
		current[finalKey] = value
	}
}

// isNumeric checks if a string represents a number
func isNumeric(s string) bool {
	// Simple check - could be enhanced
	for _, c := range s {
		if (c < '0' || c > '9') && c != '.' && c != '-' {
			return false
		}
	}
	return true
}
