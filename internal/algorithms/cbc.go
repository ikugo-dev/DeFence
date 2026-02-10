package algorithms

import (
	"crypto/rand"
	"fmt"
)

const BlockSize = 8

func encryptCBC(data []byte, key []byte) ([]byte, error) {
	iv := make([]byte, BlockSize)
	rand.Read(iv)
	if len(key) != KeySize {
		return nil, fmt.Errorf("Invalid key size")
	}

	data = pkcs7Pad(data)

	out := make([]byte, 0, len(iv)+len(data))
	out = append(out, iv...)

	prev := iv

	for i := 0; i < len(data); i += BlockSize {
		block := data[i : i+BlockSize]
		xored := xorBlocks(block, prev)
		enc, err := encryptXXTEA(xored, key)
		if err != nil {
			return nil, fmt.Errorf("Error while encrypting XXTEA in CBC; %s", err)
		}
		out = append(out, enc...)
		prev = enc
	}

	return out, nil
}

func decryptCBC(data []byte, key []byte) ([]byte, error) {
	if len(data) < BlockSize*2 {
		return nil, fmt.Errorf("Ciphertext too short")
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("Invalid key size")
	}

	iv := data[:BlockSize]
	data = data[BlockSize:]

	prev := iv
	out := make([]byte, 0, len(data))

	for i := 0; i < len(data); i += BlockSize {
		block := data[i : i+BlockSize]
		dec, err := decryptXXTEA(block, key)
		if err != nil {
			return nil, fmt.Errorf("Error while decrypting XXTEA in CBC; %s", err)
		}
		plain := xorBlocks(dec, prev)
		out = append(out, plain...)
		prev = block
	}

	result, err := pkcs7Unpad(out)
	if err != nil {
		return nil, fmt.Errorf("Error while unpadding")
	}
	return result, nil
}

func xorBlocks(a, b []byte) []byte {
	out := make([]byte, BlockSize)
	for i := range BlockSize {
		out[i] = a[i] ^ b[i]
	}
	return out
}
