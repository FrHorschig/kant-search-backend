package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/model_mapper.go -package=mocks

import (
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/errors"
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
		works = append(works, work)
	}
	mergeSummariesToWorks(works, summaries)
	// TODO handle images and tables

	fns := []dbmodel.Footnote{}
	for _, f := range footnotes {
		fns = append(fns, mapFootnote(f))
	}
	matchFnsToWorks(works, fns)

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
	work.Volume = vol
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
		dbPar := mapParagraph(par)
		section.Paragraphs = append(section.Paragraphs, dbPar)
	}
	for _, sec := range s.Sections {
		dbSec, err := mapSection(sec)
		if err.HasError {
			return section, err
		}
		dbSec.Parent = &section
		section.Sections = append(section.Sections, dbSec)
	}
	return section, errors.NilError()
}

func mapHeading(h model.Heading) (dbmodel.Heading, errors.ErrorNew) {
	heading := dbmodel.Heading{}
	lvl, err := mapLevel(h.Level)
	if err.HasError {
		return heading, err
	}
	heading.Level = lvl
	heading.TextTitle = h.TextTitle
	heading.TocTitle = h.TocTitle
	return heading, errors.NilError()
}

func mapLevel(lvl model.Level) (dbmodel.Level, errors.ErrorNew) {
	switch lvl {
	case model.H2:
		return dbmodel.H1, errors.NilError()
	case model.H3:
		return dbmodel.H2, errors.NilError()
	case model.H4:
		return dbmodel.H3, errors.NilError()
	case model.H5:
		return dbmodel.H4, errors.NilError()
	case model.H6:
		return dbmodel.H5, errors.NilError()
	case model.H7:
		return dbmodel.H7, errors.NilError()
	}
	return dbmodel.H1, errors.NewError(
		fmt.Errorf("invalid heading level %d", lvl),
		nil,
	)
}

func mapParagraph(p string) dbmodel.Paragraph {
	paragraph := dbmodel.Paragraph{}
	paragraph.Text = p
	paragraph.Pages = extractPages(p)
	paragraph.FnReferences = extractFnRefs(p)
	return paragraph
}

func extractFnRefs(text string) []string {
	// TODO
	return []string{}
}

func mergeSummariesToWorks(works []dbmodel.Work, summaries []model.Summary) {
	// TODO
}

func mapFootnote(f model.Footnote) dbmodel.Footnote {
	footnote := dbmodel.Footnote{}
	footnote.Name = fmt.Sprintf("%d.%d", f.Page, f.Nr)
	footnote.Pages = extractPages(f.Text)
	footnote.Text = f.Text
	return footnote
}

func extractPages(text string) []int32 {
	// TODO
	return []int32{}
}

func matchFnsToWorks(works []dbmodel.Work, fns []dbmodel.Footnote) {
	// TODO
}
