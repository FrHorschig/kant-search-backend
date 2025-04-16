package internal

import "strings"

type Token struct {
	IsAnd    bool
	IsOr     bool
	IsNot    bool
	IsOpen   bool
	IsClose  bool
	IsWord   bool
	IsPhrase bool
	Text     string
}

func newAnd() Token {
	return Token{IsAnd: true, Text: "&"}
}
func newOr() Token {
	return Token{IsOr: true, Text: "|"}
}
func newNot() Token {
	return Token{IsNot: true, Text: "!"}
}
func newOpen() Token {
	return Token{IsOpen: true, Text: "("}
}
func newClose() Token {
	return Token{IsClose: true, Text: ")"}
}
func newWord(text string) Token {
	return Token{IsWord: true, Text: text}
}
func newPhrase(text string) Token {
	return Token{IsPhrase: true, Text: text}
}

type astNode struct {
	Left  *astNode
	Right *astNode
	Token *Token
}

func GetSearchString(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}
	var builder strings.Builder
	for i, token := range tokens {
		if token.IsWord || token.IsPhrase {
			token.Text = escapeSpecialChars(token.Text)
		}
		if token.IsPhrase {
			builder.WriteString("(" + createPhrase(token.Text) + ")")
		} else {
			builder.WriteString(token.Text)
			if i < len(tokens)-1 {
				builder.WriteString(" ")
			}
		}
	}
	return builder.String()
}

func createPhrase(text string) string {
	var builder strings.Builder
	words := strings.Split(text, " ")
	for i, word := range words {
		builder.WriteString(word)
		if i < len(words)-1 {
			builder.WriteString(" <-> ")
		}
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
