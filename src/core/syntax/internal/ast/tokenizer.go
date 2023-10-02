package ast

import (
	"strings"
)

func Tokenize(searchTerms string) AST {
	var tokens []Token
	i := 0
	lastIsTerm := false
	for _, w := range strings.Split(searchTerms, " ") {
		switch {
		case len(w) == 0:
			continue
		case len(w) == 1:
			tokens, i, lastIsTerm = addOperator(w, tokens, i, lastIsTerm)
		case len(w) > 1:
			tokens, i, lastIsTerm = splitOrAddToken(w, tokens, i, lastIsTerm)
		}
	}

	return AST{} // TODO frhorsch: implement me
}

func addOperator(w string, tokens []Token, i int, lastIsTerm bool) ([]Token, int, bool) {
	switch {
	case w == "&":
		tokens = append(tokens, newAnd(i))
		i++
	case w == "|":
		tokens = append(tokens, newOr(i))
		i++
	case w == "!":
		tokens = append(tokens, newNot(i))
		i++
	case w == "(":
		tokens = append(tokens, newOpen(i))
		i++
	case w == ")":
		tokens = append(tokens, newClose(i))
		i++
	default:
		tokens, i, lastIsTerm = addWord(w, tokens, i, lastIsTerm)
	}
	return tokens, i, lastIsTerm
}

func splitOrAddToken(w string, tokens []Token, i int, lastIsTerm bool) ([]Token, int, bool) {
	switch {
	case strings.HasPrefix(w, "!"):
		tokens = append(tokens, newNot(i))
		i++
		tokens = append(tokens, Token{Text: w[1:], Index: i})
		i++
		lastIsTerm = true
	case strings.HasPrefix(w, "("):
		tokens = append(tokens, newOpen(i))
		i++
		tokens = append(tokens, Token{Text: w[1:], Index: i})
		i++
		lastIsTerm = true
	case strings.HasSuffix(w, ")"):
		tokens = append(tokens, newClose(i))
		i++
		tokens = append(tokens, Token{Text: w[:len(w)-1], Index: i})
		i++
		lastIsTerm = true
	default:
		tokens, i, lastIsTerm = addWord(w, tokens, i, lastIsTerm)
	}
	return tokens, i, lastIsTerm
}

func addWord(w string, tokens []Token, i int, lastIsTerm bool) ([]Token, int, bool) {
	if lastIsTerm {
		tokens = append(tokens, newAnd(i))
		i++
	}
	tokens = append(tokens, newWord(w, i))
	i++
	lastIsTerm = true
	return tokens, i, lastIsTerm
}
