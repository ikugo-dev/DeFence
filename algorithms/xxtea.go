package algorithms

import (
	"bytes"
	"encoding/binary"
	"log"
)

const DELTA = 0x9e3779b9
const MinDataSize = 16
const KeySize = 16

func MX(y, z, sum, p, e uint32, key []uint32) uint32 {
	return (((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (key[(p&3)^e] ^ z)))
}

func EncryptXXTEA(data []byte, key []byte) []byte {
	if len(data) < MinDataSize {
		log.Fatal("invalid data size")
	}
	if len(key) != KeySize {
		log.Fatal("invalid key size")
	}
	v := bytesToUint32s(data)
	n := len(v)
	k := bytesToUint32s(key)

	var y, z, sum uint32
	rounds := uint32(6 + 52/n)
	sum = 0
	z = v[n-1]

	for rounds > 0 {
		sum += DELTA
		e := (sum >> 2) & 3
		for p := 0; p < n-1; p++ {
			y = v[p+1]
			v[p] += MX(y, z, sum, uint32(p), e, k)
			z = v[p]
		}
		y = v[0]
		v[n-1] += MX(y, z, sum, uint32(n-1), e, k)
		z = v[n-1]
		rounds--
	}
	return uint32ToBytes(v)
}

func DecryptXXTEA(data []byte, key []byte) []byte {
	if len(data) < MinDataSize {
		log.Fatal("invalid data size")
	}
	if len(key) != KeySize {
		log.Fatal("invalid key size")
	}
	v := bytesToUint32s(data)
	n := len(v)
	k := bytesToUint32s(key)

	var y, z, sum uint32
	rounds := uint32(6 + 52/n)
	sum = rounds * DELTA
	y = v[0]
	for rounds > 0 {
		e := (sum >> 2) & 3
		for p := n - 1; p > 0; p-- {
			z = v[p-1]
			v[p] -= MX(y, z, sum, uint32(p), e, k)
			y = v[p]
		}
		z = v[n-1]
		v[0] -= MX(y, z, sum, 0, e, k)
		y = v[0]
		sum -= DELTA
		rounds--
	}
	return uint32ToBytes(v)
}

func bytesToUint32s(b []byte) []uint32 {
	if len(b)%4 != 0 {
		return nil // TODO: error handling
	}

	u := make([]uint32, len(b)/4)
	for i := range u {
		u[i] = binary.LittleEndian.Uint32(b[i*4 : (i+1)*4])
	}
	return u
}

func uint32ToBytes(u []uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, u)
	if err != nil {
		return nil //TODO: error handling
	}
	return buf.Bytes()
}
