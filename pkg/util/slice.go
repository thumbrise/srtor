package util

import (
	"errors"
	"strings"
)

var (
	ErrEmptyValuePassed = errors.New("empty v passed")
	ErrEmptySizePassed  = errors.New("empty size passed")
)

func SliceSplit[T any](vv []T, chunks int) ([][]T, error) {
	chunkSize := calculateChunkSize(len(vv), chunks)
	r, err := sliceSplitBySize(vv, chunkSize)

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

func sliceSplitBySize[T any](vv []T, size int) ([][]T, error) {
	result := make([][]T, 0)

	if len(vv) == 0 {
		return result, ErrEmptyValuePassed
	}

	if size == 0 {
		return result, ErrEmptySizePassed
	}

	for {
		if len(vv) == 0 {
			break
		}

		// prevent slicing beyond
		if len(vv) < size {
			size = len(vv)
		}

		result = append(result, vv[0:size])
		vv = vv[size:]
	}

	return result, nil
}
func SliceFilter(vv []string, filter func(v string) bool) []string {
	result := make([]string, 0)

	for _, v := range vv {
		ok := filter(v)
		if !ok {
			continue
		}

		result = append(result, v)
	}

	return result
}
func SliceFilterBySuffix(vv []string, ext string) []string {
	result := make([]string, 0)

	for _, v := range vv {
		ok := strings.HasSuffix(v, ext)
		if !ok {
			continue
		}

		result = append(result, v)
	}

	return result
}
func SliceFilterByContains(vv []string, str string) []string {
	return SliceFilter(vv, func(v string) bool {
		return !strings.Contains(v, str)
	})
}
