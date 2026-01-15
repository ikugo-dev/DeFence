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

	statusLabel := widget.NewLabel("Status: Stopped")
	startBtn, stopBtn := createWatcherButtons(window, state, dirEntry, statusLabel)

	dirSection := createDirectorySection(window, dirEntry, dirLabel)
	controlSection := createControlSection(statusLabel, startBtn, stopBtn)

	return container.NewVBox(
		dirSection,
		widget.NewSeparator(),
		controlSection,
		layout.NewSpacer(),
	)
}

func createDirectorySection(window fyne.Window, dirEntry *widget.Entry, dirLabel *widget.Label) *fyne.Container {
	selectBtn := widget.NewButton("Select Directory", func() {
		showDirectoryPicker(window, dirEntry, dirLabel)
	})

	return container.NewVBox(
		widget.NewLabel("Directory Selection:"),
		container.NewBorder(nil, nil, nil, selectBtn, dirEntry),
		dirLabel,
	)
}

func createControlSection(statusLabel *widget.Label, startBtn, stopBtn *widget.Button) *fyne.Container {
	return container.NewVBox(
		statusLabel,
		container.NewHBox(startBtn, stopBtn),
	)
}

func createWatcherButtons(window fyne.Window, state *AppState, dirEntry *widget.Entry, statusLabel *widget.Label) (*widget.Button, *widget.Button) {
	startBtn := widget.NewButton("Start Watching", nil)
	stopBtn := widget.NewButton("Stop Watching", nil)
	stopBtn.Disable()

	startBtn.OnTapped = func() {
		handleStartWatching(window, state, dirEntry, statusLabel, startBtn, stopBtn)
	}

	stopBtn.OnTapped = func() {
		handleStopWatching(state, statusLabel, startBtn, stopBtn)
	}

	return startBtn, stopBtn
}

func handleStartWatching(window fyne.Window, state *AppState, dirEntry *widget.Entry, statusLabel *widget.Label, startBtn, stopBtn *widget.Button) {
	dir := dirEntry.Text
	if dir == "" {
		dialog.ShowError(fmt.Errorf("please select a directory first"), window)
		return
	}

	if state.isWatching {
		dialog.ShowInformation("Already Watching", "Directory watcher is already running", window)
		return
	}

	state.watchDir = dir
	state.watcherCancel = fsw.InitWatch(dir)
	state.isWatching = true

	statusLabel.SetText("Status: Watching " + dir)
	startBtn.Disable()
	stopBtn.Enable()
	logger.Log("Started watching: " + dir)
}

func handleStopWatching(state *AppState, statusLabel *widget.Label, startBtn, stopBtn *widget.Button) {
	if !state.isWatching {
		return
	}

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

func showDirectoryPicker(window fyne.Window, dirEntry *widget.Entry, dirLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil || dir == nil {
			return
		}
		path := dir.Path()
		dirEntry.SetText(path)
		dirLabel.SetText("Selected: " + path)
	}, window)

	// Set initial directory if possible
	if fd := dialog.NewFolderOpen(nil, window); fd != nil {
		if listableURI != nil {
			fd.SetLocation(listableURI)
		}
	}
}
