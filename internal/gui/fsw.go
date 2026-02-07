package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/internal/fsw"
	"github.com/ikugo-dev/DeFence/internal/logger"
)

func createWatcherTab(window fyne.Window, state *AppState) fyne.CanvasObject {
	dirEntry := widget.NewEntry()
	dirEntry.SetPlaceHolder("Enter directory path...")
	dirEntry.Text = getCurrentDirectory()
	dirLabel := widget.NewLabel("No directory selected")

	selectBtn := widget.NewButton("Select Directory", func() {
		showDirectoryPicker(window, dirEntry, dirLabel)
	})

	cryptoUI := CreateCryptoUIComponents(false) // false = no operation radio (always encrypts)

	statusLabel := widget.NewLabel("Status: Stopped")
	startBtn := widget.NewButton("Start Watching", nil)
	stopBtn := widget.NewButton("Stop Watching", nil)
	stopBtn.Disable()

	startBtn.OnTapped = func() {
		if dirEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please select a directory first"), window)
			return
		}
		if state.isWatching {
			dialog.ShowInformation("Already Watching", "Directory watcher is already running", window)
			return
		}
		if err := cryptoUI.ValidateKey(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		algorithm := cryptoUI.GetAlgorithm()
		key := cryptoUI.GetKey()

		state.watchDir = dirEntry.Text
		state.watcherCancel = fsw.InitWatch(dirEntry.Text, key, algorithm)
		state.isWatching = true

		statusLabel.SetText("Status: Watching " + dirEntry.Text)
		startBtn.Disable()
		stopBtn.Enable()
		logger.Log("Started watching: %s", dirEntry.Text)
	}

	stopBtn.OnTapped = func() {
		if state.isWatching {
			if state.watcherCancel != nil {
				state.watcherCancel()
				state.watcherCancel = nil
			}
			state.isWatching = false
			statusLabel.SetText("Status: Stopped")
			startBtn.Enable()
			stopBtn.Disable()
			logger.Log("Stopped watching")
		}
	}

	return container.NewVBox(
		widget.NewLabel("Directory Selection:"),
		container.NewBorder(nil, nil, nil, selectBtn, dirEntry),
		dirLabel,
		CreateCryptoUISection(cryptoUI),
		widget.NewSeparator(),
		statusLabel,
		container.NewHBox(startBtn, stopBtn),
		layout.NewSpacer(),
	)
}
