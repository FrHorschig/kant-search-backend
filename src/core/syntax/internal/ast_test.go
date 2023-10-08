//go:build unit
// +build unit

package internal

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSyntax(t *testing.T) {
	testCases := []struct {
		name  string
		input []Token
		err   error
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
			err:   errors.New("unexpected end of input"),
		},
		{
			name: "NOT without following word",
			input: []Token{
				{Text: "!", IsNot: true},
			},
			err: errors.New("unexpected end of input"),
		},
		{
			name: "remaining tokens error",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "world", IsWord: true},
				{Text: "extra", IsWord: true},
			},
			err: fmt.Errorf("unexpected token: extra"),
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
			err: errors.New("unexpected token: )"),
		},
		{
			name: "missing closing parenthesis",
			input: []Token{
				{Text: "(", IsOpen: true},
				{Text: "hello", IsWord: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: errors.New("missing closing parenthesis"),
		},
		{
			name: "consecutive operators",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: errors.New("unexpected token: |"),
		},
		{
			name: "consecutive operators",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "|", IsAnd: true},
				{Text: "&", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: errors.New("unexpected token: &"),
		},
		{
			name: "starts with OR",
			input: []Token{
				{Text: "|", IsOr: true},
				{Text: "world", IsWord: true},
			},
			err: errors.New("unexpected token: |"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckSyntax(&tc.input)
			assert.Equal(t, tc.err, err)
		})
	}
}
