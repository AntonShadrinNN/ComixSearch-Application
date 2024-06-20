package stem

import (
	"regexp"
	"strings"
)

// Preproccess takes originText and remove punctuation and specific symbols from it. Returns proccessed text splited by any space.
func prepocess(origText string) []string {
	var tokens []string
	re := regexp.MustCompile(`\{\{.*?\}\}|[[:punct:]]`)
	temp := re.ReplaceAllString(origText, "")
	tokens = strings.Fields(temp)

	return tokens
}
