package algorithms

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/ikugo-dev/DeFence/metadata"
)

func EncryptFile(fileName string, key []byte, algorithm string) ([]byte, error) {
	metadata := metadata.Create(fileName, algorithm, "TigerHash")
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %s", err)
	}
	encryptedData, err := EncryptFileData(data, key, algorithm)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %s", err)
	}
	return append(metadata, encryptedData...), nil
}

func EncryptFileData(data []byte, key []byte, algorithm string) ([]byte, error) {
	switch algorithm {
	case "Railfence":
		return encryptRailfence(data, key)
	case "XXTEA":
		return encryptXXTEA(data, key)
	case "CBC":
		return encryptCBC(data, key)
	}
	return nil, fmt.Errorf("invalid algorithm selection: %s", algorithm)
}

func DecryptFile(fileName string, key []byte) ([]byte, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read encrypted file: %v", err)
	}
	metadata := metadata.Read(fileData)
	encryptedData := fileData[4+unsafe.Sizeof(metadata):]

	decrypted, err := DecryptFileData(encryptedData, key, metadata.EncryptionAlgorithm)
	if err != nil {
		return nil, fmt.Errorf("Decryption failed: %s", err)
	}

	return decrypted, nil
}

func DecryptFileData(data []byte, key []byte, algorithm string) ([]byte, error) {
	switch algorithm {
	case "Railfence Cipher":
		return decryptRailfence(data, key)
	case "XXTEA":
		return decryptXXTEA(data, key)
	case "CBC":
		return decryptCBC(data, key)
	}
	return nil, fmt.Errorf("invalid algorithm selection: %s", algorithm)
}
