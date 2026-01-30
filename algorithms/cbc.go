package algorithms

import (
	"bytes"
	"crypto/rand"
	"log"
)

const BlockSize = 16

func encryptCBC(data []byte, key []byte) []byte {
	iv := make([]byte, BlockSize)
	rand.Read(iv)
	if len(key) != KeySize {
		log.Fatal("invalid key size")
	}

	data = pkcs7Pad(data)

	out := make([]byte, 0, len(iv)+len(data))
	out = append(out, iv...)

	prev := iv

	for i := 0; i < len(data); i += BlockSize {
		block := data[i : i+BlockSize]
		xored := xorBlocks(block, prev)
		enc := encryptXXTEA(xored, key)
		out = append(out, enc...)
		prev = enc
	}

	return out
}

func decryptCBC(data []byte, key []byte) []byte {
	if len(data) < BlockSize*2 {
		log.Fatal("ciphertext too short")
	}
	if len(key) != KeySize {
		log.Fatal("invalid key size")
	}

	iv := data[:BlockSize]
	data = data[BlockSize:]

	prev := iv
	out := make([]byte, 0, len(data))

	for i := 0; i < len(data); i += BlockSize {
		block := data[i : i+BlockSize]
		dec := decryptXXTEA(block, key)
		plain := xorBlocks(dec, prev)
		out = append(out, plain...)
		prev = block
	}

	out = pkcs7Unpad(out)
	if out == nil {
		log.Fatal("invalid padding")
	}
	return out
}

func xorBlocks(a, b []byte) []byte {
	out := make([]byte, BlockSize)
	for i := range BlockSize {
		out[i] = a[i] ^ b[i]
	}
	return out
}

func pkcs7Pad(data []byte) []byte {
	padLen := BlockSize - (len(data) % BlockSize)
	if padLen == 0 {
		padLen = BlockSize
	}
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, pad...)
}

func pkcs7Unpad(data []byte) []byte {
	if len(data) == 0 || len(data)%BlockSize != 0 {
		return nil
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > BlockSize {
		return nil
	}
	for i := range padLen {
		if data[len(data)-1-i] != byte(padLen) {
			return nil
		}
	}
	return data[:len(data)-padLen]
}
