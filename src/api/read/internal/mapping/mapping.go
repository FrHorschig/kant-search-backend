package mapping

import (
	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
)

func VolumesToApiModels(in []esmodel.Volume) []models.Volume {
	out := []models.Volume{}
	for _, vIn := range in {
		vOut := models.Volume{
			VolumeNumber: vIn.VolumeNumber,
			Section:      vIn.Section,
			Title:        vIn.Title,
			Works:        []models.WorkRef{},
		}
		for _, wIn := range vIn.Works {
			vOut.Works = append(vOut.Works, models.WorkRef{
				Id:    wIn.Id,
				Code:  wIn.Code,
				Title: wIn.Title,
			})
		}
		out = append(out, vOut)
	}
	return out
}

func WorkToApiModels(in *esmodel.Work) models.Work {
	out := models.Work{
		Id:           in.Id,
		Code:         in.Code,
		Abbreviation: util.ToStrVal(in.Abbreviation),
		Title:        in.Title,
		Year:         util.ToStrVal(in.Year),
		Sections:     mapSections(in.Sections),
	}
	return out
}

func mapSections(in []esmodel.Section) []models.Section {
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

func FootnotesToApiModels(in []esmodel.Content) []models.Footnote {
	out := []models.Footnote{}
	for _, c := range in {
		out = append(out, models.Footnote{
			Id:   c.Id,
			Ref:  util.ToStrVal(c.Ref),
			Text: c.FmtText,
		})
	}
	return out
}

func HeadingsToApiModels(in []esmodel.Content) []models.Heading {
	out := []models.Heading{}
	for _, c := range in {
		out = append(out, models.Heading{
			Id:      c.Id,
			Text:    c.FmtText,
			TocText: util.ToStrVal(c.TocText),
			FnRefs:  c.FnRefs,
		})
	}
	return out
}

func ParagraphsToApiModels(in []esmodel.Content) []models.Paragraph {
	out := []models.Paragraph{}
	for _, c := range in {
		out = append(out, models.Paragraph{
			Id:         c.Id,
			Text:       c.FmtText,
			FnRefs:     c.FnRefs,
			SummaryRef: util.ToStrVal(c.SummaryRef),
		})
	}
	return out
}

func SummariesToApiModels(in []esmodel.Content) []models.Summary {
	out := []models.Summary{}
	for _, c := range in {
		out = append(out, models.Summary{
			Id:     c.Id,
			Ref:    util.ToStrVal(c.Ref),
			Text:   c.FmtText,
			FnRefs: c.FnRefs,
		})
	}
	return out
}
