//go:build unit
// +build unit

package parse

import (
	"testing"

	"github.com/FrHorschig/kant-search-backend/core/errors"
	c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"
	"github.com/stretchr/testify/assert"
)

func sp(s string) *string {
	return &s
}

func TestParseInternal(t *testing.T) {
	testCases := []struct {
		name  string
		input []c.Token
		expr  []c.Expression
		err   *errors.Error
	}{
		{
			name: "basic expression",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewClose(),
			},
			expr: []c.Expression{{
				Metadata: c.Metadata{Class: "class"},
			}},
			err: nil,
		},
		{
			name: "three basic expression",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewClose(),
				c.NewOpen(),
				c.NewClass("class2"),
				c.NewClose(),
				c.NewOpen(),
				c.NewClass("class3"),
				c.NewClose(),
			},
			expr: []c.Expression{
				{Metadata: c.Metadata{Class: "class"}},
				{Metadata: c.Metadata{Class: "class2"}},
				{Metadata: c.Metadata{Class: "class3"}},
			},
			err: nil,
		},
		{
			name: "expression with param without content",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewParam("param"),
				c.NewClose(),
			},
			expr: []c.Expression{{
				Metadata: c.Metadata{
					Class: "class",
					Param: sp("param"),
				},
			}},
			err: nil,
		},
		{
			name: "expression with content",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewSeparator(),
				c.NewText("text"),
				c.NewClose(),
			},
			expr: []c.Expression{{
				Metadata: c.Metadata{
					Class: "class",
				},
				Content: sp("text"),
			}},
			err: nil,
		},
		{
			name: "expression with param and nested content",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewParam("param"),
				c.NewSeparator(),
				c.NewOpen(),
				c.NewClass("class2"),
				c.NewClose(),
				c.NewClose(),
			},
			expr: []c.Expression{{
				Metadata: c.Metadata{
					Class: "class",
					Param: sp("param"),
				},
				Content: sp("{class2}"),
			}},
			err: nil,
		},
		{
			name: "closing brace error",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"class"},
			},
		},
		{
			name: "closing brace error with param",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewParam("Location"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"classLocation"},
			},
		},
		{
			name: "closing brace error with short content",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewSeparator(),
				c.NewText("text"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"text"},
			},
		},
		{
			name: "closing brace error with long content",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("class"),
				c.NewSeparator(),
				c.NewText("a very long test text"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"...y long test text"},
			},
		},
		{
			name: "missing class error",
			input: []c.Token{
				c.NewOpen(),
				c.NewClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg: errors.MISSING_EXPR_TYPE,
			},
		},
		{
			name: "unexpected token after expression",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("type"),
				c.NewClose(),
				c.NewClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.UNEXPECTED_TOKEN,
				Params: []string{"}"},
			},
		},
		{
			name: "unexpected token after in nested expression",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("type"),
				c.NewSeparator(),
				c.NewOpen(),
				c.NewParam("param"),
				c.NewClose(),
				c.NewClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg: errors.MISSING_EXPR_TYPE,
			},
		},
		{
			name: "not starting with OPEN",
			input: []c.Token{
				c.NewClass("type"),
				c.NewClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.UNEXPECTED_TOKEN,
				Params: []string{"type"},
			},
		},
		{
			name: "multiple expression with param and nested content",
			input: []c.Token{
				c.NewOpen(),
				c.NewClass("p"),
				c.NewParam("234"),
				c.NewClose(),
				c.NewOpen(),
				c.NewClass("paragraph"),
				c.NewSeparator(),
				c.NewText("some text "),
				c.NewOpen(),
				c.NewClass("l"),
				c.NewParam("2"),
				c.NewClose(),
				c.NewText(" more "),
				c.NewOpen(),
				c.NewClass("p"),
				c.NewParam("324"),
				c.NewClose(),
				c.NewText(" text"),
				c.NewClose(),
			},
			expr: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "paragraph"}, Content: sp("some text {l2} more {p324} text")},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := Parse(tc.input)
			assert.Len(t, tc.expr, len(expr))
			for i, e := range tc.expr {
				if e.Content != nil {
					assert.Equal(t, *e.Content, *expr[i].Content)
				} else {
					assert.Nil(t, expr[i].Content)
				}
				assert.Equal(t, e.Metadata.String(), expr[i].Metadata.String())
			}
			assert.Equal(t, tc.err, err)
		})
	}
}
