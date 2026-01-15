package gui

import (
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/logger"
)

func createSingleFileTab(window fyne.Window, state *AppState) fyne.CanvasObject {
	fileLabel := widget.NewLabel("No file selected")
	var selectedFile string

	algorithmSelect := widget.NewSelect([]string{"Railfence Cipher", "XXTEA", "CBC"}, nil)
	algorithmSelect.SetSelected("Railfence Cipher")

	operationRadio := widget.NewRadioGroup([]string{"Encrypt", "Decrypt"}, nil)
	operationRadio.SetSelected("Encrypt")

	outputLabel := widget.NewLabel("Output: Same directory as input file")
	var outputFile string

	fileSection := createFileSelectionSection(window, &selectedFile, fileLabel)
	algorithmSection := createAlgorithmSection(algorithmSelect)
	operationSection := createOperationSection(operationRadio)
	outputSection := createOutputSection(window, &outputFile, outputLabel)
	processBtn := createProcessButton(window, state, &selectedFile, &outputFile, algorithmSelect, operationRadio)

	return container.NewVBox(
		widget.NewLabel("Single File Encryption/Decryption:"),
		widget.NewSeparator(),
		fileSection,
		widget.NewSeparator(),
		algorithmSection,
		widget.NewSeparator(),
		operationSection,
		widget.NewSeparator(),
		outputSection,
		widget.NewSeparator(),
		layout.NewSpacer(),
		processBtn,
	)
}

func createFileSelectionSection(window fyne.Window, selectedFile *string, fileLabel *widget.Label) *fyne.Container {
	selectBtn := widget.NewButton("Select File", func() {
		showFilePicker(window, selectedFile, fileLabel)
	})
	return container.NewVBox(selectBtn, fileLabel)
}

func createAlgorithmSection(algorithmSelect *widget.Select) *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Select Algorithm:"),
		algorithmSelect,
	)
}

func createOperationSection(operationRadio *widget.RadioGroup) *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Select Operation:"),
		operationRadio,
	)
}

func createOutputSection(window fyne.Window, outputFile *string, outputLabel *widget.Label) *fyne.Container {
	selectBtn := widget.NewButton("Choose Output Location", func() {
		showSaveFilePicker(window, outputFile, outputLabel)
	})
	return container.NewVBox(selectBtn, outputLabel)
}

func createProcessButton(window fyne.Window, state *AppState, selectedFile, outputFile *string, algorithmSelect *widget.Select, operationRadio *widget.RadioGroup) *widget.Button {
	return widget.NewButton("Process File", func() {
		processFile(window, state, selectedFile, outputFile, algorithmSelect, operationRadio)
	})
}

func processFile(window fyne.Window, state *AppState, selectedFile, outputFile *string, algorithmSelect *widget.Select, operationRadio *widget.RadioGroup) {
	if *selectedFile == "" {
		dialog.ShowError(fmt.Errorf("please select a file first"), window)
		return
	}

	operation := operationRadio.Selected
	algorithm := algorithmSelect.Selected

	output := determineOutputPath(*selectedFile, *outputFile, operation)

	progress := dialog.NewProgressInfinite("Processing",
		fmt.Sprintf("%sing file with %s...", operation, algorithm), window)
	progress.Show()

	go func() {
		time.Sleep(1 * time.Second) // Simulate processing
		fyne.Do(func() {
			progress.Hide()
		})

		// TODO: Call actual encryption/decryption functions here
		logger.Log(fmt.Sprintf("%sed file: %s → %s (Algorithm: %s)",
			operation, filepath.Base(*selectedFile), filepath.Base(output), algorithm))

		dialog.ShowInformation("Success",
			fmt.Sprintf("File %sed successfully!\nOutput: %s", operation, output), window)
	}()
}

func determineOutputPath(selectedFile, outputFile, operation string) string {
	if outputFile != "" {
		return outputFile
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
