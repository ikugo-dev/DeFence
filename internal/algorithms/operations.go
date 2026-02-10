package algorithms

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/ikugo-dev/DeFence/internal/metadata"
)

func EncryptFile(fileName string, key []byte, algorithm string) ([]byte, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %s", err)
	}
	return EncryptFileData(fileName, fileData, key, algorithm)
}

func EncryptFileData(fileName string, fileData, key []byte, algorithm string) ([]byte, error) {
	encryptedData, err := EncryptRawData(fileData, key, algorithm)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %s", err)
	}
	hashedData := TigerHash(encryptedData)
	metadata := metadata.Create(fileName, algorithm, "TigerHash", hashedData)
	return append(metadata, encryptedData...), nil
}

func EncryptRawData(data []byte, key []byte, algorithm string) ([]byte, error) {
	switch algorithm {
	case "Railfence":
		return encryptRailfence(data, key)
	case "XXTEA":
		return encryptXXTEA(data, key)
	case "CBC":
		return encryptCBC(data, key)
	}
	return nil, fmt.Errorf("Invalid algorithm selection: %s", algorithm)
}

func DecryptFile(fileName string, key []byte) ([]byte, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read encrypted file: %v", err)
	}
	return DecryptFileData(fileData, key)
}

func DecryptFileData(data []byte, key []byte) ([]byte, error) {
	metadataLen := binary.BigEndian.Uint32(data[:4])
	metadata := metadata.Read(data)
	encryptedData := data[4+metadataLen:]

	hashedData := TigerHash(encryptedData)
	if hashedData != metadata.HashingResult {
		return nil, fmt.Errorf("Incorrect data recieved: Hashes differ")
	}

	decrypted, err := DecryptRawData(encryptedData, key, metadata.EncryptionAlgorithm)
	if err != nil {
		return nil, fmt.Errorf("Decryption failed: %s", err)
	}
	return decrypted, nil
}

func DecryptRawData(data []byte, key []byte, algorithm string) ([]byte, error) {
	switch algorithm {
	case "Railfence":
		return decryptRailfence(data, key)
	case "XXTEA":
		return decryptXXTEA(data, key)
	case "CBC":
		return decryptCBC(data, key)
	}
	return nil, fmt.Errorf("Invalid algorithm selection: %s", algorithm)
}
