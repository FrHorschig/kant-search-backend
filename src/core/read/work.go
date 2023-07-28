package read

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type WorkReader interface {
	FindAll(ctx context.Context) ([]model.Work, error)
}

type workReaderImpl struct {
	workRepo repository.WorkRepo
}

func NewWorkReader(workRepo repository.WorkRepo) WorkReader {
	impl := workReaderImpl{workRepo: workRepo}
	return &impl
}

func (rec *workReaderImpl) FindAll(ctx context.Context) ([]model.Work, error) {
	works, err := rec.workRepo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	return works, nil
}
