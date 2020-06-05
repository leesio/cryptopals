package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/leesio/cryptopals/helpers"
)

func BenchmarkPartFour(b *testing.B) {
	f, err := os.Open("data/part4.txt")
	if err != nil {
		b.Error(err)
	}
	scanner := bufio.NewScanner(f)
	candidates := make([][]byte, 0)
	for scanner.Scan() {
		candidates = append(candidates, helpers.MustDecodeHex(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		partFourComp(candidates)
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
		candidates = append(candidates, helpers.MustDecodeHex(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		b.Error(err)
	}
	for n := 0; n < b.N; n++ {
		partFourCompParallel(candidates)
	}
}
