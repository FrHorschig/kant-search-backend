//go:build unit
// +build unit

package parser

import (
	"testing"

	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseInternal(t *testing.T) {
	testCases := []struct {
		name  string
		input []Token
		expr  *Expression
		err   *errors.Error
	}{
		{
			name: "basic expression without content",
			input: []Token{
				newOpen(),
				newClass("class"),
				newClose(),
			},
			expr: &Expression{
				Metadata: Metadata{
					Class: "class",
				},
			},
			err: nil,
		},
		{
			name: "expression with param without content",
			input: []Token{
				newOpen(),
				newClass("class"),
				newParam("param"),
				newClose(),
			},
			expr: &Expression{
				Metadata: Metadata{
					Class:    "class",
					Location: &[]string{"param"}[0],
				},
			},
			err: nil,
		},
		{
			name: "expression with content",
			input: []Token{
				newOpen(),
				newClass("class"),
				newSeparator(),
				newText("text"),
				newClose(),
			},
			expr: &Expression{
				Metadata: Metadata{
					Class: "class",
				},
				Content: &Content{
					Texts: []string{"text"},
				},
			},
			err: nil,
		},
		{
			name: "expression with param and nested content",
			input: []Token{
				newOpen(),
				newClass("class"),
				newParam("param"),
				newSeparator(),
				newOpen(),
				newClass("class2"),
				newClose(),
				newClose(),
			},
			expr: &Expression{
				Metadata: Metadata{
					Class:    "class",
					Location: &[]string{"param"}[0],
				},
				Content: &Content{
					Expressions: []*Expression{
						{
							Metadata: Metadata{
								Class: "class2",
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "closing brace error",
			input: []Token{
				newOpen(),
				newClass("class"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"class"},
			},
		},
		{
			name: "closing brace error with param",
			input: []Token{
				newOpen(),
				newClass("class"),
				newParam("Location"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"classLocation"},
			},
		},
		{
			name: "closing brace error with content",
			input: []Token{
				newOpen(),
				newClass("class"),
				newSeparator(),
				newText("text"),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.MISSING_CLOSING_BRACE,
				Params: []string{"text"},
			},
		},
		{
			name: "missing class error",
			input: []Token{
				newOpen(),
				newClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg: errors.MISSING_EXPR_TYPE,
			},
		},
		{
			name: "unexpected token after expression",
			input: []Token{
				newOpen(),
				newClass("type"),
				newClose(),
				newClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.UNEXPECTED_TOKEN,
				Params: []string{"}"},
			},
		},
		{
			name: "unexpected token after in nested expression",
			input: []Token{
				newOpen(),
				newClass("type"),
				newSeparator(),
				newOpen(),
				newParam("param"),
				newClose(),
				newClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg: errors.MISSING_EXPR_TYPE,
			},
		},
		{
			name: "not starting with OPEN",
			input: []Token{
				newClass("type"),
				newClose(),
			},
			expr: nil,
			err: &errors.Error{
				Msg:    errors.UNEXPECTED_TOKEN,
				Params: []string{"type"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parse(tc.input)

			if tc.expr != nil && expr != nil {
				assert.Equal(t, tc.expr.Content, expr.Content)
				assert.Equal(t, tc.expr.Metadata, expr.Metadata)
			}
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestParsePublic(t *testing.T) {
	// TODO frhorsch
}
