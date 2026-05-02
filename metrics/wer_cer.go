package metrics

import (
	"strings"

	"github.com/bjornnorgaard/textbench/aligner"
	"github.com/bjornnorgaard/textbench/diff"
)

type Options struct {
	Aligner aligner.Options
	Diff    diff.Options
}

func CompareWords(refTokens, hypTokens []string, opt Options) (Counts, float64, diff.Tagged) {
	ops := aligner.Align(refTokens, hypTokens, opt.Aligner)

	var c Counts
	refParts := make([]string, 0, len(refTokens))
	hypParts := make([]string, 0, len(hypTokens))

	for _, op := range ops {
		switch op.Kind {
		case aligner.OpEqual:
			c.ReferenceLen++
			refParts = append(refParts, *op.Ref)
			hypParts = append(hypParts, *op.Hyp)
		case aligner.OpSubstitute:
			c.Substitutions++
			c.ReferenceLen++
			refParts = append(refParts, "<subst>"+*op.Ref+"</subst>")
			hypParts = append(hypParts, "<subst>"+*op.Hyp+"</subst>")
		case aligner.OpDelete:
			c.Deletions++
			c.ReferenceLen++
			refParts = append(refParts, "<del>"+*op.Ref+"</del>")
		case aligner.OpInsert:
			c.Insertions++
			hypParts = append(hypParts, "<ins>"+*op.Hyp+"</ins>")
		}
	}

	wer := 0.0
	if c.ReferenceLen == 0 {
		if c.Insertions == 0 {
			wer = 0
		} else {
			wer = 1
		}
	} else {
		wer = float64(c.Substitutions+c.Deletions+c.Insertions) / float64(c.ReferenceLen)
	}

	return c, wer, diff.JoinTokens(refParts, hypParts, opt.Diff)
}

func CompareRunes(refRunes, hypRunes []rune, opt Options) (Counts, float64) {
	ops := aligner.Align(refRunes, hypRunes, opt.Aligner)
	var c Counts
	for _, op := range ops {
		switch op.Kind {
		case aligner.OpEqual:
			c.ReferenceLen++
		case aligner.OpSubstitute:
			c.Substitutions++
			c.ReferenceLen++
		case aligner.OpDelete:
			c.Deletions++
			c.ReferenceLen++
		case aligner.OpInsert:
			c.Insertions++
		}
	}

	cer := 0.0
	if c.ReferenceLen == 0 {
		if c.Insertions == 0 {
			cer = 0
		} else {
			cer = 1
		}
	} else {
		cer = float64(c.Substitutions+c.Deletions+c.Insertions) / float64(c.ReferenceLen)
	}
	return c, cer
}

func TokenizeWords(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Fields(s)
}
