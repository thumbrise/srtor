package util

import "errors"

func MapSlice[T any](slice []T, callback func(item T) T) []T {
	r := append(slice)
	for i, v := range r {
		r[i] = callback(v)
	}
	return r
}

func ChunkSlice[T any](slice []T, chunkSize int) ([][]T, error) {
	var chunks [][]T

	if len(slice) == 0 {
		return chunks, errors.New("empty slice passed")
	}

	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks, nil
}
