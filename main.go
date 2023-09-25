package main

import (
	"log"
	"srtor/pkg/fsutil"
	"srtor/pkg/interaction"
	"srtor/pkg/processing"
	"srtor/pkg/util"
)

const resultDirName = "srtor-result"
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
	needReplace := interaction.AskReplace()

	files := scanDir(directory, needRecursive)
	files = filterNonTranslatable(files)
	files = filterResultDirs(files)

	processing.
		NewProcessor(languageSource, languageTarget, resultDirName).
		WithReplace(needReplace).
		Process(files)

	interaction.Bye(len(files), directory)
}

func configureLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// need for prevent recursive creating
func filterResultDirs(files []string) []string {
	return util.SliceFilterByContains(files, resultDirName)
}

func filterNonTranslatable(files []string) []string {
	return util.SliceFilterBySuffix(files, translatableExtension)
}

func scanDir(path string, recursive bool) []string {
	result, err := fsutil.DirScan(path, recursive)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
