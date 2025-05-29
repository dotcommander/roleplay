package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information - set at build time
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version of roleplay along with build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Roleplay - AI Character Chat System\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}