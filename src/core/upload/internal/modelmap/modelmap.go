package modelmap

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/core/upload/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/rs/zerolog/log"
)

func MapToModel(metadata metadata.VolumeMetadata, sections []model.TreeSection, summaries []model.TreeSummary, footnotes []model.TreeFootnote) ([]model.Work, errs.UploadError) {
	works := []model.Work{}
	latestPage := int32(1)
	for i, w := range sections {
		work, err := mapWork(w, metadata, i)
		if err.HasError {
			return nil, err
		}
		postprocessWork(&work, &latestPage)
		works = append(works, work)
	}
	// TODO (later) handle images and tables

	fns := []model.Footnote{}
	for _, f := range footnotes {
		fn, err := mapFootnote(f)
		if err.HasError {
			return nil, err
		}
		fns = append(fns, fn)
	}
	err := matchFnsToWorks(works, fns)
	if err.HasError {
		return nil, err
	}

	sms := []model.Summary{}
	for _, s := range summaries {
		summary, err := mapSummary(s)
		if err.HasError {
			return nil, err
		}
		sms = append(sms, summary)
	}
	err = mapSummariesToWorks(works, sms)
	if err.HasError {
		return nil, err
	}
	err = insertSummaryRefs(works)
	if err.HasError {
		return nil, err
	}

	return works, errs.Nil()
}

func mapWork(h0 model.TreeSection, metadata metadata.VolumeMetadata, index int) (model.Work, errs.UploadError) {
	work := model.Work{}
	work.Code = metadata.Works[index].Code
	work.Abbreviation = metadata.Works[index].Abbreviation
	work.Title = h0.Heading.TocTitle
	work.Year = &h0.Heading.Year
	for _, p := range h0.Paragraphs {
		par, err := mapParagraph(p)
		if err.HasError {
			return work, err
		}
		work.Paragraphs = append(work.Paragraphs, par)
	}
	for _, s := range h0.Sections {
		sec, err := mapSection(s)
		if err.HasError {
			return work, err
		}
		work.Sections = append(work.Sections, sec)
	}
	return work, errs.Nil()
}

func mapSection(s model.TreeSection) (model.Section, errs.UploadError) {
	section := model.Section{}
	heading, err := mapHeading(s.Heading)
	if err.HasError {
		return section, err
	}
	section.Heading = heading
	for _, par := range s.Paragraphs {
		dbPar, err := mapParagraph(par)
		if err.HasError {
			return section, err
		}
		section.Paragraphs = append(section.Paragraphs, dbPar)
	}
	for _, sec := range s.Sections {
		dbSec, err := mapSection(sec)
		if err.HasError {
			return section, err
		}
		section.Sections = append(section.Sections, dbSec)
	}
	return section, errs.Nil()
}

func mapHeading(h model.TreeHeading) (model.Heading, errs.UploadError) {
	pages, err := util.ExtractPages(h.TextTitle)
	if err.HasError {
		return model.Heading{}, err
	}
	heading := model.Heading{
		Text:    h.TextTitle,
		TocText: h.TocTitle,
		Pages:   pages,
		FnRefs:  util.ExtractFnRefs(h.TextTitle),
	}
	return heading, errs.Nil()
}

func mapParagraph(p string) (model.Paragraph, errs.UploadError) {
	pages, err := util.ExtractPages(p)
	if err.HasError {
		return model.Paragraph{}, err
	}
	paragraph := model.Paragraph{
		Text:   p,
		Pages:  pages,
		FnRefs: util.ExtractFnRefs(p),
	}
	return paragraph, errs.Nil()
}

func mapFootnote(f model.TreeFootnote) (model.Footnote, errs.UploadError) {
	pages, err := util.ExtractPages(f.Text)
	if err.HasError {
		return model.Footnote{}, err
	}
	if len(pages) == 0 {
		pages = []int32{f.Page}
	} else if !startsWithPageRef(f.Text, util.FmtPage(pages[0])) {
		pages = append([]int32{pages[0] - 1}, pages...)
	}
	if pages[0] != f.Page {
		return model.Footnote{}, errs.New(fmt.Errorf("footnote page %d does not match the first page of the footnote %d", f.Page, pages[0]), nil)
	}
	return model.Footnote{
		Ref:   fmt.Sprintf("%d.%d", f.Page, f.Nr),
		Pages: pages,
		Text:  f.Text,
	}, errs.Nil()
}

