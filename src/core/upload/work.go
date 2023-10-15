package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/parser"
	repository "github.com/FrHorschig/kant-search-backend/database"
)

type WorkUploadProcessor interface {
	Process(ctx context.Context, work model.WorkUpload) (*errors.Error, error)
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

func (rec *workUploadProcessorImpl) Process(ctx context.Context, upload model.WorkUpload) (*errors.Error, error) {
	_, err := parser.Parse(upload.Text)
	if err != nil {
		return err, nil
	}
	// TODO frhorsch: implement

	return nil, nil
}
