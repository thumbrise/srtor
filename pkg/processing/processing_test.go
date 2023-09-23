package processing

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"srtor/pkg/fs"
	"testing"
)

func TestProcessor_Process(t *testing.T) {
	removeTempDir()

	names := []string{"0.srt", "1.srt"}

	process()

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			dir, err := filepath.Abs("./testdata/sources/srtor")
			if err != nil {
				t.Error(err)
			}
			file := filepath.Join(dir, name)

			bytes, err := os.ReadFile(file)
			if err != nil {
				t.Error(err)
			}

			text := string(bytes)
			if text == "" {
				t.Errorf("Result file is empty")
			}
		})
	}

	removeTempDir()
}

func removeTempDir() {
	directory, err := filepath.Abs("./testdata/sources/srtor")
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(directory)
	if err != nil {
		log.Fatal(err)
	}
}
func process() {
	directory, err := filepath.Abs("./testdata/sources")
	if err != nil {
		panic(err)
	}

	files, err := fs.ScanDirByExtension(directory, "srt")
	if err != nil {
		panic(err)
	}

	targetDirName := "srtor"
	destination := filepath.Join(directory, targetDirName)
	err = fs.MkdirOrIgnore(destination)
	if err != nil {
		panic(err)
	}

	processor := NewProcessor("en", "ru", destination)
	processor.Process(files)
}

func hash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
