package mapping

import (
	"testing"

	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemodel"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestModelMapping(t *testing.T) {
	testCases := []struct {
		name        string
		volume      int32
		sections    []treemodel.Section
		summaries   []treemodel.Summary
		footnotes   []treemodel.Footnote
		model       []model.Work
		expectError bool
	}{
		{
			name:   "Multiple sections with simple content are mapped",
			volume: 1,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{
					Level:     treemodel.HWork,
					TocTitle:  "work title TOC text",
					TextTitle: "work title",
					Year:      "1724",
				},
				Sections: []treemodel.Section{
					{
						Heading: treemodel.Heading{
							Level:     treemodel.H1,
							TocTitle:  "subsection 1 title TOC text",
							TextTitle: "subsection 1 title",
						},
						Paragraphs: []string{"subsection 1 par text", "second text"},
					},
					{
						Heading: treemodel.Heading{
							Level:     treemodel.H1,
							TocTitle:  "subsection 2 title TOC text",
							TextTitle: "subsection 2 title",
						},
						Paragraphs: []string{"subsection 2 par text", "another second text"},
					},
				},
			}},
			model: []model.Work{{
				Title: "work title",
				Year:  commonutil.StrPtr("1724"),
				Sections: []model.Section{
					{
						Heading: model.Heading{
							Text:    util.FmtHeading(1, "subsection 1 title"),
							TocText: "subsection 1 title TOC text",
							Pages:   []int32{1},
						},
						Paragraphs: []model.Paragraph{
							{Text: "subsection 1 par text", Pages: []int32{1}},
							{Text: "second text", Pages: []int32{1}},
						},
					},
					{
						Heading: model.Heading{
							Text:    util.FmtHeading(1, "subsection 2 title"),
							TocText: "subsection 2 title TOC text",
							Pages:   []int32{1},
						},
						Paragraphs: []model.Paragraph{
							{Text: "subsection 2 par text", Pages: []int32{1}},
							{Text: "another second text", Pages: []int32{1}},
						},
					},
				},
			}},
		},
		{
			name:   "Paragraphs before the first non-work heading",
			volume: 1,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{
					Level:     treemodel.HWork,
					TocTitle:  "work title TOC text",
					TextTitle: "work title",
					Year:      "1724",
				},
				Paragraphs: []string{"paragraph 1 text", "some other paragraph text"},
				Sections:   []treemodel.Section{},
			}},
			expectError: true,
		},
		{
			name:   "Multiple nested works and sections are mapped",
			volume: 2,
			sections: []treemodel.Section{
				{
					Heading: treemodel.Heading{Level: treemodel.HWork},
					Sections: []treemodel.Section{
						{
							Heading: treemodel.Heading{Level: treemodel.H1},
							Sections: []treemodel.Section{
								{
									Heading: treemodel.Heading{Level: treemodel.H2},
									Sections: []treemodel.Section{
										{
											Heading: treemodel.Heading{Level: treemodel.H3},
											Sections: []treemodel.Section{
												{
													Heading: treemodel.Heading{Level: treemodel.H4},
													Sections: []treemodel.Section{
														{
															Heading: treemodel.Heading{Level: treemodel.H5},
															Sections: []treemodel.Section{
																{
																	Heading: treemodel.Heading{Level: treemodel.H6},
																},
															},
														},
													},
												},
											},
										},
										{Heading: treemodel.Heading{Level: treemodel.H3}},
									},
								},
								{Heading: treemodel.Heading{Level: treemodel.H2}},
								{Heading: treemodel.Heading{Level: treemodel.H2}},
							},
						},
					},
				},
				{
					Heading:  treemodel.Heading{Level: treemodel.HWork},
					Sections: []treemodel.Section{{Heading: treemodel.Heading{Level: treemodel.H1}}},
				},
				{
					Heading:  treemodel.Heading{Level: treemodel.HWork},
					Sections: []treemodel.Section{{Heading: treemodel.Heading{Level: treemodel.H1}}},
				},
			},
			model: []model.Work{
				{
					Sections: []model.Section{{
						Heading: model.Heading{Text: util.FmtHeading(1, ""), Pages: []int32{1}},
						Sections: []model.Section{
							{
								Heading: model.Heading{Text: util.FmtHeading(2, ""), Pages: []int32{1}},
								Sections: []model.Section{
									{
										Heading: model.Heading{Text: util.FmtHeading(3, ""), Pages: []int32{1}},
										Sections: []model.Section{
											{
												Heading: model.Heading{Text: util.FmtHeading(4, ""), Pages: []int32{1}},
												Sections: []model.Section{
													{
														Heading: model.Heading{Text: util.FmtHeading(5, ""), Pages: []int32{1}},
														Sections: []model.Section{
															{
																Heading: model.Heading{Text: util.FmtHeading(6, ""), Pages: []int32{1}},
															},
														},
													},
												},
											},
										},
									},
									{Heading: model.Heading{Text: util.FmtHeading(3, ""), Pages: []int32{1}}},
								},
							},
							{Heading: model.Heading{Text: util.FmtHeading(2, ""), Pages: []int32{1}}},
							{Heading: model.Heading{Text: util.FmtHeading(2, ""), Pages: []int32{1}}},
						},
					}},
				},
				{Sections: []model.Section{{Heading: model.Heading{Text: util.FmtHeading(1, "")}}}},
				{Sections: []model.Section{{Heading: model.Heading{Text: util.FmtHeading(1, "")}}}},
			},
		},
		{
			name:   "Extract pages and footnote references from paragraphs",
			volume: 3,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{Level: treemodel.HWork},
				Sections: []treemodel.Section{
					{
						Heading: treemodel.Heading{Level: treemodel.H1},
						Paragraphs: []string{
							"This paragraph ends with a page." + fnRef(482, 148) + page(3),
							page(82) + "This paragraph" + fnRef(4, 2) + "starts with a page.",
						},
						Sections: []treemodel.Section{{
							Heading: treemodel.Heading{Level: treemodel.H2},
							Paragraphs: []string{
								fnRef(2, 64) + page(120),
								"This " + fnRef(83, 3) + "is a" + page(254) + " test text.",
								"It " + fnRef(582, 1) + " continues " + page(941) + fnRef(298481, 2485) + page(942) + " in the " + fnRef(3, 5281) + " next" + page(943) + "paragraph.",
								page(12840) + fnRef(4, 23)},
						}},
					},
				},
			}},
			model: []model.Work{{
				Sections: []model.Section{{
					Heading: model.Heading{Text: util.FmtHeading(1, "")},
					Paragraphs: []model.Paragraph{
						{
							Text:   "This paragraph ends with a page." + fnRef(482, 148) + page(3),
							Pages:  []int32{2, 3},
							FnRefs: []string{"482.148"},
						},
						{
							Text:   page(82) + "This paragraph" + fnRef(4, 2) + "starts with a page.",
							Pages:  []int32{82},
							FnRefs: []string{"4.2"},
						},
					},
					Sections: []model.Section{
						{
							Heading: model.Heading{Text: util.FmtHeading(2, "")},
							Paragraphs: []model.Paragraph{
								{
									Text:   fnRef(2, 64) + page(120),
									Pages:  []int32{120},
									FnRefs: []string{"2.64"},
								},
								{
									Text:   "This " + fnRef(83, 3) + "is a" + page(254) + " test text.",
									Pages:  []int32{253, 254},
									FnRefs: []string{"83.3"},
								},
								{
									Text:   "It " + fnRef(582, 1) + " continues " + page(941) + fnRef(298481, 2485) + page(942) + " in the " + fnRef(3, 5281) + " next" + page(943) + "paragraph.",
									Pages:  []int32{940, 941, 942, 943},
									FnRefs: []string{"582.1", "298481.2485", "3.5281"},
								},
								{
									Text:   page(12840) + fnRef(4, 23),
									Pages:  []int32{12840},
									FnRefs: []string{"4.23"},
								},
							},
						},
					},
				}},
			}},
		},
		{
			name:   "Map footnote",
			volume: 1,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{Level: treemodel.HWork},
				Sections: []treemodel.Section{
					{
						Heading:    treemodel.Heading{Level: treemodel.H1, TextTitle: util.FmtPage(1)},
						Paragraphs: []string{util.FmtPage(5)},
					},
				},
			}},
			footnotes: []treemodel.Footnote{
				{
					Page: 2,
					Nr:   5,
					Text: "This is a simple footnote.",
				},
				{
					Page: 4,
					Nr:   20,
					Text: "This is a " + page(5) + " footnote with a page.",
				},
			},
			model: []model.Work{{
				Sections: []model.Section{{
					Heading: model.Heading{
						Text:  util.FmtHeading(1, util.FmtPage(1)),
						Pages: []int32{1},
					},
					Paragraphs: []model.Paragraph{{
						Text:  util.FmtPage(5),
						Pages: []int32{5},
					}},
				}},
				Footnotes: []model.Footnote{
					{
						Ref:   "2.5",
						Text:  "This is a simple footnote.",
						Pages: []int32{2},
					},
					{
						Ref:   "4.20",
						Text:  "This is a " + page(5) + " footnote with a page.",
						Pages: []int32{4, 5},
					},
				},
			}},
		},
		{
			name:   "Map footnote with non matching page numbers",
			volume: 1,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{Level: treemodel.HWork},
				Sections: []treemodel.Section{
					{
						Heading:    treemodel.Heading{Level: treemodel.H1, TextTitle: util.FmtPage(1)},
						Paragraphs: []string{util.FmtPage(5)},
					},
				},
			}},
			footnotes: []treemodel.Footnote{{
				Page: 43,
				Nr:   348,
				Text: "Summary with non " + page(56) + " matching page numbers",
			}},
			expectError: true,
		},
		{
			name:   "Map summary with non matching page numbers",
			volume: 1,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{Level: treemodel.HWork},
				Sections: []treemodel.Section{
					{
						Heading:    treemodel.Heading{Level: treemodel.H1, TextTitle: util.FmtPage(1)},
						Paragraphs: []string{util.FmtPage(5)},
					},
				},
			}},
			summaries: []treemodel.Summary{{
				Page: 43,
				Line: 348,
				Text: "Summary with non " + page(56) + " matching page numbers",
			}},
			expectError: true,
		},
		{
			name:   "Merge summaries into paragraphs",
			volume: 1,
			sections: []treemodel.Section{
				{
					Heading: treemodel.Heading{Level: treemodel.HWork},
					Sections: []treemodel.Section{{
						Heading: treemodel.Heading{Level: treemodel.H1},
						Paragraphs: []string{
							page(43) + line(348) + "I'm a paragraph.",
							line(58685) + "I'm a paragraph without a page number.",
						},
					}},
				},
				{
					Heading: treemodel.Heading{Level: treemodel.HWork},
					Sections: []treemodel.Section{{
						Heading: treemodel.Heading{Level: treemodel.H1},
						Paragraphs: []string{
							line(5) + "I'm a paragraph with " + page(483) + " a page break inside.",
						},
					}},
				},
			},
			summaries: []treemodel.Summary{
				{Page: 43, Line: 348, Text: "Summary 1"},
				{Page: 43, Line: 58685, Text: "Summary 2"},
				{Page: 482, Line: 5, Text: "Summary 3"},
			},
			model: []model.Work{
				{
					Sections: []model.Section{{
						Heading: model.Heading{Text: util.FmtHeading(1, "")},
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
					Sections: []model.Section{{
						Heading: model.Heading{Text: util.FmtHeading(1, "")},
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
		{
			name:   "Merge summary that starts in the middle of a paragraph",
			volume: 1,
			sections: []treemodel.Section{{
				Heading: treemodel.Heading{Level: treemodel.HWork},
				Sections: []treemodel.Section{
					{
						Heading: treemodel.Heading{Level: treemodel.H1, TextTitle: util.FmtPage(1)},
						Paragraphs: []string{
							page(95) + line(123) + "I'm a sentence." + page(96) + line(31) + "I'm another sentence.",
						},
					},
				},
			}},
			summaries: []treemodel.Summary{{
				Page: 96,
				Line: 31,
				Text: "Summary.",
			}},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapToModel(tc.volume, tc.sections, tc.summaries, tc.footnotes)
			if tc.expectError {
				assert.True(t, err.HasError)
				assert.Nil(t, result)
			} else {
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
	return util.FmtPage(page)
}

func fnRef(page int32, nr int32) string {
	return util.FmtFnRef(page, nr)
}
