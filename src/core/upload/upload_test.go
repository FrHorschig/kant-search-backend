//go:build unit
// +build unit

package upload

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/frhorschig/kant-search-backend/core/upload/errors"
// 	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
// 	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
// 	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
// 	"github.com/golang/mock/gomock"
// 	"gotest.tools/v3/assert"
// )

// func TestUploadProcess(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
// 	workRepo := dbMocks.NewMockWorkRepo(ctrl)
// 	contentRepo := dbMocks.NewMockContentRepo(ctrl)
// 	xmlMapper := mocks.NewMockXmlMapper(ctrl)
// 	sut := &uploadProcessorImpl{
// 		volumeRepo:  volumeRepo,
// 		workRepo:    workRepo,
// 		contentRepo: contentRepo,
// 		xmlMapper:   xmlMapper,
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()

// 	testCases := []struct {
// 		name      string
// 		xml       string
// 		err       errors.ErrorNew
// 		mockCalls func()
// 		assert    func(t *testing.T)
// 	}{
// 		{
// 			name: "Processing is successful",
// 			xml:  "",
// 			err:  errors.NilError(),
// 			mockCalls: func() {
// 				xmlMapper.EXPECT().MapWorks(gomock.Any(), gomock.Any()).Return([]model.Work{}, errors.NilError())
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.mockCalls()
// 			err := sut.Process(ctx, 1, tc.xml)
// 			assert.Equal(t, tc.err, err)
// 		})
// 	}
// }
