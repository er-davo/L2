package unpack_test

import (
	"testing"

	unpack "string-unpack"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
			err:      nil,
		},
		{
			name:     "simple string with digits",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			err:      nil,
		},
		{
			name:     "simple string without digits",
			input:    "abcd",
			expected: "abcd",
			err:      nil,
		},
		{
			name:     "string with escaped backslash",
			input:    "a\\4bc2d\\5e",
			expected: "a4bccd5e",
			err:      nil,
		},
		{
			name:     "string with only digits",
			input:    "45",
			expected: "",
			err:      unpack.ErrInvalidString,
		},
		{
			name:     "digit with escaped backslash",
			input:    "\\45",
			expected: "44444",
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := unpack.Unpack(tc.input)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
