//go:build unit
// +build unit

package upload

import (
	"context"
	"testing"

	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVolumeUploadProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockXmlMapper := mocks.NewMockXmlMapper(ctrl)
	mockParagraphRepo := dbMocks.NewMockParagraphRepo(ctrl)
	mockSentenceRepo := dbMocks.NewMockSentenceRepo(ctrl)
	processor := &volumeUploadProcessorImpl{
		xmlMapper:     mockXmlMapper,
		paragraphRepo: mockParagraphRepo,
		sentenceRepo:  mockSentenceRepo,
	}

	ctx := context.Background()

	testCases := []struct {
		name      string
		volNum    int32
		xml       abt1.Kantabt1
		err       error
		mockCalls func()
	}{
		{
			name:      "Processing is successful",
			xml:       abt1.Kantabt1{},
			volNum:    1,
			err:       nil,
			mockCalls: func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()
			err := processor.ProcessAbt1(ctx, tc.volNum, tc.xml)
			assert.Equal(t, tc.err, err)
		})
	}
}
