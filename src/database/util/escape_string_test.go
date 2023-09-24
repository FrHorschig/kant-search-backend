package util

import (
	"fmt"
	"testing"
)

func TestEscapeSpecialChars(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    `hello`,
			expected: `hello`,
		},
		{
			input:    `hello\world`,
			expected: `hello\\world`,
		},
		{
			input:    `hello&world`,
			expected: `hello\&world`,
		},
		{
			input:    `hello|world`,
			expected: `hello\|world`,
		},
		{
			input:    `hello!world`,
			expected: `hello\!world`,
		},
		{
			input:    `hello:world`,
			expected: `hello\:world`,
		},
		{
			input:    `hello*world`,
			expected: `hello\*world`,
		},
		{
			input:    `hello(world)`,
			expected: `hello\(world\)`,
		},
		{
			input:    `hello'world'`,
			expected: `hello''world''`,
		},
		{
			input:    `hello \|& world`,
			expected: `hello \\\|\& world`,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input=%s", tc.input), func(t *testing.T) {
			actual := EscapeSpecialChars(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, actual)
			}
		})
	}
}
