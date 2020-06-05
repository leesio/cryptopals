package helpers

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

var byteFrequencyMap map[byte]float64

func DecryptRepeatingKeyXOR(s, k []byte) []byte {
	result := make([]byte, len(s))
	for n, b := range s {
		result[n] = b ^ k[n%len(k)]
	}
	return result
}

func GetKeySize(cipher []byte) int {
	keySize := 1
	min := 100.0
	for k := 2; k < 41; k++ {
		total := 0.0
		rounds := 0
		for n := 0; (n + 2) < len(cipher)/k; n++ {
			a, b := cipher[n*k:(n+1)*k], cipher[(n+1)*k:(n+2)*k]
			d := FindHammingDistance(a, b)
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

func FindHammingDistance(a, b []byte) int {
	distance := 0
	for n := range a {
		a := a[n]
		b := b[n]
		// d is the XOR of the 2 things which means it's all the bits that
		// are different.
		// The count of the 1s in the binary representation of the byte is
		// the distance
		xorResult := a ^ b

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
			mask := byte(1 << i)
			if result := mask & xorResult; result == mask {
				distance++
			}
		}
	}
	return distance
}

func FindSingleByteXOR(cipher []byte) (float64, []byte, byte) {
	var bestKey byte
	var bestResult []byte
	bestScore := 1000.0
	for n := 0; n < 256; n++ {
		key := byte(n)
		result := make([]byte, len(cipher))
		for n, b := range cipher {
			result[n] = b ^ key
		}
		score := ScoreBytes(result)
		if score < bestScore {
			bestResult = result
			bestScore = score
			bestKey = key
		}
	}
	return bestScore, bestResult, bestKey
}

func ScoreBytes(b []byte) float64 {
	if byteFrequencyMap == nil {
		byteFrequencyMap = getByteFrequencyMap()
	}

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

func FixedXOR(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return []byte{}, fmt.Errorf("buffers not the same size")
	}

	c := make([]byte, len(a))
	for n, byt := range a {
		c[n] = byt ^ b[n]
	}
	return c, nil

}
func MustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func getByteFrequencyMap() map[byte]float64 {
	f, err := os.Open("data/frequency.txt")
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

func FindSingleByteXORedString(candidates [][]byte) string {
	bestScore := math.MaxFloat64
	var bestMatch []byte
	for _, c := range candidates {
		score, b, _ := FindSingleByteXOR(c)
		if score < bestScore {
			bestScore = score
			bestMatch = b
		}
	}
	return string(bestMatch)
}

func FindSingleByteXORedStringParallel(candidates [][]byte) string {
	scores := make(chan struct {
		score float64
		match []byte
	})
	var wg sync.WaitGroup
	for _, c := range candidates {
		wg.Add(1)
		go func(c []byte) {
			score, b, _ := FindSingleByteXOR(c)
			scores <- struct {
				score float64
				match []byte
			}{score, b}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(scores)
	}()
	result := make(chan []byte)
	go func() {
		bestScore := math.MaxFloat64
		var bestMatch []byte
		for score := range scores {
			if s := score.score; s < bestScore {
				bestScore = s
				bestMatch = score.match
			}
		}
		result <- bestMatch
	}()
	return string(<-result)
}
