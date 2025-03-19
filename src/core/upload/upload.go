package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type UploadProcessor interface {
	Process(ctx context.Context, volNum int32, xml string) errors.ErrorNew
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

func (rec *uploadProcessorImpl) Process(ctx context.Context, volNum int32, xml string) errors.ErrorNew {
	_, err := rec.xmlMapper.Map(xml)
	if err.HasError {
		return err
	}

	// TODO frhorschig: implement me
	//    - split sentences
	//    - add work metadata
	//    - write data to DB
	return errors.NilError()
}
