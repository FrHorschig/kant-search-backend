package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/modelmap"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemap"
)

type XmlMapper interface {
	MapVolume(volNr int32, xml string) (*model.Volume, errs.UploadError)
	MapWorks(volNr int32, xml string) ([]model.Work, errs.UploadError)
}

type xmlMapperImpl struct {
	metadata metadata.Metadata
}

func NewXmlMapper(metadata metadata.Metadata) XmlMapper {
	impl := xmlMapperImpl{
		metadata: metadata,
	}
	return &impl
}

func (rec *xmlMapperImpl) MapVolume(volNr int32, xml string) (*model.Volume, errs.UploadError) {
	metadata, mdErr := rec.metadata.Read(volNr)
	if mdErr != nil {
		return nil, errs.New(nil, mdErr)
	}

	vol := model.Volume{
		VolumeNumber: metadata.VolumeNumber,
		Title:        metadata.Title,
	}
	return &vol, errs.Nil()
}

func (rec *xmlMapperImpl) MapWorks(volNr int32, xml string) ([]model.Work, errs.UploadError) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	sections, summaries, footnotes, err := treemap.MapToTree(doc)
	if err.HasError {
		return nil, err
	}

	metadata, mdErr := rec.metadata.Read(volNr)
	if mdErr != nil {
		return nil, errs.New(nil, mdErr)
	}

	works, err := modelmap.MapToModel(metadata, sections, summaries, footnotes)
	if err.HasError {
		return nil, err
	}
	return works, errs.Nil()
}
