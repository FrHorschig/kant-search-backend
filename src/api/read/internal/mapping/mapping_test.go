package mapping

import (
	"reflect"
	"testing"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func TestVolumesToApiModels(t *testing.T) {
	in := []model.Volume{
		{
			VolumeNumber: 1,
			Title:        "Volume One",
			Works: []model.Work{
				{
					Ordinal: 1,
					Code:    "C1",
					Siglum:  util.StrPtr("abbr"),
					Title:   "The Work",
					Year:    "2024",
					Sections: []model.Section{
						{
							Heading:    1,
							Paragraphs: []int32{2, 3},
							Sections:   []model.Section{},
						},
					},
				},
			},
		},
	}
	expected := []models.Volume{
		{
			VolumeNumber: 1,
			Title:        "Volume One",
			Works: []models.Work{
				{
					Ordinal: 1,
					Code:    "C1",
					Siglum:  "abbr",
					Title:   "The Work",
					Year:    "2024",
					Sections: []models.Section{
						{
							Heading:    1,
							Paragraphs: []int32{2, 3},
							Sections:   []models.Section{},
						},
					},
				},
			},
		},
	}

	out := VolumesToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestFootnotesToApiModels(t *testing.T) {
	in := []model.Content{
		{Ordinal: 1, Ref: util.StrPtr("ref1"), FmtText: "Footnote text"},
	}
	expected := []models.Footnote{
		{Ordinal: 1, Ref: "ref1", Text: "Footnote text"},
	}

	out := FootnotesToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestHeadingsToApiModels(t *testing.T) {
	in := []model.Content{
		{Ordinal: 1, FmtText: "Heading text", TocText: util.StrPtr("toc text"), Pages: []int32{34}, FnRefs: []string{"fn1"}},
	}
	expected := []models.Heading{
		{Ordinal: 1, Text: "Heading text", TocText: "toc text", Pages: []int32{34}, FnRefs: []string{"fn1"}},
	}

	out := HeadingsToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestParagraphsToApiModels(t *testing.T) {
	in := []model.Content{
		{Ordinal: 1, FmtText: "Paragraph text", FnRefs: []string{"fn1"}, SummaryRef: util.StrPtr("s1")},
	}
	expected := []models.Paragraph{
		{Ordinal: 1, Text: "Paragraph text", FnRefs: []string{"fn1"}, SummaryRef: "s1"},
	}

	out := ParagraphsToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestSummariesToApiModels(t *testing.T) {
	in := []model.Content{
		{Ordinal: 1, Ref: util.StrPtr("ref1"), FmtText: "Summary text", FnRefs: []string{"fn1"}},
	}
	expected := []models.Summary{
		{Ordinal: 1, Ref: "ref1", Text: "Summary text", FnRefs: []string{"fn1"}},
	}

	out := SummariesToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}
