package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var (
		seed      = flag.Uint64("seed", 0, "RNG seed (default: time-based)")
		bytes     = flag.Int("bytes", 1024, "Number of bytes to generate")
		benchmark = flag.Bool("benchmark", false, "Benchmark mode (measure throughput)")
		help      = flag.Bool("help", false, "Show help")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Turmite RNG - Random Number Generator using LLLR Turmite

A cryptographically-inspired RNG based on cellular automata.
Uses an 8x8 grid with 4 colors and LLLR (Left-Left-Left-Right) movement pattern.

Usage:
  turmite-rng [options]

Seed Format (64 bits):
  bits 0-2:   x position [0-7]
  bits 3-5:   y position [0-7]
  bits 6-7:   direction [0-3] (0=N, 1=E, 2=S, 3=W)
  bits 8-31:  iterations per 32-byte block (default: 1000)
  bits 32-63: unused

Options:
  --seed N        Seed value (default: current time)
  --bytes N       Number of bytes to generate (default: 1024)
  --benchmark     Benchmark throughput instead of generating output
  --help          Show this help

Examples:
  # Generate 1KB of random data
  turmite-rng --bytes 1024 > random.bin

  # Use specific seed
  turmite-rng --seed 12345 --bytes 1048576 > random.bin

  # Benchmark throughput
  turmite-rng --benchmark

  # Test randomness with ent
  turmite-rng --bytes 1048576 | ent

  # Seed breakdown (example):
  # Seed: 0x0000000000001C5  (453 decimal)
  #   x=5, y=0, dir=3 (W), iterations=1

DNA Pattern: 1L.2L.3L.0R (LLLR)
  State 0 → paint 1, turn left
  State 1 → paint 2, turn left
  State 2 → paint 3, turn left
  State 3 → paint 0, turn right
`)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Use time-based seed if not specified
	if *seed == 0 {
		*seed = uint64(time.Now().UnixNano())
	}

	if *benchmark {
		runBenchmark(*seed)
	} else {
		generateBytes(*seed, *bytes)
	}
}

// generateBytes generates and writes random bytes to stdout
func generateBytes(seed uint64, count int) {
	rng := New(seed)

	// Parse seed components for logging
	x := int(seed & 0x7)
	y := int((seed >> 3) & 0x7)
	dir := int((seed >> 6) & 0x3)
	iterations := int((seed >> 8) & 0xFFFFFF)
	if iterations == 0 {
		iterations = 1000
	}

	dirNames := []string{"N", "E", "S", "W"}

	fmt.Fprintf(os.Stderr, "Turmite RNG initialized\n")
	fmt.Fprintf(os.Stderr, "  Seed: 0x%016X (%d)\n", seed, seed)
	fmt.Fprintf(os.Stderr, "  Position: [%d,%d]\n", x, y)
	fmt.Fprintf(os.Stderr, "  Direction: %s\n", dirNames[dir])
	fmt.Fprintf(os.Stderr, "  Iterations: %d per 32-byte block\n", iterations)
	fmt.Fprintf(os.Stderr, "  DNA: 1L.2L.3L.0R (LLLR)\n")
	fmt.Fprintf(os.Stderr, "Generating %d bytes...\n", count)

	buf := make([]byte, count)
	n, err := rng.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write to stdout
	written, err := os.Stdout.Write(buf[:n])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Generated %d bytes\n", written)
}

// runBenchmark measures RNG throughput
func runBenchmark(seed uint64) {
	rng := New(seed)

	sizes := []int{1024, 10240, 102400, 1048576} // 1KB, 10KB, 100KB, 1MB

	fmt.Println("Turmite RNG Benchmark")
	fmt.Printf("Seed: 0x%016X\n", seed)
	fmt.Println()
	fmt.Println("Size\t\tTime\t\tThroughput")
	fmt.Println("----\t\t----\t\t----------")

	for _, size := range sizes {
		buf := make([]byte, size)

		start := time.Now()
		n, err := rng.Read(buf)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		throughput := float64(n) / elapsed.Seconds() / 1024 / 1024 // MB/s

		sizeStr := formatSize(size)
		fmt.Printf("%s\t\t%v\t%.2f MB/s\n", sizeStr, elapsed.Round(time.Millisecond), throughput)
	}
}

// formatSize formats byte count for display
func formatSize(bytes int) string {
	if bytes >= 1048576 {
		return fmt.Sprintf("%d MB", bytes/1048576)
	} else if bytes >= 1024 {
		return fmt.Sprintf("%d KB", bytes/1024)
	}
	return fmt.Sprintf("%d B", bytes)
}
