package search

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type Searcher interface {
	SearchParagraphs(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchMatch, error)
}

type SearcherImpl struct {
	searchRepo repository.SearchRepo
}

func NewSearcher(searchRepo repository.SearchRepo) Searcher {
	impl := SearcherImpl{searchRepo: searchRepo}
	return &impl
}

func (rec *SearcherImpl) SearchParagraphs(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchMatch, error) {
	matches, err := rec.searchRepo.SearchParagraphs(ctx, criteria)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
