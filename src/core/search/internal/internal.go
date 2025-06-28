package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/ast_parser_mock.go -package=mocks

import (
	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/internal/model"
	"github.com/frhorschig/kant-search-backend/core/search/internal/parse"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type AstParser interface {
	Parse(searchTerms string) (*dbmodel.SearchTermNode, *errors.SyntaxError)
}

type astParserImpl struct{}

func NewAstParser() AstParser {
	impl := astParserImpl{}
	return &impl
}

func (rec *astParserImpl) Parse(searchTerms string) (*dbmodel.SearchTermNode, *errors.SyntaxError) {
	tokens, err := parse.Tokenize(searchTerms)
	if err != nil {
		return nil, err
	}
	node, err := parse.Parse(tokens)
	if err != nil {
		return nil, err
	}
	return mapNode(node), nil
}

func mapNode(node *model.AstNode) *dbmodel.SearchTermNode {
	if node == nil {
		return nil
	}
	mapped := dbmodel.SearchTermNode{
		Left:  mapNode(node.Left),
		Right: mapNode(node.Right),
		Token: &dbmodel.Token{
			IsAnd:    node.Token.IsAnd,
			IsOr:     node.Token.IsOr,
			IsNot:    node.Token.IsNot,
			IsWord:   node.Token.IsWord,
			IsPhrase: node.Token.IsPhrase,
			Text:     node.Token.Text,
		},
	}
	return &mapped
}
