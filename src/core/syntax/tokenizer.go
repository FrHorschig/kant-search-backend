package syntax

import (
	"errors"
	"strings"
)

func tokenize(input string) (*[]Token, error) {
	if wrongBeginChar(input[0]) {
		return nil, errors.New("search input must not start with &, | or )")
	}
	if wrongEndChar(input[len(input)-1]) {
		return nil, errors.New("search input must not end with &, |, ! or (")
	}

	var tokens []Token
	lastIsTerm := false
	for len(input) > 0 {
		switch {
		case strings.HasPrefix(input, "&"):
			tokens = append(tokens, newAnd())
			input = input[1:]
			lastIsTerm = false
		case strings.HasPrefix(input, "|"):
			tokens = append(tokens, newOr())
			input = input[1:]
			lastIsTerm = false
		case strings.HasPrefix(input, "!"):
			tokens = append(tokens, newNot())
			input = input[1:]
			lastIsTerm = false
		case strings.HasPrefix(input, "("):
			tokens = append(tokens, newOpen())
			input = input[1:]
			lastIsTerm = false
		case strings.HasPrefix(input, ")"):
			tokens = append(tokens, newClose())
			input = input[1:]
			lastIsTerm = false
		case strings.HasPrefix(input, "\""):
			end := strings.Index(input[1:], "\"")
			if end == -1 {
				return nil, errors.New("unterminated double quote found")
			}
			tokens = append(tokens, newPhrase(input[1:end]))
			input = input[end:]
			lastIsTerm = true
		default:
			if lastIsTerm {
				tokens = append(tokens, newAnd())
			}
			end := nextOperatorIndex(input)
			if end == -1 {
				tokens = append(tokens, newWord(input))
				input = ""
			} else {
				tokens = append(tokens, newWord(input[0:end]))
				input = input[end:]
			}
			lastIsTerm = true
		}
	}
	_, err := buildAst(tokens)
	if err != nil {
		return &tokens, err
	}
	return &tokens, nil
}

func wrongBeginChar(c byte) bool {
	return c == '&' || c == '|' || c == ')'
}

func wrongEndChar(c byte) bool {
	return c == '&' || c == '|' || c == '!' || c == '('
}

func nextOperatorIndex(s string) int {
	for i, r := range s {
		if strings.ContainsRune(`&|!()"`, r) {
			return i
		}
	}
	return -1
}
