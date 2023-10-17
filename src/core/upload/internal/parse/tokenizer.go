package parse

import (
	"regexp"
	"strings"
)

func tokenize(input string) []Token {
	rType := regexp.MustCompile(`^(paragraph|heading|footnote|fn|p|l)([|0-9])`)
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
			if matches := rType.FindStringSubmatch(input); matches != nil {
				tokens = append(tokens, newClass(matches[1]))
				input = input[len(matches[1]):]
			} else if match := rLoc.FindString(input); match != "" {
				tokens = append(tokens, newParam(match[:len(match)-1]))
				input = input[len(match)-1:]
			} else {
				match := rChar.FindString(input)
				tokens = append(tokens, newText(strings.TrimSpace(match)))
				input = input[len(match):]
			}
		}
		input = strings.TrimSpace(input)
	}
	return tokens
}
