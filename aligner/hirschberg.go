package aligner

// Hirschberg algorithm: linear-space edit script for unit costs.

func alignHirschberg[T comparable](ref, hyp []T, opt Options) []Op[T] {
	m, n := len(ref), len(hyp)
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
	if m == 1 || n == 1 {
		// Small base-case; full DP cheap and yields good script.
		return alignWagnerFischer(ref, hyp)
	}

	iMid := m / 2

	leftCost := costRow(ref[:iMid], hyp)
	rightCost := costRowRev(ref[iMid:], hyp)

	// Choose split j minimizing leftCost[j] + rightCost[n-j]
	bestJ := 0
	best := leftCost[0] + rightCost[n]
	for j := 1; j <= n; j++ {
		v := leftCost[j] + rightCost[n-j]
		if v < best {
			best = v
			bestJ = j
		}
	}

	opsL := alignHirschberg(ref[:iMid], hyp[:bestJ], opt)
	opsR := alignHirschberg(ref[iMid:], hyp[bestJ:], opt)
	return append(opsL, opsR...)
}

func costRow[T comparable](ref, hyp []T) []int {
	// Levenshtein distance last row, linear space.
	n := len(hyp)
	prev := make([]int, n+1)
	cur := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prev[j] = j
	}
	for i := 1; i <= len(ref); i++ {
		cur[0] = i
		for j := 1; j <= n; j++ {
			cost := 0
			if ref[i-1] != hyp[j-1] {
				cost = 1
			}
			diag := prev[j-1] + cost
			up := prev[j] + 1
			left := cur[j-1] + 1
			cur[j] = min3(diag, up, left)
		}
		prev, cur = cur, prev
	}
	return prev
}

func costRowRev[T comparable](ref, hyp []T) []int {
	// distance between reverse(ref) and reverse(hyp)
	rr := make([]T, len(ref))
	for i := range ref {
		rr[len(ref)-1-i] = ref[i]
	}
	rh := make([]T, len(hyp))
	for i := range hyp {
		rh[len(hyp)-1-i] = hyp[i]
	}
	return costRow(rr, rh)
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
