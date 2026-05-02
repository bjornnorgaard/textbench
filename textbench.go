package textbench

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/bjornnorgaard/textbench/metrics"
	"github.com/bjornnorgaard/textbench/normalizer"
)

type Options struct {
	Metrics metrics.Options
}

func Compare(ref, hyp string) metrics.Result {
	nref := normalizer.Normalize(ref)
	nhyp := normalizer.Normalize(hyp)

	refTokens := metrics.TokenizeWords(nref)
	hypTokens := metrics.TokenizeWords(nhyp)

	wordCounts, wer, tagged := metrics.CompareWords(refTokens, hypTokens, metrics.Options{})
	charCounts, cer := metrics.CompareRunes([]rune(nref), []rune(nhyp), metrics.Options{})

	return metrics.Result{
		WER:      wer,
		CER:      cer,
		Word:     wordCounts,
		Char:     charCounts,
		ErrorLog: tagged,
	}
}

func CompareWithOptions(ref, hyp string, opt Options) metrics.Result {
	nref := normalizer.Normalize(ref)
	nhyp := normalizer.Normalize(hyp)

	refTokens := metrics.TokenizeWords(nref)
	hypTokens := metrics.TokenizeWords(nhyp)

	wordCounts, wer, tagged := metrics.CompareWords(refTokens, hypTokens, opt.Metrics)
	charCounts, cer := metrics.CompareRunes([]rune(nref), []rune(nhyp), opt.Metrics)

	return metrics.Result{
		WER:      wer,
		CER:      cer,
		Word:     wordCounts,
		Char:     charCounts,
		ErrorLog: tagged,
	}
}

func CompareBatch(refs, hyps []string) ([]metrics.Result, error) {
	if len(refs) != len(hyps) {
		return nil, errors.New("refs/hyps length mismatch")
	}
	if len(refs) == 0 {
		return nil, nil
	}

	out := make([]metrics.Result, len(refs))

	// Simple worker pool. Avoid goroutine explosion on huge batches.
	workers := max(1, min(len(refs), 2*runtime.GOMAXPROCS(0)))

	var wg sync.WaitGroup
	jobs := make(chan int)
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := range jobs {
				out[i] = Compare(refs[i], hyps[i])
			}
		}()
	}

	for i := range refs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return out, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// EvaluateString kept for compatibility.
// Returns token edit distance (S+D+I) after spec normalizer.
func EvaluateString(a, b string) (int, error) {
	r := Compare(a, b)
	return r.Word.Substitutions + r.Word.Deletions + r.Word.Insertions, nil
}

// Evaluate is used by the CLI package. It's currently a minimal placeholder
// so the module builds; the string-level evaluation logic lives in
// EvaluateString.
func Evaluate(_ context.Context, _ string) (float64, error) {
	return 0, nil
}
