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
		"Read volumes":                    testReadVolumes,
		"Read volumes with error":         testReadVolumesError,
		"Read footnotes":                  testReadFootnotes,
		"Read footnotes with empty code":  testReadFootnotesEmptyCode,
		"Read footnotes with error":       testReadFootnotesError,
		"Read headings":                   testReadHeadings,
		"Read headings with empty code":   testReadHeadingsEmptyCode,
		"Read headings with error":        testReadHeadingsError,
		"Read paragraphs":                 testReadParagraphs,
		"Read paragraphs with empty code": testReadParagraphsEmptyCode,
		"Read paragraphs with error":      testReadParagraphsError,
		"Read summaries":                  testReadSummaries,
		"Read summaries with empty code":  testReadSummariesEmptyCode,
		"Read summaries with error":       testReadSummariesError,
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
		Works: []esmodel.Work{{
			Code:  "A123",
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

func testReadFootnotes(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	fn := esmodel.Content{
		Type:       esmodel.Footnote,
		Ref:        util.StrPtr("A121"),
		FmtText:    "formatted text 1",
		SearchText: "search text 1",
		Pages:      []int32{1, 2, 3},
		WorkCode:   workCode,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/footnotes", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().
		ProcessFootnotes(gomock.Any(), workCode).
		Return([]esmodel.Content{fn}, nil)
	// WHEN
	sut.ReadFootnotes(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), fn.FmtText)
}

func testReadFootnotesEmptyCode(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	// WHEN
	sut.ReadFootnotes(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work code")

}

func testReadFootnotesError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/footnotes", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().ProcessFootnotes(gomock.Any(), workCode).Return(nil, e)
	// WHEN
	sut.ReadFootnotes(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func testReadHeadings(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	head := esmodel.Content{
		Type:       esmodel.Heading,
		FmtText:    "formatted text 2",
		SearchText: "search text 2",
		Pages:      []int32{1, 2, 3},
		FnRefs:     []string{"fn1.2", "fn2.3"},
		WorkCode:   workCode,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/headings", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().
		ProcessHeadings(gomock.Any(), workCode).
		Return([]esmodel.Content{head}, nil)
	// WHEN
	sut.ReadHeadings(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), head.FmtText)
}

func testReadHeadingsEmptyCode(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	// WHEN
	sut.ReadHeadings(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work code")

}

func testReadHeadingsError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/headings", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().ProcessHeadings(gomock.Any(), workCode).Return(nil, e)
	// WHEN
	sut.ReadHeadings(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func testReadParagraphs(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	par := esmodel.Content{
		Type:       esmodel.Paragraph,
		Ref:        util.StrPtr("A124"),
		FmtText:    "formatted text 3",
		SearchText: "search text 3",
		Pages:      []int32{4, 5},
		FnRefs:     []string{"fn3.4", "fn4.5"},
		WorkCode:   workCode,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().
		ProcessParagraphs(gomock.Any(), workCode).
		Return([]esmodel.Content{par}, nil)
	// WHEN
	sut.ReadParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), par.FmtText)
}

func testReadParagraphsEmptyCode(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	// WHEN
	sut.ReadParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work code")

}

func testReadParagraphsError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().ProcessParagraphs(gomock.Any(), workCode).Return(nil, e)
	// WHEN
	sut.ReadParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func testReadSummaries(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	summ := esmodel.Content{
		Type:       esmodel.Summary,
		Ref:        util.StrPtr("A125"),
		FmtText:    "formatted text 5",
		SearchText: "search text 5",
		Pages:      []int32{4, 5},
		WorkCode:   workCode,
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/summaries", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().
		ProcessSummaries(gomock.Any(), workCode).
		Return([]esmodel.Content{summ}, nil)
	// WHEN
	sut.ReadSummaries(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), summ.FmtText)
}

func testReadSummariesEmptyCode(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := ""
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode, nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	// WHEN
	sut.ReadSummaries(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "empty work code")

}

func testReadSummariesError(t *testing.T, sut *readHandlerImpl, readProcessor *mocks.MockReadProcessor) {
	workCode := "A123"
	e := errors.New("test error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/"+workCode+"/summaries", nil)
	res := httptest.NewRecorder()
	ctx := createCtxWithWorkCode(req, res, workCode)
	readProcessor.EXPECT().ProcessSummaries(gomock.Any(), workCode).Return(nil, e)
	// WHEN
	sut.ReadSummaries(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "message")
}

func createCtxWithWorkCode(req *http.Request, res *httptest.ResponseRecorder, workCode string) echo.Context {
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workCode")
	ctx.SetParamValues(workCode)
	return ctx
}
