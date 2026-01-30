package tcpsocket

import (
	"fmt"
	"net"

	"github.com/ikugo-dev/DeFence/algorithms"
)

var activeListener net.Listener

func SendFile(fileName, address, algorithm string, key []byte) error {
	encryptedData := algorithms.EncryptFile(fileName, key, algorithm)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	if err := sendData(conn, encryptedData); err != nil {
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
