package textbench

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var ErrRejected = errors.New("rejected")

type PipelineFunc func(string) (string, error)

func EvaluateString(a, b string) (int, error) {
	if a == "" || b == "" {
		return 0, nil
	}

	stages := []PipelineFunc{
		// Reject inputs that need whitespace normalization (leading/trailing or
		// repeated runs of whitespace).
		WhitespaceCleaning(),
		// Unicode normalization (standardize special chars/diacritics).
		UnicodeNormalizationFunc(),
		// Reject remaining non-ASCII after normalization.
		NonASCIIFunc(),
		// Normalization stage. Currently no-op for scoring, but runs helper in pipeline.
		CaseNormalization(),
	}

	for _, stage := range stages {
		var err error
		a, err = stage(a)
		if err != nil && errors.Is(err, ErrRejected) {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
	}

	for _, stage := range stages {
		var err error
		b, err = stage(b)
		if err != nil && errors.Is(err, ErrRejected) {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func CaseNormalization() PipelineFunc {
	return func(s string) (string, error) {
		return strings.ToLower(s), nil
	}
}

func WhitespaceCleaning() PipelineFunc {
	return func(s string) (string, error) {
		// Fields splits on any whitespace and drops empty runs; joining collapses
		// repeated whitespace (including leading/trailing) into single spaces.
		normalized := strings.Join(strings.Fields(s), " ")
		if normalized != s {
			return s, ErrRejected
		}
		return s, nil
	}
}

func NonASCIIFunc() PipelineFunc {
	return func(s string) (string, error) {
		if !isASCII(s) {
			return s, ErrRejected
		}
		return s, nil
	}
}

func UnicodeNormalizationFunc() PipelineFunc {
	return func(s string) (string, error) {
		// 1) Compatibility normalize (fold some "special" forms to standard forms)
		// 2) Decompose, then drop combining marks (strip diacritics)
		n := norm.NFKC.String(s)
		out := strings.Map(func(r rune) rune {
			if unicode.Is(unicode.Mn, r) {
				return -1
			}
			return r
		}, norm.NFD.String(n))
		return out, nil
	}
}

func PunctuationRemovalFunc() PipelineFunc {
	return func(s string) (string, error) {
		out := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
				return r
			}
			return -1
		}, s)
		return out, nil
	}
}

func NumberStandardizationFunc() PipelineFunc {
	return func(s string) (string, error) {
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

		return b.String(), nil
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

// Evaluate is used by the CLI package. It's currently a minimal placeholder
// so the module builds; the string-level evaluation logic lives in
// EvaluateString.
func Evaluate(_ context.Context, _ string) (float64, error) {
	return 0, nil
}
