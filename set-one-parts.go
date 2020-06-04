package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

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
	_, r, _ := findSingleByteXOR(third)
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
		score, st, _ := findSingleByteXOR(thing)
		if score < bestScore {
			bestScore = score
			bestStr = st
		}
	}
	return bestStr
}
func partFive() string {
	s := []byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`)
	key := []byte("ICE")
	return hex.EncodeToString(decrypt(s, key))
}

func partSix() int {
	f, err := os.Open("part6.txt")
	if err != nil {
		panic(err)
	}
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
	keysize := getKeySize(cipher)

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
		score, _, k := findSingleByteXOR(hex.EncodeToString(block))
		fmt.Println("got key", k, "with score", score)
		key = append(key, k)
	}
	fmt.Println("key", string(key), len(key), key)
	fmt.Println(string(decrypt(cipher, key)))

	return 0
}
