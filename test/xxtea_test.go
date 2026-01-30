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
			input: "this is exactly 32 bytes long!!!",
			key:   []byte("examplekey123456"),
		},
		{
			input: "abcdefghijklmnopabcdefghijklmnopabcdefghijklmnopabcdefghijklmnop",
			key:   []byte("abcdefghijklmnop"),
		},
	}

	for _, tt := range tests {
		encrypted := a.EncryptFileData([]byte(tt.input), tt.key, "XXTEA")
		decrypted := a.DecryptFileData(encrypted, tt.key, "XXTEA")

		if !bytes.Equal(decrypted, []byte(tt.input)) {
			t.Errorf(
				"XXTEA round-trip failed: input=%q output=%q",
				tt.input,
				decrypted,
			)
		}
	}
}
