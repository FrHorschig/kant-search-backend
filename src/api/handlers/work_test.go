package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FrHorschig/kant-search-api/models"
	procMocks "github.com/FrHorschig/kant-search-backend/core/upload/mocks"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository/mocks"
	"github.com/FrHorschig/kant-search-backend/util"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWorkHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	volumeRepo := mocks.NewMockVolumeRepo(ctrl)
	workRepo := mocks.NewMockWorkRepo(ctrl)
	workProcessor := procMocks.NewMockWorkUploadProcessor(ctrl)
	sut := &workHandlerImpl{
		volumeRepo:    volumeRepo,
		workRepo:      workRepo,
		workProcessor: workProcessor,
	}

	for scenario, fn := range map[string]func(t *testing.T, sut *workHandlerImpl, volumeRepo *mocks.MockVolumeRepo){
		"GetVolumes database error": testSelectVolumesDatabaseError,
		"GetVolumes empty result":   testSelectVolumesEmptyResult,
		"GetVolumes single result":  testSelectVolumesSingleResult,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, volumeRepo)
		})
	}

	for scenario, fn := range map[string]func(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo){
		"GetWorks database error": testSelectWorksDatabaseError,
		"GetWorks empty result":   testSelectWorksEmptyResult,
		"GetWorks single result":  testSelectWorksSingleResult,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, workRepo)
		})
	}

	for scenario, fn := range map[string]func(t *testing.T, sut *workHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor){
		"PostWorks bind error":        testPostWorksBindError,
		"PostWorks zero workId error": testPostWorksZeroWorkId,
		"PostWorks empty text error":  testPostWorksEmptyText,
		"PostWorks process error":     testPostWorksProcessError,
		"PostWorks success":           testPostWorksSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, workProcessor)
		})
	}
}

func testSelectVolumesDatabaseError(t *testing.T, sut *workHandlerImpl, volumeRepo *mocks.MockVolumeRepo) {
	volumes := []model.Volume{}
	err := errors.New("database error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/volumes", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	volumeRepo.EXPECT().SelectAll(gomock.Any()).Return(volumes, err)
	// WHEN
	sut.GetVolumes(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSelectVolumesEmptyResult(t *testing.T, sut *workHandlerImpl, volumeRepo *mocks.MockVolumeRepo) {
	volumes := []model.Volume{}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/volumes", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	volumeRepo.EXPECT().SelectAll(gomock.Any()).Return(volumes, nil)
	// WHEN
	sut.GetVolumes(ctx)
	// THEN
	assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSelectVolumesSingleResult(t *testing.T, sut *workHandlerImpl, volumeRepo *mocks.MockVolumeRepo) {
	volumes := []model.Volume{{
		Id:      1,
		Title:   "title",
		Section: 1,
	}}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/volumes", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	volumeRepo.EXPECT().SelectAll(gomock.Any()).Return(volumes, nil)
	// WHEN
	sut.GetVolumes(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "id")
	assert.Contains(t, res.Body.String(), "title")
	assert.Contains(t, res.Body.String(), "section")
}

func testSelectWorksDatabaseError(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo) {
	works := []model.Work{}
	err := errors.New("database error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	workRepo.EXPECT().SelectAll(gomock.Any()).Return(works, err)
	// WHEN
	sut.GetWorks(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSelectWorksEmptyResult(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo) {
	works := []model.Work{}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	workRepo.EXPECT().SelectAll(gomock.Any()).Return(works, nil)
	// WHEN
	sut.GetWorks(ctx)
	// THEN
	assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSelectWorksSingleResult(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo) {
	works := []model.Work{{
		Id:           1,
		Title:        "title",
		Abbreviation: util.ToStrPtr("abbrev"),
		Ordinal:      1,
		Year:         util.ToStrPtr("1785"),
		Volume:       1,
	}}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	workRepo.EXPECT().SelectAll(gomock.Any()).Return(works, nil)
	// WHEN
	sut.GetWorks(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "id")
	assert.Contains(t, res.Body.String(), "title")
	assert.Contains(t, res.Body.String(), "abbreviation")
	assert.Contains(t, res.Body.String(), "ordinal")
	assert.Contains(t, res.Body.String(), "year")
	assert.Contains(t, res.Body.String(), "volume")
}

func testPostWorksBindError(t *testing.T, sut *workHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body, err := json.Marshal(models.Volume{Id: 1, Title: "title", Section: 1})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.Request().Body = io.NopCloser(bytes.NewReader(body))
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testPostWorksZeroWorkId(t *testing.T, sut *workHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body, err := json.Marshal(models.WorkUpload{WorkId: 0, Text: "text"})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.Request().Body = io.NopCloser(bytes.NewReader(body))
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testPostWorksEmptyText(t *testing.T, sut *workHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body, err := json.Marshal(models.WorkUpload{WorkId: 1, Text: ""})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works", bytes.NewReader(body))
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testPostWorksProcessError(t *testing.T, sut *workHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body, err := json.Marshal(models.WorkUpload{WorkId: 1, Text: "text"})
	if err != nil {
		t.Fatal(err)
	}
	processErr := errors.New("process error")
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/works", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(processErr)
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testPostWorksSuccess(t *testing.T, sut *workHandlerImpl, workProcessor *procMocks.MockWorkUploadProcessor) {
	body, err := json.Marshal(models.WorkUpload{WorkId: 1, Text: "text"})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/works", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	workProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil)
	// WHEN
	sut.PostWork(ctx)
	// THEN
	assert.Equal(t, http.StatusCreated, ctx.Response().Status)
	assert.Empty(t, res.Body.String())
}