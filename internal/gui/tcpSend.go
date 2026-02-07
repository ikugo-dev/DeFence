package gui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/internal/algorithms"
	"github.com/ikugo-dev/DeFence/internal/logger"
	"github.com/ikugo-dev/DeFence/internal/tcpsocket"
)

func createSendFileSection(window fyne.Window, state *AppState) fyne.CanvasObject {
	var selectedFile, outputFile string
	fileLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("No output location selected (file will be sent but not saved locally)")

	selectFileBtn := widget.NewButton("Select File to Send", func() {
		showFilePicker(window, &selectedFile, fileLabel)
	})

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("Recipient IP (e.g., 192.168.1.100)")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port (e.g., 8080)")
	portEntry.SetText("8080")

	cryptoUI := CreateCryptoUIComponents(false) // false = no operation radio (always encrypts)

	selectOutputBtn := widget.NewButton("Save Encrypted Copy Locally (Optional)", func() {
		showSaveFilePicker(window, &outputFile, outputLabel)
	})

	sendBtn := widget.NewButton("Send File", func() {
		if selectedFile == "" {
			dialog.ShowError(fmt.Errorf("please select a file to send"), window)
			return
		}
		if ipEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter recipient IP address"), window)
			return
		}
		if portEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter port number"), window)
			return
		}
		if err := cryptoUI.ValidateKey(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		address := ipEntry.Text
		port := portEntry.Text
		algorithm := cryptoUI.GetAlgorithm()
		key := cryptoUI.GetKey()

		progress := dialog.NewProgressInfinite("Sending", "Encrypting and sending file...", window)
		progress.Show()

		go func() {
			// Encrypt the file
			encryptedData, err := algorithms.EncryptFile(selectedFile, key, algorithm)
			if err != nil {
				fyne.Do(func() {
					progress.Hide()
					dialog.ShowError(fmt.Errorf("encryption failed: %s", err), window)
				})
				return
			}

			// Optionally save encrypted copy locally
			var output string
			if outputFile != "" {
				output = outputFile + ".enc"
				if err := os.WriteFile(output, encryptedData, 0644); err != nil {
					logger.Log("Warning: failed to save encrypted copy: %s", err)
				} else {
					logger.Log("Saved encrypted copy to: %s", output)
				}
			}

			// Send the file
			err = tcpsocket.SendFile(address, port, encryptedData)
			
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to send file: %v", err), window)
					return
				}
				
				msg := fmt.Sprintf("Sent file %s to %s:%s successfully (Algorithm: %s)", 
					filepath.Base(selectedFile), address, port, algorithm)
				
				if outputFile != "" {
					msg += fmt.Sprintf("\nEncrypted copy saved to: %s", filepath.Base(output))
				}
				
				dialog.ShowInformation("Success", msg, window)
			})
		}()
	})

	return container.NewVBox(
		widget.NewLabel("Send Encrypted File:"),
		widget.NewSeparator(),
		selectFileBtn,
		fileLabel,
		widget.NewSeparator(),
		widget.NewLabel("Recipient Details:"),
		ipEntry,
		portEntry,
		CreateCryptoUISection(cryptoUI),
		widget.NewSeparator(),
		selectOutputBtn,
		outputLabel,
		widget.NewSeparator(),
		layout.NewSpacer(),
		sendBtn,
	)
}
