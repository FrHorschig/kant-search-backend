package extract

import (
	"fmt"
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
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
			text:     "This is a malformed reference: <ks-fmt-fnref>letters</ks-fmt-fnref>.",
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
			expected: []int32{1, 2, 3, 4},
		},
		{
			name:     "No page numbers",
			text:     "Text without page numbers",
			expected: []int32{},
		},
		{
			name:      "Malformed page number",
			text:      "Invalid page: <ks-meta-page>letters234</ks-meta-page>abc.",
			expected:  []int32{},
			expectErr: true,
		},
		{
			name:      "Conversion error (int32 overflow)",
			text:      "Large page number: <ks-meta-page>1023147483648</ks-meta-page>abc.",
			expected:  []int32{},
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

func page(page int32) string {
	return fmt.Sprintf(model.PageFmt, page)
}

func fnRef(a int, b int) string {
	return fmt.Sprintf(model.FnRefFmt, a, b)
}
