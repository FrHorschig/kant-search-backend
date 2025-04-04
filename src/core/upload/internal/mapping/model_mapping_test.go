package mapping

import (
	"fmt"
	"testing"

	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestModelMapping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name        string
		volume      int32
		sections    []model.Section
		summaries   []model.Summary
		footnotes   []model.Footnote
		model       []dbmodel.Work
		expectError bool
	}{
		{
			name:   "Multiple sections with simple content are mapped",
			volume: 1,
			sections: []model.Section{{
				Heading: model.Heading{
					Level:     model.HWork,
					TocTitle:  "work title TOC text",
					TextTitle: "work title",
					Year:      "1724",
				},
				Paragraphs: []string{"paragraph 1 text", "some other paragraph text"},
				Sections: []model.Section{
					{
						Heading: model.Heading{
							Level:     model.H1,
							TocTitle:  "subsection 1 title TOC text",
							TextTitle: "subsection 1 title",
							Year:      "",
						},
						Paragraphs: []string{"subsection 1 paragraph text", "second text"},
					},
					{
						Heading: model.Heading{
							Level:     model.H1,
							TocTitle:  "subsection 2 title TOC text",
							TextTitle: "subsection 2 title",
							Year:      "",
						},
						Paragraphs: []string{"subsection 2 paragraph text", "another second text"},
					},
				},
			}},
			model: []dbmodel.Work{{
				Title: "work title",
				Year:  commonutil.ToStrPtr("1724"),
				Sections: []dbmodel.Section{{
					Heading: dbmodel.Heading{
						Text:    "work title",
						TocText: "work title TOC text",
					},
					Paragraphs: []dbmodel.Paragraph{
						{Text: "paragraph 1 text"},
						{Text: "some other paragraph text"},
					},
					Sections: []dbmodel.Section{
						{
							Heading: dbmodel.Heading{Text: "h1"},
							Paragraphs: []dbmodel.Paragraph{
								{Text: "subsection 1 paragraph text"},
								{Text: "second text"},
							},
						},
						{
							Heading: dbmodel.Heading{Text: "h1"},
							Paragraphs: []dbmodel.Paragraph{
								{Text: "subsection 2 paragraph text"},
								{Text: "another second text"},
							},
						},
					},
				}},
			}},
		},
		{
			name:   "Multiple nested works and sections are mapped",
			volume: 2,
			sections: []model.Section{
				{
					Heading: model.Heading{Level: model.HWork},
					Sections: []model.Section{
						{
							Heading: model.Heading{Level: model.H1},
							Sections: []model.Section{
								{
									Heading: model.Heading{Level: model.H2},
									Sections: []model.Section{
										{
											Heading: model.Heading{Level: model.H3},
											Sections: []model.Section{
												{
													Heading: model.Heading{Level: model.H4},
													Sections: []model.Section{
														{
															Heading: model.Heading{Level: model.H5},
															Sections: []model.Section{
																{
																	Heading: model.Heading{Level: model.H6},
																},
															},
														},
													},
												},
											},
										},
										{Heading: model.Heading{Level: model.H3}},
									},
								},
							},
						},
						{Heading: model.Heading{Level: model.H2}},
						{Heading: model.Heading{Level: model.H2}},
					},
				},
				{Heading: model.Heading{Level: model.HWork}},
				{Heading: model.Heading{Level: model.HWork}},
			},
			model: []dbmodel.Work{
				{
					Sections: []dbmodel.Section{{
						Heading: dbmodel.Heading{Text: "h1"},
						Sections: []dbmodel.Section{
							{
								Heading: dbmodel.Heading{Text: "h2"},
								Sections: []dbmodel.Section{
									{
										Heading: dbmodel.Heading{Text: "h3"},
										Sections: []dbmodel.Section{
											{
												Heading: dbmodel.Heading{Text: "h4"},
												Sections: []dbmodel.Section{
													{
														Heading: dbmodel.Heading{Text: "h5"},
														Sections: []dbmodel.Section{
															{
																Heading: dbmodel.Heading{Text: "h6"},
															},
														},
													},
												},
											},
											{Heading: dbmodel.Heading{Text: "h3"}},
										},
									},
									{Heading: dbmodel.Heading{Text: "h2"}},
									{Heading: dbmodel.Heading{Text: "h2"}},
								},
							},
						},
					}},
				},
				{Sections: []dbmodel.Section{{Heading: dbmodel.Heading{Text: "h1"}}}},
				{Sections: []dbmodel.Section{{Heading: dbmodel.Heading{Text: "h1"}}}},
			},
		},
		{
			name:   "Extract pages and footnote references from paragraphs",
			volume: 3,
			sections: []model.Section{{
				Heading: model.Heading{Level: model.HWork},
				Paragraphs: []string{
					fnRef(2, 64) + page(2),
					"This " + fnRef(83, 3) + "is a" + page(254) + " test text.",
					"It " + fnRef(582, 1) + " continues " + page(942) + fnRef(298481, 2485) + page(942) + " in the " + fnRef(3, 5281) + " next" + page(943) + "paragraph.",
					page(23) + fnRef(4, 23)},
				Sections: []model.Section{
					{
						Heading: model.Heading{Level: model.H1},
						Paragraphs: []string{
							page(28471) + "This paragraph" + fnRef(4, 2) + "starts with a page",
							"This paragraph ends with a page." + fnRef(482, 148) + page(3),
						},
					},
				},
			}},
			model: []dbmodel.Work{{
				Sections: []dbmodel.Section{{
					Paragraphs: []dbmodel.Paragraph{
						{
							Text:   fnRef(2, 64) + page(2),
							Pages:  []int32{2},
							FnRefs: []string{"2.64"},
						},
						{
							Text:   "This " + fnRef(83, 3) + "is a" + page(254) + " test text.",
							Pages:  []int32{254},
							FnRefs: []string{"83.3"},
						},
						{
							Text:   "It " + fnRef(582, 1) + " continues " + page(942) + fnRef(298481, 2485) + page(942) + " in the " + fnRef(3, 5281) + " next" + page(943) + "paragraph.",
							Pages:  []int32{941, 942, 943},
							FnRefs: []string{"582.1", "298481.2485", "3.5281"},
						},
						{
							Text:   page(23) + fnRef(4, 23),
							Pages:  []int32{23},
							FnRefs: []string{"4.23"},
						},
					},
					Sections: []dbmodel.Section{
						{
							Paragraphs: []dbmodel.Paragraph{
								{
									Text:   page(28471) + "This paragraph" + fnRef(4, 3) + " starts with a page.",
									Pages:  []int32{28471},
									FnRefs: []string{"4.3"},
								},
								{
									Text:   "This paragraph ends with a page." + fnRef(482, 148) + page(3),
									Pages:  []int32{3},
									FnRefs: []string{"482.148"},
								},
							},
						},
					},
				}},
			}},
		},
		{
			name:   "Merge summaries into paragraphs",
			volume: 4,
			sections: []model.Section{{
				Heading: model.Heading{Level: model.HWork},
				Paragraphs: []string{
					page(43) + line(348) + "I'm a paragraph.",
					"I'm a paragraph with the line " + line(58685) + " number at the end",
				},
				Sections: []model.Section{
					{
						Heading: model.Heading{Level: model.H1},
						Paragraphs: []string{
							page(95) + line(123) + "I'm a paragraph." + page(96) + line(1) + "I'm another paragraph.",
							"I'm a paragraph with the line " + page(483) + line(2) + " number at the end " + page(484) + line(5) + "that continues over multiple pages.",
						},
					},
				},
			}},
			summaries: []model.Summary{
				{Page: 43, Line: 348, Text: "Summary 1"},
				{Page: 43, Line: 58685, Text: "Summary 2"},
				{Page: 484, Line: 5, Text: "Summary 3"},
			},
			model: []dbmodel.Work{{
				Sections: []dbmodel.Section{{
					Paragraphs: []dbmodel.Paragraph{
						{
							Text:  page(43) + line(348) + sumRef(43, 348) + "I'm  line." + line(58685) + sumRef(43, 58685) + "I'm a second line.",
							Pages: []int32{43},
						},
					},
					Sections: []dbmodel.Section{
						{
							Paragraphs: []dbmodel.Paragraph{
								{
									Text:  "I'm a paragraph with the line " + page(483) + line(5) + " number at the end " + page(484) + line(5) + sumRef(484, 5) + "that continues over multiple pages.",
									Pages: []int32{482, 483, 484},
								},
							},
						},
					},
				}},
				Summaries: []dbmodel.Summary{
					{Name: "43.348", Text: "Summary 1", Pages: []int32{43}},
					{Name: "43.58685", Text: "Summary 2", Pages: []int32{43}},
					{Name: "484.5", Text: "Summary 3", Pages: []int32{484}},
				},
			}},
		},
		{
			name: "Map footnote name",
			footnotes: []model.Footnote{
				{Page: 2, Nr: 5, Text: "This is a footnote."},
				{Page: 4, Nr: 20, Text: "This is a 2nd footnote."},
			},
			model: []dbmodel.Work{{
				Footnotes: []dbmodel.Footnote{
					{Name: "2.5", Text: "This is a footnote.", Pages: []int32{2}},
					{Name: "4.20", Text: "This is a 2nd footnote.", Pages: []int32{4}},
				},
			}},
		},
		{
			name: "Extract pages from footnote text",
			footnotes: []model.Footnote{
				{
					Page: 881,
					Nr:   284,
					Text: "This " + page(9582) + "is a " + page(383) + "footnote.",
				},
				{
					Page: 2,
					Nr:   9,
					Text: "This " + page(30) + "is a 2nd footnote.",
				},
			},
			model: []dbmodel.Work{{
				Footnotes: []dbmodel.Footnote{
					{
						Name:  "881.284",
						Text:  "This " + page(9582) + "is a " + page(383) + "footnote.",
						Pages: []int32{881, 9582, 383},
					},
					{
						Name:  "2.9",
						Text:  "This " + page(30) + "is a 2nd footnote.",
						Pages: []int32{2, 30},
					},
				},
			}},
		},
	}

	sut := &modelMapperImpl{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := sut.Map(tc.volume, tc.sections, tc.summaries, tc.footnotes)
			if tc.expectError {
				assert.True(t, err.HasError)
				assert.Nil(t, result)
			}
			assert.Equal(t, len(tc.model), len(result))
			for i := range tc.model {
				assertWork(t, tc.model[i], result[i])
			}
		})
	}
}

