package modelmap

import (
	"testing"

	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestModelMapping(t *testing.T) {
	testCases := []struct {
		name        string
		volume      int32
		sections    []model.TreeSection
		summaries   []model.TreeSummary
		footnotes   []model.TreeFootnote
		model       []model.Work
		expectError bool
	}{
		{
			name:   "Merge summaries into paragraphs",
			volume: 1,
			sections: []model.TreeSection{
				{
					Heading: model.TreeHeading{Level: model.HWork, TextTitle: "work"},
					Sections: []model.TreeSection{{
						Heading: model.TreeHeading{Level: model.H1, TextTitle: "h1"},
						Paragraphs: []string{
							page(43) + line(348) + "I'm a paragraph.",
							line(58685) + "I'm a paragraph without a page number.",
						},
					}},
				},
				{
					Heading: model.TreeHeading{Level: model.HWork, TextTitle: "work2"},
					Sections: []model.TreeSection{{
						Heading: model.TreeHeading{Level: model.H1, TextTitle: page(102) + "2h1"},
						Paragraphs: []string{
							line(5) + "I'm a paragraph with " + page(483) + " a page break inside.",
						},
					}},
				},
			},
			summaries: []model.TreeSummary{
				{Page: 43, Line: 348, Text: "Summary 1"},
				{Page: 43, Line: 58685, Text: "Summary 2"},
				{Page: 482, Line: 5, Text: "Summary 3"},
			},
			model: []model.Work{
				{
					Title: "work",
					Sections: []model.Section{{
						Heading: model.Heading{Text: "h1", Pages: []int32{1}},
						Paragraphs: []model.Paragraph{
							{
								Text:  page(43) + line(348) + "I'm a paragraph.",
								Pages: []int32{43},
							},
							{
								Text:  line(58685) + "I'm a paragraph without a page number.",
								Pages: []int32{43},
							},
						},
					}},
					Summaries: []model.Summary{
						{Ref: "43.348", Text: "Summary 1", Pages: []int32{43}},
						{Ref: "43.58685", Text: "Summary 2", Pages: []int32{43}},
					},
				},
				{
					Title: "work2",
					Sections: []model.Section{{
						Heading: model.Heading{Text: page(102) + "2h1", Pages: []int32{102}},
						Paragraphs: []model.Paragraph{
							{
								Text:  line(5) + "I'm a paragraph with " + page(483) + " a page break inside.",
								Pages: []int32{482, 483},
							},
						},
					}},
					Summaries: []model.Summary{
						{Ref: "482.5", Text: "Summary 3", Pages: []int32{482}},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapToModel(tc.volume, tc.sections, tc.summaries, tc.footnotes)
			if tc.expectError {
				assert.True(t, err.HasError)
				assert.Nil(t, result)
			} else {
				assert.False(t, err.HasError)
				assert.Equal(t, len(tc.model), len(result))
				for i := range result {
					assertWork(t, tc.model[i], result[i])
				}
			}
		})
	}
}

func assertWork(t *testing.T, exp model.Work, act model.Work) {
	assert.NotNil(t, exp)
	assert.NotNil(t, act)
	assert.Equal(t, commonutil.StrVal(exp.Year), commonutil.StrVal(act.Year))

	assert.Equal(t, len(exp.Sections), len(act.Sections))
	for i := range exp.Sections {
		assertSections(t, exp.Sections[i], act.Sections[i])
	}
	assert.Equal(t, len(exp.Footnotes), len(act.Footnotes))
	for i := range exp.Footnotes {
		assertFootnote(t, exp.Footnotes[i], act.Footnotes[i])
	}
	assert.Equal(t, len(exp.Summaries), len(act.Summaries))
	for i := range exp.Summaries {
		assertSummary(t, exp.Summaries[i], act.Summaries[i])
	}
}

func assertSections(t *testing.T, exp model.Section, act model.Section) {
	assert.Equal(t, exp.Heading.Text, act.Heading.Text)
	assert.Equal(t, exp.Heading.TocText, act.Heading.TocText)
	assert.Equal(t, exp.Heading.Pages, act.Heading.Pages)
	assert.Equal(t, len(exp.Paragraphs), len(act.Paragraphs))
	for i := range exp.Paragraphs {
		assertParagraph(t, exp.Paragraphs[i], act.Paragraphs[i])
	}
	assert.Equal(t, len(exp.Sections), len(act.Sections))
	for i := range exp.Sections {
		assertSections(t, exp.Sections[i], act.Sections[i])
	}
}

func assertParagraph(t *testing.T, exp model.Paragraph, act model.Paragraph) {
	assert.Equal(t, exp.Text, act.Text)
	assert.ElementsMatch(t, exp.Pages, act.Pages)
	assert.ElementsMatch(t, exp.FnRefs, act.FnRefs)
}

func assertFootnote(t *testing.T, exp model.Footnote, act model.Footnote) {
	assert.Equal(t, exp.Ref, act.Ref)
	assert.Equal(t, exp.Text, act.Text)
	assert.ElementsMatch(t, exp.Pages, act.Pages)
}

func assertSummary(t *testing.T, exp model.Summary, act model.Summary) {
	assert.Equal(t, exp.Ref, act.Ref)
	assert.Equal(t, exp.Text, act.Text)
	assert.ElementsMatch(t, exp.Pages, act.Pages)
}

func line(line int32) string {
	return util.FmtLine(line)
}

func page(page int32) string {
	return util.FmtPage(page) + " "
}

func fnRef(page int32, nr int32) string {
	return util.FmtFnRef(page, nr)
}
