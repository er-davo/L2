package grep

import (
	"fmt"
	"regexp"
	"strings"
)

// Grep filters lines according to cfg and returns a new slice.
func Grep(lines []string, cfg Config) []string {
	var matcher func(string) bool

	if cfg.Fixed {
		pattern := cfg.Pattern
		if cfg.IgnoreCase {
			pattern = strings.ToLower(pattern)
			matcher = func(s string) bool {
				return strings.Contains(strings.ToLower(s), pattern)
			}
		} else {
			matcher = func(s string) bool {
				return strings.Contains(s, pattern)
			}
		}
	} else {
		pattern := cfg.Pattern
		if cfg.IgnoreCase {
			pattern = "(?i)" + pattern
		}
		re := regexp.MustCompile(pattern)
		matcher = func(s string) bool {
			return re.MatchString(s)
		}
	}

	matches := make(map[int]struct{})
	for i, line := range lines {
		ok := matcher(line)
		if cfg.InvertMatch {
			ok = !ok
		}
		if ok {
			matches[i] = struct{}{}
		}
	}

	if cfg.CountOnly {
		return []string{fmt.Sprintf("%d", len(matches))}
	}

	after := cfg.After
	before := cfg.Before
	if cfg.Context > 0 {
		after = cfg.Context
		before = cfg.Context
	}

	withContext := make(map[int]struct{})
	for idx := range matches {
		start := idx - before
		if start < 0 {
			start = 0
		}
		end := idx + after
		if end >= len(lines) {
			end = len(lines) - 1
		}
		for j := start; j <= end; j++ {
			withContext[j] = struct{}{}
		}
	}

	var result []string
	for i := 0; i < len(lines); i++ {
		if _, ok := withContext[i]; ok {
			if cfg.LineNum {
				result = append(result, fmt.Sprintf("%d:%s", i+1, lines[i]))
			} else {
				result = append(result, lines[i])
			}
		}
	}

	return result
}
