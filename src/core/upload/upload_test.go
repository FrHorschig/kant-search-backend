//go:build unit
// +build unit

package upload

import (
	"context"
	"fmt"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ordinal int32 = 1

func TestUploadProcessErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	xmlMapper := mocks.NewMockXmlMapper(ctrl)

	sut := &uploadProcessorImpl{
		volumeRepo:  volumeRepo,
		contentRepo: contentRepo,
		xmlMapper:   xmlMapper,
	}
	wCode := "code"
	testErr := fmt.Errorf("new error for vol num %d", 1)

	tests := []struct {
		name      string
		mockSetup func(*dbMocks.MockVolumeRepo, *dbMocks.MockContentRepo, *mocks.MockXmlMapper)
	}{
		{
			name: "MapXml fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				gomock.InOrder(
					xm.EXPECT().MapXml(gomock.Any(), gomock.Any()).
						Return(dbmodel.Volume{}, nil, errs.New(nil, testErr)),
				)
			},
		},
		{
			name: "GetByVolumeNumber fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(nil, testErr)
			},
		},
		{
			name: "ContentRepo.DeleteByWorkCode fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&dbmodel.Volume{
							Works: []dbmodel.Work{{Code: wCode}},
						}, nil),
					cr.EXPECT().DeleteByWork(gomock.Any(), wCode).
						Return(testErr),
				)
			},
		},
		{
			name: "VolumeRepo.Delete fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&dbmodel.Volume{
							Works: []dbmodel.Work{{Code: wCode}},
						}, nil),
					cr.EXPECT().DeleteByWork(gomock.Any(), wCode).
						Return(nil),
					vr.EXPECT().Delete(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
			},
		},
		{
			name: "Insert content fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				mockDeletion(vr, cr, wCode)
				cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(testErr)
				mockDeletion(vr, cr, wCode)
			},
		},
		{
			name: "Insert volume fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				mockDeletion(vr, cr, wCode)
				gomock.InOrder(
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil),
					vr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, cr, wCode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(volumeRepo, contentRepo, xmlMapper)
			err := sut.Process(context.Background(), 1, "xml")
			assert.True(t, err.HasError)
		})
	}
}

func mockXmlMapper(mapper *mocks.MockXmlMapper, wCode string) {
	mapper.EXPECT().MapXml(gomock.Any(), gomock.Any()).Return(
		dbmodel.Volume{},
		[]dbmodel.Content{{}},
		errs.Nil(),
	)
}

func mockDeletion(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, wCode string) {
	vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(
		&dbmodel.Volume{Works: []dbmodel.Work{{Code: wCode}}},
		nil,
	)
	cr.EXPECT().DeleteByWork(gomock.Any(), wCode).Return(nil)
	vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
}
