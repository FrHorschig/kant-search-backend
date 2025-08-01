//go:build unit
// +build unit

package parse

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []model.Token
		err      *errors.SyntaxError
	}{
		{
			name:     "only words success",
			input:    "hello world",
			expected: []model.Token{newWord("hello"), newAnd(), newWord("world")},
			err:      nil,
		},
		{
			name:     "AND success",
			input:    "hello & world",
			expected: []model.Token{newWord("hello"), newAnd(), newWord("world")},
			err:      nil,
		},
		{
			name:     "AND success",
			input:    "hello | world",
			expected: []model.Token{newWord("hello"), newOr(), newWord("world")},
			err:      nil,
		},
		{
			name:     "AND plus OR success",
			input:    "hello & world | kant",
			expected: []model.Token{newWord("hello"), newAnd(), newWord("world"), newOr(), newWord("kant")},
			err:      nil,
		},
		{
			name:     "NOT success",
			input:    "!world",
			expected: []model.Token{newNot(), newWord("world")},
			err:      nil,
		},
		{
			name:     "NOT plus space success",
			input:    "hello ! world",
			expected: []model.Token{newWord("hello"), newAnd(), newNot(), newWord("world")},
			err:      nil,
		},
		{
			name:     "NOT plus AND success",
			input:    "hello &! world",
			expected: []model.Token{newWord("hello"), newAnd(), newNot(), newWord("world")},
			err:      nil,
		},
		{
			name:     "parentheses success",
			input:    "hello (world)",
			expected: []model.Token{newWord("hello"), newAnd(), newOpen(), newWord("world"), newClose()},
			err:      nil,
		},
		{
			name:     "parentheses with spaces success",
			input:    "hello ( world )",
			expected: []model.Token{newWord("hello"), newAnd(), newOpen(), newWord("world"), newClose()},
			err:      nil,
		},
		{
			name:     "phrase success",
			input:    "hello \"you\" world",
			expected: []model.Token{newWord("hello"), newAnd(), newPhrase("you"), newAnd(), newWord("world")},
			err:      nil,
		},
		{
			name:     "starts with phrase success",
			input:    "\"hello\" world",
			expected: []model.Token{newPhrase("hello"), newAnd(), newWord("world")},
			err:      nil,
		},
		{
			name:     "ends with phrase success",
			input:    "hello \"world\"",
			expected: []model.Token{newWord("hello"), newAnd(), newPhrase("world")},
			err:      nil,
		},
		{
			name:     "starts with AND error",
			input:    "& hello",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongStartingChar, Params: []string{"&"}},
		},
		{
			name:     "starts with OR error",
			input:    "| hello",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongStartingChar, Params: []string{"|"}},
		},
		{
			name:     "starts with CloseParen error",
			input:    ") hello",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongStartingChar, Params: []string{")"}},
		},
		{
			name:     "ends with AND error",
			input:    "hello &",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongEndingChar, Params: []string{"&"}},
		},
		{
			name:     "ends with OR error",
			input:    "hello |",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongEndingChar, Params: []string{"|"}},
		},
		{
			name:     "ends with NOT error",
			input:    "hello !",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongEndingChar, Params: []string{"!"}},
		},
		{
			name:     "ends with OpenParen error",
			input:    "hello (",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.WrongEndingChar, Params: []string{"("}},
		},
		{
			name:     "unterminated double quote error",
			input:    "hello \"world",
			expected: nil,
			err:      &errors.SyntaxError{Msg: errors.UnterminatedDoubleQuote},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Tokenize(tc.input)
			if actual != nil {
				assert.Len(t, actual, len(tc.expected))
				for i := range tc.expected {
					assert.Equal(t, tc.expected[i], actual[i])
				}
			}
			assert.Equal(t, tc.err, err)
		})
	}
}
