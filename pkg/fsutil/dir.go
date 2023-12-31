package fsutil

import (
	"io/fs"
	"os"
	"path/filepath"
)

func DirScan(path string, recursive bool) ([]string, error) {
	entries := make([]string, 0)
	var err error

	if recursive {
		entries, err = DirScanRecursively(path)
	} else {
		entries, err = DirScanPlain(path)
	}
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func DirScanRecursively(path string) ([]string, error) {
	result := make([]string, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		result = append(result, path)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil

}

func DirScanPlain(path string) ([]string, error) {
	result := make([]string, 0)

	d, err := os.ReadDir(path)
	if err != nil {
		return result, err
	}

	for _, f := range d {
		result = append(result, filepath.Join(path, f.Name()))
	}

	return result, nil
}
