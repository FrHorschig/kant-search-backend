//go:build unit
// +build unit

package upload

import (
	"context"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVolumeUploadProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockXmlMapper := mocks.NewMockXmlMapper(ctrl)
	mockParagraphRepo := dbMocks.NewMockParagraphRepo(ctrl)
	mockSentenceRepo := dbMocks.NewMockSentenceRepo(ctrl)
	sut := &volumeUploadProcessorImpl{
		xmlMapper:     mockXmlMapper,
		paragraphRepo: mockParagraphRepo,
		sentenceRepo:  mockSentenceRepo,
	}

	ctx := context.Background()

	testCases := []struct {
		name        string
		xml         []byte
		errContains string
		mockCalls   func()
	}{
		{
			name:        "Processing is successful",
			xml:         []byte(""),
			errContains: "",
			mockCalls: func() {
				mockXmlMapper.EXPECT().MapVolume(gomock.Any(), gomock.Any()).Return([]model.Work{}, nil)
			},
		},
		{
			name:        "Processing error due to invalid xml",
			xml:         []byte("<my-tag>"),
			errContains: "error unmarshaling request body:",
			mockCalls:   func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()
			err := sut.Process(ctx, tc.xml)
			if tc.errContains != "" {
				assert.Contains(t, err.Error(), tc.errContains)
			}
		})
	}
}
