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

var ordinal int32 = 0

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
	wId := "workId"
	vol := &model.Volume{
		VolumeNumber: 1,
		Section:      2,
		Title:        "volume title",
	}
	work := model.Work{
		Code:         "code",
		Abbreviation: util.StrPtr("abbrev"),
		Title:        "work title",
		Year:         util.StrPtr("1785"),
		Sections: []model.Section{
			{
				Heading: inHead(1),
				Paragraphs: []model.Paragraph{
					inPar(2),
					inPar(3),
				},
				Sections: []model.Section{
					{
						Heading: inHead(4),
						Paragraphs: []model.Paragraph{
							inPar(5),
							inPar(6),
						},
					},
					{
						Heading: inHead(7),
						Paragraphs: []model.Paragraph{
							inPar(8),
							inPar(9),
						},
					},
				},
			},
			{
				Heading: inHead(10),
				Paragraphs: []model.Paragraph{
					inPar(11),
					inPar(12),
				},
			},
		},
		Footnotes: []model.Footnote{
			inFn(1),
			inFn(2),
			inFn(3),
			inFn(4),
			inFn(5),
			inFn(6),
			inFn(7),
			inFn(8),
			inFn(9),
			inFn(10),
			inFn(11),
			inFn(12),
			inFn(102),
			inFn(103),
			inFn(105),
			inFn(106),
			inFn(108),
			inFn(109),
			inFn(111),
			inFn(112),
		},
		Summaries: []model.Summary{
			inSumm(2),
			inSumm(3),
			inSumm(5),
			inSumm(6),
			inSumm(8),
			inSumm(9),
			inSumm(11),
			inSumm(12),
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
				Id:    wId,
				Code:  work.Code,
				Title: work.Title,
			}},
		}, nil)
	volumeRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(volNr)).Return(nil)
	workRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(wId)).Return(nil)
	contentRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Eq(wId)).Return(nil)

	// data insertion
	workRepo.EXPECT().Insert(gomock.Any(), gomock.Eq(
		&esmodel.Work{
			Code:         work.Code,
			Abbreviation: work.Abbreviation,
			Title:        work.Title,
			Year:         work.Year,
			Sections:     []esmodel.Section{},
		})).
		Do(func(ctx context.Context, w *esmodel.Work) {
			w.Id = wId
		}).Return(nil)

	expectHeading(contentRepo, 1, wId)
	expectParagraphs(contentRepo, 2, 3, wId)
	expectHeading(contentRepo, 4, wId)
	expectParagraphs(contentRepo, 5, 6, wId)
	expectHeading(contentRepo, 7, wId)
	expectParagraphs(contentRepo, 8, 9, wId)
	expectHeading(contentRepo, 10, wId)
	expectParagraphs(contentRepo, 11, 12, wId)

	workRepo.EXPECT().Update(gomock.Any(), gomock.Eq(&esmodel.Work{
		Id:           wId,
		Code:         work.Code,
		Abbreviation: work.Abbreviation,
		Title:        work.Title,
		Year:         work.Year,
		Sections: []esmodel.Section{
			{
				Heading:    "headingId1",
				Paragraphs: []string{"paragraphId2", "paragraphId3"},
				Sections: []esmodel.Section{
					{
						Heading:    "headingId4",
						Paragraphs: []string{"paragraphId5", "paragraphId6"},
						Sections:   []esmodel.Section{},
					},
					{
						Heading:    "headingId7",
						Paragraphs: []string{"paragraphId8", "paragraphId9"},
						Sections:   []esmodel.Section{},
					},
				},
			},
			{
				Heading:    "headingId10",
				Paragraphs: []string{"paragraphId11", "paragraphId12"},
				Sections:   []esmodel.Section{},
			},
		},
	})).Return(nil)
	volumeRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

	// WHEN
	err := sut.Process(ctx, volNr, "xml")

	// THEN
	assert.False(t, err.HasError)
}

func expectHeading(contentRepo *dbMocks.MockContentRepo, n int32, wId string) {
	inHead := inHeadC(n, wId)
	outHead := outHeadC(n, wId)
	inFn := inFnC(n, wId)
	outFn := outFnC(n, wId)
	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq([]esmodel.Content{inHead, inFn})).
		Do(func(ctx context.Context, c []esmodel.Content) {
			c[0] = outHead
			c[1] = outFn
		}).Return(nil)
}

