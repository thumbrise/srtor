package fsutil

import (
	"os"
	"path/filepath"
	"srtor/pkg/util"
)

func FileReadAsString(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func FileWrite(text string, path string) error {
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

func FileSwap(a, b string) error {
	_, aStatErr := os.Stat(a)
	_, bStatErr := os.Stat(b)
	if os.IsNotExist(aStatErr) {
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

func FileOpenOrCreate(path string) (*os.File, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.Create(path)
	}

	return os.Open(path)
}
