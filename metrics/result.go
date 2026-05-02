package metrics

import "github.com/bjornnorgaard/textbench/diff"

type Counts struct {
	Substitutions int // S
	Deletions     int // D
	Insertions    int // I
	ReferenceLen  int // N
}

type Result struct {
	WER      float64
	CER      float64
	Word     Counts
	Char     Counts
	ErrorLog diff.Tagged
}
