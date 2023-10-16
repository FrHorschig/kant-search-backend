package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/parser"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/transform"
	repository "github.com/FrHorschig/kant-search-backend/database"
)

type WorkUploadProcessor interface {
	Process(ctx context.Context, work model.WorkUpload) *errors.Error
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

func (rec *workUploadProcessorImpl) Process(ctx context.Context, upload model.WorkUpload) *errors.Error {
	exprs, err := parser.Parse(upload.Text)
	if err != nil {
		return err
	}

	paragraphs, err := transform.Transform(upload.WorkId, exprs)
	if err != nil {
		return err
	}
	err = persistParagraphs(ctx, rec.paragraphRepo, paragraphs)
	if err != nil {
		return err
	}

	sentences, err := transform.FindSentences(paragraphs)
	if err != nil {
		return err
	}
	return persistSentences(ctx, rec.sentenceRepo, sentences)
}

func persistParagraphs(ctx context.Context, repo repository.ParagraphRepo, paragraphs []model.Paragraph) *errors.Error {
	for i, p := range paragraphs {
		pId, err := repo.Insert(ctx, p)
		if err != nil {
			return &errors.Error{
				Msg:    errors.GO_ERR,
				Params: []string{err.Error()},
			}
		}
		paragraphs[i].Id = pId
	}
	return nil
}

func persistSentences(ctx context.Context, repo repository.SentenceRepo, sentences []model.Sentence) *errors.Error {
	_, err := repo.Insert(ctx, sentences)
	if err != nil {
		return &errors.Error{
			Msg:    errors.GO_ERR,
			Params: []string{err.Error()},
		}
	}
	return nil
}
