package interaction

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"srtor/pkg/util"
	"strings"
)

const languagesLink = "https://cloud.google.com/translate/docs/languages"
const defaultLanguageSource = "en"
const defaultLanguageTarget = "ru"

func AskLanguageSource() string {
	return AskLanguage("source", defaultLanguageSource)
}

func AskLanguageTarget() string {
	return AskLanguage("target", defaultLanguageTarget)
}

func AskLanguage(label string, defaultValue string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Click %s for look language domain abbreviation!!\n", languagesLink)
	fmt.Printf("Type %s language abbreviation. Empty for default %s\n", label, defaultValue)

	scanner.Scan()
	result := scanner.Text()
	if result == "" {
		result = defaultValue
	}

	return result
}

func AskDirectory() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Type directory absolute path which contains srt files")

	pathFromWd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Or empty for default %s\n", pathFromWd)

	scanner.Scan()
	pathFromConsole, err := util.CanonizePath(scanner.Text())
	if err != nil {
		return "", err
	}
	result := ""
	if pathFromConsole != "" {
		result = pathFromConsole
	} else {
		result = pathFromWd
	}

	return result, nil
}
func Bye(filesCount int, filesDir string) {
	fmt.Println()
	message := ""
	if filesCount > 0 {
		message = "Success"
	} else {
		abs, err := filepath.Abs(filesDir)
		if err != nil {
			log.Println(err)
		}
		message = fmt.Sprintf("No .srt files in %s", abs)
	}

	result := fmt.Sprintf("%s\n%s", message, "Press ENTER for exit")

	s := bufio.NewScanner(os.Stdin)
	fmt.Println(result)
	s.Scan()
}

func AskRecursive() bool {
	return askBool("Recursive processing Y/n (n)")
}

func AskArchive() bool {
	return askBool("Archive original subtitles Y/n (n)")
}

func AskReplace() bool {
	return askBool("Replace original subtitles by translated Y/n (n)")
}

func askBool(message string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(message)

	scanner.Scan()
	t := scanner.Text()

	return strings.ToLower(t) == "y"
}
