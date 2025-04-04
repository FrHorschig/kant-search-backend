package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/model_mapper.go -package=mocks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/extract"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type ModelMapper interface {
	Map(
		volume int32,
		sections []model.Section,
		summaries []model.Summary,
		footnotes []model.Footnote,
	) (works []dbmodel.Work, err errors.ErrorNew)
}

type modelMapperImpl struct {
}

func NewModelMapper() ModelMapper {
	impl := modelMapperImpl{}
	return &impl
}

func (rec *modelMapperImpl) Map(vol int32, sections []model.Section, summaries []model.Summary, footnotes []model.Footnote) ([]dbmodel.Work, errors.ErrorNew) {
	works := []dbmodel.Work{}
	for i, w := range sections {
		work, err := mapWork(w, vol, i)
		if err.HasError {
			return nil, err
		}
		postprocessSectionPages(&work)
		works = append(works, work)
	}
	// TODO (later) handle images and tables

	fns := []dbmodel.Footnote{}
	for _, f := range footnotes {
		fn, err := mapFootnote(f)
		if err.HasError {
			return works, err
		}
		postprocessFootnotePages(&fn, f.Page)
		fns = append(fns, fn)
	}
	matchFnsToWorks(works, fns)

	sms := []dbmodel.Summary{}
	for _, s := range summaries {
		summary, err := mapSummary(s)
		if err.HasError {
			return works, err
		}
		sms = append(sms, summary)
	}
	mapSummariesToWorks(works, sms)
	insertSummaryRefs(works)

	return works, errors.NilError()
}

func mapWork(h0 model.Section, vol int32, index int) (dbmodel.Work, errors.ErrorNew) {
	work := dbmodel.Work{}
	work.Code = model.Metadata[vol-1][index].Code
	work.Abbreviation = &model.Metadata[vol-1][index].Abbreviation
	work.Title = h0.Heading.TextTitle
	work.Year = &h0.Heading.Year
	for _, s := range h0.Sections {
		sec, err := mapSection(s)
		if err.HasError {
			return work, err
		}
		work.Sections = append(work.Sections, sec)
	}
	return work, errors.NilError()
}

func mapSection(s model.Section) (dbmodel.Section, errors.ErrorNew) {
	section := dbmodel.Section{}
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

func mapHeading(h model.Heading) (dbmodel.Heading, errors.ErrorNew) {
	pages, err := extract.ExtractPages(h.TextTitle)
	if err.HasError {
		return dbmodel.Heading{}, err
	}
	heading := dbmodel.Heading{
		Text:    util.FmtHeading(int32(h.Level), h.TextTitle),
		TocText: h.TocTitle,
		Pages:   pages,
		FnRefs:  extract.ExtractFnRefs(h.TextTitle),
	}
	return heading, errors.NilError()
}

func mapParagraph(p string) (dbmodel.Paragraph, errors.ErrorNew) {
	pages, err := extract.ExtractPages(p)
	if err.HasError {
		return dbmodel.Paragraph{}, err
	}
	paragraph := dbmodel.Paragraph{
		Text:   p,
		Pages:  pages,
		FnRefs: extract.ExtractFnRefs(p),
	}
	return paragraph, errors.NilError()
}

func mapFootnote(f model.Footnote) (dbmodel.Footnote, errors.ErrorNew) {
	pages, err := extract.ExtractPages(f.Text)
	if err.HasError {
		return dbmodel.Footnote{}, err
	}
	return dbmodel.Footnote{
		Name:  fmt.Sprintf("%d.%d", f.Page, f.Nr),
		Pages: pages,
		Text:  f.Text,
	}, errors.NilError()
}

func mapSummary(s model.Summary) (dbmodel.Summary, errors.ErrorNew) {
	pages, err := extract.ExtractPages(s.Text)
	if err.HasError {
		return dbmodel.Summary{}, err
	}
	return dbmodel.Summary{
		Name:   fmt.Sprintf("%d.%d", s.Page, s.Line),
		Text:   s.Text,
		Pages:  pages,
		FnRefs: extract.ExtractFnRefs(s.Text),
	}, errors.NilError()
}

func postprocessSectionPages(work *dbmodel.Work) {
	var maxPage int32 = 1
	for _, sec := range work.Sections {
		processSection(&sec, &maxPage)
	}
}

func processSection(section *dbmodel.Section, maxPage *int32) {
	head := section.Heading
	if len(head.Pages) > 0 {
		firstPage := head.Pages[0]
		pageRef := util.FmtPage(firstPage)
		if !startsWithPageRef(head.Text, pageRef) {
			head.Pages = append([]int32{firstPage - 1}, head.Pages...)
		}
		lastPage := head.Pages[len(head.Pages)-1]
		if lastPage > *maxPage {
			*maxPage = lastPage
		}
	} else {
		// this happens when a heading is fully inside a page and at least on line away from the page start and end
		head.Pages = []int32{*maxPage}
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
			if lastPage > *maxPage {
				*maxPage = lastPage
			}

		} else {
			// this happens when a paragraph is fully inside a page and at least on line away from the page start and end
			par.Pages = []int32{*maxPage}
		}
	}

	for i := range section.Sections {
		processSection(&section.Sections[i], maxPage)
	}
}

