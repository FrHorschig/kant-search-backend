package search

//go:generate mockgen -source=$GOFILE -destination=mocks/search_processor_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type SearchProcessor interface {
	Search(ctx context.Context, searchString string, options model.SearchOptions) ([]model.SearchResult, errors.SearchError)
}

type searchProcessorImpl struct {
	contentRepo dataaccess.ContentRepo
}

func NewSearchProcessor(contentRepo dataaccess.ContentRepo) SearchProcessor {
	impl := searchProcessorImpl{
		contentRepo: contentRepo,
	}
	return &impl
}

func (rec *searchProcessorImpl) Search(ctx context.Context, searchString string, options model.SearchOptions) ([]model.SearchResult, errors.SearchError) {
	// TODO implement me
	return nil, errors.Nil()
}
