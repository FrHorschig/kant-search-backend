package search

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type ParagraphSearcher interface {
	Search(ctx context.Context, workIds []int32) ([]model.Paragraph, error)
}

type ParagraphSearcherImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewParagraphSearcher(paragraphRepo repository.ParagraphRepo) ParagraphSearcher {
	impl := ParagraphSearcherImpl{paragraphRepo: paragraphRepo}
	return &impl
}

func (rec *ParagraphSearcherImpl) Search(ctx context.Context, workIds []int32) ([]model.Paragraph, error) {
	// TODO implement me
	return nil, nil
}
