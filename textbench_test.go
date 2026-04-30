package textbench

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_readChunk(t *testing.T) {
	expected := "Ishmael"
	path := filepath.Join(".", "data", "moby-dick-0.txt")
	f, err := os.Open(path)
	require.NoError(t, err)
	reader := bufio.NewReader(f)

	read, bytes, err := readChunk(reader, 8, len(expected))
	require.NoError(t, err)
	assert.Equal(t, len(expected), read)
	assert.Equal(t, expected, string(bytes))
}
