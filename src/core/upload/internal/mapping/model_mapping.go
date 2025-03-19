package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/model_mapper.go -package=mocks

import (
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type ModelMapper interface {
	Map([]model.Section, []model.Summary, []model.Footnote) ([]dbmodel.Work, errors.ErrorNew)
}

type modelMapperImpl struct {
}

func NewModelMapper() ModelMapper {
	impl := modelMapperImpl{}
	return &impl
}

func (rec *modelMapperImpl) Map(sections []model.Section, summaries []model.Summary, footnotes []model.Footnote) ([]dbmodel.Work, errors.ErrorNew) {
	// TODO don't forget handling of images and tables
	works := []dbmodel.Work{}
	for _, s := range sections {
		works = append(works, mapWork(s, summaries))
	}
	footnotes := []dbmodel.Footnote{}
	for _, f := range footnotes {
		footnotes = append(footnotes, mapFootnote(f))
	}
	return nil, errors.NilError()
}

func mapWork(h0 model.Section, summaries []model.Summary) dbmodel.Work {
	work := dbmodel.Work{}
	// code
	// title
	// abbrev
	// year
	// volume
	for _, s := range h0.Sections {
		sec := mapSection(s)
		work.Sections = append(work.Sections, sec)
	}
	return work
}

func mapLevel(lvl model.Level) dbmodel.Level {
	// TODO: handle levels H7 and H8 with some kind of error
	switch lvl {
	case model.H1:
		return dbmodel.H1
	case model.H2:
		return dbmodel.H1
	case model.H3:
		return dbmodel.H2
	case model.H4:
		return dbmodel.H3
	case model.H5:
		return dbmodel.H4
	case model.H6:
		return dbmodel.H5
	}
	return dbmodel.H6 // model.H7
}

func mapSection(s model.Section) dbmodel.Section {
	section := dbmodel.Section{}
	// level
	// tocTitle
	// textTitle
	for _, par := range s.Paragraphs {
		dbPar := mapParagraph(par)
		section.Paragraphs = append(section.Paragraphs, dbPar)
	}
	for _, sec := range s.Sections {
		dbSec := mapSection(sec)
		section.Sections = append(section.Sections, dbSec)
	}
	return section
}

func mapParagraph(p string) dbmodel.Paragraph {
	paragraph := dbmodel.Paragraph{}
	paragraph.Text = p
	// pages
	// fnRefs
	// sentences
	return paragraph
}

func mapFootnote(f model.Footnote) dbmodel.Footnote {
	footnote := dbmodel.Footnote{}
	footnote.Name = fmt.Sprintf("%s.%s", f.Page, f.Nr)
	footnote.Pages = extractPages(f.Text)
	footnote.Text = f.Text
	return footnote
}

func extractPages(text string) []int32 {
	// TODO
	return []int32{}
}
