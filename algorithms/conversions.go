package algorithms

import (
	"bytes"
	"encoding/binary"
)

func bytesToint(b []byte) (value int, ok bool) {
	if len(b)%4 != 0 {
		return 0, false
	}

	u := binary.LittleEndian.Uint32(b[0:4])
	return int(u), true
}

func bytesToUint32s(b []byte) (value []uint32, ok bool) {
	if len(b)%4 != 0 {
		return nil, false
	}
	u := make([]uint32, len(b)/4)
	for i := range u {
		u[i] = binary.LittleEndian.Uint32(b[i*4 : (i+1)*4])
	}
	return u, true
}

func uint32ToBytes(u []uint32) (value []byte, ok bool) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, u)
	if err != nil {
		return nil, false
	}
	return buf.Bytes(), true
}
