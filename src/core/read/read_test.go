//go:build unit
// +build unit

package read

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	sut := &readProcessorImpl{
		volumeRepo:  volumeRepo,
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
	vol := model.Volume{
		VolumeNumber: 1,
		Title:        "volume title",
		Works: []model.Work{{
			Code:  "workCode",
			Title: "work title",
		}},
	}
	// GIVEN
	volumeRepo.EXPECT().GetAll(gomock.Any()).Return([]model.Volume{vol}, nil)
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

func testProcessFootnotes(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	fn := model.Content{
		Type:       model.Footnote,
		Ref:        util.StrPtr("A121"),
		FmtText:    "formatted text 1",
		SearchText: "search text 1",
		Pages:      []int32{1, 2, 3},
		WorkCode:   workCode,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetFootnotesByWork(gomock.Any(), workCode, []int32{}).
		Return([]model.Content{fn}, nil)
	// WHEN
	res, err := sut.ProcessFootnotes(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, fn, res[0])
}

func testProcessFootnotesError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetFootnotesByWork(gomock.Any(), workCode, []int32{}).Return(nil, e)
	// WHEN
	res, err := sut.ProcessFootnotes(ctx, workCode, []int32{})
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessHeadings(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	head := model.Content{
		Type:       model.Heading,
		FmtText:    "formatted text 2",
		SearchText: "search text 2",
		Pages:      []int32{1, 2, 3},
		FnRefs:     []string{"fn1.2", "fn2.3"},
		WorkCode:   workCode,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetHeadingsByWork(gomock.Any(), workCode, []int32{}).
		Return([]model.Content{head}, nil)
	// WHEN
	res, err := sut.ProcessHeadings(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, head, res[0])
}

func testProcessHeadingsError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetHeadingsByWork(gomock.Any(), workCode, []int32{}).Return(nil, e)
	// WHEN
	res, err := sut.ProcessHeadings(ctx, workCode, []int32{})
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessParagraphs(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	par := model.Content{
		Type:       model.Paragraph,
		Ref:        util.StrPtr("A124"),
		FmtText:    "formatted text 3",
		SearchText: "search text 3",
		Pages:      []int32{4, 5},
		FnRefs:     []string{"fn3.4", "fn4.5"},
		WorkCode:   workCode,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetParagraphsByWork(gomock.Any(), workCode, []int32{}).
		Return([]model.Content{par}, nil)
	// WHEN
	res, err := sut.ProcessParagraphs(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, par, res[0])
}

func testProcessParagraphsError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetParagraphsByWork(gomock.Any(), workCode, []int32{}).Return(nil, e)
	// WHEN
	res, err := sut.ProcessParagraphs(ctx, workCode, []int32{})
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func testProcessSummaries(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	summ := model.Content{
		Type:       model.Summary,
		Ref:        util.StrPtr("A125"),
		FmtText:    "formatted text 5",
		SearchText: "search text 5",
		Pages:      []int32{4, 5},
		WorkCode:   workCode,
	}
	// GIVEN
	contentRepo.EXPECT().
		GetSummariesByWork(gomock.Any(), workCode, []int32{}).
		Return([]model.Content{summ}, nil)
	// WHEN
	res, err := sut.ProcessSummaries(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, summ, res[0])
}

func testProcessSummariesError(t *testing.T, sut *readProcessorImpl, contentRepo *mocks.MockContentRepo, ctx context.Context) {
	workCode := "workCode"
	e := errors.New("test error")
	// GIVEN
	contentRepo.EXPECT().GetSummariesByWork(gomock.Any(), workCode, []int32{}).Return(nil, e)
	// WHEN
	res, err := sut.ProcessSummaries(ctx, workCode, []int32{})
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
