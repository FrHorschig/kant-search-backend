//go:build unit
// +build unit

package texttransform

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/util"
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
			child:    elem("fett", nil, "fettText", nil),
			expected: "Test text " + util.FmtBold("fettText"),
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    elem("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text " + util.FmtTracked("gesperrtText"),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    elem("name", nil, "nameText", nil),
			expected: "Test text " + util.FmtName("nameText"),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text " + util.FmtPage(384) + "",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text " + util.FmtLine(328),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
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
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
			result, err := antiqua(el)
			assert.Equal(t, tc.expectError, err.HasError)
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
			expected: util.FmtEmph("Something"),
		},
		{
			name:     "Space is trimmed",
			text:     " Some text     ",
			expected: util.FmtEmph("Some text"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, nil)
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
			expected: util.FmtEmph2("Some emph2 text"),
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    elem("bild", nil, "", nil),
			expected: util.FmtEmph2("Test text" + util.FmtImg("src", "descr")),
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    elem("bildverweis", nil, "", nil),
			expected: util.FmtEmph2("Test text" + util.FmtImgRef("src", "descr")),
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    elem("em1", nil, "em1Text", nil),
			expected: util.FmtEmph2("Test text " + util.FmtEmph("em1Text")),
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    elem("em2", nil, "em2Text", nil),
			expected: util.FmtEmph2("Test text " + util.FmtEmph2("em2Text")),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    elem("fett", nil, "fettText", nil),
			expected: util.FmtEmph2("Test text " + util.FmtBold("fettText")),
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    elem("formel", nil, "formelText", nil),
			expected: util.FmtEmph2("Test text " + util.FmtFormula("formelText")),
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    elem("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: util.FmtEmph2("Test text " + util.FmtFnRef(1, 2)),
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    elem("fremdsprache", nil, "fremdspracheText", nil),
			expected: util.FmtEmph2("Test text " + util.FmtLang("fremdspracheText")),
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    elem("gesperrt", nil, "gesperrtText", nil),
			expected: util.FmtEmph2("Test text " + util.FmtTracked("gesperrtText")),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    elem("name", nil, "nameText", nil),
			expected: util.FmtEmph2("Test text " + util.FmtName("nameText")),
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    elem("romzahl", nil, "2.", nil),
			expected: util.FmtEmph2("Test text II."),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: util.FmtEmph2("Test text " + util.FmtPage(384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: util.FmtEmph2("Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: util.FmtEmph2("Test text " + util.FmtLine(328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			expected: util.FmtEmph2("Test text"),
		},
		{
			name:        "Text with unknown child element",
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
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
			expected: "" + util.FmtBold("Some bold text"),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "" + util.FmtBold("Test text "+util.FmtPage(384)),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "" + util.FmtBold("Test text "+util.FmtLine(328)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: "" + util.FmtBold("Test text"),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "" + util.FmtBold("Test text"),
		},
		{
			name:        "Text with unknown child element",
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
			result, err := fett(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFormel(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		child       *etree.Element
		expected    string
		expectError bool
	}{
		{
			name:     "Pure text",
			text:     "Some formula text",
			expected: util.FmtFormula("Some formula text"),
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    elem("em1", nil, "em1Text", nil),
			expected: util.FmtFormula("Test text " + util.FmtEmph("em1Text") + ""),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: util.FmtFormula("Test text"),
		},
		{
			name:        "Text with unknown child element",
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
			result, err := formel(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFr(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		attrs       map[string]string
		expected    string
		expectError bool
	}{
		{
			name:     "Page and number are extracted",
			attrs:    map[string]string{"seite": "27", "nr": "254"},
			expected: util.FmtFnRef(27, 254),
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"seite": "223845", "nr": "5"},
			expected: util.FmtFnRef(223845, 5),
		},
		{
			name:     "Attributes with leading and trailing spaces",
			attrs:    map[string]string{"seite": "  8  ", "nr": " 2     "},
			expected: util.FmtFnRef(8, 2),
		},
		{
			name:        "Error due to missing attributes",
			expectError: true,
		},
		{
			name:        "Error due to non-numeric seite attribute",
			attrs:       map[string]string{"seite": "asdf", "nr": "254"},
			expectError: true,
		},
		{
			name:        "Error due to non-numeric nr attribute",
			attrs:       map[string]string{"seite": "23", "nr": "asdfjk"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", tc.attrs, tc.text, nil)
			result, err := fr(el)
			assert.Equal(t, tc.expectError, err.HasError)
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
			expected: util.FmtLang("Some foreing language text"),
		},
		{
			name:     "Text with language attributes",
			text:     "Some foreign language text",
			attrs:    map[string]string{"sprache": "language", "zeichen": "some alphabet", "umschrift": "transcribed text"},
			expected: util.FmtLang("Some foreign language text"),
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    elem("bild", nil, "", nil),
			expected: "Test text" + util.FmtImg("src", "descr"),
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    elem("bildverweis", nil, "", nil),
			expected: "Test text" + util.FmtImgRef("src", "descr"),
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    elem("em1", nil, "em1Text", nil),
			expected: util.FmtLang("Test text " + util.FmtEmph("em1Text")),
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    elem("em2", nil, "em2Text", nil),
			expected: util.FmtLang("Test text " + util.FmtEmph2("em2Text")),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    elem("fett", nil, "fettText", nil),
			expected: util.FmtLang("Test text " + util.FmtBold("fettText")),
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    elem("formel", nil, "formelText", nil),
			expected: util.FmtLang("Test text " + util.FmtFormula("formelText")),
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    elem("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: util.FmtLang("Test text " + util.FmtFnRef(1, 2)),
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    elem("fremdsprache", nil, "fremdspracheText", nil),
			expected: util.FmtLang("Test text " + util.FmtLang("fremdspracheText")),
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    elem("gesperrt", nil, "gesperrtText", nil),
			expected: util.FmtLang("Test text " + util.FmtTracked("gesperrtText")),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    elem("name", nil, "nameText", nil),
			expected: util.FmtLang("Test text " + util.FmtName("nameText")),
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    elem("romzahl", nil, "2.", nil),
			expected: util.FmtLang("Test text II."),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: util.FmtLang("Test text " + util.FmtPage(384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: util.FmtLang("Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: util.FmtLang("Test text " + util.FmtLine(328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: util.FmtLang("Test text"),
		},
		{
			name:        "Text with unknown child element",
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
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
			expected: util.FmtTracked("Some tracked text"),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    elem("fett", nil, "fettText", nil),
			expected: util.FmtTracked("Test text " + util.FmtBold("fettText")),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    elem("name", nil, "nameText", nil),
			expected: util.FmtTracked("Test text " + util.FmtName("nameText")),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: util.FmtTracked("Test text " + util.FmtPage(384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: util.FmtTracked("Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: util.FmtTracked("Test text " + util.FmtLine(328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: util.FmtTracked("Test text"),
		},
		{
			name:        "Text with unknown child element",
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
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
			expected: util.FmtName("Some name"),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: util.FmtName("Test text " + util.FmtPage(384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: util.FmtName("Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: util.FmtName("Test text " + util.FmtLine(328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: util.FmtName("Test text"),
		},
		{
			name:        "Text with unknown child element",
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
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
			el := elem("element", nil, tc.text, nil)
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
			expected: util.FmtLine(254),
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"nr": "847"},
			expected: util.FmtLine(847),
		},
		{
			name:     "Nr attribute with leading zeros",
			attrs:    map[string]string{"nr": "00002"},
			expected: util.FmtLine(2),
		},
		{
			name:     "Nr attribute with leading and trailing spaces",
			attrs:    map[string]string{"nr": " 2     "},
			expected: util.FmtLine(2),
		},
		{
			name:        "Error due to missing number",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Error due to non-numerical attribute",
			attrs:       map[string]string{"nr": "abf382"},
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", tc.attrs, tc.text, nil)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			result, err := zeile(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expected, result)
		})
	}
}
