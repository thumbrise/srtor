package main

import (
	"srtor/pkg/fs"
	"srtor/pkg/interaction"
	"srtor/pkg/processing"
)

const resultDirName = "srtor"
const translatableExtension = "srt"

func main() {
	directory, err := interaction.AskDirectory()
	if err != nil {
		panic(err)
	}

	languageSource := interaction.AskLanguageSource()
	languageTarget := interaction.AskLanguageTarget()

	needRecursive := interaction.AskRecursive()
	files, err := fs.ScanDirByExtension(directory, translatableExtension, needRecursive)
	if err != nil {
		panic(err)
	}

	needReplace := interaction.AskReplace()
	needArchive := false
	if needReplace {
		needArchive = interaction.AskArchive()
	}

	processing.
		NewProcessor(languageSource, languageTarget, resultDirName).
		WithReplace(needReplace).
		WithArchive(needArchive).
		Process(files)

	interaction.Bye(len(files), directory)
}
