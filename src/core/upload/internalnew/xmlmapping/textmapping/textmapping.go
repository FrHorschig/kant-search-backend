package textmapping

import (
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	trafo "github.com/frhorschig/kant-search-backend/core/upload/internalnew/xmlmapping/textmapping/texttransform"
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
		mapParagraphs(works[i].Paragraphs)
		mapSections(works[i].Sections)
	}
	return errs.Nil()
}

func mapSections(sections []model.Section) errs.UploadError {
	for i := range sections {
		s := &sections[i]
		mapHeading(&s.Heading)
		mapParagraphs(s.Paragraphs)
		mapSections(s.Sections)
	}
	return errs.Nil()
}

func mapHeading(heading *model.Heading) errs.UploadError {
	fmtText, tocText, err := trafo.Hx(heading.Text)
	if err.HasError {
		return err
	}
	heading.Text = fmtText
	heading.TocText = tocText
	return errs.Nil()
}

func mapParagraphs(paragraphs []model.Paragraph) errs.UploadError {
	for i := range paragraphs {
		p := paragraphs[i]
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
