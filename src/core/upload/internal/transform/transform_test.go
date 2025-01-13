package transform

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

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
