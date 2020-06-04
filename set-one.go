package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strings"
)

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

var ()

func main() {
	fmt.Println(partOne())
	fmt.Println(partTwo())
	fmt.Println(partThree())
	fmt.Println(partFour())

}

func partOne() string {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	b, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func partTwo() string {
	a := "1c0111001f010100061a024b53535009181c"
	b := "686974207468652062756c6c277320657965"
	result, err := fixedXOR(decodeHex(a), decodeHex(b))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(result)
}

func partThree() string {
	third := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	_, r := findSingleByteXOR(third)
	return r
}
func partFour() string {
	f, err := os.Open("part4.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	things := make([]string, 0)
	for scanner.Scan() {
		things = append(things, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	bestScore := 1000.0
	bestStr := ""
	for _, thing := range things {
		score, st := findSingleByteXOR(thing)
		fmt.Println("best score", score, st)
		if score < bestScore {
			bestScore = score
			bestStr = st
		}
	}
	return bestStr
}
func findSingleByteXOR(s string) (float64, string) {
	b := decodeHex(s)
	bestScore := 1000.0
	var sentence string
	for n := 0; n < len(alphabet)*2; n++ {
		result := make([]byte, len(b))
		idx := n % len(alphabet)
		var letter byte
		if n >= len(alphabet) {
			letter = alphabet[idx][0]
		} else {
			letter = strings.ToUpper(alphabet[idx])[0]
		}
		for n, byt := range b {
			result[n] = byt ^ byte(letter)
		}
		score := scoreString(string(result))
		if score < bestScore {
			sentence = string(result)
			bestScore = score
		}
	}
	return bestScore, sentence
}
func scoreString(s string) float64 {
	m := map[string]float64{
		"a": 8.497,
		"b": 1.492,
		"c": 2.202,
		"d": 4.253,
		"e": 11.162,
		"f": 2.228,
		"g": 2.015,
		"h": 6.094,
		"i": 7.546,
		"j": 0.153,
		"k": 1.292,
		"l": 4.025,
		"m": 2.406,
		"n": 6.749,
		"o": 7.507,
		"p": 1.929,
		"q": 0.095,
		"r": 7.587,
		"s": 6.327,
		"t": 9.356,
		"u": 2.758,
		"v": 0.978,
		"w": 2.560,
		"x": 0.150,
		"y": 1.994,
		"z": 0.077,
	}
	expected := make(map[string]float64)
	actual := make(map[string]float64)
	for char, pct := range m {
		expected[char] = pct * (float64(len(s)) / 100)
	}
	for _, char := range strings.Split(s, "") {
		actual[char]++
	}
	total := 0.0
	for char, count := range actual {
		if char == " " {
			continue
		}
		e := expected[char]
		total = total + math.Abs(e-count)
	}
	return total / float64(len(s))
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
