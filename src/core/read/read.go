package read

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type ReadProcessor interface {
	Process(ctx context.Context) error
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

func (rec *readProcessorImpl) Process(ctx context.Context) error {
	// TODO implement me
	return nil
}
