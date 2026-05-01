# Textbench

## Algorithm (so far)

`EvaluateString(a, b)` does:

- Normalize both strings with same pipeline:
  - collapse whitespace
  - Unicode NFKC + strip diacritics/combining marks
  - drop remaining non-ASCII
  - remove punctuation/symbols (keep letters, digits, spaces)
  - map digits `0-9` to words (`"5"` → `"five"`)
  - lowercase
  - collapse whitespace again
- Split into word tokens (`strings.Fields`)
- Compute token-level Levenshtein edit distance (insert/delete/substitute cost 1)

Score meaning: **lower better**, **0** = identical after normalization.

## Test Data

The `/data` directory contains text files used for benchmarking and testing the
text evaluation functionality.

In `data/moby-dick/`, `moby-dick-0.md` is baseline; higher suffix files contain progressively more corruption. Tests assert evaluation score is monotonic vs baseline.

The test data files are referenced in the main application (see `cmd/main.go`)
to demonstrate and benchmark the text normalization and comparison pipeline.
