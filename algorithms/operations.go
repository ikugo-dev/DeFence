package algorithms

import (
	"os"
	"unsafe"

	"github.com/ikugo-dev/DeFence/logger"
	"github.com/ikugo-dev/DeFence/metadata"
)

func EncryptFile(fileName string, key []byte, algorithm string) []byte {
	metadata := metadata.Create(fileName, algorithm, "TigerHash")
	data, err := os.ReadFile(fileName)
	if err != nil {
		logger.Log("Failed to read file: %s", err)
		return nil
	}
	encryptedData := EncryptFileData(data, key, algorithm)
	if encryptedData == nil {
		return nil
	}
	return append(metadata, encryptedData...)
}

func EncryptFileData(data []byte, key []byte, algorithm string) []byte {
	switch algorithm {
	case "Railfence":
		return encryptRailfence(data, key)
	case "XXTEA":
		return encryptXXTEA(data, key)
	case "CBC":
		return encryptCBC(data, key)
	}
	return nil //ERROR
}

func DecryptFile(fileName string, key []byte) []byte {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		logger.Log("Failed to read encrypted file: %v", err)
		return nil
	}
	metadata := metadata.Read(fileData)
	encryptedData := fileData[4+unsafe.Sizeof(metadata):]

	decrypted := DecryptFileData(encryptedData, key, metadata.EncryptionAlgorithm)
	if decrypted == nil {
		logger.Log("Decryption failed")
		return nil
	}

	return decrypted
}

func DecryptFileData(data []byte, key []byte, algorithm string) []byte {
	switch algorithm {
	case "Railfence Cipher":
		return decryptRailfence(data, key)
	case "XXTEA":
		return decryptXXTEA(data, key)
	case "CBC":
		return decryptCBC(data, key)
	}
	return nil //ERROR
}
