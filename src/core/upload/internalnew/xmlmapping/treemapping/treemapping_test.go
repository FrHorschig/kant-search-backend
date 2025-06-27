package treemapping

import (
	"fmt"
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/testutil"
	"github.com/stretchr/testify/assert"
)

const xmlFrame = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<kant_abt1>
  <band>
    <hauptteil>%s</hauptteil>
    <fussnoten>%s</fussnoten>
    <randtexte>%s</randtexte>
  <band>
</kant_abt1>`

func TestTreeMapping(t *testing.T) {
	testCases := []struct {
		name         string
		xmlMain      string
		xmlFootnotes string
		xmlSummaries string
		expWorks     []model.Work
		expFootnotes []model.Footnote
		expSummaries []model.Summary
		expectError  bool
	}{
		{
			name: "multiple headings up to h9",
			xmlMain: `
		        <h1> work 1 </h1>
		        <h2> heading 2 </h2>
		        <h3> heading 3a </h3>
		        <h3> heading 3b </h3>
		        <h4> heading 4 </h4>
		        <h5> heading 5 </h5>
		        <h6> heading 6 </h6>
		        <h7> heading 7a </h7>
		        <h7> heading 7b </h7>
		        <h7> heading 7c </h7>
		        <h7> heading 7d </h7>
		        <h8> heading 8 </h8>
		        <h9> heading 9a </h9>
		        <h9> heading 9b </h9>

		        <h1> work 2 </h1>
		        <h2> heading 2a </h2>
		        <h2> heading 2b </h2>
		        <h3> heading 3a </h3>
		        <h3> heading 3b </h3>
		        <h4> heading 4 </h4>
		        <h2> heading 2c </h2>`,
			expWorks: []model.Work{
				{
					Title:      "<h1> work 1 </h1>",
					Paragraphs: []model.Paragraph{},
					Sections: []model.Section{{
						Heading:    model.Heading{Text: "<h2> heading 2 </h2>"},
						Paragraphs: []model.Paragraph{},
						Sections: []model.Section{
							{
								Heading:    model.Heading{Text: "<h3> heading 3a </h3>"},
								Paragraphs: []model.Paragraph{},
								Sections:   []model.Section{},
							},
							{
								Heading:    model.Heading{Text: "<h3> heading 3b </h3>"},
								Paragraphs: []model.Paragraph{},
								Sections: []model.Section{{
									Heading:    model.Heading{Text: "<h4> heading 4 </h4>"},
									Paragraphs: []model.Paragraph{},
									Sections: []model.Section{{
										Heading:    model.Heading{Text: "<h5> heading 5 </h5>"},
										Paragraphs: []model.Paragraph{},
										Sections: []model.Section{{
											Heading:    model.Heading{Text: "<h6> heading 6 </h6>"},
											Paragraphs: []model.Paragraph{},
											Sections: []model.Section{
												{
													Heading:    model.Heading{Text: "<h7> heading 7a </h7>"},
													Paragraphs: []model.Paragraph{},
													Sections:   []model.Section{},
												},
												{
													Heading:    model.Heading{Text: "<h7> heading 7b </h7>"},
													Paragraphs: []model.Paragraph{},
													Sections:   []model.Section{},
												},
												{
													Heading:    model.Heading{Text: "<h7> heading 7c </h7>"},
													Paragraphs: []model.Paragraph{},
													Sections:   []model.Section{},
												},
												{
													Heading:    model.Heading{Text: "<h7> heading 7d </h7>"},
													Paragraphs: []model.Paragraph{},
													Sections: []model.Section{{
														Heading:    model.Heading{Text: "<h8> heading 8 </h8>"},
														Paragraphs: []model.Paragraph{},
														Sections: []model.Section{
															{
																Heading:    model.Heading{Text: "<h9> heading 9a </h9>"},
																Paragraphs: []model.Paragraph{},
																Sections:   []model.Section{},
															},
															{
																Heading:    model.Heading{Text: "<h9> heading 9b </h9>"},
																Paragraphs: []model.Paragraph{},
																Sections:   []model.Section{},
															},
														},
													}},
												},
											},
										}},
									}},
								}},
							},
						},
					}},
				},
				{
					Title:      "<h1> work 2 </h1>",
					Paragraphs: []model.Paragraph{},
					Sections: []model.Section{
						{
							Heading:    model.Heading{Text: "<h2> heading 2a </h2>"},
							Paragraphs: []model.Paragraph{},
							Sections:   []model.Section{},
						},
						{
							Heading:    model.Heading{Text: "<h2> heading 2b </h2>"},
							Paragraphs: []model.Paragraph{},
							Sections: []model.Section{
								{
									Heading:    model.Heading{Text: "<h3> heading 3a </h3>"},
									Paragraphs: []model.Paragraph{},
									Sections:   []model.Section{},
								},
								{
									Heading:    model.Heading{Text: "<h3> heading 3b </h3>"},
									Paragraphs: []model.Paragraph{},
									Sections: []model.Section{{
										Heading:    model.Heading{Text: "<h4> heading 4 </h4>"},
										Paragraphs: []model.Paragraph{},
										Sections:   []model.Section{},
									}},
								},
							},
						},
						{
							Heading:    model.Heading{Text: "<h2> heading 2c </h2>"},
							Paragraphs: []model.Paragraph{},
							Sections:   []model.Section{},
						},
					},
				},
			},
		},
		{
			name: "page numbers before other elements",
			xmlMain: `
		        <seite nr="34"/>
		        <h1> first </h1>
		        <seite nr="59"/>
		        <h2> 2nd </h2>
		        <seite nr="78"/>
		        <hu> hu paragraph </hu>
		        <seite nr="99"/>
		        <p> three </p>
		        <seite nr="1038"/>
		        <table> three </table>`,
			expWorks: []model.Work{{
				Title: "<h1> first </h1>",
				Sections: []model.Section{{
					Heading: model.Heading{Text: "<h2><seite nr=\"59\"/> 2nd </h2>"},
					Paragraphs: []model.Paragraph{
						{Text: "<hu><seite nr=\"78\"/> hu paragraph </hu>"},
						{Text: "<p><seite nr=\"99\"/> three </p>"},
						{Text: "<table><seite nr=\"1038\"/> three </table>"},
					},
				}},
			}},
		},
		{
			name: "year assignment",
			xmlMain: `
			    <hj> 1234 </hj>
		        <h1> work 1 </h1>
				<h1> work 2 </h1>`,
			expWorks: []model.Work{
				{
					Title:    "<h1> work 1 </h1>",
					Year:     "1234",
					Sections: []model.Section{},
				},
				{
					Title:    "<h1> work 2 </h1>",
					Year:     "1234",
					Sections: []model.Section{},
				},
			},
		},
		{
			name: "op is ignored",
			xmlMain: `
				<op nr="2"/>
                <h1> work 1 </h1>
				<op nr="3"/>
                <h2> heading 2 </h2>
				<op nr="4"/>
				<p> paragraph </p>`,
			expWorks: []model.Work{
				{
					Title: "<h1> work 1 </h1>",
					Sections: []model.Section{{
						Heading: model.Heading{Text: "<h2> heading 2 </h2>"},
						Paragraphs: []model.Paragraph{{
							Text: "<p> paragraph </p>",
						}},
					}},
				},
			},
		},
		{
			name: "map footnotes and summaries",
			xmlFootnotes: `
	            <fn seite="6" nr="7"> footnote 1 </fn>
	            <fn seite="8" nr="9"> footnote 2 </fn>`,
			xmlSummaries: `
	            <randtext seite="2" anfang="3"> randtext 1 </randtext>
	            <randtext seite="4" anfang="5"> randtext 2 </randtext>`,
			expFootnotes: []model.Footnote{
				{Text: "<fn seite=\"6\" nr=\"7\"> footnote 1 </fn>"},
				{Text: "<fn seite=\"8\" nr=\"9\"> footnote 2 </fn>"},
			},
			expSummaries: []model.Summary{
				{Text: "<randtext seite=\"2\" anfang=\"3\"> randtext 1 </randtext>"},
				{Text: "<randtext seite=\"4\" anfang=\"5\"> randtext 2 </randtext>"},
			},
		},
		{
			name: "error on skipped heading level",
			xmlMain: `
		        <h1> work 1 </h1>
		        <h2> heading 2 </h2>
		        <h4> heading 4 </h4>`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xml := fmt.Sprintf(xmlFrame, tc.xmlMain, tc.xmlFootnotes, tc.xmlSummaries)
			works, fns, summs, err := MapTree(xml)
			if tc.expectError {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				testutil.AssertWorks(t, tc.expWorks, works)
				testutil.AssertFootnotes(t, tc.expFootnotes, fns)
				testutil.AssertSummaries(t, tc.expSummaries, summs)
			}
		})
	}
}

// func testErrorInH1(t *testing.T) {
// 	main := `<h1> <unknown/> first </h1>`
// }
// func testErrorInH2(t *testing.T) {
// 	main := `<h1> first </h1> <h2> <unknown> second </h2>`
// }
// func testErrorInHu(t *testing.T) {
// 	main := `<h1> first </h1> <hu> <unknown> hu paragraph </p>`
// }
// func testErrorInP(t *testing.T) {
// 	main := `<h1> first </h1> <p> <unknown> paragraph 1 </p>`
// }
// func testErrorInSeite(t *testing.T) {
// 	main := `<h1> first </h1> <seite/>`
// }
// func testUnknownElement(t *testing.T) {
// 	main := `<h1> first </h1> <my-custom-element> text </my-custom-element>`
// }
// func testMissingFirstH1(t *testing.T) {
// 	main := `<h2> Oh no! </h2>`
// }
// func testErrorInSummary(t *testing.T) {
// 	randtexte := `<randtext> randtext </randtext>`
// }
// func testErrorInFootnote(t *testing.T) {
// 	fussnoten := `<fn> footnote </fn>`
// }
