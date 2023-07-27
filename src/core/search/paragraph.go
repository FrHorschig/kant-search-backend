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
	paras, err := rec.paragraphRepo.Search(ctx, criteria)
	if err != nil {
		return model.ParagraphResults{}, err
	}

	results := model.ParagraphResults{
		Paragraphs:   paras,
		MatchedWords: criteria.SearchWords,
	}
	return results, nil
}