func mapSummary(s model.TreeSummary) (model.Summary, errs.UploadError) {
	pages, err := util.ExtractPages(s.Text)
	if err.HasError {
		return model.Summary{}, err
	}
	if len(pages) == 0 {
		pages = []int32{s.Page}
	} else if !startsWithPageRef(s.Text, util.FmtPage(pages[0])) {
		pages = append([]int32{pages[0] - 1}, pages...)
	}
	if pages[0] != s.Page {
		return model.Summary{}, errs.New(fmt.Errorf("summary page %d does not match the first page of the summary %d", s.Page, pages[0]), nil)
	}
	return model.Summary{
		Ref:    fmt.Sprintf("%d.%d", s.Page, s.Line),
		Text:   s.Text,
		Pages:  pages,
		FnRefs: util.ExtractFnRefs(s.Text),
	}, errs.Nil()
}

func postprocessWork(work *model.Work, latestPage *int32) {
	for i := range work.Paragraphs {
		postprocessParagraph(&work.Paragraphs[i], latestPage)
	}
	for i := range work.Sections {
		postprocessSection(&work.Sections[i], latestPage)
	}
}

func postprocessParagraph(par *model.Paragraph, latestPage *int32) {
	if len(par.Pages) > 0 {
		firstPage := par.Pages[0]
		pageRef := util.FmtPage(firstPage)
		if !startsWithPageRef(par.Text, pageRef) {
			par.Pages = append([]int32{firstPage - 1}, par.Pages...)
		}
		lastPage := par.Pages[len(par.Pages)-1]
		if lastPage > *latestPage {
			*latestPage = lastPage
		}

	} else {
		// this happens when a paragraph is fully inside a page and at least on line away from the page start and end
		par.Pages = []int32{*latestPage}
	}
}

func postprocessSection(section *model.Section, latestPage *int32) {
	head := &section.Heading
	if len(head.Pages) > 0 {
		firstPage := head.Pages[0]
		pageRef := util.FmtPage(firstPage)
		if !startsWithPageRef(head.Text, pageRef) {
			head.Pages = append([]int32{firstPage - 1}, head.Pages...)
		}
		lastPage := head.Pages[len(head.Pages)-1]
		if lastPage > *latestPage {
			*latestPage = lastPage
		}
	} else {
		// this happens when a heading is fully inside a page and at least on line away from the page start and end
		head.Pages = []int32{*latestPage}
	}

	for i := range section.Paragraphs {
		postprocessParagraph(&section.Paragraphs[i], latestPage)
	}

	for i := range section.Sections {
		postprocessSection(&section.Sections[i], latestPage)
	}
}

func matchFnsToWorks(works []model.Work, fns []model.Footnote) errs.UploadError {
	prevMax := int32(1)
	for i := range works {
		var min int32 = prevMax
		var max int32 = 1
		findMinMaxPages(works[i].Paragraphs, works[i].Sections, &min, &max)
		for j := range fns {
			pages := fns[j].Pages
			if pages[0] >= min && pages[len(pages)-1] <= max {
				works[i].Footnotes = append(works[i].Footnotes, fns[j])
			}
		}
		if min < prevMax {
			return errs.New(fmt.Errorf("minimum page number %d of work '%s' is smaller than the maximum page number %d of the previous work", min, works[i].Title, prevMax), nil)
		}
		prevMax = max
	}
	return errs.Nil()
}

func insertSummaryRefs(works []model.Work) errs.UploadError {
	for i := range works {
		w := &works[i]
		for j := range w.Summaries {
			summary := &w.Summaries[j]
			err := insertSummaryRef(summary, w.Sections)
			if err.HasError {
				return err
			}
		}
	}
	return errs.Nil()
}

func mapSummariesToWorks(works []model.Work, summaries []model.Summary) errs.UploadError {
	prevMax := int32(1)
	for i := range works {
		var min int32 = prevMax + 1
		var max int32 = 1
		findMinMaxPages(works[i].Paragraphs, works[i].Sections, &min, &max)
		for j := range summaries {
			pages := summaries[j].Pages
			if pages[0] >= min && pages[len(pages)-1] <= max {
				works[i].Summaries = append(works[i].Summaries, summaries[j])
			}
		}
		if min < prevMax {
			return errs.New(fmt.Errorf("minimum page number %d of work '%s' is smaller than the maximum page number %d of the previous work", min, works[i].Title, prevMax), nil)
		}
		prevMax = max
	}
	return errs.Nil()
}

