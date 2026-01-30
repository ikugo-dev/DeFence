package gui

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
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
		output := determineOutputPath(selectedFile, outputFile, operation)

		progress := dialog.NewProgressInfinite("Processing",
			fmt.Sprintf("%sing file with %s...", operation, algorithm), window)
		progress.Show()

		go func() {
			key := keyStringToBigEndianBytes(keyEntry.Text)

			if operation == "Encrypt" {
				algorithms.EncryptFile(selectedFile, key, algorithm)
			} else {
				algorithms.DecryptFile(selectedFile, key)
			}

			fyne.Do(func() {
				progress.Hide()
				logger.Log("%sed file: %s → %s (Algorithm: %s)", operation,
					filepath.Base(selectedFile), filepath.Base(output), algorithm)
				dialog.ShowInformation("Success",
					fmt.Sprintf("File %sed successfully!\nOutput: %s", operation, output), window)
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
		widget.NewSeparator(),
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

func showFilePicker(window fyne.Window, selectedFile *string, fileLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	fd := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if err != nil || file == nil {
			return
		}
		*selectedFile = file.URI().Path()
		fileLabel.SetText("Selected: " + filepath.Base(*selectedFile))
		file.Close()
	}, window)

	if listableURI != nil {
		fd.SetLocation(listableURI)
	}
	fd.Show()
}

func showSaveFilePicker(window fyne.Window, outputFile *string, outputLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	fd := dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) {
		if err != nil || file == nil {
			return
		}
		*outputFile = file.URI().Path()
		outputLabel.SetText("Output: " + filepath.Base(*outputFile))
		file.Close()
	}, window)

	if listableURI != nil {
		fd.SetLocation(listableURI)
	}
	fd.Show()
}

func keyStringToBigEndianBytes(s string) []byte {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 1 {
		return nil
	}

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return buf
}
