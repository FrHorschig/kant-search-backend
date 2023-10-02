package ast

import "strings"

type Token struct {
	IsBiOp   bool
	IsMonoOp bool
	IsOpen   bool
	IsClose  bool
	IsWord   bool
	IsPhrase bool

	Text  string
	Index int
}

func newAnd(index int) Token {
	return Token{IsBiOp: true, Text: "&", Index: index}
}
func newOr(index int) Token {
	return Token{IsBiOp: true, Text: "|", Index: index}
}
func newNot(index int) Token {
	return Token{IsMonoOp: true, Text: "!", Index: index}
}
func newOpen(index int) Token {
	return Token{IsOpen: true, Text: "(", Index: index}
}
func newClose(index int) Token {
	return Token{IsClose: true, Text: ")", Index: index}
}
func newWord(text string, index int) Token {
	return Token{IsWord: true, Text: text, Index: index}
}
func newPhrase(text string, index int) Token {
	return Token{IsPhrase: true, Text: text, Index: index}
}

type Node struct {
	IsExpression bool
	IsTerm       bool
	IsFactor     bool

	Left  *Node
	Right *Node
	Token Token
}

type AST struct {
	tokens []Token
	Root   *Node
}

func (ast *AST) GetSearchString() string {
	var builder strings.Builder
	for _, token := range ast.tokens {
		if token.IsWord || token.IsPhrase {
			token.Text = escapeSpecialChars(token.Text)
		}
		builder.WriteString(token.Text)
		builder.WriteString(" ")
	}
	return builder.String()
}
func escapeSpecialChars(input string) string {
	input = strings.ReplaceAll(input, `\`, `\\`)
	replacements := map[string]string{
		`&`: `\&`,
		`|`: `\|`,
		`!`: `\!`,
		`(`: `\(`,
		`)`: `\)`,
		`:`: `\:`,
		`*`: `\*`,
		`'`: `''`,
	}
	for char, replacement := range replacements {
		input = strings.ReplaceAll(input, char, replacement)
	}
	return input
}
