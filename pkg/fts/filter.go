package fts

import (
	"embed"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
	"gopkg.in/yaml.v3"

	"github.com/robotomize/2keeper/pkg/slice"
)

type FilterFunc func([]string) []string
type TokenizeFunc func(string) []string

//go:embed stopwords
var stopwordsFS embed.FS

var stopwords map[string]struct{}

func init() {
	dir, err := stopwordsFS.ReadDir(".")
	if err != nil {
		panic("embed fs ReadDir")
	}

	if stopwords == nil {
		stopwords = make(map[string]struct{})
	}

	type stopwordList struct {
		Words []string `yaml:"words"`
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			file, err := os.Open(entry.Name())
			if err != nil {
				log.Println("os.Open ", err)
				continue
			}

			var list stopwordList
			if err = yaml.NewDecoder(file).Decode(&list); err != nil {
				log.Println("yaml.NewDecoder.Decode ", err)
				continue
			}

			for _, w := range list.Words {
				stopwords[w] = struct{}{}
			}
		}
	}
}

func LowercaseFilter(tokens []string) []string {
	return slice.Map(
		tokens, func(t string) string {
			return strings.ToLower(t)
		},
	)
}

func StopwordsFilter(tokens []string) []string {
	return slice.Filter(
		tokens, func(t string) bool {
			_, ok := stopwords[t]
			return !ok
		},
	)
}

func StemmerFilter(tokens []string) []string {
	output := make([]string, 0, len(tokens))
OuterLoop:
	for _, t := range tokens {
		var lang string
		letter := []rune(t)[0]

		switch {
		case unicode.Is(unicode.Latin, letter):
			lang = "english"
		case unicode.Is(unicode.Cyrillic, letter):
			lang = "russian"
		case unicode.Is(unicode.Number, letter):
			output = append(output, t)
			continue OuterLoop
		default:
			continue OuterLoop
		}

		stemmed, err := snowball.Stem(t, lang, true)
		if err == nil {
			output = append(output, stemmed)
		}
	}

	return output
}

func StopwordFilterFunc() FilterFunc {
	return func(tokens []string) []string {
		return slice.Filter(
			tokens, func(t string) bool {
				_, ok := stopwords[t]
				return !ok
			},
		)
	}
}

func LowercaseFunc() FilterFunc {
	return func(tokens []string) []string {
		return slice.Map(
			tokens, func(t string) string {
				return strings.ToLower(t)
			},
		)
	}
}
