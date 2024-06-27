package tokenize

import (
	"regexp"
	"strings"

	"github.com/frhorschig/kant-search-backend/core/errors"
	c "github.com/frhorschig/kant-search-backend/core/upload/internal/common"
)

func Tokenize(input string) ([]c.Token, *errors.Error) {
	input = strings.TrimSpace(input)
	if input[0] != '{' {
		return nil, &errors.Error{
			Msg:    errors.UPLOAD_WRONG_STARTING_CHAR,
			Params: []string{string(input[0])},
		}
	}

	rType := regexp.MustCompile(`^(paragraph|heading|footnote|fn|p|l)([|0-9])`)
	rLoc := regexp.MustCompile(`^\d+(\.\d+)?[}|]`)
	rChar := regexp.MustCompile(`^[^{}]*`)

	var tokens []c.Token
	for len(input) > 0 {
		switch {
		case strings.HasPrefix(input, "{"):
			tokens = append(tokens, c.NewOpen())
			input = input[1:]
		case strings.HasPrefix(input, "}"):
			tokens = append(tokens, c.NewClose())
			input = input[1:]
		case strings.HasPrefix(input, "|"):
			tokens = append(tokens, c.NewSeparator())
			input = input[1:]
		default:
			if matches := rType.FindStringSubmatch(input); matches != nil {
				tokens = append(tokens, c.NewClass(matches[1]))
				input = input[len(matches[1]):]
			} else if match := rLoc.FindString(input); match != "" {
				tokens = append(tokens, c.NewParam(match[:len(match)-1]))
				input = input[len(match)-1:]
			} else {
				match := rChar.FindString(input)
				tokens = append(tokens, c.NewText(strings.TrimSpace(match)))
				input = input[len(match):]
			}
		}
		input = strings.TrimSpace(input)
	}
	return tokens, nil
}
