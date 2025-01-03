package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"
	"fmt"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type VolumeUploadProcessor interface {
	Process(ctx context.Context, body []byte) error
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

func (rec *volumeUploadProcessorImpl) Process(ctx context.Context, body []byte) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err.Error())
	}

	_, err := rec.xmlMapper.MapVolume(ctx, doc)
	if err != nil {
		return err
	}

	// TODO frhorschig: implement writing to database
	return nil
}
