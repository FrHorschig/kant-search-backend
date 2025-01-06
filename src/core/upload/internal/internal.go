package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"context"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
)

type XmlMapper interface {
	Map(ctx context.Context, xml string) ([]model.Work, error)
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

func (rec *xmlMapperImpl) Map(ctx context.Context, xml string) ([]model.Work, error) {
	xml = mapping.Simplify(xml)
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	vol := doc.FindElement("//band")
	println(etree.NewDocumentWithRoot(vol).WriteToString())

	// TODO frhorschig implement me
	return []model.Work{}, nil
}
