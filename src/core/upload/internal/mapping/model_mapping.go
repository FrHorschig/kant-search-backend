package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/model_mapper.go -package=mocks

import (
	"fmt"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/extract"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
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
	// TODO handle images and tables

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
		insertSummaryRef(summary, works)
		sms = append(sms, summary)
	}
	mapSummariesToWorks(works, sms)

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
		Text:    fmt.Sprintf(model.HeadingFmt, h.Level, h.TextTitle, h.Level),
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
		pageRef := fmt.Sprintf(model.PageFmt, firstPage)
		if !startsWithPageRef(head.Text, pageRef) {
			head.Pages = append([]int32{firstPage - 1}, head.Pages...)
		}
		lastPage := head.Pages[len(head.Pages)-1]
		if lastPage > *maxPage {
			*maxPage = lastPage
		}
	} else {
		head.Pages = []int32{*maxPage}
	}

	for i := range section.Paragraphs {
		par := &section.Paragraphs[i]
		if len(par.Pages) > 0 {
			firstPage := par.Pages[0]
			pageRef := fmt.Sprintf(model.PageFmt, firstPage)
			if !startsWithPageRef(par.Text, pageRef) {
				par.Pages = append([]int32{firstPage - 1}, par.Pages...)
			}
			lastPage := par.Pages[len(par.Pages)-1]
			if lastPage > *maxPage {
				*maxPage = lastPage
			}

		} else {
			// This happens when a paragraph is fully inside a page and does not start at the beginning of the page.
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
		pageRef := fmt.Sprintf(model.PageFmt, firstPage)
		if !startsWithPageRef(fn.Text, pageRef) {
			fn.Pages = append([]int32{firstPage - 1}, fn.Pages...)
		}
	} else {
		fn.Pages = []int32{fnStartPage}
	}
}

func matchFnsToWorks(works []dbmodel.Work, fns []dbmodel.Footnote) {

}

func insertSummaryRef(summary dbmodel.Summary, works []dbmodel.Work) {
	// TODO
}

func mapSummariesToWorks(works []dbmodel.Work, summaries []dbmodel.Summary) {
	// TODO
}

func startsWithPageRef(text, pageRef string) bool {
	index := strings.Index(text, pageRef)
	cleaned := extract.RemoveTags(text[:index])
	return cleaned == "" // text before page ref is only formatting code, so the "real text" starts with the page ref
}
