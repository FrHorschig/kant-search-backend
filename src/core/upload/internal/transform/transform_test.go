//go:build unit
// +build unit

package transform

import (
	"fmt"
	"testing"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutil/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func sp(s string) *string {
	return &s
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
					Pages: []int32{234}, HeadingLevel: 1},
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
					Pages: []int32{234}, FootnoteName: "12.45"},
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
					Pages: []int32{234}, HeadingLevel: 1},
				{WorkId: 1, Text: "some {l2} test {p324} text",
					Pages: []int32{234}, FootnoteName: "12.45"},
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
			name:  "wrong end expression l",
			pyErr: fmt.Errorf("error"),
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "paragraph"}, Content: sp("text")},
			},
			out: nil,
			err: &errors.Error{
				Msg:    errors.GO_ERR,
				Params: []string{"error"},
			},
		},
		{
			name:    "paragraph with incomplete sentence",
			pyValue: map[int32][]string{1: {"This is a sentence that is completes in the second paragraph."}},
			in: []c.Expression{
				{Metadata: c.Metadata{Class: "p", Param: sp("234")}},
				{Metadata: c.Metadata{Class: "paragraph"},
					Content: sp("This is a sentence that completes")},
				{Metadata: c.Metadata{Class: "p", Param: sp("235")}},
				{Metadata: c.Metadata{Class: "paragraph"},
					Content: sp("in the second paragraph.")},
			},
			out: []model.Paragraph{
				{WorkId: 1, Text: "{p234} This is a sentence that completes {p235} in the second paragraph.",
					Pages: []int32{234, 235}},
			},
			err: nil,
		},
	}

	ctrl := gomock.NewController(t)
	pyUtil := mocks.NewMockPythonUtil(ctrl)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.pyValue != nil || tc.pyErr != nil {
				pyUtil.EXPECT().SplitIntoSentences(gomock.Any()).Return(tc.pyValue, tc.pyErr)
			}
			out, err := Transform(1, tc.in, pyUtil)
			assert.Len(t, out, len(tc.out))
			for i := range tc.out {
				assert.Equal(t, tc.out[i], out[i])
			}
			assert.Equal(t, tc.err, err)
		})
	}
}
