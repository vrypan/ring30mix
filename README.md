# Rule 30 RND

High-performance pseudo-random number generator based on Rule 30 cellular automaton. **40% faster** than Go's standard library, with perfect entropy (8.0 bits/byte) and identical statistical quality to math/rand.

## Quick Start

### Installation

```bash
# As a command-line tool
go install github.com/vrypan/rule30rnd@latest

# As a library
go get github.com/vrypan/rule30rnd
```

### Command Line Usage

```bash
# Generate random data
./rule30 --bytes=1048576 > random.bin

# Generate specific size with dd
./rule30 --bytes=1073741824 | dd of=test.data bs=1m

# Unlimited streaming (use with head, pv, or Ctrl+C)
./rule30 --bytes=0 | head -c 1073741824 > test.data

# Benchmark throughput
./rule30 --benchmark

# Reproducible output with seed
./rule30 --seed=12345 --bytes=1024 > random.bin
```

### Library Usage

Drop-in replacement for math/rand:

```go
import "github.com/vrypan/rule30rnd/rand"

rng := rand.New(12345)

// Full math/rand interface
rng.Uint64()           // Random uint64
rng.Intn(100)          // Random int in [0, 100)
rng.Float64()          // Random float64 in [0.0, 1.0)
rng.NormFloat64()      // Normal distribution (mean=0, stddev=1)

// Also implements io.Reader
buf := make([]byte, 1024)
rng.Read(buf)
```

**API**: Compatible with `math/rand` - all methods supported (Uint32/64, Int/Intn, Float32/64, NormFloat64, ExpFloat64, Read).

## Performance

Benchmarks on Apple M1 (run `make bench` to reproduce):

### Throughput Comparison

|Algorithm       |     Read32KB |      Read1KB |       Uint64|
|----------------|--------------|--------------|-------------|
|MathRand        |  21468.00 ns |    651.10 ns |      1.78 ns|
|MathRandV2      |  12693.00 ns |    400.00 ns |      3.06 ns|
|**Rule30**      |   **3922.00 ns** |    **127.00 ns** |      **1.28 ns**|
|CryptoRand      |   6895.00 ns |    371.80 ns |     53.04 ns|

**Relative to math/rand:**

|Algorithm       |    Read32KB |     Read1KB |      Uint64|
|----------------|-------------|-------------|------------|
|MathRand        |       1.00x |       1.00x |       1.00x|
|MathRandV2      |       1.69x |       1.63x |       0.58x|
|**Rule30**      |       **5.47x** |       **5.13x** |       **1.39x**|
|CryptoRand      |       3.11x |       1.75x |       0.03x|

### vs /dev/urandom

Rule30 vs kernel CSPRNG (run `./misc/compare-urandom.sh`):

|Size    |       Rule30 |  /dev/urandom |   Speedup|
|--------|--------------|---------------|----------|
|10 MB   |     912 MB/s |      448 MB/s |    2.03x|
|100 MB  |    3364 MB/s |      538 MB/s |    6.25x|
|1 GB    |    3501 MB/s |      551 MB/s |    6.35x|

## Randomness Quality

### Basic Tests (ent)

```bash
./rule30 --bytes=10485760 | ent
```

| Test | Result | Expected | Status |
|------|--------|----------|--------|
| Entropy | 7.999982 bits/byte | 8.0000 | ✓ Perfect (99.9998%) |
| Chi-Square | 253.9 | ~255 (200-310) | ✓ Excellent |
| Arithmetic Mean | 127.4987 | 127.5 | ✓ Perfect |
| Monte Carlo π | 3.140031 | 3.141593 | ✓ 0.05% error |
| Serial Correlation | 0.000171 | 0.0000 | ✓ Uncorrelated |

### Statistical Tests (TestU01)

SmallCrush results - **passes all tests**:

| RNG | Passed | Failed | Pass Rate |
|-----|--------|--------|-----------|
| **Rule30** | **15/15** | 0/15 | **100%** ✓ |

Rule30 passes all 15 SmallCrush tests, demonstrating good statistical quality for a PRNG based on cellular automata.

Run tests yourself:

```bash
cd testu01
make smallcrush           # Test Rule30
make mathrand-smallcrush  # Compare with math/rand
```

## How It Works

**Algorithm**: 1D cellular automaton (Rule 30) on a circular 256-bit strip.

```
Evolution rule: new_bit = left XOR (center OR right)
```

**Implementation**:
- 256-bit state as 4 × 64-bit words
- Processes 64 bits in parallel per word (bitwise operations)
- Fully unrolled loop - all 4 words updated simultaneously
- XOR rotation mixing for enhanced statistical quality
- Each iteration generates 256 bits

**Why it's fast**:
- Bit-level parallelism (64 bits per operation vs 1 bit serially)
- Single-cycle operations (XOR, OR, shifts)
- Zero loop overhead (fully unrolled)
- Minimal memory allocations

**Spatial vs Temporal extraction**: Unlike Wolfram's traditional approach (1 bit per generation from a single column), we extract all 256 bits from the spatial pattern each iteration - **256× more efficient** while maintaining excellent randomness.

## Building & Testing

```bash
# Build all binaries
make all

# Run Go benchmarks
make bench

# Run comparison tools
make compare-run
./misc/compare-urandom.sh

# Run TestU01 statistical tests
cd testu01
make smallcrush
```

## Use Cases

**Good for:**
- Monte Carlo simulations
- Procedural generation (games, graphics)
- High-throughput random sampling
- Reproducible sequences (deterministic from seed)

**Not for:**
- Cryptographic applications (use `crypto/rand`)
- Security-critical randomness

## Background

Rule 30 is a Class III cellular automaton discovered by Stephen Wolfram in 1983. Despite its simple definition, it exhibits chaotic behavior and generates high-quality randomness. Wolfram Research used it as the basis for Mathematica's random number generator.

## License

MIT License - See LICENSE file for details

## References

- Wolfram, Stephen (1983). "Statistical mechanics of cellular automata"
- Wolfram, Stephen (2002). "A New Kind of Science"
- [Rule 30 on Wikipedia](https://en.wikipedia.org/wiki/Rule_30)
