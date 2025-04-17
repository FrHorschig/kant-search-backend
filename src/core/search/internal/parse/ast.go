package parse

import (
	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func Parse(tokens []model.Token) (*model.AstNode, *errors.SyntaxError) {
	node, err := parseExpression(&tokens)
	if err != nil {
		return nil, err
	}

	if len(tokens) > 0 {
		return nil, &errors.SyntaxError{
			Msg:    errors.UnexpectedToken,
			Params: []string{tokens[0].Text},
		}
	}
	return node, nil
}

func parseExpression(tokens *[]model.Token) (*model.AstNode, *errors.SyntaxError) {
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
		node = &model.AstNode{
			Left:  node,
			Right: nextNode,
			Token: opToken,
		}
	}

	return node, nil
}

func parseTerm(tokens *[]model.Token) (*model.AstNode, *errors.SyntaxError) {
	if len(*tokens) == 0 {
		return nil, &errors.SyntaxError{Msg: errors.UnexpectedEndOfInput}
	}

	if (*tokens)[0].IsNot {
		token := &(*tokens)[0]
		*tokens = (*tokens)[1:]
		node, err := parseTerm(tokens)
		if err != nil {
			return nil, err
		}
		return &model.AstNode{Left: node, Token: token}, nil
	}

	return parseFactor(tokens)
}

func parseFactor(tokens *[]model.Token) (*model.AstNode, *errors.SyntaxError) {
	token := &(*tokens)[0]
	switch {
	case token.IsWord || token.IsPhrase:
		*tokens = (*tokens)[1:]
		return &model.AstNode{Token: token}, nil
	case token.IsOpen:
		*tokens = (*tokens)[1:]
		node, err := parseExpression(tokens)
		if err != nil {
			return nil, err
		}
		if len(*tokens) == 0 || !(*tokens)[0].IsClose {
			return nil, &errors.SyntaxError{Msg: errors.MissingCloseParenthesis}
		}
		*tokens = (*tokens)[1:]
		return node, nil
	default:
		return nil, &errors.SyntaxError{
			Msg:    errors.UnexpectedToken,
			Params: []string{token.Text},
		}
	}
}
