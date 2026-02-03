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
	"github.com/ikugo-dev/DeFence/tcpsocket"
)

func createSendFileSection(window fyne.Window, state *AppState) fyne.CanvasObject {
	var selectedFile string
	fileLabel := widget.NewLabel("No file selected")
	selectFileBtn := widget.NewButton("Select File to Send", func() {
		showFilePicker(window, &selectedFile, fileLabel)
	})

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("Recipient IP (e.g., 192.168.1.100)")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port (e.g., 8080)")
	portEntry.SetText("8080")

	algorithmSelect := widget.NewSelect([]string{"Railfence Cipher", "XXTEA", "CBC"}, nil)
	algorithmSelect.SetSelected("XXTEA")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter encryption key")
	keyEntry.Password = true

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
		if keyEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Please enter a key"), window)
			return
		}

		address := ipEntry.Text
		port := portEntry.Text
		algorithm := algorithmSelect.Selected
		key := algorithms.KeyStringToBigEndianBytes(keyEntry.Text)

		progress := dialog.NewProgressInfinite("Sending", "Encrypting and sending file...", window)
		progress.Show()

		go func() {
			encryptedData, err := algorithms.EncryptFile(selectedFile, key, algorithm)
			fyne.Do(func() {
				if err != nil {
					logger.LogWithDialog(window, "Error", "Error while encrypting: %s", err)
					return
				}
				progress.Hide()
				err = tcpsocket.SendFile(address, port, encryptedData)
				if err != nil {
					logger.LogWithDialog(window, "Error", "Failed to send file: %v", err.Error())
					return
				}
				logger.LogWithDialog(window, "Success", "Sent file %s to %s successfully (Algorithm: %s)", filepath.Base(selectedFile), port, algorithm)
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
		widget.NewSeparator(),
		widget.NewLabel("Select Algorithm:"),
		algorithmSelect,
		widget.NewLabel("Encryption / Decryption Key:"),
		keyEntry,
		widget.NewSeparator(),
		layout.NewSpacer(),
		sendBtn,
	)
}
