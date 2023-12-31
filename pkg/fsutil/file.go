package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"srtor/pkg/util"
	"strconv"
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
	if FileNotExists(a) {
		message := fmt.Sprintf("File not exist %s", a)
		return errors.New(message)
	}
	if FileNotExists(b) {
		message := fmt.Sprintf("File not exist %s", a)
		return errors.New(message)
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
	if FileNotExists(path) {
		return os.Create(path)
	}

	return os.Open(path)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func FileNotExists(path string) bool {
	return !FileExists(path)
}

func FileIncrementName(path string) string {
	r := regexp.MustCompile(`(.*?)(\d*)(\..*)?$`)
	matches := r.FindStringSubmatch(path)

	name, number, extension := matches[1], matches[2], matches[3]

	numberConverted, err := strconv.Atoi(number)
	if err != nil {
		const numberDefault = 2
		return fmt.Sprintf("%s%d%s", name, numberDefault, extension)
	}

	return fmt.Sprintf("%s%d%s", name, numberConverted+1, extension)
}
