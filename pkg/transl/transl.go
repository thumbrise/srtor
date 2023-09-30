package transl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"srtor/pkg/util"
)

const EnvTranslateDebug = "TRANSLATE_DEBUG"

func Translate(source, sourceLang, targetLang string) (string, error) {
	if util.EnvGetBool(EnvTranslateDebug) {
		return source, nil
	}

	var text []string
	var result []interface{}
	u := getUrl(sourceLang, targetLang, source)
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

func getUrl(sourceLang string, targetLang string, query string) string {
	u := url.URL{
		Scheme:      "https",
		Opaque:      "",
		User:        nil,
		Host:        "translate.googleapis.com",
		Path:        "translate_a/single",
		RawPath:     "",
		OmitHost:    false,
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
	q := url.Values{
		"client": {"gtx"},
		"dt":     {"t"},
		"sl":     {sourceLang},
		"tl":     {targetLang},
		"q":      {query},
	}
	u.RawQuery = q.Encode()
	return u.String()
}
