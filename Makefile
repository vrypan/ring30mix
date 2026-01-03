# ring30mix RNG Makefile

# Binary names
RING30MIX_BIN = ring30mix
COMPARE_READ_BIN = misc/compare-read
COMPARE_UINT64_BIN = misc/compare-uint64

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

# Build flags
LDFLAGS = -s -w
BUILD_FLAGS = -ldflags "$(LDFLAGS)"

# Source files for dependency tracking
RING30MIX_SOURCES = main.go cmd/root.go cmd/raw.go cmd/ascii.go cmd/version.go rand/ring30mix.go
COMPARE_READ_SOURCES = misc/compare-read.go rand/ring30mix.go
COMPARE_UINT64_SOURCES = misc/compare-uint64.go rand/ring30mix.go

.PHONY: all compare clean fmt help compare-run test-entropy smoke deps bench

# Default target
all: $(RING30MIX_BIN) compare

# Build the ring30mix CLI tool
$(RING30MIX_BIN): $(RING30MIX_SOURCES)
	@echo "Building $(RING30MIX_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(RING30MIX_BIN) main.go
	@echo "✓ Built $(RING30MIX_BIN)"

# Build both comparison tools
compare: $(COMPARE_READ_BIN) $(COMPARE_UINT64_BIN)

# Build the Read() comparison tool
$(COMPARE_READ_BIN): $(COMPARE_READ_SOURCES)
	@echo "Building $(COMPARE_READ_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_READ_BIN) misc/compare-read.go
	@echo "✓ Built $(COMPARE_READ_BIN)"

# Build the Uint64() comparison tool
$(COMPARE_UINT64_BIN): $(COMPARE_UINT64_SOURCES)
	@echo "Building $(COMPARE_UINT64_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_UINT64_BIN) misc/compare-uint64.go
	@echo "✓ Built $(COMPARE_UINT64_BIN)"

# Run comparison benchmarks
compare-run: compare
	@echo "Running Read() benchmark..."
	./$(COMPARE_READ_BIN)
	@echo ""
	@echo "Running Uint64() benchmark..."
	./$(COMPARE_UINT64_BIN)

# Run go test benchmarks with table output
bench:
	@./misc/bench-table.sh

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(RING30MIX_BIN)
	rm -f $(COMPARE_READ_BIN)
	rm -f $(COMPARE_UINT64_BIN)
	rm -f misc/stdlib-rng
	rm -f misc/visualize-ring30mix
	rm -f *.prof
	rm -f *.test
	rm -f *.bin
	rm -f *.dat
	@echo "✓ Cleaned"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Test randomness with ent - compares all three RNGs
test-entropy: ring30mix
	@./misc/test-entropy.sh

# Quick smoke test
smoke: ring30mix
	@echo "Running smoke test..."
	@./$(RING30MIX_BIN) raw --seed=12345 --bytes=1024 > /dev/null
	@echo "✓ Smoke test passed"

# Show help
help:
	@echo "ring30mix RNG Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            Build all binaries (default)"
	@echo "  ring30mix      Build ring30mix CLI tool"
	@echo "  compare        Build comparison tools (read + uint64)"
	@echo "  compare-read   Build compare-read tool (MB/s benchmark)"
	@echo "  compare-uint64 Build compare-uint64 tool (ns/call benchmark)"
	@echo "  compare-run    Run both comparison benchmarks"
	@echo "  bench          Run go test benchmarks (table format)"
	@echo "  fmt            Format code with gofmt"
	@echo "  clean          Remove build artifacts"
	@echo "  deps           Download and tidy dependencies"
	@echo "  test-entropy   Test randomness with ent tool"
	@echo "  smoke          Quick smoke test"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make ring30mix"
	@echo "  make compare-run"
	@echo "  make test-entropy"
	@echo "  make clean ring30mix"
