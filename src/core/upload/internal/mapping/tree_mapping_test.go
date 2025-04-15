package mapping

import (
	"fmt"
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/treemodel"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestTreeMapping(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"Test nested headings":                 testNestedHeadings,
		"Test multiple equal leveled headings": testMultipleEqualHeadings,
		"Test page before heading":             testPageBeforeHeading,
		"Test pure hu heading":                 testPureHuHeading,
		"Test year assignment":                 testYearAssignment,
		"Test paragraph extraction":            testParagraphExtraction,
		"Test op is ignored":                   testOpIsIgnored,
		"Test extraction of allparts":          testMainSummaryFootnoteExtraction,
		"Test error in h1":                     testErrorInH1,
		"Test error in h2":                     testErrorInH2,
		"Test error in hu":                     testErrorInHu,
		"Test error in p":                      testErrorInP,
		"Test error in seite":                  testErrorInSeite,
		"Test unknown element":                 testUnknownElement,
		"Test missing first H1":                testMissingFirstH1,
		"Test error in summary":                testErrorInSummary,
		"Test error in footnote":               testErrorInFootnote,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testNestedHeadings(t *testing.T) {
	main := `
    <h1> heading 1 </h1>
    <h2> heading 2 </h2>
    <h3> heading 3 </h3>
    <h4> heading 4 </h4>
    <h5> heading 5 </h5>
    <h6> heading 6 </h6>
    <h7> heading 7 </h7>
    <h8> heading 8 </h8>
    <h9> heading 9 </h9>

    <h1> heading 1 </h1>
    <h2> heading 2 </h2>
    <h5> heading 5 </h5>
    <h3> heading 3 </h3>
    <h8> heading 8 </h8>
    <h2> heading 2 </h2>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 2, len(works))

	assert.Equal(t, treemodel.HWork, works[0].Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections))
	assert.Equal(t, treemodel.H1, works[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H2, works[0].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H3, works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H4, works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H5, works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H6, works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H7, works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, 1, len(works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H8, works[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Sections[0].
		Heading.Level)

	// 1
	// --2
	//   --5
	//   --3
	//     --8
	// --2
	assert.Equal(t, treemodel.HWork, works[1].Heading.Level)
	assert.Equal(t, 2, len(works[1].
		Sections))
	assert.Equal(t, treemodel.H1, works[1].
		Sections[0].
		Heading.Level)
	assert.Equal(t, treemodel.H1, works[1].
		Sections[1].
		Heading.Level)
	assert.Equal(t, 2, len(works[1].
		Sections[0].
		Sections))
	assert.Equal(t, treemodel.H2, works[1].
		Sections[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, treemodel.H2, works[1].
		Sections[0].
		Sections[1].
		Heading.Level)
	assert.Equal(t, 1, len(works[1].
		Sections[0].
		Sections[1].
		Sections))
	assert.Equal(t, treemodel.H3, works[1].
		Sections[0].
		Sections[1].
		Sections[0].
		Heading.Level)
}

func testMultipleEqualHeadings(t *testing.T) {
	main := `
    <h1> heading 1 </h1>
    <h2> heading 2 </h2>
    <h2> heading 2 </h2>
    <h3> heading 3 </h3>
    <h3> heading 3 </h3>
    <h2> heading 2 </h2>
    <h3> heading 3 </h3>
    <h3> heading 3 </h3>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 1, len(works))

	assert.Equal(t, treemodel.HWork, works[0].Heading.Level)
	assert.Equal(t, 3, len(works[0].
		Sections))
	assert.Equal(t, treemodel.H1, works[0].
		Sections[0].
		Heading.Level)
	assert.Equal(t, treemodel.H1, works[0].
		Sections[1].
		Heading.Level)
	assert.Equal(t, treemodel.H1, works[0].
		Sections[2].
		Heading.Level)
	assert.Equal(t, 2, len(works[0].
		Sections[1].
		Sections))
	assert.Equal(t, treemodel.H2, works[0].
		Sections[1].
		Sections[0].
		Heading.Level)
	assert.Equal(t, treemodel.H2, works[0].
		Sections[1].
		Sections[1].
		Heading.Level)
	assert.Equal(t, 2, len(works[0].
		Sections[2].
		Sections))
	assert.Equal(t, treemodel.H2, works[0].
		Sections[2].
		Sections[0].
		Heading.Level)
	assert.Equal(t, treemodel.H2, works[0].
		Sections[2].
		Sections[1].
		Heading.Level)
}

