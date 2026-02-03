package tcpsocket

import (
	"io"
	"net"

	"github.com/ikugo-dev/DeFence/logger"
)

func GetLocalIP() (ip string, ok bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown IP", false
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), true
			}
		}
	}
	return "Unknown IP", false
}

var activeListener net.Listener

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
