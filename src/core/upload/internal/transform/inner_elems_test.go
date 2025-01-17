//go:build unit
// +build unit

package transform

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestAntiqua(t *testing.T) {
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
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "Test text <ks-fmt-bold>fettText</ks-fmt-bold>",
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
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text <ks-meta-page>384</ks-meta-page>",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text <ks-meta-line>328</ks-meta-line>",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "Test text",
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
			result, err := antiqua(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBildBildverweis(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		attrs    map[string]string
		expected string
	}{
		{
			name:     "Default values due to missing attributes",
			expected: `{image-extract src="MISSING_IMG_SRC" desc="MISSING_IMG_DESC"}`,
		},
		{
			name:     "Bild attributes are extracted",
			attrs:    map[string]string{"src": "source string", "beschreibung": "description text"},
			expected: `{image-extract src="source string" desc="description text"}`,
		},
		{
			name:     "Text is ignored",
			attrs:    map[string]string{"src": "s", "beschreibung": "d"},
			text:     "some text",
			expected: `{image-extract src="s" desc="d"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, nil)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result := bildBildverweis(el)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestEm1(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Text is extracted",
			text:     "Something",
			expected: "<ks-fmt-emph>Something</ks-fmt-emph>",
		},
		{
			name:     "Space is trimmed",
			text:     " Some text     ",
			expected: "<ks-fmt-emph>Some text</ks-fmt-emph>",
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

func TestEm2(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some emph2 text",
			expected: "<ks-fmt-emph2>Some emph2 text</ks-fmt-emph2>",
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `<ks-fmt-emph2>Test text {image-extract src="source" desc="description text"}</ks-fmt-emph2>`,
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `<ks-fmt-emph2>Test text {image-extract src="source" desc="description text"}</ks-fmt-emph2>`,
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "<ks-fmt-emph2>Test text <ks-fmt-emph>em1Text</ks-fmt-emph></ks-fmt-emph2>",
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: "<ks-fmt-emph2>Test text <ks-fmt-emph2>em2Text</ks-fmt-emph2></ks-fmt-emph2>",
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "<ks-fmt-emph2>Test text <ks-fmt-bold>fettText</ks-fmt-bold></ks-fmt-emph2>",
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: "<ks-fmt-emph2>Test text <ks-fmt-formula>formelText</ks-fmt-formula></ks-fmt-emph2>",
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "<ks-fmt-emph2>Test text <ks-fmt-fnref>1.2</ks-fmt-fnref></ks-fmt-emph2>",
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "<ks-fmt-emph2>Test text <ks-meta-lang>fremdspracheText</ks-meta-lang></ks-fmt-emph2>",
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "<ks-fmt-emph2>Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked></ks-fmt-emph2>",
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "<ks-fmt-emph2>Test text nameText</ks-fmt-emph2>",
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: "<ks-fmt-emph2>Test text II.</ks-fmt-emph2>",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "<ks-fmt-emph2>Test text <ks-meta-page>384</ks-meta-page></ks-fmt-emph2>",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "<ks-fmt-emph2>Test text</ks-fmt-emph2>",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "<ks-fmt-emph2>Test text <ks-meta-line>328</ks-meta-line></ks-fmt-emph2>",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			expected: "<ks-fmt-emph2>Test text</ks-fmt-emph2>",
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
			result, err := em2(el)
			assert.Equal(t, tc.expectError, err.HasError)
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
			name:     "Pure text",
			text:     "Some bold text",
			expected: "<ks-fmt-bold>Some bold text</ks-fmt-bold>",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "<ks-fmt-bold>Test text <ks-meta-page>384</ks-meta-page></ks-fmt-bold>",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "<ks-fmt-bold>Test text <ks-meta-line>328</ks-meta-line></ks-fmt-bold>",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "<ks-fmt-bold>Test text</ks-fmt-bold>",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "<ks-fmt-bold>Test text</ks-fmt-bold>",
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
			assert.Equal(t, tc.expectError, err.HasError)
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
			name:     "Page and number are extracted",
			attrs:    map[string]string{"seite": "27", "nr": "254"},
			expected: "<ks-fmt-fnref>27.254</ks-fmt-fnref>",
		},
		{
			name:     "Default values are used due to missing attributes",
			expected: "<ks-fmt-fnref>MISSING_FNREF_PAGE.MISSING_FNREF_NUMBER</ks-fmt-fnref>",
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"seite": "223845", "nr": "5"},
			expected: "<ks-fmt-fnref>223845.5</ks-fmt-fnref>",
		},
		{
			name:     "Attribute is non-numerical strings",
			attrs:    map[string]string{"seite": "skdhsi", "nr": "sdk"},
			expected: "<ks-fmt-fnref>skdhsi.sdk</ks-fmt-fnref>",
		},
		{
			name:     "Attributes with leading and trailing spaces",
			attrs:    map[string]string{"seite": "  8  ", "nr": " 2     "},
			expected: "<ks-fmt-fnref>8.2</ks-fmt-fnref>",
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
		attrs       map[string]string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some foreing language text",
			expected: "<ks-meta-lang>Some foreing language text</ks-meta-lang>",
		},
		{
			name:     "Text with language attributes",
			text:     "Some foreign language text",
			attrs:    map[string]string{"sprache": "language", "zeichen": "some alphabet", "umschrift": "transcribed text"},
			expected: `<ks-meta-lang lang="language" alphabet="some alphabet" transcript="transcribed text">Some foreign language text</ks-meta-lang>`,
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `<ks-meta-lang>Test text {image-extract src="source" desc="description text"}</ks-meta-lang>`,
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: `<ks-meta-lang>Test text {image-extract src="source" desc="description text"}</ks-meta-lang>`,
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: "<ks-meta-lang>Test text <ks-fmt-emph>em1Text</ks-fmt-emph></ks-meta-lang>",
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: "<ks-meta-lang>Test text <ks-fmt-emph2>em2Text</ks-fmt-emph2></ks-meta-lang>",
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "<ks-meta-lang>Test text <ks-fmt-bold>fettText</ks-fmt-bold></ks-meta-lang>",
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: "<ks-meta-lang>Test text <ks-fmt-formula>formelText</ks-fmt-formula></ks-meta-lang>",
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "<ks-meta-lang>Test text <ks-fmt-fnref>1.2</ks-fmt-fnref></ks-meta-lang>",
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: "<ks-meta-lang>Test text <ks-meta-lang>fremdspracheText</ks-meta-lang></ks-meta-lang>",
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "<ks-meta-lang>Test text <ks-fmt-tracked>gesperrtText</ks-fmt-tracked></ks-meta-lang>",
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "<ks-meta-lang>Test text nameText</ks-meta-lang>",
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: "<ks-meta-lang>Test text II.</ks-meta-lang>",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "<ks-meta-lang>Test text <ks-meta-page>384</ks-meta-page></ks-meta-lang>",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "<ks-meta-lang>Test text</ks-meta-lang>",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "<ks-meta-lang>Test text <ks-meta-line>328</ks-meta-line></ks-meta-lang>",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "<ks-meta-lang>Test text</ks-meta-lang>",
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
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGesperrt(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some tracked text",
			expected: "<ks-fmt-tracked>Some tracked text</ks-fmt-tracked>",
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: "<ks-fmt-tracked>Test text <ks-fmt-bold>fettText</ks-fmt-bold></ks-fmt-tracked>",
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: "<ks-fmt-tracked>Test text nameText</ks-fmt-tracked>",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "<ks-fmt-tracked>Test text <ks-meta-page>384</ks-meta-page></ks-fmt-tracked>",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "<ks-fmt-tracked>Test text</ks-fmt-tracked>",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "<ks-fmt-tracked>Test text <ks-meta-line>328</ks-meta-line></ks-fmt-tracked>",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "<ks-fmt-tracked>Test text</ks-fmt-tracked>",
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
			result, err := gesperrt(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestName(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some name",
			expected: "Some name",
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
			result, err := name(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRomzahl(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		expected    string
		expectError bool
	}{
		{
			name:     "Number is converted with dot",
			text:     "14.",
			expected: "XIV.",
		},
		{
			name:     "Number is converted without dot",
			text:     "116",
			expected: "CXVI",
		},
		{
			name:     "Non-number is ignored",
			text:     "anfk.",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.text, nil)
			result, err := romzahl(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestZeile(t *testing.T) {
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
			expected: "<ks-meta-line>254</ks-meta-line>",
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"nr": "847"},
			expected: "<ks-meta-line>847</ks-meta-line>",
		},
		{
			name:     "Nr attribute with leading and trailing spaces",
			attrs:    map[string]string{"nr": " 2     "},
			expected: "<ks-meta-line>2</ks-meta-line>",
		},
		{
			name:        "Error due to missing number",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Error due to non-numerical number",
			attrs:       map[string]string{"nr": "abf382"},
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", tc.attrs, tc.text, nil)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result, err := zeile(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}
