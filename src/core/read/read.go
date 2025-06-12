package read

//go:generate mockgen -source=$GOFILE -destination=mocks/read_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type ReadProcessor interface {
	ProcessVolumes(ctx context.Context) ([]model.Volume, error)
	ProcessFootnotes(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	ProcessHeadings(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	ProcessParagraphs(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	ProcessSummaries(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
}

type readProcessorImpl struct {
	volumeRepo  dataaccess.VolumeRepo
	contentRepo dataaccess.ContentRepo
}

func NewReadProcessor(volumeRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo) ReadProcessor {
	processor := readProcessorImpl{
		volumeRepo:  volumeRepo,
		contentRepo: contentRepo,
	}
	return &processor
}

func (rec *readProcessorImpl) ProcessVolumes(ctx context.Context) ([]model.Volume, error) {
	return rec.volumeRepo.GetAll(ctx)
}

func (rec *readProcessorImpl) ProcessFootnotes(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	return rec.contentRepo.GetFootnotesByWork(ctx, workCode, ordinals)
}

func (rec *readProcessorImpl) ProcessHeadings(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	return rec.contentRepo.GetHeadingsByWork(ctx, workCode, ordinals)
}

func (rec *readProcessorImpl) ProcessParagraphs(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	return rec.contentRepo.GetParagraphsByWork(ctx, workCode, ordinals)
}

func (rec *readProcessorImpl) ProcessSummaries(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	return rec.contentRepo.GetSummariesByWork(ctx, workCode, ordinals)
}
