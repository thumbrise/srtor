package main

import (
	"fmt"
	"regexp"
	"strings"
	"subtrans/parse"
)

func main() {
	s := "1\n00:00:01,536 --> 00:00:02,560\nHello and welcome.\n\n2\n00:00:03,072 --> 00:00:05,888\nThe score is called multi treading in goal.\n\n3\n00:00:06,912 --> 00:00:09,984\nA thread is simply the stool or this abstraction\n\n4\n00:00:10,240 --> 00:00:12,032\nThat allows us to perform\n\n5\n00:00:12,288 --> 00:00:13,312\nParallel computation\n\n6\n00:00:14,336 --> 00:00:16,640\nComputing is the\n\n7\n00:00:16,896 --> 00:00:17,920\nOr Some people prefer"
	//orig := "1\n00:00:01536 --> 00:00:02560\nПривет и добро пожаловать.\n\n2\n00:00:03,072 --> 00:00:05,888\nСчет называется «многократным проникновением в ворота».\n\n3\n00:00:06,912 --> 00:00:09,984\nНить — это просто табуретка или эта абстракция.\n\n4\n00:00:10,240 --> 00:00:12,032\nЭто позволяет нам выполнять\n\n5\n00:00:12,288 --> 00:00:13,312\nПараллельные вычисления\n\n6\n00:00:14,336 --> 00:00:16,640\nВычисление – это\n\n7\n00:00:16896 --> 00:00:17920\nИли некоторые люди предпочитают"

	//fmt.Println(fixTimeBounds(orig))
	//to, err := trans.Translate(s, "en", "ru")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("From %s to %s", s, to)
	tokens := parse.Tokenize(s)
	//tokens = removeDots(tokens)
	tokens = mapSlice(tokens, func(item parse.Token) parse.Token {
		item.Text = strings.ReplaceAll(item.Text, ".", "")
		item.Text = strings.ToLower(item.Text)
		return item
	})
	for _, t := range tokens {
		fmt.Println(t.Text)
	}
}
func mapSlice[T any](s []T, callback func(item T) T) []T {
	r := append(s)
	for i, v := range r {
		r[i] = callback(v)
	}
	return r
}

func fixTimeBounds(s string) string {
	r := regexp.MustCompile(`(\d\d:\d\d:\d\d)(,)?(\d\d\d)`)
	template := "$1,$3"
	return r.ReplaceAllString(s, template)
}
