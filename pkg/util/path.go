package util

import (
	"path/filepath"
	"strings"
)

func CanonizePath(path string) (string, error) {
	result := strings.ReplaceAll(path, "\"", "")
	result, err := filepath.EvalSymlinks(result)
	return result, err
}
