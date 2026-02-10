package gui

import (
	"os"
	"path/filepath"

	"github.com/ikugo-dev/DeFence/internal/metadata"
)

func determineOutputPath(inputPath, customOutput, operation string) string {
	if operation == "Encrypt" {
		return encryptOutputPath(inputPath, customOutput)
	}

	fileData, err := os.ReadFile(inputPath)
	if err != nil {
		return inputPath + ".dec"
	}
	return decryptOutputPath(fileData, customOutput)
}

func encryptOutputPath(inputPath, customOutput string) string {
	os.Remove(inputPath)
	if customOutput == "" {
		return inputPath + ".enc"
	}
	if filepath.Ext(customOutput) == ".enc" {
		return customOutput
	}
	return customOutput + ".enc"
}

func decryptOutputPath(fileData []byte, customOutput string) string {
	if customOutput != "" {
		return customOutput
	}

	metadata := metadata.Read(fileData)
	return metadata.FileName
}
