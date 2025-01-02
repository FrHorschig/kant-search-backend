package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type VolumeUploadProcessor interface {
	ProcessAbt1(ctx context.Context, volNum int32, vol abt1.Kantabt1) error
}

type volumeUploadProcessorImpl struct {
	paragraphRepo dataaccess.ParagraphRepo
	sentenceRepo  dataaccess.SentenceRepo
	xmlMapper     internal.XmlMapper
}

func NewVolumeProcessor(paragraphRepo dataaccess.ParagraphRepo, sentenceRepo dataaccess.SentenceRepo) VolumeUploadProcessor {
	processor := volumeUploadProcessorImpl{
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
		xmlMapper:     internal.NewXmlMapper(),
	}
	return &processor
}

func (rec *volumeUploadProcessorImpl) ProcessAbt1(ctx context.Context, volNum int32, vol abt1.Kantabt1) error {
	return nil
}
