package mapping

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/extract"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemodel"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

func MapToModel(vol int32, sections []treemodel.Section, summaries []treemodel.Summary, footnotes []treemodel.Footnote) ([]model.Work, errors.ErrorNew) {
	works := []model.Work{}
	for i, w := range sections {
		work, err := mapWork(w, vol, i)
		if err.HasError {
			return nil, err
		}
		postprocessSections(&work)
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
	matchFnsToWorks(works, fns)

	sms := []model.Summary{}
	for _, s := range summaries {
		summary, err := mapSummary(s)
		if err.HasError {
			return nil, err
		}
		sms = append(sms, summary)
	}
	mapSummariesToWorks(works, sms)
	err := insertSummaryRefs(works)
	if err.HasError {
		return nil, err
	}

	return works, errors.NilError()
}

func mapWork(h0 treemodel.Section, vol int32, index int) (model.Work, errors.ErrorNew) {
	work := model.Work{}
	work.Code = Metadata[vol-1][index].Code
	work.Abbreviation = &Metadata[vol-1][index].Abbreviation
	work.Title = h0.Heading.TextTitle
	work.Year = &h0.Heading.Year
	if len(h0.Paragraphs) > 0 {
		return work, errors.NewError(fmt.Errorf("work has paragraphs before the first non-worktitle heading"), nil)
	}
	for _, s := range h0.Sections {
		sec, err := mapSection(s)
		if err.HasError {
			return work, err
		}
		work.Sections = append(work.Sections, sec)
	}
	return work, errors.NilError()
}

func mapSection(s treemodel.Section) (model.Section, errors.ErrorNew) {
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
	return section, errors.NilError()
}

func mapHeading(h treemodel.Heading) (model.Heading, errors.ErrorNew) {
	pages, err := extract.ExtractPages(h.TextTitle)
	if err.HasError {
		return model.Heading{}, err
	}
	heading := model.Heading{
		Text:    util.FmtHeading(int32(h.Level), h.TextTitle),
		TocText: h.TocTitle,
		Pages:   pages,
		FnRefs:  extract.ExtractFnRefs(h.TextTitle),
	}
	return heading, errors.NilError()
}

func mapParagraph(p string) (model.Paragraph, errors.ErrorNew) {
	pages, err := extract.ExtractPages(p)
	if err.HasError {
		return model.Paragraph{}, err
	}
	paragraph := model.Paragraph{
		Text:   p,
		Pages:  pages,
		FnRefs: extract.ExtractFnRefs(p),
	}
	return paragraph, errors.NilError()
}

func mapFootnote(f treemodel.Footnote) (model.Footnote, errors.ErrorNew) {
	pages, err := extract.ExtractPages(f.Text)
	if err.HasError {
		return model.Footnote{}, err
	}
	if len(pages) == 0 {
		pages = []int32{f.Page}
	} else if !startsWithPageRef(f.Text, util.FmtPage(pages[0])) {
		pages = append([]int32{pages[0] - 1}, pages...)
	}
	if pages[0] != f.Page {
		return model.Footnote{}, errors.NewError(fmt.Errorf("footnote page %d does not match the first page of the footnote %d", f.Page, pages[0]), nil)
	}
	return model.Footnote{
		Ref:   fmt.Sprintf("%d.%d", f.Page, f.Nr),
		Pages: pages,
		Text:  f.Text,
	}, errors.NilError()
}

func mapSummary(s treemodel.Summary) (model.Summary, errors.ErrorNew) {
	pages, err := extract.ExtractPages(s.Text)
	if err.HasError {
		return model.Summary{}, err
	}
	if len(pages) == 0 {
		pages = []int32{s.Page}
	} else if !startsWithPageRef(s.Text, util.FmtPage(pages[0])) {
		pages = append([]int32{pages[0] - 1}, pages...)
	}
	if pages[0] != s.Page {
		return model.Summary{}, errors.NewError(fmt.Errorf("summary page %d does not match the first page of the summary %d", s.Page, pages[0]), nil)
	}
	return model.Summary{
		Ref:    fmt.Sprintf("%d.%d", s.Page, s.Line),
		Text:   s.Text,
		Pages:  pages,
		FnRefs: extract.ExtractFnRefs(s.Text),
	}, errors.NilError()
}

func postprocessSections(work *model.Work) {
	var maxPage int32 = 1
	for _, sec := range work.Sections {
		postprocess(&sec, &maxPage)
	}
}

func postprocess(section *model.Section, latestPage *int32) {
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
		par := &section.Paragraphs[i]
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

	for i := range section.Sections {
		postprocess(&section.Sections[i], latestPage)
	}
}

func matchFnsToWorks(works []model.Work, fns []model.Footnote) {
	for i := range works {
		var min int32 = 0
		var max int32 = 0
		findMinMaxPages(works[i].Sections, &min, &max)
		for j := range fns {
			pages := fns[j].Pages
			if pages[0] >= min && pages[len(pages)-1] <= max {
				works[i].Footnotes = append(works[i].Footnotes, fns[j])
			}
		}
	}
}

func insertSummaryRefs(works []model.Work) errors.ErrorNew {
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
	return errors.NilError()
}

func mapSummariesToWorks(works []model.Work, summaries []model.Summary) {
	for i := range works {
		var min int32 = 0
		var max int32 = 0
		findMinMaxPages(works[i].Sections, &min, &max)
		for j := range summaries {
			pages := summaries[j].Pages
			if pages[0] >= min && pages[len(pages)-1] <= max {
				works[i].Summaries = append(works[i].Summaries, summaries[j])
			}
		}
	}
}

func startsWithPageRef(text, pageRef string) bool {
	index := strings.Index(text, pageRef)
	cleaned := extract.RemoveTags(text[:index])
	return cleaned == "" // text before page ref is only formatting code, so the "real text" starts with the page ref
}

func findMinMaxPages(sections []model.Section, min, max *int32) {
	for _, s := range sections {
		if len(s.Heading.Pages) > 0 {
			if *min == 0 || s.Heading.Pages[0] < *min {
				*min = s.Heading.Pages[0]
			}
			if s.Heading.Pages[len(s.Heading.Pages)-1] > *max {
				*max = s.Heading.Pages[len(s.Heading.Pages)-1]
			}
		}
		for _, p := range s.Paragraphs {
			if len(p.Pages) > 0 {
				if *min == 0 || p.Pages[0] < *min {
					*min = p.Pages[0]
				}
				if p.Pages[len(p.Pages)-1] > *max {
					*max = p.Pages[len(p.Pages)-1]
				}
			}
		}
		findMinMaxPages(s.Sections, min, max)
	}
}

func findPageLine(name string) (int32, int32) {
	pageLine := strings.Split(name, ".")
	// ignore errors, because we know the format
	page, _ := strconv.ParseInt(pageLine[0], 10, 32)
	line, _ := strconv.ParseInt(pageLine[1], 10, 32)
	return int32(page), int32(line)
}

func insertSummaryRef(summary *model.Summary, sections []model.Section) errors.ErrorNew {
	p, err := findSummaryParagraph(summary, sections)
	if err.HasError {
		return err
	}
	page, line := findPageLine(summary.Ref)
	if p == nil {
		return errors.NewError(fmt.Errorf("could not find a paragraph for summary on page %d line %d", page, line), nil)
	}

	// duplicate page ref in the summary, so that summary and paragraph can be displayed independently from each other without loosing the page ref
	if line == 1 && !strings.Contains(summary.Text, util.FmtPage(page)) {
		summary.Text = util.FmtPage(page) + summary.Text
	}
	// line references should already by included in the summary text

	p.Text = util.FmtSummaryRef(summary.Ref) + p.Text
	p.SummaryRef = &summary.Ref
	return errors.NilError()
}

func findSummaryParagraph(summary *model.Summary, sections []model.Section) (*model.Paragraph, errors.ErrorNew) {
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
				return p, errors.NilError()
			}
		}
		for iS := range s.Sections {
			p, err := findSummaryParagraph(summary, s.Sections[iS].Sections)
			if err.HasError {
				return nil, err
			}
			if p != nil {
				return p, errors.NilError()
			}
		}
	}
	return nil, errors.NilError()
}

func isSummaryParagraph(p *model.Paragraph, page, line int32) (bool, errors.ErrorNew) {
	if !slices.Contains(p.Pages, page) {
		return false, errors.NilError()
	}
	pageIndex := strings.Index(p.Text, util.FmtPage(page))
	if pageIndex == -1 { // paragraph starts in the middle of the page
		pageIndex = 0
	}
	lineIndex := strings.Index(p.Text[pageIndex:], util.FmtLine(line))
	if lineIndex == -1 {
		return false, errors.NilError()
	}
	index := pageIndex + lineIndex + len(util.FmtLine(line))
	if !isSummaryAtStart(p.Text, index) {
		return false, errors.NewError(fmt.Errorf("summary on page %d line %d is not at the start of paragraph", page, line), nil)
	}
	return true, errors.NilError()
}

func isSummaryAtStart(text string, startIndex int) bool {
	cleaned := extract.RemoveTags(text[:startIndex])
	return cleaned == "" // text before summary is only formatting code, so the "real text" starts with the summary
}
