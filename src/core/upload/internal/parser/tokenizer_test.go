//go:build unit
// +build unit

package parser

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
			name:     "type token success",
			input:    "type",
			expected: []Token{newClass("type")},
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
			input:    "{type}|{123.456}",
			expected: []Token{newOpen(), newClass("type"), newClose(), newSeparator(), newOpen(), newParam("123.456"), newClose()},
		},
		{
			name:     "input with spaces between tokens success",
			input:    "{type} {123.456} {a|123xyz}",
			expected: []Token{newOpen(), newClass("type"), newClose(), newOpen(), newParam("123.456"), newClose(), newOpen(), newClass("a"), newSeparator(), newText("123xyz"), newClose()},
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
