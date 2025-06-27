package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/modelmap"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemap"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/metadataextraction"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/ordering"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/referencemapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/textmapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/treemapping"
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

// =============================================================================

type XmlMapperNew interface {
	MapWorks(volNr int32, xml string) ([]model.Work, errs.UploadError)
}

type xmlMapperNewImpl struct {
}

func NewXmlMapperNew() XmlMapperNew {
	impl := xmlMapperNewImpl{}
	return &impl
}

func (rec *xmlMapperNewImpl) MapWorks(volNr int32, xml string) ([]model.Work, errs.UploadError) {
	works, fns, summs, err := treemapping.MapTree(xml)
	if err.HasError {
		return nil, err
	}
	err = textmapping.MapText(works, fns, summs)
	if err.HasError {
		return nil, err
	}
	err = metadataextraction.ExtractMetadata(works, fns, summs)
	if err.HasError {
		return nil, err
	}
	err = referencemapping.MapReferences(works, fns, summs)
	if err.HasError {
		return nil, err
	}
	err = ordering.Order(works)
	if err.HasError {
		return nil, err
	}
	return works, errs.Nil()
}
