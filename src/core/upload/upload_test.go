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
		Abbreviation: util.ToStrPtr("abbrev"),
		Title:        "work title",
		Year:         util.ToStrPtr("1785"),
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
			inFn(13),
			inFn(14),
		},
		Summaries: []model.Summary{
			inSumm(15),
			inSumm(16),
		},
	}

	// GIVEN
	// mapping
	xmlMapper.EXPECT().
		MapVolume(gomock.Eq(volNr), gomock.Any()).
		Return(vol, errors.NilError())
	xmlMapper.EXPECT().
		MapWorks(gomock.Eq(volNr), gomock.Any()).
		Return([]model.Work{work}, errors.NilError())

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
				Heading:    cId(1),
				Paragraphs: []string{cId(2), cId(3)},
				Sections: []esmodel.Section{
					{
						Heading:    cId(4),
						Paragraphs: []string{cId(5), cId(6)},
						Sections:   []esmodel.Section{},
					},
					{
						Heading:    cId(7),
						Paragraphs: []string{cId(8), cId(9)},
						Sections:   []esmodel.Section{},
					},
				},
			},
			{
				Heading:    cId(10),
				Paragraphs: []string{cId(11), cId(12)},
				Sections:   []esmodel.Section{},
			},
		},
	})).Return(nil)

	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq(
			[]esmodel.Content{inFnC(13, wId), inFnC(14, wId)},
		)).
		Do(func(ctx any, c []esmodel.Content) {
			c[0] = outFnC(13, wId)
			c[1] = outFnC(14, wId)
		}).Return(nil)
	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq(
			[]esmodel.Content{inSummC(15, wId), inSummC(16, wId)},
		)).
		Do(func(ctx any, c []esmodel.Content) {
			c[0] = outSummC(15, wId)
			c[1] = outSummC(16, wId)
		}).Return(nil)
	volumeRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

	// WHEN
	err := sut.Process(ctx, volNr, "xml")

	// THEN
	assert.False(t, err.HasError)
}

func expectHeading(contentRepo *dbMocks.MockContentRepo, n int32, wId string) {
	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq(
			[]esmodel.Content{inHeadC(n, wId)},
		)).
		DoAndReturn(func(ctx context.Context, c []esmodel.Content) {
			c[0] = outHeadC(n, wId)
		}).Return(nil)
}

func expectParagraphs(contentRepo *dbMocks.MockContentRepo, n1 int32, n2 int32, wId string) {
	contentRepo.EXPECT().
		Insert(gomock.Any(), gomock.Eq(
			[]esmodel.Content{inParC(n1, wId), inParC(n2, wId)},
		)).
		Do(func(ctx any, c []esmodel.Content) {
			c[0] = outParC(n1, wId)
			c[1] = outParC(n2, wId)
		}).Return(nil)
}

func inHead(n int32) model.Heading {
	nr := strconv.Itoa(int(n))
	return model.Heading{
		Text:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText: "toc text " + nr,
		Pages:   []int32{n, n + 100},
		FnRefs:  []string{"fnRef " + nr},
	}
}

func inHeadC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Heading,
		FmtText:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText:    util.ToStrPtr("toc text " + nr),
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		WorkId:     workId,
	}
}

func outHeadC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Id:         cId(n),
		Type:       esmodel.Heading,
		FmtText:    "<fmt-tag>heading</fmt-tag> text " + nr,
		TocText:    util.ToStrPtr("toc text " + nr),
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		WorkId:     workId,
	}
}

func inPar(n int32) model.Paragraph {
	nr := strconv.Itoa(int(n))
	return model.Paragraph{
		Text:       "<fmt-tag>paragraph</fmt-tag> text " + nr,
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		SummaryRef: util.ToStrPtr("summRef " + nr),
	}
}

func inParC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Paragraph,
		FmtText:    "<fmt-tag>paragraph</fmt-tag> text " + nr,
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		SummaryRef: util.ToStrPtr("summRef " + nr),
		WorkId:     workId,
	}
}

func outParC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Id:         cId(n),
		Type:       esmodel.Paragraph,
		FmtText:    "<fmt-tag>paragraph</fmt-tag> text " + nr,
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		SummaryRef: util.ToStrPtr("summRef " + nr),
		WorkId:     workId,
	}
}

