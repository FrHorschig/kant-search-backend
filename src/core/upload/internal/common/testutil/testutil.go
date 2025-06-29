package testutil

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/model"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func AssertWorks(t *testing.T, exp []model.Work, act []model.Work) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.NotNil(t, exp[i])
		assert.NotNil(t, act[i])
		assert.Equal(t, exp[i].Code, act[i].Code)
		if exp[i].Siglum != nil {
			assert.NotNil(t, act[i].Siglum)
			assert.Equal(t, *exp[i].Siglum, *act[i].Siglum)
		}
		assert.Equal(t, exp[i].Title, act[i].Title)
		assert.Equal(t, exp[i].Year, act[i].Year)
		assertSections(t, exp[i].Sections, act[i].Sections)
		AssertFootnotes(t, exp[i].Footnotes, act[i].Footnotes)
		AssertSummaries(t, exp[i].Summaries, act[i].Summaries)
	}
}

func assertSections(t *testing.T, exp []model.Section, act []model.Section) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i].Heading.Text, act[i].Heading.Text)
		assert.Equal(t, exp[i].Heading.TocText, act[i].Heading.TocText)
		assert.Equal(t, exp[i].Heading.Pages, act[i].Heading.Pages)
		assert.Equal(t, len(exp[i].Paragraphs), len(act[i].Paragraphs))
		assertParagraphs(t, exp[i].Paragraphs, act[i].Paragraphs)
		assert.Equal(t, len(exp[i].Sections), len(act[i].Sections))
		assertSections(t, exp[i].Sections, act[i].Sections)
	}
}

func assertParagraphs(t *testing.T, exp []model.Paragraph, act []model.Paragraph) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i].Text, act[i].Text)
		assert.ElementsMatch(t, exp[i].Pages, act[i].Pages)
		assert.ElementsMatch(t, exp[i].FnRefs, act[i].FnRefs)
	}
}

func AssertFootnotes(t *testing.T, exp []model.Footnote, act []model.Footnote) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i].Ref, act[i].Ref)
		assert.Equal(t, exp[i].Text, act[i].Text)
		assert.ElementsMatch(t, exp[i].Pages, act[i].Pages)
	}
}

func AssertSummaries(t *testing.T, exp []model.Summary, act []model.Summary) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i].Ref, act[i].Ref)
		assert.Equal(t, exp[i].Text, act[i].Text)
		assert.ElementsMatch(t, exp[i].Pages, act[i].Pages)
	}
}

func AssertDbVolume(t *testing.T, exp dbmodel.Volume, act dbmodel.Volume) {
	assert.Equal(t, exp.VolumeNumber, act.VolumeNumber)
	assert.Equal(t, exp.Title, act.Title)
	assertDbWorks(t, exp.Works, act.Works)
}

func assertDbWorks(t *testing.T, exp []dbmodel.Work, act []dbmodel.Work) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.NotNil(t, exp[i])
		assert.NotNil(t, act[i])
		assert.Equal(t, exp[i].Code, act[i].Code)
		if exp[i].Siglum != nil {
			assert.NotNil(t, act[i].Siglum)
			assert.Equal(t, *exp[i].Siglum, *act[i].Siglum)
		}
		assert.Equal(t, exp[i].Title, act[i].Title)
		assert.Equal(t, exp[i].Year, act[i].Year)
		assert.Equal(t, exp[i].Ordinal, act[i].Ordinal)
		assertDbSections(t, exp[i].Sections, act[i].Sections)
	}
}

func assertDbSections(t *testing.T, exp []dbmodel.Section, act []dbmodel.Section) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i], act[i])
		assertDbParagraphs(t, exp[i].Paragraphs, act[i].Paragraphs)
		assertDbSections(t, exp[i].Sections, act[i].Sections)
	}
}

func assertDbParagraphs(t *testing.T, exp []int32, act []int32) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i], act[i])
	}
}

func AssertDbContents(t *testing.T, exp []dbmodel.Content, act []dbmodel.Content) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i].FmtText, act[i].FmtText)
		if exp[i].TocText != nil {
			assert.NotNil(t, act[i].TocText)
			assert.Equal(t, *exp[i].TocText, *act[i].TocText)
		}
		assert.Equal(t, exp[i].SearchText, act[i].SearchText)
		assert.Equal(t, exp[i].Type, act[i].Type)
		assert.Equal(t, exp[i].Ordinal, act[i].Ordinal)
		assert.Equal(t, exp[i].WorkCode, act[i].WorkCode)
		assert.Equal(t, len(exp[i].Pages), len(act[i].Pages))
		for j := range exp[i].Pages {
			assert.Equal(t, exp[i].Pages[j], act[i].Pages[j])
		}
		assertContentMaps(t, exp[i], act[i])
		assertContentReferences(t, exp[i], act[i])
	}
}

func assertContentMaps(t *testing.T, exp dbmodel.Content, act dbmodel.Content) {
	assert.Equal(t, len(exp.PageByIndex), len(act.PageByIndex))
	for j := range exp.PageByIndex {
		assert.Equal(t, exp.PageByIndex[j].I, act.PageByIndex[j].I)
		assert.Equal(t, exp.PageByIndex[j].Num, act.PageByIndex[j].Num)
	}
	assert.Equal(t, len(exp.LineByIndex), len(act.LineByIndex))
	for j := range exp.LineByIndex {
		assert.Equal(t, exp.LineByIndex[j].I, act.LineByIndex[j].I)
		assert.Equal(t, exp.LineByIndex[j].Num, act.LineByIndex[j].Num)
	}
	assert.Equal(t, len(exp.WordIndexMap), len(act.WordIndexMap))
	for k, v := range exp.WordIndexMap {
		assert.Equal(t, v, act.WordIndexMap[k])
	}
}

func assertContentReferences(t *testing.T, exp dbmodel.Content, act dbmodel.Content) {
	assert.Equal(t, len(exp.FnRefs), len(act.FnRefs))
	for j := range exp.FnRefs {
		assert.Equal(t, exp.FnRefs[j], act.FnRefs[j])
	}
	if exp.SummaryRef != nil {
		assert.Equal(t, *exp.SummaryRef, *act.SummaryRef)
	}
	if exp.Ref != nil {
		assert.Equal(t, exp.Ref, act.Ref)
	}
}
