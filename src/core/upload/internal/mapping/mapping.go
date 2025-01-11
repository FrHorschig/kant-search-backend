package mapping

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	commonmodel "github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/transform"
)

func MapToTree(doc *etree.Document) ([]model.Section, errors.ErrorNew) {
	vol := doc.FindElement("//band")
	works, _ := findSections(vol.FindElement("//hauptteil"))
	// randtexte := findRandtexte(vol.FindElement("//randtexte"))
	// if len(randtexte) > 0 {
	// 	mergeRandtexte(works, randtexte)
	// }
	// footnotes := findFootnotes(vol.FindElement("//fussnoten"))
	// mergeFootnotes(works, footnotes)

	return works, errors.NilError()
}

func findSections(hauptteil *etree.Element) ([]model.Section, errors.ErrorNew) {
	secs := make([]model.Section, 0)
	var currentSec *model.Section
	pagePrefix := ""
	for _, el := range hauptteil.ChildElements() {
		switch el.Tag {
		case "hj":
			// We use the year data from https://kant.bbaw.de/de/akademieausgabe, so we ignore these year elements
			continue

		case "h1", "h2", "h3", "h4", "h5", "h6", "h7", "h8", "h9":
			hx, err := transform.Hx(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hx.TextTitle = pagePrefix + hx.TextTitle
				pagePrefix = ""
			}
			if hx.TocTitle == "" { // this happens if the hx only consists of an hu element
				currentSec.Paragraphs = append(currentSec.Paragraphs, hx.TextTitle)
				continue
			}

			sec := model.Section{Heading: hx, Paragraphs: []string{}, Sections: []model.Section{}}
			if hx.Level == commonmodel.H1 {
				secs = append(secs, sec)
				currentSec = &secs[len(secs)-1]
				continue
			}

			parent := findParent(hx, currentSec)
			sec.Parent = parent
			parent.Sections = append(parent.Sections, sec)
			currentSec = &parent.Sections[len(parent.Sections)-1]

		case "hu":
			hu, err := transform.Hu(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hu = pagePrefix + hu
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, hu)

		case "p":
			p, err := transform.P(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				p = pagePrefix + p
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, p)

		case "seite":
			pagePrefix += transform.Seite(el)
		case "op":
			continue
		case "table":
			// TODO implement me

		default:
			return nil, errors.NewError(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
		}
	}
	return secs, errors.NilError()
}

func findRandtexte(randtexte *etree.Element) []model.Randtext {
	panic("unimplemented")
}

func findFootnotes(fussnoten *etree.Element) []model.Footnote {
	panic("unimplemented")
}

func mergeRandtexte(sections []model.Section, randtexte []model.Randtext) {
	panic("unimplemented")
}

func mergeFootnotes(sections []model.Section, footnotes []model.Footnote) {
	panic("unimplemented")
}

func findParent(hx model.Heading, current *model.Section) *model.Section {
	if hx.Level > current.Heading.Level { // new heading is lower in hierarchy
		return current
	} else if hx.Level == current.Heading.Level {
		return current.Parent
	} else { // new heading is higher in hierarchy
		return findParent(hx, current.Parent)
	}
}

func prnt(el *etree.Element) {
	println(etree.NewDocumentWithRoot(el).WriteToString())
}
