package algorithms

import (
	"os"
	"unsafe"

	"github.com/ikugo-dev/DeFence/logger"
)

func EncryptFile(fileName string, key []byte, algorithm, hash string) []byte {
	metadata := createMetadata(fileName, algorithm, "TigerHash")
	data, err := os.ReadFile(fileName)
	if err != nil {
		logger.Log("Failed to read file: %s", err)
		return nil
	}
	encryptedData := EncryptFileData(data, key, algorithm, hash)
	if encryptedData == nil {
		return nil
	}
	return append(metadata, encryptedData...)
}

func EncryptFileData(data []byte, key []byte, algorithm, hash string) []byte {
	switch algorithm {
	case "Railfence Cipher":
		return EncryptRailfence(data, key)
	case "XXTEA":
		return EncryptXXTEA(data, key)
	case "CBC":
		return EncryptCBC(data, key)
	}
	return nil //ERROR
}

func DecryptFile(fileName string, key []byte) []byte {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		logger.Log("Failed to read encrypted file: %v", err)
		return nil
	}
	metadata := readMetadata(fileData)
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
		return DecryptRailfence(data, key)
	case "XXTEA":
		return DecryptXXTEA(data, key)
	case "CBC":
		return DecryptCBC(data, key)
	}
	return nil //ERROR
}
