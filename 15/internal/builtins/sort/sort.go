package sort

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// ErrNotSorted is returned when input is not sorted.
var ErrNotSorted = errors.New("input is not sorted")

// Sort sorts lines according to cfg and returns a new slice.
func Sort(lines []string, cfg Config) ([]string, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	if cfg.IgnoreBlanks {
		for i := 0; i < len(lines); i++ {
			lines[i] = strings.TrimSpace(lines[i])
		}
	}

	// copy lines to not modify original
	out := make([]string, len(lines))
	copy(out, lines)

	less := func(a, b string) bool {
		aKey, bKey := strings.Clone(a), strings.Clone(b)

		if cfg.Column > 0 {
			aCols := strings.Split(a, cfg.Delimiter)
			bCols := strings.Split(b, cfg.Delimiter)
			if cfg.Column <= len(aCols) {
				aKey = aCols[cfg.Column-1]
			}
			if cfg.Column <= len(bCols) {
				bKey = bCols[cfg.Column-1]
			}
		}

		if cfg.IgnoreBlanks {
			aKey = strings.TrimSpace(aKey)
			bKey = strings.TrimSpace(bKey)
		}

		if cfg.Month {
			ma, okA := parseMonth(aKey)
			mb, okB := parseMonth(bKey)
			if okA && okB {
				if cfg.Reverse {
					return ma > mb
				}
				return ma < mb
			}
		}

		if cfg.Numeric {
			na, errA := parseNumber(aKey, cfg.Human)
			nb, errB := parseNumber(bKey, cfg.Human)
			if errA == nil && errB == nil {
				if cfg.Reverse {
					return na > nb
				}
				return na < nb
			}
		}

		if cfg.Reverse {
			return bKey < aKey
		}
		return aKey < bKey
	}

	if cfg.CheckSorted {
		for i := 0; i < len(out)-1; i++ {
			if less(out[i+1], out[i]) {
				return nil, ErrNotSorted
			}
		}
		return out, nil
	}

	sort.Slice(out, func(i, j int) bool {
		return less(out[i], out[j])
	})

	if cfg.Unique {
		uniq := out[:0]
		var prev string
		for i, line := range out {
			if i == 0 || line != prev {
				uniq = append(uniq, line)
			}
			prev = line
		}
		out = uniq
	}

	return out, nil
}

func parseMonth(s string) (int, bool) {
	months := map[string]int{
		"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
		"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
		"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
	}
	if val, ok := months[s]; ok {
		return val, true
	}
	return 0, false
}

func parseNumber(s string, human bool) (float64, error) {
	if !human {
		return strconv.ParseFloat(s, 64)
	}

	mult := 1.0
	last := len(s) - 1
	if last >= 0 {
		switch s[last] {
		case 'K', 'k':
			mult = 1 << 10
			s = s[:last]
		case 'M', 'm':
			mult = 1 << 20
			s = s[:last]
		case 'G', 'g':
			mult = 1 << 30
			s = s[:last]
		}
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return val * mult, nil
}
