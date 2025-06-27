package metadataextraction

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetadataExtraction(t *testing.T) {
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
			name: "paragraph and heading page extraction",
			works: []model.Work{
				{
					Title: "Work 1",
					Sections: []model.Section{{
						Heading: model.Heading{
							Text:    "<ks-meta-page>4</ks-meta-page>heading 2",
							TocText: "Heading 2",
						},
						Paragraphs: []model.Paragraph{{
							Text: "par <ks-meta-page>5</ks-meta-page> text",
						}},
						Sections: []model.Section{},
					}},
				},
			},
			expWorks: []model.Work{
				{
					Title: "Work 1",
					Sections: []model.Section{{
						Heading: model.Heading{
							Text:    "<ks-meta-page>4</ks-meta-page>heading 2",
							TocText: "Heading 2",
							Pages:   []int32{4},
						},
						Paragraphs: []model.Paragraph{{
							Text:  "par <ks-meta-page>5</ks-meta-page> text",
							Pages: []int32{4, 5},
						}},
						Sections: []model.Section{},
					}},
				},
			},
		},
		{
			name: "footnote and summary page extraction",
			footnotes: []model.Footnote{
				{Text: "footnote <ks-meta-page>7</ks-meta-page> 1", Ref: "6.7"},
			},
			summaries: []model.Summary{
				{Text: "randtext <ks-meta-page>3</ks-meta-page> 1", Ref: "2.3"},
			},
			expFootnotes: []model.Footnote{
				{
					Text:  "footnote <ks-meta-page>7</ks-meta-page> 1",
					Ref:   "6.7",
					Pages: []int32{6, 7},
				},
			},
			expSummaries: []model.Summary{
				{
					Text:  "randtext <ks-meta-page>3</ks-meta-page> 1",
					Ref:   "2.3",
					Pages: []int32{2, 3},
				},
			},
		},
		{
			name: "paragraph and heading fnRef extraction",
			works: []model.Work{
				{
					Title: "Work 1",
					Sections: []model.Section{{
						Heading: model.Heading{
							Text:    "heading 2<ks-meta-fnref>23.5</ks-meta-fnref>",
							TocText: "Heading 2",
						},
						Paragraphs: []model.Paragraph{{
							Text: "par <ks-meta-fnref>283.17</ks-meta-fnref> text",
						}},
						Sections: []model.Section{},
					}},
				},
			},
			expWorks: []model.Work{
				{
					Title: "Work 1",
					Sections: []model.Section{{
						Heading: model.Heading{
							Text:    "heading 2<ks-meta-fnref>23.5</ks-meta-fnref>",
							TocText: "Heading 2",
							Pages:   []int32{1},
							FnRefs:  []string{"23.5"},
						},
						Paragraphs: []model.Paragraph{{
							Text:   "par <ks-meta-fnref>283.17</ks-meta-fnref> text",
							Pages:  []int32{1},
							FnRefs: []string{"283.17"},
						}},
						Sections: []model.Section{},
					}},
				},
			},
		},
		{
			name: "footnote and summary fnRef extraction",
			footnotes: []model.Footnote{
				{Text: "footnote <ks-meta-fnref>2842.218</ks-meta-fnref> 1", Ref: "6.7"},
			},
			summaries: []model.Summary{
				{Text: "randtext <ks-meta-fnref>3.2</ks-meta-fnref> 1", Ref: "2.3"},
			},
			expFootnotes: []model.Footnote{
				{
					Text:  "footnote <ks-meta-fnref>2842.218</ks-meta-fnref> 1",
					Ref:   "6.7",
					Pages: []int32{6},
				},
			},
			expSummaries: []model.Summary{
				{
					Text:   "randtext <ks-meta-fnref>3.2</ks-meta-fnref> 1",
					Ref:    "2.3",
					Pages:  []int32{2},
					FnRefs: []string{"2.3"},
				},
			},
		},
		{
			name: "footnote textPage-refPage mismatch",
			footnotes: []model.Footnote{
				{Text: "footnote <ks-meta-page>2842</ks-meta-page> 1", Ref: "6.7"},
			},
			expectError: true,
		},
		{
			name: "summary textPage-refPage mismatch",
			summaries: []model.Summary{
				{Text: "randtext <ks-meta-page>8</ks-meta-page> 1", Ref: "2.3"},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ExtractMetadata(tc.works, tc.footnotes, tc.summaries)
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
