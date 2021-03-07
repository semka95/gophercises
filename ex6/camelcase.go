package ex6

import (
	"unicode"
)

func camelcase(s string) int {
	var wordCount int = 1

	for _, v := range s {
		if unicode.IsUpper(v) {
			wordCount += 1
		}
	}

	return wordCount
}
