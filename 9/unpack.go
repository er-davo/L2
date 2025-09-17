package unpack

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString is returned when the input string is invalid.
var ErrInvalidString = fmt.Errorf("invalid string")

// Unpack returns unpacked string or error if s is invalid
func Unpack(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	if unicode.IsDigit(rune(s[0])) {
		return "", fmt.Errorf("%w: digit at the start", ErrInvalidString)
	}

	builder := strings.Builder{}

	runed := []rune(s)

	for i := 0; i < len(runed); i++ {
		if runed[i] == '\\' {
			i++
			if i < len(runed) {
				builder.WriteRune(runed[i])
			}
			continue
		}
		if unicode.IsDigit(runed[i]) {
			// i is guaranteed to be >= 1 cause of the check at the start of func
			symbolToUnpack := runed[i-1]
			stringNum := []rune{}
			for i < len(runed) && unicode.IsDigit(runed[i]) {
				stringNum = append(stringNum, runed[i])
				i++
			}
			i-- // i is incremented at the end of loop

			num, _ := strconv.Atoi(string(stringNum))

			// num-1 because unpacked symbol is written before
			builder.WriteString(strings.Repeat(string(symbolToUnpack), num-1))
		} else {
			builder.WriteRune(runed[i])
		}
	}

	return builder.String(), nil
}
