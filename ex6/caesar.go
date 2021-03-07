package ex6

import (
	"strings"
	"unicode"
)

var alphabet = [26]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var alphabetMap = map[rune]int{'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7, 'i': 8, 'j': 9, 'k': 10, 'l': 11, 'm': 12, 'n': 13, 'o': 14, 'p': 15, 'q': 16, 'r': 17, 's': 18, 't': 19, 'u': 20, 'v': 21, 'w': 22, 'x': 23, 'y': 24, 'z': 25}

func caesarCipher(s string, k int) string {
	k = k % len(alphabet)
	shifted := append(alphabet[k:], alphabet[:k]...)
	var b strings.Builder
	for _, r := range s {
		i, ok := alphabetMap[unicode.ToLower(r)]
		if !ok {
			_, err := b.WriteRune(r)
			if err != nil {
				panic(err)
			}
			continue
		}
		encR := shifted[i]
		if unicode.IsUpper(r) {
			encR = unicode.ToUpper(encR)
		}
		_, err := b.WriteRune(encR)
		if err != nil {
			panic(err)
		}
	}
	return b.String()
}
