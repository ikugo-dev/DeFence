package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/internal/algorithms"
	"github.com/ikugo-dev/DeFence/internal/logger"
	"github.com/ikugo-dev/DeFence/internal/tcpsocket"
)

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
		showIPBtn,
		ipLabel,
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
