package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type AppState struct {
	watcherCancel context.CancelFunc
	isWatching    bool
	watchDir      string
}

func Start() {
	myApp := app.NewWithID("DeFence")
	myWindow := myApp.NewWindow("File Encryption/Decryption Tool")
	myWindow.Resize(fyne.NewSize(800, 600))

	state := &AppState{}

	tabs := container.NewAppTabs(
		container.NewTabItem("Directory Watcher", createWatcherTab(myWindow, state)),
		container.NewTabItem("Single File", createSingleFileTab(myWindow, state)),
		container.NewTabItem("File Transfer", container.NewAppTabs(
			container.NewTabItem("Send File", createSendFileSection(myWindow, state)),
			container.NewTabItem("Receive File", createReceiveFileSection(myWindow, state)),
		)),
		container.NewTabItem("Activity Log", createActivityLogTab()),
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
