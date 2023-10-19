//go:build unit
// +build unit

package transform

import (
	"testing"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"
	"github.com/stretchr/testify/assert"
)

func sp(s string) *string {
	return &s
}

func ip(s int) *int32 {
	return &[]int32{int32(s)}[0]
}

func TestTransform(t *testing.T) {
	testCases := []struct {
		name    string
		pyValue map[int32][]string
		pyErr   error
		in      []c.Expression
		out     []model.Paragraph
		err     *errors.Error
	}{
		{
			name:    "paragraph expression",
			pyValue: map[int32][]string{},
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "paragraph"},
					Content: sp("some {l2} test {p324} text")},
			},
			out: []model.Paragraph{
				{WorkId: 1, Text: "{p234} some {l2} test {p324} text",
					Pages: []int32{234}},
			},
			err: nil,
		},
		{
			name:    "heading expression",
			pyValue: map[int32][]string{},
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "heading", Param: sp("1")},
					Content: sp("some {l2} test {p324} text")},
			},
			out: []model.Paragraph{
				{WorkId: 1, Text: "{p234} some {l2} test {p324} text",
					Pages: []int32{234}, HeadingLevel: ip(1)},
			},
			err: nil,
		},
		{
			name:    "footnote expression",
			pyValue: map[int32][]string{},
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "footnote", Param: sp("12.45")},
					Content: sp("some {l2} test {p324} text")},
			},
			out: []model.Paragraph{
				{WorkId: 1, Text: "{p234} some {l2} test {p324} text",
					Pages: []int32{234}, FootnoteName: sp("12.45")},
			},
			err: nil,
		},
		{
			name:    "all three expression types",
			pyValue: map[int32][]string{},
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "paragraph"},
					Content: sp("some {l2} test {p324} text")},
				{Metadata: c.Metadata{Class: "paragraph"},
					Content: sp("some {l2} test {p324} text")},
				{Metadata: c.Metadata{Class: "heading", Param: sp("1")},
					Content: sp("some {l2} test {p324} text")},
				{Metadata: c.Metadata{Class: "footnote", Param: sp("12.45")},
					Content: sp("some {l2} test {p324} text")},
			},
			out: []model.Paragraph{
				{WorkId: 1, Text: "{p234} some {l2} test {p324} text",
					Pages: []int32{234}},
				{WorkId: 1, Text: "some {l2} test {p324} text",
					Pages: []int32{234}},
				{WorkId: 1, Text: "some {l2} test {p324} text",
					Pages: []int32{234}, HeadingLevel: ip(1)},
				{WorkId: 1, Text: "some {l2} test {p324} text",
					Pages: []int32{234}, FootnoteName: sp("12.45")},
			},
			err: nil,
		},
		{
			name: "wrong start expression",
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "l", Param: sp("234")}},
			},
			out: nil,
			err: &errors.Error{
				Msg:    errors.WRONG_START_EXPRESSION,
				Params: []string{"l"},
			},
		},
		{
			name: "wrong end expression p",
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
			},
			out: nil,
			err: &errors.Error{
				Msg:    errors.WRONG_END_EXPRESSION,
				Params: []string{"p"},
			},
		},
		{
			name: "wrong end expression l",
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "l", Param: sp("234")}},
			},
			out: nil,
			err: &errors.Error{
				Msg:    errors.WRONG_END_EXPRESSION,
				Params: []string{"l"},
			},
		},
		{
			name: "unexpected expression class",
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "something"}, Content: sp("text")},
			},
			out: nil,
			err: &errors.Error{
				Msg:    errors.UNKNOWN_EXPRESSION_CLASS,
				Params: []string{"something"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Transform(1, tc.in)
			assert.Len(t, out, len(tc.out))
			for i := range tc.out {
				assert.Equal(t, tc.out[i], out[i])
			}
			assert.Equal(t, tc.err, err)
		})
	}
}
