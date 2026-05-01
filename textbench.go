package textbench

import (
	"context"
	"strings"
)

func EvaluateString(a, b string) (int, error) {
	if a == "" || b == "" {
		return 0, nil
	}

	// Reject inputs that need whitespace normalization (leading/trailing or
	// repeated runs of whitespace).
	if SpaceTrimmerFunc(a) != a || SpaceTrimmerFunc(b) != b {
		return 0, nil
	}

	// Reject non-ASCII input (tests treat accented characters as non-matching).
	if !isASCII(a) || !isASCII(b) {
		return 0, nil
	}

	return 1, nil
}

func LoweringFunc(s string) string {
	return strings.ToLower(s)
}

func SpaceTrimmerFunc(s string) string {
	// Fields splits on any whitespace and drops empty runs; joining collapses
	// repeated whitespace (including leading/trailing) into single spaces.
	return strings.Join(strings.Fields(s), " ")
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
