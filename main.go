package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"subtrans/pkg/trans"
	"unicode/utf8"
)

const targetDirName = "subtrans"

func main() {
	file, _ := os.ReadFile("E:\\go-projects\\subtrans\\testdata\\sources\\subtrans\\01 - Understanding Parallel Computing.srt")
	if !utf8.Valid(file) {
		panic("Invalit")
	}
	d := pollSourceDirectory()

	ff, err := scanSourceDirectory(d)
	if err != nil {
		panic(err)
	}

	err = ensureTargetDirectory(d)
	if err != nil {
		panic(err)
	}

	for _, f := range ff {
		err = processFile(f)
		if err != nil {
			panic(err)
		}
	}
}
func processFile(path string) error {
	source, err := readFile(path)
	if err != nil {
		return err
	}
	target, err := trans.Translate(source, "en", "ru")
	if err != nil {
		return err
	}

	sourceDir := filepath.Dir(path)
	sourceName := filepath.Base(path)
	targetPath := filepath.Join(sourceDir, targetDirName, sourceName)
	targetBytes := []byte(target)
	unicodeReplacement := []byte{0xef, 0xbf, 0xbd}
	targetBytes = bytes.ToValidUTF8(targetBytes, unicodeReplacement)

	err = os.WriteFile(targetPath, targetBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
func ensureTargetDirectory(path string) error {
	dirTarget := filepath.Join(path, targetDirName)
	if _, err := os.Stat(dirTarget); os.IsNotExist(err) {
		err = os.Mkdir(dirTarget, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := string(b)
	return s, nil
}
func mapSlice[T any](s []T, callback func(item T) T) []T {
	r := append(s)
	for i, v := range r {
		r[i] = callback(v)
	}
	return r
}

func fixTimeBounds(s []byte) []byte {
	r := regexp.MustCompile("(\\d\\d:\\d\\d:\\d\\d)(,)?(\\d\\d\\d)")
	template := "$1,$3"
	result := r.ReplaceAllString(string(s), template)
	return []byte(result)
}

func pollSourceDirectory() string {
	s := bufio.NewScanner(os.Stdin)
	fmt.Println("Type directory absolute path which contains srt files")
	s.Scan()
	return s.Text()
}
func scanSourceDirectory(path string) ([]string, error) {
	result := make([]string, 0)
	d, err := os.ReadDir(path)
	if err != nil {
		return result, err
	}
	for _, f := range d {
		name := f.Name()
		if !strings.HasSuffix(name, ".srt") {
			continue
		}
		result = append(result, filepath.Join(path, name))
	}
	return result, nil
}
