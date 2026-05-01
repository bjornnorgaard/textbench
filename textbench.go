package textbench

import (
	"context"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type PipelineFunc func(string) string

func EvaluateString(a, b string) (int, error) {
	if a == "" || b == "" {
		return 0, nil
	}

	stages := []PipelineFunc{
		// Whitespace normalization.
		WhitespaceCleaning(),
		// Unicode normalization (standardize special chars/diacritics).
		UnicodeNormalizationFunc(),
		// Strip any remaining non-ASCII (post-normalization).
		NonASCIIFunc(),
		// Punctuation removal.
		PunctuationRemovalFunc(),
		// Digits -> words.
		NumberStandardizationFunc(),
		// Case normalization.
		CaseNormalization(),
		// Final whitespace normalization after transforms.
		WhitespaceCleaning(),
	}

	for _, stage := range stages {
		a = stage(a)
	}

	for _, stage := range stages {
		b = stage(b)
	}

	// Lower score better. 0 means identical after normalization.
	return wordEditDistance(strings.Fields(a), strings.Fields(b)), nil
}

func CaseNormalization() PipelineFunc {
	return func(s string) string {
		return strings.ToLower(s)
	}
}

func WhitespaceCleaning() PipelineFunc {
	return func(s string) string {
		// Fields splits on any whitespace and drops empty runs; joining collapses
		// repeated whitespace (including leading/trailing) into single spaces.
		return strings.Join(strings.Fields(s), " ")
	}
}

func NonASCIIFunc() PipelineFunc {
	return func(s string) string {
		// Strip any remaining non-ASCII runes.
		return strings.Map(func(r rune) rune {
			if r > 0x7F {
				return -1
			}
			return r
		}, s)
	}
}

func UnicodeNormalizationFunc() PipelineFunc {
	return func(s string) string {
		// 1) Compatibility normalize (fold some "special" forms to standard forms)
		// 2) Decompose, then drop combining marks (strip diacritics)
		n := norm.NFKC.String(s)
		out := strings.Map(func(r rune) rune {
			if unicode.Is(unicode.Mn, r) {
				return -1
			}
			return r
		}, norm.NFD.String(n))
		return out
	}
}

func PunctuationRemovalFunc() PipelineFunc {
	return func(s string) string {
		out := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
				return r
			}
			return -1
		}, s)
		return out
	}
}

func NumberStandardizationFunc() PipelineFunc {
	return func(s string) string {
		runes := []rune(s)
		var b strings.Builder
		b.Grow(len(s))

		for i, r := range runes {
			if r < '0' || r > '9' {
				b.WriteRune(r)
				continue
			}

			word := digitWord(r)
			prev := rune(0)
			next := rune(0)
			if i > 0 {
				prev = runes[i-1]
			}
			if i+1 < len(runes) {
				next = runes[i+1]
			}

			// Add leading space only when digit adjacent to letter/number on left
			// and left side not digit (avoid double spaces for consecutive digits).
			if prev != 0 && (prev < '0' || prev > '9') {
				if (unicode.IsLetter(prev) || unicode.IsNumber(prev)) && !unicode.IsSpace(prev) {
					b.WriteByte(' ')
				}
			}
			b.WriteString(word)

			// Between consecutive digits, force single space.
			if next >= '0' && next <= '9' {
				b.WriteByte(' ')
				continue
			}

			// Add trailing space only when digit adjacent to letter/number on right.
			if next != 0 && (unicode.IsLetter(next) || unicode.IsNumber(next)) && !unicode.IsSpace(next) {
				b.WriteByte(' ')
			}
		}

		return b.String()
	}
}

func digitWord(r rune) string {
	switch r {
	case '0':
		return "zero"
	case '1':
		return "one"
	case '2':
		return "two"
	case '3':
		return "three"
	case '4':
		return "four"
	case '5':
		return "five"
	case '6':
		return "six"
	case '7':
		return "seven"
	case '8':
		return "eight"
	case '9':
		return "nine"
	default:
		return string(r)
	}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 0x7F {
			return false
		}
	}
	return true
}

func wordEditDistance(a, b []string) int {
	// Levenshtein distance over tokens.
	// dp[j] holds distance for prefix a[:i] vs b[:j].
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	dp := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		dp[j] = j
	}

	for i := 1; i <= len(a); i++ {
		prevDiag := dp[0]
		dp[0] = i
		for j := 1; j <= len(b); j++ {
			old := dp[j]
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			ins := dp[j] + 1
			del := dp[j-1] + 1
			sub := prevDiag + cost
			dp[j] = min3(ins, del, sub)
			prevDiag = old
		}
	}
	return dp[len(b)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Evaluate is used by the CLI package. It's currently a minimal placeholder
// so the module builds; the string-level evaluation logic lives in
// EvaluateString.
func Evaluate(_ context.Context, _ string) (float64, error) {
	return 0, nil
}
