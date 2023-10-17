//go:build unit
// +build unit

package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  []Token
	}{
		{name: "open token success", in: "{", out: []Token{newOpen()}},
		{name: "close token success", in: "}", out: []Token{newClose()}},
		{name: "separator token success", in: "|", out: []Token{newSeparator()}},
		{name: "class token success", in: "paragraph|", out: []Token{newClass("paragraph"), newSeparator()}},
		{name: "param token ending with close success", in: "123.456}", out: []Token{newParam("123.456"), newClose()}},
		{name: "param token ending with separator success", in: "123.456|", out: []Token{newParam("123.456"), newSeparator()}},
		{name: "char token success", in: `123abc()\n[]<i> Text</i>`, out: []Token{newText(`123abc()\n[]<i> Text</i>`)}},
		{
			name: "multiple tokens success",
			in:   "{p2}{fn123.456|text}",
			out: []Token{newOpen(),
				newClass("p"),
				newParam("2"),
				newClose(),
				newOpen(),
				newClass("fn"),
				newParam("123.456"),
				newSeparator(),
				newText("text"),
				newClose()},
		},
		{
			name: "input with spaces between tokens success",
			in:   "{p2} {fn123.456| text}",
			out: []Token{newOpen(),
				newClass("p"),
				newParam("2"),
				newClose(),
				newOpen(),
				newClass("fn"),
				newParam("123.456"),
				newSeparator(),
				newText("text"),
				newClose()},
		},
		{
			name: "multiple expression with params and nested content",
			in:   "{p234} {paragraph|some text {l2} more {p324} text}",
			out: []Token{
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
			actual := tokenize(tc.in)
			assert.Len(t, actual, len(tc.out))
			for i := range tc.out {
				assert.Equal(t, tc.out[i], actual[i])
			}
		})
	}
}
