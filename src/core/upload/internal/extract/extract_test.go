package extract

import (
	"fmt"
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestExtractNumericAttribute(t *testing.T) {
	tests := []struct {
		name          string
		element       *etree.Element
		attribute     string
		expectedValue int32
		expectErr     bool
	}{
		{
			name: "Valid numeric attribute",
			element: func() *etree.Element {
				el := etree.NewElement("test")
				el.CreateAttr("number", "42")
				return el
			}(),
			attribute:     "number",
			expectedValue: 42,
			expectErr:     false,
		},
		{
			name: "Missing attribute",
			element: func() *etree.Element {
				el := etree.NewElement("test")
				return el
			}(),
			attribute:     "number",
			expectedValue: 0,
			expectErr:     true,
		},
		{
			name: "Invalid numeric value",
			element: func() *etree.Element {
				el := etree.NewElement("test")
				el.CreateAttr("number", "not-a-number")
				return el
			}(),
			attribute:     "number",
			expectedValue: 0,
			expectErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, err := ExtractNumericAttribute(tc.element, tc.attribute)
			if tc.expectErr {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				assert.Equal(t, tc.expectedValue, value)
			}
		})
	}
}

func TestExtractFnRefs(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Valid footnote references",
			text:     "Text" + fnRef(2, 4) + " with " + fnRef(203, 238) + "fnRefs.",
			expected: []string{"2.4", "203.238"},
		},
		{
			name:     "No footnote references",
			text:     "Text without fn refs.",
			expected: []string{},
		},
		{
			name:     "Malformed footnote reference",
			text:     "This is a malformed reference: <ks-meta-fnref>letters</ks-meta-fnref>.",
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractFnRefs(tc.text)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestExtractPages(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		expected  []int32
		expectErr bool
	}{
		{
			name:     "Valid page numbers",
			text:     "Text with" + page(4) + " page " + page(23) + "numbers.",
			expected: []int32{4, 23},
		},
		{
			name:     "Valid page numbers",
			text:     page(9248) + "starting with" + page(284) + " a number",
			expected: []int32{9248, 284},
		},
		{
			name:     "No page numbers",
			text:     "Text without page numbers",
			expected: []int32{},
		},
		{
			name:      "Conversion error (int32 overflow)",
			text:      "Large page number: <ks-meta-page>1023147483648</ks-meta-page>.",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ExtractPages(tc.text)
			if tc.expectErr {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestRemoveTags(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Text with footnote references",
			text:     "This is a text with " + fnRef(1, 2) + " and " + fnRef(3, 4) + ".",
			expected: "This is a text with and .",
		},
		{
			name:     "Text with line matches",
			text:     "This is a text with " + line(12) + " tags.",
			expected: "This is a text with tags.",
		},
		{
			name:     "Text with page matches",
			text:     "This is a text with " + page(5) + " and " + page(10) + ".",
			expected: "This is a text with and .",
		},
		{
			name:     "Text with summary matches",
			text:     "This is a text with " + summ(5, 7) + " and " + summ(10, 22) + ".",
			expected: "This is a text with and .",
		},
		{
			name:     "Text with HTML tags",
			text:     "<div>This is <b>bold</b> and <i>italic</i>.</div>",
			expected: "This is bold and italic.",
		},
		{
			name:     "Text with mixed tags",
			text:     "Mixed " + fnRef(1, 2) + " and <b>HTML</b> tags " + page(3) + ".",
			expected: "Mixed and HTML tags .",
		},
		{
			name:     "Text with no tags",
			text:     "Plain text without any tags.",
			expected: "Plain text without any tags.",
		},
		{
			name:     "Empty string",
			text:     "",
			expected: "",
		},
		{
			name:     "Text with malformed HTML tags",
			text:     "Malformed <tag text.",
			expected: "Malformed <tag text.",
		},
		{
			name:     "Text with self-closing HTML tags",
			text:     "Image: <img src='image.jpg'/> here.",
			expected: "Image: here.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := RemoveTags(tc.text)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func line(line int32) string {
	return util.FmtLine(line)
}

func page(page int32) string {
	return util.FmtPage(page)
}

func fnRef(page, nr int32) string {
	return util.FmtFnRef(page, nr)
}

func summ(page, line int32) string {
	return util.FmtSummaryRef(fmt.Sprintf("%d.%d", page, line))
}
