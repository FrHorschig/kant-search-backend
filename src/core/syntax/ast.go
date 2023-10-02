package syntax

import (
	"errors"
	"fmt"
)

func buildAst(tokens []Token) (*Node, error) {
	node, err := parseExpression(&tokens)
	if err != nil {
		return &Node{}, err
	}

	if len(tokens) > 0 {
		return &Node{}, fmt.Errorf("unexpected token after parsing")
	}
	return node, nil
}

func parseExpression(tokens *[]Token) (*Node, error) {
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

		if opToken.IsAnd {
			node = &Node{
				Left:  node,
				Right: nextNode,
				Token: opToken,
			}
		} else {
			node = &Node{
				Left:  node,
				Right: nextNode,
				Token: opToken,
			}
		}
	}

	return node, nil
}

func parseTerm(tokens *[]Token) (*Node, error) {
	if len(*tokens) == 0 {
		return nil, errors.New("unexpected end of input")
	}

	if (*tokens)[0].IsNot {
		token := &(*tokens)[0]
		*tokens = (*tokens)[1:]
		node, err := parseTerm(tokens)
		if err != nil {
			return nil, err
		}
		return &Node{Left: node, Token: token}, nil
	}

	return parseFactor(tokens)
}

func parseFactor(tokens *[]Token) (*Node, error) {
	if len(*tokens) == 0 {
		return nil, errors.New("unexpected end of input")
	}

	token := &(*tokens)[0]
	switch {
	case token.IsWord:
		*tokens = (*tokens)[1:]
		return &Node{Token: token}, nil
	case token.IsPhrase:
		*tokens = (*tokens)[1:]
		return &Node{Token: token}, nil
	case token.IsOpen:
		*tokens = (*tokens)[1:]
		node, err := parseExpression(tokens)
		if err != nil {
			return nil, err
		}
		if len(*tokens) == 0 || !(*tokens)[0].IsClose {
			return nil, errors.New("missing closing parenthesis")
		}
		*tokens = (*tokens)[1:]
		return node, nil
	default:
		return nil, fmt.Errorf("unexpected token")
	}
}
