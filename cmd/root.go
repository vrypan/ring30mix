package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "r30r2",
	Short: "R30R2 - Random Number Generator using Rule 30 Cellular Automaton",
	Long: `R30R2 - Random Number Generator using Rule 30 Cellular Automaton

A deterministic RNG based on 1D cellular automata (Rule 30).
Uses a circular 256-bit strip with radius-2 Rule 30 evolution rules.

Known for generating high-quality pseudo-randomness.
Passes all 319 TestU01 tests including complete BigCrush suite.`,
	// If no subcommand is provided, run the raw command by default
	Run: func(cmd *cobra.Command, args []string) {
		// If no args and no flags, run raw command with default flags
		// This handles: r30r2 (with no arguments)
		rawCmd.Run(rawCmd, args)
	},
}

// Execute runs the root command
func Execute() error {
	// Check if first argument is a known subcommand
	// If not, prepend "raw" to make it the default subcommand
	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		// Check if it's a known subcommand or help/version flag
		if firstArg != "raw" && firstArg != "ascii" &&
		   firstArg != "version" && firstArg != "help" && firstArg != "completion" &&
		   firstArg != "-h" && firstArg != "--help" {
			// Not a subcommand, so prepend "raw"
			os.Args = append([]string{os.Args[0], "raw"}, os.Args[1:]...)
		}
	}

	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(rawCmd)
	rootCmd.AddCommand(asciiCmd)
	rootCmd.AddCommand(versionCmd)
}
