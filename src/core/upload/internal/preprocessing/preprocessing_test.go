package preprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocess(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{"&alpha;", "α"},
		{"&beta;", "β"},
		{"&gamma;", "γ"},
		{"&delta;", "δ"},
		{"&epsilon;", "ε"},
		{"&omega;", "ω"},

		{"&Alpha;", "Α"},
		{"&Beta;", "Β"},
		{"&Gamma;", "Γ"},
		{"&Delta;", "Δ"},
		{"&Epsilon;", "Ε"},
		{"&Omega;", "Ω"},

		{"&sigma2;", "ς"},
		{"&theata;", "θ"},

		{"This is &alpha; and &beta;.", "This is α and β."},
		{"&Gamma;&Delta;&Epsilon;", "ΓΔΕ"},
		{"&lt;div&gt;Hello &amp; World&lt;/div&gt;", "<div>Hello & World</div>"},
		{"No encodings here!", "No encodings here!"},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, ReplaceHtml(tc.input), tc.expected)
		})
	}
}

func TestSimplify(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want string
	}{
		{
			name: "valid xml with <zeile> and <seite>",
			xml:  `<zeile nr="1"/><seite nr="2" />`,
			want: `{l1}{p2}`,
		},
		{
			name: "valid xml with <zeile>, <seite> and <trenn>",
			xml:  `<zeile nr="1"/><seite nr="2" />Text<trenn/>abschnitt<trenn />`,
			want: `{l1}{p2}Textabschnitt`,
		},
		{
			name: "valid xml with no match",
			xml:  `<note><to>Note</to></note>`,
			want: `<note><to>Note</to></note>`,
		},
		{
			name: "valid xml with multiple <zeile> and <seite>",
			xml:  `<zeile nr="10"/><seite nr="20" /><zeile nr="30"/>`,
			want: `{l10}{p20}{l30}`,
		},
		{
			name: "empty xml",
			xml:  ``,
			want: ``,
		},
		{
			name: "malformed <seite> tag",
			xml:  `<seite nr="42"`,
			want: `<seite nr="42"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Simplify(tc.xml)
			assert.Equal(t, got, tc.want)
		})
	}
}
