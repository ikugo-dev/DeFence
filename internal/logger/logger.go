package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var logChan chan string = make(chan string, 100)
var logs []string = make([]string, 0)
var mu sync.Mutex
var logFile = "log.txt"

func Log(format string, args ...any) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf("[%s] %s", timestamp, fmt.Sprintf(format, args...))

	mu.Lock()
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString(message + "\n")
		f.Close()
	}
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

	data, err := os.ReadFile(logFile)
	if err != nil { // read session logs if file cant be read
		var result strings.Builder
		for _, log := range logs {
			result.WriteString(log + "\n")
		}
		return result.String()
	}
	return string(data)
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
