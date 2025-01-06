package mapping

import (
	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

func MapToSections(xml string) ([]model.Section, []model.Footnote, error) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)

	vol := doc.FindElement("//band")
	sections := findSections(vol.FindElement("//hauptteil"))
	randtexte := findRandtexte(vol.FindElement("//randtexte"))
	footnotes := findFootnotes(vol.FindElement("//fussnoten"))
	if len(randtexte) > 0 {
		mergeRandtexteIntoSections(sections, randtexte)
	}

	return sections, footnotes, nil
}

func findSections(element *etree.Element) []model.Section {
	panic("unimplemented")
}

func findRandtexte(element *etree.Element) []model.Randtext {
	panic("unimplemented")
}

func findFootnotes(element *etree.Element) []model.Footnote {
	panic("unimplemented")
}

func mergeRandtexteIntoSections(sections, randtexte any) {
	panic("unimplemented")
}
