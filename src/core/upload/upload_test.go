//go:build unit
// +build unit

package upload

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ordinal int32 = 1

func TestUploadProcessSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
	workRepo := dbMocks.NewMockWorkRepo(ctrl)
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	xmlMapper := mocks.NewMockXmlMapper(ctrl)
	sut := &uploadProcessorImpl{
		volumeRepo:  volumeRepo,
		workRepo:    workRepo,
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
		Year:         util.StrPtr("1785"),
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
		Return(vol, errors.Nil())
	xmlMapper.EXPECT().
		MapWorks(gomock.Eq(volNr), gomock.Any()).
		Return([]model.Work{work}, errors.Nil())

	// data deletion
	volumeRepo.EXPECT().
		GetByVolumeNumber(gomock.Any(), gomock.Eq(volNr)).
		Return(&esmodel.Volume{
			VolumeNumber: vol.VolumeNumber,
			Section:      vol.Section,
			Title:        vol.Title,
			Works: []esmodel.WorkRef{{
				Code:         work.Code,
				Abbreviation: work.Abbreviation,
				Title:        work.Title,
			}},
		}, nil)
	contentRepo.EXPECT().DeleteByWorkCode(gomock.Any(), gomock.Eq(wCode)).Return(nil)
	workRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(wCode)).Return(nil)
	volumeRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(volNr)).Return(nil)

	// data insertion
	expectHeading(contentRepo, 1, wCode)
	expectParagraphs(contentRepo, 2, 3, wCode)
	expectHeading(contentRepo, 4, wCode)
	expectParagraphs(contentRepo, 5, 6, wCode)
	expectHeading(contentRepo, 7, wCode)
	expectParagraphs(contentRepo, 8, 9, wCode)
	expectHeading(contentRepo, 10, wCode)
	expectParagraphs(contentRepo, 11, 12, wCode)
	workRepo.EXPECT().Insert(gomock.Any(), gomock.Eq(
		&esmodel.Work{
			Code:         work.Code,
			Abbreviation: work.Abbreviation,
			Title:        work.Title,
			Year:         work.Year,
			Ordinal:      1,
			Paragraphs:   []int32{},
			Sections: []esmodel.Section{
				{
					Heading:    1,
					Paragraphs: []int32{5, 9},
					Sections: []esmodel.Section{
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
		})).Return(nil)
	volumeRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

	// WHEN
	err := sut.Process(ctx, volNr, "xml")

	// THEN
	assert.False(t, err.HasError)
}

func expectHeading(contentRepo *dbMocks.MockContentRepo, n int32, wCode string) {
	head := esHead(n, wCode)
	fn := esFn(n, wCode)
	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq([]esmodel.Content{head, fn})).
		Return(nil)
}

func expectParagraphs(contentRepo *dbMocks.MockContentRepo, n1 int32, n2 int32, wCode string) {
	summ1 := esSumm(n1, wCode)
	fn100 := esFn(n1+100, wCode)
	par1 := esPar(n1, wCode)
	fn1 := esFn(n1, wCode)

	summ2 := esSumm(n2, wCode)
	fn200 := esFn(n2+100, wCode)
	par2 := esPar(n2, wCode)
	fn2 := esFn(n2, wCode)

	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq([]esmodel.Content{summ1, fn100, par1, fn1, summ2, fn200, par2, fn2})).
		Return(nil)

}

func head(n int32) model.Heading {
	nr := strconv.Itoa(int(n))
	return model.Heading{
		Text:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText: "toc text " + nr,
		Pages:   []int32{n},
		FnRefs:  []string{"fnRef" + nr},
	}
}

func esHead(n int32, workCode string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	head := esmodel.Content{
		Type:       esmodel.Heading,
		Ordinal:    ordinal,
		FmtText:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText:    util.StrPtr("toc text " + nr),
		SearchText: "heading text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		WorkCode:   workCode,
	}
	ordinal += 1
	return head
}

func par(n int32) model.Paragraph {
	nr := strconv.Itoa(int(n))
	return model.Paragraph{
		Text:       "<fmt-tag>paragraph</fmt-tag> text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		SummaryRef: util.StrPtr("summRef" + nr),
	}
}

func esPar(n int32, workCode string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	par := esmodel.Content{
		Type:       esmodel.Paragraph,
		Ordinal:    ordinal,
		FmtText:    "<fmt-tag>paragraph</fmt-tag> text " + nr,
		SearchText: "paragraph text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		SummaryRef: util.StrPtr("summRef" + nr),
		WorkCode:   workCode,
	}
	ordinal += 1
	return par
}

func fn(n int32) model.Footnote {
	nr := strconv.Itoa(int(n))
	return model.Footnote{
		Text:  "<fmt-tag>footnote</fmt-tag> text " + nr,
		Ref:   "fnRef" + nr,
		Pages: []int32{n},
	}
}

func esFn(n int32, workCode string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	fn := esmodel.Content{
		Type:       esmodel.Footnote,
		Ordinal:    ordinal,
		Ref:        util.StrPtr("fnRef" + nr),
		FmtText:    "<fmt-tag>footnote</fmt-tag> text " + nr,
		SearchText: "footnote text " + nr,
		Pages:      []int32{n},
		WorkCode:   workCode,
	}
	ordinal += 1
	return fn
}

