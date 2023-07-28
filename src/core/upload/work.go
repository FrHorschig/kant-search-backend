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

	rec.insertParagraphs(ctx, paras)
	if err != nil {
		return err
	}

	return nil
}

func (rec *workProcessorImpl) insertParagraphs(ctx context.Context, paragraphs []model.Paragraph) error {
	for _, p := range paragraphs {
		text := p.Text
		p.Text = processing.RemoveFormatting(p.Text)
		id, err := rec.paragraphRepo.Insert(ctx, p)
		if err != nil {
			return err
		}
		p.Id = id
		p.Text = text
		err = rec.paragraphRepo.UpdateText(ctx, p, false)
		if err != nil {
			return err
		}
	}
	return nil
}
