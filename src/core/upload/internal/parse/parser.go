package parse

import (
	"regexp"
	"strings"

	"github.com/FrHorschig/kant-search-backend/core/errors"
)

func Parse(input string) ([]Expression, *errors.Error) {
	input = strings.TrimSpace(input)
	if input[0] != '{' {
		return nil, &errors.Error{
			Msg:    errors.WRONG_STARTING_CHAR,
			Params: []string{string(input[0])},
		}
	}
	tokens := tokenize(input)
	return parse(tokens)
}

func parse(tokens []Token) ([]Expression, *errors.Error) {
	tk := &tokenIterator{tokens: tokens}
	results := make([]Expression, 0)
	for tk.hasNext() {
		expr, err := parseExpression(tk)
		if err != nil {
			return nil, err
		}
		results = append(results, expr)
	}

	return results, nil
}

func parseExpression(tk *tokenIterator) (Expression, *errors.Error) {
	if !tk.consume(OPEN) {
		return Expression{}, &errors.Error{
			Msg:    errors.UNEXPECTED_TOKEN,
			Params: []string{tk.peek().Text},
		}
	}

	meta, err := parseMetadata(tk)
	if err != nil {
		return Expression{}, err
	}
	expr := Expression{Metadata: *meta}

	if tk.consume(SEPARATOR) {
		content, err := parseContent(tk)
		if err != nil {
			return Expression{}, err
		}
		expr.Content = content
	}

	if !tk.consume(CLOSE) {
		var errText string
		if expr.Content != nil {
			if len(*expr.Content) < 16 {
				errText = *expr.Content
			} else {
				errText = (*expr.Content)[len(*expr.Content)-16 : len(*expr.Content)]
			}
		} else {
			errText = expr.Metadata.Class
			if expr.Metadata.Param != nil {
				errText += *expr.Metadata.Param
			}
		}
		return Expression{}, &errors.Error{
			Msg:    errors.MISSING_CLOSING_BRACE,
			Params: []string{errText},
		}
	}

	return expr, nil
}

func parseMetadata(tk *tokenIterator) (*Metadata, *errors.Error) {
	text, ok := tk.consumeWithText(CLASS)
	if !ok {
		return nil, &errors.Error{
			Msg: errors.MISSING_EXPR_TYPE,
		}
	}

	meta := &Metadata{Class: text}
	if loc, ok := tk.consumeWithText(PARAM); ok {
		meta.Param = &loc
	}

	return meta, nil
}

func parseContent(tk *tokenIterator) (*string, *errors.Error) {
	content := ""
	for tk.hasNext() && tk.peek().Type != CLOSE {
		if tk.peek().Type == OPEN {
			expr, err := parseExpression(tk)
			if err != nil {
				return nil, err
			}
			content += " {" + expr.String() + "}"
		} else if text, ok := tk.consumeWithText(TEXT); ok {
			content += " " + text
		}
	}
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	if len(content) > 0 {
		content = content[1:]
	}

	return &content, nil
}
