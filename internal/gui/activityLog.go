package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ikugo-dev/DeFence/internal/logger"
)

func createActivityLogTab() fyne.CanvasObject {

	logText := widget.NewMultiLineEntry()
	logText.Disable()
	logText.SetText(logger.GetLogs())

	clearBtn := widget.NewButton("Clear Log", func() {
		logger.Clear()
		logText.SetText("")
	})

	go func() {
		for msg := range logger.Subscribe() {
			fyne.Do(func() {
				logText.Append(msg + "\n")
			})
		}
	}()

	logScroll := container.NewScroll(logText)

	return container.NewBorder(
		container.NewHBox(widget.NewLabel("Activity Log"), clearBtn),
		nil, nil, nil,
		logScroll,
	)
}
