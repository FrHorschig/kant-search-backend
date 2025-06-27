package referencemapping

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/testutil"
	"github.com/stretchr/testify/assert"
)

func TestReferenceMapping(t *testing.T) {
	testCases := []struct {
		name        string
		works       []model.Work
		footnotes   []model.Footnote
		summaries   []model.Summary
		expWorks    []model.Work
		expectError bool
	}{
		{
			name: "basic reference mapping",
			works: []model.Work{
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							Pages:  []int32{2},
							FnRefs: []string{"2.3"},
						},
						Paragraphs: []model.Paragraph{{
							Text:  "<ks-meta-line>2</ks-meta-line> text",
							Pages: []int32{8},
						}},
						Sections: []model.Section{},
					}},
				},
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							Pages:  []int32{94},
							FnRefs: []string{"94.21"},
						},
						Paragraphs: []model.Paragraph{{
							Text:  "<ks-meta-line>3</ks-meta-line> bla",
							Pages: []int32{284},
						}},
						Sections: []model.Section{},
					}},
				},
			},
			footnotes: []model.Footnote{
				{Ref: "2.3", Pages: []int32{2}},
				{Ref: "94.21", Pages: []int32{94}},
				{Ref: "284.1", Pages: []int32{283}},
			},
			summaries: []model.Summary{
				{Ref: "8.2", Pages: []int32{8}},
				{Ref: "284.3", Pages: []int32{284}, FnRefs: []string{"284.1"}},
			},
			expWorks: []model.Work{
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							Pages:  []int32{2},
							FnRefs: []string{"2.3"},
						},
						Paragraphs: []model.Paragraph{{
							Text:       "<ks-meta-line>2</ks-meta-line> text",
							Pages:      []int32{8},
							SummaryRef: util.StrPtr("8.2"),
						}},
						Sections: []model.Section{},
					}},
					Footnotes: []model.Footnote{{
						Text: "", Ref: "2.3", Pages: []int32{2},
					}},
					Summaries: []model.Summary{{
						Text: "", Ref: "8.2", Pages: []int32{8},
					}},
				},
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							Pages:  []int32{94},
							FnRefs: []string{"94.21"},
						},
						Paragraphs: []model.Paragraph{{
							Text:       "<ks-meta-line>3</ks-meta-line> bla",
							Pages:      []int32{284},
							SummaryRef: util.StrPtr("284.3"),
						}},
						Sections: []model.Section{},
					}},
					Footnotes: []model.Footnote{
						{Ref: "94.21", Pages: []int32{94}},
						{Ref: "284.1", Pages: []int32{283}},
					},
					Summaries: []model.Summary{{
						Ref: "284.3", Pages: []int32{284}, FnRefs: []string{"284.1"},
					}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := MapReferences(tc.works, tc.footnotes, tc.summaries)
			if tc.expectError {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				testutil.AssertWorks(t, tc.expWorks, tc.works)
			}
		})
	}
}
