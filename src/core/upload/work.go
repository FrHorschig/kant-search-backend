package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

import (
	"context"
	"strings"

	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
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
	input := strings.TrimSpace(upload.Text)
	if input[0] != '{' {
		return &errors.Error{
			Msg:    errors.WRONG_STARTING_CHAR,
			Params: []string{string(input[0])},
		}, nil
	}
	tokens := internal.Tokenize(input)
	_, err := internal.Parse(tokens)
	if err != nil {
		return err, nil
	}
	// TODO frhorsch: implement

	return nil, nil
}
