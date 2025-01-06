package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"context"
	"fmt"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
)

type XmlMapper interface {
	Map(ctx context.Context, doc *etree.Document) ([]model.Work, error)
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

func (rec *xmlMapperImpl) Map(ctx context.Context, doc *etree.Document) ([]model.Work, error) {
	xmlStr, err := doc.WriteToString()
	if err != nil {
		return nil, fmt.Errorf("error when writing xml to string: %v", err.Error())
	}
	xmlStr = mapping.Simplify(xmlStr)
	println(xmlStr)
	doc.ReadFromString(xmlStr)
	// vol := doc.FindElement("//band")
	// println(etree.NewDocumentWithRoot(vol).WriteToString())

	// TODO frhorschig implement me
	return []model.Work{}, nil
}
