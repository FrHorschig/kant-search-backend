package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/transform"
)

type XmlMapper interface {
	Map(xml string) ([]model.Work, errors.ErrorNew)
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

func (rec *xmlMapperImpl) Map(xml string) ([]model.Work, errors.ErrorNew) {
	xml = transform.Simplify(xml)
	_, err := mapping.MapToWorks(xml)
	if err.HasError {
		return nil, err
	}
	// TODO frhorsch implement me
	return nil, errors.NilError()
}
