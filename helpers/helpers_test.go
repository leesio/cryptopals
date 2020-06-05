package helpers

import (
	"bufio"
	"os"
	"testing"
)

func TestGetHammingDistance(t *testing.T) {
	actual := FindHammingDistance([]byte("this is a test"), []byte("wokka wokka!!!"))
	// https://cryptopals.com/sets/1/challenges/6
	exp := 37
	if actual != exp {
		t.Errorf("Got %d as distance between reference strings, expected: %d", actual, exp)
	}
}

func BenchmarkFindSingleByteXORedString(b *testing.B) {
	f, err := os.Open("data/part4.txt")
	if err != nil {
		b.Error(err)
	}
	scanner := bufio.NewScanner(f)
	candidates := make([][]byte, 0)
	for scanner.Scan() {
		candidates = append(candidates, MustDecodeHex(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		FindSingleByteXORedString(candidates)
	}
}
func BenchmarkPartFourParallel(b *testing.B) {
	f, err := os.Open("data/part4.txt")
	if err != nil {
		b.Error(err)
	}
	scanner := bufio.NewScanner(f)
	candidates := make([][]byte, 0)
	for scanner.Scan() {
		candidates = append(candidates, MustDecodeHex(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		FindSingleByteXORedStringParallel(candidates)
	}
}
