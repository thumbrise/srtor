package fs

import (
	"os"
	"path/filepath"
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
func ScanDir(path string, filter func(f os.DirEntry) bool) ([]string, error) {
	result := make([]string, 0)
	d, err := os.ReadDir(path)
	if err != nil {
		return result, err
	}
	for _, f := range d {
		ok := filter(f)
		if ok {
			result = append(result, filepath.Join(path, f.Name()))
		}
	}
	return result, nil
}
