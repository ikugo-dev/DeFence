package metadata

import (
	"encoding/binary"
	"encoding/json"
	"os"

	"github.com/ikugo-dev/DeFence/logger"
)

type Metadata struct {
	FileName            string   `json:"fileName"`
	FileSize            int64    `json:"fileSize"`
	CreationDateTime    string   `json:"creationDateTime"`
	EncryptionAlgorithm string   `json:"encryptionAlgorithm"`
	HashingAlgorithm    string   `json:"hashingAlgorithm"`
	HashingResult       [24]byte `json:"hashingResult"`
}

func Create(fileName, algorithm, hashingAlgorithm string, hashingResult [24]byte) []byte {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		logger.Log("Could not read file %s: %v", fileName, err)
		return nil
	}

	metadata := &Metadata{
		FileName:            fileName,
		FileSize:            fileInfo.Size(),
		CreationDateTime:    fileInfo.ModTime().String(),
		EncryptionAlgorithm: algorithm,
		HashingAlgorithm:    hashingAlgorithm,
		HashingResult:       hashingResult,
	}
	jsonContent, err := json.Marshal(metadata)

	if err != nil {
		logger.Log("Could not encode metadata for file %s: %v", fileName, err)
		return nil
	}

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(jsonContent)))
	return append(lenBuf, jsonContent...)
}

func Read(data []byte) Metadata {
	metadataLen := binary.BigEndian.Uint32(data[:4])
	metadataBytes := data[4 : 4+metadataLen]

	var metadata Metadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		logger.Log("Failed to decode metadata: %v", err)
		return Metadata{}
	}
	return metadata
}
