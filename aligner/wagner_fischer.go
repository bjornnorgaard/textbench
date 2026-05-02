package aligner

// alignWagnerFischer computes full DP + traceback.
// Cost: ins=del=sub=1, eq=0.
func alignWagnerFischer[T comparable](ref, hyp []T) []Op[T] {
	m, n := len(ref), len(hyp)

	// back[i*(n+1)+j] gives move for cell (i,j):
	// 0 diag, 1 up (delete), 2 left (insert)
	back := make([]uint8, (m+1)*(n+1))

	prev := make([]int, n+1)
	cur := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prev[j] = j
		if j > 0 {
			back[j] = 2
		}
	}
	for i := 1; i <= m; i++ {
		cur[0] = i
		back[i*(n+1)] = 1
		for j := 1; j <= n; j++ {
			cost := 0
			if ref[i-1] != hyp[j-1] {
				cost = 1
			}
			diag := prev[j-1] + cost
			up := prev[j] + 1
			left := cur[j-1] + 1

			best := diag
			move := uint8(0)
			if up < best {
				best = up
				move = 1
			}
			if left < best {
				best = left
				move = 2
			}
			cur[j] = best
			back[i*(n+1)+j] = move
		}
		prev, cur = cur, prev
	}

	// Traceback.
	i, j := m, n
	ops := make([]Op[T], 0, m+n)
	for i > 0 || j > 0 {
		move := back[i*(n+1)+j]
		switch move {
		case 0: // diag
			r := ref[i-1]
			h := hyp[j-1]
			if r == h {
				ops = append(ops, Op[T]{Kind: OpEqual, Ref: &r, Hyp: &h})
			} else {
				ops = append(ops, Op[T]{Kind: OpSubstitute, Ref: &r, Hyp: &h})
			}
			i--
			j--
		case 1: // up
			r := ref[i-1]
			ops = append(ops, Op[T]{Kind: OpDelete, Ref: &r})
			i--
		default: // left
			h := hyp[j-1]
			ops = append(ops, Op[T]{Kind: OpInsert, Hyp: &h})
			j--
		}
	}

	// reverse in-place
	for l, r := 0, len(ops)-1; l < r; l, r = l+1, r-1 {
		ops[l], ops[r] = ops[r], ops[l]
	}
	return ops
}
