package aligner

type Options struct {
	// SmallMaxLen picks Wagner-Fischer when both sequences <= this length.
	// Default 1000 words/chars per spec.
	SmallMaxLen int

	// SimilarityBandMax enables banded alignment when estimated edit distance
	// small. Set 0 to disable.
	SimilarityBandMax int
}

func (o Options) withDefaults() Options {
	if o.SmallMaxLen <= 0 {
		o.SmallMaxLen = 1000
	}
	if o.SimilarityBandMax < 0 {
		o.SimilarityBandMax = 0
	}
	if o.SimilarityBandMax == 0 {
		o.SimilarityBandMax = 128
	}
	return o
}

// Align returns edit script mapping ref -> hyp.
// Strategy:
// - small: Wagner-Fischer with traceback
// - highly similar: banded DP (Ukkonen-style) with traceback
// - large: Hirschberg (linear space) producing full script
func Align[T comparable](ref, hyp []T, opt Options) []Op[T] {
	opt = opt.withDefaults()

	m, n := len(ref), len(hyp)
	if m == 0 && n == 0 {
		return nil
	}
	if m == 0 {
		ops := make([]Op[T], 0, n)
		for i := range hyp {
			v := hyp[i]
			ops = append(ops, Op[T]{Kind: OpInsert, Hyp: &v})
		}
		return ops
	}
	if n == 0 {
		ops := make([]Op[T], 0, m)
		for i := range ref {
			v := ref[i]
			ops = append(ops, Op[T]{Kind: OpDelete, Ref: &v})
		}
		return ops
	}

	if m <= opt.SmallMaxLen && n <= opt.SmallMaxLen {
		return alignWagnerFischer(ref, hyp)
	}

	// Try banded alignment when lengths close (high similarity heuristic).
	diff := m - n
	if diff < 0 {
		diff = -diff
	}
	if diff <= opt.SimilarityBandMax {
		t := minInt(opt.SimilarityBandMax, diff+16)
		if t > 0 {
			if ops, ok := alignBanded(ref, hyp, t); ok {
				return ops
			}
		}
	}

	return alignHirschberg(ref, hyp, opt)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
