package transform

import (
	stderr "errors"
	"testing"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/stretchr/testify/assert"
)

func TestExtractText(t *testing.T) {
	testCases := []struct {
		name           string
		before         string
		child          *etree.Element
		switchFnString string
		switchFnErr    errors.ErrorNew
		after          string
		expected       string
		expectError    bool
	}{
		{
			name:     "only before text",
			before:   "Some text before",
			expected: "Some text before",
		},
		{
			name:     "only after text",
			after:    "Some text after",
			expected: "Some text after",
		},
		{
			name:     "text before and after",
			before:   "first text",
			after:    "second text",
			expected: "first text second text",
		},
		{
			name:     "text with leading and trailing spaces",
			before:   "   first text ",
			after:    " second text     ",
			expected: "first text second text",
		},
		{
			name:           "switchFn returns success",
			before:         "text one",
			child:          createElement("my-tag", nil, "", nil),
			switchFnString: "switch fn result",
			after:          "text two",
			expected:       "text one switch fn result text two",
		},
		{
			name:        "switchFn returns error",
			before:      "text one",
			child:       createElement("my-tag", nil, "", nil),
			switchFnErr: errors.NewError(stderr.New("some error"), nil),
			after:       "text two",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			el := createElement("element", nil, tc.before, tc.child)
			if tc.after != "" {
				el.CreateText(tc.after)
			}
			switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
				return tc.switchFnString, tc.switchFnErr
			}

			result, err := extractText(el, switchFn)
			if tc.expectError {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}
