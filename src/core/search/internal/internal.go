package internal

import (
	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/internal/parse"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CheckSyntax(searchTerms string) (*model.AstNode, *errors.SyntaxError) {
	tokens, err := parse.Tokenize(searchTerms)
	if err != nil {
		return nil, err
	}
	node, err := parse.Parse(tokens)
	if err != nil {
		return nil, err
	}
	return node, nil
}
