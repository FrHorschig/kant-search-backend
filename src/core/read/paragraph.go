package read

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type ParagraphReader interface {
	FindOfPages(ctx context.Context, workId int32, startPage int32, endPage int32) ([]model.Paragraph, error)
}

type paragraphReaderImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewParagraphReader(paragraphRepo repository.ParagraphRepo) ParagraphReader {
	impl := paragraphReaderImpl{paragraphRepo: paragraphRepo}
	return &impl
}

func (rec *paragraphReaderImpl) FindOfPages(ctx context.Context, workId int32, startPage int32, endPage int32) ([]model.Paragraph, error) {
	paras, err := rec.paragraphRepo.SelectOfPages(ctx, workId, startPage, endPage)
	if err != nil {
		return nil, err
	}
	return paras, nil
}
