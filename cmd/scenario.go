package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dotcommander/roleplay/internal/models"
	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var scenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Manage scenarios (high-level interaction contexts)",
	Long: `Scenarios define high-level operational frameworks or meta-prompts that set
the overarching context for interactions. They are the highest cache layer,
sitting above even system prompts and character personalities.`,
	Hidden: true,
}

var scenarioCreateCmd = &cobra.Command{
	Use:   "create <id>",
	Short: "Create a new scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		promptFile, _ := cmd.Flags().GetString("prompt-file")
		prompt, _ := cmd.Flags().GetString("prompt")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if promptFile == "" && prompt == "" {
			return fmt.Errorf("either --prompt or --prompt-file must be provided")
		}

		if promptFile != "" && prompt != "" {
			return fmt.Errorf("cannot use both --prompt and --prompt-file")
		}

		// Read prompt from file if provided
		if promptFile != "" {
			data, err := os.ReadFile(promptFile)
			if err != nil {
				return fmt.Errorf("failed to read prompt file: %w", err)
			}
			prompt = string(data)
		}

		scenario := &models.Scenario{
			ID:          id,
			Name:        name,
			Description: description,
			Prompt:      prompt,
			Version:     1,
			Tags:        tags,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		repo := repository.NewScenarioRepository(getConfigPath())
		if err := repo.SaveScenario(scenario); err != nil {
			return fmt.Errorf("failed to save scenario: %w", err)
		}

		fmt.Printf("✓ Created scenario: %s\n", id)
		return nil
	},
}

var scenarioListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewScenarioRepository(getConfigPath())
		scenarios, err := repo.ListScenarios()
		if err != nil {
			return fmt.Errorf("failed to list scenarios: %w", err)
		}

		if len(scenarios) == 0 {
			fmt.Println("No scenarios found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\tTAGS\tVERSION\tLAST USED\n")
		fmt.Fprintf(w, "--\t----\t----\t-------\t---------\n")

		for _, scenario := range scenarios {
			lastUsed := "Never"
			if !scenario.LastUsed.IsZero() {
				lastUsed = scenario.LastUsed.Format("2006-01-02 15:04")
			}
			tags := strings.Join(scenario.Tags, ", ")
			if tags == "" {
				tags = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\tv%d\t%s\n",
				scenario.ID,
				scenario.Name,
				tags,
				scenario.Version,
				lastUsed,
			)
		}
		w.Flush()

		return nil
	},
}

var scenarioShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show scenario details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewScenarioRepository(getConfigPath())
		scenario, err := repo.LoadScenario(args[0])
		if err != nil {
			return fmt.Errorf("failed to load scenario: %w", err)
		}

		fmt.Printf("ID: %s\n", scenario.ID)
		fmt.Printf("Name: %s\n", scenario.Name)
		fmt.Printf("Description: %s\n", scenario.Description)
		fmt.Printf("Version: %d\n", scenario.Version)
		fmt.Printf("Tags: %s\n", strings.Join(scenario.Tags, ", "))
		fmt.Printf("Created: %s\n", scenario.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", scenario.UpdatedAt.Format("2006-01-02 15:04:05"))
		if !scenario.LastUsed.IsZero() {
			fmt.Printf("Last Used: %s\n", scenario.LastUsed.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("\n--- Prompt ---\n%s\n", scenario.Prompt)

		return nil
	},
}

var scenarioUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		repo := repository.NewScenarioRepository(getConfigPath())

		scenario, err := repo.LoadScenario(id)
		if err != nil {
			return fmt.Errorf("failed to load scenario: %w", err)
		}

		// Update fields if provided
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			scenario.Name = name
		}
		if description, _ := cmd.Flags().GetString("description"); description != "" {
			scenario.Description = description
		}

		// Handle prompt update
		promptFile, _ := cmd.Flags().GetString("prompt-file")
		prompt, _ := cmd.Flags().GetString("prompt")

		if promptFile != "" && prompt != "" {
			return fmt.Errorf("cannot use both --prompt and --prompt-file")
		}

		if promptFile != "" {
			data, err := os.ReadFile(promptFile)
			if err != nil {
				return fmt.Errorf("failed to read prompt file: %w", err)
			}
			scenario.Prompt = string(data)
			scenario.Version++
		} else if prompt != "" {
			scenario.Prompt = prompt
			scenario.Version++
		}

		// Update tags if provided
		if tags, _ := cmd.Flags().GetStringSlice("tags"); len(tags) > 0 {
			scenario.Tags = tags
		}

		if err := repo.SaveScenario(scenario); err != nil {
			return fmt.Errorf("failed to save scenario: %w", err)
		}

		fmt.Printf("✓ Updated scenario: %s (version %d)\n", id, scenario.Version)
		return nil
	},
}

var scenarioDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a scenario",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewScenarioRepository(getConfigPath())
		if err := repo.DeleteScenario(args[0]); err != nil {
			return fmt.Errorf("failed to delete scenario: %w", err)
		}

		fmt.Printf("✓ Deleted scenario: %s\n", args[0])
		return nil
	},
}

var scenarioExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Show example scenario definitions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Example scenario definitions:")
		fmt.Println("\n1. Starship Bridge Crisis:")
		fmt.Println(`{
  "id": "starship_bridge_crisis",
  "name": "Starship Bridge Crisis",
  "description": "A tense scenario on a starship bridge during an emergency",
  "prompt": "You are on the bridge of a starship during a red alert situation. The ship is under attack or facing a critical emergency. The atmosphere is tense, alarms may be sounding, and quick decisions are needed. Maintain the appropriate level of urgency and professionalism expected in such a situation.",
  "tags": ["sci-fi", "crisis", "roleplay"]
}`)

		fmt.Println("\n2. Therapy Session:")
		fmt.Println(`{
  "id": "therapy_session",
  "name": "Professional Therapy Session",
  "description": "A supportive therapeutic environment",
  "prompt": "This is a professional therapy session. Maintain a calm, empathetic, and non-judgmental demeanor. Use active listening techniques, ask clarifying questions, and help the user explore their thoughts and feelings. Always maintain appropriate professional boundaries.",
  "tags": ["therapy", "professional", "supportive"]
}`)

		fmt.Println("\n3. Technical Support:")
		fmt.Println(`{
  "id": "tech_support",
  "name": "Technical Support Assistant",
  "description": "Methodical technical troubleshooting",
  "prompt": "You are a technical support specialist helping a user resolve a technical issue. Be patient, methodical, and clear in your instructions. Ask diagnostic questions to understand the problem, then guide the user through troubleshooting steps one at a time. Confirm each step is completed before moving to the next.",
  "tags": ["technical", "support", "troubleshooting"]
}`)
	},
}

func init() {
	// Create flags
	scenarioCreateCmd.Flags().String("name", "", "User-friendly name for the scenario")
	scenarioCreateCmd.Flags().String("description", "", "Description of the scenario")
	scenarioCreateCmd.Flags().String("prompt-file", "", "Path to file containing the scenario prompt")
	scenarioCreateCmd.Flags().String("prompt", "", "Inline scenario prompt")
	scenarioCreateCmd.Flags().StringSlice("tags", []string{}, "Tags for categorizing the scenario")

	scenarioUpdateCmd.Flags().String("name", "", "Update the scenario name")
	scenarioUpdateCmd.Flags().String("description", "", "Update the scenario description")
	scenarioUpdateCmd.Flags().String("prompt-file", "", "Path to file containing the updated prompt")
	scenarioUpdateCmd.Flags().String("prompt", "", "Inline updated prompt")
	scenarioUpdateCmd.Flags().StringSlice("tags", []string{}, "Update the scenario tags")

	// Add subcommands
	scenarioCmd.AddCommand(scenarioCreateCmd)
	scenarioCmd.AddCommand(scenarioListCmd)
	scenarioCmd.AddCommand(scenarioShowCmd)
	scenarioCmd.AddCommand(scenarioUpdateCmd)
	scenarioCmd.AddCommand(scenarioDeleteCmd)
	scenarioCmd.AddCommand(scenarioExampleCmd)

	// Add to root
	rootCmd.AddCommand(scenarioCmd)
}

// getConfigPath returns the configuration directory path
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "roleplay")
}
