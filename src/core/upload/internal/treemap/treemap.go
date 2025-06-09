package treemap

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemap/transform"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

func MapToTree(doc *etree.Document) ([]model.TreeSection, []model.TreeSummary, []model.TreeFootnote, errs.UploadError) {
	works, err := findSections(doc.FindElement("//hauptteil"))
	if err.HasError {
		return nil, nil, nil, err
	}
	summaries, err := findSummaries(doc.FindElement("//randtexte"))
	if err.HasError {
		return nil, nil, nil, err
	}
	footnotes, err := findFootnotes(doc.FindElement("//fussnoten"))
	if err.HasError {
		return nil, nil, nil, err
	}
	return works, summaries, footnotes, errs.Nil()
}

func findSections(hauptteil *etree.Element) ([]model.TreeSection, errs.UploadError) {
	secs := make([]model.TreeSection, 0)
	var currentSec *model.TreeSection
	currentYear := ""
	pagePrefix := ""
	for _, el := range hauptteil.ChildElements() {
		// TODO (later): If <seite> occurs after a word with <trenn>, the page break is inside the word. Improve the handling of these situations.
		switch el.Tag {
		case "h1":
			hx, err := transform.Hx(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hx.TextTitle = pagePrefix + " " + hx.TextTitle
				pagePrefix = ""
			}

			sec := model.TreeSection{Heading: hx, Paragraphs: []string{}, Sections: []model.TreeSection{}}
			sec.Heading.Year = currentYear
			secs = append(secs, sec)
			currentSec = &secs[len(secs)-1]

		case "h2", "h3", "h4", "h5", "h6", "h7", "h8", "h9":
			hx, err := transform.Hx(el)
			if err.HasError {
				return nil, err
			}
			if hx.TocTitle == "" {
				// this happens if hx only has an hu element, which is the part of the heading that is not displayed in the TOC
				currentSec.Paragraphs = append(currentSec.Paragraphs, util.FmtParHeading(hx.TextTitle))
				continue
			}

			sec := model.TreeSection{Heading: hx, Paragraphs: []string{}, Sections: []model.TreeSection{}}
			if len(secs) == 0 {
				return nil, errs.New(fmt.Errorf("the first heading is '%s', but must be h1", el.Tag), nil)
			}

			parent := findParent(hx, currentSec)
			sec.Parent = parent
			sec.Heading.Level = parent.Heading.Level + 1 // this ensures a level difference of 1 in parent-child headings, even if by mistake a level is skipped
			sec.Heading.TextTitle = util.FmtHeading(int32(sec.Heading.Level), sec.Heading.TextTitle)
			if pagePrefix != "" {
				sec.Heading.TextTitle = pagePrefix + " " + sec.Heading.TextTitle
				pagePrefix = ""
			}

			sec.Parent.Sections = append(sec.Parent.Sections, sec)
			currentSec = &parent.Sections[len(parent.Sections)-1]

		case "hj":
			currentYear = strings.TrimSpace(el.Text())

		case "hu":
			hu, err := transform.Hu(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				hu = pagePrefix + " " + hu
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, util.FmtParHeading(hu))

		case "op":
			continue

		case "p":
			p, err := transform.P(el)
			if err.HasError {
				return nil, err
			}
			if pagePrefix != "" {
				p = pagePrefix + " " + p
				pagePrefix = ""
			}
			currentSec.Paragraphs = append(currentSec.Paragraphs, p)

		case "seite":
			page, err := transform.Seite(el)
			if err.HasError {
				return nil, err
			}
			pagePrefix += page

		case "table":
			currentSec.Paragraphs = append(currentSec.Paragraphs, transform.Table())

		default:
			return nil, errs.New(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
		}
	}
	return secs, errs.Nil()
}

func findSummaries(randtexte *etree.Element) ([]model.TreeSummary, errs.UploadError) {
	if randtexte == nil {
		return []model.TreeSummary{}, errs.Nil()
	}
	result := make([]model.TreeSummary, 0)
	for _, el := range randtexte.ChildElements() {
		rt, err := transform.Summary(el)
		if err.HasError {
			return nil, err
		}
		result = append(result, rt)
	}
	return result, errs.Nil()
}

func findFootnotes(fussnoten *etree.Element) ([]model.TreeFootnote, errs.UploadError) {
	if fussnoten == nil {
		return []model.TreeFootnote{}, errs.Nil()
	}
	result := make([]model.TreeFootnote, 0)
	for _, el := range fussnoten.ChildElements() {
		rt, err := transform.Footnote(el)
		if err.HasError {
			return nil, err
		}
		result = append(result, rt)
	}
	return result, errs.Nil()
}

func findParent(hx model.TreeHeading, current *model.TreeSection) *model.TreeSection {
	if hx.Level > current.Heading.Level { // new heading is lower in hierarchy
		return current
	} else if hx.Level == current.Heading.Level {
		return current.Parent
	} else { // new heading is higher in hierarchy
		return findParent(hx, current.Parent)
	}
}
