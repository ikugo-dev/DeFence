package tcpsocket

import (
	"io"
	"net"

	"github.com/ikugo-dev/DeFence/logger"
)

var activeListener net.Listener

func SendFile(address, port string, encryptedData []byte) error {
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write(encryptedData); err != nil {
		return err
	}
	return nil
}

func StartListening(port string, dataCh chan<- []byte) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer ln.Close()
	logger.Log("Listening on :%s", port)
	activeListener = ln

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		chunk := make([]byte, n)
		copy(chunk, buf[:n])
		dataCh <- chunk

		if err != nil {
			if err == io.EOF {
				logger.Log("Connection closed by sender")
				close(dataCh)
				return nil
			}
			logger.Log("Stopped listening on :%s", port)
			return err
		}
	}
}

func CollectAll(dataCh <-chan []byte) []byte {
	var all []byte
	for chunk := range dataCh {
		all = append(all, chunk...)
	}
	return all
}

func StopListening(onStatus func(string)) {
	if activeListener != nil {
		activeListener.Close()
		activeListener = nil
		onStatus("Stopped listening")
	}
}
