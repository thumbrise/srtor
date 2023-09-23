package main

import (
	"os"
	"path/filepath"
	"srtor/pkg/fs"
	"srtor/pkg/interaction"
	"srtor/pkg/processing"
	"strings"
)

const targetDirName = "srtor"

func main() {
	directory, err := interaction.AskDirectory()
	if err != nil {
		panic(err)
	}

	files, err := fs.ScanDir(directory, func(f os.DirEntry) bool {
		return strings.HasSuffix(f.Name(), ".srt")
	})
	if err != nil {
		panic(err)
	}

	err = fs.MkdirOrIgnore(filepath.Join(directory, targetDirName))
	if err != nil {
		panic(err)
	}

	languageSource := interaction.AskLanguageSource()
	languageTarget := interaction.AskLanguageTarget()

	processor := processing.NewProcessor()
	processor.LangSource = languageSource
	processor.LangTarget = languageTarget
	processor.Process(files)

	interaction.Bye(len(files), directory)
}
