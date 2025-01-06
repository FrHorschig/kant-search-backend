package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/transform"
)

type XmlMapper interface {
	Map(xml string) ([]model.Work, error)
}

type xmlMapperImpl struct {
	pyUtil pyutil.PythonUtil
}

func NewXmlMapper() XmlMapper {
	impl := xmlMapperImpl{
		pyUtil: pyutil.NewPythonUtil(),
	}
	return &impl
}

func (rec *xmlMapperImpl) Map(xml string) ([]model.Work, error) {
	xml = transform.Simplify(xml)
	_, _, err := mapping.MapToSections(xml)
	if err != nil {
		return nil, err
	}
	// TODO frhorsch implement me
	return nil, err
}
