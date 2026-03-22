package utils

import (
	"regexp"
	"strings"
)

var SlugRegexSanitizer = regexp.MustCompile(`[^\p{L}\p{N}_-]+`)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = SlugRegexSanitizer.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
