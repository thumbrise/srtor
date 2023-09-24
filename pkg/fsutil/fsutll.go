package fsutil

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"srtor/pkg/util"
	"strings"
)

func ScanDirByExtension(path string, ext string, recursive bool) ([]string, error) {
	entries := make([]string, 0)
	var err error

	if recursive {
		entries, err = scanDirRecursively(path)
	} else {
		entries, err = scanDir(path)
	}
	if err != nil {
		log.Println(err)
	}

	entries = filterByExtension(entries, ext)

	return entries, nil
}

func scanDirRecursively(path string) ([]string, error) {
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

func scanDir(path string) ([]string, error) {
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

func ReadFileAsString(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func WriteFile(text string, path string) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	bytes := util.ToUTF8FixedBytes(text)

	err = os.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
func SwapFiles(a, b string) error {
	_, aStatErr := os.Stat(a)
	_, bStatErr := os.Stat(b)
	if os.IsNotExist(aStatErr) {
		log.Fatal("OK")
		return bStatErr
	}
	if os.IsNotExist(bStatErr) {
		return bStatErr
	}

	aTemp := a + ".temp"

	err := os.Rename(a, aTemp)
	if err != nil {
		return err
	}

	err = os.Rename(b, a)
	if err != nil {
		return err
	}

	err = os.Rename(aTemp, b)
	if err != nil {
		return err
	}

	return nil
}
