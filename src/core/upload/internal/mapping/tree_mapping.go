package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/tree_mapper.go -package=mocks

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/transform"
)

type TreeMapper interface {
	Map(doc *etree.Document) (
		sections []model.Section,
		summaries []model.Summary,
		footnotes []model.Footnote,
		err errors.ErrorNew,
	)
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

func (rec *TreeMapperImpl) Map(doc *etree.Document) ([]model.Section, []model.Summary, []model.Footnote, errors.ErrorNew) {
	works, err := rec.findSections(doc.FindElement("//hauptteil"))
	if err.HasError {
		return nil, nil, nil, err
	}
	summaries, err := rec.findSummaries(doc.FindElement("//randtexte"))
	if err.HasError {
		return nil, nil, nil, err
	}
	footnotes, err := rec.findFootnotes(doc.FindElement("//fussnoten"))
	if err.HasError {
		return nil, nil, nil, err
	}
	return works, summaries, footnotes, errors.NilError()
}

func (rec *TreeMapperImpl) findSections(hauptteil *etree.Element) ([]model.Section, errors.ErrorNew) {
	secs := make([]model.Section, 0)
	var currentSec *model.Section
	currentYear := ""
	pagePrefix := ""
	for _, el := range hauptteil.ChildElements() {
		switch el.Tag {
		case "h1":
			hx, err := rec.trafo.Hx(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hx.TextTitle = pagePrefix + " " + hx.TextTitle
				pagePrefix = ""
			}

			sec := model.Section{Heading: hx, Paragraphs: []string{}, Sections: []model.Section{}}
			if hx.Level == model.HWork {
				sec.Heading.Year = currentYear
				secs = append(secs, sec)
				currentSec = &secs[len(secs)-1]
			}

		case "h2", "h3", "h4", "h5", "h6", "h7", "h8", "h9":
			hx, err := rec.trafo.Hx(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hx.TextTitle = pagePrefix + " " + hx.TextTitle
				pagePrefix = ""
			}
			if hx.TocTitle == "" {
				// this happens if hx only has an hu element, which is the part of the heading that is not displayed in the TOC
				currentSec.Paragraphs = append(currentSec.Paragraphs, fmt.Sprintf(model.ParHeadFmt, hx.TextTitle))
				continue
			}

			sec := model.Section{Heading: hx, Paragraphs: []string{}, Sections: []model.Section{}}
			if len(secs) == 0 {
				return nil, errors.NewError(fmt.Errorf("the first heading is '%s', but must be h1", el.Tag), nil)
			}

			parent := findParent(hx, currentSec)
			sec.Parent = parent
			// this ensures a level difference of 1 in parent-child headings
			sec.Heading.Level = parent.Heading.Level + 1
			sec.Heading.TextTitle = fmt.Sprintf(model.HeadingFmt, sec.Heading.Level, sec.Heading.TextTitle, sec.Heading.Level)

			currentSec = &parent.Sections[len(parent.Sections)-1]

		case "hj":
			currentYear = strings.TrimSpace(el.Text())

		case "hu":
			hu, err := rec.trafo.Hu(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hu = pagePrefix + " " + hu
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, fmt.Sprintf(model.ParHeadFmt, hu))

		case "op":
			continue

		case "p":
			p, err := rec.trafo.P(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				p = pagePrefix + " " + p
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, p)

		case "seite":
			page, err := rec.trafo.Seite(el)
			if err.HasError {
				return nil, err
			}
			pagePrefix += page

		case "table":
			currentSec.Paragraphs = append(currentSec.Paragraphs, rec.trafo.Table())

		default:
			return nil, errors.NewError(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
		}
	}
	return secs, errors.NilError()
}

func (rec *TreeMapperImpl) findSummaries(randtexte *etree.Element) ([]model.Summary, errors.ErrorNew) {
	if randtexte == nil {
		return []model.Summary{}, errors.NilError()
	}
	result := make([]model.Summary, 0)
	for _, el := range randtexte.ChildElements() {
		rt, err := rec.trafo.Summary(el)
		if err.HasError {
			return nil, err
		}
		result = append(result, rt)
	}
	return result, errors.NilError()
}

func (rec *TreeMapperImpl) findFootnotes(fussnoten *etree.Element) ([]model.Footnote, errors.ErrorNew) {
	if fussnoten == nil {
		return []model.Footnote{}, errors.NilError()
	}
	result := make([]model.Footnote, 0)
	for _, el := range fussnoten.ChildElements() {
		rt, err := rec.trafo.Footnote(el)
		if err.HasError {
			return nil, err
		}
		result = append(result, rt)
	}
	return result, errors.NilError()
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
