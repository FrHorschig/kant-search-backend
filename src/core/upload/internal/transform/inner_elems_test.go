//go:build unit
// +build unit

package transform

import (
	"fmt"
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
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
			expected: "Test text " + fmt.Sprintf(model.BoldFmt, "fettText"),
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: "Test text " + fmt.Sprintf(model.TrackedFmt, "gesperrtText"),
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
			expected: "Test text " + fmt.Sprintf(model.PageFmt, 384) + "",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text " + fmt.Sprintf(model.LineFmt, 328),
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
			expected: `{extract-image src="MISSING_IMG_SRC" desc="MISSING_IMG_DESC"}`,
		},
		{
			name:     "Bild attributes are extracted",
			attrs:    map[string]string{"src": "source string", "beschreibung": "description text"},
			expected: `{extract-image src="source string" desc="description text"}`,
		},
		{
			name:     "Text is ignored",
			attrs:    map[string]string{"src": "s", "beschreibung": "d"},
			text:     "some text",
			expected: `{extract-image src="s" desc="d"}`,
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
			expected: fmt.Sprintf(model.EmphFmt, "Something"),
		},
		{
			name:     "Space is trimmed",
			text:     " Some text     ",
			expected: fmt.Sprintf(model.EmphFmt, "Some text"),
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
			expected: fmt.Sprintf(model.Emph2Fmt, "Some emph2 text"),
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, `Test text {extract-image src="source" desc="description text"}`),
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, `Test text {extract-image src="source" desc="description text"}`),
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.EmphFmt, "em1Text")),
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.Emph2Fmt, "em2Text")),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.BoldFmt, "fettText")),
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.FormulaFmt, "formelText")),
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.FnRefFmt, 1, 2)),
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.LangFmt, "", "fremdspracheText")),
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.TrackedFmt, "gesperrtText")),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text nameText"),
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text II."),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.PageFmt, 384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text "+fmt.Sprintf(model.LineFmt, 328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			expected: fmt.Sprintf(model.Emph2Fmt, "Test text"),
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
			expected: "" + fmt.Sprintf(model.BoldFmt, "Some bold text"),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "" + fmt.Sprintf(model.BoldFmt, "Test text "+fmt.Sprintf(model.PageFmt, 384)),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "" + fmt.Sprintf(model.BoldFmt, "Test text "+fmt.Sprintf(model.LineFmt, 328)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: "" + fmt.Sprintf(model.BoldFmt, "Test text"),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: "" + fmt.Sprintf(model.BoldFmt, "Test text"),
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
			expected: fmt.Sprintf(model.FormulaFmt, "Some formula text"),
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: fmt.Sprintf(model.FormulaFmt, "Test text "+fmt.Sprintf(model.EmphFmt, "em1Text")+""),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: fmt.Sprintf(model.FormulaFmt, "Test text"),
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
			expected: fmt.Sprintf(model.FnRefFmt, 27, 254),
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"seite": "223845", "nr": "5"},
			expected: fmt.Sprintf(model.FnRefFmt, 223845, 5),
		},
		{
			name:     "Attributes with leading and trailing spaces",
			attrs:    map[string]string{"seite": "  8  ", "nr": " 2     "},
			expected: fmt.Sprintf(model.FnRefFmt, 8, 2),
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
			el := createElement("element", tc.attrs, tc.text, nil)
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
			expected: fmt.Sprintf(model.LangFmt, "", ">Some foreing language text"),
		},
		{
			name:     "Text with language attributes",
			text:     "Some foreign language text",
			attrs:    map[string]string{"sprache": "language", "zeichen": "some alphabet", "umschrift": "transcribed text"},
			expected: fmt.Sprintf(model.LangFmt, "", `lang="language" alphabet="some alphabet" transcript="transcribed text">Some foreign language text`),
		},
		{
			name:     "Text with bild child element",
			text:     "Test text",
			child:    createElement("bild", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: fmt.Sprintf(model.LangFmt, "", `Test text {extract-image src="source" desc="description text"}`),
		},
		{
			name:     "Text with bildverweis child element",
			text:     "Test text",
			child:    createElement("bildverweis", map[string]string{"src": "source", "beschreibung": "description text"}, "", nil),
			expected: fmt.Sprintf(model.LangFmt, "", `Test text {extract-image src="source" desc="description text"}`),
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    createElement("em1", nil, "em1Text", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.EmphFmt, "em1Text")),
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    createElement("em2", nil, "em2Text", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.Emph2Fmt, "em2Text")),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.BoldFmt, "fettText")),
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    createElement("formel", nil, "formelText", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.FormulaFmt, "formelText")),
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    createElement("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.FnRefFmt, 1, 2)),
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    createElement("fremdsprache", nil, "fremdspracheText", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.LangFmt, "", "fremdspracheText")),
		},
		{
			name:     "Text with gesperrt child element",
			text:     "Test text",
			child:    createElement("gesperrt", nil, "gesperrtText", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.TrackedFmt, "gesperrtText")),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text nameText"),
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    createElement("romzahl", nil, "2.", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text II."),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.PageFmt, 384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: fmt.Sprintf(model.LangFmt, "", "Test text "+fmt.Sprintf(model.LineFmt, 328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: fmt.Sprintf(model.LangFmt, "", "Test text"),
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
			expected: fmt.Sprintf(model.TrackedFmt, "Some tracked text"),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    createElement("fett", nil, "fettText", nil),
			expected: fmt.Sprintf(model.TrackedFmt, "Test text "+fmt.Sprintf(model.BoldFmt, "fettText")),
		},
		{
			name:     "Text with name child element",
			text:     "Test text",
			child:    createElement("name", nil, "nameText", nil),
			expected: fmt.Sprintf(model.TrackedFmt, "Test text nameText"),
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    createElement("seite", map[string]string{"nr": "384"}, "", nil),
			expected: fmt.Sprintf(model.TrackedFmt, "Test text "+fmt.Sprintf(model.PageFmt, 384)),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    createElement("trenn", nil, "trennText", nil),
			expected: fmt.Sprintf(model.TrackedFmt, "Test text"),
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    createElement("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: fmt.Sprintf(model.TrackedFmt, "Test text "+fmt.Sprintf(model.LineFmt, 328)),
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     "   Test text       ",
			child:    nil,
			expected: fmt.Sprintf(model.TrackedFmt, "Test text"),
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
			expected: "Test text " + fmt.Sprintf(model.PageFmt, 384),
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
			expected: "Test text " + fmt.Sprintf(model.LineFmt, 328),
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
			expected: fmt.Sprintf(model.LineFmt, 254),
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"nr": "847"},
			expected: fmt.Sprintf(model.LineFmt, 847),
		},
		{
			name:     "Nr attribute with leading zeros",
			attrs:    map[string]string{"nr": "00002"},
			expected: fmt.Sprintf(model.LineFmt, 2),
		},
		{
			name:     "Nr attribute with leading and trailing spaces",
			attrs:    map[string]string{"nr": " 2     "},
			expected: fmt.Sprintf(model.LineFmt, 2),
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
