package fs

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func MkdirOrIgnore(dirTarget string) error {
	if _, err := os.Stat(dirTarget); os.IsNotExist(err) {
		err = os.Mkdir(dirTarget, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
func ScanDirByExtension(path string, ext string, recursive bool) ([]string, error) {
	entries, err := scanDir(path, recursive)
	if err != nil {
		log.Println(err)
	}

	entries = filterByExtension(entries, ext)

	return entries, nil
}
func scanDir(path string, recursive bool) ([]string, error) {
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

func filterByExtension(paths []string, ext string) []string {
	result := make([]string, 0)

	for _, p := range paths {
		ok := strings.HasSuffix(p, ext)
		if !ok {
			continue
		}

		result = append(result, p)
	}

	return result
}
