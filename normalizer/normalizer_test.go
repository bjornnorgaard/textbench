package normalizer

import "testing"

func TestNormalize(t *testing.T) {
	got := Normalize("Héllo,\nWORLD!\t")
	// NFKC + lowercase + remove punctuation + collapse whitespace.
	if got != "héllo world" {
		t.Fatalf("got %q", got)
	}
}
