package internal

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMapVolume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sut := NewXmlMapper()

	tests := []struct {
		name           string
		xmlInput       string
		inputVolNr     int32
		expectedVolNr  int32
		expectedTitle  string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name: "success",
			xmlInput: `
				<band nr="3">
					<titel>Kritik der reinen Vernunft</titel>
				</band>
			`,
			inputVolNr:    3,
			expectedVolNr: 3,
			expectedTitle: "Kritik der reinen Vernunft",
		},
		{
			name: "mismatched volume numbers",
			xmlInput: `
				<band nr="2">
					<titel>Prolegomena</titel>
				</band>
			`,
			inputVolNr:     3,
			expectError:    true,
			expectedErrMsg: "non matching volume numbers",
		},
		{
			name: "missing volume number attribute",
			xmlInput: `
				<band>
					<titel>Ohne Nummer</titel>
				</band>
			`,
			inputVolNr:     1,
			expectError:    true,
			expectedErrMsg: "missing 'nr' attribute in 'band' element",
		},
		{
			name: "volume number out of range",
			xmlInput: `
				<band nr="99">
					<titel>Unbekanntes Werk</titel>
				</band>
			`,
			inputVolNr:     99,
			expectError:    true,
			expectedErrMsg: "invalid volume number 99",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vol, err := sut.MapVolume(tc.inputVolNr, tc.xmlInput)

			if tc.expectError {
				assert.True(t, err.HasError)
				if err.DomainError != nil {
					assert.Contains(t, err.DomainError.Error(), tc.expectedErrMsg)
				} else {
					assert.Contains(t, err.TechnicalError.Error(), tc.expectedErrMsg)
				}
			} else {
				assert.False(t, err.HasError)
				assert.NotNil(t, vol)
				assert.Equal(t, tc.expectedVolNr, vol.VolumeNumber)
				assert.Equal(t, int32(1), vol.Section)
				assert.Equal(t, tc.expectedTitle, vol.Title)
			}
		})
	}
}
