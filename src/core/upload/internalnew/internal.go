package internalnew

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/dbmapping/flattening"
	dbmetadataextraction "github.com/frhorschig/kant-search-backend/core/upload/internalnew/dbmapping/metadataextraction"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/metadatamapping/metadatamapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/metadatamapping/metadatamapping/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/metadatamapping/ordering"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/xmlmapping/metadataextraction"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/xmlmapping/referencemapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/xmlmapping/textmapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/xmlmapping/treemapping"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type XmlMapper interface {
	MapXml(volNr int32, xml string) (dbmodel.Volume, []dbmodel.Content, errs.UploadError)
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

func (rec *xmlMapperImpl) MapXml(volNr int32, xml string) (dbmodel.Volume, []dbmodel.Content, errs.UploadError) {
	// map xml to model
	works, fns, summs, err := treemapping.MapTree(xml)
	if err.HasError {
		return dbmodel.Volume{}, nil, err
	}
	err = textmapping.MapText(works, fns, summs)
	if err.HasError {
		return dbmodel.Volume{}, nil, err
	}
	err = metadataextraction.ExtractMetadata(works, fns, summs)
	if err.HasError {
		return dbmodel.Volume{}, nil, err
	}
	err = referencemapping.MapReferences(works, fns, summs)
	if err.HasError {
		return dbmodel.Volume{}, nil, err
	}

	// add non-xml metadata to model
	err = ordering.Order(works)
	if err.HasError {
		return dbmodel.Volume{}, nil, err
	}
	// use common/model/Volume instead of string
	volTitle, err := metadatamapping.MapMetadata(volNr, works, rec.metadata)
	if err.HasError {
		return dbmodel.Volume{}, nil, err
	}

	// map to db model
	dbWorks, contents := flattening.Flatten(works)
	dbmetadataextraction.ExtractMetadata(contents)
	return dbmodel.Volume{
		VolumeNumber: volNr,
		Title:        volTitle,
		Works:        dbWorks,
	}, contents, errs.Nil()
}
