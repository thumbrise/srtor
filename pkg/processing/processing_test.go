package processing

import (
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
			dir, err := filepath.Abs("./testdata/srtor")
			if err != nil {
				t.Error(err)
			}
			file := filepath.Join(dir, name)

			bytes, err := os.ReadFile(file)
			if err != nil {
				t.Errorf("Readeng file error\n%v", err)
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
	directory, err := filepath.Abs("./testdata/srtor")
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll(directory)
	if err != nil {
		log.Fatal(err)
	}
}
func process() {
	directory, err := filepath.Abs("./testdata")
	if err != nil {
		log.Fatal(err)
	}

	files, err := fs.ScanDirByExtension(directory, "srt", false)
	if err != nil {
		log.Fatal(err)
	}

	targetDirName := "srtor"

	processor := NewProcessor("en", "ru", targetDirName)
	processor.Process(files)
}
