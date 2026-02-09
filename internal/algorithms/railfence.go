package algorithms

import "fmt"

func encryptRailfence(text []byte, byteKey []byte) ([]byte, error) {
	key, err := bytesToint(byteKey)
	if err != nil || key <= 1 {
		return nil, fmt.Errorf("invalid key")
	}
	if key > len(text) {
		key = len(text)
	}
	// create the matrix to cipher plain text
	// key = rows , length(text) = columns
	rail := make([][]byte, key)
	filled := make([][]bool, key)
	for i := range rail {
		rail[i] = make([]byte, len(text))
		// filling the rail matrix to distinguish filled
		// spaces from blank ones
		filled[i] = make([]bool, len(text))
	}

	// to find the direction
	dir_down := false
	row, col := 0, 0

	for i := range text {
		// check the direction of flow
		// reverse the direction if we've just
		// filled the top or bottom rail
		if row == 0 || row == key-1 {
			dir_down = !dir_down
		}

		// fill the corresponding alphabet
		rail[row][col] = byte(text[i])
		filled[row][col] = true
		col++

		// find the next row using direction flag
		if dir_down {
			row++
		} else {
			row--
		}
	}

	var result []byte
	for i := range rail {
		for j := range rail[i] {
			if filled[i][j] {
				result = append(result, rail[i][j])
			}
		}
	}
	return result, nil
}

func decryptRailfence(cipher []byte, byteKey []byte) ([]byte, error) {
	key, err := bytesToint(byteKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key size")
	}

	rail := make([][]byte, key)
	marked := make([][]bool, key)
	for i := range rail {
		rail[i] = make([]byte, len(cipher))
		marked[i] = make([]bool, len(cipher))
	}

	dir_down := false
	row, col := 0, 0

	for range cipher {
		// check the direction of flow
		if row == 0 {
			dir_down = true
		}
		if row == key-1 {
			dir_down = false
		}

		marked[row][col] = true
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
			if marked[i][j] {
				rail[i][j] = cipher[index]
				index++
			}
		}
	}

	var result []byte
	row, col = 0, 0
	for range cipher {
		// check the direction of flow
		if row == 0 {
			dir_down = true
		}
		if row == key-1 {
			dir_down = false
		}

		// place the marker
		if marked[row][col] {
			result = append(result, rail[row][col])
			col++
		}

		// find the next row using direction flag
		if dir_down {
			row++
		} else {
			row--
		}
	}
	return result, nil
}
