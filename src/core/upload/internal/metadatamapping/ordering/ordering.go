package ordering

import (
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/model"
)

func Order(works []model.Work) errs.UploadError {
	for i := range works {
		ordinal := int32(0)
		fnByRef := findFootnoteByRef(works[i].Footnotes)
		summByRef := findSummaryByRef(works[i].Summaries)
		err := addParagraphOrdinals(works[i].Paragraphs, &ordinal, fnByRef, summByRef)
		if err.HasError {
			return err
		}
		err = addSectionOrdinals(works[i].Sections, &ordinal, fnByRef, summByRef)
		if err.HasError {
			return err
		}
	}

	return errs.Nil()
}

func findFootnoteByRef(footnotes []model.Footnote) map[string]*model.Footnote {
	result := make(map[string]*model.Footnote)
	for i := range footnotes {
		fn := &footnotes[i]
		result[fn.Ref] = fn
	}
	return result
}

func findSummaryByRef(summaries []model.Summary) map[string]*model.Summary {
	result := make(map[string]*model.Summary)
	for i := range summaries {
		summ := &summaries[i]
		result[summ.Ref] = summ
	}
	return result
}

func addSectionOrdinals(sections []model.Section, ordinal *int32, fnByRef map[string]*model.Footnote, summByRef map[string]*model.Summary) errs.UploadError {
	for i := range sections {
		s := &sections[i]
		err := addHeadingOrdinals(&s.Heading, ordinal, fnByRef)
		if err.HasError {
			return err
		}
		err = addParagraphOrdinals(s.Paragraphs, ordinal, fnByRef, summByRef)
		if err.HasError {
			return err
		}
		err = addSectionOrdinals(s.Sections, ordinal, fnByRef, summByRef)
		if err.HasError {
			return err
		}
	}
	return errs.Nil()
}

func addHeadingOrdinals(heading *model.Heading, ordinal *int32, fnByRef map[string]*model.Footnote) errs.UploadError {
	heading.Ordinal = *ordinal
	*ordinal += 1
	for _, ref := range heading.FnRefs {
		fn := fnByRef[ref]
		if fn == nil {
			return errs.New(fmt.Errorf("the footnote matching the heading footnote reference '%s' ('seite.nr') is missing", ref), nil)
		}
		fn.Ordinal = *ordinal
		*ordinal += 1
	}
	return errs.Nil()
}
func addParagraphOrdinals(paragraphs []model.Paragraph, ordinal *int32, fnByRef map[string]*model.Footnote, summByRef map[string]*model.Summary) errs.UploadError {
	for i := range paragraphs {
		p := &paragraphs[i]
		if p.SummaryRef != nil {
			summ := summByRef[*p.SummaryRef]
			summ.Ordinal = *ordinal
			*ordinal += 1
			for j := range summ.FnRefs {
				fnRef := summ.FnRefs[j]
				fn := fnByRef[fnRef]
				if fn == nil {
					return errs.New(fmt.Errorf("the footnote matching the summary footnote reference '%s' ('seite.nr') is missing", fnRef), nil)
				}
				fn.Ordinal = *ordinal
				fn.Ordinal = *ordinal
				*ordinal += 1
			}
		}
		p.Ordinal = *ordinal
		*ordinal += 1
		for j := range p.FnRefs {
			fnRef := p.FnRefs[j]
			fn := fnByRef[fnRef]
			if fn == nil {
				return errs.New(fmt.Errorf("the footnote matching the paragraph footnote reference '%s' ('seite.nr') is missing", fnRef), nil)
			}
			fn.Ordinal = *ordinal
			*ordinal += 1
		}
	}
	return errs.Nil()
}
