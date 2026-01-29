package algorithms

import "log"

const DELTA = 0x9e3779b9
const MinDataSize = 16
const KeySize = 16

func MX(y, z, sum, p, e uint32, key []uint32) uint32 {
	return (((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (key[(p&3)^e] ^ z)))
}

func EncryptXXTEA(data []byte, key []byte) []byte {
	if len(data)%4 != 0 {
		log.Fatal("data length must be multiple of 4 bytes")
	}
	if len(data)/4 < 2 {
		log.Fatal("data must contain at least two uint32 words")
	}
	if len(data) < MinDataSize {
		log.Fatal("invalid data size")
	}
	if len(key) != KeySize {
		log.Fatal("invalid key size")
	}

	v, ok := bytesToUint32s(data)
	if !ok {
		log.Fatal("invalid key")
	}
	n := len(v)
	k, ok := bytesToUint32s(key)
	if !ok {
		log.Fatal("invalid key")
	}

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
	result, ok := uint32ToBytes(v)
	if !ok {
		log.Fatal("cannot convert encripted result back to bytes")
	}
	return result
}

func DecryptXXTEA(data []byte, key []byte) []byte {
	if len(data) < MinDataSize {
		log.Fatal("invalid data size")
	}
	if len(key) != KeySize {
		log.Fatal("invalid key size")
	}

	v, ok := bytesToUint32s(data)
	if !ok {
		log.Fatal("invalid key")
	}
	n := len(v)
	k, ok := bytesToUint32s(key)
	if !ok {
		log.Fatal("invalid key")
	}

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
	result, ok := uint32ToBytes(v)
	if !ok {
		log.Fatal("cannot convert encripted result back to bytes")
	}
	return result
}
