package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
)

type XmlMapper interface {
	Map(xml string) ([]model.Work, errors.ErrorNew)
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

func (rec *xmlMapperImpl) Map(xml string) ([]model.Work, errors.ErrorNew) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)

	_, _, _, err := rec.treeMapper.Map(doc)
	if err.HasError {
		return nil, err
	}
	// TODO frhorsch implement me
	return nil, errors.NilError()
}
