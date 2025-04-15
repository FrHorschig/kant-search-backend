package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

type XmlMapper interface {
	Map(xml string) (works []model.Work, err errors.ErrorNew)
}

type xmlMapperImpl struct {
}

func NewXmlMapper() XmlMapper {
	impl := xmlMapperImpl{}
	return &impl
}

func (rec *xmlMapperImpl) Map(xml string) ([]model.Work, errors.ErrorNew) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	sections, summaries, footnotes, err := mapping.MapToTree(doc)
	if err.HasError {
		return nil, err
	}

	vol := doc.FindElement("//band")
	volNo, err := util.ExtractNumericAttribute(vol, "nr")
	if err.HasError {
		return nil, err
	}
	works, err := mapping.MapToModel(volNo, sections, summaries, footnotes)
	if err.HasError {
		return nil, err
	}

	return works, errors.NilError()
}
