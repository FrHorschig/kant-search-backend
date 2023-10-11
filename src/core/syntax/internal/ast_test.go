//go:build unit
// +build unit

package internal

import (
	"testing"

	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestCheckSyntax(t *testing.T) {
	testCases := []struct {
		name  string
		input []Token
		err   *errors.Error
	}{
		{
			name: "success",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "world", IsWord: true},
			},
			err: nil,
		},
		{
			name:  "empty input",
			input: []Token{},
			err:   &errors.Error{Msg: errors.UNEXPECTED_END_OF_INPUT},
		},
		{
			name: "NOT without following word",
			input: []Token{
				{Text: "!", IsNot: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_END_OF_INPUT},
		},
		{
			name: "remaining tokens error",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "world", IsWord: true},
				{Text: "extra", IsWord: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_TOKEN, Params: []string{"extra"}},
		},
		{
			name: "phrase OR word",
			input: []Token{
				{Text: "\"hello world\"", IsPhrase: true},
				{Text: "|", IsOr: true},
				{Text: "friend", IsWord: true},
			},
			err: nil,
		},
		{
			name: "escaped special characters in phrase",
			input: []Token{
				{Text: "\"hello \\\"world\\\"!\"", IsPhrase: true},
			},
			err: nil,
		},
		{
			name: "NOT word",
			input: []Token{
				{Text: "!", IsNot: true},
				{Text: "enemy", IsWord: true},
			},
			err: nil,
		},
		{
			name: "grouped expression",
			input: []Token{
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
			input: []Token{
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
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "(", IsOpen: true},
				{Text: ")", IsClose: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_TOKEN, Params: []string{")"}},
		},
		{
			name: "missing closing parenthesis",
			input: []Token{
				{Text: "(", IsOpen: true},
				{Text: "hello", IsWord: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.Error{Msg: errors.MISSING_CLOSING_PARENTHESIS},
		},
		{
			name: "OR following AND",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_TOKEN, Params: []string{"|"}},
		},
		{
			name: "AND following OR",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "|", IsAnd: true},
				{Text: "&", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_TOKEN, Params: []string{"&"}},
		},
		{
			name: "starts with OR",
			input: []Token{
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_TOKEN, Params: []string{"|"}},
		},
		{
			name: "starts with AND",
			input: []Token{
				{Text: "&", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_TOKEN, Params: []string{"&"}},
		},
		{
			name: "ends with OR",
			input: []Token{
				{Text: "world", IsWord: true},
				{Text: "|", IsOr: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_END_OF_INPUT},
		},
		{
			name: "ends with AND",
			input: []Token{
				{Text: "world", IsWord: true},
				{Text: "&", IsOr: true},
			},
			err: &errors.Error{Msg: errors.UNEXPECTED_END_OF_INPUT},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckSyntax(tc.input)
			assert.Equal(t, tc.err, err)
		})
	}
}
