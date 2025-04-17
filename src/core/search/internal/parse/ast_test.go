//go:build unit
// +build unit

package parse

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestCheckSyntax(t *testing.T) {
	testCases := []struct {
		name  string
		input []model.Token
		err   *errors.SyntaxError
	}{
		{
			name: "success",
			input: []model.Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "world", IsWord: true},
			},
			err: nil,
		},
		{
			name:  "empty input",
			input: []model.Token{},
			err:   &errors.SyntaxError{Msg: errors.UnexpectedEndOfInput},
		},
		{
			name: "NOT without following word",
			input: []model.Token{
				{Text: "!", IsNot: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedEndOfInput},
		},
		{
			name: "remaining tokens error",
			input: []model.Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "world", IsWord: true},
				{Text: "extra", IsWord: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedToken, Params: []string{"extra"}},
		},
		{
			name: "phrase OR word",
			input: []model.Token{
				{Text: "\"hello world\"", IsPhrase: true},
				{Text: "|", IsOr: true},
				{Text: "friend", IsWord: true},
			},
			err: nil,
		},
		{
			name: "escaped special characters in phrase",
			input: []model.Token{
				{Text: "\"hello \\\"world\\\"!\"", IsPhrase: true},
			},
			err: nil,
		},
		{
			name: "NOT word",
			input: []model.Token{
				{Text: "!", IsNot: true},
				{Text: "enemy", IsWord: true},
			},
			err: nil,
		},
		{
			name: "grouped expression",
			input: []model.Token{
				{Text: "(", IsOpen: true},
				{Text: "hello", IsWord: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
				{Text: ")", IsClose: true},
				{Text: "&", IsAnd: true},
				{Text: "friend", IsWord: true},
			},
			err: nil,
		},
		{
			name: "nested groups",
			input: []model.Token{
				{Text: "(", IsOpen: true},
				{Text: "(", IsOpen: true},
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "world", IsWord: true},
				{Text: ")", IsClose: true},
				{Text: "|", IsOr: true},
				{Text: "universe", IsWord: true},
				{Text: ")", IsClose: true},
			},
			err: nil,
		},
		{
			name: "empty expression in parenthesis",
			input: []model.Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "(", IsOpen: true},
				{Text: ")", IsClose: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedToken, Params: []string{")"}},
		},
		{
			name: "missing closing parenthesis",
			input: []model.Token{
				{Text: "(", IsOpen: true},
				{Text: "hello", IsWord: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.SyntaxError{Msg: errors.MissingCloseParenthesis},
		},
		{
			name: "OR following AND",
			input: []model.Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedToken, Params: []string{"|"}},
		},
		{
			name: "AND following OR",
			input: []model.Token{
				{Text: "hello", IsWord: true},
				{Text: "|", IsAnd: true},
				{Text: "&", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedToken, Params: []string{"&"}},
		},
		{
			name: "starts with OR",
			input: []model.Token{
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedToken, Params: []string{"|"}},
		},
		{
			name: "starts with AND",
			input: []model.Token{
				{Text: "&", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedToken, Params: []string{"&"}},
		},
		{
			name: "ends with OR",
			input: []model.Token{
				{Text: "world", IsWord: true},
				{Text: "|", IsOr: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedEndOfInput},
		},
		{
			name: "ends with AND",
			input: []model.Token{
				{Text: "world", IsWord: true},
				{Text: "&", IsOr: true},
			},
			err: &errors.SyntaxError{Msg: errors.UnexpectedEndOfInput},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.input)
			assert.Equal(t, tc.err, err)
		})
	}
}
