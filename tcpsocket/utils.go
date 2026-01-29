package tcpsocket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func sendData(conn net.Conn, data []byte) error {
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to send file data: %v", err)
	}
	return nil
}

func acceptConnections(saveDir string, onStatus func(string)) {
	for {
		conn, err := activeListener.Accept()
		if err != nil {
			return
		}
		onStatus("New connection from " + conn.RemoteAddr().String())
		go handleIncomingFile(conn, saveDir, onStatus)
	}
}

func handleIncomingFile(conn net.Conn, saveDir string, onStatus func(string)) {
	defer conn.Close()

	metadata, err := receiveMetadata(conn)
	if err != nil {
		onStatus("Metadata error: " + err.Error())
		return
	}

	onStatus("Receiving file: " + metadata.FileName)

	data, err := receiveFileData(conn, metadata.FileSize)
	if err != nil {
		onStatus("Receive error: " + err.Error())
		return
	}

	outputPath := filepath.Join(saveDir, filepath.Base(metadata.FileName))
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		onStatus("Save error: " + err.Error())
		return
	}

	onStatus("File saved: " + outputPath)
}

func receiveMetadata(conn net.Conn) (*FileMetadata, error) {
	metadataLen, err := readUint32(conn)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata length: %v", err)
	}

	metadataBuf := make([]byte, metadataLen)
	if _, err := io.ReadFull(conn, metadataBuf); err != nil {
		return nil, fmt.Errorf("error reading metadata: %v", err)
	}

	var metadata FileMetadata
	if err := json.Unmarshal(metadataBuf, &metadata); err != nil {
		return nil, fmt.Errorf("error parsing metadata: %v", err)
	}

	return &metadata, nil
}

func receiveFileData(conn net.Conn, size int64) ([]byte, error) {
	data := make([]byte, size)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, fmt.Errorf("error reading file data: %v", err)
	}
	return data, nil
}

// Network Helper Functions
func writeUint32(w io.Writer, val uint32) error {
	buf := []byte{
		byte(val >> 24),
		byte(val >> 16),
		byte(val >> 8),
		byte(val),
	}
	_, err := w.Write(buf)
	return err
}

func readUint32(r io.Reader) (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3]), nil
}
