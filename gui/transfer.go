package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
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

func createReceiveFileSection(window fyne.Window, state *AppState) fyne.CanvasObject {
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port to listen on (e.g., 8080)")
	portEntry.SetText("8080")

	saveDir := "./received_files/"
	saveDirLabel := widget.NewLabel("Save to: " + saveDir)

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter encryption key")
	keyEntry.Password = true

	selectDirBtn := widget.NewButton("Choose Save Directory", func() {
		currentDir := getCurrentDirectory()
		listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

		fd := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}
			saveDir = dir.Path()
			saveDirLabel.SetText("Save to: " + saveDir)
		}, window)

		if listableURI != nil {
			fd.SetLocation(listableURI)
		}
		fd.Show()
	})

	statusLabel := widget.NewLabel("Status: Not listening")
	startBtn := widget.NewButton("Start Listening", nil)
	stopBtn := widget.NewButton("Stop Listening", nil)
	stopBtn.Disable()

	startBtn.OnTapped = func() {
		if keyEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Please enter a key"), window)
			return
		}
		key := algorithms.KeyStringToBigEndianBytes(keyEntry.Text)

		dataCh := make(chan []byte)
		go func() {
			err := tcpsocket.StartListening(portEntry.Text, dataCh)
			if err != nil {
				logger.LogWithDialog(window, "Error", "Error while trying to listen; %s", err)
			}
		}()
		go func() {
			allData := tcpsocket.CollectAll(dataCh)

			decrypted, err := algorithms.DecryptFileData(allData, key)
			if err != nil {
				logger.LogWithDialog(window, "Error", "Decrypting failed: %v", err)
				return
			}
			logger.LogWithDialog(window, "Success", "Received %d bytes total", len(decrypted))
		}()

		startBtn.Disable()
		stopBtn.Enable()
		portEntry.Disable()
	}

	stopBtn.OnTapped = func() {
		tcpsocket.StopListening(func(status string) {
			logger.Log("%s", status)
			statusLabel.SetText("Status: " + status)
		})

		startBtn.Enable()
		stopBtn.Disable()
		portEntry.Enable()
	}

	return container.NewVBox(
		widget.NewLabel("Receive Encrypted File:"),
		widget.NewSeparator(),
		widget.NewLabel("Listening Port:"),
		portEntry,
		widget.NewSeparator(),
		selectDirBtn,
		saveDirLabel,
		widget.NewSeparator(),
		widget.NewLabel("Encryption / Decryption Key:"),
		keyEntry,
		widget.NewSeparator(),
		statusLabel,
		container.NewHBox(startBtn, stopBtn),
		layout.NewSpacer(),
	)
}
