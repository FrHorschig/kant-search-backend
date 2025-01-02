package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt2"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt31"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt32"
	"github.com/frhorschig/kant-search-backend/core/upload/model/vol14"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type VolumeUploadProcessor interface {
	ProcessAbt1(ctx context.Context, volNum int32, vol abt1.Band) error
	ProcessAbt2(ctx context.Context, volNum int32, vol abt2.Band) error
	ProcessVol14(ctx context.Context, volNum int32, vol vol14.Band) error
	ProcessAbt31(ctx context.Context, volNum int32, vol abt31.Band) error
	ProcessAbt32(ctx context.Context, volNum int32, vol abt32.Band) error
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

func (rec *volumeUploadProcessorImpl) ProcessAbt1(ctx context.Context, volNum int32, vol abt1.Band) error {
	return nil
}

func (rec *volumeUploadProcessorImpl) ProcessAbt2(ctx context.Context, volNum int32, vol abt2.Band) error {
	return nil
}

func (rec *volumeUploadProcessorImpl) ProcessVol14(ctx context.Context, volNum int32, vol vol14.Band) error {
	return nil
}

func (rec *volumeUploadProcessorImpl) ProcessAbt31(ctx context.Context, volNum int32, vol abt31.Band) error {
	return nil
}

func (rec *volumeUploadProcessorImpl) ProcessAbt32(ctx context.Context, volNum int32, vol abt32.Band) error {
	return nil
}
