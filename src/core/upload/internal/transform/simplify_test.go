package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplify(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want string
	}{
		{
			name: "valid xml with <zeile> and <seite>",
			xml:  `<zeile nr="1"/><seite nr="2" />`,
			want: `{l1}<seite nr="2" />`,
		},
		{
			name: "valid xml with <zeile>, <seite> and <trenn>",
			xml:  `<zeile nr="1"/><seite nr="2" />Text<trenn/>abschnitt<trenn />`,
			want: `{l1}<seite nr="2" />Textabschnitt`,
		},
		{
			name: "valid xml with no match",
			xml:  `<note><to>Note</to></note>`,
			want: `<note><to>Note</to></note>`,
		},
		{
			name: "valid xml with multiple <zeile> and <seite>",
			xml:  `<zeile nr="10"/><seite nr="20" /><zeile nr="30"/>`,
			want: `{l10}<seite nr="20" />{l30}`,
		},
		{
			name: "empty xml",
			xml:  ``,
			want: ``,
		},
		{
			name: "malformed <zeile> tag",
			xml:  `<zeile nr="42"`,
			want: `<zeile nr="42"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Simplify(tc.xml)
			assert.Equal(t, got, tc.want)
		})
	}
}
