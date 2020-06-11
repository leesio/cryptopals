package helpers

import (
	"crypto/aes"
	"encoding/hex"
)

const blockSize = 16

func DecryptAES128ECB(data, key []byte) []byte {
	blocks := len(data) / blockSize
	cipher, _ := aes.NewCipher([]byte(key))

	for i := 0; i < blocks; i++ {
		cipher.Decrypt(
			data[i*blockSize:(i+1)*blockSize],
			data[i*blockSize:(i+1)*blockSize],
		)
	}
	return data
}

func FindRepeatingBlocks(b []byte) int {
	blockSize := 16
	m := make(map[string]int)
	for i := 0; i < len(b)/blockSize; i++ {
		m[hex.EncodeToString(b[i*blockSize:(i+1)*blockSize])]++
	}
	max := 0
	for _, val := range m {
		if val > max {
			max = val
		}
	}
	return max
}
