package upload

import (
	"context"
	"regexp"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/core/processing"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type WorkProcessor interface {
	Process(ctx context.Context, work model.Work) error
}

type workProcessorImpl struct {
	workRepo      repository.WorkRepo
	paragraphRepo repository.ParagraphRepo
	sentenceRepo  repository.SentenceRepo
}

func NewWorkProcessor(workRepo repository.WorkRepo, paragraphRepo repository.ParagraphRepo, sentenceRepo repository.SentenceRepo) WorkProcessor {
	processor := workProcessorImpl{
		workRepo:      workRepo,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}
	return &processor
}

func (rec *workProcessorImpl) Process(ctx context.Context, work model.Work) error {
	workId, err := rec.workRepo.Insert(ctx, work)
	if err != nil {
		return err
	}
	paras, err := processing.BuildParagraphModels(work.Text, workId)
	if err != nil {
		return err
	}

	// For now remove all line numbering
	r, _ := regexp.Compile(`\s*\{l\d+\}\s*`)
	for i := range paras {
		paras[i].Text = r.ReplaceAllString(paras[i].Text, " ")
	}

	for _, p := range paras {
		_, err := rec.paragraphRepo.Insert(ctx, p)
		if err != nil {
			return err
		}
	}

	return nil
}
