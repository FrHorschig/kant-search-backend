package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"context"
	"fmt"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
)

type XmlMapper interface {
	Map(ctx context.Context, xml []byte) ([]model.Work, error)
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

func (rec *xmlMapperImpl) Map(ctx context.Context, xml []byte) ([]model.Work, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xml); err != nil {
		return nil, fmt.Errorf("error unmarshaling request body: %v", err.Error())
	}
	println(doc.WriteToString())

	// TODO frhorschig implement me
	return []model.Work{}, nil
}
