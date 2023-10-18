package parse

import c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"

type tokenIterator struct {
	tokens []c.Token
	index  int
}

func (tk *tokenIterator) hasNext() bool {
	return tk.index < len(tk.tokens)
}

func (tk *tokenIterator) peek() c.Token {
	if tk.index < len(tk.tokens) {
		return tk.tokens[tk.index]
	}
	return c.Token{}
}

func (tk *tokenIterator) consume(expected c.Type) bool {
	if tk.hasNext() && tk.tokens[tk.index].Type == expected {
		tk.index++
		return true
	}
	return false
}

func (tk *tokenIterator) consumeWithText(expected c.Type) (string, bool) {
	if tk.consume(expected) {
		tk.index--
		text := tk.tokens[tk.index].Text
		tk.index++
		return text, true
	}
	return "", false
}
