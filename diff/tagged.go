package diff

import (
	"strings"
)

type Tagged struct {
	Ref string
	Hyp string
}

func (t Tagged) String() string {
	if t.Ref == "" && t.Hyp == "" {
		return ""
	}
	return "REF: " + t.Ref + "\nHYP: " + t.Hyp
}

type Options struct {
	// MaxChars limits output size; 0 means default 4096.
	MaxChars int
}

func (o Options) withDefaults() Options {
	if o.MaxChars <= 0 {
		o.MaxChars = 4096
	}
	return o
}

// JoinTokens builds simple tagged diff for token sequences.
// Input ops kinds: "equal", "sub", "ins", "del".
func JoinTokens(refParts, hypParts []string, opt Options) Tagged {
	opt = opt.withDefaults()
	ref := strings.Join(refParts, " ")
	hyp := strings.Join(hypParts, " ")
	out := Tagged{Ref: ref, Hyp: hyp}
	return truncate(out, opt.MaxChars)
}

func truncate(t Tagged, max int) Tagged {
	s := t.String()
	if len(s) <= max {
		return t
	}
	keep := max - len("\n... truncated ...")
	if keep < 0 {
		keep = 0
	}
	s = s[:keep] + "\n... truncated ..."

	// Re-split back into Ref/Hyp best-effort.
	lines := strings.SplitN(s, "\n", 3)
	var ref, hyp string
	if len(lines) > 0 {
		ref = strings.TrimPrefix(lines[0], "REF: ")
	}
	if len(lines) > 1 {
		hyp = strings.TrimPrefix(lines[1], "HYP: ")
	}
	return Tagged{Ref: ref, Hyp: hyp}
}
