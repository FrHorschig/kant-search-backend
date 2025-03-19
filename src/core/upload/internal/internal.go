package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/transform"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type XmlMapper interface {
	Map(xml string) (works []dbmodel.Work, err errors.ErrorNew)
}

type xmlMapperImpl struct {
	treeMapper  mapping.TreeMapper
	modelMapper mapping.ModelMapper
	pyUtil      pyutil.PythonUtil
}

func NewXmlMapper() XmlMapper {
	impl := xmlMapperImpl{
		treeMapper:  mapping.NewTreeMapper(),
		modelMapper: mapping.NewModelMapper(),
		pyUtil:      pyutil.NewPythonUtil(),
	}
	return &impl
}

func (rec *xmlMapperImpl) Map(xml string) ([]dbmodel.Work, errors.ErrorNew) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	sections, summaries, footnotes, err := rec.treeMapper.Map(doc)
	if err.HasError {
		return nil, err
	}
	works, err := rec.modelMapper.Map(sections, summaries, footnotes)
	if err.HasError {
		return nil, err
	}

	vol := doc.FindElement("//band")
	volNo, err := transform.ExtractNumericAttribute(vol, "nr")
	if err.HasError {
		return nil, err
	}
	for i := range works {
		works[i].Volume = volNo
	}

	return works, errors.NilError()
}
