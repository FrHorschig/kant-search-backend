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
			name: "valid xml with <zeile>",
			xml:  `<zeile nr="1"/>`,
			want: `{l1}`,
		},
		{
			name: "valid xml with <trenn>",
			xml:  `Test<trenn/>string`,
			want: `Teststring`,
		},
		{
			name: "valid xml with <romzahl>",
			xml:  `<romzahl> 2.</romzahl>`,
			want: ` II. `,
		},
		{
			name: "valid xml with multiple tags",
			xml:  `<zeile nr="10"/>Etwas<trenn/><romzahl>12. </romzahl><romzahl> 234. </romzahl><zeile nr="30"/>Text<trenn/>abschnitt<romzahl> 1. </romzahl>`,
			want: `{l10}Etwas XII. CCXXXIV. {l30}Textabschnitt I. `,
		},
		{
			name: "malformed <zeile> tag",
			xml:  `<zeile nr="42"`,
			want: `<zeile nr="42"`,
		},
		{
			name: "malformed <trenn> tag",
			xml:  `<trenn>`,
			want: `<trenn>`,
		},
		{
			name: "malformed <romzahl> tag",
			xml:  `<romzahl> 23. <romzahl>`,
			want: `<romzahl> 23. <romzahl>`,
		},
		{
			name: "valid xml with no match",
			xml:  `<note><to>Note</to></note>`,
			want: `<note><to>Note</to></note>`,
		},
		{
			name: "empty xml",
			xml:  ``,
			want: ``,
		},
		{
			name: "multiple spaces to single space",
			xml:  `   `,
			want: ` `,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Simplify(tc.xml)
			assert.Equal(t, tc.want, got)
		})
	}
}
