//go:build unit
// +build unit

package read

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/api/common/util"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWorkHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	volumeRepo := mocks.NewMockVolumeRepo(ctrl)
	workRepo := mocks.NewMockWorkRepo(ctrl)
	sut := NewWorkHandler(volumeRepo, workRepo).(*workHandlerImpl)

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
	assertErrorResponse(t, res, string(models.INTERNAL_SERVER_ERROR))
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
	assertErrorResponse(t, res, string(models.NOT_FOUND_VOLUMES))
}

func testSelectVolumesSingleResult(t *testing.T, sut *workHandlerImpl, volumeRepo *mocks.MockVolumeRepo) {
	volumes := []model.Volume{{
		Id:      1,
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
	assert.Contains(t, res.Body.String(), "section")
}

func testSelectWorksDatabaseError(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo) {
	works := []model.Work{}
	err := errors.New("database error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/1", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	workRepo.EXPECT().SelectAll(gomock.Any()).Return(works, err)
	// WHEN
	sut.GetWorks(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res, string(models.INTERNAL_SERVER_ERROR))
}

func testSelectWorksEmptyResult(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo) {
	works := []model.Work{}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/1", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	workRepo.EXPECT().SelectAll(gomock.Any()).Return(works, nil)
	// WHEN
	sut.GetWorks(ctx)
	// THEN
	assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	assertErrorResponse(t, res, string(models.NOT_FOUND_WORKS))
}

func testSelectWorksSingleResult(t *testing.T, sut *workHandlerImpl, workRepo *mocks.MockWorkRepo) {
	works := []model.Work{{
		Id:           1,
		Code:         "code",
		Abbreviation: util.ToStrPtr("abbrev"),
		Ordinal:      1,
		Year:         util.ToStrPtr("1785"),
		Volume:       1,
	}}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/works/1", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	workRepo.EXPECT().SelectAll(gomock.Any()).Return(works, nil)
	// WHEN
	sut.GetWorks(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "id")
	assert.Contains(t, res.Body.String(), "code")
	assert.Contains(t, res.Body.String(), "abbreviation")
	assert.Contains(t, res.Body.String(), "ordinal")
	assert.Contains(t, res.Body.String(), "year")
	assert.Contains(t, res.Body.String(), "volume")
}

func assertErrorResponse(t *testing.T, res *httptest.ResponseRecorder, errStr string) {
	assert.Contains(t, res.Body.String(), "code")
	assert.Contains(t, res.Body.String(), "message")
	assert.Contains(t, res.Body.String(), errStr)
}
