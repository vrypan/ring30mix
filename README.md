# Turmite RNG

A random number generator based on cellular automata, specifically using an LLLR (Left-Left-Left-Right) turmite pattern on an 8×8 grid with 4 colors.

## Features

- **Deterministic**: Same seed produces same output
- **Fast**: ~20-27 MB/s throughput
- **Compact**: 8×8 grid with 2-bit cells (16 bytes storage)
- **Standard Interface**: Implements `io.Reader` interface
- **CLI Tool**: Generate random bytes from command line
- **Library**: Use as Go package in your projects

## Installation

```bash
go build -o turmite-rng
```

## CLI Usage

### Basic Examples

```bash
# Generate 1KB of random data
./turmite-rng --bytes 1024 > random.bin

# Use specific seed
./turmite-rng --seed=12345 --bytes=1048576 > random.bin

# Benchmark throughput
./turmite-rng --benchmark

# Pipe to analysis tools
./turmite-rng --bytes 1048576 | ent
./turmite-rng --bytes 10000000 | dieharder -a -g 200
```

### CLI Options

| Option | Description | Default |
|--------|-------------|---------|
| `--seed N` | RNG seed (64-bit unsigned integer) | Current time |
| `--bytes N` | Number of bytes to generate | 1024 |
| `--benchmark` | Run throughput benchmark | false |
| `--help` | Show help message | false |

### Seed Format

The 64-bit seed encodes all RNG parameters:

```
Bits 0-2:   x position [0-7]
Bits 3-5:   y position [0-7]
Bits 6-7:   direction [0-3] (0=N, 1=E, 2=S, 3=W)
Bits 8-31:  iterations per 32-byte block (default: 1000 if 0)
Bits 32-63: unused (reserved for future use)
```

**Example seed breakdown:**
```bash
$ ./turmite-rng --seed=12345 --bytes=100 > /dev/null

Seed: 0x0000000000003039 (12345)
  Position: [1,7]
  Direction: N
  Iterations: 48 per 32-byte block
```

### Calculating Seeds

```bash
# Seed formula: x + (y << 3) + (dir << 6) + (iterations << 8)

# Example: x=2, y=3, dir=1 (East), iterations=5000
seed=$((2 + (3 << 3) + (1 << 6) + (5000 << 8)))
./turmite-rng --seed=$seed --bytes=1024 > output.bin
```

## API Usage

### As io.Reader

```go
package main

import (
    "fmt"
    turmiterng "path/to/turmite-rng"
)

func main() {
    // Create RNG with seed
    rng := turmiterng.New(12345)

    // Read random bytes
    buf := make([]byte, 1024)
    n, err := rng.Read(buf)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated %d random bytes\n", n)
}
```

### As crypto/rand.Reader Replacement

```go
package main

import (
    "fmt"
    turmiterng "path/to/turmite-rng"
)

func main() {
    // Initialize global Reader
    turmiterng.InitReader(12345)

    // Use like crypto/rand.Reader
    buf := make([]byte, 32)
    n, err := turmiterng.Reader.Read(buf)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated %d bytes\n", n)
}
```

### Stream Processing

```go
package main

import (
    "encoding/hex"
    "fmt"
    turmiterng "path/to/turmite-rng"
)

func main() {
    rng := turmiterng.New(54321)

    // Generate multiple blocks
    for i := 0; i < 5; i++ {
        buf := make([]byte, 16)
        rng.Read(buf)
        fmt.Printf("Block %d: %s\n", i, hex.EncodeToString(buf))
    }
}
```

### Custom Iterations

```go
package main

import (
    turmiterng "path/to/turmite-rng"
)

func main() {
    // More iterations = more mixing (slower but potentially better randomness)
    // Seed with 10,000 iterations per 32-byte block
    seed := uint64(5) + (3 << 3) + (2 << 6) + (10000 << 8)
    rng := turmiterng.New(seed)

    buf := make([]byte, 1024)
    rng.Read(buf)
}
```

## Performance

Benchmarks on typical hardware:

| Block Size | Time | Throughput |
|------------|------|------------|
| 1 KB | <1ms | ~15 MB/s |
| 10 KB | <1ms | ~25 MB/s |
| 100 KB | ~5ms | ~20 MB/s |
| 1 MB | ~40ms | ~27 MB/s |

**Factors affecting performance:**
- **Iterations**: More iterations per block = slower but more mixing
- **Block size**: Larger reads amortize setup overhead
- **CPU**: Performance scales with CPU speed

**Typical configurations:**
```bash
# Fast (low iterations)
--seed=$((0 + (0 << 3) + (0 << 6) + (100 << 8)))   # ~100 MB/s

# Balanced (default)
--seed=$((0 + (0 << 3) + (0 << 6) + (1000 << 8)))  # ~25 MB/s

# High quality (many iterations)
--seed=$((0 + (0 << 3) + (0 << 6) + (10000 << 8))) # ~3 MB/s
```

## DNA Pattern

The turmite follows the **LLLR** (Left-Left-Left-Right) pattern:

```
State 0 → paint 1, turn left
State 1 → paint 2, turn left
State 2 → paint 3, turn left
State 3 → paint 0, turn right
```

This pattern creates complex, chaotic behavior suitable for random number generation.

## Testing Randomness

### Quick test with ent

```bash
./turmite-rng --bytes 1048576 | ent

# Good output should show:
# - Entropy: ~7.99+ bits per byte
# - Chi-square: passes (>0.01, <0.99)
# - Monte Carlo: π estimate ~3.14
```

### Comprehensive test with dieharder

```bash
# Generate 10MB test file
./turmite-rng --bytes 10000000 > testdata.bin

# Run dieharder test suite
dieharder -a -g 201 -f testdata.bin
```

### Visual test

```bash
# Generate PGM image to visualize randomness
./turmite-rng --bytes 65536 | convert -size 256x256 -depth 8 gray:- output.png
```

## Use Cases

**Good for:**
- ✓ Simulations requiring reproducible randomness
- ✓ Testing and debugging (deterministic)
- ✓ Procedural generation with seeds
- ✓ Educational/research into CA-based RNGs

**Not recommended for:**
- ✗ Cryptographic keys (not cryptographically secure)
- ✗ Security tokens or passwords
- ✗ High-stakes gambling/lottery

## Comparison

| RNG | Speed | Crypto-secure | Deterministic | Quality |
|-----|-------|---------------|---------------|---------|
| `crypto/rand` | Fast | ✓ Yes | ✗ No | Excellent |
| `math/rand` | Very Fast | ✗ No | ✓ Yes | Good |
| `turmite-rng` | Medium | ✗ No | ✓ Yes | Good* |

*Quality depends on iteration count; needs statistical validation for your use case.

## Algorithm

See [ALGORITHM.md](ALGORITHM.md) for detailed explanation of the internal implementation.

## License

Same as the parent turmites project.

## References

- [Turmite on Wikipedia](https://en.wikipedia.org/wiki/Turmite)
- [Langton's Ant](https://en.wikipedia.org/wiki/Langton%27s_ant) (similar concept)
- Parent project: [go-turmites](../)
