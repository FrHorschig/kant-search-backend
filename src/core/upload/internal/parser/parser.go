package parser

import (
	"strings"

	"github.com/FrHorschig/kant-search-backend/core/errors"
)

func Parse(input string) (*Expression, *errors.Error) {
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

func parse(tokens []Token) (*Expression, *errors.Error) {
	tk := &tokenIterator{tokens: tokens}
	expr, err := parseExpression(tk)
	if err != nil {
		return nil, err
	}

	if tk.hasNext() {
		return nil, &errors.Error{
			Msg:    errors.UNEXPECTED_TOKEN,
			Params: []string{tk.peek().Text},
		}
	}

	return expr, nil
}

func parseExpression(tk *tokenIterator) (*Expression, *errors.Error) {
	if !tk.consume(OPEN) {
		return nil, &errors.Error{
			Msg:    errors.UNEXPECTED_TOKEN,
			Params: []string{tk.peek().Text},
		}
	}

	meta, err := parseMetadata(tk)
	if err != nil {
		return nil, err
	}
	expr := &Expression{Metadata: *meta}

	if tk.consume(SEPARATOR) {
		content, err := parseContent(tk)
		if err != nil {
			return nil, err
		}
		expr.Content = content
	}

	if !tk.consume(CLOSE) {
		var errText string
		if expr.Content != nil {
			errText = expr.Content.Texts[len(expr.Content.Texts)-1]
		} else {
			errText = expr.Metadata.Class
			if expr.Metadata.Location != nil {
				errText += *expr.Metadata.Location
			}
		}
		return nil, &errors.Error{
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
		meta.Location = &loc
	}

	return meta, nil
}

func parseContent(tk *tokenIterator) (*Content, *errors.Error) {
	content := &Content{}

	for tk.hasNext() && tk.peek().Type != CLOSE {
		if tk.peek().Type == OPEN {
			expr, err := parseExpression(tk)
			if err != nil {
				return nil, err
			}
			content.Expressions = append(content.Expressions, expr)
		} else if text, ok := tk.consumeWithText(TEXT); ok {
			content.Texts = append(content.Texts, text)
		}
	}

	return content, nil
}
