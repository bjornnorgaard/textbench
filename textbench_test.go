package textbench

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluateString(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		want int
		args args
	}{
		{"empty strings", 0, args{
			"",
			""},
		},
		{"identical strings", 1, args{
			"hello",
			"hello"},
		},
		{"whitespace difference", 0, args{
			"hello ",
			"hello"},
		},
		{"non-empty strings", 1, args{
			"hello",
			"world"},
		},
		{"casing only difference", 1, args{
			"Hello",
			"hello"},
		},
		{"pluralization only difference", 1, args{
			"world",
			"worlds"},
		},
		{"symbol difference", 1, args{
			"héllo",
			"hello"},
		},
		{"symbol difference", 1, args{
			"héllo",
			"hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateString(tt.args.a, tt.args.b)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
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
			s, err := CaseNormalization()(tt.args.s)
			require.NoError(t, err)
			require.Equal(t, tt.want, s)
		})
	}
}

func TestWhitespaceCleaning(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"empty string", args{""}, nil},
		{"no spaces", args{"hello"}, nil},
		{"single space between words", args{"hello world"}, nil},
		{"leading spaces", args{"  hello"}, ErrRejected},
		{"trailing spaces", args{"hello  "}, ErrRejected},
		{"leading and trailing spaces", args{"  hello  "}, ErrRejected},
		{"multiple spaces between words", args{"hello  world"}, ErrRejected},
		{"leading, trailing and multiple spaces", args{"  hello  world  "}, ErrRejected},
		{"multiple consecutive spaces in multiple places", args{"  hello   world  test  "}, ErrRejected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := WhitespaceCleaning()(tt.args.s)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.args.s, s)
		})
	}
}

func TestNonASCIIFunc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{"empty string", args{""}, "", nil},
		{"ascii", args{"hello world"}, "hello world", nil},
		{"accented", args{"héllo"}, "héllo", ErrRejected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NonASCIIFunc()(tt.args.s)
			require.ErrorIs(t, err, tt.wantErr)
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
			s, err := UnicodeNormalizationFunc()(tt.args.s)
			require.NoError(t, err)
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
			s, err := PunctuationRemovalFunc()(tt.args.s)
			require.NoError(t, err)
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
			s, err := NumberStandardizationFunc()(tt.args.s)
			require.NoError(t, err)
			require.Equal(t, tt.want, s)
		})
	}
}
