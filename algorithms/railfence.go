package algorithms

import (
	"log"
)

func EncryptRailfence(text []byte, byteKey []byte) []byte {
	key, err := bytesToint(byteKey)
	if err != nil || key <= 1 {
		log.Fatal("invalid key")
	}
	if key > len(text) {
		key = len(text)
	}
	// create the matrix to cipher plain text
	// key = rows , length(text) = columns
	rail := make([][]byte, key)
	for i := range rail {
		rail[i] = make([]byte, len(text))
	}
	// filling the rail matrix to distinguish filled
	// spaces from blank ones
	for i := range key {
		for j := range text {
			rail[i][j] = '\n'
		}
	}

	// to find the direction
	dir_down := false
	row := 0
	col := 0

	for i := range text {
		// check the direction of flow
		// reverse the direction if we've just
		// filled the top or bottom rail
		if row == 0 || row == key-1 {
			dir_down = !dir_down
		}

		// fill the corresponding alphabet
		rail[row][col] = byte(text[i])
		col++

		// find the next row using direction flag
		if dir_down {
			row++
		} else {
			row--
		}
	}

	var result []byte
	for _, row := range rail {
		for _, r := range row {
			if r != '\n' {
				result = append(result, r)
			}
		}
	}
	return result
}

func DecryptRailfence(cipher []byte, byteKey []byte) []byte {
	key, err := bytesToint(byteKey)
	if err != nil {
		log.Fatal("invalid key size")
	}

	rail := make([][]byte, key)
	for i := range rail {
		rail[i] = make([]byte, len(cipher))
	}
	for i := range key {
		for j := range cipher {
			rail[i][j] = '\n'
		}
	}

	dir_down := false
	row := 0
	col := 0

	for range cipher {
		// check the direction of flow
		if row == 0 {
			dir_down = true
		}
		if row == key-1 {
			dir_down = false
		}

		rail[row][col] = '*'
		col++

		if dir_down {
			row++
		} else {
			row--
		}
	}

	index := 0
	for i := range key {
		for j := range cipher {
			if rail[i][j] == '*' && index < len(cipher) {
				rail[i][j] = cipher[index]
				index++
			}
		}
	}

	var result []byte
	row = 0
	col = 0
	for range cipher {
		// check the direction of flow
		if row == 0 {
			dir_down = true
		}
		if row == key-1 {
			dir_down = false
		}

		// place the marker
		if rail[row][col] != '*' {
			result = append(result, (rail[row][col]))
			col++
		}

		// find the next row using direction flag
		if dir_down {
			row++
		} else {
			row--
		}
	}
	return result
}
