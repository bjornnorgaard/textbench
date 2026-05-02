package textbench

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateString(t *testing.T) {
	r := Compare("the cat", "the bat")
	require.InEpsilon(t, 0.5, r.WER, 1e-9) // 1 sub / 2 ref words
	assert.Contains(t, r.ErrorLog.String(), "<subst>")
}

func TestCompare_WhitespaceAndPunct(t *testing.T) {
	r := Compare("Hi,\nworld!\t", "hi world")
	require.Equal(t, 0.0, r.WER)
	require.Equal(t, 0.0, r.CER)
}