func inFn(n int32) model.Footnote {
	nr := strconv.Itoa(int(n))
	return model.Footnote{
		Text:  "<fmt-tag>footnote</fmt-tag> text " + nr,
		Ref:   "footnote ref " + nr,
		Pages: []int32{n, n + 100},
	}
}

func inFnC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Footnote,
		Ref:        util.ToStrPtr("footnote ref " + nr),
		FmtText:    "<fmt-tag>footnote</fmt-tag> text " + nr,
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		WorkId:     workId,
	}
}

func outFnC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Id:         cId(n),
		Type:       esmodel.Footnote,
		Ref:        util.ToStrPtr("footnote ref " + nr),
		FmtText:    "<fmt-tag>footnote</fmt-tag> text " + nr,
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		WorkId:     workId,
	}
}

func inSumm(n int32) model.Summary {
	nr := strconv.Itoa(int(n))
	return model.Summary{
		Text:   "<fmt-tag>summary</fmt-tag> text " + nr,
		Ref:    "summary ref " + nr,
		Pages:  []int32{n, n + 100},
		FnRefs: []string{"fnRef " + nr},
	}
}

func inSummC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Type:       esmodel.Summary,
		Ref:        util.ToStrPtr("summary ref " + nr),
		FmtText:    "<fmt-tag>summary</fmt-tag> text " + nr,
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		WorkId:     workId,
	}
}
func outSummC(n int32, workId string) esmodel.Content {
	nr := strconv.Itoa(int(n))
	return esmodel.Content{
		Id:         cId(n),
		Type:       esmodel.Summary,
		Ref:        util.ToStrPtr("summary ref " + nr),
		FmtText:    "<fmt-tag>summary</fmt-tag> text " + nr,
		SearchText: "", // TODO
		Pages:      []int32{n, n + 100},
		FnRefs:     []string{"fnRef " + nr},
		WorkId:     workId,
	}
}

func cId(n int32) string {
	nr := strconv.Itoa(int(n))
	return "contentId" + nr
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
					Return(nil, errors.NewError(nil, testErr))
			},
		},
		{
			name: "MapWorks fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				gomock.InOrder(
					xm.EXPECT().MapVolume(gomock.Any(), gomock.Any()).
						Return(&model.Volume{}, errors.NilError()),
					xm.EXPECT().MapWorks(gomock.Any(), gomock.Any()).
						Return(nil, errors.NewError(nil, testErr)),
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
			name: "Insert paragraphs fails",
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
						Times(2).Return(nil),
					wr.EXPECT().Update(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Insert footnotes fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Times(2).Return(nil),
					wr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Return(testErr),
				)
				mockDeletion(vr, wr, cr, wId)
			},
		},
		{
			name: "Insert summaries fails",
			mockSetup: func(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, xm *mocks.MockXmlMapper) {
				mockXmlMapper(xm)
				mockDeletion(vr, wr, cr, wId)
				gomock.InOrder(
					wr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, w *esmodel.Work) {
							w.Id = wId
						}).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Times(2).Return(nil),
					wr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
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
						Times(2).Return(nil),
					wr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil),
					cr.EXPECT().Insert(gomock.Any(), gomock.Any()).
						Times(2).Return(nil),
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
	mapper.EXPECT().MapVolume(gomock.Any(), gomock.Any()).Return(&model.Volume{}, errors.NilError())
	mapper.EXPECT().MapWorks(gomock.Any(), gomock.Any()).Return([]model.Work{
		{
			Code:         "c",
			Title:        "t",
			Abbreviation: util.ToStrPtr("abbr"),
			Year:         util.ToStrPtr("2024"),
			Sections: []model.Section{{
				Heading: inHead(1),
				Paragraphs: []model.Paragraph{
					inPar(2),
				},
			}},
			Footnotes: []model.Footnote{inFn(3)},
			Summaries: []model.Summary{inSumm(4)},
		},
	}, errors.NilError())
}

func mockDeletion(vr *dbMocks.MockVolumeRepo, wr *dbMocks.MockWorkRepo, cr *dbMocks.MockContentRepo, wId string) {
	vr.EXPECT().GetByVolumeNumber(gomock.Any(), gomock.Any()).Return(&esmodel.Volume{
		Works: []esmodel.WorkRef{{Id: wId}},
	}, nil)
	vr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	wr.EXPECT().Delete(gomock.Any(), wId).Return(nil)
	cr.EXPECT().DeleteByWorkId(gomock.Any(), wId).Return(nil)
}
