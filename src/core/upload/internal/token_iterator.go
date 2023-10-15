package internal

type tokenIterator struct {
	tokens []Token
	index  int
}

func (tk *tokenIterator) hasNext() bool {
	return tk.index < len(tk.tokens)
}

func (tk *tokenIterator) peek() Token {
	if tk.index < len(tk.tokens) {
		return tk.tokens[tk.index]
	}
	return Token{}
}

func (tk *tokenIterator) consume(expected Type) bool {
	if tk.hasNext() && tk.tokens[tk.index].Type == expected {
		tk.index++
		return true
	}
	return false
}

func (tk *tokenIterator) consumeWithText(expected Type) (string, bool) {
	if tk.consume(expected) {
		tk.index--
		text := tk.tokens[tk.index].Text
		tk.index++
		return text, true
	}
	return "", false
}
