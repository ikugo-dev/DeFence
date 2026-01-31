package test

import (
	"encoding/binary"
	"testing"

	a "github.com/ikugo-dev/DeFence/algorithms"
)

type RailFenceTest struct {
	input  string
	rails  int
	output string
}

func toByteArray(i int) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(i))
	return arr
}

func TestEncryptRailfence(t *testing.T) {
	testDataEncrypt := []RailFenceTest{
		{"attack at once", 2, "atc toctaka ne"},
		{"GeeksforGeeks ", 3, "GsGsekfrek eoe"},
		{"defend the east wall", 3, "dnhaweedtees alf  tl"},
	}
	for _, testData := range testDataEncrypt {
		result, _ := a.EncryptFileData([]byte(testData.input), toByteArray(testData.rails), "Railfence")
		if string(result) != testData.output {
			t.Errorf(`algorithms.EncryptRailfence("%s", %d) = "%s", want "%s"`,
				testData.input, testData.rails, result, testData.output)
		}
	}
}

func TestDecryptRailfence(t *testing.T) {
	testDataDecrypt := []RailFenceTest{
		{"GsGsekfrek eoe", 3, "GeeksforGeeks "},
		{"atc toctaka ne", 2, "attack at once"},
		{"dnhaweedtees alf  tl", 3, "defend the east wall"},
	}
	for _, testData := range testDataDecrypt {
		result, _ := a.DecryptRawData([]byte(testData.input), toByteArray(testData.rails), "Railfence")
		if string(result) != testData.output {
			t.Errorf(`algorithms.EncryptRailfence("%s", %d) = "%s", want "%s"`,
				testData.input, testData.rails, result, testData.output)
		}
	}
}
