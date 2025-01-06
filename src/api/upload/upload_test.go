//go:build unit
// +build unit

package upload

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	procMocks "github.com/frhorschig/kant-search-backend/core/upload/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const xmlBase = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`
const abt1Xml = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<kant_abt1>
  <band nr="01">
    <titel>1</titel>
    <hauptteil>Hauptteil</hauptteil>
  </band>
</kant_abt1>
`

func TestUploadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	volumeProcessor := procMocks.NewMockUploadProcessor(ctrl)
	sut := NewUploadHandler(volumeProcessor).(*uploadHandlerImpl)

	testCases := []struct {
		name        string
		xml         string
		mockSuccess bool
		mockError   error
		wantCode    int
		wantMsg     string
	}{
		{
			name:        "Processing success",
			xml:         abt1Xml,
			mockSuccess: true,
			wantCode:    http.StatusCreated,
		},
		{
			name:     "Unknown encoding",
			xml:      `<?xml version="1.0" encoding="unknown" standalone="no"?><root><band nr="5"></band></root>`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "error unmarshaling request body",
		},
		{
			name:     "Invalid XML syntax",
			xml:      xmlBase + `<root`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "error unmarshaling request body",
		},
		{
			name:     "No band element with nr attribute",
			xml:      xmlBase + "<root><band></band></root>",
			wantCode: http.StatusBadRequest,
			wantMsg:  "missing element 'band' with attribute 'nr'",
		},
		{
			name:     "Nr attribute is not a number",
			xml:      xmlBase + `<root><band nr="abc"></band></root>`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "attribute 'nr' of element 'band' can't be converted to a number",
		},
		{
			name:     "Volume number is zero",
			xml:      xmlBase + `<root><band nr="0"></band></root>`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "the volume number is 0, but it must be between 1 and 9",
		},
		{
			name:     "Volume number too low",
			xml:      xmlBase + `<root><band nr="-1"></band></root>`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "the volume number is -1, but it must be between 1 and 9",
		},
		{
			name:     "Volume number too high",
			xml:      xmlBase + `<root><band nr="10"></band></root>`,
			wantCode: http.StatusNotImplemented,
			wantMsg:  "uploading volumes greater than 9 is not yet implemented",
		},
		{
			name:      "Error processing XML",
			xml:       xmlBase + `<root><band nr="5"></band></root>`,
			mockError: errors.New("processing error"),
			wantCode:  http.StatusInternalServerError,
			wantMsg:   "error processing XML data for volume 5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := []byte(tc.xml)
			// GIVEN
			req := httptest.NewRequest(echo.GET, "/api/write/v1/volumes", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)
			if tc.mockSuccess {
				volumeProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			}
			if tc.mockError != nil {
				volumeProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.mockError)
			}

			// WHEN
			sut.PostVolume(ctx)

			// THEN

			if tc.wantCode == http.StatusCreated {
				assert.Equal(t, tc.wantCode, rec.Code)
			} else {
				assert.Equal(t, tc.wantCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.wantMsg)
			}
			assert.Equal(t, rec.Code, tc.wantCode)
		})
	}
}

func TestReplaceCustomEntities(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `<?xml version="1.0"?><root>&kreis; &quadrat;</root>`,
			expected: `<?xml version="1.0"?><root>○ ■</root>`,
		},
		{
			input:    `<root>&kreis;</root>`,
			expected: `<root>○</root>`,
		},
		{
			input:    `<root>&quadrat;</root>`,
			expected: `<root>■</root>`,
		},
		{
			input:    `<root>&kreis; &quadrat; &kreis;</root>`,
			expected: `<root>○ ■ ○</root>`,
		},
		{
			input:    `<root>no custom entities here</root>`,
			expected: `<root>no custom entities here</root>`,
		},
		{
			input:    ``,
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			actual := replaceCustomEntities(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}

func TestReplaceHtml(t *testing.T) {
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
			assert.Equal(t, replaceHtml(tc.input), tc.expected)
		})
	}
}
