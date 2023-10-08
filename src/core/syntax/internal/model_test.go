package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeSpecialChars(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    `hello`,
			expected: `hello`,
		},
		{
			input:    `hello\world`,
			expected: `hello\\world`,
		},
		{
			input:    `hello&world`,
			expected: `hello\&world`,
		},
		{
			input:    `hello|world`,
			expected: `hello\|world`,
		},
		{
			input:    `hello!world`,
			expected: `hello\!world`,
		},
		{
			input:    `hello:world`,
			expected: `hello\:world`,
		},
		{
			input:    `hello*world`,
			expected: `hello\*world`,
		},
		{
			input:    `hello(world)`,
			expected: `hello\(world\)`,
		},
		{
			input:    `hello'world'`,
			expected: `hello''world''`,
		},
		{
			input:    `hello \|& world`,
			expected: `hello \\\|\& world`,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input=%s", tc.input), func(t *testing.T) {
			actual := escapeSpecialChars(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, actual)
			}
		})
	}
}

func TestNewAnd(t *testing.T) {
	actual := newAnd()
	assert.True(t, actual.IsAnd)
	assert.Equal(t, "&", actual.Text)
}

func TestNewOr(t *testing.T) {
	actual := newOr()
	assert.True(t, actual.IsOr)
	assert.Equal(t, "|", actual.Text)
}

func TestNewNot(t *testing.T) {
	actual := newNot()
	assert.True(t, actual.IsNot)
	assert.Equal(t, "!", actual.Text)
}

func TestNewOpen(t *testing.T) {
	actual := newOpen()
	assert.True(t, actual.IsOpen)
	assert.Equal(t, "(", actual.Text)
}

func TestNewClose(t *testing.T) {
	actual := newClose()
	assert.True(t, actual.IsClose)
	assert.Equal(t, ")", actual.Text)
}

func TestNewPhrase(t *testing.T) {
	actual := newPhrase("phrase")
	assert.True(t, actual.IsPhrase)
	assert.Equal(t, "phrase", actual.Text)
}

func TestNewWord(t *testing.T) {
	actual := newWord("word")
	assert.True(t, actual.IsWord)
	assert.Equal(t, "word", actual.Text)
}

func TestGetSearchString(t *testing.T) {
	testCases := []struct {
		name     string
		input    []Token
		expected string
	}{
		{
			name:     "empty input",
			input:    []Token{},
			expected: "",
		},
		{
			name: "words input",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "world", IsWord: true},
			},
			expected: "hello world",
		},
		{
			name: "phrase input",
			input: []Token{
				{Text: "hello world goodbye world", IsPhrase: true},
			},
			expected: "(hello <-> world <-> goodbye <-> world)",
		},
		{
			name: "words, phrase and operators",
			input: []Token{
				{Text: "hello", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "(", IsOpen: true},
				{Text: "world", IsWord: true},
				{Text: "|", IsOr: true},
				{Text: "goodbye", IsWord: true},
				{Text: "&", IsAnd: true},
				{Text: "!", IsNot: true},
				{Text: "world", IsWord: true},
				{Text: ")", IsClose: true},
				{Text: "&", IsAnd: true},
				{Text: "hello again world", IsPhrase: true},
			},
			expected: "hello & ( world | goodbye & ! world ) & (hello <-> again <-> world)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetSearchString(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
