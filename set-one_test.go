package main

import (
	"fmt"
	"testing"
)

func TestScoring(t *testing.T) {
	problem := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	a, b, c := findSingleByteXOR(problem)
	fmt.Println(a, b, c)
}
func TestDistance(t *testing.T) {
	a := "this is a test"
	b := "wokka wokka!!!"
	fmt.Println(findDistance(a, b))
	fmt.Println(findDistance(b, a))
}
