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

func (c *CryptoUIComponents) GetKey() []byte {
	algorithm := c.AlgorithmSelect.Selected
	keyText := c.KeyEntry.Text

	switch algorithm {
	case "Railfence":
		return algorithms.KeyStringTo4Bytes(keyText)
	case "XXTEA", "CBC":
		key, _ := algorithms.KeyHexStringTo16Bytes(keyText)
		//TODO handle errors
		return key
	default:
		return algorithms.KeyStringTo4Bytes(keyText)
	}
}

func (c *CryptoUIComponents) GetAlgorithm() string {
	return c.AlgorithmSelect.Selected
}

func (c *CryptoUIComponents) GetOperation() string {
	if c.OperationRadio != nil {
		return c.OperationRadio.Selected
	}
	return ""
}

func (c *CryptoUIComponents) ValidateKey() error {
	if c.KeyEntry.Text == "" {
		return fmt.Errorf("Please enter a key")
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

func ProcessAndSaveFile(inputPath, outputPath, operation string, key []byte, algorithm string) error {
	var data []byte
	var err error

	switch operation {
	case "Encrypt":
		data, err = algorithms.EncryptFile(inputPath, key, algorithm)
		if err != nil {
			return fmt.Errorf("Encryption failed: %w", err)
		}
	case "Decrypt":
		data, err = algorithms.DecryptFile(inputPath, key)
		if err != nil {
			return fmt.Errorf("Decryption failed: %w", err)
		}
	default:
		return fmt.Errorf("Invalid operation: %s", operation)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("Failed to write output file: %w", err)
	}

	logger.Log("%sed %s -> %s (Algorithm: %s)",
		operation,
		filepath.Base(inputPath),
		filepath.Base(outputPath),
		algorithm)

	return nil
}

func ProcessAndSaveData(data []byte, outputPath, operation string, key []byte, algorithm string) error {
	var processedData []byte
	var err error

	switch operation {
	case "Encrypt":
		processedData, err = algorithms.EncryptRawData(data, key, algorithm)
		if err != nil {
			return fmt.Errorf("Encryption failed: %w", err)
		}
	case "Decrypt":
		processedData, err = algorithms.DecryptFileData(data, key)
		if err != nil {
			return fmt.Errorf("Decryption failed: %w", err)
		}
	default:
		return fmt.Errorf("Invalid operation: %s", operation)
	}

	if err := os.WriteFile(outputPath, processedData, 0644); err != nil {
		return fmt.Errorf("Failed to write output file: %w", err)
	}

	logger.Log("%sed data and saved to %s (%d bytes)",
		operation, outputPath, len(processedData))

	return nil
}
