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

	testCases := []struct {
		name      string
		xml       string
		volNum    string
		mockCalls func()
		assert    func(t *testing.T, ctx echo.Context, res *httptest.ResponseRecorder)
	}{
		{
			name:      "Processing error due to invalid volume number",
			xml:       abt1Xml,
			volNum:    "x",
			mockCalls: func() {},
			assert: func(t *testing.T, ctx echo.Context, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
		{
			name:      "Processing error due to invalid xml",
			xml:       "<my-tag>",
			volNum:    "1",
			mockCalls: func() {},
			assert: func(t *testing.T, ctx echo.Context, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
		{
			name:   "Processing success",
			xml:    abt1Xml,
			volNum: "1",
			mockCalls: func() {
				volumeProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			assert: func(t *testing.T, ctx echo.Context, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, ctx.Response().Status)
				assert.Empty(t, res.Body.String())
			},
		},
		{
			name:   "Processing error",
			xml:    abt1Xml,
			volNum: "1",
			mockCalls: func() {
				volumeProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			assert: func(t *testing.T, ctx echo.Context, res *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := []byte(tc.xml)
			// GIVEN
			req := httptest.NewRequest(echo.GET, "/api/write/v1/volumes/"+tc.volNum, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationXML)
			res := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, res)
			ctx.SetParamNames("volumeNumber")
			ctx.SetParamValues(tc.volNum)
			tc.mockCalls()

			// WHEN
			sut.PostVolume(ctx)

			// THEN
			tc.assert(t, ctx, res)
		})
	}
}
