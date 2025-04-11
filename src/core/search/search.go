package search

//go:generate mockgen -source=$GOFILE -destination=mocks/search_processor_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type SearchProcessor interface {
	Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
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

func (rec *searchProcessorImpl) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	// TODO implement me
	return nil, nil
}
