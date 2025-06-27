package testutil

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
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
