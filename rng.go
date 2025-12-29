package main

import (
	"io"
)

// TurmiteRNG implements io.Reader using a turmite cellular automaton
type TurmiteRNG struct {
	turmite         *Turmite
	iterationsSplit int // iterations for first block (second block gets 64-N)
	buffer          []byte
}

// New creates a new TurmiteRNG from a seed
// Seed format (64 bits):
//   bits 0-2:   x position [0-7]
//   bits 3-5:   y position [0-7]
//   bits 6-7:   direction [0-3] (N/E/S/W)
//   bits 8-13:  iteration split [0-64] (first block iterations, second block gets 64-N)
//   bits 14-63: grid initialization seed
func New(seed uint64) *TurmiteRNG {
	x := int(seed & 0x7)
	y := int((seed >> 3) & 0x7)
	dir := int((seed >> 6) & 0x3)
	iterationsSplit := int((seed >> 8) & 0x3F) // 6 bits: 0-63

	// Clamp to valid range [0-64]
	if iterationsSplit > 64 {
		iterationsSplit = 64
	}

	turmite := NewTurmite(x, y, dir)

	// Initialize grid from seed
	turmite.InitGrid(seed)

	return &TurmiteRNG{
		turmite:         turmite,
		iterationsSplit: iterationsSplit,
		buffer:          make([]byte, 0, 32),
	}
}

// Read implements io.Reader interface
// Generates random bytes by running the turmite simulation
func (r *TurmiteRNG) Read(buf []byte) (n int, err error) {
	needed := len(buf)
	copied := 0

	// Use leftover buffer first
	if len(r.buffer) > 0 {
		n := copy(buf, r.buffer)
		r.buffer = r.buffer[n:]
		copied += n
		if copied >= needed {
			return copied, nil
		}
	}

	// Generate more bytes as needed
	for copied < needed {
		output := r.generate32Bytes()

		n := copy(buf[copied:], output)
		copied += n

		// Save leftover bytes
		if n < len(output) {
			r.buffer = append(r.buffer[:0], output[n:]...)
		}
	}

	return copied, nil
}

// generate32Bytes runs the turmite and extracts 32 bytes
// Always runs exactly 64 iterations total, split between two blocks
func (r *TurmiteRNG) generate32Bytes() []byte {
	result := make([]byte, 32)

	// Run turmite for first block iterations
	for i := 0; i < r.iterationsSplit; i++ {
		r.turmite.Step()
	}

	// Extract first 16 bytes from grid state
	grid1 := r.turmite.CopyGrid()
	copy(result[0:16], grid1[:])

	// Run remaining iterations (64 - N)
	remaining := 64 - r.iterationsSplit
	for i := 0; i < remaining; i++ {
		r.turmite.Step()
	}

	// Extract next 16 bytes from new grid state
	grid2 := r.turmite.CopyGrid()
	copy(result[16:32], grid2[:])

	return result
}

// Reader is a global Reader compatible with crypto/rand.Reader interface
// Can be initialized with a seed using InitReader()
var Reader io.Reader

// InitReader initializes the global Reader with a seed
func InitReader(seed uint64) {
	Reader = New(seed)
}
