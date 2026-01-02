package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version = "1.0.0"
	GitRepo = "github.com/vrypan/r30r2"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version number and build information for r30r2.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("r30r2 version %s\n", Version)
		fmt.Printf("Random Number Generator using Rule 30 Cellular Automaton (Radius-2)\n")
		fmt.Printf("\n")
		fmt.Printf("Repository: %s\n", GitRepo)
		fmt.Printf("Statistical Quality: 160/160 BigCrush tests passed\n")
		fmt.Printf("Performance: 3.86× faster than math/rand for bulk operations\n")
	},
}