func expectParagraphs(contentRepo *dbMocks.MockContentRepo, n1 int32, n2 int32, wId string) {
	inSumm1 := inSummC(n1, wId)
	outSumm1 := outSummC(n1, wId)
	inFn100 := inFnC(n1+100, wId)
	outFn100 := outFnC(n1+100, wId)
	inPar1 := inParC(n1, wId)
	outPar1 := outParC(n1, wId)
	inFn1 := inFnC(n1, wId)
	outFn1 := outFnC(n1, wId)
	inSumm2 := inSummC(n2, wId)
	outSumm2 := outSummC(n2, wId)
	inFn200 := inFnC(n2+100, wId)
	outFn200 := outFnC(n2+100, wId)
	inPar2 := inParC(n2, wId)
	outPar2 := outParC(n2, wId)
	inFn2 := inFnC(n2, wId)
	outFn2 := outFnC(n2, wId)

	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq([]esmodel.Content{inSumm1, inFn100, inFn1, inSumm2, inFn200, inFn2})).
		Do(func(ctx context.Context, c []esmodel.Content) {
			c[0] = outSumm1
			c[1] = outFn100
			c[2] = outFn1
			c[3] = outSumm2
			c[4] = outFn200
			c[5] = outFn2
		}).Return(nil)
	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq([]esmodel.Content{inPar1, inPar2})).
		Do(func(ctx any, c []esmodel.Content) {
			c[0] = outPar1
			c[1] = outPar2
		}).Return(nil)

}

func inHead(n int32) model.Heading {
	nr := strconv.Itoa(int(n))
	return model.Heading{
		Text:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText: "toc text " + nr,
		Pages:   []int32{n},
		FnRefs:  []string{"fnRef" + nr},
	}
}

func inHeadC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Heading,
		Ordinal:    ordinal,
		FmtText:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText:    util.StrPtr("toc text " + nr),
		SearchText: "heading text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		WorkId:     workId,
	}
}

func outHeadC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	head := esmodel.Content{
		Id:         "headingId" + nr,
		Ordinal:    ordinal,
		Type:       esmodel.Heading,
		FmtText:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText:    util.StrPtr("toc text " + nr),
		SearchText: "heading text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		WorkId:     workId,
	}
	ordinal += 1
	return head
}

func inPar(n int32) model.Paragraph {
	nr := strconv.Itoa(int(n))
	return model.Paragraph{
		Text:       "<fmt-tag>paragraph</fmt-tag> text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		SummaryRef: util.StrPtr("summRef" + nr),
	}
}

func inParC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Paragraph,
		Ordinal:    ordinal,
		FmtText:    "<fmt-tag>paragraph</fmt-tag> text " + nr,
		SearchText: "paragraph text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		SummaryRef: util.StrPtr("summRef" + nr),
		WorkId:     workId,
	}
}

func outParC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	par := esmodel.Content{
		Type:       esmodel.Paragraph,
		Id:         "paragraphId" + nr,
		Ordinal:    ordinal,
		FmtText:    "<fmt-tag>paragraph</fmt-tag> text " + nr,
		SearchText: "paragraph text " + nr,
		Pages:      []int32{n},
		FnRefs:     []string{"fnRef" + nr},
		SummaryRef: util.StrPtr("summRef" + nr),
		WorkId:     workId,
	}
	ordinal += 1
	return par
}

func inFn(n int32) model.Footnote {
	nr := strconv.Itoa(int(n))
	return model.Footnote{
		Text:  "<fmt-tag>footnote</fmt-tag> text " + nr,
		Ref:   "fnRef" + nr,
		Pages: []int32{n},
	}
}

func inFnC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Footnote,
		Ordinal:    ordinal,
		Ref:        util.StrPtr("fnRef" + nr),
		FmtText:    "<fmt-tag>footnote</fmt-tag> text " + nr,
		SearchText: "footnote text " + nr,
		Pages:      []int32{n},
		WorkId:     workId,
	}
}

