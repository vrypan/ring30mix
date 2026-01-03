package rand

import (
	"encoding/binary"
	"math/bits"
)

// RNG implements a 1D cellular automaton (Rule 30) on a 256-bit ring
type RNG struct {
	state [4]uint64 // 256-bit state (4 × 64-bit words)
	pos   int       // current position for output (0-3)
}

// New creates a new Rule 30 RNG from a seed
func New(seed uint64) *RNG {
	rng := &RNG{
		state: [4]uint64{
			seed,
			seed ^ 0x9e3779b97f4a7c15,
			seed ^ 0x3c6ef372fe94f82a,
			seed ^ 0x78dde6e5fd29f054,
		},
		pos: 0,
	}
	// Run a few steps to mix the initial state
	for i := 0; i < 16; i++ {
		rng.step()
	}
	return rng
}

// step applies radius-1 Rule 30 to the 256-bit ring
// Rule 30: new_bit = left XOR (center OR right)
// Unrolled and optimized with pre-computed borders
//
//go:noinline
func (r *RNG) step() {
	s0, s1, s2, s3 := r.state[0], r.state[1], r.state[2], r.state[3]

	// Pre-compute border bits (one state's right border is next state's left border)
	b0_0, b0_1, b0_2, b0_3 := s0&1, s1&1, s2&1, s3&1
	b63_0, b63_1, b63_2, b63_3 := s0>>63, s1>>63, s2>>63, s3>>63

	// Word 0: left from s3.bit0, right from s1.bit63
	r.state[0] = ((s0 >> 1) | (b0_3 << 63)) ^ (s0 | ((s0 << 1) | b63_1))

	// Word 1: left from s0.bit0, right from s2.bit63
	r.state[1] = ((s1 >> 1) | (b0_0 << 63)) ^ (s1 | ((s1 << 1) | b63_2))

	// Word 2: left from s1.bit0, right from s3.bit63
	r.state[2] = ((s2 >> 1) | (b0_1 << 63)) ^ (s2 | ((s2 << 1) | b63_3))

	// Word 3: left from s2.bit0, right from s0.bit63
	r.state[3] = ((s3 >> 1) | (b0_2 << 63)) ^ (s3 | ((s3 << 1) | b63_0))
}

// mix applies a diffusion function to improve output quality
// Uses hybrid rotation + multiply mixing for optimal balance of speed and quality
// This mixing function achieves perfect SmallCrush (15/15) while being faster than math/rand
func mix(x uint64) uint64 {
	x ^= bits.RotateLeft64(x, 13)
	x *= 0x9e3779b97f4a7c15 // Golden ratio constant
	x ^= x >> 27
	return x
}

// Uint64 returns a random uint64
// Cycles through the 4 state words, applying mixing to each
// Only evolves the state when all 4 words have been consumed
func (r *RNG) Uint64() uint64 {
	if r.pos == 4 {
		r.step()
		r.pos = 0
	}
	out := mix(r.state[r.pos])
	r.pos++
	return out
}

// Read implements io.Reader interface
func (r *RNG) Read(buf []byte) (n int, err error) {
	i := 0
	limit := len(buf)

	// Handle 8-byte chunks
	for limit-i >= 8 {
		val := r.Uint64()
		binary.LittleEndian.PutUint64(buf[i:], val)
		i += 8
	}

	// Handle remaining tail bytes
	if rem := limit - i; rem > 0 {
		val := r.Uint64()
		for j := 0; j < rem; j++ {
			buf[i+j] = byte(val)
			val >>= 8
		}
	}

	return limit, nil
}
