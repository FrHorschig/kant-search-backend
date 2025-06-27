package textmapping

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTextMapping(t *testing.T) {
	testCases := []struct {
		name         string
		works        []model.Work
		footnotes    []model.Footnote
		summaries    []model.Summary
		expWorks     []model.Work
		expFootnotes []model.Footnote
		expSummaries []model.Summary
		expectError  bool
	}{
		{
			name: "paragraph and heading transforms",
			works: []model.Work{
				{
					Title: "<h1> work 1 </h1>",
					Sections: []model.Section{{
						Heading: model.Heading{Text: "<h2> heading 2 </h2>"},
						Paragraphs: []model.Paragraph{
							{Text: "<p><seite nr=\"4\"/> p paragraph </p>"},
							{Text: "<hu> hu paragraph </hu>"},
							{Text: "<table> table paragraph </table>"},
						},
						Sections: []model.Section{{
							Heading:    model.Heading{Text: "<h3> heading 3 </h3>"},
							Paragraphs: []model.Paragraph{},
							Sections:   []model.Section{},
						}},
					}},
				},
			},
			expWorks: []model.Work{
				{
					Title: "Work 1",
					Sections: []model.Section{{
						Heading: model.Heading{
							Text:    "heading 2",
							TocText: "Heading 2",
						},
						Paragraphs: []model.Paragraph{
							{Text: "<ks-meta-page>4</ks-meta-page> p paragraph"},
							{Text: "hu paragraph"},
							{Text: ""},
						},
						Sections: []model.Section{{
							Heading: model.Heading{
								Text:    "heading 3",
								TocText: "Heading 3",
							},
							Paragraphs: []model.Paragraph{},
							Sections:   []model.Section{},
						}},
					}},
				},
			},
		},
		{
			name: "footnote and summary transforms",
			footnotes: []model.Footnote{
				{Text: "<fn seite=\"6\" nr=\"7\"> footnote 1 </fn>"},
				{Text: "<fn seite=\"8\" nr=\"9\"> footnote 2 </fn>"},
			},
			summaries: []model.Summary{
				{Text: "<randtext seite=\"2\" anfang=\"3\"> randtext 1 </randtext>"},
				{Text: "<randtext seite=\"4\" anfang=\"5\"> randtext 2 </randtext>"},
			},
			expFootnotes: []model.Footnote{
				{Text: "footnote 1", Ref: "6.7"},
				{Text: "footnote 2", Ref: "8.9"},
			},
			expSummaries: []model.Summary{
				{Text: "randtext 1", Ref: "2.3"},
				{Text: "randtext 2", Ref: "4.5"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := MapText(tc.works, tc.footnotes, tc.summaries)
			if tc.expectError {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				testutil.AssertWorks(t, tc.expWorks, tc.works)
				testutil.AssertFootnotes(t, tc.expFootnotes, tc.footnotes)
				testutil.AssertSummaries(t, tc.expSummaries, tc.summaries)
			}
		})
	}
}
