package ordering

import (
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
)

func Order(works []model.Work) errs.UploadError {
	for i := range works {
		ordinal := int32(0)
		fnByRef := findFootnoteByRef(works[i].Footnotes)
		summByRef := findSummaryByRef(works[i].Summaries)
		addParagraphOrdinals(works[i].Paragraphs, &ordinal, fnByRef, summByRef)
		addSectionOrdinals(works[i].Sections, &ordinal, fnByRef, summByRef)
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

func addSectionOrdinals(sections []model.Section, ordinal *int32, fnByRef map[string]*model.Footnote, summByRef map[string]*model.Summary) {
	for i := range sections {
		s := &sections[i]
		addHeadingOrdinals(&s.Heading, ordinal, fnByRef)
		addParagraphOrdinals(s.Paragraphs, ordinal, fnByRef, summByRef)
		addSectionOrdinals(s.Sections, ordinal, fnByRef, summByRef)
	}
}

func addHeadingOrdinals(heading *model.Heading, ordinal *int32, fnByRef map[string]*model.Footnote) {
	heading.Ordinal = *ordinal
	*ordinal += 1
	for i := range heading.FnRefs {
		fn := fnByRef[heading.FnRefs[i]]
		fn.Ordinal = *ordinal
		*ordinal += 1
	}
}
func addParagraphOrdinals(paragraphs []model.Paragraph, ordinal *int32, fnByRef map[string]*model.Footnote, summByRef map[string]*model.Summary) {
	for i := range paragraphs {
		p := &paragraphs[i]
		if p.SummaryRef != nil {
			summ := summByRef[*p.SummaryRef]
			summ.Ordinal = *ordinal
			*ordinal += 1
			for j := range summ.FnRefs {
				fn := fnByRef[summ.FnRefs[j]]
				fn.Ordinal = *ordinal
				*ordinal += 1
			}
		}
		p.Ordinal = *ordinal
		*ordinal += 1
		for j := range p.FnRefs {
			fn := fnByRef[p.FnRefs[j]]
			fn.Ordinal = *ordinal
			*ordinal += 1
		}
	}
}
