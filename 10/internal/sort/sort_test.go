package sort_test

import (
	"testing"

	"mysort/internal/config"
	"mysort/internal/sort"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		cfg     config.Config
		want    []string
		wantErr bool
	}{
		{
			name:  "simple alphabetical",
			lines: []string{"banana", "apple", "cherry"},
			cfg:   config.Config{},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "numeric sort",
			lines: []string{"10", "2", "1"},
			cfg:   config.Config{Numeric: true},
			want:  []string{"1", "2", "10"},
		},
		{
			name:  "reverse sort",
			lines: []string{"a", "c", "b"},
			cfg:   config.Config{Reverse: true},
			want:  []string{"c", "b", "a"},
		},
		{
			name:  "unique values",
			lines: []string{"apple", "apple", "banana"},
			cfg:   config.Config{Unique: true},
			want:  []string{"apple", "banana"},
		},
		{
			name:  "month sort",
			lines: []string{"Mar", "Jan", "Feb"},
			cfg:   config.Config{Month: true},
			want:  []string{"Jan", "Feb", "Mar"},
		},
		{
			name:    "check sorted - not sorted",
			lines:   []string{"a", "c", "b"},
			cfg:     config.Config{CheckSorted: true},
			wantErr: true,
		},
		{
			name:  "sort by column",
			lines: []string{"3\tb", "1\ta", "2\tc"},
			cfg:   config.Config{Column: 2, Delimiter: "\t"},
			want:  []string{"1\ta", "3\tb", "2\tc"},
		},
		{
			name:  "ignore blanks",
			lines: []string{"  apple", "banana  ", " cherry "},
			cfg:   config.Config{IgnoreBlanks: true},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "human readable numbers",
			lines: []string{"10K", "2M", "500"},
			cfg:   config.Config{Numeric: true, Human: true},
			want:  []string{"500", "10K", "2M"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sort.Sort(tt.lines, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
