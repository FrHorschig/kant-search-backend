//go:build unit
// +build unit

package upload

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	procMocks "github.com/frhorschig/kant-search-backend/core/upload/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const abt1Xml = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<kant_abt1>
  <band nr="01">
    <titel>1</titel>
    <hauptteil>Hauptteil</hauptteil>
  </band>
</kant_abt1>
`

func TestUploadHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	volumeProcessor := procMocks.NewMockVolumeUploadProcessor(ctrl)
	sut := NewUploadHandler(volumeProcessor).(*uploadHandlerImpl)

	for scenario, fn := range map[string]func(t *testing.T, sut *uploadHandlerImpl, volumeProcessor *procMocks.MockVolumeUploadProcessor){
		"PostVolumes invalid volume number error": testPostVolumesInvalidVolumeNumber,
		"PostVolumes error reading body":          testPostVolumesErrorReadingBody,
		"PostVolumes error processing abt1":       testPostVolumesErrorProcessingAbt1,
		"PostVolumes success processing abt1":     testPostVolumesSuccessProcessingAbt1,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, volumeProcessor)
		})
	}
}

func testPostVolumesInvalidVolumeNumber(t *testing.T, sut *uploadHandlerImpl, volumeProcessor *procMocks.MockVolumeUploadProcessor) {
	body := []byte("text")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/write/v1/volumes/x", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("volumeNumber")
	ctx.SetParamValues("x")
	// WHEN
	sut.PostVolume(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

func testPostVolumesErrorReadingBody(t *testing.T, sut *uploadHandlerImpl, volumeProcessor *procMocks.MockVolumeUploadProcessor) {
	body := []byte("")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/write/v1/volumes/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("volumeNumber")
	ctx.SetParamValues("1")
	// WHEN
	sut.PostVolume(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
}

func testPostVolumesErrorProcessingAbt1(t *testing.T, sut *uploadHandlerImpl, volumeProcessor *procMocks.MockVolumeUploadProcessor) {
	body := []byte(abt1Xml)
	processErr := fmt.Errorf("error")
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/write/v1/volumes/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("volumeNumber")
	ctx.SetParamValues("1")
	volumeProcessor.EXPECT().ProcessAbt1(gomock.Any(), gomock.Any(), gomock.Any()).Return(processErr)
	// WHEN
	sut.PostVolume(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
}

func testPostVolumesSuccessProcessingAbt1(t *testing.T, sut *uploadHandlerImpl, volumeProcessor *procMocks.MockVolumeUploadProcessor) {
	body := []byte(abt1Xml)
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/write/v1/volumes/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("volumeNumber")
	ctx.SetParamValues("1")
	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	volumeProcessor.EXPECT().ProcessAbt1(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	// WHEN
	sut.PostVolume(ctx)
	// THEN
	assert.Equal(t, http.StatusCreated, ctx.Response().Status)
	assert.Empty(t, res.Body.String())
}
