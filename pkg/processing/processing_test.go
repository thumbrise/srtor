package processing

import (
	"log"
	"os"
	"path/filepath"
	"srtor/pkg/fsutil"
	"srtor/pkg/transl"
	"srtor/pkg/util"
	"testing"
)

const resultDirName = "srtor-result"

func TestProcessor_Process(t *testing.T) {
	removeTempDir()

	names := []string{"0.srt", "1.srt"}

	process()

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			dir, err := filepath.Abs(filepath.Join("./testdata", resultDirName))
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
	directory, err := filepath.Abs(filepath.Join("./testdata", resultDirName))
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

	files, err := fsutil.DirScan(directory, false)
	if err != nil {
		log.Fatal(err)
	}

	files = util.SliceFilterBySuffix(files, "srt")
	processor := NewProcessor(transl.NewEnvBasedTranslator(), "en", "ru", resultDirName)
	processor.Process(files)
}
