package textbench

import (
	"strings"
	"testing"
)

func BenchmarkCompare_Small(b *testing.B) {
	ref := "the quick brown fox jumps over the lazy dog"
	hyp := "the quick brown fax jumps over lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Compare(ref, hyp)
	}
}

func BenchmarkCompare_1kWords(b *testing.B) {
	refTok := strings.Repeat("word ", 1000)
	hypTok := strings.Repeat("word ", 999) + "ward "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Compare(refTok, hypTok)
	}
}

func BenchmarkCompare_100kChars(b *testing.B) {
	ref := strings.Repeat("hello world ", 9000) // ~108k chars
	hyp := strings.Repeat("hello world ", 8999) + "hello wurld "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Compare(ref, hyp)
	}
}
