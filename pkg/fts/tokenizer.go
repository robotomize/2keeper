package fts

import (
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	return strings.FieldsFunc(
		text, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		},
	)
}
