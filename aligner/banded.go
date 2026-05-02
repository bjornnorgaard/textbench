package aligner

// Banded (Ukkonen-style) DP with traceback within band t.
// Returns ok=false when optimal path exits band.

const big = int(^uint(0) >> 1) // max int

func alignBanded[T comparable](ref, hyp []T, t int) ([]Op[T], bool) {
	m, n := len(ref), len(hyp)
	if t <= 0 {
		return nil, false
	}
	// If band too narrow to even reach (m,n).
	if abs(m-n) > t {
		return nil, false
	}

	// For each i, j ranges [lo(i), hi(i)] inclusive.
	lo := func(i int) int {
		if i-t > 0 {
			return i - t
		}
		return 0
	}
	hi := func(i int) int {
		if i+t < n {
			return i + t
		}
		return n
	}

	// Store backpointers for each cell in band as uint8:
	// 0 diag, 1 up (delete), 2 left (insert)
	type rowMeta struct {
		lo int
		hi int
	}
	meta := make([]rowMeta, m+1)
	back := make([][]uint8, m+1)

	prev := make([]int, 0)
	prevLo, prevHi := 0, 0

	for i := 0; i <= m; i++ {
		curLo, curHi := lo(i), hi(i)
		meta[i] = rowMeta{lo: curLo, hi: curHi}
		width := curHi - curLo + 1
		cur := make([]int, width)
		back[i] = make([]uint8, width)

		for k := 0; k < width; k++ {
			cur[k] = big
		}

		if i == 0 {
			// dp[0][j] = j within band
			for j := curLo; j <= curHi; j++ {
				cur[j-curLo] = j
				if j > 0 {
					back[i][j-curLo] = 2
				}
			}
			prev, prevLo, prevHi = cur, curLo, curHi
			continue
		}

		for j := curLo; j <= curHi; j++ {
			best := big
			move := uint8(0)

			// diag: (i-1,j-1)
			if j-1 >= prevLo && j-1 <= prevHi && j-1 >= 0 {
				cost := 0
				if ref[i-1] != hyp[j-1] {
					cost = 1
				}
				v := prev[(j-1)-prevLo] + cost
				best = v
				move = 0
			}

			// up: (i-1,j) delete
			if j >= prevLo && j <= prevHi {
				v := prev[j-prevLo] + 1
				if v < best {
					best = v
					move = 1
				}
			}

			// left: (i,j-1) insert
			if j-1 >= curLo && j-1 <= curHi && j-1 >= 0 {
				v := cur[(j-1)-curLo] + 1
				if v < best {
					best = v
					move = 2
				}
			}

			cur[j-curLo] = best
			back[i][j-curLo] = move
		}

		prev, prevLo, prevHi = cur, curLo, curHi
	}

	// If (m,n) outside band, fail.
	last := meta[m]
	if n < last.lo || n > last.hi {
		return nil, false
	}

	// Traceback from (m,n)
	i, j := m, n
	ops := make([]Op[T], 0, m+n)
	for i > 0 || j > 0 {
		rm := meta[i]
		if j < rm.lo || j > rm.hi {
			return nil, false
		}
		move := back[i][j-rm.lo]
		switch move {
		case 0:
			if i == 0 || j == 0 {
				return nil, false
			}
			r := ref[i-1]
			h := hyp[j-1]
			if r == h {
				ops = append(ops, Op[T]{Kind: OpEqual, Ref: &r, Hyp: &h})
			} else {
				ops = append(ops, Op[T]{Kind: OpSubstitute, Ref: &r, Hyp: &h})
			}
			i--
			j--
		case 1:
			if i == 0 {
				return nil, false
			}
			r := ref[i-1]
			ops = append(ops, Op[T]{Kind: OpDelete, Ref: &r})
			i--
		default:
			if j == 0 {
				return nil, false
			}
			h := hyp[j-1]
			ops = append(ops, Op[T]{Kind: OpInsert, Hyp: &h})
			j--
		}
	}

	for l, r := 0, len(ops)-1; l < r; l, r = l+1, r-1 {
		ops[l], ops[r] = ops[r], ops[l]
	}
	return ops, true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
