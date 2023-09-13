package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const targetDirName = "subtrans"

func main() {
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
	//source, err := readFile(path)
	//if err != nil {
	//	return err
	//}
	//target, err := trans.Translate(source, "en", "ru")
	//if err != nil {
	//	return err
	//}

	//fmt.Printf("From %s to %s", source, target)
	//tokens := parse.Tokenize(s)
	//tokens = removeDots(tokens)
	//tokens = mapSlice(tokens, func(item parse.Token) parse.Token {
	//	item.Text = strings.ReplaceAll(item.Text, ".", "")
	//	item.Text = strings.ToLower(item.Text)
	//	return item
	//})
	//for _, t := range tokens {
	//	fmt.Println(t.Text)
	//}
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

func fixTimeBounds(s string) string {
	r := regexp.MustCompile(`(\d\d:\d\d:\d\d)(,)?(\d\d\d)`)
	template := "$1,$3"
	return r.ReplaceAllString(s, template)
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
