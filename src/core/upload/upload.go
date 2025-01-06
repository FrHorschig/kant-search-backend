package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type UploadProcessor interface {
	Process(ctx context.Context, volNum int32, xml string) error
}

type uploadProcessorImpl struct {
	paragraphRepo dataaccess.ParagraphRepo
	sentenceRepo  dataaccess.SentenceRepo
	xmlMapper     internal.XmlMapper
}

func NewUploadProcessor(paragraphRepo dataaccess.ParagraphRepo, sentenceRepo dataaccess.SentenceRepo) UploadProcessor {
	processor := uploadProcessorImpl{
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
		xmlMapper:     internal.NewXmlMapper(),
	}
	return &processor
}

func (rec *uploadProcessorImpl) Process(ctx context.Context, volNum int32, xml string) error {
	_, err := rec.xmlMapper.Map(ctx, xml)
	if err != nil {
		return err
	}

	// TODO frhorschig: implement writing to database
	return nil
}
