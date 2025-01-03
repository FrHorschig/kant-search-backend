package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"context"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
)

type XmlMapper interface {
	MapVolume(ctx context.Context, vol *etree.Document) ([]model.Work, error)
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

func (rec *xmlMapperImpl) MapVolume(ctx context.Context, vol *etree.Document) ([]model.Work, error) {
	// TODO frhorschig implement me
	return []model.Work{}, nil
}
