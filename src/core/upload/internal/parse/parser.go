package parse

import (
	"regexp"

	"github.com/FrHorschig/kant-search-backend/core/errors"
	c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"
)

func Parse(tokens []c.Token) ([]c.Expression, *errors.Error) {
	tk := &tokenIterator{tokens: tokens}
	results := make([]c.Expression, 0)
	for tk.hasNext() {
		expr, err := parseExpression(tk)
		if err != nil {
			return nil, err
		}
		results = append(results, expr)
	}

	return results, nil
}

func parseExpression(tk *tokenIterator) (c.Expression, *errors.Error) {
	if !tk.consume(c.OPEN) {
		return c.Expression{}, &errors.Error{
			Msg:    errors.UNEXPECTED_TOKEN,
			Params: []string{tk.peek().Text},
		}
	}

	meta, err := parseMetadata(tk)
	if err != nil {
		return c.Expression{}, err
	}
	expr := c.Expression{Metadata: *meta}

	if tk.consume(c.SEPARATOR) {
		content, err := parseContent(tk)
		if err != nil {
			return c.Expression{}, err
		}
		expr.Content = content
	}

	if !tk.consume(c.CLOSE) {
		var errText string
		if expr.Content != nil {
			if len(*expr.Content) < 16 {
				errText = *expr.Content
			} else {
				errText = "..." + (*expr.Content)[len(*expr.Content)-16:len(*expr.Content)]
			}
		} else {
			errText = expr.Metadata.Class
			if expr.Metadata.Param != nil {
				errText += *expr.Metadata.Param
			}
		}
		return c.Expression{}, &errors.Error{
			Msg:    errors.MISSING_CLOSING_BRACE,
			Params: []string{errText},
		}
	}

	return expr, nil
}

func parseMetadata(tk *tokenIterator) (*c.Metadata, *errors.Error) {
	text, ok := tk.consumeWithText(c.CLASS)
	if !ok {
		return nil, &errors.Error{
			Msg: errors.MISSING_EXPR_TYPE,
		}
	}

	meta := &c.Metadata{Class: text}
	if loc, ok := tk.consumeWithText(c.PARAM); ok {
		meta.Param = &loc
	}

	return meta, nil
}

func parseContent(tk *tokenIterator) (*string, *errors.Error) {
	content := ""
	for tk.hasNext() && tk.peek().Type != c.CLOSE {
		if tk.peek().Type == c.OPEN {
			expr, err := parseExpression(tk)
			if err != nil {
				return nil, err
			}
			content += " {" + expr.String() + "}"
		} else if text, ok := tk.consumeWithText(c.TEXT); ok {
			content += " " + text
		}
	}
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	if len(content) > 0 {
		content = content[1:]
	}

	return &content, nil
}
