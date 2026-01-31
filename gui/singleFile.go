package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/algorithms"
	"github.com/ikugo-dev/DeFence/logger"
)

func createSingleFileTab(window fyne.Window, state *AppState) fyne.CanvasObject {
	var selectedFile, outputFile string
	fileLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("Output: Same directory as input file")

	selectFileBtn := widget.NewButton("Select File", func() {
		showFilePicker(window, &selectedFile, fileLabel)
	})

	algorithmSelect := widget.NewSelect([]string{"Railfence Cipher", "XXTEA", "CBC"}, nil)
	algorithmSelect.SetSelected("Railfence Cipher")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter encryption key")
	keyEntry.Password = true

	operationRadio := widget.NewRadioGroup([]string{"Encrypt", "Decrypt"}, nil)
	operationRadio.SetSelected("Encrypt")

	selectOutputBtn := widget.NewButton("Choose Output Location", func() {
		showSaveFilePicker(window, &outputFile, outputLabel)
	})

	processBtn := widget.NewButton("Process File", func() {
		if selectedFile == "" {
			dialog.ShowError(fmt.Errorf("Please select a file first"), window)
			return
		}
		if keyEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Please enter a key"), window)
			return
		}

		operation := operationRadio.Selected
		algorithm := algorithmSelect.Selected
		key := algorithms.KeyStringToBigEndianBytes(keyEntry.Text)
		output := determineOutputPath(selectedFile, outputFile, operation)

		progress := dialog.NewProgressInfinite("Processing", fmt.Sprintf("%sing file with %s...", operation, algorithm), window)
		progress.Show()

		go func() {

			switch operation {
			case "Encrypt":
				_, err := algorithms.EncryptFile(selectedFile, key, algorithm)
				if err != nil {
					logger.LogWithDialog(window, "Error", "Error while encrypting: %s", err)
				}
			case "Decrypt":
				_, err := algorithms.DecryptFile(selectedFile, key)
				if err != nil {
					logger.LogWithDialog(window, "Error", "Error while decrypting: %s", err)
				}
			}

			fyne.Do(func() {
				progress.Hide()
				logger.LogWithDialog(window, "Success", "%sed file: %s → %s successfully. (Algorithm: %s)", operation, filepath.Base(selectedFile), filepath.Base(output), algorithm)
			})
		}()
	})

	return container.NewVBox(
		widget.NewLabel("Single File Encryption/Decryption:"),
		widget.NewSeparator(),
		selectFileBtn,
		fileLabel,
		widget.NewSeparator(),
		widget.NewLabel("Select Algorithm:"),
		algorithmSelect,
		widget.NewLabel("Encryption / Decryption Key:"),
		keyEntry,
		widget.NewSeparator(),
		widget.NewLabel("Select Operation:"),
		operationRadio,
		widget.NewSeparator(),
		selectOutputBtn,
		outputLabel,
		widget.NewSeparator(),
		layout.NewSpacer(),
		processBtn,
	)
}

func determineOutputPath(selectedFile, outputFile, operation string) string {
	if outputFile != "" {
		if operation == "Encrypt" {
			return outputFile + ".enc"
		}
		return outputFile + ".dec"
	}

	if operation == "Encrypt" {
		return selectedFile + ".enc"
	}
	return selectedFile + ".dec"
}
