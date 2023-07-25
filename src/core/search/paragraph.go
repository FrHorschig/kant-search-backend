package search

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type ParagraphSearcher interface {
	Search(ctx context.Context, criteria model.SearchCriteria) (model.ParagraphResults, error)
}

type ParagraphSearcherImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewParagraphSearcher(paragraphRepo repository.ParagraphRepo) ParagraphSearcher {
	impl := ParagraphSearcherImpl{paragraphRepo: paragraphRepo}
	return &impl
}

func (rec *ParagraphSearcherImpl) Search(ctx context.Context, criteria model.SearchCriteria) (model.ParagraphResults, error) {
	// TODO implement me
	return model.ParagraphResults{}, nil
}
