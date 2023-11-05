//go:build unit
// +build unit

package transform

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeParagraphs(t *testing.T) {
	testCases := []struct {
		name string
		in   []model.Paragraph
		out  []model.Paragraph
	}{
		{
			name: "no merge for end punctuations",
			in: []model.Paragraph{
				{Text: "This is a sentence.", Pages: []int32{1}},
				{Text: "This is another sentence!", Pages: []int32{2}},
				{Text: "Is this a question?", Pages: []int32{3}},
				{Text: "This is the final sentence.", Pages: []int32{4}},
			},
			out: []model.Paragraph{
				{Text: "This is a sentence.", Pages: []int32{1}},
				{Text: "This is another sentence!", Pages: []int32{2}},
				{Text: "Is this a question?", Pages: []int32{3}},
				{Text: "This is the final sentence.", Pages: []int32{4}},
			},
		},
		{
			name: "merge on non-end punctuations",
			in: []model.Paragraph{
				{Text: "This is a sentence,", Pages: []int32{1}},
				{Text: "this is another sentence:", Pages: []int32{2}},
				{Text: "Is this a question;", Pages: []int32{3}},
				{Text: "this is the final sentence", Pages: []int32{4}},
			},
			out: []model.Paragraph{
				{
					Text:  "This is a sentence, this is another sentence: Is this a question; this is the final sentence",
					Pages: []int32{1, 2, 3, 4},
				},
			},
		},
		{
			name: "ignore merging with in-between heading",
			in: []model.Paragraph{
				{Text: "This is a sentence", Pages: []int32{1}},
				{Text: "This is a heading", Pages: []int32{2}, HeadingLevel: &[]int32{1}[0]},
				{Text: "Is this a question?", Pages: []int32{3}},
			},
			out: []model.Paragraph{
				{Text: "This is a sentence", Pages: []int32{1}},
				{Text: "This is a heading", Pages: []int32{2}, HeadingLevel: &[]int32{1}[0]},
				{Text: "Is this a question?", Pages: []int32{3}},
			},
		},
		{
			name: "merge with in-between footnote",
			in: []model.Paragraph{
				{Text: "This is a sentence", Pages: []int32{1}},
				{Text: "This is a footnote", Pages: []int32{2}, FootnoteName: &[]string{"12.34"}[0]},
				{Text: "that continues {fn12.34}.", Pages: []int32{3}},
			},
			out: []model.Paragraph{
				{Text: "This is a sentence that continues {fn12.34}.", Pages: []int32{1, 3}},
				{Text: "This is a footnote", Pages: []int32{2}, FootnoteName: &[]string{"12.34"}[0]},
			},
		},
		{
			name: "ignore merging inside the same page",
			in: []model.Paragraph{
				{Text: "This is a sentence", Pages: []int32{1}},
				{Text: "This is a heading", Pages: []int32{1}},
			},
			out: []model.Paragraph{
				{Text: "This is a sentence", Pages: []int32{1}},
				{Text: "This is a heading", Pages: []int32{1}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := MergeParagraphs(tc.in)
			assert.Equal(t, out, tc.out)
		})
	}
}
