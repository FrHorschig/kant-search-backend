package syntax

import "strings"

func tokenize(searchTerms string) ([]operator, []parenthesis, []word, int) {
	var ops []operator
	var parens []parenthesis
	var words []word
	i := 0
	lastIsTerm := false
	for _, w := range strings.Split(searchTerms, " ") {
		switch {
		case len(w) == 0:
			continue
		case len(w) == 1:
			ops, parens, words, i, lastIsTerm = addOperator(w, ops, parens, words, i, lastIsTerm)
		case len(w) > 1:
			ops, parens, words, i, lastIsTerm = splitOrAddWord(w, ops, parens, words, i, lastIsTerm)
		}
	}
	return ops, parens, words, i - 1
}

func addOperator(w string, ops []operator, parens []parenthesis, words []word, i int, lastIsTerm bool) ([]operator, []parenthesis, []word, int, bool) {
	switch {
	case w == "&":
		ops = append(ops, newAnd(i))
		i++
	case w == "|":
		ops = append(ops, newOr(i))
		i++
	case w == "!":
		ops = append(ops, newNot(i))
		i++
	case w == "(":
		parens = append(parens, newOpen(i))
		i++
	case w == ")":
		parens = append(parens, newClose(i))
		i++
	default:
		ops, words, i, lastIsTerm = addWord(w, ops, words, i, lastIsTerm)
	}
	return ops, parens, words, i, lastIsTerm
}

func splitOrAddWord(w string, ops []operator, parens []parenthesis, words []word, i int, lastIsTerm bool) ([]operator, []parenthesis, []word, int, bool) {
	switch {
	case strings.HasPrefix(w, "!"):
		ops = append(ops, newNot(i))
		i++
		words = append(words, word{w[1:], i})
		i++
		lastIsTerm = true
	case strings.HasPrefix(w, "("):
		parens = append(parens, newOpen(i))
		i++
		words = append(words, word{w[1:], i})
		i++
		lastIsTerm = true
	case strings.HasSuffix(w, ")"):
		parens = append(parens, newClose(i))
		i++
		words = append(words, word{w[:len(w)-1], i})
		i++
		lastIsTerm = true
	default:
		ops, words, i, lastIsTerm = addWord(w, ops, words, i, lastIsTerm)
	}
	return ops, parens, words, i, lastIsTerm
}

func addWord(w string, ops []operator, words []word, i int, lastIsTerm bool) ([]operator, []word, int, bool) {
	if lastIsTerm {
		ops = append(ops, operator{isAnd: true, index: i})
		i++
	}
	words = append(words, word{w, i})
	i++
	lastIsTerm = true
	return ops, words, i, lastIsTerm
}
