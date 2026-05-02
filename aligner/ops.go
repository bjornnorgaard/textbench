package aligner

type OpKind uint8

const (
	OpEqual OpKind = iota
	OpSubstitute
	OpInsert
	OpDelete
)

type Op[T comparable] struct {
	Kind OpKind
	Ref  *T
	Hyp  *T
}

func equalOp[T comparable](v T) Op[T] { return Op[T]{Kind: OpEqual, Ref: &v, Hyp: &v} }
