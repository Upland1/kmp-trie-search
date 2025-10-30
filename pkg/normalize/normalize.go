package normalize

import (
	"regexp"
	"strings"
)

func NormalizeText(book string) string {
	lowerText := strings.ToLower(book)
	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	cleanedText := reg.ReplaceAllString(lowerText, "")

	return cleanedText
}
