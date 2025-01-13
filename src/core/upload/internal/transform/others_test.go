//go:build unit
// +build unit

package transform

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestEm1(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "text is extracted",
			text:     "Something",
			expected: "<em>Something</em>",
		},
		{
			name:     "space is trimmed",
			text:     " Some text     ",
			expected: "<em>Some text</em>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, nil)
			result := em1(el)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFett(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "pure text",
			text:     "Some bold text",
			expected: "<b>Some bold text</b>",
		},
		{
			name:     "text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "<b>Test text <ks-page>384</ks-page></b>",
		},
		{
			name:     "text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "<b>Test text <ks-line>328</ks-line></b>",
		},
		{
			name:     "text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "", nil),
			expected: "<b>Test text</b>",
		},
		{
			name:     "text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    createElement("zeile", map[string]string{"nr": "842"}, "", nil),
			expected: "<b>Test text <ks-line>842</ks-line></b>",
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
			result, err := fett(el)
			if tc.expectError {
				assert.NotNil(t, err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFr(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		attrs    map[string]string
		expected string
	}{
		{
			name:     "page and number are extracted",
			attrs:    map[string]string{"seite": "27", "nr": "254"},
			expected: "<ks-fn-ref>27.254</ks-fn-ref>",
		},
		{
			name:     "default values are used due to missing attributes",
			expected: "<ks-fn-ref>MISSING_FR_PAGE.MISSING_FR_NUMBER</ks-fn-ref>",
		},
		{
			name:     "text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"seite": "223845", "nr": "5"},
			expected: "<ks-fn-ref>223845.5</ks-fn-ref>",
		},
		{
			name:     "attribute is non-numerical strings",
			attrs:    map[string]string{"seite": "skdhsi", "nr": "sdk"},
			expected: "<ks-fn-ref>skdhsi.sdk</ks-fn-ref>",
		},
		{
			name:     "attributes with leading and trailing spaces",
			attrs:    map[string]string{"seite": "  8  ", "nr": " 2     "},
			expected: "<ks-fn-ref>8.2</ks-fn-ref>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", tc.attrs, tc.text, nil)
			result := fr(el)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFremdsprache(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "pure text",
			text:     "Some bold text",
			expected: "Some bold text",
		},
		{
			name:     "text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "abc"}, "", nil),
			expected: `Test text <ks-img>abc</ks-img>`,
		},
		{
			name:     "text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "abc"}, "", nil),
			expected: `Test text <ks-img-ref>abc</ks-img-ref>`,
		},
		{
			name:     "text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "Test text <em>em1Text</em>",
		},
		{
			name:     "text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: "Test text <ks-tracked>em2Text</ks-tracked>",
		},
		{
			name:     "text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "Test text <b>fettText</b>",
		},
		{
			name:     "text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: "Test text <ks-formula>formelText</ks-formula>",
		},
		{
			name:     "text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text <ks-fn-ref>1.2</ks-fn-ref>",
		},
		{
			name:     "text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text fremdspracheText",
		},
		{
			name:     "text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text <ks-tracked>gesperrtText</ks-tracked>",
		},
		{
			name:     "text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "Test text nameText",
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
			expected: "Test text <ks-page>384</ks-page>",
		},
		{
			name:     "text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "", nil),
			expected: "Test text",
		},
		{
			name:     "text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text <ks-line>328</ks-line>",
		},
		{
			name:     "text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    createElement("zeile", map[string]string{"nr": "842"}, "", nil),
			expected: "Test text <ks-line>842</ks-line>",
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
			result, err := fremdsprache(el)
			if tc.expectError {
				assert.NotNil(t, err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}
