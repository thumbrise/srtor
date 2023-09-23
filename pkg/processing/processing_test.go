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
			pathActual, err := filepath.Abs("./testdata/sources/srtor")
			pathActual = filepath.Join(pathActual, name)
			if err != nil {
				t.Error(err)
			}

			pathExpected, err := filepath.Abs("./testdata/expected")
			if err != nil {
				t.Error(err)
			}
			pathExpected = filepath.Join(pathExpected, name)

			hashActual := hash(pathActual)
			hashExpected := hash(pathExpected)

			if hashExpected != hashActual {
				t.Errorf("Hashes not equals")
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
