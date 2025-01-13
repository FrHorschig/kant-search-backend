package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/tree_mapper.go -package=mocks

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	commonmodel "github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/transform"
)

type TreeMapper interface {
	Map(doc *etree.Document) ([]model.Section, errors.ErrorNew)
}

type TreeMapperImpl struct {
	trafo transform.XmlTransformator
}

func NewTreeMapper() TreeMapper {
	impl := TreeMapperImpl{
		trafo: transform.NewXmlTransformator(),
	}
	return &impl
}

func (rec *TreeMapperImpl) Map(doc *etree.Document) ([]model.Section, errors.ErrorNew) {
	// TODO implement me
	vol := doc.FindElement("//band")
	works, _ := rec.findSections(vol.FindElement("//hauptteil"))
	// randtexte := rec.findRandtexte(vol.FindElement("//randtexte"))
	// if len(randtexte) > 0 {
	// 	mergeRandtexte(works, randtexte)
	// }
	// footnotes := rec.findFootnotes(vol.FindElement("//fussnoten"))
	// mergeFootnotes(works, footnotes)

	return works, errors.NilError()
}

func (rec *TreeMapperImpl) findSections(hauptteil *etree.Element) ([]model.Section, errors.ErrorNew) {
	secs := make([]model.Section, 0)
	var currentSec *model.Section
	pagePrefix := ""
	for _, el := range hauptteil.ChildElements() {
		switch el.Tag {
		case "h1", "h2", "h3", "h4", "h5", "h6", "h7", "h8", "h9":
			hx, err := rec.trafo.Hx(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hx.TextTitle = pagePrefix + hx.TextTitle
				pagePrefix = ""
			}
			if hx.TocTitle == "" { // this happens if hx only has an hu element
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

		case "hj":
			// We use the year data from https://kant.bbaw.de/de/akademieausgabe, so we ignore these year elements
			continue

		case "hu":
			hu, err := rec.trafo.Hu(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hu = pagePrefix + hu
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, hu)

		case "op":
			continue

		case "p":
			p, err := rec.trafo.P(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				p = pagePrefix + p
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, p)

		case "seite":
			pagePrefix += rec.trafo.Seite(el)

		case "table":
			// TODO implement me

		default:
			return nil, errors.NewError(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
		}
	}
	return secs, errors.NilError()
}

func (rec *TreeMapperImpl) findRandtexte(randtexte *etree.Element) []model.Randtext {
	panic("unimplemented")
}

func (rec *TreeMapperImpl) findFootnotes(fussnoten *etree.Element) []model.Footnote {
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
