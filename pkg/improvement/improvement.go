package improvement

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	Number int
	Start  time.Time
	End    time.Time
	Text   string
}

func Tokenize(s string) []Token {
	tokensRaw := strings.Split(s, "\n\n")
	tokensResult := make([]Token, 0)
	for _, raw := range tokensRaw {
		parts := strings.Split(raw, "\n")
		numberRaw := parts[0]
		number, _ := strconv.Atoi(numberRaw)
		timeRaw := parts[1]
		timeParts := strings.Split(timeRaw, " --> ")
		layout := "15:04:05,000"
		start, _ := time.Parse(layout, timeParts[0])
		end, _ := time.Parse(layout, timeParts[1])
		text := parts[2]
		t := Token{
			Number: number,
			Start:  start,
			End:    end,
			Text:   text,
		}
		tokensResult = append(tokensResult, t)
	}
	return tokensResult
}
func FixTimeBounds(s []byte) []byte {
	r := regexp.MustCompile("(\\d\\d:\\d\\d:\\d\\d)(,)?(\\d\\d\\d)")
	template := "$1,$3"
	result := r.ReplaceAllString(string(s), template)
	return []byte(result)
}
