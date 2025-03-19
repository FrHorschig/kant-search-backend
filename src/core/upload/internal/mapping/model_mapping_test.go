package mapping

import (
	"fmt"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
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
					Level:     model.H1,
					TocTitle:  "work title TOC text",
					TextTitle: "work title",
					Year:      "1724",
				},
				Paragraphs: []string{"paragraph 1 text", "some other paragraph text"},
				Sections: []model.Section{
					{
						Heading: model.Heading{
							Level:     model.H2,
							TocTitle:  "subsection 1 title TOC text",
							TextTitle: "subsection 1 title",
							Year:      "",
						},
						Paragraphs: []string{"subsection 1 paragraph text", "second text"},
					},
					{
						Heading: model.Heading{
							Level:     model.H2,
							TocTitle:  "subsection 2 title TOC text",
							TextTitle: "subsection 2 title",
							Year:      "",
						},
						Paragraphs: []string{"subsection 2 paragraph text", "another second text"},
					},
				},
			}},
			model: []dbmodel.Work{{
				Title:  "work title",
				Year:   util.ToStrPtr("1724"),
				Volume: 1,
				Sections: []dbmodel.Section{{
					Heading: dbmodel.Heading{
						Level:     dbmodel.H1,
						TocTitle:  "work title TOC text",
						TextTitle: "work title",
					},
					Paragraphs: []dbmodel.Paragraph{
						{Text: "paragraph 1 text"},
						{Text: "some other paragraph text"},
					},
					Sections: []dbmodel.Section{
						{
							Heading: dbmodel.Heading{Level: dbmodel.H1},
							Paragraphs: []dbmodel.Paragraph{
								{Text: "subsection 1 paragraph text"},
								{Text: "second text"},
							},
						},
						{
							Heading: dbmodel.Heading{Level: dbmodel.H1},
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
															Sections: []model.Section{
																{
																	Heading: model.Heading{Level: model.H7},
																},
															},
														},
													},
												},
											},
										},
										{Heading: model.Heading{Level: model.H4}},
									},
								},
								{Heading: model.Heading{Level: model.H3}},
								{Heading: model.Heading{Level: model.H3}},
							},
						},
					},
				},
				{Heading: model.Heading{Level: model.H1}},
				{Heading: model.Heading{Level: model.H1}},
			},
			model: []dbmodel.Work{
				{
					Volume: 2,
					Sections: []dbmodel.Section{{
						Heading: dbmodel.Heading{Level: dbmodel.H1},
						Sections: []dbmodel.Section{
							{
								Heading: dbmodel.Heading{Level: dbmodel.H1},
								Sections: []dbmodel.Section{
									{
										Heading: dbmodel.Heading{Level: dbmodel.H2},
										Sections: []dbmodel.Section{
											{
												Heading: dbmodel.Heading{Level: dbmodel.H3},
												Sections: []dbmodel.Section{
													{
														Heading: dbmodel.Heading{Level: dbmodel.H4},
														Sections: []dbmodel.Section{
															{
																Heading: dbmodel.Heading{Level: dbmodel.H5},
																Sections: []dbmodel.Section{
																	{
																		Heading: dbmodel.Heading{Level: dbmodel.H6},
																	},
																},
															},
														},
													},
												},
											},
											{Heading: dbmodel.Heading{Level: dbmodel.H3}},
										},
									},
									{Heading: dbmodel.Heading{Level: dbmodel.H2}},
									{Heading: dbmodel.Heading{Level: dbmodel.H2}},
								},
							},
						},
					}},
				},
				{Sections: []dbmodel.Section{{Heading: dbmodel.Heading{Level: dbmodel.H1}}}},
				{Sections: []dbmodel.Section{{Heading: dbmodel.Heading{Level: dbmodel.H1}}}},
			},
		},
		{
			name:   "Extract pages and footnote references from paragraphs",
			volume: 3,
			sections: []model.Section{{
				Heading: model.Heading{Level: model.H1},
				Paragraphs: []string{
					fnr(2, 64) + page(2),
					"This " + fnr(83, 3) + "is a" + page(254) + " test text.",
					"It " + fnr(582, 1) + " continues " + page(942) + fnr(298481, 2485) + page(942) + " in the " + fnr(3, 5281) + " next" + page(943) + "paragraph.",
					page(23) + fnr(4, 23)},
				Sections: []model.Section{
					{
						Heading: model.Heading{Level: model.H2},
						Paragraphs: []string{
							page(28471) + "This paragraph" + fnr(4, 2) + "starts with a page",
							"This paragraph ends with a page." + fnr(482, 148) + page(3),
						},
					},
				},
			}},
			model: []dbmodel.Work{{
				Volume: 3,
				Sections: []dbmodel.Section{{
					Heading: dbmodel.Heading{Level: dbmodel.H1},
					Paragraphs: []dbmodel.Paragraph{
						{
							Text:         fnr(2, 64) + page(2),
							Pages:        []int32{2},
							FnReferences: []string{"2.64"},
						},
						{
							Text:         "This " + fnr(83, 3) + "is a" + page(254) + " test text.",
							Pages:        []int32{254},
							FnReferences: []string{"83.3"},
						},
						{
							Text:         "It " + fnr(582, 1) + " continues " + page(942) + fnr(298481, 2485) + page(942) + " in the " + fnr(3, 5281) + " next" + page(943) + "paragraph.",
							Pages:        []int32{941, 942, 943},
							FnReferences: []string{"582.1", "298481.2485", "3.5281"},
						},
						{
							Text:         page(23) + fnr(4, 23),
							Pages:        []int32{23},
							FnReferences: []string{"4.23"},
						},
					},
					Sections: []dbmodel.Section{
						{
							Heading: dbmodel.Heading{Level: dbmodel.H1},
							Paragraphs: []dbmodel.Paragraph{
								{
									Text:         page(28471) + "This paragraph" + fnr(4, 3) + " starts with a page.",
									Pages:        []int32{28471},
									FnReferences: []string{"4.3"},
								},
								{
									Text:         "This paragraph ends with a page." + fnr(482, 148) + page(3),
									Pages:        []int32{3},
									FnReferences: []string{"482.148"},
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
				Heading: model.Heading{Level: model.H1},
				Paragraphs: []string{
					page(43) + line(348) + "I'm a paragraph.",
					"I'm a paragraph with the line " + line(58685) + " number at the end",
				},
				Sections: []model.Section{
					{
						Heading: model.Heading{Level: model.H2},
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
				{Page: 96, Line: 1, Text: "Summary 3"},
				{Page: 484, Line: 5, Text: "Summary 4"},
			},
			model: []dbmodel.Work{{
				Volume: 4,
				Sections: []dbmodel.Section{{
					Heading: dbmodel.Heading{Level: dbmodel.H1},
					Paragraphs: []dbmodel.Paragraph{
						{
							Text:  summ("Summary 1") + page(43) + line(348) + "I'm  paragraph.",
							Pages: []int32{43},
						},
						{
							Text:  summ("Summary 2") + page(43) + line(348) + "I'm a paragraph.",
							Pages: []int32{58685},
						},
					},
					Sections: []dbmodel.Section{
						{
							Heading: dbmodel.Heading{Level: dbmodel.H1},
							Paragraphs: []dbmodel.Paragraph{
								{
									Text:  summ("Summary 3") + page(95) + line(123) + "I'm a paragraph.",
									Pages: []int32{95, 96},
								},
								{
									Text:  summ("Summary 4") + "I'm a paragraph with the line " + page(483) + line(2) + " number at the end " + page(484) + line(5) + "that continues over multiple pages.",
									Pages: []int32{483, 484},
								},
							},
						},
					},
				}},
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
	assert.Equal(t, util.ToStrVal(exp.Abbreviation), util.ToStrVal(act.Abbreviation))
	assert.Equal(t, util.ToStrVal(exp.Year), util.ToStrVal(act.Year))
	assert.Equal(t, exp.Volume, act.Volume)

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
	assert.Equal(t, exp.Heading.Level, act.Heading.Level)
	assert.Equal(t, exp.Heading.TocTitle, act.Heading.TocTitle)
	assert.Equal(t, exp.Heading.TextTitle, act.Heading.TextTitle)
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
	assert.ElementsMatch(t, exp.FnReferences, act.FnReferences)
	if len(exp.Sentences) > 0 {
		assert.Equal(t, len(exp.Sentences), len(act.Sentences))
		for i := range exp.Sentences {
			sExp := exp.Sentences[i]
			sAct := act.Sentences[i]
			assert.Equal(t, sExp.Text, sAct.Text)
			assert.ElementsMatch(t, sExp.Pages, sAct.Pages)
		}
	}

}

func assertFootnote(t *testing.T, exp dbmodel.Footnote, act dbmodel.Footnote) {
	assert.Equal(t, exp.Name, act.Name)
	assert.Equal(t, exp.Text, act.Text)
	assert.ElementsMatch(t, exp.Pages, act.Pages)
}

func line(line int32) string {
	return fmt.Sprintf(model.LineFmt, line)
}

func page(page int32) string {
	return fmt.Sprintf(model.PageFmt, page)
}

func fnr(page int32, nr int32) string {
	return fmt.Sprintf(model.FnRefFmt, page, nr)
}

func summ(text string) string {
	return fmt.Sprintf(model.SummaryFmt, text)
}
