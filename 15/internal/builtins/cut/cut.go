package cut

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func parseFields(fields string) ([]int, error) {
	set := map[int]struct{}{}

	parts := strings.Split(fields, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			bounds := strings.Split(p, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start > end {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			for i := start; i <= end; i++ {
				set[i] = struct{}{}
			}
		} else {
			val, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid field: %s", p)
			}
			set[val] = struct{}{}
		}
	}

	result := make([]int, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Ints(result)
	return result, nil
}

// Cut filters lines according to config.
func Cut(lines []string, cfg Config) ([]string, error) {
	delimiter := cfg.Delimiter
	if delimiter == "" {
		delimiter = "\t"
	}

	idxs, err := parseFields(cfg.Fields)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, line := range lines {
		if cfg.Separated && !strings.Contains(line, delimiter) {
			continue
		}

		cols := strings.Split(line, delimiter)
		var selected []string
		for _, idx := range idxs {
			if idx-1 < len(cols) {
				selected = append(selected, cols[idx-1])
			}
		}
		if len(selected) > 0 {
			result = append(result, strings.Join(selected, delimiter))
		}
	}

	return result, nil
}
