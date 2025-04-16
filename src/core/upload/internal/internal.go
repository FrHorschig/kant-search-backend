package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_mapper_mock.go -package=mocks

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/extract"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

type XmlMapper interface {
	MapVolume(volNr int32, xml string) (*model.Volume, errors.UploadError)
	MapWorks(volNr int32, xml string) ([]model.Work, errors.UploadError)
}

type xmlMapperImpl struct {
}

func NewXmlMapper() XmlMapper {
	impl := xmlMapperImpl{}
	return &impl
}

func (rec *xmlMapperImpl) MapWorks(volNr int32, xml string) ([]model.Work, errors.UploadError) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	sections, summaries, footnotes, err := mapping.MapToTree(doc)
	if err.HasError {
		return nil, err
	}
	works, err := mapping.MapToModel(volNr, sections, summaries, footnotes)
	if err.HasError {
		return nil, err
	}
	return works, errors.Nil()
}

func (rec *xmlMapperImpl) MapVolume(volNr int32, xml string) (*model.Volume, errors.UploadError) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	volXml := doc.FindElement("//band")
	xmlVolNr, err := extract.ExtractNumericAttribute(volXml, "nr")
	if err.HasError {
		return nil, err
	}
	if volNr != xmlVolNr {
		return nil, errors.New(fmt.Errorf("non matching volume numbers: is %d in URL, but %d in XML", volNr, xmlVolNr), nil)
	}

	section, e := getSection(volNr)
	if e != nil {
		return nil, errors.New(e, nil)
	}

	vol := model.Volume{
		VolumeNumber: volNr,
		Section:      section,
		Title:        volXml.FindElement("//titel").Text(),
	}
	return &vol, errors.Nil()
}

func getSection(volNr int32) (int32, error) {
	switch {
	case volNr > 0 && volNr <= 9:
		return 1, nil
	}
	return 0, fmt.Errorf("invalid volume number %d, must be >0 and <=9", volNr)
}
