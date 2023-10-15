package parser

import (
	"regexp"
	"strings"
)

func tokenize(input string) []Token {
	rType := regexp.MustCompile(`^[a-z]+`)
	rLoc := regexp.MustCompile(`^\d+(\.\d+)?[}|]`)
	rChar := regexp.MustCompile(`^[^{}]*`)

	var tokens []Token
	for len(input) > 0 {
		switch {
		case strings.HasPrefix(input, "{"):
			tokens = append(tokens, newOpen())
			input = input[1:]
		case strings.HasPrefix(input, "}"):
			tokens = append(tokens, newClose())
			input = input[1:]
		case strings.HasPrefix(input, "|"):
			tokens = append(tokens, newSeparator())
			input = input[1:]
		default:
			if match := rType.FindString(input); match != "" {
				tokens = append(tokens, newClass(match))
				input = input[len(match):]
			} else if match := rLoc.FindString(input); match != "" {
				tokens = append(tokens, newParam(match[:len(match)-1]))
				input = input[len(match)-1:]
			} else {
				match := rChar.FindString(input)
				tokens = append(tokens, newText(match))
				input = input[len(match):]
			}
		}
		input = strings.TrimSpace(input)
	}
	return tokens
}
