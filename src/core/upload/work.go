package upload

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type WorkProcessor interface {
	Process(ctx context.Context, work model.Work) error
}

type WorkProcessorImpl struct {
	workRepo      repository.WorkRepo
	paragraphRepo repository.ParagraphRepo
	sentenceRepo  repository.SentenceRepo
}

func NewWorkProcessor(workRepo repository.WorkRepo, paragraphRepo repository.ParagraphRepo, sentenceRepo repository.SentenceRepo) WorkProcessor {
	processor := WorkProcessorImpl{
		workRepo:      workRepo,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}
	return &processor
}

func (proc *WorkProcessorImpl) Process(ctx context.Context, work model.Work) error {
	// TODO implement me
	return nil
}
