package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
)

type XmlMapper interface {
	MapAbt1(ctx context.Context, volNum int32, vol abt1.Kantabt1) (model.Volume, error)
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

func (rec *xmlMapperImpl) MapAbt1(ctx context.Context, volNum int32, vol abt1.Kantabt1) (model.Volume, error) {
	return model.Volume{}, nil
}
