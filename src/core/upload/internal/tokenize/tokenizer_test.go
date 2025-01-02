//go:build unit
// +build unit

package tokenize

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/common/errors"
	c "github.com/frhorschig/kant-search-backend/core/upload/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  []c.Token
		err  *errors.Error
	}{
		{
			name: "open token success",
			in:   "{",
			out:  []c.Token{c.NewOpen()},
			err:  nil,
		},
		{
			name: "close token success",
			in:   "{}",
			out:  []c.Token{c.NewOpen(), c.NewClose()},
			err:  nil,
		},
		{
			name: "separator token success",
			in:   "{|",
			out:  []c.Token{c.NewOpen(), c.NewSeparator()},
			err:  nil,
		},
		{
			name: "class token success",
			in:   "{paragraph|",
			out:  []c.Token{c.NewOpen(), c.NewClass("paragraph"), c.NewSeparator()},
			err:  nil,
		},
		{
			name: "param token ending with close success",
			in:   "{123.456}",
			out:  []c.Token{c.NewOpen(), c.NewParam("123.456"), c.NewClose()},
			err:  nil,
		},
		{
			name: "param token ending with separator success",
			in:   "{123.456|",
			out:  []c.Token{c.NewOpen(), c.NewParam("123.456"), c.NewSeparator()},
			err:  nil,
		},
		{
			name: "char token success",
			in:   `{123abc()\n[]<i> Text</i>`,
			out:  []c.Token{c.NewOpen(), c.NewText(`123abc()\n[]<i> Text</i>`)},
			err:  nil,
		},
		{
			name: "multiple tokens success",
			in:   "{p2}{fn123.456|text}",
			out: []c.Token{
				c.NewOpen(),
				c.NewClass("p"),
				c.NewParam("2"),
				c.NewClose(),
				c.NewOpen(),
				c.NewClass("fn"),
				c.NewParam("123.456"),
				c.NewSeparator(),
				c.NewText("text"),
				c.NewClose()},
			err: nil,
		},
		{
			name: "input with spaces between tokens success",
			in:   "{p2} {fn123.456| text}",
			out: []c.Token{
				c.NewOpen(),
				c.NewClass("p"),
				c.NewParam("2"),
				c.NewClose(),
				c.NewOpen(),
				c.NewClass("fn"),
				c.NewParam("123.456"),
				c.NewSeparator(),
				c.NewText("text"),
				c.NewClose()},
			err: nil,
		},
		{
			name: "wrong starting char",
			in:   "}paragraph|some text}",
			out:  nil,
			err: &errors.Error{
				Msg:    errors.UPLOAD_WRONG_STARTING_CHAR,
				Params: []string{"}"},
			},
		},
		{
			name: "multiple expression with params and nested content",
			in:   "{p234} {paragraph|some text {l2} more {p324} text}",
			out: []c.Token{
				c.NewOpen(),
				c.NewClass("p"),
				c.NewParam("234"),
				c.NewClose(),
				c.NewOpen(),
				c.NewClass("paragraph"),
				c.NewSeparator(),
				c.NewText("some text"),
				c.NewOpen(),
				c.NewClass("l"),
				c.NewParam("2"),
				c.NewClose(),
				c.NewText("more"),
				c.NewOpen(),
				c.NewClass("p"),
				c.NewParam("324"),
				c.NewClose(),
				c.NewText("text"),
				c.NewClose(),
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Tokenize(tc.in)
			assert.Len(t, out, len(tc.out))
			for i := range tc.out {
				assert.Equal(t, tc.out[i], out[i])
			}
			assert.Equal(t, err, tc.err)
		})
	}
}
