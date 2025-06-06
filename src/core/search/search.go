package search

//go:generate mockgen -source=$GOFILE -destination=mocks/search_processor_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/internal"
	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type SearchProcessor interface {
	Search(ctx context.Context, searchString string, options model.SearchOptions) ([]model.SearchResult, errors.SearchError)
}

type searchProcessorImpl struct {
	astParser   internal.AstParser
	contentRepo dataaccess.ContentRepo
}

func NewSearchProcessor(contentRepo dataaccess.ContentRepo) SearchProcessor {
	impl := searchProcessorImpl{
		astParser:   internal.NewAstParser(),
		contentRepo: contentRepo,
	}
	return &impl
}

func (rec *searchProcessorImpl) Search(ctx context.Context, searchTerms string, options model.SearchOptions) ([]model.SearchResult, errors.SearchError) {
	ast, syntaxErr := rec.astParser.Parse(searchTerms)
	if syntaxErr != nil {
		return nil, errors.New(syntaxErr, nil)
	}
	results, err := rec.contentRepo.Search(ctx, ast, options)
	if err != nil {
		return nil, errors.New(nil, err)
	}
	return results, errors.Nil()
}
