package util

import (
	"testing"

	"github.com/beevik/etree"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := ExtractNumericAttribute(tt.element, tt.attribute)
			if tt.expectErr {
				assert.True(t, err.HasError)
			} else {
				assert.False(t, err.HasError)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}
