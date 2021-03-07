package ex6

import (
	"io"
	"strings"
	"unicode"
)

func camelcase(s string) int32 {
	var wordCount int32 = 1
	r := strings.NewReader(s)

	for {
		rune, _, err := r.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		if unicode.IsUpper(rune) {
			wordCount += 1
		}
	}

	return wordCount
}
