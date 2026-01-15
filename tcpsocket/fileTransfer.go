package tcpsocket

import (
	"fmt"
	"net"
	"os"
	"time"
)

type FileMetadata struct {
	FileName            string    `json:"fileName"`
	FileSize            int64     `json:"fileSize"`
	CreationDateTime    time.Time `json:"creationDateTime"`
	EncryptionAlgorithm string    `json:"encryptionAlgorithm"`
	HashingAlgorithm    string    `json:"hashingAlgorithm"`
}

var activeListener net.Listener

func SendFile(filePath, address, algorithm, hash string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	metadata := createMetadata(filePath, algorithm, hash)
	encryptedData := encryptFileData(fileData, algorithm, hash)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	if err := sendData(conn, metadata, encryptedData); err != nil {
		return err
	}

	return nil
}

func StartListening(port string, saveDir string, onStatus func(string)) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	activeListener = listener
	onStatus("Listening on port " + port)

	go acceptConnections(saveDir, onStatus)
	return nil
}

func StopListening(onStatus func(string)) {
	if activeListener != nil {
		activeListener.Close()
		activeListener = nil
		onStatus("Stopped listening")
	}
}
