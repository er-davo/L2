package grep_test

import (
	"testing"

	"mygrep/internal/config"
	"mygrep/internal/grep"

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
		cfg  config.Config
		want []string
	}{
		{
			name: "simple match",
			cfg:  config.Config{Pattern: "Hello"},
			want: []string{"Hello World"},
		},
		{
			name: "ignore case",
			cfg:  config.Config{Pattern: "hello", IgnoreCase: true},
			want: []string{"Hello World", "HELLO again", "foo HELLO"},
		},
		{
			name: "invert match",
			cfg:  config.Config{Pattern: "foo", InvertMatch: true},
			want: []string{"Hello World", "HELLO again", "something else"},
		},
		{
			name: "fixed string match",
			cfg:  config.Config{Pattern: "foo", Fixed: true},
			want: []string{"foo bar", "foo HELLO"},
		},
		{
			name: "count only",
			cfg:  config.Config{Pattern: "HELLO", IgnoreCase: true, CountOnly: true},
			want: []string{"3"},
		},
		{
			name: "line numbers",
			cfg:  config.Config{Pattern: "foo", Fixed: true, LineNum: true},
			want: []string{"2:foo bar", "5:foo HELLO"},
		},
		{
			name: "after context",
			cfg:  config.Config{Pattern: "foo bar", After: 1},
			want: []string{"foo bar", "HELLO again"},
		},
		{
			name: "before context",
			cfg:  config.Config{Pattern: "foo HELLO", Before: 1},
			want: []string{"something else", "foo HELLO"},
		},
		{
			name: "context both sides (C)",
			cfg:  config.Config{Pattern: "HELLO again", Context: 1},
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
