package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var logChan chan string = make(chan string, 100)
var logs []string = make([]string, 0)
var mu sync.Mutex

func Log(format string, args ...any) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf("[%s] %s", timestamp, fmt.Sprintf(format, args...))

	mu.Lock()
	logs = append(logs, message)
	mu.Unlock()

	select {
	case logChan <- message:
	default:
	}
}

func GetLogs() string {
	mu.Lock()
	defer mu.Unlock()

	var result strings.Builder
	for _, log := range logs {
		result.WriteString(log + "\n")
	}
	return result.String()
}

func Subscribe() <-chan string {
	return logChan
}

func Clear() {
	mu.Lock()
	logs = logs[:0]
	mu.Unlock()
}

func LogWithDialog(window fyne.Window, dialogType, format string, args ...any) {
	switch dialogType {
	case "Error":
		dialog.ShowError(fmt.Errorf(format, args...), window)
	default:
		dialog.ShowInformation(dialogType, fmt.Sprintf(format, args...), window)
	}
	Log(format, args...)
}
