package main

import (
	"github.com/ikugo-dev/DeFence/gui"
)

func main() {
	gui.Start()
}

// TODO: Implement encryption/decryption functions
func EncryptRailfence(data []byte, key int) ([]byte, error) { return data, nil }
func DecryptRailfence(data []byte, key int) ([]byte, error) { return data, nil }
func EncryptXXTEA(data []byte, key []byte) ([]byte, error)  { return data, nil }
func DecryptXXTEA(data []byte, key []byte) ([]byte, error)  { return data, nil }
func EncryptCBC(data []byte, key []byte) ([]byte, error)    { return data, nil }
func DecryptCBC(data []byte, key []byte) ([]byte, error)    { return data, nil }
