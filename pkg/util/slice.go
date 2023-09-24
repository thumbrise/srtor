package util

import (
	"errors"
)

var (
	ErrEmptyValuePassed = errors.New("empty v passed")
	ErrEmptySizePassed  = errors.New("empty size passed")
)

func SliceSplit[T any](slice []T, chunks int) ([][]T, error) {
	chunkSize := calculateChunkSize(len(slice), chunks)
	r, err := sliceSplitBySize(slice, chunkSize)

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

func sliceSplitBySize[T any](v []T, size int) ([][]T, error) {
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
func SliceFilter[T any](vv []T, filter func(v T) bool) []T {
	result := make([]T, 0)

	for _, v := range vv {
		ok := filter(v)
		if !ok {
			continue
		}

		result = append(result, v)
	}

	return result
}
