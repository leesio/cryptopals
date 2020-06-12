package settwo

import (
	"crypto/aes"

	"github.com/leesio/cryptopals/helpers"
)

func PartTwo() string {
	b, err := helpers.ReadAndDecodeBase64("./data/part10.txt")
	if err != nil {
		panic(err)
	}
	iv := make([]byte, 16)
	for i := 0; i < 16; i++ {
		iv[i] = 0
	}
	return string(decryptCBC(iv, []byte("YELLOW SUBMARINE"), b))
}

func decryptCBC(iv, key, ciphertext []byte) []byte {
	blockSize := len(key)
	c, _ := aes.NewCipher(key)
	result := make([]byte, len(ciphertext))
	xorBlock := iv
	for i := 0; i < len(ciphertext)/blockSize; i++ {
		cipherBlock := ciphertext[i*blockSize : (i+1)*blockSize]
		resultBlock := result[i*blockSize : (i+1)*blockSize]
		c.Decrypt(resultBlock, cipherBlock)
		for n, b := range resultBlock {
			resultBlock[n] = b ^ xorBlock[n]
		}
		xorBlock = cipherBlock
	}
	return result
}

func pad(k int, b []byte) []byte {
	paddingByte := k - (len(b) % k)
	padding := make([]byte, paddingByte)
	for i := 0; i < paddingByte; i++ {
		padding[i] = byte(paddingByte)
	}
	return append(b, padding...)
}
