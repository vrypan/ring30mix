package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/r30r2/rand"
)

var (
	rawSeed  uint64
	rawBytes int
)

var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Generate raw random bytes",
	Long: `Generate raw random bytes to stdout.

This is the default subcommand if none is specified.

Examples:
  # Generate 1KB of random data
  r30r2 raw --bytes 1024 > random.bin

  # Use specific seed
  r30r2 raw --seed 12345 --bytes 1048576 > random.bin

  # Generate specific size with dd
  r30r2 raw --bytes 1073741824 | dd of=test.data bs=1m

  # Unlimited streaming (use with head, pv, or Ctrl+C)
  r30r2 raw --bytes 0 | head -c 1073741824 > test.data

  # Test randomness with ent
  r30r2 raw --bytes 1048576 | ent

  # Default behavior (no subcommand)
  r30r2 --bytes 1024 > random.bin`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use time-based seed if not specified
		if rawSeed == 0 {
			rawSeed = uint64(time.Now().UnixNano())
		}

		generateBytes(rawSeed, rawBytes)
	},
}

func init() {
	rawCmd.Flags().Uint64Var(&rawSeed, "seed", 0, "RNG seed (default: time-based)")
	rawCmd.Flags().IntVar(&rawBytes, "bytes", 1024, "Number of bytes to generate (0 = unlimited)")
}

// generateBytes generates and writes random bytes to stdout
func generateBytes(seed uint64, count int) {
	rng := rand.New(seed)

	if count == 0 {
		// Unlimited mode: stream chunks until pipe breaks
		buf := make([]byte, 1024*1024) // 1MB chunks
		for {
			n, err := rng.Read(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Write all bytes, handling partial writes
			written := 0
			for written < n {
				w, err := os.Stdout.Write(buf[written:n])
				if err != nil {
					// Pipe closed (e.g., dd finished) - exit gracefully
					os.Exit(0)
				}
				written += w
			}
		}
	} else {
		// Fixed size: stream in chunks to avoid huge allocations
		const chunkSize = 1024 * 1024 // 1MB chunks
		buf := make([]byte, chunkSize)
		remaining := count

		for remaining > 0 {
			toRead := chunkSize
			if remaining < chunkSize {
				toRead = remaining
			}

			n, err := rng.Read(buf[:toRead])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Write all bytes from this read, handling partial writes
			written := 0
			for written < n {
				w, err := os.Stdout.Write(buf[written:n])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
					os.Exit(1)
				}
				written += w
			}

			remaining -= n
		}
	}
}
