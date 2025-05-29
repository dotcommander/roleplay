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
	// Show where config is loaded from
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Configuration file: %s\n\n", configFile)
	} else {
		fmt.Println("No configuration file found, using defaults and environment variables")
		fmt.Println()
	}

	// Key settings to display with their sources
	keySettings := []string{
		"provider",
		"model",
		"api_key",
		"base_url",
	}

	fmt.Println("Active Configuration:")
	fmt.Println("--------------------")

	// Display key settings with sources
	for _, key := range keySettings {
		displaySettingWithSource(key)
	}

	// Display cache settings
	fmt.Println("\nCache Settings:")
	cacheSettings := []string{
		"cache.default_ttl",
		"cache.cleanup_interval",
		"cache.adaptive_ttl",
		"cache.max_entries",
	}
	for _, key := range cacheSettings {
		displaySettingWithSource(key)
	}

	// Display user profile settings
	fmt.Println("\nUser Profile Settings:")
	profileSettings := []string{
		"user_profile.enabled",
		"user_profile.update_frequency",
		"user_profile.turns_to_consider",
		"user_profile.confidence_threshold",
	}
	for _, key := range profileSettings {
		displaySettingWithSource(key)
	}

	// Show which environment variables are set
	fmt.Println("\nEnvironment Variables Detected:")
	fmt.Println("-------------------------------")
	envVars := []string{
		"ROLEPLAY_API_KEY",
		"ROLEPLAY_BASE_URL",
		"ROLEPLAY_MODEL",
		"ROLEPLAY_PROVIDER",
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
			fmt.Printf("%s = %s\n", env, val)
			anySet = true
		}
	}

	if !anySet {
		fmt.Println("(none set)")
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

// displaySettingWithSource shows a setting value and where it came from
func displaySettingWithSource(key string) {
	if !viper.IsSet(key) {
		return
	}

	value := viper.Get(key)
	source := getConfigSource(key)

	// Format the value
	valStr := fmt.Sprintf("%v", value)
	if strings.Contains(strings.ToLower(key), "key") && len(valStr) > 8 {
		valStr = valStr[:4] + "****" + valStr[len(valStr)-4:]
	}

	// Print with aligned formatting
	fmt.Printf("%-30s = %-25s (from %s)\n", key, valStr, source)
}

// getConfigSource determines where a configuration value came from
func getConfigSource(key string) string {
	// Check if it was set by a command-line flag
	if flag := rootCmd.PersistentFlags().Lookup(key); flag != nil && flag.Changed {
		return "command-line flag"
	}

	// Map of keys to their environment variable names
	envMapping := map[string][]string{
		"api_key": {"ROLEPLAY_API_KEY", "OPENAI_API_KEY", "ANTHROPIC_API_KEY", "GEMINI_API_KEY", "GROQ_API_KEY"},
		"base_url": {"ROLEPLAY_BASE_URL", "OPENAI_BASE_URL", "OLLAMA_HOST"},
		"model": {"ROLEPLAY_MODEL"},
		"provider": {"ROLEPLAY_PROVIDER"},
	}

	// Check environment variables
	if envVars, exists := envMapping[key]; exists {
		for _, env := range envVars {
			if os.Getenv(env) != "" {
				return fmt.Sprintf("environment variable %s", env)
			}
		}
	}

	// For nested keys, check the root key env vars
	rootKey := strings.Split(key, ".")[0]
	if envVars, exists := envMapping[rootKey]; exists {
		for _, env := range envVars {
			if os.Getenv(env) != "" {
				return fmt.Sprintf("environment variable %s", env)
			}
		}
	}

	// Check if it's in the config file
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		// Read the config file to check if key exists
		if data, err := os.ReadFile(configFile); err == nil {
			var config map[string]interface{}
			if err := yaml.Unmarshal(data, &config); err == nil {
				if keyExistsInMap(config, key) {
					return "config file"
				}
			}
		}
	}

	// Otherwise it's a default value
	return "default value"
}

// keyExistsInMap checks if a dot-separated key exists in a nested map
func keyExistsInMap(m map[string]interface{}, key string) bool {
	parts := strings.Split(key, ".")
	current := m

	for i, part := range parts {
		if val, exists := current[part]; exists {
			if i == len(parts)-1 {
				return true
			}
			if next, ok := val.(map[string]interface{}); ok {
				current = next
			} else {
				return false
			}
		} else {
			return false
		}
	}

	return false
}