func outFnC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	fn := esmodel.Content{
		Type:       esmodel.Footnote,
		Id:         "footnoteId" + nr,
		Ordinal:    ordinal,
		Ref:        util.StrPtr("fnRef" + nr),
		FmtText:    "<fmt-tag>footnote</fmt-tag> text " + nr,
		SearchText: "footnote text " + nr,
		Pages:      []int32{n},
		WorkId:     workId,
	}
	ordinal += 1
	return fn
}

func inSumm(n int32) model.Summary {
	nr := strconv.Itoa(int(n))
	nr100 := strconv.Itoa(int(n) + 100)
	return model.Summary{
		Text:   "<fmt-tag>summary</fmt-tag> text " + nr,
		Ref:    "summRef" + nr,
		Pages:  []int32{n, n + 100},
		FnRefs: []string{"fnRef" + nr100},
	}
}

func inSummC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	nr100 := strconv.Itoa(int(n) + 100)
	return esmodel.Content{
		Type:       esmodel.Summary,
		Ordinal:    ordinal,
		Ref:        util.StrPtr("summRef" + nr),
		FmtText:    "<fmt-tag>summary</fmt-tag> text " + nr,
		SearchText: "summary text " + nr,
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef" + nr100},
		WorkId:     workId,
	}
}
func outSummC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	nr100 := strconv.Itoa(int(n) + 100)
	summ := esmodel.Content{
		Type:       esmodel.Summary,
		Id:         "summaryId" + nr,
		Ordinal:    ordinal,
		Ref:        util.StrPtr("summRef" + nr),
		FmtText:    "<fmt-tag>summary</fmt-tag> text " + nr,
		SearchText: "summary text " + nr,
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef" + nr100},
		WorkId:     workId,
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
	wId := "workId"
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
				mockXmlMapper(xm)
				vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(nil, testErr)
			},
		},
		{
			name: "VolumeRepo.Delete fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&esmodel.Volume{
							Works: []esmodel.WorkRef{{Id: wId}},
						}, nil),
					vr.EXPECT().Delete(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
			},
		},
		{
			name: "WorkRepo.Delete fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&esmodel.Volume{
							Works: []esmodel.WorkRef{{Id: wId}},
						}, nil),
					vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
					wr.EXPECT().Delete(gomock.Any(), wId).Return(testErr),
				)
			},
		},
		{
			name: "ContentRepo.DeleteByWorkId fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				gomock.InOrder(
					vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).
						Return(&esmodel.Volume{
							Works: []esmodel.WorkRef{{Id: wId}},
						}, nil),
					vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
					wr.EXPECT().Delete(gomock.Any(), wId).Return(nil),
					cr.EXPECT().DeleteByWorkId(gomock.Any(), wId).
						Return(testErr),
				)
			},
		},
		{
			name: "Insert work fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				wr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(testErr)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Insert heading fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Insert summary and footnote fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Insert paragraphs fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(2),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Update work fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Times(3).Return(nil),
					wr.EXPECT().Update(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Insert volume fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Times(3).Return(nil),
					wr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
					vr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
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

func mockXmlMapper(mapper *mocks.MockXmlMapper) {
	mapper.EXPECT().MapVolume(gomock.Any(), gomock.Any()).Return(&model.Volume{}, errors.Nil())
	mapper.EXPECT().MapWorks(gomock.Any(), gomock.Any()).Return([]model.Work{
		{
			Code:         "c",
			Title:        "t",
			Abbreviation: util.StrPtr("abbr"),
			Year:         util.StrPtr("2024"),
			Sections: []model.Section{{
				Heading: inHead(1),
				Paragraphs: []model.Paragraph{
					inPar(2),
				},
			}},
			Footnotes: []model.Footnote{inFn(3)},
			Summaries: []model.Summary{inSumm(4)},
		},
	}, errors.Nil())
}

func mockDeletion(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, wId string) {
	vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(&esmodel.Volume{
		Works: []esmodel.WorkRef{{Id: wId}},
	}, nil)
	vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	wr.EXPECT().Delete(gomock.Any(), wId).Return(nil)
	cr.EXPECT().DeleteByWorkId(gomock.Any(), wId).Return(nil)
}
