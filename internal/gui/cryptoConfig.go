package gui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/internal/algorithms"
	"github.com/ikugo-dev/DeFence/internal/logger"
)

type CryptoConfig struct {
	Algorithm string
	Key       []byte
	Operation string // "Encrypt" | "Decrypt"
}

type CryptoUIComponents struct {
	AlgorithmSelect *widget.Select
	KeyEntry        *widget.Entry
	OperationRadio  *widget.RadioGroup
}

func CreateCryptoUIComponents(includeOperation bool) *CryptoUIComponents {
	components := &CryptoUIComponents{
		AlgorithmSelect: widget.NewSelect([]string{"Railfence", "XXTEA", "CBC"}, nil),
		KeyEntry:        widget.NewEntry(),
	}

	components.AlgorithmSelect.SetSelected("Railfence")
	components.KeyEntry.SetPlaceHolder("Enter encryption key")
	components.KeyEntry.Password = true

	if includeOperation {
		components.OperationRadio = widget.NewRadioGroup([]string{"Encrypt", "Decrypt"}, nil)
		components.OperationRadio.SetSelected("Encrypt")
	}

	return components
}

func (c *CryptoUIComponents) GetConfig() CryptoConfig {
	config := CryptoConfig{
		Algorithm: c.AlgorithmSelect.Selected,
		Key:       algorithms.KeyStringToBigEndianBytes(c.KeyEntry.Text),
	}

	if c.OperationRadio != nil {
		config.Operation = c.OperationRadio.Selected
	}

	return config
}

func (c *CryptoUIComponents) ValidateKey() error {
	if c.KeyEntry.Text == "" {
		return fmt.Errorf("please enter a key")
	}
	return nil
}

func CreateCryptoUISection(components *CryptoUIComponents) *fyne.Container {
	section := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("Select Algorithm:"),
		components.AlgorithmSelect,
		widget.NewLabel("Encryption / Decryption Key:"),
		components.KeyEntry,
	)

	if components.OperationRadio != nil {
		section.Add(widget.NewSeparator())
		section.Add(widget.NewLabel("Select Operation:"))
		section.Add(components.OperationRadio)
	}

	return section
}

func ProcessAndSaveFile(inputPath string, outputPath string, config CryptoConfig) error {
	var data []byte
	var err error

	switch config.Operation {
	case "Encrypt":
		data, err = algorithms.EncryptFile(inputPath, config.Key, config.Algorithm)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
	case "Decrypt":
		data, err = algorithms.DecryptFile(inputPath, config.Key)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}
	default:
		return fmt.Errorf("invalid operation: %s", config.Operation)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	logger.Log("%sed %s -> %s (Algorithm: %s)",
		config.Operation,
		filepath.Base(inputPath),
		filepath.Base(outputPath),
		config.Algorithm)

	return nil
}

func EncryptAndSave(inputPath string, outputPath string, key []byte, algorithm string) error {
	data, err := algorithms.EncryptFile(inputPath, key, algorithm)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	logger.Log("encrypted %s -> %s (Algorithm: %s)",
		filepath.Base(inputPath),
		filepath.Base(outputPath),
		algorithm)

	return nil
}

func DecryptAndSave(data []byte, outputPath string, key []byte) error {
	decrypted, err := algorithms.DecryptFileData(data, key)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	if err := os.WriteFile(outputPath, decrypted, 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	logger.Log("decrypted and saved to %s (%d bytes)", outputPath, len(decrypted))

	return nil
}

func DetermineOutputPath(inputPath, customOutput, operation string) string {
	if operation == "Encrypt" {
		// Encryption: add .enc extension
		if customOutput != "" {
			return customOutput + ".enc"
		}
		return inputPath + ".enc"
	}

	if customOutput != "" {
		return customOutput
	}

	if filepath.Ext(inputPath) == ".enc" {
		return inputPath[:len(inputPath)-4] // Remove .enc
	}
	return inputPath + ".dec"
}
