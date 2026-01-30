package gui

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
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

func showFilePicker(window fyne.Window, selectedFile *string, fileLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	fd := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if err != nil || file == nil {
			return
		}
		*selectedFile = file.URI().Path()
		fileLabel.SetText("Selected: " + filepath.Base(*selectedFile))
		file.Close()
	}, window)

	if listableURI != nil {
		fd.SetLocation(listableURI)
	}
	fd.Show()
}

func showSaveFilePicker(window fyne.Window, outputFile *string, outputLabel *widget.Label) {
	currentDir := getCurrentDirectory()
	listableURI, _ := storage.ListerForURI(storage.NewFileURI(currentDir))

	fd := dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) {
		if err != nil || file == nil {
			return
		}
		*outputFile = file.URI().Path()
		outputLabel.SetText("Output: " + filepath.Base(*outputFile))
		file.Close()
	}, window)

	if listableURI != nil {
		fd.SetLocation(listableURI)
	}
	fd.Show()
}
