package search

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

type SearchProcessor interface {
	Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
}

type searchProcessorImpl struct {
	searchRepo repository.SearchRepo
}

func NewSearchProcessor(searchRepo repository.SearchRepo) SearchProcessor {
	impl := searchProcessorImpl{searchRepo: searchRepo}
	return &impl
}

func (rec *searchProcessorImpl) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	// TODO implement me
	return nil, nil
}
