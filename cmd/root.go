package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dotcommander/roleplay/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "roleplay",
	Short: "A sophisticated character bot with psychological modeling",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Set custom usage function for root command only
	defaultUsageFunc := rootCmd.UsageFunc()
	rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		if cmd == rootCmd {
			// Custom help for root command
			fmt.Println(`🎭 Roleplay - AI Character Chat System

Roleplay is a character bot system that implements psychologically-realistic 
AI characters with personality evolution, emotional states, and multi-layered memory systems.

Quick Start:
  roleplay quickstart              Start chatting immediately (no setup required)
  roleplay setup                   Interactive setup wizard for configuration

Chat Commands:
  roleplay chat <message>          Send a single message to a character (alias: c)
  roleplay interactive             Start interactive chat session (alias: i)

Character Management:
  roleplay character list          List available characters (alias: ls)
  roleplay character create        Create from JSON file
  roleplay character import        Import from markdown file

Configuration:
  roleplay config status           Show current settings
  roleplay config test             Test API connection
  roleplay session list            View chat history
  roleplay profile show            Manage user profiles

Use "roleplay [command] --help" for more information about a command.`)
			return nil
		}
		// Use default for other commands
		return defaultUsageFunc(cmd)
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.roleplay.yaml)")
	rootCmd.PersistentFlags().String("provider", "openai", "AI provider to use (anthropic, openai)")
	rootCmd.PersistentFlags().String("model", "", "Model to use (e.g., gpt-4o-mini, gpt-4 for OpenAI)")
	rootCmd.PersistentFlags().String("api-key", "", "API key for the AI provider")
	rootCmd.PersistentFlags().String("base-url", "", "Base URL for OpenAI-compatible API")
	rootCmd.PersistentFlags().Duration("cache-ttl", 10*time.Minute, "Default cache TTL")
	rootCmd.PersistentFlags().Bool("adaptive-ttl", true, "Enable adaptive TTL for cache")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	if err := viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding provider flag: %v\n", err)
	}
	if err := viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding model flag: %v\n", err)
	}
	if err := viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding api_key flag: %v\n", err)
	}
	if err := viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding base_url flag: %v\n", err)
	}
	if err := viper.BindPFlag("cache.default_ttl", rootCmd.PersistentFlags().Lookup("cache-ttl")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding cache.default_ttl flag: %v\n", err)
	}
	if err := viper.BindPFlag("cache.adaptive_ttl", rootCmd.PersistentFlags().Lookup("adaptive-ttl")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding cache.adaptive_ttl flag: %v\n", err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, ".config", "roleplay"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("ROLEPLAY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Resolve configuration with proper priority: Flags > Config File > Environment Variables
	profileName := viper.GetString("provider")

	// API Key Resolution
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		// Check ROLEPLAY_API_KEY first
		apiKey = os.Getenv("ROLEPLAY_API_KEY")
		if apiKey == "" {
			// Check provider-specific environment variables
			switch profileName {
			case "openai":
				apiKey = os.Getenv("OPENAI_API_KEY")
			case "anthropic", "anthropic_compatible":
				apiKey = os.Getenv("ANTHROPIC_API_KEY")
			case "gemini", "gemini_compatible":
				apiKey = os.Getenv("GEMINI_API_KEY")
			case "groq":
				apiKey = os.Getenv("GROQ_API_KEY")
			}
		}
	}

	// Base URL Resolution
	baseURL := viper.GetString("base_url")
	if baseURL == "" {
		// Check ROLEPLAY_BASE_URL first
		baseURL = os.Getenv("ROLEPLAY_BASE_URL")
		if baseURL == "" {
			// Check common environment variables
			baseURL = os.Getenv("OPENAI_BASE_URL")
			if baseURL == "" && (profileName == "ollama" || os.Getenv("OLLAMA_HOST") != "") {
				// Handle Ollama special case
				ollamaHost := os.Getenv("OLLAMA_HOST")
				if ollamaHost == "" {
					ollamaHost = "http://localhost:11434"
				}
				baseURL = ollamaHost + "/v1"
			}
		}
	}

	// Model Resolution
	model := viper.GetString("model")
	if model == "" {
		// Check ROLEPLAY_MODEL environment variable
		model = os.Getenv("ROLEPLAY_MODEL")
	}

	cfg = &config.Config{
		DefaultProvider: viper.GetString("provider"),
		Model:           model,
		APIKey:          apiKey,
		BaseURL:         baseURL,
		ModelAliases:    viper.GetStringMapString("model_aliases"),
		CacheConfig: config.CacheConfig{
			MaxEntries:        viper.GetInt("cache.max_entries"),
			CleanupInterval:   viper.GetDuration("cache.cleanup_interval"),
			DefaultTTL:        viper.GetDuration("cache.default_ttl"),
			EnableAdaptiveTTL: viper.GetBool("cache.adaptive_ttl"),
		},
		MemoryConfig: config.MemoryConfig{
			ShortTermWindow:    viper.GetInt("memory.short_term_window"),
			MediumTermDuration: viper.GetDuration("memory.medium_term_duration"),
			ConsolidationRate:  viper.GetFloat64("memory.consolidation_rate"),
		},
		PersonalityConfig: config.PersonalityConfig{
			EvolutionEnabled:   viper.GetBool("personality.evolution_enabled"),
			MaxDriftRate:       viper.GetFloat64("personality.max_drift_rate"),
			StabilityThreshold: viper.GetFloat64("personality.stability_threshold"),
		},
		UserProfileConfig: config.UserProfileConfig{
			Enabled:             viper.GetBool("user_profile.enabled"),
			UpdateFrequency:     viper.GetInt("user_profile.update_frequency"),
			TurnsToConsider:     viper.GetInt("user_profile.turns_to_consider"),
			ConfidenceThreshold: viper.GetFloat64("user_profile.confidence_threshold"),
			PromptCacheTTL:      viper.GetDuration("user_profile.prompt_cache_ttl"),
		},
	}

	// Set defaults if not configured
	if cfg.CacheConfig.MaxEntries == 0 {
		cfg.CacheConfig.MaxEntries = 10000
	}
	if cfg.CacheConfig.CleanupInterval == 0 {
		cfg.CacheConfig.CleanupInterval = 5 * time.Minute
	}
	if cfg.MemoryConfig.ShortTermWindow == 0 {
		cfg.MemoryConfig.ShortTermWindow = 20
	}
	if cfg.MemoryConfig.MediumTermDuration == 0 {
		cfg.MemoryConfig.MediumTermDuration = 24 * time.Hour
	}
	if cfg.MemoryConfig.ConsolidationRate == 0 {
		cfg.MemoryConfig.ConsolidationRate = 0.1
	}
	if cfg.PersonalityConfig.MaxDriftRate == 0 {
		cfg.PersonalityConfig.MaxDriftRate = 0.02
	}
	if cfg.PersonalityConfig.StabilityThreshold == 0 {
		cfg.PersonalityConfig.StabilityThreshold = 10
	}

	// Set defaults for UserProfileConfig
	if cfg.UserProfileConfig.UpdateFrequency == 0 {
		cfg.UserProfileConfig.UpdateFrequency = 5 // Update every 5 messages
	}
	if cfg.UserProfileConfig.TurnsToConsider == 0 {
		cfg.UserProfileConfig.TurnsToConsider = 20 // Analyze last 20 turns
	}
	if cfg.UserProfileConfig.ConfidenceThreshold == 0 {
		cfg.UserProfileConfig.ConfidenceThreshold = 0.5 // Include facts with >50% confidence
	}
	if cfg.UserProfileConfig.PromptCacheTTL == 0 {
		cfg.UserProfileConfig.PromptCacheTTL = 1 * time.Hour // Cache user profiles for 1 hour
	}
}

func GetConfig() *config.Config {
	return cfg
}
