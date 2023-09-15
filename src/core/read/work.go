package read

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/database/model"
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
	return rec.workRepo.SelectAll(ctx)
}
