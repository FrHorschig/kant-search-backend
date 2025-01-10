package mapping

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

func MapToWorks(xml string) ([]model.Work, errors.ErrorNew) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)

	vol := doc.FindElement("//band")
	works, _ := findWorks(vol.FindElement("//hauptteil"))
	randtexte := findRandtexte(vol.FindElement("//randtexte"))
	if len(randtexte) > 0 {
		mergeRandtexte(works, randtexte)
	}
	footnotes := findFootnotes(vol.FindElement("//fussnoten"))
	mergeFootnotes(works, footnotes)

	return works, errors.NilError()
}

func findWorks(hauptteil *etree.Element) ([]model.Work, error) {
	panic("unimplemented")
}

func findRandtexte(randtexte *etree.Element) []model.Randtext {
	panic("unimplemented")
}

func findFootnotes(fussnoten *etree.Element) []model.Footnote {
	panic("unimplemented")
}

func mergeRandtexte(works []model.Work, randtexte []model.Randtext) {
	panic("unimplemented")
}

func mergeFootnotes(works []model.Work, footnotes []model.Footnote) {
	panic("unimplemented")
}
