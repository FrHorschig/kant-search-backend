package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/parse"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/tokenize"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/transform"
	"github.com/FrHorschig/kant-search-backend/database"
)

type WorkUploadProcessor interface {
	Process(ctx context.Context, work model.WorkUpload) *errors.Error
}

type workUploadProcessorImpl struct {
	workRepo      database.WorkRepo
	paragraphRepo database.ParagraphRepo
	sentenceRepo  database.SentenceRepo
	pyUtil        pyutil.PythonUtil
}

func NewWorkProcessor(workRepo database.WorkRepo, paragraphRepo database.ParagraphRepo, sentenceRepo database.SentenceRepo) WorkUploadProcessor {
	processor := workUploadProcessorImpl{
		workRepo:      workRepo,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
		pyUtil:        pyutil.NewPythonUtil(),
	}
	return &processor
}

func (rec *workUploadProcessorImpl) Process(ctx context.Context, upload model.WorkUpload) *errors.Error {
	tokens, err := tokenize.Tokenize(upload.Text)
	if err != nil {
		return err
	}
	exprs, err := parse.Parse(tokens)
	if err != nil {
		return err
	}

	paragraphs, err := transform.Transform(upload.WorkId, exprs, rec.pyUtil)
	if err != nil {
		return err
	}
	err = persistParagraphs(ctx, rec.paragraphRepo, paragraphs)
	if err != nil {
		return err
	}

	sentences, err := transform.FindSentences(paragraphs, rec.pyUtil)
	if err != nil {
		return err
	}
	return persistSentences(ctx, rec.sentenceRepo, sentences)
}

func persistParagraphs(ctx context.Context, repo database.ParagraphRepo, paragraphs []model.Paragraph) *errors.Error {
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

func persistSentences(ctx context.Context, repo database.SentenceRepo, sentences []model.Sentence) *errors.Error {
	_, err := repo.Insert(ctx, sentences)
	if err != nil {
		return &errors.Error{
			Msg:    errors.GO_ERR,
			Params: []string{err.Error()},
		}
	}
	return nil
}
