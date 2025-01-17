package transform

import (
	"testing"

	"github.com/beevik/etree"
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
			name:              "pure text",
			text:              "Some text",
			expectedTextTitle: "Some text",
		},
		{
			name:              "text with fett child element",
			text:              "Test text",
			child:             createElement("fett", nil, "fettText", nil),
			expectedTocTitle:  "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
			expectedTextTitle: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
		},
		{
			name:              "text with fr child element",
			text:              "Test text",
			child:             createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text <ks-fmt-fnref>1.2</ks-fmt-fnref>",
		},
		{
			name:              "text with fremdsprache child element",
			text:              "Test text",
			child:             createElement("fremdsprache", nil, "fremdspracheText", nil),
			expectedTocTitle:  "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
			expectedTextTitle: "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
		},
		{
			name:              "text with gesperrt child element",
			text:              "Test text",
			child:             createElement("gesperrt", nil, "gesperrtText", nil),
			expectedTocTitle:  "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
			expectedTextTitle: "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
		},
		{
			name:              "text with hi child element",
			text:              "Test text",
			child:             createElement("hi", nil, "hiText", nil),
			expectedTocTitle:  "Test text hiText",
			expectedTextTitle: "Test text hiText",
		},
		{
			name:              "text with hu child element",
			text:              "Test text",
			child:             createElement("hi", nil, "huText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text huText",
		},
		{
			name:              "text with name child element",
			text:              "Test text",
			child:             createElement("name", nil, "nameText", nil),
			expectedTocTitle:  "Test text nameText",
			expectedTextTitle: "Test text nameText",
		},
		{
			name:              "text with op child element",
			text:              "Test text",
			child:             createElement("op", nil, "opText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:              "text with romzahl child element",
			text:              "Test text",
			child:             createElement("romzahl", nil, "2.", nil),
			expectedTocTitle:  "Test text II.",
			expectedTextTitle: "Test text II.",
		},
		{
			name:              "text with seite child element",
			text:              "Test text",
			child:             createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expectedTocTitle:  "Test text <ks-meta-page>384</ks-meta-page>",
			expectedTextTitle: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			name:              "text with trenn child element",
			text:              "Test text",
			child:             createElement("trenn", nil, "trennText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:              "text with zeile child element",
			text:              "Test text",
			child:             createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:              "text with leading and trailing spaces",
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
			if tc.expectError {
				assert.NotNil(t, err)
			}
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
			name:     "pure text",
			text:     "Some text",
			expected: "Some text",
		},
		{
			name:     "text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "Test text <ks-fmt-emph>em1Text</ks-fmt-emph>",
		},
		{
			name:     "text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
		},
		{
			name:     "text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text <ks-fmt-fnref>1.2</ks-fmt-fnref>",
		},
		{
			name:     "text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
		},
		{
			name:     "text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
		},
		{
			name:     "text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "Test text nameText",
		},
		{
			name:     "text with op child element",
			text:     "Test text",
			child:    createElement("op", nil, "opText", nil),
			expected: "Test text",
		},
		{
			name:     "text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: "Test text II.",
		},
		{
			name:     "text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			name:     "text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "Test text",
		},
		{
			name:     "text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:     "text with leading and trailing spaces",
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
			if tc.expectError {
				assert.NotNil(t, err)
			}
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
			name:     "pure text",
			text:     "Some paragraph text",
			expected: "Some paragraph text",
		},
		{
			name:     "text with antiqua child element",
			text:     "Test text",
			child:    createElement("antiqua", nil, "antiquaText", nil),
			expected: "Test text",
		},
		{
			name:     "text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `Test text{image-extract src="source" description="description text"}`,
		},
		{
			name:     "text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `Test text {image-extract src="source" description="description text"}"`,
		},
		{
			name:     "text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "Test text <ks-fmt-emph>em1Text</ks-fmt-emph>",
		},
		{
			name:     "text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: "Test text <ks-fmt-tracked>em2Text</ks-fmt-tracked>",
		},
		{
			name:     "text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
		},
		{
			name:     "text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: "Test text <ks-fmt-formula>formelText</ks-fmt-formula>",
		},
		{
			name:     "text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text <ks-fmt-fnref>1.2</ks-fmt-fnref>",
		},
		{
			name:     "text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text <ks-meta-lang>fremdspracheText</ks-meta-lang>",
		},
		{
			name:     "text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked>",
		},
		{
			name:     "text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "Test text nameText",
		},
		{
			name:     "text with op child element",
			text:     "Test text",
			child:    createElement("op", nil, "opText", nil),
			expected: "Test text",
		},
		{
			name:     "text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: "Test text II.",
		},
		{
			name:     "text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			// TODO implement me
			name:     "text with table child element",
			text:     "Test text",
			child:    createElement("table", nil, "tableText", nil),
			expected: "Test text",
		},
		{
			name:     "text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "Test text",
		},
		{
			name:     "text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:     "text with leading and trailing spaces",
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
			result, err := fremdsprache(el)
			if tc.expectError {
				assert.NotNil(t, err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSeite(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		attrs    map[string]string
		expected string
	}{
		{
			name:     "number is extracted",
			attrs:    map[string]string{"nr": "254"},
			expected: "<ks-page>254</ks-page>",
		},
		{
			name:     "default value is used due to missing number",
			expected: "<ks-page>MISSING_PAGE_NUMBER</ks-page>",
		},
		{
			name:     "text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"nr": "847"},
			expected: "<ks-page>847</ks-page>",
		},
		{
			name:     "nr attribute is non-numerical string",
			attrs:    map[string]string{"nr": "kdfghsd"},
			expected: "<ks-page>kdfghsd</ks-page>",
		},
		{
			name:     "nr attribute with leading and trailing spaces",
			attrs:    map[string]string{"nr": " 2     "},
			expected: "<ks-page>2</ks-page>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", tc.attrs, tc.text, nil)
			result := seite(el)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTable(t *testing.T) {
	// TODO implement me
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
