package read

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type WorkReader interface {
	FindAll(ctx context.Context) ([]model.Work, error)
}

type WorkReaderImpl struct {
	workRepo repository.WorkRepo
}

func NewWorkReader(workRepo repository.WorkRepo) WorkReader {
	impl := WorkReaderImpl{workRepo: workRepo}
	return &impl
}

func (rec *WorkReaderImpl) FindAll(ctx context.Context) ([]model.Work, error) {
	works, err := rec.workRepo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	return works, nil
}
