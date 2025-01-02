//go:build unit
// +build unit

package upload

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/errors"
	procMocks "github.com/frhorschig/kant-search-backend/core/upload/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workProcessor := procMocks.NewMockWorkUploadProcessor(ctrl)
	volumeProcessor := procMocks.NewMockVolumeUploadProcessor(ctrl)
	sut := NewUploadHandler(workProcessor, volumeProcessor).(*uploadHandlerImpl)

	for scenario, fn := range map[string]func(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor){
		"PostWorks invalid workId error": testPostWorksInvalidWorkId,
		"PostWorks empty text error":     testPostWorksEmptyText,
		"PostWorks process error":        testPostWorksProcessError,
		"PostWorks success":              testPostWorksSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, workProcessor)
		})
	}
}

func testPostWorksInvalidWorkId(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body := []byte("text")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/x", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("x")
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

func testPostWorksEmptyText(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body := []byte("")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

func testPostWorksProcessError(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body := []byte("text")
	processErr := &errors.Error{
		Msg:    errors.GO_ERR,
		Params: []string{"detail"},
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/works/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(processErr)
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
}

func testPostWorksParseError(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body := []byte("text")
	parseErr := &errors.Error{
		Msg:    errors.WRONG_STARTING_CHAR,
		Params: []string{string("detail")},
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/works/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(parseErr)
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

func testPostWorksSuccess(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body := []byte("text")
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/works/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusCreated, ctx.Response().Status)
	assert.Empty(t, res.Body.String())
}
