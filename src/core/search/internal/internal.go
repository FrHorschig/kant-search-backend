package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/ast_parser_mock.go -package=mocks

import (
	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/internal/parse"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type AstParser interface {
	Parse(searchTerms string) (*model.AstNode, *errors.SyntaxError)
}

type astParserImpl struct{}

func NewAstParser() AstParser {
	impl := astParserImpl{}
	return &impl
}

func (rec *astParserImpl) Parse(searchTerms string) (*model.AstNode, *errors.SyntaxError) {
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
