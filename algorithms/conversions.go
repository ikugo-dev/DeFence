package algorithms

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func bytesToint(b []byte) (int, error) {
	if len(b)%4 != 0 {
		return 0, fmt.Errorf("invalid key size")
	}

	u := binary.BigEndian.Uint32(b[0:4])
	return int(u), nil
}

func bytesToUint32s(b []byte) ([]uint32, error) {
	if len(b)%4 != 0 {
		return nil, fmt.Errorf("invalid key size")
	}
	u := make([]uint32, len(b)/4)
	for i := range u {
		u[i] = binary.BigEndian.Uint32(b[i*4 : (i+1)*4])
	}
	return u, nil
}

func uint32ToBytes(u []uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, u)
	if err != nil {
		return nil, fmt.Errorf("invalid key size")
	}
	return buf.Bytes(), nil
}
