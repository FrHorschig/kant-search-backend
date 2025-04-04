package extract

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/stretchr/testify/assert"
)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ExtractFnRefs(tt.text)
			assert.Equal(t, tt.expected, actual)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ExtractPages(tt.text)
			if tt.expectErr {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				assert.Equal(t, tt.expected, actual)
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
			name:     "Text with nested tags",
			text:     "<div><span>Nested</span> tags</div>",
			expected: "Nested tags",
		},
		{
			name:     "Text with no tags",
			text:     "Plain text without tags.",
			expected: "Plain text without tags.",
		},
		{
			name:     "Empty string",
			text:     "",
			expected: "",
		},
		{
			name:     "Text with malformed tags",
			text:     "Malformed <tag text.",
			expected: "Malformed <tag text.",
		},
		{
			name:     "Text with self-closing tags",
			text:     "Image: <img src='image.jpg'/> here.",
			expected: "Image:  here.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := RemoveTags(tt.text)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func page(page int32) string {
	return util.FmtPage(page)
}

func fnRef(a int32, b int32) string {
	return util.FmtFnRef(a, b)
}
