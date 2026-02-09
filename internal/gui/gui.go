package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func Start() {
	myApp := app.NewWithID("DeFence")
	myWindow := myApp.NewWindow("DeFence")
	myWindow.Resize(fyne.NewSize(800, 600))

	tabs := container.NewAppTabs(
		container.NewTabItem("Single File", createSingleFileTab(myWindow)),
		container.NewTabItem("Directory Watcher", createWatcherTab(myWindow)),
		container.NewTabItem("File Transfer", container.NewAppTabs(
			container.NewTabItem("Send File", createSendFileSection(myWindow)),
			container.NewTabItem("Receive File", createReceiveFileSection(myWindow)),
		)),
		container.NewTabItem("Activity Log", createActivityLogTab()),
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
