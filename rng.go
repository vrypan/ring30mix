package main

import (
	"io"
)

// TurmiteRNG implements io.Reader using a turmite cellular automaton
type TurmiteRNG struct {
	turmite    *Turmite
	iterations int
	buffer     []byte
}

// New creates a new TurmiteRNG from a seed
// Seed format (64 bits):
//   bits 0-2:   x position [0-7]
//   bits 3-5:   y position [0-7]
//   bits 6-7:   direction [0-3] (N/E/S/W)
//   bits 8-31:  iterations per 32-byte block
//   bits 32-63: unused (for future extensions)
func New(seed uint64) *TurmiteRNG {
	x := int(seed & 0x7)
	y := int((seed >> 3) & 0x7)
	dir := int((seed >> 6) & 0x3)
	iterations := int((seed >> 8) & 0xFFFFFF)

	// Default to 1000 iterations if not specified
	if iterations == 0 {
		iterations = 1000
	}

	return &TurmiteRNG{
		turmite:    NewTurmite(x, y, dir),
		iterations: iterations,
		buffer:     make([]byte, 0, 32),
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
// Uses Option A: run twice and concatenate grid states
func (r *TurmiteRNG) generate32Bytes() []byte {
	result := make([]byte, 32)

	// Run turmite for N iterations
	for i := 0; i < r.iterations; i++ {
		r.turmite.Step()
	}

	// Extract first 16 bytes from grid state
	grid1 := r.turmite.CopyGrid()
	copy(result[0:16], grid1[:])

	// Run another N iterations
	for i := 0; i < r.iterations; i++ {
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
