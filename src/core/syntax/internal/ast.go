package internal

import (
	"github.com/frhorschig/kant-search-backend/core/errors"
)

func Parse(tokens []Token) *errors.Error {
	_, err := parseExpression(&tokens)
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		return &errors.Error{
			Msg:    errors.UNEXPECTED_TOKEN,
			Params: []string{tokens[0].Text},
		}
	}
	return nil
}

func parseExpression(tokens *[]Token) (*astNote, *errors.Error) {
	node, err := parseTerm(tokens)
	if err != nil {
		return nil, err
	}

	for len(*tokens) > 0 && ((*tokens)[0].IsAnd || (*tokens)[0].IsOr) {
		opToken := &(*tokens)[0]
		*tokens = (*tokens)[1:]
		nextNode, err := parseTerm(tokens)
		if err != nil {
			return nil, err
		}
		node = &astNote{
			Left:  node,
			Right: nextNode,
			Token: opToken,
		}
	}

	return node, nil
}

func parseTerm(tokens *[]Token) (*astNote, *errors.Error) {
	if len(*tokens) == 0 {
		return nil, &errors.Error{Msg: errors.UNEXPECTED_END_OF_INPUT}
	}

	if (*tokens)[0].IsNot {
		token := &(*tokens)[0]
		*tokens = (*tokens)[1:]
		node, err := parseTerm(tokens)
		if err != nil {
			return nil, err
		}
		return &astNote{Left: node, Token: token}, nil
	}

	return parseFactor(tokens)
}

func parseFactor(tokens *[]Token) (*astNote, *errors.Error) {
	token := &(*tokens)[0]
	switch {
	case token.IsWord || token.IsPhrase:
		*tokens = (*tokens)[1:]
		return &astNote{Token: token}, nil
	case token.IsOpen:
		*tokens = (*tokens)[1:]
		node, err := parseExpression(tokens)
		if err != nil {
			return nil, err
		}
		if len(*tokens) == 0 || !(*tokens)[0].IsClose {
			return nil, &errors.Error{Msg: errors.MISSING_CLOSING_PARENTHESIS}
		}
		*tokens = (*tokens)[1:]
		return node, nil
	default:
		return nil, &errors.Error{
			Msg:    errors.UNEXPECTED_TOKEN,
			Params: []string{token.Text},
		}
	}
}
