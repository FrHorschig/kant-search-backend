//go:build unit
// +build unit

package read

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
	workRepo := dbMocks.NewMockWorkRepo(ctrl)
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	sut := &readProcessorImpl{
		volumeRepo:  volumeRepo,
		workRepo:    workRepo,
		contentRepo: contentRepo,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for scenario, fn := range map[string]func(*testing.T, *readProcessorImpl, *mocks.MockVolumeRepo, context.Context){
		"Process volumes":            testProcessVolumes,
		"Process volumes with error": testProcessVolumesError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, volumeRepo, ctx)
		})
	}
	for scenario, fn := range map[string]func(*testing.T, *readProcessorImpl, *mocks.MockWorkRepo, context.Context){
		"Process work":            testProcessWork,
		"Process work with error": testProcessWorkError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, workRepo, ctx)
		})
	}
	for scenario, fn := range map[string]func(*testing.T, *readProcessorImpl, *mocks.MockContentRepo, context.Context){
		"Process footnotes":             testProcessFootnotes,
		"Process footnotes with error":  testProcessFootnotesError,
		"Process headings":              testProcessHeadings,
		"Process headings with error":   testProcessHeadingsError,
		"Process paragraphs":            testProcessParagraphs,
		"Process paragraphs with error": testProcessParagraphsError,
		"Process summaries":             testProcessSummaries,
		"Process summaries with error":  testProcessSummariesError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, contentRepo, ctx)
		})
	}
}

func testProcessVolumes(t *testing.T, sut *readProcessorImpl, volumeRepo *mocks.MockVolumeRepo, ctx context.Context) {
	vol := esmodel.Volume{
		VolumeNumber: 1,
		Section:      2,
		Title:        "volume title",
		Works: []esmodel.WorkRef{{
			Id:    "workId",
			Code:  "code",
			Title: "work title",
		}},
	}
	// GIVEN
	volumeRepo.EXPECT().GetAll(gomock.Any()).Return([]esmodel.Volume{vol}, nil)
	// WHEN
	res, err := sut.ProcessVolumes(ctx)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, vol, res[0])
}

func testProcessVolumesError(t *testing.T, sut *readProcessorImpl, volumeRepo *mocks.MockVolumeRepo, ctx context.Context) {
	e := errors.New("test error")
	// GIVEN
	volumeRepo.EXPECT().GetAll(gomock.Any()).Return(nil, e)
	// WHEN
	res, err := sut.ProcessVolumes(ctx)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessWork(t *testing.T, sut *readProcessorImpl, workRepo *mocks.MockWorkRepo, ctx context.Context) {
	workId := "workId"
	work := esmodel.Work{
		Id:           workId,
		Code:         "GMS",
		Abbreviation: util.ToStrPtr("GMS"),
		Title:        "Grundlegung zur Metaphysik der Sitten",
		Year:         util.ToStrPtr("1785"),
	}
	// GIVEN
	workRepo.EXPECT().Get(gomock.Any(), workId).Return(&work, nil)
	// WHEN
	res, err := sut.ProcessWork(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Equal(t, work, *res)
}

func testProcessWorkError(t *testing.T, sut *readProcessorImpl, workRepo *mocks.MockWorkRepo, ctx context.Context) {
	workId := "workId"
	e := errors.New("test error")
	// GIVEN
	workRepo.EXPECT().Get(gomock.Any(), workId).Return(nil, e)
	// WHEN
	res, err := sut.ProcessWork(ctx, workId)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessFootnotes(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	fn := esmodel.Content{
		Type:       esmodel.Footnote,
		Ref:        util.ToStrPtr("A121"),
		FmtText:    "formatted text 1",
		SearchText: "search text 1",
		Pages:      []int32{1, 2, 3},
		WorkId:     workId,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetFootnotesByWorkId(gomock.Any(), workId).
		Return([]esmodel.Content{fn}, nil)
	// WHEN
	res, err := sut.ProcessFootnotes(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, fn, res[0])
}

func testProcessFootnotesError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetFootnotesByWorkId(gomock.Any(), workId).Return(nil, e)
	// WHEN
	res, err := sut.ProcessFootnotes(ctx, workId)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessHeadings(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	head := esmodel.Content{
		Type:       esmodel.Heading,
		FmtText:    "formatted text 2",
		SearchText: "search text 2",
		Pages:      []int32{1, 2, 3},
		FnRefs:     []string{"fn1.2", "fn2.3"},
		WorkId:     workId,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetHeadingsByWorkId(gomock.Any(), workId).
		Return([]esmodel.Content{head}, nil)
	// WHEN
	res, err := sut.ProcessHeadings(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, head, res[0])
}

func testProcessHeadingsError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetHeadingsByWorkId(gomock.Any(), workId).Return(nil, e)
	// WHEN
	res, err := sut.ProcessHeadings(ctx, workId)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessParagraphs(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	par := esmodel.Content{
		Type:       esmodel.Paragraph,
		Ref:        util.ToStrPtr("A124"),
		FmtText:    "formatted text 3",
		SearchText: "search text 3",
		Pages:      []int32{4, 5},
		FnRefs:     []string{"fn3.4", "fn4.5"},
		WorkId:     workId,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetParagraphsByWorkId(gomock.Any(), workId).
		Return([]esmodel.Content{par}, nil)
	// WHEN
	res, err := sut.ProcessParagraphs(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, par, res[0])
}

func testProcessParagraphsError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetParagraphsByWorkId(gomock.Any(), workId).Return(nil, e)
	// WHEN
	res, err := sut.ProcessParagraphs(ctx, workId)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessSummaries(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	summ := esmodel.Content{
		Type:       esmodel.Summary,
		Ref:        util.ToStrPtr("A125"),
		FmtText:    "formatted text 5",
		SearchText: "search text 5",
		Pages:      []int32{4, 5},
		WorkId:     workId,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetSummariesByWorkId(gomock.Any(), workId).
		Return([]esmodel.Content{summ}, nil)
	// WHEN
	res, err := sut.ProcessSummaries(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, summ, res[0])
}

func testProcessSummariesError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workId := "workId"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetSummariesByWorkId(gomock.Any(), workId).Return(nil, e)
	// WHEN
	res, err := sut.ProcessSummaries(ctx, workId)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
