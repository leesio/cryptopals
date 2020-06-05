package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"time"

	"github.com/leesio/cryptopals/helpers"
)

func main() {
	start := time.Now()

	fmt.Println(partOne())
	fmt.Println("one done in", time.Now().Sub(start))
	start = time.Now()

	fmt.Println(partTwo())
	fmt.Println("two done in", time.Now().Sub(start))
	start = time.Now()

	fmt.Println(partThree())
	fmt.Println("three done in", time.Now().Sub(start))
	start = time.Now()

	fmt.Println(partFour())
	fmt.Println("four done in", time.Now().Sub(start))
	start = time.Now()

	// fmt.Println(partFive())
	// fmt.Println("five done in", time.Now().Sub(start))
	// start = time.Now()

	// fmt.Println(partSix())
	// fmt.Println("six done in", time.Now().Sub(start))
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
	result, err := helpers.FixedXOR(
		helpers.MustDecodeHex(a),
		helpers.MustDecodeHex(b),
	)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(result)
}

func partThree() string {
	input := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	_, r, _ := helpers.FindSingleByteXOR(helpers.MustDecodeHex(input))
	return string(r)
}
func partFour() string {
	f, err := os.Open("data/part4.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	candidates := make([][]byte, 0)
	for scanner.Scan() {
		candidates = append(candidates, helpers.MustDecodeHex(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	scores := make(chan struct {
		score float64
		match []byte
	})

	var wg sync.WaitGroup
	for _, c := range candidates {
		wg.Add(1)
		go func(c []byte) {
			score, b, _ := helpers.FindSingleByteXOR(c)
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
			if score.score < bestScore {
				bestScore = score.score
				bestMatch = score.match
			}
		}
		result <- bestMatch
	}()
	return string(<-result)
}
func partFive() string {
	s := []byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`)
	key := []byte("ICE")
	return hex.EncodeToString(helpers.DecryptRepeatingKeyXOR(s, key))
}

func partSix() string {
	f, err := os.Open("data/part6.txt")
	if err != nil {
		panic(err)
	}
	// replace with decoder
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	cipher := make([]byte, 10000)
	n, err := base64.StdEncoding.Decode(cipher, raw)
	if err != nil {
		panic(err)
	}
	cipher = cipher[:n]
	keysize := helpers.GetKeySize(cipher)

	b := make([]byte, len(cipher))
	copy(b, cipher)
	blocks := make([][]byte, 0)
	for len(b) > 0 {
		var chunk int
		if keysize > len(b) {
			chunk = len(b)
		} else {
			chunk = keysize
		}
		blocks = append(blocks, b[:chunk])
		b = b[chunk:]
	}

	transposed := make([][]byte, keysize)
	for t := range transposed {
		transposed[t] = make([]byte, 0)
		for _, block := range blocks {
			if len(block) > t {
				transposed[t] = append(transposed[t], block[t])
			}
		}
	}

	key := make([]byte, 0)
	for _, block := range transposed {
		_, _, k := helpers.FindSingleByteXOR(block)
		key = append(key, k)
	}
	return string(helpers.DecryptRepeatingKeyXOR(cipher, key))
}
