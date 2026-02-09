package algorithms

import (
	"bytes"
	"fmt"
)

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
