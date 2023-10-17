//go:build unit
// +build unit

package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "open token success",
			input:    "{",
			expected: []Token{newOpen()},
		},
		{
			name:     "close token success",
			input:    "}",
			expected: []Token{newClose()},
		},
		{
			name:     "separator token success",
			input:    "|",
			expected: []Token{newSeparator()},
		},
		{
			name:     "class token success",
			input:    "paragraph|",
			expected: []Token{newClass("paragraph"), newSeparator()},
		},
		{
			name:     "param token ending with close success",
			input:    "123.456}",
			expected: []Token{newParam("123.456"), newClose()},
		},
		{
			name:     "param token ending with separator success",
			input:    "123.456|",
			expected: []Token{newParam("123.456"), newSeparator()},
		},
		{
			name:     "char token success",
			input:    `123abc()\n[]<i> Text</i>`,
			expected: []Token{newText(`123abc()\n[]<i> Text</i>`)},
		},
		{
			name:     "multiple tokens success",
			input:    "{p2}{fn123.456|text}",
			expected: []Token{newOpen(), newClass("p"), newParam("2"), newClose(), newOpen(), newClass("fn"), newParam("123.456"), newSeparator(), newText("text"), newClose()},
		},
		{
			name:     "input with spaces between tokens success",
			input:    "{p2} {fn123.456| text}",
			expected: []Token{newOpen(), newClass("p"), newParam("2"), newClose(), newOpen(), newClass("fn"), newParam("123.456"), newSeparator(), newText("text"), newClose()},
		},
		{
			name:  "multiple expression with params and nested content",
			input: "{p234} {paragraph|some text {l2} more {p324} text}",
			expected: []Token{
				newOpen(),
				newClass("p"),
				newParam("234"),
				newClose(),
				newOpen(),
				newClass("paragraph"),
				newSeparator(),
				newText("some text"),
				newOpen(),
				newClass("l"),
				newParam("2"),
				newClose(),
				newText("more"),
				newOpen(),
				newClass("p"),
				newParam("324"),
				newClose(),
				newText("text"),
				newClose(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tokenize(tc.input)
			assert.Len(t, actual, len(tc.expected))
			for i := range tc.expected {
				assert.Equal(t, tc.expected[i], actual[i])
			}
		})
	}
}
