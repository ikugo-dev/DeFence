package test

import (
	"bytes"
	"testing"

	a "github.com/ikugo-dev/DeFence/algorithms"
)

type XXTEATest struct {
	input string
	key   []byte
}

func TestEncryptDecryptXXTEA(t *testing.T) {
	tests := []XXTEATest{
		{
			input: "abcdefghijklmnop",
			key:   []byte("1234567890abcdef"),
		},
		{
			input: "this is exactly 32 bytes long!!",
			key:   []byte("examplekey123456"),
		},
		{
			input: "four blocks of data....four blocks of data....",
			key:   []byte("deadbeefdeadbeef"),
		},
	}

	for _, tt := range tests {
		encrypted := a.EncryptXXTEA([]byte(tt.input), tt.key)
		decrypted := a.DecryptXXTEA(encrypted, tt.key)

		if !bytes.Equal(decrypted, []byte(tt.input)) {
			t.Errorf(
				"XXTEA round-trip failed: input=%q output=%q",
				tt.input,
				decrypted,
			)
		}
	}
}
