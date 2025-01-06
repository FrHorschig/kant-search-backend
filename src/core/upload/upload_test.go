//go:build unit
// +build unit

package upload

import (
	"context"
	"errors"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUploadProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	xmlMapper := mocks.NewMockXmlMapper(ctrl)
	paragraphRepo := dbMocks.NewMockParagraphRepo(ctrl)
	sentenceRepo := dbMocks.NewMockSentenceRepo(ctrl)
	sut := &uploadProcessorImpl{
		xmlMapper:     xmlMapper,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}

	ctx := context.Background()
	e := errors.New("Mock error")

	testCases := []struct {
		name      string
		xml       string
		err       error
		mockCalls func()
		assert    func(t *testing.T)
	}{
		{
			name: "Processing is successful",
			xml:  "",
			err:  nil,
			mockCalls: func() {
				xmlMapper.EXPECT().Map(gomock.Any()).Return([]model.Work{}, nil)
			},
		},
		{
			name: "Processing error due to failed mapping",
			xml:  "",
			err:  e,
			mockCalls: func() {
				xmlMapper.EXPECT().Map(gomock.Any()).Return(nil, e)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()
			err := sut.Process(ctx, 1, tc.xml)
			assert.Equal(t, tc.err, err)
		})
	}
}
