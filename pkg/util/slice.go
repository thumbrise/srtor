package util

import (
	"errors"
)

var (
	ErrEmptyValuePassed = errors.New("empty v passed")
	ErrEmptySizePassed  = errors.New("empty size passed")
)

func MapSlice[T any](slice []T, callback func(item T) T) []T {
	r := append(slice)
	for i, v := range r {
		r[i] = callback(v)
	}

	return r
}

func SplitSlice[T any](slice []T, chunks int) ([][]T, error) {
	chunkSize := calculateChunkSize(len(slice), chunks)
	r, err := chunkSliceBySize(slice, chunkSize)

	return r, err
}

func calculateChunkSize(total, divider int) int {
	// prevent negative divider
	dividerReal := Abs(divider)
	// prevent total < divider
	dividerReal = Min(dividerReal, total)
	// prevent zero division
	if dividerReal == 0 {
		return total
	}

	return total / dividerReal
}

func chunkSliceBySize[T any](v []T, size int) ([][]T, error) {
	result := make([][]T, 0)

	if len(v) == 0 {
		return result, ErrEmptyValuePassed
	}

	if size == 0 {
		return result, ErrEmptySizePassed
	}

	for {
		if len(v) == 0 {
			break
		}

		// prevent slicing beyond
		if len(v) < size {
			size = len(v)
		}

		result = append(result, v[0:size])
		v = v[size:]
	}

	return result, nil
}
