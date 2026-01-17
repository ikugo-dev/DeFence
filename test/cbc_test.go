package test

import (
	"bytes"
	"testing"

	a "github.com/ikugo-dev/DeFence/algorithms"
)

type CBCTest struct {
	input string
	key   []byte
	iv    []byte
}

func TestEncryptDecryptCBC(t *testing.T) {
	tests := []CBCTest{
		{
			input: "short message",
			key:   []byte("1234567890abcdef"),
			iv:    []byte("abcdefghijklmnop"),
		},
		{
			input: "this message spans multiple blocks and needs padding",
			key:   []byte("examplekey123456"),
			iv:    []byte("iviviviviviviviv"),
		},
		{
			input: "exactly sixteen byte",
			key:   []byte("deadbeefdeadbeef"),
			iv:    []byte("feedfacecafebeef"),
		},
	}

	for _, tt := range tests {
		encrypted := a.EncryptCBC([]byte(tt.input), tt.key, tt.iv)
		decrypted := a.DecryptCBC(encrypted, tt.key)

		if !bytes.Equal(decrypted, []byte(tt.input)) {
			t.Errorf(
				"CBC round-trip failed: input=%q output=%q",
				tt.input,
				decrypted,
			)
		}
	}
}
