//go:build unit
// +build unit

package upload

import (
	"context"
	"testing"

	"github.com/beevik/etree"
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
		name      string
		volNum    int32
		xml       *etree.Document
		err       error
		mockCalls func()
	}{
		{
			name:      "Processing is successful",
			xml:       etree.NewDocument(),
			volNum:    1,
			err:       nil,
			mockCalls: func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()
			err := sut.Process(ctx, tc.volNum, tc.xml)
			assert.Equal(t, tc.err, err)
		})
	}
}
