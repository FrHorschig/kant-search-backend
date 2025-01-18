package transform

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestHx(t *testing.T) {
	testCases := []struct {
		name              string
		text              string
		child             *etree.Element
		expectedTocTitle  string
		expectedTextTitle string
		expectError       bool
	}{
		{
			name:              "Pure text",
			text:              "Some text",
			expectedTocTitle:  "Some text",
			expectedTextTitle: "Some text",
		},
		{
			name:              "Text with fett child element",
			text:              "Test text",
			child:             createElement("fett", nil, "fettText", nil),
			expectedTocTitle:  "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
			expectedTextTitle: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
		},
		{
			name:              "Text with fr child element",
			text:              "Test text",
			child:             createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text <ks-fmt-fnref>1.2</ks-fmt-fnref>",
		},
		{
			name:              "Text with fremdsprache child element",
			text:              "Test text",
			child:             createElement("fremdsprache", nil, "fremdspracheText", nil),
			expectedTocTitle:  "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
			expectedTextTitle: "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
		},
		{
			name:              "Text with gesperrt child element",
			text:              "Test text",
			child:             createElement("gesperrt", nil, "gesperrtText", nil),
			expectedTocTitle:  "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
			expectedTextTitle: "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
		},
		{
			name:              "Text with hi child element",
			text:              "Test text",
			child:             createElement("hi", nil, "hiText", nil),
			expectedTocTitle:  "Test text hiText",
			expectedTextTitle: "Test text",
		},
		{
			name:              "Text with hu child element",
			text:              "Test text",
			child:             createElement("hu", nil, "huText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text huText",
		},
		{
			name:              "Text with name child element",
			text:              "Test text",
			child:             createElement("name", nil, "nameText", nil),
			expectedTocTitle:  "Test text nameText",
			expectedTextTitle: "Test text nameText",
		},
		{
			name:              "Text with op child element",
			text:              "Test text",
			child:             createElement("op", nil, "opText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:              "Text with romzahl child element",
			text:              "Test text",
			child:             createElement("romzahl", nil, "2.", nil),
			expectedTocTitle:  "Test text II.",
			expectedTextTitle: "Test text II.",
		},
		{
			name:              "Text with seite child element",
			text:              "Test text",
			child:             createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			name:              "Text with trenn child element",
			text:              "Test text",
			child:             createElement("trenn", nil, "trennText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:              "Text with zeile child element",
			text:              "Test text",
			child:             createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:              "Text with leading and trailing spaces",
			text:              "   Test text       ",
			child:             nil,
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:        "Text with unknown child element",
			child:       createElement("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, tc.child)
			result, err := hx(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expectedTocTitle, result.TocTitle)
			assert.Equal(t, tc.expectedTextTitle, result.TextTitle)
		})
	}
}

func TestHu(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some text",
			expected: "Some text",
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "Test text <ks-fmt-emph>em1Text</ks-fmt-emph>",
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text <ks-fmt-fnref>1.2</ks-fmt-fnref>",
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "Test text nameText",
		},
		{
			name:     "Text with op child element",
			text:     "Test text",
			child:    createElement("op", nil, "opText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: "Test text II.",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "Test text",
		},
		{
			name:        "Text with unknown child element",
			child:       createElement("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, tc.child)
			result, err := hu(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestP(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		attrs       map[string]string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some paragraph text",
			expected: "Some paragraph text",
		},
		{
			name:     "Text with antiqua child element",
			text:     "Test text",
			child:    createElement("antiqua", nil, "antiquaText", nil),
			expected: "Test text antiquaText",
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `Test text {image-extract src="source" desc="description text"}`,
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `Test text {image-extract src="source" desc="description text"}`,
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "Test text <ks-fmt-emph>em1Text</ks-fmt-emph>",
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: "Test text <ks-fmt-emph2>em2Text</ks-fmt-emph2>",
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: "Test text <ks-fmt-formula>formelText</ks-fmt-formula>",
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text <ks-fmt-fnref>1.2</ks-fmt-fnref>",
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "Test text nameText",
		},
		{
			name:     "Text with op child element",
			text:     "Test text",
			child:    createElement("op", nil, "opText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: "Test text II.",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			name:     "Text with table child element",
			text:     "Test text",
			child:    createElement("table", nil, "tableText", nil),
			expected: "Test text {table-extract}",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "Test text",
		},
		{
			name:        "Text with unknown child element",
			child:       createElement("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, tc.child)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result, err := p(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSeite(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		attrs       map[string]string
		expected    string
		expectError bool
	}{
		{
			name:     "Number is extracted",
			attrs:    map[string]string{"nr": "254"},
			expected: "<ks-meta-page>254</ks-meta-page>",
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"nr": "847"},
			expected: "<ks-meta-page>847</ks-meta-page>",
		},
		{
			name:     "Nr attribute with leading zeros",
			attrs:    map[string]string{"nr": "00045"},
			expected: "<ks-meta-page>45</ks-meta-page>",
		},
		{
			name:     "Nr attribute with leading and trailing spaces",
			attrs:    map[string]string{"nr": " 2     "},
			expected: "<ks-meta-page>2</ks-meta-page>",
		},
		{
			name:        "Error due to missing number",
			expectError: true,
		},
		{
			name:        "Error due to non-numeric number attribute",
			attrs:       map[string]string{"nr": "asdh234"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", tc.attrs, tc.text, nil)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result, err := seite(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTable(t *testing.T) {
	result := table()
	assert.Equal(t, "{table-extract}", result)
}

func TestRandtext(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		attrs       map[string]string
		child       *etree.Element
		expected    model.Summary
		expectError bool
	}{
		{
			name:     "Text with randtext attributes",
			text:     "Some text",
			attrs:    map[string]string{"seite": "123", "anfang": "567"},
			expected: model.Summary{Page: 123, Line: 567, Text: "Some text"},
		},
		{
			name:     "Text with p child element",
			text:     "Some text",
			attrs:    map[string]string{"seite": "123", "anfang": "567"},
			child:    createElement("p", nil, "pText", nil),
			expected: model.Summary{Page: 123, Line: 567, Text: "Some text pText"},
		},
		{
			name:        "text with unknown child element",
			text:        "Some text",
			attrs:       map[string]string{"seite": "123", "anfang": "567"},
			child:       createElement("my-custom-tag", nil, "", nil),
			expectError: true,
		},
		{
			name:        "Error due to missing attributes",
			text:        "Some text",
			expectError: true,
		},
		{
			name:        "Error due to non-numerical seite attribute",
			text:        "Some text",
			attrs:       map[string]string{"seite": "s812k", "anfang": "567"},
			expectError: true,
		},
		{
			name:        "Error due to non-numerical anfang attribute",
			text:        "Some text",
			attrs:       map[string]string{"seite": "234", "anfang": "s3j2"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, tc.child)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result, err := summary(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected.Page, result.Page)
			assert.Equal(t, tc.expected.Line, result.Line)
			assert.Equal(t, tc.expected.Text, result.Text)
		})
	}
}

func TestFootnote(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		attrs       map[string]string
		child       *etree.Element
		expected    model.Footnote
		expectError bool
	}{
		{
			name:     "Text with footnote attributes",
			text:     "Some text",
			attrs:    map[string]string{"seite": "123", "nr": "567"},
			expected: model.Footnote{Page: 123, Nr: 567, Text: "Some text"},
		},
		{
			name:     "Text with p child element",
			text:     "Some text",
			attrs:    map[string]string{"seite": "123", "nr": "567"},
			child:    createElement("p", nil, "pText", nil),
			expected: model.Footnote{Page: 123, Nr: 567, Text: "Some text pText"},
		},
		{
			name:        "text with unknown child element",
			text:        "Some text",
			attrs:       map[string]string{"seite": "123", "nr": "567"},
			child:       createElement("my-custom-tag", nil, "", nil),
			expectError: true,
		},
		{
			name:        "Error due to missing attributes",
			text:        "Some text",
			expectError: true,
		},
		{
			name:        "Error due to non-numerical seite attribute",
			text:        "Some text",
			attrs:       map[string]string{"seite": "s812k", "anfang": "567"},
			expectError: true,
		},
		{
			name:        "Error due to non-numerical nr attribute",
			text:        "Some text",
			attrs:       map[string]string{"seite": "234", "nr": "s3j2"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, tc.child)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result, err := footnote(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected.Page, result.Page)
			assert.Equal(t, tc.expected.Nr, result.Nr)
			assert.Equal(t, tc.expected.Text, result.Text)
		})
	}
}

func createElement(tag string, attrs map[string]string, text string, child *etree.Element) *etree.Element {
	el := etree.NewElement(tag)
	for k, v := range attrs {
		el.CreateAttr(k, v)
	}
	if text != "" {
		el.CreateText(text)
	}
	if child != nil {
		el.AddChild(child)
	}
	return el
}
