package cut_test

import (
	"mycut/internal/config"
	"mycut/internal/cut"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		cfg     config.Config
		want    []string
		wantErr bool
	}{
		{
			name:  "single column",
			lines: []string{"a\tb\tc", "1\t2\t3"},
			cfg:   config.Config{Fields: "1", Delimiter: "\t"},
			want:  []string{"a", "1"},
		},
		{
			name:  "multiple columns",
			lines: []string{"a:b:c:d", "1:2:3:4"},
			cfg:   config.Config{Fields: "1,3-4", Delimiter: ":"},
			want:  []string{"a:c:d", "1:3:4"},
		},
		{
			name:  "range columns",
			lines: []string{"x y z w", "1 2 3 4"},
			cfg:   config.Config{Fields: "2-3", Delimiter: " "},
			want:  []string{"y z", "2 3"},
		},
		{
			name:  "ignore lines without delimiter -s",
			lines: []string{"a:b", "no_delim"},
			cfg:   config.Config{Fields: "1", Delimiter: ":", Separated: true},
			want:  []string{"a"},
		},
		{
			name:  "fields out of range",
			lines: []string{"a:b:c"},
			cfg:   config.Config{Fields: "2,5", Delimiter: ":"},
			want:  []string{"b"},
		},
		{
			name:    "invalid field",
			lines:   []string{"a:b:c"},
			cfg:     config.Config{Fields: "x", Delimiter: ":"},
			wantErr: true,
		},
		{
			name:    "invalid range",
			lines:   []string{"a:b:c"},
			cfg:     config.Config{Fields: "3-1", Delimiter: ":"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cut.Cut(tt.lines, tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