func summ(n int32) model.Summary {
	nr := strconv.Itoa(int(n))
	nr100 := strconv.Itoa(int(n) + 100)
	return model.Summary{
		Text:   "<fmt-tag>summary</fmt-tag> text " + nr,
		Ref:    "summRef" + nr,
		Pages:  []int32{n},
		FnRefs: []string{"fnRef" + nr100},
	}
}

func esSumm(n int32, workCode string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	nr100 := strconv.Itoa(int(n) + 100)
	summ := esmodel.Content{
		Type:       esmodel.Summary,
		Ordinal:    ordinal,
		Ref:        util.StrPtr("summRef" + nr),
		FmtText:    "<fmt-tag>summary</fmt-tag> text " + nr,
		SearchText: "summary text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr100},
		WorkCode:   workCode,
	}
	ordinal += 1
	return summ
}
func TestUploadProcessErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
	workRepo := dbMocks.NewMockWorkRepo(ctrl)
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	xmlMapper := mocks.NewMockXmlMapper(ctrl)

	sut := &uploadProcessorImpl{
		volumeRepo:  volumeRepo,
		workRepo:    workRepo,
		contentRepo: contentRepo,
		xmlMapper:   xmlMapper,
	}
	wCode := "code"
	testErr := fmt.Errorf("new error for vol num %d", 1)

	tests := []struct {
		name      string
		mockSetup func(*dbMocks.MockVolumeRepo, *dbMocks.MockWorkRepo, *dbMocks.MockContentRepo, *mocks.MockXmlMapper)
	}{
		{
			name: "MapVolume fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				xm.EXPECT().MapVolume(gomock.Any(), gomock.Any()).
					Return(nil, errors.New(nil, testErr))
			},
		},
		{
			name: "MapWorks fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				gomock.InOrder(
					xm.EXPECT().MapVolume(gomock.Any(), gomock.Any()).
						Return(&model.Volume{}, errors.Nil()),
					xm.EXPECT().MapWorks(gomock.Any(), gomock.Any()).
						Return(nil, errors.New(nil, testErr)),
				)
			},
		},
		{
			name: "GetByVolumeNumber fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(nil, testErr)
			},
		},
		{
			name: "ContentRepo.DeleteByWorkCode fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&esmodel.Volume{
							Works: []esmodel.WorkRef{{Code: wCode}},
						}, nil),
					cr.EXPECT().DeleteByWorkCode(gomock.Any(), wCode).
						Return(testErr),
				)
			},
		},
		{
			name: "WorkRepo.Delete fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&esmodel.Volume{
							Works: []esmodel.WorkRef{{Code: wCode}},
						}, nil),
					cr.EXPECT().DeleteByWorkCode(gomock.Any(), wCode).
						Return(nil),
					wr.EXPECT().Delete(gomock.Any(), wCode).Return(testErr),
				)
			},
		},
		{
			name: "VolumeRepo.Delete fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&esmodel.Volume{
							Works: []esmodel.WorkRef{{Code: wCode}},
						}, nil),
					cr.EXPECT().DeleteByWorkCode(gomock.Any(), wCode).
						Return(nil),
					wr.EXPECT().Delete(gomock.Any(), wCode).Return(nil),
					vr.EXPECT().Delete(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
			},
		},
		{
			name: "Insert heading fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				mockDeletion(vr, wr, cr, wCode)
				cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(testErr)
				mockDeletion(vr, wr, cr, wCode)
			},
		},
		{
			name: "Insert paragraphs fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				mockDeletion(vr, wr, cr, wCode)
				gomock.InOrder(
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wCode)
			},
		},
		{
			name: "Insert work fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				mockDeletion(vr, wr, cr, wCode)
				gomock.InOrder(
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(2),
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(testErr),
				)
				mockDeletion(vr, wr, cr, wCode)
			},
		},
		{
			name: "Insert volume fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm, wCode)
				mockDeletion(vr, wr, cr, wCode)
				gomock.InOrder(
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Times(2).Return(nil),
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil),
					vr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wCode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(volumeRepo, workRepo, contentRepo, xmlMapper)
			err := sut.Process(context.Background(), 1, "xml")
			assert.True(t, err.HasError)
		})
	}
}

func mockXmlMapper(mapper *mocks.MockXmlMapper, wCode string) {
	mapper.EXPECT().MapVolume(gomock.Any(), gomock.Any()).Return(&model.Volume{}, errors.Nil())
	mapper.EXPECT().MapWorks(gomock.Any(), gomock.Any()).Return([]model.Work{
		{
			Code:         wCode,
			Title:        "t",
			Abbreviation: util.StrPtr("abbr"),
			Year:         util.StrPtr("2024"),
			Sections: []model.Section{{
				Heading: head(1),
				Paragraphs: []model.Paragraph{
					par(2),
				},
			}},
		},
	}, errors.Nil())
}

func mockDeletion(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, wCode string) {
	vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(&esmodel.Volume{
		Works: []esmodel.WorkRef{{Code: wCode}},
	}, nil)
	cr.EXPECT().DeleteByWorkCode(gomock.Any(), wCode).Return(nil)
	wr.EXPECT().Delete(gomock.Any(), wCode).Return(nil)
	vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
}
