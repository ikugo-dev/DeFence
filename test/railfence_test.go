package test

import (
	"testing"

	a "github.com/ikugo-dev/DeFence/algorithms"
)

type RailFenceTest struct {
	input  string
	rails  int
	output string
}

func TestEncryptRailfence(t *testing.T) {
	testDataEncrypt := []RailFenceTest{
		{"attack at once", 2, "atc toctaka ne"},
		{"GeeksforGeeks ", 3, "GsGsekfrek eoe"},
		{"defend the east wall", 3, "dnhaweedtees alf  tl"},
	}
	for _, testData := range testDataEncrypt {
		result := a.EncryptRailfence([]byte(testData.input), testData.rails)
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
		result := a.DecryptRailfence([]byte(testData.input), testData.rails)
		if string(result) != testData.output {
			t.Errorf(`algorithms.EncryptRailfence("%s", %d) = "%s", want "%s"`,
				testData.input, testData.rails, result, testData.output)
		}
	}
}
