package normalizer

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Normalize standardizes text for fair comparison.
//
// Spec:
// - lowercase
// - Unicode normalization (NFKC)
// - remove all punctuation (unicode.IsPunct)
// - normalize whitespace (tabs/newlines/multi-space -> single space)
func Normalize(input string) string {
	if input == "" {
		return ""
	}

	// Unicode normalization first, so punctuation classification stable.
	s := norm.NFKC.String(input)
	s = strings.ToLower(s)

	// Remove punctuation only. Keep letters/digits/symbols.
	s = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, s)

	// Collapse all whitespace runs (includes tabs/newlines) to single spaces.
	return strings.Join(strings.Fields(s), " ")
}
