//go:build unit
// +build unit

package upload

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ordinal int32 = 1

func TestUploadProcessSuccess(t *testing.T) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	volNr := int32(1)
	wCode := "workCode"
	vol := &model.Volume{
		VolumeNumber: volNr,
		Section:      2,
		Title:        "volume title",
	}
	work := model.Work{
		Code:         wCode,
		Abbreviation: util.StrPtr("abbrev"),
		Title:        "work title",
		Year:         "1785",
		Sections: []model.Section{
			{
				Heading: head(1),
				Paragraphs: []model.Paragraph{
					par(2),
					par(3),
				},
				Sections: []model.Section{
					{
						Heading: head(4),
						Paragraphs: []model.Paragraph{
							par(5),
							par(6),
						},
					},
					{
						Heading: head(7),
						Paragraphs: []model.Paragraph{
							par(8),
							par(9),
						},
					},
				},
			},
			{
				Heading: head(10),
				Paragraphs: []model.Paragraph{
					par(11),
					par(12),
				},
			},
		},
		Footnotes: []model.Footnote{
			fn(1),
			fn(2),
			fn(3),
			fn(4),
			fn(5),
			fn(6),
			fn(7),
			fn(8),
			fn(9),
			fn(10),
			fn(11),
			fn(12),
			fn(102),
			fn(103),
			fn(105),
			fn(106),
			fn(108),
			fn(109),
			fn(111),
			fn(112),
		},
		Summaries: []model.Summary{
			summ(2),
			summ(3),
			summ(5),
			summ(6),
			summ(8),
			summ(9),
			summ(11),
			summ(12),
		},
	}

	// GIVEN
	// mapping
	xmlMapper.EXPECT().
		MapVolume(gomock.Eq(volNr), gomock.Any()).
		Return(vol, errs.Nil())
	xmlMapper.EXPECT().
		MapWorks(gomock.Eq(volNr), gomock.Any()).
		Return([]model.Work{work}, errs.Nil())

	// data deletion
	volumeRepo.EXPECT().
		GetByVolumeNumber(gomock.Any(), gomock.Eq(volNr)).
		Return(&dbmodel.Volume{
			VolumeNumber: vol.VolumeNumber,
			Section:      vol.Section,
			Title:        vol.Title,
			Works: []dbmodel.Work{{
				Code:         work.Code,
				Abbreviation: work.Abbreviation,
				Title:        work.Title,
				Year:         work.Year,
				Ordinal:      1,
				Paragraphs:   []int32{},
				Sections: []dbmodel.Section{
					{
						Heading:    1,
						Paragraphs: []int32{5, 9},
						Sections: []dbmodel.Section{
							{
								Heading:    11,
								Paragraphs: []int32{15, 19},
							},
							{
								Heading:    21,
								Paragraphs: []int32{25, 29},
							},
						},
					},
					{
						Heading:    31,
						Paragraphs: []int32{35, 39},
					},
				},
			}},
		}, nil)
	contentRepo.EXPECT().DeleteByWork(gomock.Any(), gomock.Eq(wCode)).Return(nil)
	volumeRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(volNr)).Return(nil)

	// data insertion
	contentRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	volumeRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

	// WHEN
	err := sut.Process(ctx, volNr, "xml")

	// THEN
	assert.False(t, err.HasError)
}

func head(n int32) model.Heading {
	nr := strconv.Itoa(int(n))
	return model.Heading{
		Text:    "<ks-fmt-h1>heading</ks-fmt-h1> text " + nr,
		TocText: "toc text " + nr,
		Pages:   []int32{n},
		FnRefs:  []string{"fnRef" + nr},
	}
}

func par(n int32) model.Paragraph {
	nr := strconv.Itoa(int(n))
	return model.Paragraph{
		Text:       "<ks-fmt-bold>paragraph</ks-fmt-bold> text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		SummaryRef: util.StrPtr("summRef" + nr),
	}
}

func fn(n int32) model.Footnote {
	nr := strconv.Itoa(int(n))
	return model.Footnote{
		Text:  "<ks-fmt-emph>footnote</ks-fmt-emph> text " + nr,
		Ref:   "fnRef" + nr,
		Pages: []int32{n},
	}
}

func summ(n int32) model.Summary {
	nr := strconv.Itoa(int(n))
	nr100 := strconv.Itoa(int(n) + 100)
	return model.Summary{
		Text:   "<ks-fmt-tracked>summary</ks-fmt-tracked> text " + nr,
		Ref:    "summRef" + nr,
		Pages:  []int32{n},
		FnRefs: []string{"fnRef" + nr100},
	}
}

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
			name: "MapVolume fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				xm.EXPECT().MapVolume(gomock.Any(), gomock.Any()).
					Return(nil, errs.New(nil, testErr))
			},
		},
		{
			name: "MapWorks fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				gomock.InOrder(
					xm.EXPECT().MapVolume(gomock.Any(), gomock.Any()).
						Return(&model.Volume{}, errs.Nil()),
					xm.EXPECT().MapWorks(gomock.Any(), gomock.Any()).
						Return(nil, errs.New(nil, testErr)),
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
	mapper.EXPECT().MapVolume(gomock.Any(), gomock.Any()).Return(&model.Volume{}, errs.Nil())
	mapper.EXPECT().MapWorks(gomock.Any(), gomock.Any()).Return([]model.Work{
		{
			Code:         wCode,
			Title:        "t",
			Abbreviation: util.StrPtr("abbr"),
			Year:         "2024",
			Sections: []model.Section{{
				Heading: head(1),
				Paragraphs: []model.Paragraph{
					par(2),
				},
			}},
		},
	}, errs.Nil())
}

func mockDeletion(vr *dbMocks.MockVolumeRepo, cr *dbMocks.MockContentRepo, wCode string) {
	vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(
		&dbmodel.Volume{Works: []dbmodel.Work{{Code: wCode}}},
		nil,
	)
	cr.EXPECT().DeleteByWork(gomock.Any(), wCode).Return(nil)
	vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
}
