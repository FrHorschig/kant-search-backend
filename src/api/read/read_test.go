//go:build unit
// +build unit

package read

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/read/mocks"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestReadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	readProcessor := mocks.NewMockReadProcessor(ctrl)
	sut := &readHandlerImpl{
		readProcessor: readProcessor,
	}

	for scenario, fn := range map[string]func(*testing.T, *readHandlerImpl, *mocks.MockReadProcessor){
		"Read volumes":                  testReadVolumes,
		"Read volumes with error":       testReadVolumesError,
		"Read work":                     testReadWork,
		"Read work not found":           testReadWorkNotFound,
		"Read work with empty id":       testReadWorkEmptyId,
		"Read work with error":          testReadWorkError,
		"Read footnotes":                testReadFootnotes,
		"Read footnotes with empty id":  testReadFootnotesEmptyId,
		"Read footnotes with error":     testReadFootnotesError,
		"Read headings":                 testReadHeadings,
		"Read headings with empty id":   testReadHeadingsEmptyId,
		"Read headings with error":      testReadHeadingsError,
		"Read paragraphs":               testReadParagraphs,
		"Read paragraphs with empty id": testReadParagraphsEmptyId,
		"Read paragraphs with error":    testReadParagraphsError,
		"Read summaries":                testReadSummaries,
		"Read summaries with empty id":  testReadSummariesEmptyId,
		"Read summaries with error":     testReadSummariesError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, readProcessor)
		})
	}
}

func testReadVolumes(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	vol := esmodel.Volume{
		VolumeNumber: 1,
		Section:      2,
		Title:        "volume title",
		Works: []esmodel.WorkRef{{
			Id:    "A123",
			Code:  "code",
			Title: "work title",
		}},
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/volumes", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	readProcessor.EXPECT().ProcessVolumes(gomock.Any()).Return([]esmodel.Volume{vol}, nil)
	// WHEN
	sut.ReadVolumes(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), vol.Title)
}

func testReadVolumesError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/volumes", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	readProcessor.EXPECT().ProcessVolumes(gomock.Any()).Return(nil, e)
	// WHEN
	sut.ReadVolumes(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")

}

func testReadWork(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	work := esmodel.Work{
		Id:           workId,
		Code:         "GMS",
		Abbreviation: util.ToStrPtr("GMS"),
		Title:        "Grundlegung zur Metaphysik der Sitten",
		Year:         util.ToStrPtr("1785"),
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessWork(gomock.Any(), workId).Return(&work, nil)
	// WHEN
	sut.ReadWork(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), work.Title)
}

func testReadWorkEmptyId(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	// WHEN
	sut.ReadWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work ID")

}

func testReadWorkNotFound(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessWork(gomock.Any(), workId).Return(nil, nil)
	// WHEN
	sut.ReadWork(ctx)
	// THEN
	assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")

}

func testReadWorkError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessWork(gomock.Any(), workId).Return(nil, e)
	// WHEN
	sut.ReadWork(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")

}

func testReadFootnotes(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	fn := esmodel.Content{
		Type:       esmodel.Footnote,
		Ref:        util.ToStrPtr("A121"),
		FmtText:    "formatted text 1",
		SearchText: "search text 1",
		Pages:      []int32{1, 2, 3},
		WorkId:     workId,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/footnotes", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().
		ProcessFootnotes(gomock.Any(), workId).
		Return([]esmodel.Content{fn}, nil)
	// WHEN
	sut.ReadFootnotes(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), fn.FmtText)
}

func testReadFootnotesEmptyId(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	// WHEN
	sut.ReadFootnotes(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work ID")

}

func testReadFootnotesError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/footnotes", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessFootnotes(gomock.Any(), workId).Return(nil, e)
	// WHEN
	sut.ReadFootnotes(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func testReadHeadings(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	head := esmodel.Content{
		Type:       esmodel.Heading,
		FmtText:    "formatted text 2",
		SearchText: "search text 2",
		Pages:      []int32{1, 2, 3},
		FnRefs:     []string{"fn1.2", "fn2.3"},
		WorkId:     workId,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/headings", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().
		ProcessHeadings(gomock.Any(), workId).
		Return([]esmodel.Content{head}, nil)
	// WHEN
	sut.ReadHeadings(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), head.FmtText)
}

func testReadHeadingsEmptyId(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	// WHEN
	sut.ReadHeadings(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work ID")

}

func testReadHeadingsError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/headings", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessHeadings(gomock.Any(), workId).Return(nil, e)
	// WHEN
	sut.ReadHeadings(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func testReadParagraphs(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
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
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().
		ProcessParagraphs(gomock.Any(), workId).
		Return([]esmodel.Content{par}, nil)
	// WHEN
	sut.ReadParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), par.FmtText)
}

func testReadParagraphsEmptyId(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	// WHEN
	sut.ReadParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work ID")

}

func testReadParagraphsError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessParagraphs(gomock.Any(), workId).Return(nil, e)
	// WHEN
	sut.ReadParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func testReadSummaries(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	summ := esmodel.Content{
		Type:       esmodel.Summary,
		Ref:        util.ToStrPtr("A125"),
		FmtText:    "formatted text 5",
		SearchText: "search text 5",
		Pages:      []int32{4, 5},
		WorkId:     workId,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/summaries", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().
		ProcessSummaries(gomock.Any(), workId).
		Return([]esmodel.Content{summ}, nil)
	// WHEN
	sut.ReadSummaries(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), summ.FmtText)
}

func testReadSummariesEmptyId(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	// WHEN
	sut.ReadSummaries(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work ID")

}

func testReadSummariesError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workId := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workId+"/summaries", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkId(req, res, workId)
	readProcessor.EXPECT().ProcessSummaries(gomock.Any(), workId).Return(nil, e)
	// WHEN
	sut.ReadSummaries(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func createCtxWithWorkId(req *http.Request, res *httptest.ResponseRecorder, workId string) echo.Context {
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues(workId)
	return ctx
}
