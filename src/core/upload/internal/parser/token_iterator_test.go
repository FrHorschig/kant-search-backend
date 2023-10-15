package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenIteratorPeek(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
		result Token
	}{
		{
			name: "success",
			tokens: []Token{
				{Type: CLASS, Text: "foo"},
				{Type: TEXT, Text: "bar"},
			},
			result: Token{Type: CLASS, Text: "foo"},
		},
		{
			name:   "failure on empty token list",
			tokens: []Token{},
			result: Token{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tk := &tokenIterator{tokens: tc.tokens}
			result := tk.peek()
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestTokenIteratorConsume(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
		input  Type
		ok     bool
	}{
		{
			name: "success",
			tokens: []Token{
				{Type: CLASS, Text: "foo"},
				{Type: TEXT, Text: "bar"},
			},
			input: CLASS,
			ok:    true,
		},
		{
			name: "failure on wrong type",
			tokens: []Token{
				{Type: CLASS, Text: "foo"},
				{Type: TEXT, Text: "bar"},
			},
			input: LOCATION,
			ok:    false,
		},
		{
			name:   "failure on empty token list",
			tokens: []Token{},
			input:  LOCATION,
			ok:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tk := &tokenIterator{tokens: tc.tokens}
			ok := tk.consume(tc.input)
			assert.Equal(t, tc.ok, ok)
		})
	}
}

func TestTokenIteratorConsumeWithText(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
		input  Type
		text   string
		ok     bool
	}{
		{
			name: "success",
			tokens: []Token{
				{Type: CLASS, Text: "foo"},
				{Type: TEXT, Text: "bar"},
			},
			input: CLASS,
			text:  "foo",
			ok:    true,
		},
		{
			name: "failure on wrong type",
			tokens: []Token{
				{Type: CLASS, Text: "foo"},
				{Type: TEXT, Text: "bar"},
			},
			input: LOCATION,
			text:  "",
			ok:    false,
		},
		{
			name:   "failure on empty token list",
			tokens: []Token{},
			input:  LOCATION,
			text:   "",
			ok:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tk := &tokenIterator{tokens: tc.tokens}
			text, ok := tk.consumeWithText(tc.input)
			assert.Equal(t, tc.text, text)
			assert.Equal(t, tc.ok, ok)
		})
	}
}
