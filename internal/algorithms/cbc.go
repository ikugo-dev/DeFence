package algorithms

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

const BlockSize = 16

func encryptCBC(data []byte, key []byte) ([]byte, error) {
	iv := make([]byte, BlockSize)
	rand.Read(iv)
	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size")
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
			return nil, fmt.Errorf("error while encrypting XXTEA in CBC; %s", err)
		}
		out = append(out, enc...)
		prev = enc
	}

	return out, nil
}

func decryptCBC(data []byte, key []byte) ([]byte, error) {
	if len(data) < BlockSize*2 {
		return nil, fmt.Errorf("ciphertext too short")
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size")
	}

	iv := data[:BlockSize]
	data = data[BlockSize:]

	prev := iv
	out := make([]byte, 0, len(data))

	for i := 0; i < len(data); i += BlockSize {
		block := data[i : i+BlockSize]
		dec, err := decryptXXTEA(block, key)
		if err != nil {
			return nil, fmt.Errorf("error while decrypting XXTEA in CBC; %s", err)
		}
		plain := xorBlocks(dec, prev)
		out = append(out, plain...)
		prev = block
	}

	result, err := pkcs7Unpad(out)
	if err != nil {
		return nil, fmt.Errorf("invalid padding")
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

func pkcs7Pad(data []byte) []byte {
	padLen := BlockSize - (len(data) % BlockSize)
	if padLen == 0 {
		padLen = BlockSize
	}
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, pad...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%BlockSize != 0 {
		return nil, fmt.Errorf("PKCS7 Unpad: data is invalid size")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > BlockSize {
		return nil, fmt.Errorf("PKCS7 Unpad: lenght of padding is invalid size")
	}
	for i := range padLen {
		if data[len(data)-1-i] != byte(padLen) {
			return nil, fmt.Errorf("PKCS7 Unpad: idek")
		}
	}
	return data[:len(data)-padLen], nil
}
