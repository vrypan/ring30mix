package main

import (
	cryptorand "crypto/rand"
	"io"
	mathrand "math/rand"
	"testing"
)

// Benchmark TurmiteRNG
func BenchmarkTurmiteRNG_1KB(b *testing.B) {
	benchmarkRNG(b, New(12345), 1024)
}

func BenchmarkTurmiteRNG_10KB(b *testing.B) {
	benchmarkRNG(b, New(12345), 10*1024)
}

func BenchmarkTurmiteRNG_100KB(b *testing.B) {
	benchmarkRNG(b, New(12345), 100*1024)
}

func BenchmarkTurmiteRNG_1MB(b *testing.B) {
	benchmarkRNG(b, New(12345), 1024*1024)
}

// Benchmark crypto/rand
func BenchmarkCryptoRand_1KB(b *testing.B) {
	benchmarkRNG(b, cryptorand.Reader, 1024)
}

func BenchmarkCryptoRand_10KB(b *testing.B) {
	benchmarkRNG(b, cryptorand.Reader, 10*1024)
}

func BenchmarkCryptoRand_100KB(b *testing.B) {
	benchmarkRNG(b, cryptorand.Reader, 100*1024)
}

func BenchmarkCryptoRand_1MB(b *testing.B) {
	benchmarkRNG(b, cryptorand.Reader, 1024*1024)
}

// Benchmark math/rand (via wrapper)
type mathRandReader struct {
	rng *mathrand.Rand
}

func (m *mathRandReader) Read(p []byte) (n int, err error) {
	return m.rng.Read(p)
}

func newMathRandReader(seed int64) io.Reader {
	return &mathRandReader{
		rng: mathrand.New(mathrand.NewSource(seed)),
	}
}

func BenchmarkMathRand_1KB(b *testing.B) {
	benchmarkRNG(b, newMathRandReader(12345), 1024)
}

func BenchmarkMathRand_10KB(b *testing.B) {
	benchmarkRNG(b, newMathRandReader(12345), 10*1024)
}

func BenchmarkMathRand_100KB(b *testing.B) {
	benchmarkRNG(b, newMathRandReader(12345), 100*1024)
}

func BenchmarkMathRand_1MB(b *testing.B) {
	benchmarkRNG(b, newMathRandReader(12345), 1024*1024)
}

// Helper function to benchmark any io.Reader
func benchmarkRNG(b *testing.B, r io.Reader, size int) {
	buf := make([]byte, size)
	b.SetBytes(int64(size))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
