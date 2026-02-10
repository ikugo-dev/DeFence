package algorithms

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

func bytesToint(b []byte) (int, error) {
	if len(b)%4 != 0 {
		return 0, fmt.Errorf("Invalid key size")
	}

	u := binary.BigEndian.Uint32(b[0:4])
	return int(u), nil
}

func bytesToUint32s(b []byte) ([]uint32, error) {
	if len(b)%4 != 0 {
		return nil, fmt.Errorf("Invalid key size")
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
		return nil, fmt.Errorf("Invalid key size")
	}
	return buf.Bytes(), nil
}

func KeyStringTo4Bytes(s string) []byte {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 1 {
		return nil
	}

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return buf
}

func KeyHexStringTo16Bytes(s string) ([]byte, error) {
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		s = s[2:]
	}

	if len(s) == 0 {
		return nil, fmt.Errorf("Hex string is empty")
	}
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("Hex string must have even length")
	}
	if len(s) > 32 {
		return nil, fmt.Errorf("Hex string too long, %d", len(s))
	}

	if len(s) < KeySize*2 {
		pad := make([]byte, (KeySize*2)-len(s))
		for i := range pad {
			pad[i] = '0'
		}
		s = string(pad) + s
	}

	buf := make([]byte, KeySize)
	if _, err := hex.Decode(buf, []byte(s)); err != nil {
		return nil, fmt.Errorf("Invalid hex string: %w", err)
	}

	return buf, nil
}
