package textbench

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateString(t *testing.T) {
	var (
		correct = "Captain Ahab stood on the deck and watched the endless sea as the crew prepared the ropes and sails."
		minor   = "Captain Ahab stood on the desk and watched the endless see as the crew prepared the rope and sails."
		medium  = "Captain Ahab stood on deck and watched endless sea as crew prepared ropes and sails."
		bad     = "Captn Ahhab stood deck watched endles sea; crew prepare rope sail."
		worst   = "Yesterday we drove to town and forgot the groceries while the radio played loudly."
	)

	s0, err := EvaluateString(correct, correct)
	require.NoError(t, err)
	require.Equal(t, 0, s0)

	s1, err := EvaluateString(correct, minor)
	require.NoError(t, err)
	s2, err := EvaluateString(correct, medium)
	require.NoError(t, err)
	s3, err := EvaluateString(correct, bad)
	require.NoError(t, err)
	s4, err := EvaluateString(correct, worst)
	require.NoError(t, err)

	// Expect strict improvement ordering in broad strokes:
	assert.True(t, s0 < s1)
	assert.True(t, s1 < s2)
	assert.True(t, s2 < s3)
	assert.True(t, s3 < s4)
}

func TestEvaluateString_DataFilesOrdering(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("data", "moby-dick", "moby-dick-*.md"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	re := regexp.MustCompile(`moby-dick-(\d+)\.md$`)
	type item struct {
		path  string
		level int
		text  string
	}

	var items []item
	for _, p := range files {
		m := re.FindStringSubmatch(p)
		require.NotNil(t, m, "unexpected file name: %s", p)
		level := 0
		for i := 0; i < len(m[1]); i++ {
			level = level*10 + int(m[1][i]-'0')
		}

		b, readErr := os.ReadFile(p)
		require.NoError(t, readErr)
		items = append(items, item{path: p, level: level, text: string(b)})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].level < items[j].level })
	require.GreaterOrEqual(t, len(items), 2)

	correct := items[0].text
	same, err := EvaluateString(correct, correct)
	require.NoError(t, err)
	require.Equal(t, 0, same)

	prev := -1
	for _, it := range items {
		score, evalErr := EvaluateString(correct, it.text)
		require.NoError(t, evalErr, "file: %s", it.path)
		if prev >= 0 {
			require.GreaterOrEqual(t, score, prev, "non-monotonic at %s", it.path)
		}
		prev = score
	}

	firstScore, err := EvaluateString(correct, items[0].text)
	require.NoError(t, err)
	lastScore, err := EvaluateString(correct, items[len(items)-1].text)
	require.NoError(t, err)
	require.GreaterOrEqual(t, lastScore, firstScore)
}

func TestCaseNormalization(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty string", args{""}, ""},
		{"already lowercase", args{"hello"}, "hello"},
		{"uppercase string", args{"HELLO"}, "hello"},
		{"mixed case string", args{"HeLLo WoRLd"}, "hello world"},
		{"with numbers and special chars", args{"Hello123!@#"}, "hello123!@#"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CaseNormalization()(tt.args.s)
			require.Equal(t, tt.want, s)
		})
	}
}

func TestWhitespaceCleaning(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty string", args{""}, ""},
		{"no spaces", args{"hello"}, "hello"},
		{"single space between words", args{"hello world"}, "hello world"},
		{"leading spaces", args{"  hello"}, "hello"},
		{"trailing spaces", args{"hello  "}, "hello"},
		{"leading and trailing spaces", args{"  hello  "}, "hello"},
		{"multiple spaces between words", args{"hello  world"}, "hello world"},
		{"leading, trailing and multiple spaces", args{"  hello  world  "}, "hello world"},
		{"multiple consecutive spaces in multiple places", args{"  hello   world  test  "}, "hello world test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := WhitespaceCleaning()(tt.args.s)
			require.Equal(t, tt.want, s)
		})
	}
}

func TestNonASCIIFunc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty string", args{""}, ""},
		{"ascii", args{"hello world"}, "hello world"},
		{"accented stripped", args{"héllo"}, "hllo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NonASCIIFunc()(tt.args.s)
			require.Equal(t, tt.want, s)
		})
	}
}

func TestUnicodeNormalizationFunc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty string", args{""}, ""},
		{"ascii passthrough", args{"hello world"}, "hello world"},
		{"strip diacritics", args{"héllo"}, "hello"},
		{"strip combining marks", args{"e\u0301"}, "e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UnicodeNormalizationFunc()(tt.args.s)
			require.Equal(t, tt.want, s)
		})
	}
}

func TestPunctuationRemovalFunc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty string", args{""}, ""},
		{"no punctuation", args{"hello world 123"}, "hello world 123"},
		{"strip punctuation", args{"hello, world!"}, "hello world"},
		{"strip symbols", args{"abc-123 + xyz"}, "abc123  xyz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PunctuationRemovalFunc()(tt.args.s)
			require.Equal(t, tt.want, s)
		})
	}
}

func TestNumberStandardizationFunc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty string", args{""}, ""},
		{"single digit", args{"5"}, "five"},
		{"digits in word boundaries", args{"x5y"}, "x five y"},
		{"multiple digits", args{"123"}, "one two three"},
		{"keep non-digits", args{"v2.0"}, "v two.zero"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NumberStandardizationFunc()(tt.args.s)
			require.Equal(t, tt.want, s)
		})
	}
}
