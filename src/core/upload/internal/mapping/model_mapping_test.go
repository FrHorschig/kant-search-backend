package mapping

import (
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
		sections    []model.Section
		summaries   []model.Summary
		footnotes   []model.Footnote
		model       []dbmodel.Work
		expectError bool
	}{
		{
			name: "Multiple sections with simple content are mapped",
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
				Title: "work title",
				Year:  util.ToStrPtr("1724"),
				TextData: dbmodel.Section{
					Heading: dbmodel.Heading{
						Level:     dbmodel.HWork,
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
				},
			}},
		},
		{
			name: "Multiple nested works and sections are mapped",
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
					TextData: dbmodel.Section{
						Heading: dbmodel.Heading{Level: dbmodel.HWork},
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
					},
				},
				{TextData: dbmodel.Section{Heading: dbmodel.Heading{Level: dbmodel.HWork}}},
				{TextData: dbmodel.Section{Heading: dbmodel.Heading{Level: dbmodel.HWork}}},
			},
		},
		{
			name: "Extract pages from paragraphs",
			// TODO
		},
		{
			name: "Extract footnote references from paragraphs",
			// TODO
		},
		{
			name: "Merge summaries into paragraphs",
			// TODO
		},
		{
			name: "Map footnotes",
			// TODO
		},
		{
			name: "Extract pages from footnotes",
			// TODO
		},
		// {
		// 	name: "Test name",
		// 	sections: []model.Section{{
		// 		Heading:    model.Heading{},
		// 		Paragraphs: []string{},
		// 		Sections:   []model.Section{},
		// 	}},
		// 	summaries: []model.Summary{},
		// 	footnotes: []model.Footnote{},
		// 	model: []dbmodel.Work{{
		// 		TextData: dbmodel.Section{
		// 			Heading:    dbmodel.Heading{},
		// 			Paragraphs: []dbmodel.Paragraph{},
		// 			Sentences:  []dbmodel.Sentence{},
		// 			Sections:   []dbmodel.Section{},
		// 		},
		// 		Footnotes: []dbmodel.Footnote{},
		// 	}},
		// },
	}

	sut := &modelMapperImpl{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := sut.Map(tc.sections, tc.summaries, tc.footnotes)
			if tc.expectError {
				assert.True(t, err.HasError)
				assert.Nil(t, result)
			}
			assertWorks(t, tc.model, result)
		})
	}
}

func assertWorks(t *testing.T, expected []dbmodel.Work, actual []dbmodel.Work) {
	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		exp := expected[i]
		act := actual[i]
		assert.NotNil(t, exp)
		assert.NotNil(t, act)

		assert.Equal(t, exp.Code, act.Code)
		assert.Equal(t, util.ToStrVal(exp.Abbreviation), util.ToStrVal(act.Abbreviation))
		assert.Equal(t, util.ToStrVal(exp.Year), util.ToStrVal(act.Year))
		assert.Equal(t, exp.Volume, act.Volume)

		assertSections(t, exp.TextData, act.TextData)
		assertFootnotes(t, exp.Footnotes, act.Footnotes)
	}
}

func assertSections(t *testing.T, exp dbmodel.Section, act dbmodel.Section) {
	assert.Equal(t, exp.Heading.Level, act.Heading.Level)
	assert.Equal(t, exp.Heading.TocTitle, act.Heading.TocTitle)
	assert.Equal(t, exp.Heading.TextTitle, act.Heading.TextTitle)

	assert.Equal(t, len(exp.Paragraphs), len(act.Paragraphs))
	for i := range exp.Paragraphs {
		pExp := exp.Paragraphs[i]
		pAct := act.Paragraphs[i]
		assert.Equal(t, pExp.Text, pAct.Text)
		assert.ElementsMatch(t, pExp.Pages, pAct.Pages)
		assert.ElementsMatch(t, pExp.FnReferences, pAct.FnReferences)

		assert.Equal(t, len(pExp.Sentences), len(pAct.Sentences))
		for i := range pExp.Sentences {
			sExp := pExp.Sentences[i]
			sAct := pAct.Sentences[i]
			assert.Equal(t, sExp.Text, sAct.Text)
			assert.ElementsMatch(t, sExp.Pages, sAct.Pages)
		}

	}

	assert.Equal(t, len(exp.Sections), len(act.Sections))
	for i := range exp.Sections {
		assertSections(t, exp.Sections[i], act.Sections[i])
	}
}

func assertFootnotes(t *testing.T, exp []dbmodel.Footnote, act []dbmodel.Footnote) {
	assert.Equal(t, len(exp), len(act))
	for i := range exp {
		assert.Equal(t, exp[i].Name, act[i].Name)
		assert.Equal(t, exp[i].Text, act[i].Text)
		assert.ElementsMatch(t, exp[i].Pages, act[i].Pages)
	}
}
