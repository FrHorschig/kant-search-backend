//go:build unit
// +build unit

package upload

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	procMocks "github.com/frhorschig/kant-search-backend/core/upload/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	volumeProcessor := procMocks.NewMockVolumeUploadProcessor(ctrl)
	sut := NewUploadHandler(volumeProcessor).(*uploadHandlerImpl)

	for scenario, fn := range map[string]func(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockVolumeUploadProcessor){
		"PostVolumes invalid workId error": testPostVolumesInvalidWorkId,
		"PostVolumes empty text error":     testPostVolumesEmptyText,
		// "PostVolumes process error":        testPostVolumesProcessError,
		// "PostVolumes success":              testPostVolumesSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, volumeProcessor)
		})
	}
}

func testPostVolumesInvalidWorkId(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockVolumeUploadProcessor) {
	body := []byte("text")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/x", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("x")
	// WHEN
	sut.PostVolume(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

func testPostVolumesEmptyText(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockVolumeUploadProcessor) {
	body := []byte("")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	// WHEN
	sut.PostVolume(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

// func testPostVolumesProcessError(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockVolumeUploadProcessor) {
// 	body := []byte("text")
// 	processErr := &errors.Error{
// 		Msg:    errors.GO_ERR,
// 		Params: []string{"detail"},
// 	}
// 	// GIVEN
// 	req := httptest.NewRequest(echo.POST, "/api/v1/works/1", bytes.NewReader(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
// 	res := httptest.NewRecorder()
// 	ctx := echo.New().NewContext(req, res)
// 	ctx.SetParamNames("workId")
// 	ctx.SetParamValues("1")
// 	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(processErr)
// 	// WHEN
// 	sut.PostVolume(ctx)
// 	// THEN
// 	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
// }

// func testPostVolumesParseError(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockVolumeUploadProcessor) {
// 	body := []byte("text")
// 	parseErr := &errors.Error{
// 		Msg:    errors.WRONG_STARTING_CHAR,
// 		Params: []string{string("detail")},
// 	}
// 	// GIVEN
// 	req := httptest.NewRequest(echo.POST, "/api/v1/works/1", bytes.NewReader(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
// 	res := httptest.NewRecorder()
// 	ctx := echo.New().NewContext(req, res)
// 	ctx.SetParamNames("workId")
// 	ctx.SetParamValues("1")
// 	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(parseErr)
// 	// WHEN
// 	sut.PostVolume(ctx)
// 	// THEN
// 	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
// }

// func testPostVolumesSuccess(t *testing.T, sut *uploadHandlerImpl, workProcessor *procMocks.MockVolumeUploadProcessor) {
// 	body := []byte("text")
// 	// GIVEN
// 	req := httptest.NewRequest(echo.POST, "/api/v1/works/1", bytes.NewReader(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
// 	res := httptest.NewRecorder()
// 	ctx := echo.New().NewContext(req, res)
// 	ctx.SetParamNames("workId")
// 	ctx.SetParamValues("1")
// 	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
// 	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
// 	// WHEN
// 	sut.PostVolume(ctx)
// 	// THEN
// 	assert.Equal(t, http.StatusCreated, ctx.Response().Status)
// 	assert.Empty(t, res.Body.String())
// }
