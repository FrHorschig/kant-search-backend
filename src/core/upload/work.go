package upload

import (
	"context"
	"regexp"

	"github.com/FrHorschig/kant-search-backend/core/processing"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type WorkUploadProcessor interface {
	Process(ctx context.Context, work model.WorkUpload) error
}

type workUploadProcessorImpl struct {
	workRepo      repository.WorkRepo
	paragraphRepo repository.ParagraphRepo
	sentenceRepo  repository.SentenceRepo
}

func NewWorkProcessor(workRepo repository.WorkRepo, paragraphRepo repository.ParagraphRepo, sentenceRepo repository.SentenceRepo) WorkUploadProcessor {
	processor := workUploadProcessorImpl{
		workRepo:      workRepo,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}
	return &processor
}

func (rec *workUploadProcessorImpl) Process(ctx context.Context, upload model.WorkUpload) error {
	err := rec.workRepo.UpdateText(ctx, upload)
	if err != nil {
		return err
	}
	paras, err := processing.BuildParagraphModels(upload.Text, upload.WorkId)
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
