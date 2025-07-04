package parse

import (
	"strings"

	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/internal/model"
)

func Tokenize(input string) ([]model.Token, *errors.SyntaxError) {
	input = strings.TrimSpace(input)
	if wrongBeginChar(input[0]) {
		return nil, &errors.SyntaxError{
			Msg:    errors.WrongStartingChar,
			Params: []string{string(input[0])},
		}
	}
	if wrongEndChar(input[len(input)-1]) {
		return nil, &errors.SyntaxError{
			Msg:    errors.WrongEndingChar,
			Params: []string{string(input[len(input)-1])},
		}
	}

	tokens, err := createTokens(input)
	if err != nil {
		return nil, err
	}
	return addInBetweenAnds(tokens), nil
}

func wrongBeginChar(c byte) bool {
	return c == '&' || c == '|' || c == ')'
}

func wrongEndChar(c byte) bool {
	return c == '&' || c == '|' || c == '!' || c == '('
}

func createTokens(input string) ([]model.Token, *errors.SyntaxError) {
	var tokens []model.Token
	for len(input) > 0 {
		switch {
		case strings.HasPrefix(input, "&"):
			tokens = append(tokens, newAnd())
			input = input[1:]
		case strings.HasPrefix(input, "|"):
			tokens = append(tokens, newOr())
			input = input[1:]
		case strings.HasPrefix(input, "!"):
			tokens = append(tokens, newNot())
			input = input[1:]
		case strings.HasPrefix(input, "("):
			tokens = append(tokens, newOpen())
			input = input[1:]
		case strings.HasPrefix(input, ")"):
			tokens = append(tokens, newClose())
			input = input[1:]
		case strings.HasPrefix(input, "\""):
			token, newInput, err := findPhrase(input)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
			input = newInput
		default:
			token, newInput := findWord(input)
			if token != nil {
				tokens = append(tokens, *token)
			}
			input = newInput
		}
		input = strings.TrimSpace(input)
	}
	return tokens, nil
}

func findPhrase(input string) (*model.Token, string, *errors.SyntaxError) {
	var token model.Token
	end := strings.Index(input[1:], "\"")
	if end == -1 {
		return nil, "", &errors.SyntaxError{Msg: errors.UnterminatedDoubleQuote}
	}
	end += 1
	token = newPhrase(strings.TrimSpace(input[1:end]))
	if len(input) > end+1 {
		input = input[end+1:]
	} else {
		input = ""
	}
	return &token, input, nil
}

func findWord(input string) (*model.Token, string) {
	var token model.Token
	end := nextNonWordCharIndex(input)
	if end == -1 {
		token = newWord(input)
		input = ""
	} else {
		word := strings.TrimSpace(input[0:end])
		if len(word) > 0 {
			token = newWord(word)
		}
		input = input[end:]
	}
	return &token, input
}

func nextNonWordCharIndex(s string) int {
	for i, r := range s {
		if strings.ContainsRune(`&|!()" `, r) {
			return i
		}
	}
	return -1
}

func addInBetweenAnds(tokens []model.Token) []model.Token {
	var result []model.Token
	for i, t := range tokens {
		if i == len(tokens)-1 {
			result = append(result, t)
			break
		}
		result = append(result, t)
		if isLhs(t) && isRhs(tokens[i+1]) {
			result = append(result, newAnd())
		}
	}
	return result
}

func isLhs(t model.Token) bool {
	return t.IsWord || t.IsPhrase || t.IsClose
}

func isRhs(t model.Token) bool {
	return t.IsWord || t.IsPhrase || t.IsNot || t.IsOpen
}
