package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadatamapping/metadatamapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/modelmap"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemap"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/metadataextraction"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/ordering"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/referencemapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/textmapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/treemapping"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
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

type XmlMapperNew interface { // TODO find appropriate name
	MapXml(volNr int32, xml string) (dbmodel.Volume, []dbmodel.Content, errs.UploadError)
}

type xmlMapperNewImpl struct {
	metadata metadata.Metadata
}

func NewXmlMapperNew(metadata metadata.Metadata) XmlMapperNew {
	impl := xmlMapperNewImpl{
		metadata: metadata,
	}
	return &impl
}

func (rec *xmlMapperNewImpl) MapXml(volNr int32, xml string) (dbmodel.Volume, []dbmodel.Content, errs.UploadError) {
	// map xml to model
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

	// add non-xml metadata to model
	err = ordering.Order(works)
	if err.HasError {
		return nil, err
	}
	volTitle, err := metadatamapping.MapMetadata(volNr, works, rec.metadata)
	if err.HasError {
		return nil, err
	}

	// map to db model
	panic("implement me")
}
