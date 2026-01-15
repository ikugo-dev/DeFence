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
	"github.com/ikugo-dev/DeFence/logger"
	"github.com/ikugo-dev/DeFence/tcpsocket"
)

// Send File Section
func createSendFileSection(window fyne.Window, state *AppState) fyne.CanvasObject {
	var selectedFile string
	fileLabel := widget.NewLabel("No file selected")

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("Recipient IP (e.g., 192.168.1.100)")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port (e.g., 8080)")
	portEntry.SetText("8080")

	algorithmSelect := widget.NewSelect([]string{"Railfence Cipher", "XXTEA", "CBC"}, nil)
	algorithmSelect.SetSelected("XXTEA")

	fileSection := createSendFileSelection(window, &selectedFile, fileLabel)
	recipientSection := createRecipientSection(ipEntry, portEntry)
	algorithmSection := createAlgorithmSection(algorithmSelect)
	sendBtn := createSendButton(window, state, &selectedFile, ipEntry, portEntry, algorithmSelect)

	return container.NewVBox(
		widget.NewLabel("Send Encrypted File:"),
		widget.NewSeparator(),
		fileSection,
		widget.NewSeparator(),
		recipientSection,
		widget.NewSeparator(),
		algorithmSection,
		widget.NewSeparator(),
		layout.NewSpacer(),
		sendBtn,
	)
}

func createSendFileSelection(window fyne.Window, selectedFile *string, fileLabel *widget.Label) *fyne.Container {
	selectBtn := widget.NewButton("Select File to Send", func() {
		showFilePicker(window, selectedFile, fileLabel)
	})
	return container.NewVBox(selectBtn, fileLabel)
}

func createRecipientSection(ipEntry, portEntry *widget.Entry) *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Recipient Details:"),
		ipEntry,
		portEntry,
	)
}

func createSendButton(window fyne.Window, state *AppState, selectedFile *string, ipEntry, portEntry *widget.Entry, algorithmSelect *widget.Select) *widget.Button {
	return widget.NewButton("Send File", func() {
		handleSendFile(window, state, selectedFile, ipEntry, portEntry, algorithmSelect)
	})
}

func handleSendFile(window fyne.Window, state *AppState, selectedFile *string, ipEntry, portEntry *widget.Entry, algorithmSelect *widget.Select) {
	if *selectedFile == "" {
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

	address := ipEntry.Text + ":" + portEntry.Text
	algorithm := algorithmSelect.Selected

	progress := dialog.NewProgressInfinite("Sending", "Encrypting and sending file...", window)
	progress.Show()

	go func() {
		err := tcpsocket.SendFile(*selectedFile, address, algorithm, "tigerHash")
		progress.Hide()

		if err != nil {
			logger.Log(fmt.Sprintf("Failed to send file: %v", err))
			dialog.ShowError(fmt.Errorf("failed to send file: %v", err), window)
		} else {
			logger.Log(fmt.Sprintf("Sent file %s to %s (Algorithm: %s)",
				filepath.Base(*selectedFile), address, algorithm))
			dialog.ShowInformation("Success", "File sent successfully!", window)
		}
	}()
}

// Receive File Section
func createReceiveFileSection(window fyne.Window, state *AppState) fyne.CanvasObject {
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port to listen on (e.g., 8080)")
	portEntry.SetText("8080")

	saveDirLabel := widget.NewLabel("Save to: ./received_files/")
	saveDir := "./received_files/"

	statusLabel := widget.NewLabel("Status: Not listening")
	startBtn, stopBtn := createReceiveButtons(window, state, portEntry, &saveDir, saveDirLabel, statusLabel)

	portSection := createPortSection(portEntry)
	saveSection := createSaveDirectorySection(window, &saveDir, saveDirLabel)
	controlSection := createReceiveControlSection(statusLabel, startBtn, stopBtn)

	return container.NewVBox(
		widget.NewLabel("Receive Encrypted File:"),
		widget.NewSeparator(),
		portSection,
		widget.NewSeparator(),
		saveSection,
		widget.NewSeparator(),
		controlSection,
		layout.NewSpacer(),
	)
}

func createPortSection(portEntry *widget.Entry) *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Listening Port:"),
		portEntry,
	)
}

func createSaveDirectorySection(window fyne.Window, saveDir *string, saveDirLabel *widget.Label) *fyne.Container {
	selectBtn := widget.NewButton("Choose Save Directory", func() {
		showDirectoryPickerForSave(window, saveDir, saveDirLabel)
	})
	return container.NewVBox(selectBtn, saveDirLabel)
}

func createReceiveControlSection(statusLabel *widget.Label, startBtn, stopBtn *widget.Button) *fyne.Container {
	return container.NewVBox(
		statusLabel,
		container.NewHBox(startBtn, stopBtn),
	)
}

func createReceiveButtons(window fyne.Window, state *AppState, portEntry *widget.Entry, saveDir *string, saveDirLabel, statusLabel *widget.Label) (*widget.Button, *widget.Button) {
	startBtn := widget.NewButton("Start Listening", nil)
	stopBtn := widget.NewButton("Stop Listening", nil)
	stopBtn.Disable()

	startBtn.OnTapped = func() {
		err := tcpsocket.StartListening(
			portEntry.Text,
			*saveDir,
			func(status string) {
				logger.Log(status)
				statusLabel.SetText("Status: " + status)
			},
		)

		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		startBtn.Disable()
		stopBtn.Enable()
		portEntry.Disable()
	}

	stopBtn.OnTapped = func() {
		tcpsocket.StopListening(func(status string) {
			logger.Log(status)
			statusLabel.SetText("Status: " + status)
		})

		startBtn.Enable()
		stopBtn.Disable()
		portEntry.Enable()
	}

	return startBtn, stopBtn
}

func showDirectoryPickerForSave(window fyne.Window, saveDir *string, saveDirLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	fd := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil || dir == nil {
			return
		}
		*saveDir = dir.Path()
		saveDirLabel.SetText("Save to: " + *saveDir)
	}, window)

	if listableURI != nil {
		fd.SetLocation(listableURI)
	}
	fd.Show()
}
