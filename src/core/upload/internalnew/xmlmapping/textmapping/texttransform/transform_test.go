package texttransform

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/util"
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
			child:             elem("fett", nil, "fettText", nil),
			expectedTocTitle:  "Test text " + "fettText",
			expectedTextTitle: "Test text " + util.FmtBold("fettText"),
		},
		{
			name:              "Text with fr child element",
			text:              "Test text",
			child:             elem("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text " + util.FmtFnRef(1, 2),
		},
		{
			name:              "Text with fremdsprache child element",
			text:              "Test text",
			child:             elem("fremdsprache", nil, "fremdspracheText", nil),
			expectedTocTitle:  "Test text " + util.FmtLang("fremdspracheText"),
			expectedTextTitle: "Test text " + util.FmtLang("fremdspracheText"),
		},
		{
			name:              "Text with gesperrt child element",
			text:              "Test text",
			child:             elem("gesperrt", nil, "gesperrtText", nil),
			expectedTocTitle:  "Test text " + "gesperrtText",
			expectedTextTitle: "Test text " + util.FmtTracked("gesperrtText"),
		},
		{
			name:              "Text with hi child element",
			text:              "Test text",
			child:             elem("hi", nil, "hiText", nil),
			expectedTocTitle:  "Test text hiText",
			expectedTextTitle: "Test text",
		},
		{
			name:              "Text with hu child element",
			text:              "Test text",
			child:             elem("hu", nil, "huText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text huText",
		},
		{
			name:              "Text with name child element",
			text:              "Test text",
			child:             elem("name", nil, "nameText", nil),
			expectedTocTitle:  "Test text nameText",
			expectedTextTitle: "Test text " + util.FmtName("nameText"),
		},
		{
			name:              "Text with op child element",
			text:              "Test text",
			child:             elem("op", nil, "opText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:              "Text with romzahl child element",
			text:              "Test text",
			child:             elem("romzahl", nil, "2.", nil),
			expectedTocTitle:  "Test text II",
			expectedTextTitle: "Test text II.",
		},
		{
			name:              "Text with seite child element",
			text:              "Test text",
			child:             elem("seite", map[string]string{"nr": "384"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text " + util.FmtPage(384) + "",
		},
		{
			name:              "Text with trenn child element",
			text:              "Test text",
			child:             elem("trenn", nil, "trennText", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text",
		},
		{
			name:              "Text with zeile child element",
			text:              "Test text",
			child:             elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expectedTocTitle:  "Test text",
			expectedTextTitle: "Test text " + util.FmtLine(328),
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
			child:       elem("my-custom-tag", nil, "", nil),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := elem("element", nil, tc.text, tc.child)
			fmtResult, tocResult, err := hx(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expectedTocTitle, tocResult)
			assert.Equal(t, tc.expectedTextTitle, fmtResult)
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
			child:    elem("em1", nil, "em1Text", nil),
			expected: "Test text " + util.FmtEmph("em1Text"),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    elem("fett", nil, "fettText", nil),
			expected: "Test text " + util.FmtBold("fettText"),
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    elem("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text " + util.FmtFnRef(1, 2),
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    elem("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text " + util.FmtLang("fremdspracheText"),
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
			name:     "Text with op child element",
			text:     "Test text",
			child:    elem("op", nil, "opText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    elem("romzahl", nil, "2.", nil),
			expected: "Test text II.",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text " + util.FmtPage(384),
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text " + util.FmtLine(328),
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
			child:    elem("antiqua", nil, "antiquaText", nil),
			expected: "Test text antiquaText",
		},
		{
			name:     "Text with em1 child element",
			text:     "Test text",
			child:    elem("em1", nil, "em1Text", nil),
			expected: "Test text " + util.FmtEmph("em1Text"),
		},
		{
			name:     "Text with em2 child element",
			text:     "Test text",
			child:    elem("em2", nil, "em2Text", nil),
			expected: "Test text " + util.FmtEmph2("em2Text"),
		},
		{
			name:     "Text with fett child element",
			text:     "Test text",
			child:    elem("fett", nil, "fettText", nil),
			expected: "Test text " + util.FmtBold("fettText"),
		},
		{
			name:     "Text with formel child element",
			text:     "Test text",
			child:    elem("formel", nil, "formelText", nil),
			expected: "Test text " + util.FmtFormula("formelText"),
		},
		{
			name:     "Text with fr child element",
			text:     "Test text",
			child:    elem("fr", map[string]string{"seite": "1", "nr": "2"}, "", nil),
			expected: "Test text " + util.FmtFnRef(1, 2),
		},
		{
			name:     "Text with fremdsprache child element",
			text:     "Test text",
			child:    elem("fremdsprache", nil, "fremdspracheText", nil),
			expected: "Test text " + util.FmtLang("fremdspracheText"),
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
			name:     "Text with op child element",
			text:     "Test text",
			child:    elem("op", nil, "opText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with romzahl child element",
			text:     "Test text",
			child:    elem("romzahl", nil, "2.", nil),
			expected: "Test text II.",
		},
		{
			name:     "Text with seite child element",
			text:     "Test text",
			child:    elem("seite", map[string]string{"nr": "384"}, "", nil),
			expected: "Test text " + util.FmtPage(384),
		},
		{
			name:     "Text with table child element",
			text:     "Test text",
			child:    elem("table", nil, "tableText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with trenn child element",
			text:     "Test text",
			child:    elem("trenn", nil, "trennText", nil),
			expected: "Test text",
		},
		{
			name:     "Text with zeile child element",
			text:     "Test text",
			child:    elem("zeile", map[string]string{"nr": "328"}, "", nil),
			expected: "Test text " + util.FmtLine(328),
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
			expected: util.FmtPage(254),
		},
		{
			name:     "Text is ignored",
			text:     "Some text",
			attrs:    map[string]string{"nr": "847"},
			expected: util.FmtPage(847),
		},
		{
			name:     "Nr attribute with leading zeros",
			attrs:    map[string]string{"nr": "00045"},
			expected: util.FmtPage(45),
		},
		{
			name:     "Nr attribute with leading and trailing spaces",
			attrs:    map[string]string{"nr": " 2     "},
			expected: util.FmtPage(2),
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
			el := elem("element", tc.attrs, tc.text, nil)
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
	result, err := Table("xml")
	assert.Equal(t, "", result)
	assert.False(t, err.HasError)
}

func TestRandtext(t *testing.T) {
	testCases := []struct {
		name         string
		text         string
		attrs        map[string]string
		child        *etree.Element
		expectedText string
		expectedRef  string
		expectError  bool
	}{
		{
			name:         "Text with randtext attributes",
			text:         "Some text",
			attrs:        map[string]string{"seite": "123", "anfang": "567"},
			expectedText: "Some text",
			expectedRef:  "123.567",
		},
		{
			name:         "Text with p child element",
			text:         "Some text",
			attrs:        map[string]string{"seite": "123", "anfang": "567"},
			child:        elem("p", nil, "pText", nil),
			expectedText: "Some text pText",
			expectedRef:  "123.567",
		},
		{
			name:        "text with unknown child element",
			text:        "Some text",
			attrs:       map[string]string{"seite": "123", "anfang": "567"},
			child:       elem("my-custom-tag", nil, "", nil),
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
			el := elem("element", nil, tc.text, tc.child)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			text, ref, err := summary(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expectedText, text)
			assert.Equal(t, tc.expectedRef, ref)
		})
	}
}

func TestFootnote(t *testing.T) {
	testCases := []struct {
		name         string
		text         string
		attrs        map[string]string
		child        *etree.Element
		expectedText string
		expectedRef  string
		expectError  bool
	}{
		{
			name:         "Text with footnote attributes",
			text:         "Some text",
			attrs:        map[string]string{"seite": "123", "nr": "567"},
			expectedText: "Some text",
			expectedRef:  "123.567",
		},
		{
			name:         "Text with p child element",
			text:         "Some text",
			attrs:        map[string]string{"seite": "123", "nr": "567"},
			child:        elem("p", nil, "pText", nil),
			expectedText: "Some text pText",
			expectedRef:  "123.567",
		},
		{
			name:        "text with unknown child element",
			text:        "Some text",
			attrs:       map[string]string{"seite": "123", "nr": "567"},
			child:       elem("my-custom-tag", nil, "", nil),
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
			el := elem("element", nil, tc.text, tc.child)
			for k, v := range tc.attrs {
				el.CreateAttr(k, v)
			}
			text, ref, err := footnote(el)
			assert.Equal(t, tc.expectError, err.HasError)
			assert.Equal(t, tc.expectedText, text)
			assert.Equal(t, tc.expectedRef, ref)
		})
	}
}

func elem(tag string, attrs map[string]string, text string, child *etree.Element) *etree.Element {
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
