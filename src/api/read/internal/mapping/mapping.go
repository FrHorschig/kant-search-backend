package mapping

import (
	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func VolumesToApiModels(in []model.Volume) []models.Volume {
	out := []models.Volume{}
	for _, vIn := range in {
		vOut := models.Volume{
			VolumeNumber: vIn.VolumeNumber,
			Title:        vIn.Title,
			Works:        []models.Work{},
		}
		for _, wIn := range vIn.Works {
			vOut.Works = append(vOut.Works,
				models.Work{
					Ordinal:      wIn.Ordinal,
					Code:         wIn.Code,
					Abbreviation: util.StrVal(wIn.Abbreviation),
					Title:        wIn.Title,
					Year:         wIn.Year,
					Paragraphs:   wIn.Paragraphs,
					Sections:     mapSections(wIn.Sections),
				},
			)
		}
		out = append(out, vOut)
	}
	return out
}

func mapSections(in []model.Section) []models.Section {
	out := []models.Section{}
	for _, sIn := range in {
		out = append(out, models.Section{
			Heading:    sIn.Heading,
			Paragraphs: sIn.Paragraphs,
			Sections:   mapSections(sIn.Sections),
		})
	}
	return out
}

func FootnotesToApiModels(in []model.Content) []models.Footnote {
	out := []models.Footnote{}
	for _, c := range in {
		out = append(out, models.Footnote{
			Ordinal: c.Ordinal,
			Ref:     util.StrVal(c.Ref),
			Text:    c.FmtText,
		})
	}
	return out
}

func HeadingsToApiModels(in []model.Content) []models.Heading {
	out := []models.Heading{}
	for _, c := range in {
		out = append(out, models.Heading{
			Ordinal: c.Ordinal,
			Text:    c.FmtText,
			TocText: util.StrVal(c.TocText),
			Pages:   c.Pages,
			FnRefs:  c.FnRefs,
		})
	}
	return out
}

func ParagraphsToApiModels(in []model.Content) []models.Paragraph {
	out := []models.Paragraph{}
	for _, c := range in {
		out = append(out, models.Paragraph{
			Ordinal:    c.Ordinal,
			Text:       c.FmtText,
			FnRefs:     c.FnRefs,
			SummaryRef: util.StrVal(c.SummaryRef),
		})
	}
	return out
}

func SummariesToApiModels(in []model.Content) []models.Summary {
	out := []models.Summary{}
	for _, c := range in {
		out = append(out, models.Summary{
			Ordinal: c.Ordinal,
			Ref:     util.StrVal(c.Ref),
			Text:    c.FmtText,
			FnRefs:  c.FnRefs,
		})
	}
	return out
}
