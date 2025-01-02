package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type WorkUploadProcessor interface {
	Process(ctx context.Context, workId int32, text string) *errors.Error
}

type workUploadProcessorImpl struct {
	paragraphRepo dataaccess.ParagraphRepo
	sentenceRepo  dataaccess.SentenceRepo
	textMapper    internal.TextMapper
}

func NewWorkProcessor(paragraphRepo dataaccess.ParagraphRepo, sentenceRepo dataaccess.SentenceRepo) WorkUploadProcessor {
	processor := workUploadProcessorImpl{
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
		textMapper:    internal.NewTextMapper(),
	}
	return &processor
}

func (rec *workUploadProcessorImpl) Process(ctx context.Context, workId int32, text string) *errors.Error {
	paragraphs, err := rec.textMapper.FindParagraphs(workId, text)
	if err != nil {
		return err
	}

	// TODO frhorschig: use transaction
	err = deleteExistingData(ctx, rec.sentenceRepo, rec.paragraphRepo, workId)
	if err != nil {
		return err
	}
	err = persistParagraphs(ctx, rec.paragraphRepo, paragraphs)
	if err != nil {
		return err
	}

	sentences, err := rec.textMapper.FindSentences(paragraphs)
	if err != nil {
		return err
	}
	return persistSentences(ctx, rec.sentenceRepo, sentences)
}

func deleteExistingData(ctx context.Context, sentenceRepo dataaccess.SentenceRepo, paragraphRepo dataaccess.ParagraphRepo, workId int32) *errors.Error {
	err := sentenceRepo.DeleteByWorkId(ctx, workId)
	if err != nil {
		return &errors.Error{
			Msg:    errors.UPLOAD_GO_ERR,
			Params: []string{err.Error()},
		}
	}
	err = paragraphRepo.DeleteByWorkId(ctx, workId)
	if err != nil {
		return &errors.Error{
			Msg:    errors.UPLOAD_GO_ERR,
			Params: []string{err.Error()},
		}
	}
	return nil
}

func persistParagraphs(ctx context.Context, repo dataaccess.ParagraphRepo, paragraphs []model.Paragraph) *errors.Error {
	for i, p := range paragraphs {
		// TODO frhorschig: write and use bulk insert
		pId, err := repo.Insert(ctx, p)
		if err != nil {
			return &errors.Error{
				Msg:    errors.UPLOAD_GO_ERR,
				Params: []string{err.Error()},
			}
		}
		paragraphs[i].Id = pId
	}
	return nil
}

func persistSentences(ctx context.Context, repo dataaccess.SentenceRepo, sentences []model.Sentence) *errors.Error {
	_, err := repo.Insert(ctx, sentences)
	if err != nil {
		return &errors.Error{
			Msg:    errors.UPLOAD_GO_ERR,
			Params: []string{err.Error()},
		}
	}
	return nil
}
