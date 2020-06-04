package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var byteFrequencyMap map[byte]float64

func main() {
	byteFrequencyMap = getByteFrequencyMap()
	// fmt.Println("one")
	// fmt.Println(partOne())
	// fmt.Println("two")
	// fmt.Println(partTwo())
	// fmt.Println("three")
	// fmt.Println(partThree())
	// fmt.Println("four")
	// fmt.Println(partFour())
	// fmt.Println("five")
	// fmt.Println(partFive())
	// fmt.Println("six")
	fmt.Println(partSix())
}

func decrypt(s, k []byte) []byte {
	result := make([]byte, len(s))
	for n, b := range s {
		result[n] = b ^ k[n%len(k)]
	}
	fmt.Println("decrypted result", len(result), len(s))
	return result
}

func getKeySize(cipher []byte) int {
	keySize := 1
	min := 100.0
	for k := 2; k < 41; k++ {
		total := 0.0
		rounds := 0
		for n := 0; (n + 2) < len(cipher)/k; n++ {
			a, b := cipher[n*k:(n+1)*k], cipher[(n+1)*k:(n+2)*k]
			d := findDistance(string(a), string(b))
			fmt.Println("k", k, "n", n, "d", d, "d", d/k)
			total = total + float64(d)/float64(k)
			rounds++
		}
		mean := total / float64(rounds)
		if mean < min {
			min = mean
			keySize = k
		}
	}
	return keySize
}

func findDistance(a, b string) int {
	sum := 0
	for n := range a {
		a := byte(a[n])
		b := byte(b[n])
		// d is the XOR of the 2 things which means it's all the bits that
		// are different.
		// The count of the 1s in the binary representation of the byte is
		// the distance
		d := a ^ b

		// Starting with the right most bit, calculate the bitwise AND of the
		// byte and 1, if this equals 1, it means the right most bit is 1.
		// Repeat for the each subsequent bit, bitshifting the 1 a single place
		// left each time.
		// e.g. for 1001010
		// 1001010 & 0000001 = 0000000 so the 1st bit is *not* a 1
		// bitshift the 1, 1 place left 1<<1 is 00000010
		// 1001010 & 0000010 = 0000010 == 2 ^ 1 so the 1st bit *is* a 1)
		// repeat for all 8 bits
		for i := 0; i < 8; i++ {
			bit := byte(1 << i)
			if result := bit & d; result == bit {
				sum = sum + 1
			}
		}
	}
	return sum
}

func findSingleByteXOR(s string) (float64, string, byte) {
	b := decodeHex(s)
	bestScore := 1000.0
	var sentence string
	var key byte
	for n := 0; n < 256; n++ {
		x := byte(n)
		result := make([]byte, len(b))
		for n, byt := range b {
			result[n] = byt ^ x
		}
		score := scoreBytes(result)
		if score < bestScore {
			sentence = string(result)
			bestScore = score
			key = x
		}
	}

	return bestScore, sentence, key
}
func getByteFrequencyMap() map[byte]float64 {
	f, err := os.Open("frequency.txt")
	if err != nil {
		panic(err)
	}
	m := make(map[byte]float64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		b, f := parts[0], parts[1]
		by, err := strconv.Atoi(b)
		if err != nil {
			panic(err)
		}
		fre, err := strconv.ParseFloat(f, 64)
		if err != nil {
			panic(err)
		}
		m[byte(by)] = fre
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return m
}
func scoreBytes(b []byte) float64 {
	expected := make(map[byte]float64)
	actual := make(map[byte]float64)
	for char, pct := range byteFrequencyMap {
		expected[char] = pct * (float64(len(b)) / 100)
	}
	for _, char := range b {
		actual[char]++
	}
	total := 0.0
	for char, count := range actual {
		total = total + math.Abs(expected[char]-count)
	}
	return math.Sqrt(total)
}

func fixedXOR(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return []byte{}, fmt.Errorf("buffers not the same size")
	}

	c := make([]byte, len(a))
	for n, byt := range a {
		c[n] = byt ^ b[n]
	}
	return c, nil

}
func decodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
