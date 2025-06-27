package textmapping

import (
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	trafo "github.com/frhorschig/kant-search-backend/core/upload/internal/xmlmapping/textmapping/texttransform"
)

func MapText(works []model.Work, footnotes []model.Footnote, summaries []model.Summary) errs.UploadError {
	err := mapWorks(works)
	if err.HasError {
		return err
	}
	err = mapSummaries(summaries)
	if err.HasError {
		return err
	}
	err = mapFootnotes(footnotes)
	if err.HasError {
		return err
	}
	return errs.Nil()
}

func mapWorks(works []model.Work) errs.UploadError {
	for i := range works {
		_, tocText, err := trafo.Hx(works[i].Title)
		if err.HasError {
			return err
		}
		works[i].Title = tocText
		for j := range works[i].Paragraphs {
			mapParagraph(&works[i].Paragraphs[j])
		}
		for j := range works[i].Sections {
			mapSection(&works[i].Sections[j])
		}
	}
	return errs.Nil()
}

func mapParagraph(p *model.Paragraph) errs.UploadError {
	pText, err := p.Text, errs.Nil()
	if strings.HasPrefix(pText, "<hu>") {
		pText, err = trafo.Hu(pText)
	} else if strings.HasPrefix(pText, "<table>") {
		pText, err = trafo.Table(pText)
	} else {
		pText, err = trafo.P(pText)
	}
	if err.HasError {
		return err
	}
	p.Text = pText
	return errs.Nil()
}

func mapSection(s *model.Section) errs.UploadError {
	fmtText, tocText, err := trafo.Hx(s.Heading.Text)
	if err.HasError {
		return err
	}
	s.Heading.Text = fmtText
	s.Heading.TocText = tocText

	for i := range s.Paragraphs {
		mapParagraph(&s.Paragraphs[i])
	}
	for i := range s.Sections {
		mapSection(&s.Sections[i])
	}
	return errs.Nil()
}

func mapFootnotes(footnotes []model.Footnote) errs.UploadError {
	for i := range footnotes {
		text, ref, err := trafo.Footnote(footnotes[i].Text)
		if err.HasError {
			return err
		}
		footnotes[i].Text = text
		footnotes[i].Ref = ref
	}
	return errs.Nil()
}
func mapSummaries(summaries []model.Summary) errs.UploadError {
	for i := range summaries {
		text, ref, err := trafo.Summary(summaries[i].Text)
		if err.HasError {
			return err
		}
		summaries[i].Text = text
		summaries[i].Ref = ref
	}
	return errs.Nil()
}
