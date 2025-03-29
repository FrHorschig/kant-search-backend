package extract

import (
	"fmt"
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
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
			expected: []int32{3, 4, 23},
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

func TestFindParagraph(t *testing.T) {
	tests := []struct {
		name     string
		section  dbmodel.Section
		page     int32
		line     int32
		expectId int32
	}{
		{
			name: "Find Paragraph in top-level section",
			section: dbmodel.Section{
				Paragraphs: []dbmodel.Paragraph{
					{
						Id:    1,
						Text:  page(1) + line(194),
						Pages: []int32{1},
					},
					{
						Id:    2,
						Text:  page(2) + line(193),
						Pages: []int32{2},
					},
					{
						Id:    3,
						Text:  page(2) + line(194),
						Pages: []int32{2},
					},
					{
						Id:    4,
						Text:  page(2) + line(195),
						Pages: []int32{2},
					},
				},
			},
			page:     2,
			line:     194,
			expectId: 3,
		},
		{
			name: "Find Paragraph in subsection",
			section: dbmodel.Section{
				Paragraphs: []dbmodel.Paragraph{
					{
						Id:    1,
						Text:  page(1) + line(194),
						Pages: []int32{1},
					},
				},
				Sections: []dbmodel.Section{{
					Paragraphs: []dbmodel.Paragraph{{
						Id:    2,
						Text:  page(2) + line(194),
						Pages: []int32{2},
					}},
				}},
			},
			page:     2,
			line:     194,
			expectId: 2,
		},
		{
			name: "Find Paragraph with text before page and line",
			section: dbmodel.Section{
				Paragraphs: []dbmodel.Paragraph{
					{
						Id:    1,
						Text:  "Text before " + page(2) + " text continued" + line(194) + " text after.",
						Pages: []int32{2},
					},
				},
			},
			page:     2,
			line:     194,
			expectId: 1,
		},
		{
			name: "Find Paragraph with pages before and after",
			section: dbmodel.Section{
				Paragraphs: []dbmodel.Paragraph{
					{
						Id:    1,
						Text:  page(5) + line(5) + page(6) + line(5) + page(7) + line(5) + page(8) + line(5),
						Pages: []int32{7},
					},
				},
			},
			page:     7,
			line:     5,
			expectId: 1,
		},
		{
			name: "Find Paragraph with lines before and after",
			section: dbmodel.Section{
				Paragraphs: []dbmodel.Paragraph{
					{
						Id:    1,
						Text:  page(7) + line(3) + line(4) + line(5) + line(6),
						Pages: []int32{7},
					},
				},
			},
			page:     7,
			line:     5,
			expectId: 1,
		},
		{
			name: "",
			section: dbmodel.Section{
				Paragraphs: []dbmodel.Paragraph{},
				Sections:   []dbmodel.Section{},
			},
			page:     1,
			line:     1,
			expectId: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			par := FindParagraph(&tt.section, tt.page, tt.line)
			assert.NotNil(t, par)
			assert.Equal(t, tt.expectId, par.Id)
		})
	}
}

func page(page int32) string {
	return fmt.Sprintf(model.PageFmt, page)
}

func line(line int32) string {
	return fmt.Sprintf(model.LineFmt, line)
}

func fnRef(a int, b int) string {
	return fmt.Sprintf(model.FnRefFmt, a, b)
}
