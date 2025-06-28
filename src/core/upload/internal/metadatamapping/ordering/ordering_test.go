package ordering

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/testutil"
	"github.com/stretchr/testify/assert"
)

func TestOrdering(t *testing.T) {
	testCases := []struct {
		name        string
		works       []model.Work
		expWorks    []model.Work
		expectError bool
	}{
		{
			name: "ordinal ordering",
			works: []model.Work{
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							FnRefs: []string{"1.1"},
						},
						Paragraphs: []model.Paragraph{
							{},
							{
								FnRefs:     []string{"2.1", "2.2"},
								SummaryRef: util.StrPtr("2.1"),
							},
							{SummaryRef: util.StrPtr("3.1")},
							{},
						},
						Sections: []model.Section{},
					}},
					Footnotes: []model.Footnote{
						{Ref: "1.1"},
						{Ref: "2.1"},
						{Ref: "2.2"},
					},
					Summaries: []model.Summary{{Ref: "2.1"}, {Ref: "3.1"}},
				},
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							FnRefs: []string{"4.1", "4.2"},
						},
						Paragraphs: []model.Paragraph{{
							FnRefs:     []string{"5.1"},
							SummaryRef: util.StrPtr("5.1"),
						}},
						Sections: []model.Section{},
					}},
					Footnotes: []model.Footnote{
						{Ref: "4.1"},
						{Ref: "4.2"},
						{Ref: "5.1"},
					},
					Summaries: []model.Summary{{Ref: "5.1"}},
				},
			},
			expWorks: []model.Work{
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							Ordinal: 1,
							FnRefs:  []string{"1.1"},
						},
						Paragraphs: []model.Paragraph{
							{Ordinal: 3},
							{
								Ordinal:    5,
								FnRefs:     []string{"2.1", "2.2"},
								SummaryRef: util.StrPtr("2.1"),
							},
							{Ordinal: 9, SummaryRef: util.StrPtr("3.1")},
							{Ordinal: 10},
						},
						Sections: []model.Section{},
					}},
					Footnotes: []model.Footnote{
						{Ordinal: 2, Ref: "1.1"},
						{Ordinal: 6, Ref: "2.1"},
						{Ordinal: 7, Ref: "2.2"},
					},
					Summaries: []model.Summary{
						{Ordinal: 4, Ref: "2.1"},
						{Ordinal: 8, Ref: "3.1"},
					},
				},
				{
					Sections: []model.Section{{
						Heading: model.Heading{
							Ordinal: 1,
							FnRefs:  []string{"4.1", "4.2"},
						},
						Paragraphs: []model.Paragraph{{
							Ordinal:    5,
							FnRefs:     []string{"5.1"},
							SummaryRef: util.StrPtr("5.1"),
						}},
						Sections: []model.Section{},
					}},
					Footnotes: []model.Footnote{
						{Ordinal: 2, Ref: "4.1"},
						{Ordinal: 3, Ref: "4.2"},
						{Ordinal: 6, Ref: "5.1"},
					},
					Summaries: []model.Summary{{Ordinal: 4, Ref: "5.1"}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Order(tc.works)
			if tc.expectError {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				testutil.AssertWorks(t, tc.expWorks, tc.works)
			}
		})
	}
}
