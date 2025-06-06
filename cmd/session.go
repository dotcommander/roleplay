package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/dotcommander/roleplay/internal/repository"
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage conversation sessions",
	Long:  `List, resume, and analyze conversation sessions with cache metrics.`,
}

var sessionListCmd = &cobra.Command{
	Use:   "list [character-id]",
	Short: "List all sessions for a character",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSessionList,
}

var sessionStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show caching statistics across all sessions",
	RunE:  runSessionStats,
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(sessionListCmd)
	sessionCmd.AddCommand(sessionStatsCmd)
}

func runSessionList(cmd *cobra.Command, args []string) error {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo := repository.NewSessionRepository(dataDir)

	if len(args) == 0 {
		// List all characters with sessions
		charRepo, err := repository.NewCharacterRepository(dataDir)
		if err != nil {
			return err
		}

		chars, err := charRepo.ListCharacters()
		if err != nil {
			return err
		}

		fmt.Println("Available characters with sessions:")
		for _, charID := range chars {
			sessions, _ := repo.ListSessions(charID)
			if len(sessions) > 0 {
				char, _ := charRepo.LoadCharacter(charID)
				fmt.Printf("\n%s (%s) - %d sessions\n", char.Name, charID, len(sessions))
			}
		}
		return nil
	}

	// List sessions for specific character
	characterID := args[0]
	sessions, err := repo.ListSessions(characterID)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Printf("No sessions found for character %s\n", characterID)
		return nil
	}

	// Display sessions in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SESSION ID\tSTARTED\tLAST ACTIVE\tMESSAGES\tCACHE HIT RATE")

	for _, session := range sessions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%.1f%%\n",
			session.ID[:8],
			session.StartTime.Format("Jan 2 15:04"),
			formatDuration(time.Since(session.LastActivity)),
			session.MessageCount,
			session.CacheHitRate*100,
		)
	}

	w.Flush()
	return nil
}

func runSessionStats(cmd *cobra.Command, args []string) error {
	dataDir := filepath.Join(os.Getenv("HOME"), ".config", "roleplay")
	repo := repository.NewSessionRepository(dataDir)
	charRepo, err := repository.NewCharacterRepository(dataDir)
	if err != nil {
		return err
	}

	chars, err := charRepo.ListCharacters()
	if err != nil {
		return err
	}

	var totalRequests, totalHits, totalTokensSaved, totalCachedTokens int
	var totalCostSaved, totalPromptCacheRate float64
	var sessionCount int

	fmt.Println("Cache Performance Statistics")
	fmt.Println("===========================")

	for _, charID := range chars {
		sessions, err := repo.ListSessions(charID)
		if err != nil || len(sessions) == 0 {
			continue
		}

		char, _ := charRepo.LoadCharacter(charID)
		fmt.Printf("\n%s (%s):\n", char.Name, charID)

		var charRequests, charHits, charTokensSaved, charCachedTokens int
		var charCostSaved, charPromptCacheRate float64
		var charSessionCount int

		for _, sessionInfo := range sessions {
			session, err := repo.LoadSession(charID, sessionInfo.ID)
			if err != nil {
				continue
			}

			charRequests += session.CacheMetrics.TotalRequests
			charHits += session.CacheMetrics.CacheHits
			charTokensSaved += session.CacheMetrics.TokensSaved
			charCachedTokens += session.CacheMetrics.CachedTokensTotal
			charCostSaved += session.CacheMetrics.CostSaved
			
			if session.CacheMetrics.TotalRequests > 0 {
				charPromptCacheRate += session.CacheMetrics.PromptCacheHitRate
				charSessionCount++
			}
		}

		if charRequests > 0 {
			hitRate := float64(charHits) / float64(charRequests) * 100
			avgPromptCacheRate := 0.0
			if charSessionCount > 0 {
				avgPromptCacheRate = charPromptCacheRate / float64(charSessionCount) * 100
			}
			fmt.Printf("  Sessions: %d\n", len(sessions))
			fmt.Printf("  Total Requests: %d\n", charRequests)
			fmt.Printf("  Response Cache Hit Rate: %.1f%%\n", hitRate)
			
			// Only show OpenAI prompt cache stats if available
			if charCachedTokens > 0 {
				fmt.Printf("  OpenAI Prompt Cache Rate: %.1f%%\n", avgPromptCacheRate)
				fmt.Printf("  Cached Tokens (OpenAI): %d\n", charCachedTokens)
			}
			
			fmt.Printf("  Tokens Saved: %d\n", charTokensSaved)
			fmt.Printf("  Cost Saved: $%.2f\n", charCostSaved)
		}

		totalRequests += charRequests
		totalHits += charHits
		totalTokensSaved += charTokensSaved
		totalCachedTokens += charCachedTokens
		totalCostSaved += charCostSaved
		if charSessionCount > 0 {
			totalPromptCacheRate += charPromptCacheRate
			sessionCount += charSessionCount
		}
	}

	if totalRequests > 0 {
		avgTotalPromptCacheRate := 0.0
		if sessionCount > 0 {
			avgTotalPromptCacheRate = totalPromptCacheRate / float64(sessionCount) * 100
		}
		fmt.Println("\nOverall Statistics:")
		fmt.Printf("  Total Requests: %d\n", totalRequests)
		fmt.Printf("  Overall Response Cache Hit Rate: %.1f%%\n", float64(totalHits)/float64(totalRequests)*100)
		
		// Only show OpenAI prompt cache stats if available
		if totalCachedTokens > 0 {
			fmt.Printf("  Overall OpenAI Prompt Cache Rate: %.1f%%\n", avgTotalPromptCacheRate)
			fmt.Printf("  Total Cached Tokens (OpenAI): %d\n", totalCachedTokens)
		}
		
		fmt.Printf("  Total Tokens Saved: %d\n", totalTokensSaved)
		fmt.Printf("  Total Cost Saved: $%.2f\n", totalCostSaved)
	}

	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(d.Hours()/24))
}
