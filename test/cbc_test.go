package test

import (
	"bytes"
	"testing"

	a "github.com/ikugo-dev/DeFence/internal/algorithms"
)

type CBCTest struct {
	input string
	key   []byte
}

func TestEncryptDecryptCBC(t *testing.T) {
	tests := []CBCTest{
		{
			input: "short message",
			key:   []byte("1234567890abcdef"),
		},
		{
			input: "this message spans multiple blocks and needs padding",
			key:   []byte("examplekey123456"),
		},
		{
			input: "exactly sixteen byte",
			key:   []byte("deadbeefdeadbeef"),
		},
	}

	for _, tt := range tests {
		encrypted, _ := a.EncryptRawData([]byte(tt.input), tt.key, "CBC")
		decrypted, _ := a.DecryptRawData(encrypted, tt.key, "CBC")

		if !bytes.Equal(decrypted, []byte(tt.input)) {
			t.Errorf(
				"CBC round-trip failed: input=%q output=%q",
				tt.input,
				decrypted,
			)
		}
	}
}
