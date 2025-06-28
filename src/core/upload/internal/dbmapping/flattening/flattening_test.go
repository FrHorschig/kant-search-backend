package flattening

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/testutil"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func TestFlattening(t *testing.T) {
	testCases := []struct {
		name       string
		volume     model.Volume
		works      []model.Work
		expVolume  dbmodel.Volume
		expContent []dbmodel.Content
	}{
		{
			name:   "map volume data",
			volume: model.Volume{VolumeNumber: 2, Title: "vol title"},
			works: []model.Work{
				{
					Paragraphs: []model.Paragraph{{Ordinal: 1}, {Ordinal: 2}},
					Sections: []model.Section{
						{
							Heading: model.Heading{Ordinal: 3},
							Paragraphs: []model.Paragraph{
								{Ordinal: 4},
								{Ordinal: 5},
								{Ordinal: 6},
							},
							Sections: []model.Section{{
								Heading:    model.Heading{Ordinal: 7},
								Paragraphs: []model.Paragraph{{Ordinal: 8}},
								Sections:   []model.Section{},
							}},
						},
						{
							Heading: model.Heading{Ordinal: 9},
							Paragraphs: []model.Paragraph{
								{Ordinal: 10},
								{Ordinal: 11},
							},
							Sections: []model.Section{},
						},
					},
				},
				{
					Paragraphs: []model.Paragraph{{Ordinal: 1}},
					Sections: []model.Section{
						{
							Heading:    model.Heading{Ordinal: 2},
							Paragraphs: []model.Paragraph{{Ordinal: 3}},
							Sections:   []model.Section{},
						},
					},
				},
			},
			expVolume: dbmodel.Volume{
				VolumeNumber: 2,
				Title:        "vol title",
				Works: []dbmodel.Work{
					{
						Ordinal:    1,
						Paragraphs: []int32{1, 2},
						Sections: []dbmodel.Section{
							{
								Heading:    3,
								Paragraphs: []int32{4, 5, 6},
								Sections: []dbmodel.Section{
									{
										Heading:    7,
										Paragraphs: []int32{8},
										Sections:   []dbmodel.Section{},
									},
								},
							},
							{
								Heading:    9,
								Paragraphs: []int32{10, 11},
								Sections:   []dbmodel.Section{},
							},
						},
					},
					{
						Ordinal:    2,
						Paragraphs: []int32{1},
						Sections: []dbmodel.Section{{
							Heading:    2,
							Paragraphs: []int32{3},
							Sections:   []dbmodel.Section{},
						}},
					},
				},
			},
			expContent: []dbmodel.Content{
				{Type: dbmodel.Paragraph, Ordinal: 1},
				{Type: dbmodel.Paragraph, Ordinal: 2},
				{Type: dbmodel.Heading, Ordinal: 3},
				{Type: dbmodel.Paragraph, Ordinal: 4},
				{Type: dbmodel.Paragraph, Ordinal: 5},
				{Type: dbmodel.Paragraph, Ordinal: 6},
				{Type: dbmodel.Heading, Ordinal: 7},
				{Type: dbmodel.Paragraph, Ordinal: 8},
				{Type: dbmodel.Heading, Ordinal: 9},
				{Type: dbmodel.Paragraph, Ordinal: 10},
				{Type: dbmodel.Paragraph, Ordinal: 11},
				{Type: dbmodel.Paragraph, Ordinal: 1},
				{Type: dbmodel.Heading, Ordinal: 2},
				{Type: dbmodel.Paragraph, Ordinal: 3},
			},
		},
		{
			name:   "map contents",
			volume: model.Volume{VolumeNumber: 2, Title: "vol title"},
			works: []model.Work{
				{
					Sections: []model.Section{
						{
							Heading: model.Heading{
								Ordinal: 1,
								Text:    "heading 1 text",
								TocText: "heading 1 toc text",
								Pages:   []int32{1},
								FnRefs:  []string{"1.1"},
							},
							Paragraphs: []model.Paragraph{{
								Ordinal:    2,
								Text:       "paragraph 2 text",
								Pages:      []int32{2},
								FnRefs:     []string{"2.1"},
								SummaryRef: util.StrPtr("2.2"),
							}},
							Sections: []model.Section{{
								Heading: model.Heading{
									Ordinal: 3,
									Text:    "heading 3 text",
									TocText: "heading 3 toc text",
									Pages:   []int32{3},
									FnRefs:  []string{"3.1"},
								},
								Paragraphs: []model.Paragraph{{
									Ordinal:    4,
									Text:       "paragraph 4 text",
									Pages:      []int32{4},
									FnRefs:     []string{"4.1"},
									SummaryRef: util.StrPtr("4.2"),
								}},
								Sections: []model.Section{},
							}},
						},
					},
					Footnotes: []model.Footnote{{
						Ordinal: 5,
						Ref:     "123.4",
						Text:    "footnote 5 text",
						Pages:   []int32{5},
					}},
					Summaries: []model.Summary{{
						Ordinal: 6,
						Ref:     "98.7",
						Text:    "summary 6 text",
						Pages:   []int32{6},
						FnRefs:  []string{"6.1"},
					}},
				},
			},
			expVolume: dbmodel.Volume{
				VolumeNumber: 2,
				Title:        "vol title",
				Works: []dbmodel.Work{
					{
						Ordinal: 1,
						Sections: []dbmodel.Section{
							{
								Heading:    1,
								Paragraphs: []int32{2},
								Sections: []dbmodel.Section{
									{
										Heading:    3,
										Paragraphs: []int32{4},
										Sections:   []dbmodel.Section{},
									},
								},
							},
						},
					},
				},
			},
			expContent: []dbmodel.Content{
				{
					Type:       dbmodel.Heading,
					Ordinal:    1,
					FmtText:    "heading 1 text",
					TocText:    util.StrPtr("heading 1 toc text"),
					SearchText: "heading 1 text",
					Pages:      []int32{1},
					FnRefs:     []string{"1.1"},
				},
				{
					Type:       dbmodel.Paragraph,
					Ordinal:    2,
					FmtText:    "paragraph 2 text",
					SearchText: "paragraph 2 text",
					Pages:      []int32{2},
					FnRefs:     []string{"2.1"},
					SummaryRef: util.StrPtr("2.2"),
				},
				{
					Type:       dbmodel.Heading,
					Ordinal:    3,
					FmtText:    "heading 3 text",
					SearchText: "heading 3 text",
					TocText:    util.StrPtr("heading 3 toc text"),
					Pages:      []int32{3},
					FnRefs:     []string{"3.1"},
				},
				{
					Type:       dbmodel.Paragraph,
					Ordinal:    4,
					FmtText:    "paragraph 4 text",
					SearchText: "paragraph 4 text",
					Pages:      []int32{4},
					FnRefs:     []string{"4.1"},
					SummaryRef: util.StrPtr("4.2"),
				},
				{
					Type:       dbmodel.Footnote,
					Ordinal:    5,
					Ref:        util.StrPtr("123.4"),
					FmtText:    "footnote 5 text",
					SearchText: "footnote 5 text",
					Pages:      []int32{5},
				},
				{
					Type:       dbmodel.Summary,
					Ordinal:    6,
					Ref:        util.StrPtr("98.7"),
					FmtText:    "summary 6 text",
					SearchText: "summary 6 text",
					Pages:      []int32{6},
					FnRefs:     []string{"6.1"},
				},
			},
		},
		{
			name:   "remove tags from search text",
			volume: model.Volume{VolumeNumber: 2, Title: "vol title"},
			works: []model.Work{
				{
					Paragraphs: []model.Paragraph{{Ordinal: 1, Text: "<op nr=\"2\"/>paragraph <ks-meta-page>2</ks-meta-page> with <some/> tags <end tag=\"with\"></attributes>"}},
					Sections:   []model.Section{},
				},
			},
			expVolume: dbmodel.Volume{
				VolumeNumber: 2,
				Title:        "vol title",
				Works: []dbmodel.Work{
					{
						Ordinal:    1,
						Paragraphs: []int32{1},
						Sections:   []dbmodel.Section{},
					},
				},
			},
			expContent: []dbmodel.Content{
				{
					Type:       dbmodel.Paragraph,
					Ordinal:    1,
					FmtText:    "<op nr=\"2\"/>paragraph <ks-meta-page>2</ks-meta-page> with <some/> tags <end tag=\"with\"></attributes>",
					SearchText: "paragraph with tags",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vols, contents := Flatten(tc.volume, tc.works)
			testutil.AssertDbVolumes(t, tc.expVolume, vols)
			testutil.AssertDbContents(t, tc.expContent, contents)
		})
	}
}
