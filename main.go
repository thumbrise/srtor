package main

import (
	"log"
	"srtor/pkg/fs"
	"srtor/pkg/interaction"
	"srtor/pkg/processing"
)

const resultDirName = "srtor"
const translatableExtension = "srt"

func main() {
	configureLogger()

	directory, err := interaction.AskDirectory()
	if err != nil {
		log.Fatal(err)
	}

	languageSource := interaction.AskLanguageSource()
	languageTarget := interaction.AskLanguageTarget()

	needRecursive := interaction.AskRecursive()
	files, err := fs.ScanDirByExtension(directory, translatableExtension, needRecursive)
	if err != nil {
		log.Fatal(err)
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

func configureLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
