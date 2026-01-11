package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		running := os.Args[0]
		if strings.Contains(running, "/tmp/") {
			running = "go run *.go"
		}
		log.Fatalf("Usage: %s <input file>", running)
	}
	fileName := os.Args[1]
	_, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Could not read file " + fileName)
	}
	cancelWatcherFunc := InitWatch("./testDir")

	<-make(chan struct{})
	cancelWatcherFunc()
}

func createMetadata(fileName string) []byte {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Fatal("Could not read file " + fileName)
	}
	jsonContent, err := json.Marshal(map[string]any{
		"fileName":            fileName,
		"fileSize":            fileInfo.Size(),
		"creationDateTime":    fileInfo.ModTime(),
		"encryptionAlgorithm": "",
		"hashingAlgorithm":    "",
	})
	if err != nil {
		log.Fatal("Could not encode metadata for file " + fileName)
	}
	return jsonContent
}

// TODO: Implement these functions
func RailfenceCipher() {}
func XXTEA()           {}
func CBC()             {}