func startsWithPageRef(text, pageRef string) bool {
	index := strings.Index(text, pageRef)
	cleaned := util.RemoveTags(text[:index])
	return cleaned == "" // in this case the text before page ref is only formatting code, so the "real" text starts with the page ref
}

func findMinMaxPages(paragraphs []model.Paragraph, sections []model.Section, min, max *int32) {
	for _, p := range paragraphs {
		if p.Pages[0] < *min {
			*min = p.Pages[0]
		}
		if p.Pages[len(p.Pages)-1] > *max {
			*max = p.Pages[len(p.Pages)-1]
		}
	}
	for _, s := range sections {
		if len(s.Heading.Pages) > 0 {
			if s.Heading.Pages[0] < *min {
				*min = s.Heading.Pages[0]
			}
			if s.Heading.Pages[len(s.Heading.Pages)-1] > *max {
				*max = s.Heading.Pages[len(s.Heading.Pages)-1]
			}
		}
		for _, p := range s.Paragraphs {
			if len(p.Pages) > 0 {
				if p.Pages[0] < *min {
					*min = p.Pages[0]
				}
				if p.Pages[len(p.Pages)-1] > *max {
					*max = p.Pages[len(p.Pages)-1]
				}
			}
		}
		findMinMaxPages([]model.Paragraph{}, s.Sections, min, max)
	}
}

func findPageLine(name string) (int32, int32) {
	pageLine := strings.Split(name, ".")
	// ignore errs, because we know the format
	page, _ := strconv.ParseInt(pageLine[0], 10, 32)
	line, _ := strconv.ParseInt(pageLine[1], 10, 32)
	return int32(page), int32(line)
}

func insertSummaryRef(summary *model.Summary, sections []model.Section) errs.UploadError {
	page, line := findPageLine(summary.Ref)
	p, err := findSummaryParagraph(summary, sections)
	if err.HasError {
		// in this case the summary starts in the middle of a paragraph, this is probably an error in the text, so we ignore the summary
		log.Debug().Msgf("found summary in the middle of a paragraph: %d.%d", page, line)
		return errs.Nil()
	}
	if p == nil {
		return errs.New(fmt.Errorf("could not find a paragraph for summary on page %d line %d", page, line), nil)
	}

	// duplicate page ref in the summary, so that summary and paragraph can be displayed independently from each other without loosing the page ref
	if line == 1 && !strings.Contains(summary.Text, util.FmtPage(page)) {
		summary.Text = util.FmtPage(page) + summary.Text
	}
	// line references should already by included in the summary text

	p.SummaryRef = &summary.Ref
	return errs.Nil()
}

func findSummaryParagraph(summary *model.Summary, sections []model.Section) (*model.Paragraph, errs.UploadError) {
	page, line := findPageLine(summary.Ref)
	for i := range sections {
		s := &sections[i]
		for iP := range s.Paragraphs {
			p := &s.Paragraphs[iP]
			ok, err := isSummaryParagraph(p, page, line)
			if err.HasError {
				return nil, err
			}
			if ok {
				return p, errs.Nil()
			}
		}
		for iS := range s.Sections {
			p, err := findSummaryParagraph(summary, s.Sections[iS].Sections)
			if err.HasError {
				return nil, err
			}
			if p != nil {
				return p, errs.Nil()
			}
		}
	}
	return nil, errs.Nil()
}

func isSummaryParagraph(p *model.Paragraph, page, line int32) (bool, errs.UploadError) {
	if !slices.Contains(p.Pages, page) {
		return false, errs.Nil()
	}
	pageIndex := strings.Index(p.Text, util.FmtPage(page))
	if pageIndex == -1 { // paragraph starts in the middle of the page
		pageIndex = 0
	}
	lineIndex := strings.Index(p.Text[pageIndex:], util.FmtLine(line))
	if lineIndex == -1 {
		return false, errs.Nil()
	}
	index := pageIndex + lineIndex + len(util.FmtLine(line))
	if !isSummaryAtStart(p.Text, index) {
		return false, errs.New(fmt.Errorf("summary on page %d line %d is not at the start of paragraph", page, line), nil)
	}
	return true, errs.Nil()
}

func isSummaryAtStart(text string, startIndex int) bool {
	cleaned := util.RemoveTags(text[:startIndex])
	return cleaned == "" // text before summary is only formatting code, so the "real text" starts with the summary
}