func testPageBeforeHeading(t *testing.T) {
	main := `
    <seite nr="34"/>
    <h1> first </h1>
    <seite nr="59"/>
    <h2> one </h2>
    <h2> two </h2>
    <seite nr="78"/>
    <hu> hu paragraph </hu>
    <seite nr="99"/>
    <h2> three </h2>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 1, len(works))
	assert.Equal(t, 3, len(works[0].Sections))
	assert.Equal(t, page(34)+" first", works[0].Heading.TextTitle)
	assert.Equal(t, util.FmtHeading(1, page(59)+" one"), works[0].Sections[0].Heading.TextTitle)
	assert.Equal(t, util.FmtHeading(1, "two"), works[0].Sections[1].Heading.TextTitle)
	assert.Equal(t, 1, len(works[0].Sections[1].Paragraphs))
	assert.Equal(t, util.FmtParHeading(page(78)+" hu paragraph"), works[0].Sections[1].Paragraphs[0])
	assert.Equal(t, util.FmtHeading(1, page(99)+" three"), works[0].Sections[2].Heading.TextTitle)
}

func testPureHuHeading(t *testing.T) {
	main := `
    <h1> first </h1>
    <h2> <hu> hu paragraph </hu> </h2>
    <h3> h3 text </h3>
    `
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 1, len(works))
	assert.Equal(t, 1, len(works[0].Paragraphs))
	assert.Equal(t, util.FmtParHeading("hu paragraph"), works[0].Paragraphs[0])

	assert.Equal(t, 1, len(works[0].Sections))
	assert.Equal(t, treemodel.H1, works[0].Sections[0].Heading.Level)
	assert.Equal(t, util.FmtHeading(1, "h3 text"), works[0].Sections[0].Heading.TextTitle)
}

func testYearAssignment(t *testing.T) {
	main := `
    <h1> first </h1>
	<hj> 1724 </hj>
    <h1> second </h1>
	<hj> 1804 </hj>
    <h2> two </h2>
    <h1> third </h1>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 3, len(works))
	assert.Equal(t, "", works[0].Heading.Year)
	assert.Equal(t, "1724", works[1].Heading.Year)
	assert.Equal(t, "1804", works[2].Heading.Year)
	assert.Equal(t, 1, len(works[1].Sections))
	assert.Equal(t, "", works[1].Sections[0].Heading.Year)
}

func testParagraphExtraction(t *testing.T) {
	main := `
    <h1> first </h1>
    <p> paragraph 1.1 </p>
    <p> paragraph 1.2 </p>
    <h2> second </h2>
    <p> paragraph 2.1 </p>
    <p> paragraph 2.2 </p>
    <h7> second </h7>
    <p> paragraph 7.1 </p>
    <h2> second </h2>
    <p> paragraph 22.1 </p>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 1, len(works))
	assert.Equal(t, 2, len(works[0].Paragraphs))
	assert.Equal(t, "paragraph 1.1", works[0].Paragraphs[0])
	assert.Equal(t, "paragraph 1.2", works[0].Paragraphs[1])

	assert.Equal(t, 2, len(works[0].Sections))
	assert.Equal(t, 2, len(works[0].Sections[0].Paragraphs))
	assert.Equal(t, "paragraph 2.1", works[0].Sections[0].Paragraphs[0])
	assert.Equal(t, "paragraph 2.2", works[0].Sections[0].Paragraphs[1])
	assert.Equal(t, 1, len(works[0].Sections[0].Sections))
	assert.Equal(t, 1, len(works[0].Sections[0].Sections[0].Paragraphs))
	assert.Equal(t, "paragraph 7.1", works[0].Sections[0].Sections[0].Paragraphs[0])
	assert.Equal(t, 1, len(works[0].Sections[1].Paragraphs))
	assert.Equal(t, "paragraph 22.1", works[0].Sections[1].Paragraphs[0])
}

func testOpIsIgnored(t *testing.T) {
	main := `
    <op> op text </op>
    <h1> first </h1>
    <op> op text </op>
    <p> paragraph 1.1 </p>
    <op> op text </op>
    <p> paragraph 1.2 </p>
    <op> op text </op>
    <h2> second </h2>
    <op> op text </op>
    <p> paragraph 2.1 </p>
    <op> op text </op>
    <p> paragraph 2.2 </p>
    <op> op text </op>
    <h7> second </h7>
    <op> op text </op>
    <p> paragraph 7.1 </p>
    <op> op text </op>
    <h2> second </h2>
    <op> op text </op>
    <p> paragraph 22.1 </p>
    <op> op text </op>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 1, len(works))
	assert.Equal(t, 2, len(works[0].Paragraphs))
	assert.Equal(t, "paragraph 1.1", works[0].Paragraphs[0])
	assert.Equal(t, "paragraph 1.2", works[0].Paragraphs[1])

	assert.Equal(t, 2, len(works[0].Sections))
	assert.Equal(t, 2, len(works[0].Sections[0].Paragraphs))
	assert.Equal(t, "paragraph 2.1", works[0].Sections[0].Paragraphs[0])
	assert.Equal(t, "paragraph 2.2", works[0].Sections[0].Paragraphs[1])
	assert.Equal(t, 1, len(works[0].Sections[0].Sections))
	assert.Equal(t, 1, len(works[0].Sections[0].Sections[0].Paragraphs))
	assert.Equal(t, "paragraph 7.1", works[0].Sections[0].Sections[0].Paragraphs[0])
	assert.Equal(t, 1, len(works[0].Sections[1].Paragraphs))
	assert.Equal(t, "paragraph 22.1", works[0].Sections[1].Paragraphs[0])
}

