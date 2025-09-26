package grep_test

import (
	"minishell/internal/builtins/grep"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	lines := []string{
		"Hello World",
		"foo bar",
		"HELLO again",
		"something else",
		"foo HELLO",
	}

	tests := []struct {
		name string
		cfg  grep.Config
		want []string
	}{
		{
			name: "simple match",
			cfg:  grep.Config{Pattern: "Hello"},
			want: []string{"Hello World"},
		},
		{
			name: "ignore case",
			cfg:  grep.Config{Pattern: "hello", IgnoreCase: true},
			want: []string{"Hello World", "HELLO again", "foo HELLO"},
		},
		{
			name: "invert match",
			cfg:  grep.Config{Pattern: "foo", InvertMatch: true},
			want: []string{"Hello World", "HELLO again", "something else"},
		},
		{
			name: "fixed string match",
			cfg:  grep.Config{Pattern: "foo", Fixed: true},
			want: []string{"foo bar", "foo HELLO"},
		},
		{
			name: "count only",
			cfg:  grep.Config{Pattern: "HELLO", IgnoreCase: true, CountOnly: true},
			want: []string{"3"},
		},
		{
			name: "line numbers",
			cfg:  grep.Config{Pattern: "foo", Fixed: true, LineNum: true},
			want: []string{"2:foo bar", "5:foo HELLO"},
		},
		{
			name: "after context",
			cfg:  grep.Config{Pattern: "foo bar", After: 1},
			want: []string{"foo bar", "HELLO again"},
		},
		{
			name: "before context",
			cfg:  grep.Config{Pattern: "foo HELLO", Before: 1},
			want: []string{"something else", "foo HELLO"},
		},
		{
			name: "context both sides (C)",
			cfg:  grep.Config{Pattern: "HELLO again", Context: 1},
			want: []string{"foo bar", "HELLO again", "something else"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := grep.Grep(lines, tt.cfg)
			assert.Equal(t, tt.want, got)
		})
	}
}
