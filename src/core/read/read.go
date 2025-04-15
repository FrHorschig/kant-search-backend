package read

//go:generate mockgen -source=$GOFILE -destination=mocks/read_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
)

type ReadProcessor interface {
	ProcessVolumes(ctx context.Context) ([]esmodel.Volume, error)
	ProcessWork(ctx context.Context, workId string) (*esmodel.Work, error)
	ProcessFootnotes(ctx context.Context, workId string) ([]esmodel.Content, error)
	ProcessHeadings(ctx context.Context, workId string) ([]esmodel.Content, error)
	ProcessParagraphs(ctx context.Context, workId string) ([]esmodel.Content, error)
	ProcessSummaries(ctx context.Context, workId string) ([]esmodel.Content, error)
}

type readProcessorImpl struct {
	volumeRepo  dataaccess.VolumeRepo
	workRepo    dataaccess.WorkRepo
	contentRepo dataaccess.ContentRepo
}

func NewReadProcessor(volumeRepo dataaccess.VolumeRepo, workRepo dataaccess.WorkRepo, contentRepo dataaccess.ContentRepo) ReadProcessor {
	processor := readProcessorImpl{
		volumeRepo:  volumeRepo,
		workRepo:    workRepo,
		contentRepo: contentRepo,
	}
	return &processor
}

func (rec *readProcessorImpl) ProcessVolumes(ctx context.Context) ([]esmodel.Volume, error) {
	return rec.volumeRepo.GetAll(ctx)
}

func (rec *readProcessorImpl) ProcessWork(ctx context.Context, workId string) (*esmodel.Work, error) {
	return rec.workRepo.Get(ctx, workId)
}

func (rec *readProcessorImpl) ProcessFootnotes(ctx context.Context, workId string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetFootnotesByWorkId(ctx, workId)
}

func (rec *readProcessorImpl) ProcessHeadings(ctx context.Context, workId string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetHeadingsByWorkId(ctx, workId)
}

func (rec *readProcessorImpl) ProcessParagraphs(ctx context.Context, workId string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetParagraphsByWorkId(ctx, workId)
}

func (rec *readProcessorImpl) ProcessSummaries(ctx context.Context, workId string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetSummariesByWorkId(ctx, workId)
}
