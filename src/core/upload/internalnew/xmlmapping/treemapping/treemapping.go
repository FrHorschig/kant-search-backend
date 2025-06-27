package treemapping

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
)

func MapTree(xml string) ([]model.Work, []model.Footnote, []model.Summary, errs.UploadError) {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	works, err := findWorks(doc.FindElement("//hauptteil"))
	if err.HasError {
		return nil, nil, nil, err
	}
	footnotes, err := findFootnotes(doc.FindElement("//fussnoten"))
	if err.HasError {
		return nil, nil, nil, err
	}
	summaries, err := findSummaries(doc.FindElement("//randtexte"))
	if err.HasError {
		return nil, nil, nil, err
	}
	return works, footnotes, summaries, errs.Nil()
}

func findWorks(hauptteil *etree.Element) ([]model.Work, errs.UploadError) {
	works := []model.Work{}
	currentYear := ""
	sectionList := make([]*[]model.Section, 10)
	var currPars *[]model.Paragraph = nil
	pagePrefix := ""

	for _, el := range hauptteil.ChildElements() {
		elStr, err := elemToString(el)
		if err != nil {
			return nil, errs.New(nil, fmt.Errorf("error writing element '%s' to string: %v", el.Tag, err))
		}

		switch el.Tag {
		case "hj":
			currentYear = strings.TrimSpace(el.Text())

		case "h1":
			work := model.Work{
				Title:      elStr,
				Year:       currentYear,
				Paragraphs: []model.Paragraph{},
				Sections:   []model.Section{},
			}
			works = append(works, work)
			w := &works[len(works)-1]
			currPars = &w.Paragraphs
			updateSectionList(&sectionList, &w.Sections, 1)

		case "h2", "h3", "h4", "h5", "h6", "h7", "h8", "h9":
			if pagePrefix != "" {
				elStr = elStr[0:len(el.Tag)+2] + pagePrefix + elStr[len(el.Tag)+2:]
				pagePrefix = ""
			}
			sec := model.Section{
				Heading:    model.Heading{Text: elStr},
				Paragraphs: []model.Paragraph{},
				Sections:   []model.Section{},
			}
			lvl := int(el.Tag[1] - '0')
			if sectionList[lvl-1] == nil {
				return nil, errs.New(fmt.Errorf("the difference between the heading level %s for heading '%s' and the heading before is greater than 1", el.Tag, sec.Heading.Text), nil)
			}
			*sectionList[lvl-1] = append(*sectionList[lvl-1], sec)
			s := &((*(sectionList[lvl-1]))[len(*sectionList[lvl-1])-1])
			currPars = &s.Paragraphs
			updateSectionList(&sectionList, &s.Sections, lvl)

		case "hu", "p", "table":
			if pagePrefix != "" {
				elStr = elStr[0:len(el.Tag)+2] + pagePrefix + elStr[len(el.Tag)+2:]
				pagePrefix = ""
			}
			*currPars = append(*currPars, model.Paragraph{Text: elStr})

		case "op":
			continue

		case "seite":
			pagePrefix = elStr

		default:
			return nil, errs.New(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
		}
	}
	return works, errs.Nil()
}

func findFootnotes(fussnoten *etree.Element) ([]model.Footnote, errs.UploadError) {
	if fussnoten == nil {
		return []model.Footnote{}, errs.Nil()
	}
	result := make([]model.Footnote, 0)
	for _, el := range fussnoten.ChildElements() {
		elStr, err := elemToString(el)
		if err != nil {
			return nil, errs.New(nil, fmt.Errorf("error writing element '%s' to string: %v", el.Tag, err))
		}
		result = append(result, model.Footnote{Text: elStr})
	}
	return result, errs.Nil()
}

func findSummaries(randtexte *etree.Element) ([]model.Summary, errs.UploadError) {
	if randtexte == nil {
		return []model.Summary{}, errs.Nil()
	}
	result := make([]model.Summary, 0)
	for _, el := range randtexte.ChildElements() {
		elStr, err := elemToString(el)
		if err != nil {
			return nil, errs.New(nil, fmt.Errorf("error writing element '%s' to string: %v", el.Tag, err))
		}
		result = append(result, model.Summary{Text: elStr})
	}
	return result, errs.Nil()
}

func elemToString(el *etree.Element) (string, error) {
	doc := etree.NewDocument()
	doc.SetRoot(el.Copy())
	return doc.WriteToString()
}

func updateSectionList(list *[]*[]model.Section, newSections *[]model.Section, level int) {
	(*list)[level] = newSections
	for i := level + 1; i < len(*list); i++ {
		(*list)[i] = nil
	}
}
