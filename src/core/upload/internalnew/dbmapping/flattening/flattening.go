package flattening

import (
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/util"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func Flatten(volume model.Volume, works []model.Work) (dbmodel.Volume, []dbmodel.Content) {
	dbWorks := mapWorks(works)
	contents := mapContents(works)
	return dbmodel.Volume{
		VolumeNumber: volume.VolumeNumber, Title: volume.Title, Works: dbWorks,
	}, contents
}

func mapWorks(works []model.Work) []dbmodel.Work {
	results := make([]dbmodel.Work, len(works))
	for i, w := range works {
		pars := mapParagraphOrdinals(w.Paragraphs)
		secs := mapSectionOrdinals(w.Sections)
		results[i] = dbmodel.Work{
			Ordinal:    int32(i + 1),
			Code:       w.Code,
			Siglum:     w.Siglum,
			Title:      w.Title,
			Year:       w.Year,
			Paragraphs: pars,
			Sections:   secs,
		}
	}
	return results
}

func mapSectionOrdinals(sections []model.Section) []dbmodel.Section {
	results := make([]dbmodel.Section, len(sections))
	for i, s := range sections {
		results[i] = dbmodel.Section{
			Heading:    s.Heading.Ordinal,
			Paragraphs: mapParagraphOrdinals(s.Paragraphs),
			Sections:   mapSectionOrdinals(s.Sections),
		}
	}
	return results
}

func mapParagraphOrdinals(paragraphs []model.Paragraph) []int32 {
	results := make([]int32, len(paragraphs))
	for i, p := range paragraphs {
		results[i] = p.Ordinal
	}
	return results
}

func mapContents(works []model.Work) []dbmodel.Content {
	contents := []dbmodel.Content{}
	for _, w := range works {
		addParagraphs(w.Paragraphs, &contents, w.Code)
		addSections(w.Sections, &contents, w.Code)
		addFootnotes(w.Footnotes, &contents, w.Code)
		addSummaries(w.Summaries, &contents, w.Code)
	}
	return contents
}

func addParagraphs(paragraphs []model.Paragraph, contents *[]dbmodel.Content, workCode string) {
	for _, p := range paragraphs {
		*contents = append(*contents, dbmodel.Content{
			Type:       dbmodel.Paragraph,
			Ordinal:    p.Ordinal,
			FmtText:    p.Text,
			SearchText: util.RemoveTags(p.Text),
			Pages:      p.Pages,
			FnRefs:     p.FnRefs,
			SummaryRef: p.SummaryRef,
			WorkCode:   workCode,
		})
	}
}

func addSections(sections []model.Section, contents *[]dbmodel.Content, workCode string) {
	for _, s := range sections {
		h := s.Heading
		*contents = append(*contents, dbmodel.Content{
			Type:       dbmodel.Heading,
			Ordinal:    h.Ordinal,
			FmtText:    h.Text,
			TocText:    &h.TocText,
			SearchText: util.RemoveTags(h.Text),
			Pages:      h.Pages,
			FnRefs:     h.FnRefs,
			WorkCode:   workCode,
		})
		addParagraphs(s.Paragraphs, contents, workCode)
		addSections(s.Sections, contents, workCode)
	}
}

func addFootnotes(footnotes []model.Footnote, contents *[]dbmodel.Content, workCode string) {
	for _, f := range footnotes {
		*contents = append(*contents, dbmodel.Content{
			Type:       dbmodel.Footnote,
			Ordinal:    f.Ordinal,
			Ref:        &f.Ref,
			FmtText:    f.Text,
			SearchText: util.RemoveTags(f.Text),
			Pages:      f.Pages,
			WorkCode:   workCode,
		})
	}
}

func addSummaries(summaries []model.Summary, contents *[]dbmodel.Content, workCode string) {
	for _, s := range summaries {
		*contents = append(*contents, dbmodel.Content{
			Type:       dbmodel.Summary,
			Ordinal:    s.Ordinal,
			Ref:        &s.Ref,
			FmtText:    s.Text,
			SearchText: util.RemoveTags(s.Text),
			Pages:      s.Pages,
			FnRefs:     s.FnRefs,
			WorkCode:   workCode,
		})
	}
}
