package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage user profiles",
	Long:  `View, list, and delete AI-extracted user profiles that characters maintain about users.`,
}

var profileShowCmd = &cobra.Command{
	Use:   "show <user-id> <character-id>",
	Short: "Show a specific user profile",
	Long:  `Display the AI-extracted profile that a character has built about a user.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		characterID := args[1]

		home, _ := os.UserHomeDir()
		profilesDir := filepath.Join(home, ".config", "roleplay", "user_profiles")
		repo := repository.NewUserProfileRepository(profilesDir)

		profile, err := repo.LoadUserProfile(userID, characterID)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no profile found for user '%s' with character '%s'", userID, characterID)
			}
			return fmt.Errorf("failed to load profile: %w", err)
		}

		// Pretty print the profile
		data, err := json.MarshalIndent(profile, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format profile: %w", err)
		}

		fmt.Println(string(data))
		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list <user-id>",
	Short: "List all profiles for a user",
	Long:  `Display all character profiles that exist for a specific user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]

		home, _ := os.UserHomeDir()
		profilesDir := filepath.Join(home, ".config", "roleplay", "user_profiles")
		repo := repository.NewUserProfileRepository(profilesDir)

		profiles, err := repo.ListUserProfiles(userID)
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		if len(profiles) == 0 {
			fmt.Printf("No profiles found for user '%s'\n", userID)
			return nil
		}

		fmt.Printf("Profiles for user '%s':\n\n", userID)
		for _, profile := range profiles {
			fmt.Printf("Character: %s\n", profile.CharacterID)
			fmt.Printf("  Version: %d\n", profile.Version)
			fmt.Printf("  Last Analyzed: %s\n", profile.LastAnalyzed.Format("2006-01-02 15:04:05"))
			if profile.OverallSummary != "" {
				fmt.Printf("  Summary: %s\n", profile.OverallSummary)
			}
			if profile.InteractionStyle != "" {
				fmt.Printf("  Interaction Style: %s\n", profile.InteractionStyle)
			}
			fmt.Printf("  Facts Count: %d\n", len(profile.Facts))
			fmt.Println()
		}

		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <user-id> <character-id>",
	Short: "Delete a user profile",
	Long:  `Remove the profile that a character has built about a user.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		characterID := args[1]

		// Confirm deletion
		if !force {
			fmt.Printf("Are you sure you want to delete the profile for user '%s' with character '%s'? (y/N): ", userID, characterID)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil || (response != "y" && response != "Y") {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		home, _ := os.UserHomeDir()
		profilesDir := filepath.Join(home, ".config", "roleplay", "user_profiles")
		repo := repository.NewUserProfileRepository(profilesDir)

		err := repo.DeleteUserProfile(userID, characterID)
		if err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		fmt.Printf("Profile for user '%s' with character '%s' has been deleted.\n", userID, characterID)
		return nil
	},
}

var force bool

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDeleteCmd)

	// Add force flag to delete command
	profileDeleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
}