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
		{"symbol difference", 0, args{
			"héllo",
			"hello"},
		},
		{"symbol difference", 0, args{
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

func TestLoweringFunc(t *testing.T) {
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
			if got := LoweringFunc(tt.args.s); got != tt.want {
				t.Errorf("LoweringFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpaceTrimmerFunc(t *testing.T) {
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
		{"leading spaces", args{"  hello"}, "hello"},
		{"trailing spaces", args{"hello  "}, "hello"},
		{"leading and trailing spaces", args{"  hello  "}, "hello"},
		{"multiple spaces between words", args{"hello  world"}, "hello world"},
		{"leading, trailing and multiple spaces", args{"  hello  world  "}, "hello world"},
		{"multiple consecutive spaces in multiple places", args{"  hello   world  test  "}, "hello world test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpaceTrimmerFunc(tt.args.s); got != tt.want {
				t.Errorf("SpaceTrimmerFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
