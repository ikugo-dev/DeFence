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
	"github.com/ikugo-dev/DeFence/internal/logger"
	"github.com/ikugo-dev/DeFence/internal/tcpsocket"
)

func createReceiveFileSection(window fyne.Window, state *AppState) fyne.CanvasObject {
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port to listen on (e.g., 8080)")
	portEntry.SetText("8080")

	saveDir := "./received_files/"
	saveDirLabel := widget.NewLabel("Save to: " + saveDir)

	// Create shared crypto UI components (no operation radio - always decrypts)
	cryptoUI := CreateCryptoUIComponents(false)

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

	ipLabel := widget.NewLabel("")
	ipLabel.Hide()

	showIPBtn := widget.NewButton("Show Local IP", nil)
	ipVisible := false

	showIPBtn.OnTapped = func() {
		ipVisible = !ipVisible
		if ipVisible {
			ip, ok := tcpsocket.GetLocalIP()
			if ok {
				ipLabel.SetText("Local IP: " + ip)
			} else {
				ipLabel.SetText("Local IP: could not be determined")
			}
			ipLabel.Show()
			showIPBtn.SetText("Hide Local IP")
		} else {
			ipLabel.Hide()
			showIPBtn.SetText("Show Local IP")
		}
	}

	statusLabel := widget.NewLabel("Status: Not listening")
	startBtn := widget.NewButton("Start Listening", nil)
	stopBtn := widget.NewButton("Stop Listening", nil)
	stopBtn.Disable()

	startBtn.OnTapped = func() {
		if err := cryptoUI.ValidateKey(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		config := cryptoUI.GetConfig()

		dataCh := make(chan []byte)
		go func() {
			err := tcpsocket.StartListening(portEntry.Text, dataCh)
			if err != nil {
				logger.LogWithDialog(window, "Error", "Error while trying to listen: %s", err)
			}
		}()

		go func() {
			allData := tcpsocket.CollectAll(dataCh)

			// Generate output filename with timestamp
			timestamp := time.Now().Format("20060102_150405")
			outputPath := filepath.Join(saveDir, fmt.Sprintf("received_%s.dec", timestamp))

			// Decrypt and save the file
			err := DecryptAndSave(allData, outputPath, config.Key)
			if err != nil {
				logger.LogWithDialog(window, "Error", "Decrypting failed: %v", err)
				return
			}

			logger.LogWithDialog(window, "Success",
				"Received and decrypted file successfully!\nSaved to: %s\nSize: %d bytes",
				outputPath, len(allData))
		}()

		statusLabel.SetText("Status: Listening on port " + portEntry.Text)
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
		showIPBtn,
		ipLabel,
		widget.NewSeparator(),
		widget.NewLabel("Listening Port:"),
		portEntry,
		widget.NewSeparator(),
		selectDirBtn,
		saveDirLabel,
		CreateCryptoUISection(cryptoUI),
		widget.NewSeparator(),
		statusLabel,
		container.NewHBox(startBtn, stopBtn),
		layout.NewSpacer(),
	)
}
