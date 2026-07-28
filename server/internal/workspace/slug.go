package workspace

import (
	"regexp"
	"strings"
)

// nonSlugChars matches any run of characters that are NOT lowercase letters or digits
var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)


func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonSlugChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
