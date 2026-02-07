package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createSingleFileTab(window fyne.Window, state *AppState) fyne.CanvasObject {
	var selectedFile, outputFile string
	fileLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("Output: Same directory as input file")

	selectFileBtn := widget.NewButton("Select File", func() {
		showFilePicker(window, &selectedFile, fileLabel)
	})

	// Create shared crypto UI components
	cryptoUI := CreateCryptoUIComponents(true) // true = include operation radio

	selectOutputBtn := widget.NewButton("Choose Output Location", func() {
		showSaveFilePicker(window, &outputFile, outputLabel)
	})

	processBtn := widget.NewButton("Process File", func() {
		if selectedFile == "" {
			dialog.ShowError(fmt.Errorf("please select a file first"), window)
			return
		}
		if err := cryptoUI.ValidateKey(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		config := cryptoUI.GetConfig()
		output := DetermineOutputPath(selectedFile, outputFile, config.Operation)

		progress := dialog.NewProgressInfinite("Processing", 
			fmt.Sprintf("%sing file with %s...", config.Operation, config.Algorithm), window)
		progress.Show()

		go func() {
			err := ProcessAndSaveFile(selectedFile, output, config)
			
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(fmt.Errorf("error while processing: %s", err), window)
					return
				}
				dialog.ShowInformation("Success", 
					fmt.Sprintf("%sed file: %s → %s successfully (Algorithm: %s)", 
						config.Operation, 
						filepath.Base(selectedFile), 
						filepath.Base(output), 
						config.Algorithm), 
					window)
			})
		}()
	})

	return container.NewVBox(
		widget.NewLabel("Single File Encryption/Decryption:"),
		widget.NewSeparator(),
		selectFileBtn,
		fileLabel,
		CreateCryptoUISection(cryptoUI),
		widget.NewSeparator(),
		selectOutputBtn,
		outputLabel,
		widget.NewSeparator(),
		layout.NewSpacer(),
		processBtn,
	)
}
