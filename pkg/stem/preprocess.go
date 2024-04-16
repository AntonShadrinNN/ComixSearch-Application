package stem

import (
	"regexp"
	"strings"
)

func Prepocess(origText string) []string {
	var tokens []string
	re := regexp.MustCompile(`\{\{.*?\}\}|[[:punct:]]`)
	temp := re.ReplaceAllString(origText, "")
	tokens = strings.Fields(temp)

	return tokens
}