func postprocessFootnotePages(fn *dbmodel.Footnote, fnStartPage int32) {
	if len(fn.Pages) > 0 {
		firstPage := fn.Pages[0]
		pageRef := util.FmtPage(firstPage)
		if !startsWithPageRef(fn.Text, pageRef) {
			fn.Pages = append([]int32{firstPage - 1}, fn.Pages...)
		}
	} else {
		fn.Pages = []int32{fnStartPage}
	}
}

func matchFnsToWorks(works []dbmodel.Work, fns []dbmodel.Footnote) {
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

func insertSummaryRefs(works []dbmodel.Work) {
	for i := range works {
		w := &works[i]
		for j := range w.Summaries {
			summary := &w.Summaries[j]
			page, line := findPageLine(summary.Name)
			insertSummaryRef(summary, page, line, w.Sections)
		}
	}
}

func mapSummariesToWorks(works []dbmodel.Work, summaries []dbmodel.Summary) {
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

func findMinMaxPages(sections []dbmodel.Section, min, max *int32) {
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

func insertSummaryRef(summary *dbmodel.Summary, page, line int32, sections []dbmodel.Section) errors.ErrorNew {
	for _, s := range sections {
		if len(s.Paragraphs) == 0 {
			return insertSummaryRef(summary, page, line, s.Sections)
		}

		for i := range s.Paragraphs {
			p := &s.Paragraphs[i]
			pageIndex := strings.Index(p.Text, util.FmtPage(page))
			if pageIndex == -1 {
				continue
			}
			lineIndex := strings.Index(p.Text[pageIndex:], util.FmtLine(line))
			if lineIndex == -1 {
				continue
			}
			index := pageIndex + lineIndex + len(util.FmtLine(line))

			// sanity check: summary should be at the start of the paragraph
			if !isSummaryAtStart(p.Text, index) {
				return errors.NewError(fmt.Errorf("summary on page %d line %d is not at the start of paragraph", page, line), nil)
			}
			if line == 1 {
				// move the page reference to the summary
				summary.Text = util.FmtPage(page) + summary.Text
				textWithoutPage := strings.Replace(p.Text, util.FmtPage(page), "", 1)
				p.Text = util.FmtSummaryRef(summary.Name) + textWithoutPage
			} else {
				p.Text = p.Text[:index] +
					util.FmtSummaryRef(summary.Name) +
					p.Text[index:]
			}
			return errors.NilError()
		}
		return insertSummaryRef(summary, page, line, s.Sections)
	}
	return errors.NewError(fmt.Errorf("could not find a paragraph for summary on page %d line %d", page, line), nil)
}

func isSummaryAtStart(text string, startIndex int) bool {
	cleaned := extract.RemoveTags(text[:startIndex])
	return cleaned == "" // text before summary is only formatting code, so the "real text" starts with the summary
}
