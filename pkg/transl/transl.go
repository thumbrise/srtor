package transl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"srtor/pkg/util"
	"strings"
)

const EnvTranslateDebug = "TRANSLATE_DEBUG"

func Translate(source, sourceLang, targetLang string) (string, error) {
	if util.EnvGetBool(EnvTranslateDebug) {
		return source, nil
	}

	var text []string
	var result []interface{}
	encodedSource := url.QueryEscape(source)
	u := fmt.Sprintf(
		"https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		sourceLang,
		targetLang,
		encodedSource,
	)
	r, err := http.Get(u)
	if err != nil {
		return "err", errors.New("error getting translate.googleapis.com")
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "err", errors.New("error reading response body")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("\nstatus: %v\nresponse:\n%v", r.Status, string(body))
		return "err", errors.New("error unmarshaling data")
	}

	if len(result) > 0 {
		inner := result[0]
		for _, slice := range inner.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		cText := strings.Join(text, "")

		return cText, nil
	} else {
		return "err", errors.New("no translated data in responce")
	}
}
