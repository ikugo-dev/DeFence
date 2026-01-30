package gui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/fsw"
	"github.com/ikugo-dev/DeFence/logger"
)

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		logger.Log("Error getting current directory: %v", err)
		return "."
	}
	return dir
}

func createWatcherTab(window fyne.Window, state *AppState) fyne.CanvasObject {
	dirEntry := widget.NewEntry()
	dirEntry.SetPlaceHolder("Enter directory path...")
	dirEntry.Text = getCurrentDirectory()
	dirLabel := widget.NewLabel("No directory selected")

	selectBtn := widget.NewButton("Select Directory", func() {
		showDirectoryPicker(window, dirEntry, dirLabel)
	})

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

		state.watchDir = dirEntry.Text
		state.watcherCancel = fsw.InitWatch(dirEntry.Text)
		state.isWatching = true

		statusLabel.SetText("Status: Watching " + dirEntry.Text)
		startBtn.Disable()
		stopBtn.Enable()
		logger.Log("Started watching: " + dirEntry.Text)
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
		widget.NewSeparator(),
		statusLabel,
		container.NewHBox(startBtn, stopBtn),
		layout.NewSpacer(),
	)
}

func showDirectoryPicker(window fyne.Window, dirEntry *widget.Entry, dirLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	fd := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil || dir == nil {
			return
		}
		path := dir.Path()
		dirEntry.SetText(path)
		dirLabel.SetText("Selected: " + path)
	}, window)

	if listableURI != nil {
		fd.SetLocation(listableURI)
	}
	fd.Show()
}
