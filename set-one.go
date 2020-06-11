package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/leesio/cryptopals/helpers"
)

func main() {
	fmt.Printf(partOne())
	fmt.Println(partTwo())
	fmt.Println(partThree())
	fmt.Println(partFour())
	fmt.Println(partFive())
	fmt.Println(partSix())
	fmt.Println(partSeven())
	fmt.Println(partEight())
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
	return helpers.FindSingleByteXORedStringParallel(candidates)

}
func partFive() string {
	s := []byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`)
	key := []byte("ICE")
	return hex.EncodeToString(helpers.DecryptRepeatingKeyXOR(s, key))
}

func partSix() string {
	cipher, err := helpers.ReadAndDecodeBase64("data/part6.txt")
	if err != nil {
		panic(err)
	}
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
func partSeven() string {
	cipher, err := helpers.ReadAndDecodeBase64("data/part7.txt")
	if err != nil {
		panic(err)
	}
	r := helpers.DecryptAES128ECB(cipher, []byte("YELLOW SUBMARINE"))
	return string(r)
}

func partEight() string {
	f, err := os.Open("data/part8.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	max := 0
	bestMatch := []byte{}
	for scanner.Scan() {
		b := scanner.Bytes()
		if n := helpers.FindRepeatingBlocks(b); n > max {
			max = n
			bestMatch = b
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bestMatch)
}
