package main

import (
	"path/filepath"
	"srtor/pkg/fs"
	"srtor/pkg/interaction"
	"srtor/pkg/processing"
)

const targetDirName = "srtor"
const translatableExtension = "srt"

func main() {
	directory, err := interaction.AskDirectory()
	if err != nil {
		panic(err)
	}

	files, err := fs.ScanDirByExtension(directory, translatableExtension)
	if err != nil {
		panic(err)
	}

	destination := filepath.Join(directory, targetDirName)
	err = fs.MkdirOrIgnore(destination)
	if err != nil {
		panic(err)
	}

	languageSource := interaction.AskLanguageSource()
	languageTarget := interaction.AskLanguageTarget()

	processor := processing.NewProcessor(languageSource, languageTarget, destination)
	processor.Process(files)

	interaction.Bye(len(files), directory)
}
