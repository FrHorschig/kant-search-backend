package parse

import (
	"testing"

	c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestTokenIteratorPeek(t *testing.T) {
	tests := []struct {
		name   string
		tokens []c.Token
		result c.Token
	}{
		{
			name: "success",
			tokens: []c.Token{
				{Type: c.CLASS, Text: "foo"},
				{Type: c.TEXT, Text: "bar"},
			},
			result: c.Token{Type: c.CLASS, Text: "foo"},
		},
		{
			name:   "failure on empty token list",
			tokens: []c.Token{},
			result: c.Token{},
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
		tokens []c.Token
		input  c.Type
		ok     bool
	}{
		{
			name: "success",
			tokens: []c.Token{
				{Type: c.CLASS, Text: "foo"},
				{Type: c.TEXT, Text: "bar"},
			},
			input: c.CLASS,
			ok:    true,
		},
		{
			name: "failure on wrong type",
			tokens: []c.Token{
				{Type: c.CLASS, Text: "foo"},
				{Type: c.TEXT, Text: "bar"},
			},
			input: c.PARAM,
			ok:    false,
		},
		{
			name:   "failure on empty token list",
			tokens: []c.Token{},
			input:  c.PARAM,
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
		tokens []c.Token
		input  c.Type
		text   string
		ok     bool
	}{
		{
			name: "success",
			tokens: []c.Token{
				{Type: c.CLASS, Text: "foo"},
				{Type: c.TEXT, Text: "bar"},
			},
			input: c.CLASS,
			text:  "foo",
			ok:    true,
		},
		{
			name: "failure on wrong type",
			tokens: []c.Token{
				{Type: c.CLASS, Text: "foo"},
				{Type: c.TEXT, Text: "bar"},
			},
			input: c.PARAM,
			text:  "",
			ok:    false,
		},
		{
			name:   "failure on empty token list",
			tokens: []c.Token{},
			input:  c.PARAM,
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
