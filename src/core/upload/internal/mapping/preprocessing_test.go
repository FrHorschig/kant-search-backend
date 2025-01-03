package mapping

import (
	"testing"
)

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name    string
		xml     []byte
		want    string
		wantErr bool
	}{
		{
			name:    "valid xml with <zeile> and <seite>",
			xml:     []byte(`<zeile nr="1"/><seite nr="2" />`),
			want:    `{l1}{p2}`,
			wantErr: false,
		},
		{
			name:    "valid xml with <zeile>, <seite> and <trenn>",
			xml:     []byte(`<zeile nr="1"/><seite nr="2" />Text<trenn/>abschnitt<trenn />`),
			want:    `{l1}{p2}Textabschnitt`,
			wantErr: false,
		},
		{
			name:    "valid xml with no match",
			xml:     []byte(`<note><to>Note</to></note>`),
			want:    `<note><to>Note</to></note>`,
			wantErr: false,
		},
		{
			name:    "valid xml with multiple <zeile> and <seite>",
			xml:     []byte(`<zeile nr="10"/><seite nr="20" /><zeile nr="30"/>`),
			want:    `{l10}{p20}{l30}`,
			wantErr: false,
		},
		{
			name:    "empty xml",
			xml:     []byte(``),
			want:    ``,
			wantErr: false,
		},
		{
			name:    "malformed <seite> tag",
			xml:     []byte(`<seite nr="42"`),
			want:    `<seite nr="42"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Preprocess(tt.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf("preprocess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("preprocess() = %v, want %v", got, tt.want)
			}
		})
	}
}
