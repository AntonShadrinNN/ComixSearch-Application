package stem

import (
	"strings"

	"github.com/kljensen/snowball"
)

type Stem struct {
	language  string
	stopWords bool
}

func NewStem(language string, stopWords bool) Stem {
	return Stem{
		language:  language,
		stopWords: stopWords,
	}
}

func (s Stem) Normalize(data string) (string, error) {

	tokens := Prepocess(data)

	for i := 0; i < len(tokens); i++ {
		stemmed, err := snowball.Stem(tokens[i], s.language, s.stopWords)
		tokens[i] = stemmed
		if err != nil {
			return "", err
		}
	}

	stemmedData := strings.Join(tokens, " ")

	return stemmedData, nil
}
