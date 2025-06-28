//go:build unit
// +build unit

package internal

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAstParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sut := NewAstParser()

	tests := []struct {
		name     string
		input    string
		expected *model.SearchTermNode
	}{
		{
			name:  "Multiple space separated words",
			input: "hello parsing world",
			expected: &model.SearchTermNode{
				Token: newAnd(),
				Left: &model.SearchTermNode{
					Token: newWord("hello"),
				},
				Right: &model.SearchTermNode{
					Token: newAnd(),
					Left:  &model.SearchTermNode{Token: newWord("parsing")},
					Right: &model.SearchTermNode{Token: newWord("world")},
				},
			},
		},
		{
			name:  "Simple AND",
			input: "hello & world",
			expected: &model.SearchTermNode{
				Token: newAnd(),
				Left:  &model.SearchTermNode{Token: newWord("hello")},
				Right: &model.SearchTermNode{Token: newWord("world")},
			},
		},
		{
			name:  "Simple OR",
			input: "hello & world",
			expected: &model.SearchTermNode{
				Token: newOr(),
				Left:  &model.SearchTermNode{Token: newWord("hello")},
				Right: &model.SearchTermNode{Token: newWord("world")},
			},
		},
		{
			name:  "Simple NOT",
			input: "!hello",
			expected: &model.SearchTermNode{
				Token: newNot(),
				Left:  &model.SearchTermNode{Token: newWord("hello")},
			},
		},
		{
			name:  "Complex search query",
			input: "(dog | cat) & !mouse & \"night bird\"",
			expected: &model.SearchTermNode{
				Token: newAnd(),
				Left: &model.SearchTermNode{
					Token: newAnd(),
					Left: &model.SearchTermNode{
						Token: newOr(),
						Left:  &model.SearchTermNode{Token: newWord("dog")},
						Right: &model.SearchTermNode{Token: newWord("cat")},
					},
					Right: &model.SearchTermNode{
						Token: newNot(),
						Left:  &model.SearchTermNode{Token: newWord("mouse")},
					},
				},
				Right: &model.SearchTermNode{Token: newPhrase("night bird")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := sut.Parse(tc.input)
			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	}
}

func newAnd() *model.Token {
	return &model.Token{IsAnd: true, Text: "&"}
}
func newOr() *model.Token {
	return &model.Token{IsOr: true, Text: "|"}
}
func newNot() *model.Token {
	return &model.Token{IsNot: true, Text: "!"}
}
func newWord(text string) *model.Token {
	return &model.Token{IsWord: true, Text: text}
}
func newPhrase(text string) *model.Token {
	return &model.Token{IsPhrase: true, Text: text}
}
