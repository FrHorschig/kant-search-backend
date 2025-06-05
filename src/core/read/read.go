package read

//go:generate mockgen -source=$GOFILE -destination=mocks/read_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
)

type ReadProcessor interface {
	ProcessVolumes(ctx context.Context) ([]esmodel.Volume, error)
	ProcessWork(ctx context.Context, workCode string) (*esmodel.Work, error)
	ProcessFootnotes(ctx context.Context, workCode string) ([]esmodel.Content, error)
	ProcessHeadings(ctx context.Context, workCode string) ([]esmodel.Content, error)
	ProcessParagraphs(ctx context.Context, workCode string) ([]esmodel.Content, error)
	ProcessSummaries(ctx context.Context, workCode string) ([]esmodel.Content, error)
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

func (rec *readProcessorImpl) ProcessWork(ctx context.Context, workCode string) (*esmodel.Work, error) {
	return rec.workRepo.Get(ctx, workCode)
}

func (rec *readProcessorImpl) ProcessFootnotes(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetFootnotesByWorkCode(ctx, workCode)
}

func (rec *readProcessorImpl) ProcessHeadings(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetHeadingsByWorkCode(ctx, workCode)
}

func (rec *readProcessorImpl) ProcessParagraphs(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetParagraphsByWorkCode(ctx, workCode)
}

func (rec *readProcessorImpl) ProcessSummaries(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	return rec.contentRepo.GetSummariesByWorkCode(ctx, workCode)
}