func testMainSummaryFootnoteExtraction(t *testing.T) {
	main := `
    <h1> first </h1>
    <p> paragraph 1.1 </p>`
	randtexte := `
	<randtext seite="2" anfang="3"> randtext 1 </randtext>
	<randtext seite="4" anfang="5"> randtext 2 </randtext>`
	fussnoten := `
	<fn seite="6" nr="7"> footnote 1 </fn>
	<fn seite="8" nr="9"> footnote 2 </fn>`
	doc := createNewDocument(main, randtexte, fussnoten)

	// WHEN
	works, summaries, footnotes, err := MapToTree(doc)

	// THEN
	assert.False(t, err.HasError)
	assert.Equal(t, 1, len(works))
	assert.Equal(t, 1, len(works[0].Paragraphs))
	assert.Equal(t, "paragraph 1.1", works[0].Paragraphs[0])

	assert.Equal(t, 2, len(summaries))
	assert.Equal(t, "randtext 1", summaries[0].Text)
	assert.Equal(t, "randtext 2", summaries[1].Text)
	assert.Equal(t, 2, len(footnotes))
	assert.Equal(t, "footnote 1", footnotes[0].Text)
	assert.Equal(t, "footnote 2", footnotes[1].Text)
}

func testErrorInH1(t *testing.T) {
	main := `<h1> <unknown/> first </h1>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testErrorInH2(t *testing.T) {
	main := `<h1> first </h1> <h2> <unknown> second </h2>
	`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testErrorInHu(t *testing.T) {
	main := `<h1> first </h1> <hu> <unknown> hu paragraph </p>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testErrorInP(t *testing.T) {
	main := `<h1> first </h1> <p> <unknown> paragraph 1 </p>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testErrorInSeite(t *testing.T) {
	main := `<h1> first </h1> <seite/>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testUnknownElement(t *testing.T) {
	main := `<h1> first </h1> <my-custom-element> text </my-custom-element>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testMissingFirstH1(t *testing.T) {
	main := `<h2> Oh no! </h2>`
	doc := createNewDocument(main, "", "")

	// WHEN
	works, _, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, works)
}

func testErrorInSummary(t *testing.T) {
	randtexte := `<randtext> randtext </randtext>`
	doc := createNewDocument("", randtexte, "")

	// WHEN
	_, summaries, _, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, summaries)
}

func testErrorInFootnote(t *testing.T) {
	fussnoten := `<fn> footnote </fn>`
	doc := createNewDocument("", "", fussnoten)

	// WHEN
	_, _, footnotes, err := MapToTree(doc)

	// THEN
	assert.True(t, err.HasError)
	assert.Nil(t, footnotes)
}

const xmlFrame = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<kant_abt1>
  <band>
    <hauptteil>%s</hauptteil>
    <randtexte>%s</randtexte>
    <fussnoten>%s</fussnoten>
  <band>
</kant_abt1>`

func createNewDocument(main string, summaries string, footnotes string) *etree.Document {
	xml := fmt.Sprintf(xmlFrame, main, summaries, footnotes)
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	return doc
}
