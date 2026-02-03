package tcpsocket

import "net"

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
