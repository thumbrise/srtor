package transl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"srtor/pkg/util"
)

const EnvTranslateDebug = "TRANSLATE_DEBUG"
const host = "translate.googleapis.com"
const path = "translate_a/single"
const scheme = "https"

func Translate(source, sourceLang, targetLang string) (string, error) {
	if util.EnvGetBool(EnvTranslateDebug) {
		return source, nil
	}

	respBytes, err := doRequest(sourceLang, targetLang, source)
	if err != nil {
		return "", err
	}

	result, err := parse(respBytes)
	if err != nil {
		return "", err
	}

	return result, nil
}

// expected structure s[0][0][0], s[0][1][0], s[0][2][0]...
//
//	[
//		[	<- Target list of elements
//			[
//				translatedText,	<- Target field in each element
//				trash,
//				trash,
//				trash...
//			]...
//		],
//		trash,
//		trash,
//		trash...
//	]
func parse(bytes []byte) (string, error) {
	var wrapped []interface{}

	err := json.Unmarshal(bytes, &wrapped)
	if err != nil {
		return "", err
	}

	if len(wrapped) <= 0 {
		msg := fmt.Sprintf("no data in %s response", host)
		return "err", errors.New(msg)
	}

	result := ""

	unwrapped := wrapped[0].([]interface{})
	for _, tData := range unwrapped {
		tConverted := tData.([]interface{})
		text := fmt.Sprintf("%v", tConverted[0])
		result += text
	}

	return result, nil
}

func doRequest(sourceLang string, targetLang string, query string) ([]byte, error) {
	q := url.Values{
		"client": {"gtx"},
		"dt":     {"t"},
		"sl":     {sourceLang},
		"tl":     {targetLang},
		"q":      {query},
	}

	u := url.URL{
		Host:     host,
		Scheme:   scheme,
		Path:     path,
		RawQuery: q.Encode(),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		msg := fmt.Sprintf("error getting %s", host)
		return nil, errors.New(msg)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	if resp.StatusCode >= 300 {
		msg := fmt.Sprintf("bad status (%d) of %s\nresponse:\n%s", resp.StatusCode, host, string(bytes))
		return nil, errors.New(msg)
	}

	return bytes, nil
}