func assertWork(t *testing.T, exp dbmodel.Work, act dbmodel.Work) {
	assert.NotNil(t, exp)
	assert.NotNil(t, act)

	assert.Equal(t, exp.Code, act.Code)
	assert.Equal(t, commonutil.ToStrVal(exp.Abbreviation), commonutil.ToStrVal(act.Abbreviation))
	assert.Equal(t, commonutil.ToStrVal(exp.Year), commonutil.ToStrVal(act.Year))

	assert.Equal(t, len(exp.Sections), len(act.Sections))
	for j := range exp.Sections {
		assertSections(t, exp.Sections[j], act.Sections[j])
	}
	assert.Equal(t, len(exp.Footnotes), len(act.Footnotes))
	for j := range exp.Sections {
		assertFootnote(t, exp.Footnotes[j], act.Footnotes[j])
	}
}

func assertSections(t *testing.T, exp dbmodel.Section, act dbmodel.Section) {
	assert.Equal(t, exp.Heading.Text, act.Heading.Text)
	assert.Equal(t, exp.Heading.TocText, act.Heading.TocText)
	assert.Equal(t, len(exp.Paragraphs), len(act.Paragraphs))
	for i := range exp.Paragraphs {
		assertParagraph(t, exp.Paragraphs[i], act.Paragraphs[i])
	}
	assert.Equal(t, len(exp.Sections), len(act.Sections))
	for i := range exp.Sections {
		assertSections(t, exp.Sections[i], act.Sections[i])
	}
}

func assertParagraph(t *testing.T, exp dbmodel.Paragraph, act dbmodel.Paragraph) {
	assert.Equal(t, exp.Text, act.Text)
	assert.ElementsMatch(t, exp.Pages, act.Pages)
	assert.ElementsMatch(t, exp.FnRefs, act.FnRefs)
}

func assertFootnote(t *testing.T, exp dbmodel.Footnote, act dbmodel.Footnote) {
	assert.Equal(t, exp.Name, act.Name)
	assert.Equal(t, exp.Text, act.Text)
	assert.ElementsMatch(t, exp.Pages, act.Pages)
}

func line(line int32) string {
	return util.FmtLine(line)
}

func page(page int32) string {
	return util.FmtPage(page)
}

func fnRef(page int32, nr int32) string {
	return util.FmtFnRef(page, nr)
}

func sumRef(page, line int32) string {
	return util.FmtSummaryRef(fmt.Sprintf("%d.%d", page, line))
}
