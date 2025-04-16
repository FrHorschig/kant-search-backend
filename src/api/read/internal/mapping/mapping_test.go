package mapping

import (
	"reflect"
	"testing"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
)

func TestVolumesToApiModels(t *testing.T) {
	in := []esmodel.Volume{
		{
			VolumeNumber: 1,
			Section:      2,
			Title:        "Volume One",
			Works: []esmodel.WorkRef{
				{Id: "w1", Code: "C1", Title: "Work One"},
			},
		},
	}
	expected := []models.Volume{
		{
			VolumeNumber: 1,
			Section:      2,
			Title:        "Volume One",
			Works: []models.WorkRef{
				{Id: "w1", Code: "C1", Title: "Work One"},
			},
		},
	}

	out := VolumesToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestWorkToApiModels(t *testing.T) {
	in := &esmodel.Work{
		Id:           "w1",
		Code:         "C1",
		Abbreviation: util.ToStrPtr("abbr"),
		Title:        "The Work",
		Year:         util.ToStrPtr("2024"),
		Sections: []esmodel.Section{
			{
				Heading:    "Section 1",
				Paragraphs: []string{"par1Id", "par2Id"},
				Sections:   []esmodel.Section{},
			},
		},
	}

	expected := models.Work{
		Id:           "w1",
		Code:         "C1",
		Abbreviation: "abbr",
		Title:        "The Work",
		Year:         "2024",
		Sections: []models.Section{
			{
				Heading:    "Section 1",
				Paragraphs: []string{"par1Id", "par2Id"},
				Sections:   []models.Section{},
			},
		},
	}

	out := WorkToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestFootnotesToApiModels(t *testing.T) {
	in := []esmodel.Content{
		{Id: "f1", Ref: util.ToStrPtr("ref1"), FmtText: "Footnote text"},
	}
	expected := []models.Footnote{
		{Id: "f1", Ref: "ref1", Text: "Footnote text"},
	}

	out := FootnotesToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestHeadingsToApiModels(t *testing.T) {
	in := []esmodel.Content{
		{Id: "h1", FmtText: "Heading text", TocText: util.ToStrPtr("toc text"), FnRefs: []string{"fn1"}},
	}
	expected := []models.Heading{
		{Id: "h1", Text: "Heading text", TocText: "toc text", FnRefs: []string{"fn1"}},
	}

	out := HeadingsToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestParagraphsToApiModels(t *testing.T) {
	in := []esmodel.Content{
		{Id: "p1", FmtText: "Paragraph text", FnRefs: []string{"fn1"}, SummaryRef: util.ToStrPtr("s1")},
	}
	expected := []models.Paragraph{
		{Id: "p1", Text: "Paragraph text", FnRefs: []string{"fn1"}, SummaryRef: "s1"},
	}

	out := ParagraphsToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}

func TestSummariesToApiModels(t *testing.T) {
	in := []esmodel.Content{
		{Id: "s1", Ref: util.ToStrPtr("ref1"), FmtText: "Summary text", FnRefs: []string{"fn1"}},
	}
	expected := []models.Summary{
		{Id: "s1", Ref: "ref1", Text: "Summary text", FnRefs: []string{"fn1"}},
	}

	out := SummariesToApiModels(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Expected %+v, got %+v", expected, out)
	}
}
