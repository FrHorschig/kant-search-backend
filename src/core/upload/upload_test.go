//go:build unit
// +build unit

package upload

import (
	"context"
	stderr "errors"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
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
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
		xmlMapper:     xmlMapper,
	}

	ctx := context.Background()
	e := errors.NewError(stderr.New("domain error"), nil)

	testCases := []struct {
		name      string
		xml       string
		err       errors.ErrorNew
		mockCalls func()
		assert    func(t *testing.T)
	}{
		{
			name: "Processing is successful",
			xml:  "",
			err:  errors.NilError(),
			mockCalls: func() {
				xmlMapper.EXPECT().Map(gomock.Any()).Return([]model.Work{}, errors.NilError())
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
