package normalizer

import (
	"github.com/kljensen/snowball"
)

const (
	Language  = "english"
	StopWords = true
)

func Normalize(data string) ([]string, error) {

	tokens := Prepocess(data)
	for i := 0; i < len(tokens); i++ {
		stemmed, err := snowball.Stem(tokens[i], Language, StopWords)
		tokens[i] = stemmed
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}
